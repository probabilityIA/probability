#!/bin/bash

# Script para instalar Terraform en Ubuntu/Debian
# Probability - Instalación de Terraform

echo "📦 Instalando Terraform..."

# Agregar repositorio de HashiCorp
sudo apt-get update
sudo apt-get install -y gnupg software-properties-common

# Agregar clave GPG de HashiCorp
wget -O- https://apt.releases.hashicorp.com/gpg | \
    gpg --dearmor | \
    sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null

# Agregar repositorio
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
    sudo tee /etc/apt/sources.list.d/hashicorp.list

# Instalar Terraform
sudo apt-get update
sudo apt-get install -y terraform

# Verificar instalación
echo ""
echo "✅ Verificando instalación..."
terraform version

echo ""
echo "✅ Terraform instalado correctamente!"
echo "Ahora puedes ejecutar: cd infra/terraform && terraform init"


