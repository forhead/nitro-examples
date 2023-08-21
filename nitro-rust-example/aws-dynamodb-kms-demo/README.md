Provides a list of commands to interact with AWS DynamoDB, KMS using the Rust.

### Commands:

#### Create a DynamoDB table:
```sh
cargo run --bin aws-dynamodb-create-table -- -r ap-east-1 -t AccountTable -k name
```

#### Add an item to the DynamoDB table:
```sh
cargo run --bin aws-dynamodb-add-item -- \
  -r ap-east-1 \
  -t AccountTable \
  -n UID10000 \
  -k 0a5713c9-7d29-4b8e-aa0e-8e27e58ac6e1 \
  -e encrypted_private_key \
  -a 0x1f9090aae28b8a3dceadf281b0f12828e676c326 \
  -d encrypted_private_key
```

#### Get an item from the DynamoDB table:
```sh
cargo run --bin aws-dynamodb-get-item -- -r ap-east-1 -t AccountTable -n UID10000 
cargo run --bin aws-dynamodb-get-item -- -r ap-east-1 -t AccountTable -n UID10001
```

#### Create a KMS key:
```sh
cargo run --bin aws-kms-create-key -- -r ap-east-1
```

#### Generate a data key using the KMS key:
```sh
cargo run --bin aws-kms-generate-data-key -- -r ap-east-1 -k 0a5713c9-7d29-4b8e-aa0e-8e27e58ac6e1
```

#### Encrypt a private key using the KMS Data key:
```sh
cargo run --bin aws-kms-encrypt-by-data-key -- \
  -d gnbLKIT08rBb6beAGBqyhBb+usZKBSQ3DgAyDFEolzs= \
  -p 0x3a1076bf45ab87712ad64ccb3b10217737f7faacbf2872e88fdd9a537d8fe266
```

#### Decrypt a private key using the KMS Data key:
```sh
cargo run --bin aws-kms-decrypt-by-data-key -- \
  -d gnbLKIT08rBb6beAGBqyhBb+usZKBSQ3DgAyDFEolzs= \
  -p K20TvEKnpVnfHIcRDG4i+8hT6VoquHlZH7PuDtzxqvazlsYcZFkOj9eUpi5/kMo4LyK95Fdv6wfcwki0Fw1pl/2Z22He41dZCaZxhX98NvBiXA==
```

#### Decrypt KMS Data key:
```sh
cargo run --bin aws-kms-decrypt-data-key -- \
  -k 0a5713c9-7d29-4b8e-aa0e-8e27e58ac6e1 \
  -d AQIDAHg7GwN+gKAgDZ0/L6q90F9t0vUeNYMZyeRjqvSQcWMlKwE2iKZnUfZ5oCWW2rSqQZFjAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQM/y0MA93t1BITErGiAgEQgDuKtv5ZuVbeTQ14QHsR5f82IpD8zhQdejMNc8+0FZUByt8f4mw/tO/+KZEiBOdl09YfCSjndlrFh6XaTw==
```

#### (deprecated) Encrypt data using the KMS key:
```sh
cargo run --bin aws-kms-encrypt -- -r ap-east-1 -k 0a5713c9-7d29-4b8e-aa0e-8e27e58ac6e1 -o /tmp/kms-encrypt.txt -t KeyId1
```

#### (deprecated) Decrypt data using the KMS key:
```sh
cargo run --bin aws-kms-decrypt -- -r ap-east-1 -k 0a5713c9-7d29-4b8e-aa0e-8e27e58ac6e1 -i /tmp/kms-encrypt.txt
```

