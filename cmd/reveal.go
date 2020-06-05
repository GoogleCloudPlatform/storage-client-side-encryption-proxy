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

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/data"
	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/env"

	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// revealCmd represents the reveal command
var revealCmd = &cobra.Command{
	Use:   "reveal",
	Short: "Make encrypted data appear",
	Long: `Using a Tink enabled KMS backend, decrypt the data. If the outputFile is provided,
	the plaintext is saved there.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reveal called")

		config, err := env.Get()
		if err != nil {
			log.Fatal(err)
		}
		logger := config.Logger()

		keyURI := config.KmsMkekURI
		gcpclient, err := gcpkms.NewClient(keyURI)
		if err != nil {
			err := errors.Wrap(err, "gcp client creation failed")
			logger.Fatalf("%+v", err)
		}

		registry.RegisterKMSClient(gcpclient)

		f, errRead := ioutil.ReadFile(args[0])
		if errRead != nil {
			errRead := errors.Wrap(errRead, "check file")
			logger.Fatalf("%+v", errRead)
		}
		var b data.EncryptedData
		if errUnmarshal := json.Unmarshal(f, &b); errUnmarshal != nil {
			errUnmarshal := errors.Wrap(errUnmarshal, "cannot unmarshal")
			logger.Fatalf("%+v", errUnmarshal)
		}

		ee := data.NewEncryptionEngine(b.KekName, b.WdekName, gcpclient, logger)
		ee.Load(b)

		cipher, errDecode := base64.StdEncoding.DecodeString(b.EncryptedData)
		if errDecode != nil {
			errDecode := errors.Wrap(errDecode, "unexpected wdek format. Use Tink")
			logger.Fatalf("%+v", errDecode)
		}
		ciphertext := ee.Reveal(cipher)

		if outputFile != "" {
			if err := ioutil.WriteFile(outputFile, ciphertext, 0644); err != nil {
				err := errors.Wrap(err, "check file or disk space")
				logger.Fatalf("%+v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(revealCmd)
}
