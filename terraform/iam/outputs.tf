output "iam_instance_profile_name" {
  description = "The ARN of the IAM role"
  value       = aws_iam_role.enclave_role.name
}
