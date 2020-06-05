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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/data"
	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/env"

	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// vanishCmd represents the vanish command
var vanishCmd = &cobra.Command{
	Use:   "vanish",
	Short: "make data vanish by encrypting the entire file or directory",
	Long: `Using a Tink enabled KMS backend, encrypt the data. If the outputFile is provided,
the ciphertext saved there.  outputFile is ignored, if  a directory is being encrypted.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vanish called")

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

		wDekPathName := config.DekPathName
		registry.RegisterKMSClient(gcpclient)

		ee := data.NewEncryptionEngine(keyURI, wDekPathName, gcpclient, logger)

		// support both a file or directory
		sourceItem := args[0]
		fi, err := os.Stat(sourceItem)
		if err != nil {
			err := errors.Wrap(err, "check the specified file/directory.")
			logger.Fatalf("%+v", err)
		}
		switch mode := fi.Mode(); {
		case mode.IsRegular():
			handleFile(sourceItem, ee, logger)
		case mode.IsDir():
			handleDir(sourceItem, ee, logger)
		case mode&os.ModeSymlink != 0:
			logger.Fatalf("%+v", errors.New("cannot handle symLink"))
		case mode&os.ModeNamedPipe != 0:
			logger.Fatalf("%+v", errors.New("cannot handle pipes"))
		}
	},
}

func handleDir(dir string, ee *data.EncryptionEngine, logger *logrus.Logger) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		err := errors.Wrap(err, "cannot read directory")
		logger.Fatalf("%+v", err)
	}

	if outputFile != "" {
		logger.Info("Ignoring output file name. Processing a directory.")
	}

	// Create a wDEK
	ee.WriteWdek()

	for _, file := range files {
		// Skip encrypted files
		if strings.Contains(file.Name(), "enc") {
			continue
		}

		fmt.Println(file.Name())
		f, err := ioutil.ReadFile(dir + "/" + file.Name())
		if err != nil {
			err := errors.Wrap(err, "cannot read file")
			logger.Fatalf("%+v", err)
		}

		ciphertext := ee.Obfuscate(f)
		encryptedBlob := ee.Package(ciphertext)
		b, errMarshal := json.Marshal(encryptedBlob)
		if errMarshal != nil {
			errMarshal := errors.Wrap(errMarshal, "marshal encrypted data to package for writing")
			logger.Fatalf("%+v", errMarshal)
		}

		if errWrite := ioutil.WriteFile(dir+"/"+file.Name()+".enc", b, 0644); err != nil {
			errWrite := errors.Wrap(errWrite, "cannot write wdek")
			logger.Fatalf("%+v", errWrite)
		}
	}
}

func handleFile(file string, ee *data.EncryptionEngine, logger *logrus.Logger) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	ee.WriteWdek()
	ciphertext := ee.Obfuscate(f)

	encryptedBlob := ee.Package(ciphertext)
	b, errMarshal := json.Marshal(encryptedBlob)
	if errMarshal != nil {
		errMarshal := errors.Wrap(errMarshal, "marshal encrypted data to package for writing")
		logger.Fatalf("%+v", errMarshal)
	}
	if outputFile != "" {
		if err := ioutil.WriteFile(outputFile, b, 0644); err != nil {
			err := errors.Wrap(err, "check specified output")
			logger.Fatalf("%+v", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(vanishCmd)
}
