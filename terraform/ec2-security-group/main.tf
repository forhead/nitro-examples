provider "aws" {
  region = "ap-east-1"
}

resource "aws_security_group" "nitro_enclave_group" {
  name = "nitro_enclave_group_${replace(timestamp(), "/[^A-Za-z0-9]/", "")}"
  description = "nitro enclave security group"

  ingress {
    description = "SSH from anywhere"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}