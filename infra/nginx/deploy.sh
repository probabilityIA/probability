#!/bin/bash

# Script para construir y subir la imagen personalizada de nginx a ECR (ARM64)
# Probability - Nginx
# Uso: ./deploy.sh [tag]
set -e

IMAGE_NAME="probability-nginx"
# Mismo repositorio que backend y frontend, diferentes etiquetas
ECR_REPO="public.ecr.aws/c1l9h7c9/probability"
VERSION=${1:-"latest"}
DOCKERFILE_PATH="Dockerfile"
PLATFORM="linux/arm64"
AWS_PROFILE_NAME="probability"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 Construyendo imagen nginx personalizada para ARM64...${NC}"
echo -e "${YELLOW}Versión: ${VERSION}${NC}"
echo -e "${YELLOW}Perfil de AWS: ${AWS_PROFILE_NAME}${NC}"

# Verificar Docker
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker no está corriendo${NC}"
    exit 1
fi

# Verificar AWS CLI con el perfil correcto
if ! aws --profile "${AWS_PROFILE_NAME}" sts get-caller-identity > /dev/null 2>&1; then
    echo -e "${RED}❌ AWS CLI no está configurado correctamente${NC}"
    exit 1
fi

# Verificar que buildx esté disponible
if ! docker buildx version > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker buildx no está disponible${NC}"
    echo -e "${YELLOW}💡 Instala buildx: https://docs.docker.com/buildx/working-with-buildx/${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Verificaciones completadas${NC}"

# Crear builder multi-arquitectura si no existe
echo -e "${YELLOW}🔧 Configurando builder multi-arquitectura...${NC}"
if ! docker buildx inspect multiarch-builder > /dev/null 2>&1; then
    docker buildx create --name multiarch-builder --driver docker-container --use
else
    docker buildx use multiarch-builder
fi

# Build ARM64
echo -e "${YELLOW}🔨 Construyendo imagen Docker para ${PLATFORM}...${NC}"
echo -e "${BLUE}   Esto puede tomar varios minutos...${NC}"
docker buildx build \
    --platform $PLATFORM \
    -f $DOCKERFILE_PATH \
    -t $IMAGE_NAME:$VERSION \
    --load \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Imagen construida exitosamente${NC}"
else
    echo -e "${RED}❌ Error construyendo la imagen${NC}"
    exit 1
fi

# Etiquetar para ECR con nombres descriptivos
echo -e "${YELLOW}🏷️ Etiquetando imagen para ECR...${NC}"

# Crear tags descriptivos
if [ "${VERSION}" = "latest" ]; then
    # Para latest, crear múltiples tags descriptivos
    TIMESTAMP=$(date +%Y%m%d)
    DESCRIPTIVE_TAG="nginx-latest"
    DATED_TAG="nginx-${TIMESTAMP}"
    
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DESCRIPTIVE_TAG}
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DATED_TAG}
    
    echo -e "${GREEN}📅 Tags creados: ${DESCRIPTIVE_TAG}, ${DATED_TAG}${NC}"
else
    # Para versiones específicas, crear tag descriptivo
    DESCRIPTIVE_TAG="nginx-${VERSION}"
    
    docker tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DESCRIPTIVE_TAG}
    
    echo -e "${GREEN}🏷️ Tags creados: ${DESCRIPTIVE_TAG}${NC}"
fi

# Login ECR público
echo -e "${YELLOW}🔐 Haciendo login a ECR público con el perfil '${AWS_PROFILE_NAME}'...${NC}"
aws --profile "${AWS_PROFILE_NAME}" ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

# Push de las imágenes
echo -e "${YELLOW}⬆️ Subiendo imágenes a ECR...${NC}"
echo -e "${BLUE}   Esto puede tomar varios minutos dependiendo de tu conexión...${NC}"

if [ "${VERSION}" = "latest" ]; then
    # Subir todos los tags para latest
    docker push ${ECR_REPO}:${DESCRIPTIVE_TAG}
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Imagen subida con tag: ${DESCRIPTIVE_TAG}${NC}"
    else
        echo -e "${RED}❌ Error subiendo imagen con tag: ${DESCRIPTIVE_TAG}${NC}"
        exit 1
    fi
    
    docker push ${ECR_REPO}:${DATED_TAG}
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Imagen subida con tag: ${DATED_TAG}${NC}"
    else
        echo -e "${RED}❌ Error subiendo imagen con tag: ${DATED_TAG}${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Todas las imágenes subidas exitosamente${NC}"
else
    # Subir tags para versiones específicas
    docker push ${ECR_REPO}:${DESCRIPTIVE_TAG}
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Imagen subida con tag: ${DESCRIPTIVE_TAG}${NC}"
    else
        echo -e "${RED}❌ Error subiendo imagen con tag: ${DESCRIPTIVE_TAG}${NC}"
        exit 1
    fi
fi

echo ""
echo -e "${GREEN}🎉 Despliegue completado exitosamente!${NC}"
echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📋 Información de la imagen desplegada:${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if [ "${VERSION}" = "latest" ]; then
    echo -e "${BLUE}🔖 Tags disponibles:${NC}"
    echo -e "   • ${ECR_REPO}:${DESCRIPTIVE_TAG}"
    echo -e "   • ${ECR_REPO}:${DATED_TAG}"
else
    echo -e "${BLUE}🔖 Tag disponible:${NC}"
    echo -e "   • ${ECR_REPO}:${DESCRIPTIVE_TAG}"
fi

echo ""
echo -e "${BLUE}🐳 Para ejecutar en producción (ARM64):${NC}"
echo -e "   docker run -d \\"
echo -e "     --name probability-nginx \\"
echo -e "     --restart unless-stopped \\"
echo -e "     --network app-network \\"
echo -e "     -p 80:80 -p 443:443 \\"
echo -e "     -e DOMAIN=tu-dominio.com \\"
echo -e "     -e SSL_CERT_PATH=/etc/letsencrypt/live/tu-dominio.com/fullchain.pem \\"
echo -e "     -e SSL_KEY_PATH=/etc/letsencrypt/live/tu-dominio.com/privkey.pem \\"
echo -e "     -v /etc/letsencrypt:/etc/letsencrypt:ro \\"
echo -e "     ${ECR_REPO}:${DESCRIPTIVE_TAG}"

echo ""
echo -e "${BLUE}🌐 Repositorio ECR:${NC}"
echo -e "   https://gallery.ecr.aws/c1l9h7c9/probability"
echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}✨ ¡Listo para desplegar en tu servidor ARM64!${NC}" 