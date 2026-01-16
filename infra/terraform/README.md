# 🏗️ Infraestructura como Código - Terraform

Este directorio contiene la configuración de Terraform para gestionar la infraestructura de AWS.

## 📋 Estructura

```
terraform/
├── main.tf              # Configuración principal y provider
├── backend.tf           # Backend remoto en S3
├── ecr.tf               # Repositorios ECR privados
├── iam_roles.tf         # IAM Roles y políticas
├── outputs.tf           # Outputs de los recursos
└── terraform.tfstate     # Estado local (backup)
```

## 🚀 Inicio Rápido

### 1. Inicializar Terraform

```bash
cd infra/terraform
terraform init
```

### 2. Revisar el plan

```bash
terraform plan
```

### 3. Aplicar los cambios

```bash
terraform apply
```

## 📦 Recursos Creados

### Repositorios ECR

- `monorepo-auth`: Servicio de autenticación
- `monorepo-api`: Servicio API principal
- `monorepo-worker`: Servicio worker

Cada repositorio tiene una política de ciclo de vida que mantiene solo 1 imagen (la última). Cada nueva imagen reemplaza automáticamente la anterior.

### IAM Role para EC2

- **Role**: `ec2-ecr-pull-role`
- **Política**: `AmazonEC2ContainerRegistryReadOnly` + permisos adicionales para ECR
- **Instance Profile**: `ec2-ecr-pull-profile`

## 🔧 Configuración

### Asociar el IAM Role a la EC2

Después de aplicar Terraform, asocia el Instance Profile a tu instancia EC2:

```bash
# Opción 1: Desde la consola de AWS
# EC2 > Instances > Select instance > Actions > Security > Modify IAM role

# Opción 2: Desde AWS CLI
aws ec2 associate-iam-instance-profile \
  --instance-id i-0123456789abcdef0 \
  --iam-instance-profile Name=ec2-ecr-pull-profile
```

### Variables de Entorno

El workflow de GitHub Actions necesita estos secrets:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `EC2_SSH_KEY`
- `EC2_HOST`
- `EC2_USER`

## 📝 Notas Importantes

1. **Account ID**: El Account ID está hardcodeado en algunos archivos (`334689162817`). Ajusta según tu cuenta AWS.

2. **EC2 Instance**: El data source en `iam_roles.tf` busca una instancia por tag `Name=rupu-production`. Ajusta el filtro según tu configuración.

3. **Rutas**: Ajusta las rutas en el workflow de GitHub Actions y en `podman-compose.yml` según tu estructura de directorios.

4. **Backend S3**: El estado de Terraform se guarda en `terraform-state-bucket-rupu`. Este bucket ya existe y está configurado.

## 🔄 Workflow de CI/CD

El workflow `.github/workflows/deploy.yml` se activa cuando hay cambios en `services/` y:

1. Construye las imágenes de los 3 servicios usando Podman
2. Sube las imágenes a ECR privado
3. Se conecta vía SSH a la EC2
4. Ejecuta `podman-compose pull` y `podman-compose up -d`

## 🐛 Troubleshooting

### Error: "No EC2 instance found"

Ajusta el filtro en `iam_roles.tf`:

```hcl
data "aws_instance" "existing_ec2" {
  filter {
    name   = "tag:Name"
    values = ["tu-nombre-de-instancia"]
  }
}
```

### Error: "ECR repository not found"

Asegúrate de haber aplicado `terraform apply` antes de ejecutar el workflow.

### Error: "Permission denied" en EC2

Verifica que el Instance Profile esté asociado a la instancia EC2.

## 📚 Recursos

- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [AWS ECR Documentation](https://docs.aws.amazon.com/ecr/)
- [Podman Documentation](https://docs.podman.io/)
