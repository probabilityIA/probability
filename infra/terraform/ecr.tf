# Repositorios ECR privados para frontend, backend y nginx - Probability

# Repositorio para el frontend (probability-frontend)
resource "aws_ecr_repository" "frontend" {
  name                 = "probability-frontend"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = merge(
    var.tags,
    {
      Name    = "probability-frontend"
      Service = "frontend"
    }
  )
}

# Repositorio para el backend (probability-backend)
resource "aws_ecr_repository" "backend" {
  name                 = "probability-backend"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = merge(
    var.tags,
    {
      Name    = "probability-backend"
      Service = "backend"
    }
  )
}

# Repositorio para nginx (probability-nginx)
resource "aws_ecr_repository" "nginx" {
  name                 = "probability-nginx"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = merge(
    var.tags,
    {
      Name    = "probability-nginx"
      Service = "nginx"
    }
  )
}

# Lifecycle Policy para frontend - Mantener ultimas N imagenes
resource "aws_ecr_lifecycle_policy" "frontend_lifecycle" {
  repository = aws_ecr_repository.frontend.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Mantener solo ${var.ecr_image_retention_count} imagenes, eliminar las demas"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = var.ecr_image_retention_count
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# Lifecycle Policy para backend - Mantener ultimas N imagenes
resource "aws_ecr_lifecycle_policy" "backend_lifecycle" {
  repository = aws_ecr_repository.backend.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Mantener solo ${var.ecr_image_retention_count} imagenes, eliminar las demas"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = var.ecr_image_retention_count
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# Lifecycle Policy para nginx - Mantener ultimas N imagenes
resource "aws_ecr_lifecycle_policy" "nginx_lifecycle" {
  repository = aws_ecr_repository.nginx.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Mantener solo ${var.ecr_image_retention_count} imagenes, eliminar las demas"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = var.ecr_image_retention_count
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# Outputs para facilitar el uso en otros modulos
output "ecr_repository_frontend_url" {
  description = "URL del repositorio ECR para frontend"
  value       = aws_ecr_repository.frontend.repository_url
}

output "ecr_repository_backend_url" {
  description = "URL del repositorio ECR para backend"
  value       = aws_ecr_repository.backend.repository_url
}

output "ecr_repository_nginx_url" {
  description = "URL del repositorio ECR para nginx"
  value       = aws_ecr_repository.nginx.repository_url
}
