provider "aws" {
  region = var.aws_region
}

# Bucket S3 para el estado de Terraform de Probability
resource "aws_s3_bucket" "terraform_state" {
  bucket = "terraform-state-bucket-probability"

  # Evita que se borre por accidente
  lifecycle {
    prevent_destroy = true
  }
}

# Habilitamos el versionado para poder recuperar estados anteriores si algo falla
resource "aws_s3_bucket_versioning" "enabled" {
  bucket = aws_s3_bucket.terraform_state.id
  versioning_configuration {
    status = "Enabled"
  }
}
