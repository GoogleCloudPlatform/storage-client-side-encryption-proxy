# storage-client-side-encryption-proxy
This tool uses Tink to perform client side encryption operations backed against Google Cloud KMS.  It is both a client side encryption/decryption tool as well as a decrypting proxy for Google Cloud Storage.

**This is not an officially supported Google product**

## Pre-Setup
1.  You have golang 1.14 or newer installed

## Setup
1.  Create a service account that can encrypt and decrypt
2.  Use application credential: https://cloud.google.com/docs/authentication/production
3.  Setup a bucket to store your encrypted files.
4.  Setup KMS and create a key in the same region as your bucket
5.  acquire certificate and corresponding key for TLS. Place in the `tools` directory, and name them as follows.  Note: you can change the names and location by altering the `scriptes/variables.sh` file
    1.  `tools/cert.pem`
    2.  `tools/key.pem`
    3.  Note: for testing, consider creating a self signed cert: https://golang.org/src/crypto/tls/generate_cert.go
6.  edit `scripts/variables.sh` with your GCP information (i.e bucket name and key name)

## Building
1.  `go build -o tinkproxy`

## Running
This example uses the binary built named `tinkproxy` as described in the previous step. The tool uses Tink backed by Google Cloud KMS to encrypt a data encryption key (DEK) per directory, which is then
uploaded to your GCS bucket.  After the encrypted files are uploaded, a single file is then downlaod through the decrypting proxy, which decrypts using the appropriate KMS key.
1.  `./tinkproxy --help`
2.  `source scripts/variables.sh`  # be sure to edit the configuration to match your environment
3. `./scripts/uploadDirectory.sh samples`
4. `./scripts/getObject.sh`
5. `./scripts/cleanup.sh`

Also, you can encrypt and decrypt individual files
1. `./tinkproxy vanish samples/gettysburg.pdf -o demo.cipher`
2. `./tinkproxy reveal demo.cipher -o cleartext.pdf`

## Production Considerations
Consider the following items when using for production.
1. build and version the binary
2. deploy the proxy on trusted compute such as shielded VMs, use tmpfs, and private network access along with other controls to mitigate data exfiltration (e.g VPC-SC)
3. use certificates and AAD that meet your security governance requirements
4. use appropriate logging levels and client timeouts
