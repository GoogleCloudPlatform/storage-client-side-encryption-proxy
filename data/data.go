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
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/keyset"
	"github.com/sirupsen/logrus"
)

// EncryptionEngine specifies necessary details to use Tink.
type EncryptionEngine struct {
	kekName      string
	wDekPathName string // path name to file containing wdek
	aad          string
	dekHandle    *keyset.Handle // dek in memory (in clear)
	gcpClient    registry.KMSClient
	logger       *logrus.Logger
}

// Encryptor defines methods to support data encryption
type Encryptor interface {
	Obfuscate(data io.Reader)
	Reveal() io.Writer

	ReadWdek()
	WriteWdek()

	Package(data []byte) EncryptedData
	Load(data EncryptedData)
}

//NewEncryptionEngine creates engines with required parameters
func NewEncryptionEngine(kekName string, wDekPathName string, gcpClient registry.KMSClient, logger *logrus.Logger) *EncryptionEngine {
	return &EncryptionEngine{
		kekName:      kekName,
		wDekPathName: wDekPathName,
		gcpClient:    gcpClient,
		logger:       logger,
	}
}

// ReadWdek loads the wdek using KMS
func (ee *EncryptionEngine) ReadWdek() {
	fo, err := os.Open(ee.wDekPathName)
	if err != nil {
		err := errors.Wrapf(err, "cannot read wdek to unmarshal it %s", ee.wDekPathName)
		ee.logger.Fatalf("%+v", err)
	}
	defer fo.Close()

	backend, errClient := ee.gcpClient.GetAEAD(ee.kekName)
	if errClient != nil {
		errClient := errors.Wrap(errClient, "cannot retrieve dek from KMS")
		ee.logger.Fatalf("%+v", errClient)
	}

	jreader := keyset.NewJSONReader(fo)
	masterKey := aead.NewKMSEnvelopeAEAD(*aead.AES256GCMKeyTemplate(), backend)

	// Read the encrypted keyset handle back from the io.Reader implementation
	// and decrypt it using the master key.
	ee.dekHandle, err = keyset.Read(jreader, masterKey)
	if err != nil {
		err := errors.Wrap(err, "cannot created wdek handle")
		ee.logger.Fatalf("%+v", err)
	}
}

// WriteWdek outputs JSON file with wDEK (encrypted)
func (ee *EncryptionEngine) WriteWdek() {
	if ee.dekHandle == nil {
		ee.logger.Info("creating a new DEK, since none found")
		dek := aead.AES256GCMKeyTemplate()
		kh, err := keyset.NewHandle(dek)
		if err != nil {
			err := errors.Wrap(err, "creating dek key handle failed")
			ee.logger.Fatalf("%+v", err)
		}
		ee.dekHandle = kh
	}

	f, errCreate := os.Create(ee.wDekPathName)
	if errCreate != nil {
		errCreate := errors.Wrapf(errCreate, "writing wdek failed %s", ee.wDekPathName)
		ee.logger.Fatalf("%+v", errCreate)
	}
	defer f.Close()

	jwriter := keyset.NewJSONWriter(f)

	backend, errClient := ee.gcpClient.GetAEAD(ee.kekName)
	if errClient != nil {
		errClient := errors.Wrap(errClient, "cannot retrieve dek")
		ee.logger.Fatalf("%+v", errClient)
	}

	wdek := aead.NewKMSEnvelopeAEAD(*aead.AES256GCMKeyTemplate(), backend)

	// Write encrypts the keyset handle with the master key and writes to the
	// io.Writer implementation (memKeyset).
	if err := ee.dekHandle.Write(jwriter, wdek); err != nil {
		err := errors.Wrap(err, "cannot write JSON marshalled wdek")
		ee.logger.Fatalf("%+v", err)
	}
}

// Obfuscate encrypts data using the underlying encryption engine
func (ee *EncryptionEngine) Obfuscate(dataPlain []byte) []byte {
	ee.logger.Infof("...encrypting using this master KEK %s\n", ee.kekName)

	if ee.dekHandle == nil {
		ee.ReadWdek()
	}

	a, err := aead.New(ee.dekHandle)
	if err != nil {
		err := errors.Wrap(err, "AEAD encryption object creation failed")
		ee.logger.Fatalf("%+v", err)
	}

	ct, err := a.Encrypt(dataPlain, []byte(ee.aad))
	if err != nil {
		err := errors.Wrap(err, "cannot encrypt")
		ee.logger.Fatalf("%+v", err)
	}

	return ct
}

// Reveal decrypts data using the underlying encryption engine
func (ee *EncryptionEngine) Reveal(cipherData []byte) []byte {
	ee.logger.Infof("...decrypting with this master KEK %s\n", ee.kekName)

	if ee.dekHandle == nil {
		ee.ReadWdek()
	}

	a, err := aead.New(ee.dekHandle)
	if err != nil {
		err := errors.Wrap(err, "cannot create AEAD object from wdek")
		ee.logger.Fatalf("%+v", err)
	}

	pt, err := a.Decrypt(cipherData, []byte(ee.aad))
	if err != nil {
		err := errors.Wrap(err, "cannot decrypt data")
		ee.logger.Fatalf("%+v", err)
	}

	return pt
}

// Load grabs the wDek and reads it in.
func (ee *EncryptionEngine) Load(data EncryptedData) {
	if err := ioutil.WriteFile(ee.wDekPathName, []byte(data.Wdek), 0644); err != nil {
		err := errors.Wrapf(err, "cannot create file for loading: %s", ee.wDekPathName)
		ee.logger.Fatalf("%+v", err)
	}

	ee.ReadWdek()
}

// Package marshalls the encrypted data with key hierarchy information to be stored as a blob of structured data
func (ee *EncryptionEngine) Package(data []byte) EncryptedData {
	ee.WriteWdek()
	wdek, err := ioutil.ReadFile(ee.wDekPathName)
	if err != nil {
		err := errors.Wrap(err, "cannot open wdek")
		ee.logger.Fatalf("%+v", err)
	}

	return NewEncryptedData(ee.kekName, ee.wDekPathName, string(wdek), data)
}
