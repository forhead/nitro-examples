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

- **IAM Role** : IAM Role should assign the policy to allow call KMS encrypt,generateDataKey,generateRandom
- **vsock-proxy** : before you start the enclave application, you should start the vsock-proxy for kms. below command with run the proxy on parent instance which will forward request on port 8000 to endpoint *kms.ap-southeast-1.amazonaws.com* on port443
```
vsock-proxy 8000 kms.ap-southeast-1.amazonaws.com 443 &
```
- **http forward in Nitro Enclave** : nitro enclave communicate with parent via vsock. and the kms call in the Nitro Enclave will call the kms endpoint. so need below configuration to forward kms call to vsock. the **traffic-forwarder.py** is an application do the traffic forward
```
ifconfig lo 127.0.0.1
# Add a hosts record, pointing API endpoint to local loopback
echo "127.0.0.1   kms.ap-southeast-1.amazonaws.com" >> /etc/hosts

# Run traffic forwarder in background and start the server
nohup python3 /app/traffic-forwarder.py 443 3 8000 &
```
- **CID** : vsock connection use CID not address. when the vsock-proxy start, CID **3** is default parent CID, according to this  [doc](!https://github.com/aws/aws-nitro-enclaves-sdk-c/blob/main/bin/kmstool-enclave-cli/main.c#L18)
- **KMS** : you need setup a SYMMETRIC_DEFAULT kms
- **Region** : you should set the region of your KMS at

