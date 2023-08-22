/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0.
 */

#![allow(clippy::result_large_err)]

use aws_dynamodb_kms_demo::{
    make_config, scenario::error::Error, scenario::query::query_item, BaseOpt,
};
use aws_sdk_dynamodb::{error::DisplayErrorContext, Client};
use clap::Parser;
use std::process;

#[derive(Debug, Parser)]
struct Opt {
    /// account name for this account, used for identify wallet
    #[structopt(short, long)]
    name: String,

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

async fn run_example(Opt { name, table, base }: Opt) -> Result<(), Error> {
    let shared_config = make_config(base).await?;
    let client = Client::new(&shared_config);

    query_item(&client, &name, &table).await?;

    Ok(())
}
