/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0.
 */

#![allow(clippy::result_large_err)]

use aws_config::meta::region::RegionProviderChain;
use aws_sdk_kms::primitives::Blob;
use aws_sdk_kms::{config::Region, meta::PKG_VERSION, Client, Error};
use base64::{engine::general_purpose, Engine as _};
use clap::Parser;

#[derive(Debug, Parser)]
struct Opt {
    /// The AWS Region.
    #[structopt(short, long)]
    region: Option<String>,

    /// The encryption key.
    #[structopt(short, long)]
    key: String,

    /// The name of the input file with encrypted text to decrypt.
    #[structopt(short, long)]
    data_key: String,

    /// Whether to display additional information.
    #[structopt(short, long)]
    verbose: bool,
}

// Decrypt a string.
// snippet-start:[kms.rust.decrypt]
async fn decrypt_key(client: &Client, key: &str, data_key: &str) -> Result<(), Error> {
    // Open input text file and get contents as a string
    // input is a base-64 encoded string, so decode it:
    let data = general_purpose::STANDARD
        .decode(data_key)
        .map(Blob::new)
        .expect("Input file does not contain valid base 64 characters.");

    let resp = client
        .decrypt()
        .key_id(key)
        .ciphertext_blob(data)
        .send()
        .await?;

    let inner = resp.plaintext.unwrap();
    let bytes = inner.as_ref();

    // The 'bytes' variable should be equivalent to
    // the result of decoding the base64 string, e.g: "l3p994w+hdqFzKwA3zuii6Lb9DsIfcpLfcAHl11goTY="
    dbg!(&bytes);

    Ok(())
}
// snippet-end:[kms.rust.decrypt]

/// Decrypts a string encrypted by AWS KMS.
/// # Arguments
///
/// * `-k KEY` - The encryption key.
/// * `-d Encrypted Data KEY` -
/// * `[-r REGION]` - The Region in which the client is created.
///    If not supplied, uses the value of the **AWS_REGION** environment variable.
///    If the environment variable is not set, defaults to **us-west-2**.
/// * `[-v]` - Whether to display additional information.
#[tokio::main]
async fn main() -> Result<(), Error> {
    let Opt {
        key,
        data_key,
        region,
        verbose,
    } = Opt::parse();

    let region_provider = RegionProviderChain::first_try(region.map(Region::new))
        .or_default_provider()
        .or_else(Region::new("us-west-2"));
    println!();

    if verbose {
        println!("KMS client version: {}", PKG_VERSION);
        println!(
            "Region:             {}",
            region_provider.region().await.unwrap().as_ref()
        );
        println!("Key:                {}", &key);
        println!("Input:              {}", &data_key);
        println!();
    }

    let shared_config = aws_config::from_env().region(region_provider).load().await;
    let client = Client::new(&shared_config);

    decrypt_key(&client, &key, &data_key).await
}
