#!/bin/bash

# Script para construir y subir la imagen personalizada de nginx a ECR (ARM64)
# Uso: ./deploy.sh [tag]
set -e

IMAGE_NAME="nginx-custom"
ECR_REPO="public.ecr.aws/d3a6d4r1/cam/reserve"
TAG=${1:-"nginx-latest"}
DOCKERFILE_PATH="Dockerfile"
PLATFORM="linux/arm64"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 Construyendo imagen nginx personalizada para ARM64...${NC}"

# Verificar Docker
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker no está corriendo${NC}"
    exit 1
fi

# Verificar AWS CLI
if ! aws sts get-caller-identity > /dev/null 2>&1; then
    echo -e "${RED}❌ AWS CLI no está configurado correctamente${NC}"
    exit 1
fi

# Build ARM64
docker buildx build --platform $PLATFORM -f $DOCKERFILE_PATH -t $IMAGE_NAME:$TAG . --load

# Tag para ECR
docker tag $IMAGE_NAME:$TAG $ECR_REPO:$TAG

# Login ECR público
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

# Push
docker push $ECR_REPO:$TAG

echo -e "${GREEN}✅ Imagen subida a ECR: $ECR_REPO:$TAG${NC}" 