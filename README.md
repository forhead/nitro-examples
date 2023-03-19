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

# important configurations
when you try to run the Nitro Enclave application, you should configure below things

- **Region** : you need choose a region where you deploy your application, and you need set **ENV** in Dockerfile 
- **KMS** : you need create a **Symmetric** kms key, which used for **Encrypt and decrypt**, you need copy this
- **IAM Role** : IAM Role should assign the policy to allow call KMS encrypt,generateDataKey
- **vsock-proxy** : before you start the enclave application, you should start the vsock-proxy for kms. below command with run the proxy on parent instance which will forward request on port 8000 to endpoint *kms.ap-southeast-1.amazonaws.com* on port443
```
vsock-proxy 8000 kms.ap-southeast-1.amazonaws.com 443 &
```
- **cid** : vsock connection use CID not address. when the vsock-proxy start, CID **3** is default parent CID, according to this  [doc](!https://github.com/aws/aws-nitro-enclaves-sdk-c/blob/main/bin/kmstool-enclave-cli/main.c#L18)