/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0.
 */

#![allow(clippy::result_large_err)]

use aes_gcm::{
    aead::{Aead, KeyInit},
    Aes256Gcm,
    Key, // Or `Aes128Gcm`
    Nonce,
};
use base64::{engine::general_purpose, Engine as _};
use clap::Parser;

#[derive(Debug, Parser)]
struct Opt {
    /// The encryption key.
    #[structopt(short, long)]
    data_key: String,

    /// The private key.
    #[structopt(short, long)]
    private_key: String,
}

// datakey_cipher_text:
// AQIDAHg7GwN+gKAgDZ0/L6q90F9t0vUeNYMZyeRjqvSQcWMlKwEkaR21ie8PGJMr7ruZBuAhAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMH1787Rqf8BLz4lKtAgEQgDu5NjaLlVYtoAdd2H+x/4UXsqYlUxe15vMhdO7U9anAphxZzwPoIXaYPGb7gdm4lrta76KM8FrrdSfPjA==
// datakey_plaintext:
// gnbLKIT08rBb6beAGBqyhBb+usZKBSQ3DgAyDFEolzs=

// encrypt decrypt by data key.
fn encrypt_by_data_key(datakey_plaintext_base64: &str, private_key: &str) {
    let datakey_bytes = general_purpose::STANDARD
        .decode(datakey_plaintext_base64)
        .expect("Input file does not contain valid base 64 characters.");

    // Create a key for AES256
    let key = Key::<Aes256Gcm>::from_slice(&datakey_bytes);

    // Create a new AES256 cipher
    let cipher = Aes256Gcm::new(key);

    // let nonce = Aes256Gcm::generate_nonce(&mut OsRng); // 96-bits; unique per message
    let nonce_bytes = [204, 92, 172, 44, 119, 145, 175, 178, 245, 248, 89, 193];
    let nonce = Nonce::from_slice(&nonce_bytes);

    let ciphertext = cipher.encrypt(nonce, private_key.as_bytes()).unwrap();
    let ciphertext_base64 = general_purpose::STANDARD.encode(ciphertext);
    dbg!(&ciphertext_base64);
}

/// Creates an AWS KMS data key.
/// # Arguments
///
/// * `[-d DataKey]` - The encryption key.
/// * `[-p Private_key]` - The private key.
///
fn main() {
    let Opt {
        data_key,
        private_key,
    } = Opt::parse();

    encrypt_by_data_key(&data_key, &private_key);
}
