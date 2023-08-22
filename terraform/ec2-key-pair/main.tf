resource "aws_key_pair" "deployer" {
  key_name   = "nitro-enclave-deployer-key"
  public_key = file("~/.ssh/id_rsa.pub")
}
