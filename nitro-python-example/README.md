# how to deploy

# 1.AWS Configuration
you need configure below things on AWS Cloud
### dynamodb table 
this demo uses dynamodb to store the wallet related content, code automatically creates the table if not exists

Table name
- AccountTable

Colume
- KeyId: kms key id which used for encryption on the wallet private key
- Name: account name for this account, used for identify wallet
- EncryptedPrivateKey: encrypted wallet private key
- Address: the address key of the wallet
- EncryptedDataKey: the data key used to encrypt the private key

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

# 2. how to run 
golang demo code has two ways of deployment. **EC2** and **EKS**, we provide scripts to run the demo

## 2.1 run on EC2
before you run this demo, you need setup Nitro Enclave service on EC2, you can follow this [doc](https://docs.aws.amazon.com/enclaves/latest/user/nitro-enclave-cli-install.html)

you can follow below steps to run the demo code on EC2 

- 1. go to enclave folder to run enclave server
```
cd enclave
sh buildEnclaveAndRunOnEC2.sh
```
- 2. open another terminal, and go to parent folder to run app client
```
cd parent
python appClient.py
```

### check dynamodb in your region 
in dynamodb you can see there's a new row inserted
![dynamodb result](/image/dynamodb_query_result.png)
