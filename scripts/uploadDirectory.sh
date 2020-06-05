#!/usr/bin/env bash
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

DIR=$1

echo ""
echo "==================================================="
echo "Encrypting client side - all files in directory: ${DIR:?}"
read -p "press enter to continue"

./tinkproxy vanish ${DIR}
echo "==================================================="
echo ""
read -p "press enter to continue"

echo "==================================================="
echo "Uploading all files into bucket: ${TINKPROXY_BUCKET_NAME:?}"
read -p "press enter to continue" 
# TINKPROXY encrypts files and encode it changing suffix to .enc
gsutil -m cp ${DIR}/*.enc gs://${TINKPROXY_BUCKET_NAME}
echo "==================================================="
read -p "press enter to continue"

echo ""
echo "Success"
echo ""
