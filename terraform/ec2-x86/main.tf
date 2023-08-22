
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

module "iam" {
  source = "../iam"
}

module "ec2_security_group" {
  source = "../ec2-security-group"
}

provider "aws" {
  region = "ap-east-1"
}

resource "aws_iam_instance_profile" "demo_profile" {
  name = "demo_profile"
  role = module.iam.iam_instance_profile_name
}

resource "aws_instance" "app_server" {
  ami                    = "ami-042268d06b381c465"
  instance_type          = "m5.xlarge"
  key_name               = "nitro-enclave-deployer-key"
  vpc_security_group_ids = [module.ec2_security_group.security_group_id]
  iam_instance_profile   = aws_iam_instance_profile.demo_profile.name
  tags = {
    Name = "EC2Instance-NitroEnclave-X86"
  }
  root_block_device {
    volume_size = 20 # in GB 
    encrypted   = true
  }

  enclave_options {
    enabled = true
  }

  provisioner "remote-exec" {
    inline = [
      "sudo yum update -y",
      "sudo amazon-linux-extras install aws-nitro-enclaves-cli -y",
      "sudo yum install aws-nitro-enclaves-cli-devel -y",
      "sudo yum install openssl openssl-devel -y",
      "sudo usermod -aG ne ec2-user",
      "sudo usermod -aG docker ec2-user",
      "sudo systemctl start nitro-enclaves-allocator.service",
      "sudo systemctl enable --now nitro-enclaves-allocator.service",
      "sudo systemctl start nitro-enclaves-vsock-proxy.service",
      "sudo systemctl enable --now nitro-enclaves-vsock-proxy.service",
      "sudo systemctl start docker",
      "sudo systemctl enable --now docker",
      "curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh -s -- -y",
      "sudo yum install gcc -y",
    ]
  }

  connection {
    type        = "ssh"
    user        = "ec2-user"
    private_key = file("~/.ssh/id_rsa.pem")
    host        = self.public_ip
  }
}

