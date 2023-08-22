output "kms_arn" {
  value       = aws_kms_key.enclave_key.arn
  description = "The ARN of the KMS key"
}
