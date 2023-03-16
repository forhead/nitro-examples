# what's this repo for
this repo contains AWS Nitro Enclave demo with golang and python code

# workflow implemented in this repo

this demo code implement two use case, generateWallet and signature

## generateWallet
generateWallet API call will generate the wallet content in Nitro Enclave and will use KMS to encrypt it with envelope encryption.

below is the process of the workflow

![generateWallet process](/image/generateWallet.png)

## sign
signature API call will sign a message with the private key of the wallet

below is the process of the signature
![sign](/image/sign.png)

## something need to know
parent instance with cid 3
must start vsock for kms

demo with cid 16 for enclave server

