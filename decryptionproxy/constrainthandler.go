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
	"net/http"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

// ConstraintHandler middleware enforces limitations that the proxy currently has
// 1. proxy only supports GET
// 2. must have at least a bucket name
func ConstraintHandler(logger *logrus.Logger) Decorator {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			///<1> supported methods
			if r.Method != http.MethodGet &&
				r.Method != http.MethodHead {
				err := errors.New("method not yet supported")
				logger.Errorf("%s: %+v", r.Method, err)
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
				return
			}

			///<2> valid objects
			if r.URL.Path == "/" {
				err := errors.New("must specify a valid object, not root directory")
				logger.Errorf("%+v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
