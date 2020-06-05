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
	"encoding/base64"
	"reflect"
	"testing"
)

func TestNewEncryptedData(t *testing.T) {
	kekName := "fake/resource/keyring/key"
	wdekName := "wdek.json"
	wdek := "12345678"
	data := []byte{1, 2, 3, 4}
	type args struct {
		kekName  string
		wdekName string
		wdek     string
		data     []byte
	}
	tests := []struct {
		name string
		args args
		want EncryptedData
	}{
		{
			name: "basic test",
			args: args{
				kekName:  kekName,
				wdekName: wdekName,
				wdek:     wdek,
				data:     data,
			},
			want: EncryptedData{
				KekName:       kekName,
				WdekName:      wdekName,
				Wdek:          wdek,
				EncryptedData: base64.StdEncoding.EncodeToString(data),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEncryptedData(tt.args.kekName, tt.args.wdekName, tt.args.wdek, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEncryptedData() = %v, want %v", got, tt.want)
			}
		})
	}
}
