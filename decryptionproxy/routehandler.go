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
	"strings"
)

// RouteHandler middleware does basic path checks.  Enforces a root path
func RouteHandler() Decorator {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// handle route construction.
			if r.URL.Path != "" {
				if !strings.HasPrefix(r.URL.Path, "/") {
					r.URL.Path = "/" + r.URL.Path
				}
				handler.ServeHTTP(w, r)
			} else {
				http.Error(w, "not found", http.StatusNotFound)
			}
		})
	}
}
