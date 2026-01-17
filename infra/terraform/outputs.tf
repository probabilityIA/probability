# Outputs centralizados para facilitar el acceso a los recursos creados - Probability

output "ecr_repositories" {
  description = "URLs de los repositorios ECR"
  value = {
    frontend = aws_ecr_repository.frontend.repository_url
    backend  = aws_ecr_repository.backend.repository_url
    nginx    = aws_ecr_repository.nginx.repository_url
    website  = aws_ecr_repository.website.repository_url
  }
}

output "iam_role_info" {
  description = "Informacion del IAM Role para EC2"
  value = {
    role_arn         = aws_iam_role.ec2_ecr_role.arn
    instance_profile = aws_iam_instance_profile.ec2_ecr_profile.name
  }
}

output "terraform_state_bucket" {
  description = "Bucket S3 para el estado de Terraform"
  value       = aws_s3_bucket.terraform_state.bucket
}
