/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0.
 */

#![allow(clippy::result_large_err)]

use aws_dynamodb_kms_demo::{
    make_config,
    scenario::add::{add_item, Item},
    scenario::error::Error,
    BaseOpt,
};
use aws_sdk_dynamodb::{error::DisplayErrorContext, Client};
use clap::Parser;
use std::process;

#[derive(Debug, Parser)]
struct Opt {
    /// account name for this account, used for identify wallet
    #[structopt(short, long)]
    name: String,

    /// kms key id which used for encryption on the wallet private key
    #[structopt(short, long)]
    key_id: String,

    /// encrypted private key of the wallet
    #[structopt(short, long)]
    encrypted_private_key: String,

    /// the address key of the wallet
    #[structopt(short, long)]
    address: String,

    /// the data key used to encrypt the private key
    #[structopt(short = 'd', long)]
    encrypted_data_key: String,

    /// The table name.
    #[structopt(short, long)]
    table: String,

    #[structopt(flatten)]
    base: BaseOpt,
}

/// Adds an item to an Amazon DynamoDB table.
/// The table schema must use one of username, p_type, age, first, or last as the primary key.
/// # Arguments
///
/// * `-t TABLE` - The name of the table.
/// * `-n name` -
/// * `-k key_id` -
/// * `-e encrypted_private_key` -
/// * `-a address` -
/// * `-d encrypted_data_key` -
/// * `[-r REGION]` - The region in which the table is created.
///   If not supplied, uses the value of the **AWS_REGION** environment variable.
///   If the environment variable is not set, defaults to **us-west-2**.
/// * `[-v]` - Whether to display additional information.
#[tokio::main]
async fn main() {
    if let Err(err) = run_example(Opt::parse()).await {
        eprintln!("Error: {}", DisplayErrorContext(err));
        process::exit(1);
    }
}

async fn run_example(
    Opt {
        name,
        key_id,
        encrypted_private_key,
        address,
        encrypted_data_key,
        table,
        base,
    }: Opt,
) -> Result<(), Error> {
    let shared_config = make_config(base).await?;
    let client = Client::new(&shared_config);

    add_item(
        &client,
        Item {
            name,
            key_id,
            encrypted_private_key,
            address,
            encrypted_data_key,
        },
        &table,
    )
    .await?;

    Ok(())
}
