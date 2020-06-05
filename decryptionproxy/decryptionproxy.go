/**
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package decryptionproxy

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/data"
	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/env"

	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/gcpkms"

	"github.com/sirupsen/logrus"
)

// GcsEndpoint specifies REST base endpoint for GCS
const GcsEndpoint = "storage.googleapis.com"

type handler struct {
	logger     *logrus.Logger
	config     env.Config
	restClient *http.Client
}

//ServeHTTP overwrites behavior to handle GET and performs decryption through Tink
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// err and wrapper variables are used to communicate errors from calling GCP APIs
	var err error
	resp := RespWrapper{
		Status:         http.StatusOK,
		ResponseWriter: w,
	}

	// At completion, provide details on bad response from REST calls to GCP
	defer func() {
		if err != nil {
			h.logger.WithFields(logrus.Fields{
				"response": resp.Status,
			}).WithError(err).Errorf("failed while decrypting %+v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), h.config.Proxy.Timeout)
	defer cancel()

	bucketEndpoint := h.config.BucketName + "." + GcsEndpoint
	url := url.URL{
		Scheme: "https",
		Host:   bucketEndpoint,
		Path:   r.URL.RequestURI(),
	}

	proxyToGCSReq, err := http.NewRequest(r.Method, url.String(), nil)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		resp.SaveStatus(http.StatusInternalServerError)
		return
	}

	proxyToGCSReq = proxyToGCSReq.WithContext(ctx)

	copyReqHeader(proxyToGCSReq.Header, r.Header)

	gcsResp, err := h.restClient.Do(proxyToGCSReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		resp.SaveStatus(http.StatusInternalServerError)
		return
	}
	defer gcsResp.Body.Close()

	copyRespHeader(resp, gcsResp)

	/// Decrypt response using Tink
	// from this point forward, unless otherwise specified, kill the proxy for any decryption erreor, since something
	// seriously wrong.
	// client will timeout for long operations based on environment variables
	keyURI := h.config.KmsMkekURI
	kmsClient, err := gcpkms.NewClient(keyURI)
	if err != nil {
		err := errors.Wrap(err, "gcp client creation failed")
		h.logger.Fatalf("%+v", err)
	}

	registry.RegisterKMSClient(kmsClient)

	bodyBytes, err := ioutil.ReadAll(gcsResp.Body)
	if err != nil {
		err := errors.Wrap(err, "cannot read ciphertext from GCS")
		h.logger.Fatalf("%+v", err)
	}

	var b data.EncryptedData
	if errUnmarshal := json.Unmarshal(bodyBytes, &b); errUnmarshal != nil {
		err = errors.Wrap(errUnmarshal, "either bad object name or unexpected structure from GCS")
		http.Error(resp, err.Error(), http.StatusBadRequest)
		resp.SaveStatus(http.StatusBadRequest)
		return
	}

	ee := data.NewEncryptionEngine(b.KekName, b.WdekName, kmsClient, h.logger)
	ee.Load(b)

	cipher, errDecode := base64.StdEncoding.DecodeString(b.EncryptedData)
	if errDecode != nil {
		err := errors.Wrap(errDecode, "possible incomplete GCS transfer.")
		h.logger.Fatalf("%+v", err)
	}

	plaintext := ee.Reveal(cipher)

	length, errRespWrite := resp.Write(plaintext)
	if errRespWrite != nil {
		err := errors.Wrap(errRespWrite, "could not write plaintext into client response")
		h.logger.Fatalf("%+v", err)
	}
	h.logger.Debugf("writer length: %v", length)

	// the size of the file in the bucket is different due to encryption, so when decrypted, its size doesn't match what
	// a client (i.e curl) might expect.  Avoids getting an error such as "(18) transfer closed with NN bytes remaining to read" where NN is the difference.
	resp.Header().Add("Content-Length", string(length))
}

// New returns a tink proxy handler
func New(c env.Config, client *http.Client) http.Handler {
	logger := c.Logger()
	return &handler{logger: logger, restClient: client, config: c}
}

// from httputil
func copyReqHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// copies GCS response headers to proxy to client
func copyRespHeader(dst RespWrapper, src *http.Response) {
	for k, vv := range src.Header {
		for _, v := range vv {
			// skip content length since encrypted and decrypted lengths will be different.
			// set it later after decryption
			if strings.Contains(k, "Content-Length") {
				continue
			}
			dst.Header().Add(k, v)
		}
	}
	dst.WriteHeader(src.StatusCode)
}
