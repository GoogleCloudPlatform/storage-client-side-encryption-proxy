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

package env

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	ghttp "google.golang.org/api/transport/http"
)

// ClientConfig for Google Cloud Storage
type ClientConfig struct {
	Timeout         time.Duration `split_words:"true" default:"3s"`
	IdleConnTimeout time.Duration `split_words:"true" default:"60s"`
	MaxIdleConns    int           `split_words:"true" default:"30"`
}

// BasicTLSClient sets up TLS, default application credentials, and timeouts.
func (c ClientConfig) BasicTLSClient() (*http.Client, error) {
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	t := http.Transport{
		IdleConnTimeout: c.IdleConnTimeout,
		MaxIdleConns:    c.MaxIdleConns,
		TLSClientConfig: cfg,
	}
	gTransport, err := ghttp.NewTransport(context.Background(), &t, option.WithScopes(storage.ScopeReadOnly))

	return &http.Client{
		Timeout:   c.Timeout,
		Transport: gTransport,
	}, err
}
