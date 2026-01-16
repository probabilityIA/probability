# IAM Role para la instancia EC2 existente - Probability
# Este role permite que la EC2 haga pull de imagenes de ECR privado

# Data source para obtener la instancia EC2 existente
data "aws_instance" "existing_ec2" {
  filter {
    name   = "tag:Name"
    values = [var.ec2_instance_name]
  }
}

# IAM Role para EC2
resource "aws_iam_role" "ec2_ecr_role" {
  name = "probability-ec2-ecr-pull-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(
    var.tags,
    {
      Name    = "probability-ec2-ecr-pull-role"
      Purpose = "ECR access for EC2"
    }
  )
}

# Politica para permitir pull de ECR (usando la politica managed de AWS)
resource "aws_iam_role_policy_attachment" "ec2_ecr_readonly" {
  role       = aws_iam_role.ec2_ecr_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

# Politica adicional para permitir autenticacion en ECR (privado y publico)
resource "aws_iam_role_policy" "ec2_ecr_auth" {
  name = "probability-ec2-ecr-auth-policy"
  role = aws_iam_role.ec2_ecr_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "ecr:DescribeRepositories",
          "ecr:DescribeImages",
          "ecr:ListImages"
        ]
        Resource = [
          aws_ecr_repository.frontend.arn,
          aws_ecr_repository.backend.arn,
          aws_ecr_repository.nginx.arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "ecr-public:GetAuthorizationToken",
          "ecr-public:BatchCheckLayerAvailability",
          "ecr-public:GetDownloadUrlForLayer",
          "ecr-public:BatchGetImage",
          "ecr-public:DescribeRepositories",
          "ecr-public:DescribeImages",
          "ecr-public:ListImages"
        ]
        Resource = "*"
      }
    ]
  })
}

# Instance Profile para asociar el role a la EC2
resource "aws_iam_instance_profile" "ec2_ecr_profile" {
  name = "probability-ec2-ecr-pull-profile"
  role = aws_iam_role.ec2_ecr_role.name

  tags = var.tags
}

# Outputs
output "ec2_iam_role_arn" {
  description = "ARN del IAM Role para EC2"
  value       = aws_iam_role.ec2_ecr_role.arn
}

output "ec2_instance_profile_name" {
  description = "Nombre del Instance Profile para asociar a la EC2"
  value       = aws_iam_instance_profile.ec2_ecr_profile.name
}

output "ec2_instance_id" {
  description = "ID de la instancia EC2 encontrada"
  value       = try(data.aws_instance.existing_ec2.id, "No encontrada - ajusta los filtros en iam_roles.tf")
}
