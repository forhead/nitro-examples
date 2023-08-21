use serde_json::{Map, Value};
use std::env;
use subprocess::{Popen, PopenConfig, Redirection};

pub fn call_kms_generate_data_key(credential: &Map<String, Value>, key_id: &str) -> String {
    let aws_access_key_id = credential["aws_access_key_id"].as_str().unwrap();
    let aws_secret_access_key = credential["aws_secret_access_key"].as_str().unwrap();
    let aws_session_token = credential["aws_session_token"].as_str().unwrap();

    let mut p = Popen::create(
        &[
            "/app/kmstool_enclave_cli",
            "genkey",
            "--region",
            &env::var("REGION").unwrap(),
            "--proxy-port",
            "8000",
            "--aws-access-key-id",
            aws_access_key_id,
            "--aws-secret-access-key",
            aws_secret_access_key,
            "--aws-session-token",
            aws_session_token,
            "--key-id",
            key_id,
            "--key-spec",
            "AES-256",
        ],
        PopenConfig {
            stdout: Redirection::Pipe,
            ..Default::default()
        },
    )
    .unwrap();

    // Obtain the output from the standard streams.
    let (out, _err) = p.communicate(None).unwrap();

    if let Some(_exit_status) = p.poll() {
        // the pocess has finished
    } else {
        // it is still running, terminate it
        p.terminate().unwrap();
    }

    out.unwrap()
}

pub fn parse_raw_data_key(input: String) -> (String, String) {
    // base64-encoded datakey, 0 ciphertext,1 plaintext
    // CIPHERTEXT: ciphertext \n PLAINTEXT: plaintext
    let parts: Vec<&str> = input.split("\n").collect();

    let ciphertext = parts[0].trim_start_matches("CIPHERTEXT: ");
    let plaintext = parts[1].trim_start_matches("PLAINTEXT: ");

    (ciphertext.to_string(), plaintext.to_string())
}

