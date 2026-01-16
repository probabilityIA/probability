# Backend remoto en S3 para el estado de Terraform
terraform {
  backend "s3" {
    bucket         = "terraform-state-bucket-probability"
    key            = "probability/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = null # Opcional: puedes crear una tabla DynamoDB para state locking
  }
}
