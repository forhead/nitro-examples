use aws_dynamodb_kms_demo::{
    make_config, scenario::create::create_table, scenario::error::Error, BaseOpt,
};
use {
    aws_sdk_dynamodb::{error::DisplayErrorContext, Client},
    clap::Parser,
    std::process,
};

#[derive(Debug, Parser)]
struct Opt {
    /// The table name
    #[structopt(short, long)]
    table: String,

    /// The primary key
    #[structopt(short, long)]
    key: String,

    #[structopt(flatten)]
    base: BaseOpt,
}

/// Creates a DynamoDB table.
/// # Arguments
///
/// * `-k KEY` - The primary key for the table.
/// * `-t TABLE` - The name of the table.
/// * `[-r DEFAULT-REGION]` - The region in which the client is created.
///    If not supplied, uses the value of the **AWS_DEFAULT_REGION** environment variable.
///    If the environment variable is not set, defaults to **us-west-2**.
/// * `[-v]` - Whether to display additional information.
#[tokio::main]
async fn main() {
    if let Err(err) = run_example(Opt::parse()).await {
        eprintln!("Error: {}", DisplayErrorContext(err));
        process::exit(1);
    }
}

async fn run_example(Opt { table, key, base }: Opt) -> Result<(), Error> {
    let shared_config = make_config(base).await?;
    let client = Client::new(&shared_config);

    create_table(&client, &table, &key).await?;

    Ok(())
}
