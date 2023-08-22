module "kms" {
  source = "../kms"
}

module "dynamodb" {
  source = "../dynamodb"
}

resource "aws_iam_policy" "enclave_policy_template" {
  name        = "EnclavePolicyTemplate"
  path        = "/"
  description = "Your policy description"

  policy = <<EOF
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
            "Resource": "${module.kms.kms_arn}"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "dynamodb:PutItem",
                "dynamodb:GetItem"
            ],
            "Resource": "${module.dynamodb.dynamodb_arn}"
        }
    ]
}
EOF
}

resource "aws_iam_role" "enclave_role" {
  name               = "EnclaveRole-${replace(timestamp(), "/[^A-Za-z0-9]/", "")}"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "enclave_role_policy_attachment" {
  role       = aws_iam_role.enclave_role.name
  policy_arn = aws_iam_policy.enclave_policy_template.arn
}