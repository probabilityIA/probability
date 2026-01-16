# Variables de Terraform - Probability

variable "aws_region" {
  description = "Region de AWS donde se desplegaran los recursos"
  type        = string
  default     = "us-east-1"
}

variable "ec2_instance_name" {
  description = "Nombre (tag) de la instancia EC2 existente"
  type        = string
  default     = "probability"
}

variable "environment" {
  description = "Ambiente de despliegue"
  type        = string
  default     = "production"
}

variable "ecr_image_retention_count" {
  description = "Numero de imagenes a mantener en ECR (lifecycle policy)"
  type        = number
  default     = 3
}

variable "tags" {
  description = "Tags comunes para todos los recursos"
  type        = map(string)
  default = {
    Project     = "probability"
    ManagedBy   = "terraform"
    Environment = "production"
  }
}
