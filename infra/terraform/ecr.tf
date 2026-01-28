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

# Lifecycle Policy para frontend - Mantener últimas 5 imágenes con CalVer
resource "aws_ecr_lifecycle_policy" "frontend_lifecycle" {
  repository = aws_ecr_repository.frontend.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Mantener tag latest siempre"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["latest"]
          countType     = "imageCountMoreThan"
          countNumber   = 999
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Mantener últimas 5 imágenes tagueadas"
        selection = {
          tagStatus   = "tagged"
          countType   = "imageCountMoreThan"
          countNumber = 5
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 3
        description  = "Eliminar imágenes sin tag después de 7 días"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 7
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# Lifecycle Policy para backend - Mantener últimas 5 imágenes con CalVer
resource "aws_ecr_lifecycle_policy" "backend_lifecycle" {
  repository = aws_ecr_repository.backend.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Mantener tag latest siempre"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["latest"]
          countType     = "imageCountMoreThan"
          countNumber   = 999
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Mantener últimas 5 imágenes tagueadas"
        selection = {
          tagStatus   = "tagged"
          countType   = "imageCountMoreThan"
          countNumber = 5
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 3
        description  = "Eliminar imágenes sin tag después de 7 días"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 7
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# Lifecycle Policy para nginx - Mantener últimas 5 imágenes con CalVer
resource "aws_ecr_lifecycle_policy" "nginx_lifecycle" {
  repository = aws_ecr_repository.nginx.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Mantener tag latest siempre"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["latest"]
          countType     = "imageCountMoreThan"
          countNumber   = 999
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Mantener últimas 5 imágenes tagueadas"
        selection = {
          tagStatus   = "tagged"
          countType   = "imageCountMoreThan"
          countNumber = 5
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 3
        description  = "Eliminar imágenes sin tag después de 7 días"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 7
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# Repositorio para el website (probability-website)
resource "aws_ecr_repository" "website" {
  name                 = "probability-website"
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
      Name    = "probability-website"
      Service = "website"
    }
  )
}

# Lifecycle Policy para website - Mantener últimas 5 imágenes con CalVer
resource "aws_ecr_lifecycle_policy" "website_lifecycle" {
  repository = aws_ecr_repository.website.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Mantener tag latest siempre"
        selection = {
          tagStatus     = "tagged"
          tagPrefixList = ["latest"]
          countType     = "imageCountMoreThan"
          countNumber   = 999
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Mantener últimas 5 imágenes tagueadas"
        selection = {
          tagStatus   = "tagged"
          countType   = "imageCountMoreThan"
          countNumber = 5
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 3
        description  = "Eliminar imágenes sin tag después de 7 días"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 7
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

output "ecr_repository_website_url" {
  description = "URL del repositorio ECR para website"
  value       = aws_ecr_repository.website.repository_url
}
