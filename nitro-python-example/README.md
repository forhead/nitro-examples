# Nitro Enclave demo with python code

## dynamodb table 
we use dynamodb to store the wallet related content, table designed as below.

table name
- walletTable

colume
- walletName: wallet name for this wallet
- encryptedPrivateKey: encrypted wallet private key
- publicKey: the public key of the wallet
- encryptedDatakey: the data key used to encrypt the private key
- keyId: kms alias id which used for encryption for the private key

## 