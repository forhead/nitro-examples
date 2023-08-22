use clap::Parser;
use parent::{
    dynamodb_helper::{make_config, BaseOpt},
    iam::get_iam_token,
    protocol_helper::{build_payload, recv_message, send_message},
    scenario::add::{add_item, Item},
};

use aws_sdk_dynamodb::Client;
use serde_json::{Map, Value};
use vsock::{VsockAddr, VsockStream};

#[derive(Debug, Parser)]
struct Opt {
    /// CID
    #[structopt(short, long)]
    cid: u32,

    /// The port
    #[structopt(short, long)]
    port: u32,

    /// KMS key id
    #[structopt(short, long)]
    key_id: String,

    /// DynamoDB table name.
    #[structopt(short, long)]
    table: String,
}

#[tokio::main]
async fn main() -> Result<(), anyhow::Error> {
    let Opt {
        cid,
        port,
        key_id,
        table,
    } = Opt::parse();

    // Initiate a connection on an AF_VSOCK socket
    let mut stream = VsockStream::connect(&VsockAddr::new(cid, port)).expect("connection failed");

    // build payload
    let credential = get_iam_token().await.unwrap();

    let base_opt = BaseOpt {
        region: Some("ap-east-1".to_string()),
        verbose: false,
    };
    let shared_config = make_config(base_opt).await?;
    let dynamodb_client = Client::new(&shared_config);

    let credential: Map<String, Value> = credential
        .into_iter()
        .map(|(k, v)| (k, Value::String(v)))
        .collect();
    let user_id = "UID10000".to_string();
    let payload = build_payload("generateAccount", credential, user_id.clone(), key_id.clone());

    // send payload
    send_message(&mut stream, payload)?;

    // recv response
    let response = recv_message(&mut stream).map_err(|err| anyhow::anyhow!("{:?}", err))?;

    // Decode the payload as JSON
    let json: Value =
        serde_json::from_slice(&response).map_err(|err| anyhow::anyhow!("{:?}", err))?;
    println!("response {}", json);

    let content = json["content"].as_object().unwrap();

    // write to dynamodb
    add_item(
        &dynamodb_client,
        Item {
            name: user_id,
            key_id,
            encrypted_private_key: content["encryptedPrivateKey"].as_str().unwrap().to_string(),
            address: content["address"].as_str().unwrap().to_string(),
            encrypted_data_key: content["encryptedDataKey"].as_str().unwrap().to_string(),
        },
        &table,
    )
    .await?;
    Ok(())
}
