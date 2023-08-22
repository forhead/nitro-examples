# aws-nitro-enclaves-rs

## Generate Account
![generateAccount.png](..%2Fimage%2FgenerateAccount.png)
### Account Generation Workflow

1. parent instance client receive generateccount API call
2. call getlAMToken function which get credential of IAM Role
3. send credential and generateAccount API call via vsock
4. generate account in Nitro Enclave
5. call KMS API generateDataKey with credential and attestation
6. encrypt the account content with datakey
7. send the encrypted account content to parent instance via vsock
8. call API to save the encrypted account content to dynamodb

## Sign Signature by Private Key
![sign.png](..%2Fimage%2Fsign.png)

### Transaction Signature Workflow

1. parent instance client receive sign API call
2. call getlAMToken function which get credential of IAM Role
3. send credential and sign API call via vsock
4. decrypt the encrypted datakey with KMS API decrypt
5. decrypt the encrypted wallet's private key with datakey
6. sign the message with wallet's private key
7. send signature to parent instance via vsock

## Core Components

### enclave server
enclave/bin/enclave-server

How to build:
```sh
cd enclave
make server
```

### parent client
```sh
cd enclave
make client
```

# important configurations
when you try to run the Nitro Enclave application, you should configure below things

- **AWS Region** : you need choose a region where you deploy your application, and you need set **ENV** in Dockerfile
```
ENV REGION ap-east-1
```

- **AWS KMS**
you need create a **Symmetric** kms key, which used for **Encrypt and decrypt**
```sh
cargo run --bin aws-kms-create-key -- -r ap-east-1
```

- **AWS IAM Role**
 
you need create a IAM Role which will be attached to EC2/EKS, it need have the access for kms and dynamodb. you need update this policy after your enclave image created with condition check of PCR
`EnclavePolicyKmsDynamodbTemplate`
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "kms:Decrypt",
                "kms:GenerateDataKey"
            ],
            "Resource": "<Your KMS ARN>"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "dynamodb:PutItem",
                "dynamodb:GetItem"
            ],
            "Resource": "*"
        }
    ]
}
```

`EnclaveRole`, attach policy EnclavePolicyKmsDynamodbTemplate, attach the role on EC2

- **AWS DynamoDB**

```sh
cargo run --bin aws-dynamodb-create-table -- -r ap-east-1 -t AccountTable -k name
```

 Table Name:  AccountTable 


| Column               | Description                                         |
| -------------------- | --------------------------------------------------- |
| KeyId                | KMS key id used for encryption on the wallet private key |
| Name                 | Account name for this account, used to identify wallet |
| EncryptedPrivateKey  | Encrypted wallet private key                        |
| Address              | The address key of the wallet                       |
| EncryptedDataKey     | The data key used to encrypt the private key        |


- **vsock-proxy** : before you start the enclave application, you should start the vsock-proxy for kms. below command with run the proxy on parent instance which will forward request on port 8000 to endpoint kms.ap-east-1.amazonaws.com on port 443. you should run it before run enclave
```sh
sudo systemctl start nitro-enclaves-allocator.service
sudo systemctl enable --now nitro-enclaves-vsock-proxy.service
``` 
or
```sh
vsock-proxy 8000 kms.east-northeast-1.amazonaws.com 443 &
```

### Anychain Ethereum
### Terraform AWS Provider
