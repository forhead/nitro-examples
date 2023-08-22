resource "aws_dynamodb_table" "AccountTable" {
  name           = "AccountTable1"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "name"

  attribute {
    name = "name"
    type = "S"
  }
}
