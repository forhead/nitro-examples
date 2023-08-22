use reqwest::Error;
use serde_json::Value;
use std::collections::HashMap;

pub async fn get_iam_token() -> Result<HashMap<String, String>, Error> {
    let client = reqwest::Client::new();
    let instance_profile_name = client
        .get("http://169.254.169.254/latest/meta-data/iam/security-credentials/")
        .send()
        .await?
        .text()
        .await?;

    let url = format!(
        "http://169.254.169.254/latest/meta-data/iam/security-credentials/{}",
        instance_profile_name
    );
    let response: Value = client.get(&url).send().await?.json().await?;

    let mut credential = HashMap::new();
    credential.insert(
        "aws_access_key_id".to_string(),
        response["AccessKeyId"].as_str().unwrap().to_string(),
    );
    credential.insert(
        "aws_secret_access_key".to_string(),
        response["SecretAccessKey"].as_str().unwrap().to_string(),
    );
    credential.insert(
        "aws_session_token".to_string(),
        response["Token"].as_str().unwrap().to_string(),
    );

    Ok(credential)
}
