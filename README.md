# what is in this repo
- this repo contains demo of use Amazon Nitro Enclave and Amazon KMS which implement the crypto wallet on AWS Cloud
- this repo contains sample code for golang, python, and rust. You can find them in **nitro-golang-example**, **nitro-python-example**, and **nitro-rust-example** respectively.
- this repo supports to deploy on Nitro Enclave on EKS, you can follow the steps in README.md file in each sub-folder

# two main workflow 

this demo code implement two use case, **generateAccount** and **signature**

## generateAccount
- generateAccount API generates wallet in Nitro Enclave and uses KMS to encrypt it with envelope encryption.

- below is the process of the workflow

![generateAccount process](/image/generateAccount.png)

## sign
- signature API signs a transaction with the private key of the wallet

- below is the process of the signature

![sign](/image/sign.png)

# important configurations
when you try to run the Nitro Enclave application, you should configure below things

- **Region** : you need choose a region where you deploy your application, and you need set **ENV** in Dockerfile
```
ENV REGION ap-northeast-1
```
- **KMS** : you need create a **Symmetric** kms key, which used for **Encrypt and decrypt**, you need copy the kms id. In this demo, it is hardcode in the appClient code
- **IAM Role** : IAM Role should assign the policy to allow call KMS encrypt,generateDataKey. In this demo, we attach the role on EC2
- **vsock-proxy** : before you start the enclave application, you should start the vsock-proxy for kms. below command with run the proxy on parent instance which will forward request on port 8000 to endpoint *kms.ap-southeast-1.amazonaws.com* on port 443. you should run it before run enclave
```
vsock-proxy 8000 kms.ap-northeast-1.amazonaws.com 443 &
```
- **cid** : vsock connection use CID not address. when the vsock-proxy start, CID **3** is default parent CID, according to this  [doc](!https://github.com/aws/aws-nitro-enclaves-sdk-c/blob/main/bin/kmstool-enclave-cli/main.c#L18)