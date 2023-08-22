output "security_group_id" {
  description = "The ARN of the security group"
  value       = aws_security_group.nitro_enclave_group.id
}
