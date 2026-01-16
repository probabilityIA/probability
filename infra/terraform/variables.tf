# Variables de Terraform

variable "aws_region" {
  description = "Región de AWS donde se desplegarán los recursos"
  type        = string
  default     = "us-east-1"
}

variable "ec2_instance_name" {
  description = "Nombre (tag) de la instancia EC2 existente"
  type        = string
  default     = "cam"
}

variable "environment" {
  description = "Ambiente de despliegue"
  type        = string
  default     = "production"
}

variable "ecr_image_retention_count" {
  description = "Número de imágenes a mantener en ECR (lifecycle policy) - Actualmente fijado en 1"
  type        = number
  default     = 1
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
