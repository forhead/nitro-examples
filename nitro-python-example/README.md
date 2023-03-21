# Nitro Enclave demo with python code

## dynamodb table 
we use dynamodb to store the wallet related content, table designed as below.

Table name
- AccountTable

Colume
- keyId: kms alias id which used for encryption for the private key
- name: account name for this account
- encryptedPrivateKey: encrypted wallet private key
- address: the address key of the wallet
- encryptedDataKey: the data key used to encrypt the private key


## IAM Role
you need create a IAM Role which will be attached to EC2/EKS, it need have the access for kms and dynamodb. you need update this policy after your enclave image created with condition check of PCR
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "kms:Decrypt",
                "kms:Encrypt",
                "kms:GenerateDataKey"
            ],
            "Resource": "your kms arn"
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

## kms key
you need create a **Symmetric** kms key, which used for **Encrypt and decrypt**, you need copy this
