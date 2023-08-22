output "dynamodb_arn" {
  value       = aws_dynamodb_table.AccountTable.arn
  description = "The ARN of the AccountTable"
}
