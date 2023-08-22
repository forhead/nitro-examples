resource "aws_kms_key" "enclave_key" {
  description             = "Symmetric KMS key for Nitro Enclaves"
  deletion_window_in_days = 10
  key_usage               = "ENCRYPT_DECRYPT"
  is_enabled              = true
}
