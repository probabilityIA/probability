provider "aws" {
  region = var.aws_region
}

# Este es el bucket donde vivirá el estado de TODA tu infraestructura
resource "aws_s3_bucket" "terraform_state" {
  bucket = "terraform-state-bucket-probability" # DEBE SER ÚNICO
  
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