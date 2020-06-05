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
	"net/http/httptest"
	"testing"
)

func Test_respWrapper_WriteHeader(t *testing.T) {
	type fields struct {
		status         int
		ResponseWriter http.ResponseWriter
	}
	type args struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "simple wrapper test",
			fields: fields{
				status:         http.StatusOK,
				ResponseWriter: httptest.NewRecorder(),
			},
			args: args{
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &RespWrapper{
				Status:         tt.fields.status,
				ResponseWriter: tt.fields.ResponseWriter,
			}
			w.WriteHeader(tt.args.status)
		})
	}
}
