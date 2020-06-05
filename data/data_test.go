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

package data

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/env"

	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/google/tink/go/keyset"
	"github.com/pkg/errors"
)

func TestEncryptionEngine_WdekOps(t *testing.T) {
	config, err := env.Get()
	if err != nil {
		log.Fatal(err)
	}

	keyURI := config.KmsMkekURI
	gcpclient, err := gcpkms.NewClient(keyURI)
	if err != nil {
		err := errors.Wrap(err, "gcp client creation failed")
		log.Fatal(err)
	}

	registry.RegisterKMSClient(gcpclient)

	type fields struct {
		wDekPathName string
		dekHandle    *keyset.Handle
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"simpleTest",
			fields{
				"wdek2.json",
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EncryptionEngine{
				kekName:      keyURI,
				wDekPathName: tt.fields.wDekPathName,
				dekHandle:    tt.fields.dekHandle,
				gcpClient:    gcpclient,
				logger:       config.Logger(),
			}
			ee.WriteWdek()
			ee.ReadWdek()
		})
	}
}

func TestEncryptionEngine_EncryptDecrypt(t *testing.T) {
	config, err := env.Get()
	if err != nil {
		log.Fatal(err)
	}

	keyURI := config.KmsMkekURI

	gcpclient, err := gcpkms.NewClient(keyURI)
	if err != nil {
		err := errors.Wrap(err, "gcp client creation failed")
		log.Fatal(err)
	}

	registry.RegisterKMSClient(gcpclient)

	gettysburgFile, err := ioutil.ReadFile("../samples/gettysburg.pdf")
	if err != nil {
		log.Fatal(err)
	}

	type fields struct {
		wDekPathName string
		dekHandle    *keyset.Handle
	}
	type args struct {
		dataPlain []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"simpleTest",
			fields{
				"wdek2.json",
				nil,
			},
			args{
				[]byte("encrypt this"),
			},
		},
		{
			"FileTest",
			fields{
				"wdek2.json",
				nil,
			},
			args{
				gettysburgFile,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EncryptionEngine{
				kekName:      keyURI,
				wDekPathName: tt.fields.wDekPathName,
				dekHandle:    tt.fields.dekHandle,
				gcpClient:    gcpclient,
				logger:       config.Logger(),
			}
			got := ee.Obfuscate(tt.args.dataPlain)
			plaintext := ee.Reveal(got)
			if !reflect.DeepEqual(plaintext, tt.args.dataPlain) {
				t.Errorf("EncryptionEngine.Reveal() = %v, want %v", plaintext, tt.args.dataPlain)
			}
		})
	}
}
