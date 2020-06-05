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
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/decryptionproxy"
	"github.com/GoogleCloudPlatform/storage-client-side-encryption-proxy/env"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts a decrypting proxy server to GCS",
	Long: `Supports a GET request to retrieve encrypted files using Tink to decrypt it.  
	Defaults to localhost:8080 unless otherwise specified by environment variables`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")
		config, err := env.Get()
		if err != nil {
			log.Fatal(err)
		}
		logger := config.Logger()

		hc, errClient := config.Client.BasicTLSClient()
		if err != nil {
			errClient := errors.Wrap(errClient, "cannot init https client")
			logger.Fatalf("%+v", errClient)
		}

		tinkProxyHandler := getHandler(config, hc)

		middlewareHandlers := decryptionproxy.Decorate(tinkProxyHandler, decryptionproxy.LoggerHandler(logger),
			decryptionproxy.RouteHandler(),
			decryptionproxy.ConstraintHandler(logger),
		)

		s := &http.Server{
			Addr:           config.Proxy.Listen,
			Handler:        middlewareHandlers,
			ReadTimeout:    2 * time.Second,           // handle slow clients
			WriteTimeout:   2 * config.Client.Timeout, // give room for GCS request to complete
			MaxHeaderBytes: 1 << 20,
		}

		logger.Infof("Starting proxy %s", config.Proxy.Listen)
		logger.Fatal(s.ListenAndServeTLS(config.Proxy.CertFilePath, config.Proxy.CertKeyFilePath))
	},
}

func getHandler(c env.Config, hc *http.Client) http.HandlerFunc {
	proxyHandler := decryptionproxy.New(c, hc)

	return func(w http.ResponseWriter, r *http.Request) {
		proxyHandler.ServeHTTP(w, r)
	}
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
