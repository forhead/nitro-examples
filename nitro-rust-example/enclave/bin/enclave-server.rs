use clap::Parser;
use enclave::{
    chain_eth::generate_random_secret_key,
    kms::call_kms_generate_data_key,
    protocol_helper::{build_response, recv_message, send_message},
    crypto_utils::encrypt_by_data_key,
};
use serde_json::{Map, Value};
use vsock::{VsockAddr, VsockListener, VsockStream};

fn handle_client(mut stream: VsockStream) -> Result<(), anyhow::Error> {
    let payload_buffer = recv_message(&mut stream).map_err(|err| anyhow::anyhow!("{:?}", err))?;

    // Decode the payload as JSON
    let payload: Value =
        serde_json::from_slice(&payload_buffer).map_err(|err| anyhow::anyhow!("{:?}", err))?;

    if let Some(api_request) = payload["apiRequest"].as_str() {
        if api_request == "generateAccount" {
            let raw_data_key = call_kms_generate_data_key(
                payload["credential"].as_object().unwrap(),
                payload["key_id"].as_str().unwrap(),
            );

            let parts: Vec<&str> = raw_data_key.split("\n").collect();

            let data_key_ciphertext = parts[0].trim_start_matches("CIPHERTEXT: ");
            let data_key_plaintext = parts[1].trim_start_matches("PLAINTEXT: ");

            let (secret_key , ethereum_address) = generate_random_secret_key();
            let encrypted_private_key: String = encrypt_by_data_key(data_key_plaintext, &secret_key);

            let mut content: Map<String, Value> = Map::new();
            content.insert("encryptedPrivateKey".to_string(), Value::String(encrypted_private_key));
            content.insert("address".to_string(), Value::String(ethereum_address));
            content.insert("encryptedDataKey".to_string(), Value::String(data_key_ciphertext.to_string()));

            let response = build_response("generateResponse", content);

            send_message(&mut stream, response)?;
        }
    }

    Ok(())
}

#[derive(Debug, Parser)]
struct Opt {
    /// server virtio port
    #[structopt(short, long)]
    port: u32,
}

fn main() -> Result<(), anyhow::Error> {
    let Opt { port } = Opt::parse();

    let listener = VsockListener::bind(&VsockAddr::new(libc::VMADDR_CID_ANY, port))
        .expect("bind and listen failed");

    for stream in listener.incoming() {
        let stream = stream.unwrap();

        // write your own code here
        let _ = handle_client(stream);
    }

    Ok(())
}
