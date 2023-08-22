# Manual AWS Configuration Steps

1. **Create an EC2 Security Group**
   - Navigate to the EC2 dashboard in the AWS Management Console.
   - Click on 'Security Groups' under 'Network & Security'.
   - Click 'Create security group', provide a name, description, and select the VPC you want to create it in.

2. **Create an EC2 Instance**
   - From the EC2 dashboard, click 'Instances'.
   - Click 'Launch instances' and follow the wizard to create a new instance.

3. **Create an IAM Policy**
   - Navigate to the IAM dashboard.
   - Click on 'Policies' and then 'Create policy'.
   - Define the policy according to your requirements.

4. **Create an IAM Role**
   - In the IAM dashboard, click 'Roles' and then 'Create role'.
   - Select the service that will use this role and define the permissions.

5. **Attach IAM Role to EC2 Instance**
   - Go back to the EC2 dashboard and select your instance.
   - Click 'Actions', then 'Security', and then 'Modify IAM role'.
   - Select the role you created and save.

6. **Create a Key in KMS (Key Management Service)**
   - Navigate to the KMS dashboard.
   - Click 'Create key' and follow the wizard to create a new key.

7. **Create a Table in DynamoDB**
   - Go to the DynamoDB dashboard.
   - Click 'Create table', provide a name and define the primary key.

# Terraform Dependencies
kms/main.tf
dynamodb/main.tf
iam/main.tf
ec2-security-group/main.tf
ec2-x86/main.tf

# how to build
1. create aws_access_key_id,aws_secret_access_key in AWS console security_credentials.
2. create key pair, modify public_key as your own.
```sh
cd ec2-key-pair
terraform apply 
```
3. create ec2 instance which depends on special(kms, dynamodb, iam, EC2 Security Group)

# Build Instructions

Follow these steps to set up your AWS environment:

1. **Generate AWS Access Keys**
   - Go to the AWS console and navigate to your security credentials.
   - Create a new `aws_access_key_id` and `aws_secret_access_key`.

2. **Install Terraform**
   - [Terraform Installation Guide](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
   - To use your IAM credentials to authenticate the Terraform AWS provider, set the AWS_ACCESS_KEY_ID environment variable.
```sh
   export AWS_ACCESS_KEY_ID=<your ak>
   export AWS_SECRET_ACCESS_KEY=<your sk>
```

2. **Create and Configure Key Pair**
   - Navigate to the `ec2-key-pair` directory.
   - Run the following command to create a new key pair and modify the public key to your own: ```cd ec2-key-pair && terraform apply```

3. **Create EC2 Instance**
   - Create an EC2 instance that depends on specific AWS services (KMS, DynamoDB, IAM, EC2 Security Group):  ```cd ec2-x86 && terraform apply``` 
