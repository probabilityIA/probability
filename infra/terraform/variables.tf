# Variables de Terraform - Probability

variable "aws_region" {
  description = "Region de AWS donde se desplegaran los recursos"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "Perfil de AWS a usar (para seleccionar la cuenta correcta)"
  type        = string
  default     = "probability"
}

variable "ec2_instance_name" {
  description = "Nombre (tag) de la instancia EC2 existente"
  type        = string
  default     = "Probability"
}

variable "ec2_instance_id" {
  description = "ID de la instancia EC2 (opcional, si se proporciona se usa en lugar del nombre)"
  type        = string
  default     = ""
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
