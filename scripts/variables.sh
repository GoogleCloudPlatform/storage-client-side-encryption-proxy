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

#===============================================================
# The following MUST be set to your project specific 
# resource names and environment (i.e TLS certificates)
#===============================================================

# Provide project name
export GOOGLE_CLOUD_PROJECT=""

# name of GCS bucket where encrypted files are stored
# do not include gs:// as part of the name
export TINKPROXY_BUCKET_NAME="" 

# include full resource name to GCP KMS
# example "gcp-kms://projects/<myproject>/locations/us/keyRings/<myKeyRing>/cryptoKeys/<myKey>" 
export TINKPROXY_KMS_MKEK_URI=""

# proxy only runs in TLS
export TINKPROXY_PROXY_CERT_FILE_PATH="tools/cert.pem"
export TINKPROXY_PROXY_CERT_KEY_FILE_PATH="tools/key.pem"

# AAD for AES encryption.  Pick a value that is meaningful and do not lose it. 
# data cannot be decrypted unless the same AAD value is supplied at encryption time.
export TINKPROXY_AAD="this aad"  

#===============================================================
# The following are configurable, but do not need to be changed.
#===============================================================

# path and file name of where to find DEK
export TINKPROXY_DEK_PATH_NAME="wdek.json"  

# store log files
export TINKPROXY_LOG_FILE="logs.out"  

#===============================================================
# Output for visual inspection
#===============================================================
echo ""
echo "Environment setup with values:"
echo "==================================================="
echo "Credentials Path: ${GOOGLE_APPLICATION_CREDENTIALS:?}"
echo "Project: ${GOOGLE_CLOUD_PROJECT:?}"
echo "Bucket Name: ${TINKPROXY_BUCKET_NAME:?}"
echo "KMS Master KEK URI: ${TINKPROXY_KMS_MKEK_URI:?}"
echo "DEK Pathname: ${TINKPROXY_DEK_PATH_NAME:?}"
echo "Log Level: ${TINKPROXY_LOG_LEVEL:DEBUG}"
echo "Log File: ${TINKPROXY_LOG_FILE:?}"
echo "Certificate File Path: ${TINKPROXY_PROXY_CERT_FILE_PATH:?}"
echo "Certificate Key File Path: ${TINKPROXY_PROXY_CERT_KEY_FILE_PATH:?}"
echo "==================================================="
echo ""
