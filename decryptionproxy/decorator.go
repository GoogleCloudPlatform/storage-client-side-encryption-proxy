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
)

// Decorator allows creation of uniform middleware for the proxy
type Decorator func(http.Handler) http.Handler

// Decorate executes each middleware in the order it's specified
func Decorate(handler http.Handler, decorators ...Decorator) http.Handler {
	for i := len(decorators) - 1; i >= 0; i-- {
		handler = decorators[i](handler)
	}
	return handler
}
