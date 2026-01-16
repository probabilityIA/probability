# IAM Role para la instancia EC2 existente - Probability
# Este role permite que la EC2 haga pull de imagenes de ECR privado

# Data source para obtener la instancia EC2 existente
# Si se proporciona instance_id, usarlo directamente; si no, buscar por tag Name
data "aws_instance" "existing_ec2" {
  count = var.ec2_instance_id != "" ? 0 : 1
  
  filter {
    name   = "tag:Name"
    values = [var.ec2_instance_name]
  }
  
  filter {
    name   = "instance-state-name"
    values = ["running"]
  }
}

# Data source alternativo usando instance ID directamente
data "aws_instance" "existing_ec2_by_id" {
  count       = var.ec2_instance_id != "" ? 1 : 0
  instance_id = var.ec2_instance_id
}

# Local para obtener el instance ID
locals {
  instance_id = var.ec2_instance_id != "" ? data.aws_instance.existing_ec2_by_id[0].id : data.aws_instance.existing_ec2[0].id
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

# Asociar el Instance Profile a la instancia EC2 automáticamente
# Nota: AWS no tiene un recurso directo para asociar instance profiles a instancias existentes
# Usamos null_resource con local-exec para ejecutar el comando AWS CLI
resource "null_resource" "associate_iam_profile" {
  depends_on = [
    aws_iam_instance_profile.ec2_ecr_profile,
    aws_iam_role_policy.ec2_ecr_auth
  ]

  # Desasociar rol anterior si existe
  provisioner "local-exec" {
    command = <<-EOT
      set -e
      
      # Esperar un momento para que el instance profile se propague
      echo "Esperando propagación del instance profile..."
      sleep 5
      
      # Obtener association ID actual si existe
      ASSOC_ID=$(aws ec2 describe-iam-instance-profile-associations \
        --filters "Name=instance-id,Values=${local.instance_id}" \
        --query 'IamInstanceProfileAssociations[0].AssociationId' \
        --output text 2>/dev/null || echo "")
      
      # Desasociar si existe
      if [ "$ASSOC_ID" != "None" ] && [ -n "$ASSOC_ID" ]; then
        echo "Desasociando rol anterior: $ASSOC_ID"
        aws ec2 disassociate-iam-instance-profile --association-id "$ASSOC_ID" || true
        sleep 3
      fi
      
      # Asociar nuevo rol usando el ARN
      PROFILE_ARN="${aws_iam_instance_profile.ec2_ecr_profile.arn}"
      echo "Asociando nuevo rol: $PROFILE_ARN"
      aws ec2 associate-iam-instance-profile \
        --instance-id ${local.instance_id} \
        --iam-instance-profile Arn="$PROFILE_ARN"
      
      echo "✅ Rol asociado correctamente"
    EOT
    interpreter = ["bash", "-c"]
  }

  # Trigger cuando cambie el instance profile
  triggers = {
    instance_profile_name = aws_iam_instance_profile.ec2_ecr_profile.name
    instance_id = local.instance_id
  }
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
  value       = local.instance_id
}
