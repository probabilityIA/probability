#!/bin/bash

# Script de despliegue para ECR público
# Probability - Backend Central

set -e

# Variables
IMAGE_NAME="probability-back-central"
ECR_REPO="public.ecr.aws/d3a6d4r1/cam/probability-back-central"
VERSION=${1:-"latest"}
DOCKERFILE_PATH="docker/Dockerfile"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Iniciando despliegue de Probability Backend Central${NC}"
echo -e "${YELLOW}Versión: ${VERSION}${NC}"

# Verificar que estamos en el directorio correcto
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: No se encontró go.mod. Ejecuta desde el directorio raíz del proyecto${NC}"
    exit 1
fi

# Verificar que Docker esté corriendo
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Docker no está corriendo${NC}"
    exit 1
fi

# Verificar que AWS CLI esté configurado
if ! aws sts get-caller-identity > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: AWS CLI no está configurado correctamente${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Verificaciones completadas${NC}"

# Limpiar dependencias
echo -e "${YELLOW}📦 Limpiando dependencias...${NC}"
go mod tidy

# Construir la imagen
echo -e "${YELLOW}🔨 Construyendo imagen Docker...${NC}"
docker buildx build --platform linux/amd64 -f ${DOCKERFILE_PATH} -t ${IMAGE_NAME}:${VERSION} --load .

# Etiquetar para ECR con nombres descriptivos
echo -e "${YELLOW}🏷️ Etiquetando imagen para ECR...${NC}"

# Crear tags descriptivos
if [ "${VERSION}" = "latest" ]; then
    # Para latest, crear múltiples tags descriptivos
    TIMESTAMP=$(date +%Y%m%d)
    DESCRIPTIVE_TAG="backend-latest"
    DATED_TAG="backend-${TIMESTAMP}"
    
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${VERSION}
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DESCRIPTIVE_TAG}
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DATED_TAG}
    
    echo -e "${GREEN}📅 Tags creados: latest, ${DESCRIPTIVE_TAG}, ${DATED_TAG}${NC}"
else
    # Para versiones específicas, crear tag descriptivo
    DESCRIPTIVE_TAG="backend-${VERSION}"
    
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${VERSION}
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DESCRIPTIVE_TAG}
    
    echo -e "${GREEN}🏷️ Tags creados: ${VERSION}, ${DESCRIPTIVE_TAG}${NC}"
fi

# Login a ECR público
echo -e "${YELLOW}🔐 Haciendo login a ECR público...${NC}"
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

# Push de las imágenes
echo -e "${YELLOW}⬆️ Subiendo imágenes a ECR...${NC}"

if [ "${VERSION}" = "latest" ]; then
    # Subir todos los tags para latest
    docker push ${ECR_REPO}:${VERSION}
    docker push ${ECR_REPO}:${DESCRIPTIVE_TAG}
    docker push ${ECR_REPO}:${DATED_TAG}
    echo -e "${GREEN}✅ Imágenes subidas con tags: latest, ${DESCRIPTIVE_TAG}, ${DATED_TAG}${NC}"
else
    # Subir tags para versiones específicas
    docker push ${ECR_REPO}:${VERSION}
    docker push ${ECR_REPO}:${DESCRIPTIVE_TAG}
    echo -e "${GREEN}✅ Imágenes subidas con tags: ${VERSION}, ${DESCRIPTIVE_TAG}${NC}"
fi

echo -e "${GREEN}🎉 Despliegue completado exitosamente!${NC}"
echo -e "${YELLOW}📋 Para usar la imagen:${NC}"
if [ "${VERSION}" = "latest" ]; then
    echo -e "docker run --env-file .env -p 8080:8080 ${ECR_REPO}:${DESCRIPTIVE_TAG}"
    echo -e "${YELLOW}🔖 Opciones de tags disponibles:${NC}"
    echo -e "  - ${ECR_REPO}:latest"
    echo -e "  - ${ECR_REPO}:${DESCRIPTIVE_TAG}"
    echo -e "  - ${ECR_REPO}:${DATED_TAG}"
else
    echo -e "docker run --env-file .env -p 8080:8080 ${ECR_REPO}:${DESCRIPTIVE_TAG}"
    echo -e "${YELLOW}🔖 Tags disponibles:${NC}"
    echo -e "  - ${ECR_REPO}:${VERSION}"
    echo -e "  - ${ECR_REPO}:${DESCRIPTIVE_TAG}"
fi
echo -e "${YELLOW}🌐 URL del repositorio ECR:${NC}"
echo -e "https://gallery.ecr.aws/d3a6d4r1/cam/probability-back-central"