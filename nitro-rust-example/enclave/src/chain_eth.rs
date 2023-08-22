use anychain_ethereum::{address::EthereumAddress, public_key::EthereumPublicKey};
use base64::{engine::general_purpose, Engine as _};
use libsecp256k1::{PublicKey, SecretKey};
use rand::rngs::OsRng;
pub fn generate_random_secret_key() -> (String, String) {
    let mut rng = OsRng;

    let secret_key = SecretKey::random(&mut rng);
    let secret_key_bytes = secret_key.serialize();
    let secret_key_base64 = general_purpose::STANDARD.encode(secret_key_bytes);
    // println!("Generated secret key: {:?}", secret_key);
    // println!("Generated secret key (base64): {:?}", secret_key_base64);

    let public_key = PublicKey::from_secret_key(&secret_key);
    // println!("Generated public key: {:?}", public_key);

    let ethereum_public_key = EthereumPublicKey::from_secp256k1_public_key(public_key);
    let ethereum_address = EthereumAddress::checksum_address(&ethereum_public_key);
    // println!(
    //     "Generated Ethereum address: {:?}",
    //     ethereum_address.to_string()
    // );

    (secret_key_base64, ethereum_address.to_string())
}
