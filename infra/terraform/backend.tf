# Backend remoto en S3 para el estado de Terraform
# NOTA: Comentar temporalmente para crear el bucket primero
# Después de crear el bucket, descomentar esto y ejecutar: terraform init -migrate-state
terraform {
  # backend "s3" {
  #   bucket         = "terraform-state-bucket-probability"
  #   key            = "probability/terraform.tfstate"
  #   region         = "us-east-1"
  #   encrypt        = true
  #   dynamodb_table = null # Opcional: puedes crear una tabla DynamoDB para state locking
  # }
}
