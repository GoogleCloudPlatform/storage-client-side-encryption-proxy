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
)

// EncryptedData is the object stored in the bucket.
//    EncryptedData is base64 encoded for transfer.
//    Wdek from Tink is json and encrypted.
//    WdekName is the primaryKeyID for the dek
//    KekName is the key stored in GCP KMS
type EncryptedData struct {
	KekName       string `json:"kek"`
	WdekName      string `json:"wdekName"`
	Wdek          string `json:"wdek"`
	EncryptedData string `json:"data"`
}

// NewEncryptedData constructs an object to send to GCS
func NewEncryptedData(kekName string, wdekName string, wdek string, data []byte) EncryptedData {
	return EncryptedData{
		KekName:       kekName,
		WdekName:      wdekName,
		Wdek:          wdek,
		EncryptedData: base64.StdEncoding.EncodeToString(data),
	}
}
