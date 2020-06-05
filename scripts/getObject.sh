#/usr/bin/env bash
# 
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ENCRYPED_OBJ="gettysburg.pdf.enc"
OUTPUT_FILENAME="plaintext.pdf"

echo ""
echo "==================================================="
echo "starting proxy to GCS: ${TINKPROXY_BUCKET_NAME:?}"
./tinkproxy proxy &
echo "==================================================="
read -p "press enter to continue"

echo ""
echo "==================================================="
echo "downloading file: ${ENCRYPED_OBJ} through proxy from ${TINKPROXY_BUCKET_NAME}"
read -p "press enter to continue"
echo ""
curl -k https://localhost:8080/${ENCRYPED_OBJ} -o ${OUTPUT_FILENAME}
echo ""
echo ""
echo "retrieved ${ENCRYPED_OBJ} as ${OUTPUT_FILENAME} and decrypted it"
echo "==================================================="
read -p "press enter to continue"

echo ""
echo "==================================================="
echo "comparing the decrypted file with original file"
diff -s  ${OUTPUT_FILENAME} samples/gettysburg.pdf
echo "==================================================="
read -p "press enter to continue"

echo ""
echo "==================================================="
echo "stopping proxy"
kill %1
echo "==================================================="

echo ""
echo "Success"
echo ""
