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

import "time"

// ProxyConfig contains configuration for the proxy mode.
type ProxyConfig struct {
	Listen          string        `default:":8080"`
	Timeout         time.Duration `split_words:"true" default:"10s"`
	CertFilePath    string        `split_words:"true" required:"true"`
	CertKeyFilePath string        `split_words:"true" required:"true"`
}
