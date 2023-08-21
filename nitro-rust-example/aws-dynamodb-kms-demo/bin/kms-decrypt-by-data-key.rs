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

// encrypt decrypt by data key.
fn decrypt_by_data_key(datakey_plaintext_base64: &str, private_key_base64: &str) {
    let datakey_bytes = general_purpose::STANDARD
        .decode(datakey_plaintext_base64)
        .expect("Input file does not contain valid base 64 characters.");

    let private_key_bytes = general_purpose::STANDARD
        .decode(private_key_base64)
        .expect("Input file does not contain valid base 64 characters.");

    // Create a key for AES256
    let key = Key::<Aes256Gcm>::from_slice(&datakey_bytes);

    // Create a new AES256 cipher
    let cipher = Aes256Gcm::new(key);

    // hardcode nonce
    let nonce_bytes = [204, 92, 172, 44, 119, 145, 175, 178, 245, 248, 89, 193];
    let nonce = Nonce::from_slice(&nonce_bytes);

    let private_key_origin = cipher.decrypt(nonce, private_key_bytes.as_ref()).unwrap();
    let private_key_origin_str = std::str::from_utf8(&private_key_origin).unwrap();
    dbg!(&private_key_origin_str);
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

    decrypt_by_data_key(&data_key, &private_key);
}
