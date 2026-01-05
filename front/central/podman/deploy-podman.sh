#!/bin/bash

# Script de despliegue para ECR público usando Podman
# Probability - Frontend Central
# Migración de Docker a Podman

set -e

# Variables
IMAGE_NAME="probability-front-central"
# Mismo repositorio que el backend, diferentes etiquetas
ECR_REPO="public.ecr.aws/c1l9h7c9/probability"
VERSION=${1:-"latest"}
DOCKERFILE_PATH="docker/Dockerfile"
AWS_PROFILE_NAME="probability"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Iniciando despliegue de Probability Frontend Central (Podman)${NC}"
echo -e "${YELLOW}Versión: ${VERSION}${NC}"
echo -e "${YELLOW}Perfil de AWS: ${AWS_PROFILE_NAME}${NC}"

# Verificar que estamos en el directorio correcto
# Este script debe ejecutarse desde front/central/
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ Error: No se encontró package.json. Ejecuta desde front/central/${NC}"
    exit 1
fi

# Verificar que Podman esté instalado y corriendo
if ! command -v podman > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Podman no está instalado${NC}"
    echo -e "${YELLOW}💡 Instala Podman: https://podman.io/getting-started/installation${NC}"
    exit 1
fi

if ! podman info > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Podman no está corriendo o no está configurado correctamente${NC}"
    exit 1
fi

# Verificar y configurar emulación para ARM64 si es necesario
ARCH=$(uname -m)
if [ "$ARCH" != "aarch64" ] && [ "$ARCH" != "arm64" ]; then
    echo -e "${YELLOW}⚠️  Sistema x86_64 detectado, verificando emulación para ARM64...${NC}"
    
    # Verificar si binfmt_misc ya está configurado para qemu-aarch64
    if [ ! -f /proc/sys/fs/binfmt_misc/qemu-aarch64 ] 2>/dev/null; then
        echo -e "${YELLOW}📦 Configurando emulación QEMU para ARM64 (requiere sudo)...${NC}"
        
        # Configurar emulación usando sudo
        if sudo podman run --rm --privileged multiarch/qemu-user-static --reset -p yes > /dev/null 2>&1; then
            echo -e "${GREEN}   ✅ Emulación QEMU configurada${NC}"
        else
            echo -e "${RED}   ❌ Error configurando emulación${NC}"
            echo -e "${YELLOW}   Asegúrate de tener privilegios sudo y que Podman esté instalado${NC}"
            exit 1
        fi
    else
        echo -e "${GREEN}   ✅ Emulación QEMU ya configurada${NC}"
    fi
fi

# Verificar que AWS CLI esté configurado con el perfil correcto
if ! aws --profile "${AWS_PROFILE_NAME}" sts get-caller-identity > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: AWS CLI no está configurado correctamente${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Verificaciones completadas${NC}"

# Verificar que el lockfile existe
# El Dockerfile usará --frozen-lockfile para garantizar versiones exactas (React 19.2.3, styled-jsx 5.1.7)
echo -e "${YELLOW}📦 Verificando lockfile de dependencias...${NC}"
if [ -f "pnpm-lock.yaml" ]; then
    echo -e "${GREEN}   ✅ pnpm-lock.yaml encontrado${NC}"
    echo -e "${BLUE}   Podman usará --frozen-lockfile para respetar versiones exactas${NC}"
elif [ -f "package-lock.json" ]; then
    echo -e "${GREEN}   ✅ package-lock.json encontrado${NC}"
    echo -e "${BLUE}   Podman usará npm ci --legacy-peer-deps para respetar versiones exactas${NC}"
else
    echo -e "${YELLOW}   ⚠️  No se encontró lockfile. Podman generará uno nuevo durante el build${NC}"
fi

# URLs del API
# NEXT_PUBLIC_API_BASE_URL = Cliente (SSE, dominio público)
# API_BASE_URL = Servidor (Server Actions, red interna Podman)
PUBLIC_API_URL=${NEXT_PUBLIC_API_BASE_URL:-"https://app.probabilityia.com.co/api/v1"}
SERVER_API_URL=${API_BASE_URL:-"http://back-central:3050/api/v1"}

echo -e "${BLUE}🌐 URLs del API:${NC}"
echo -e "   Cliente (SSE):  ${PUBLIC_API_URL}"
echo -e "   Servidor (Actions): ${SERVER_API_URL}"
echo ""

# Construir la imagen para ARM64
echo -e "${YELLOW}🔨 Construyendo imagen Podman para ARM64...${NC}"
if [ "$ARCH" != "aarch64" ] && [ "$ARCH" != "arm64" ]; then
    echo -e "${BLUE}   Esto puede tomar varios minutos (usando emulación QEMU)...${NC}"
else
    echo -e "${BLUE}   Esto puede tomar varios minutos...${NC}"
fi

podman build \
    --platform linux/arm64 \
    --build-arg NEXT_PUBLIC_API_BASE_URL=${PUBLIC_API_URL} \
    --build-arg API_BASE_URL=${SERVER_API_URL} \
    -f ${DOCKERFILE_PATH} \
    -t ${IMAGE_NAME}:${VERSION} \
    .

echo -e "${GREEN}✅ Imagen construida exitosamente${NC}"

# Etiquetar para ECR con nombres más descriptivos
echo -e "${YELLOW}🏷️  Etiquetando imagen para ECR...${NC}"

# Crear tags descriptivos
if [ "${VERSION}" = "latest" ]; then
    # Para latest, crear múltiples tags descriptivos
    TIMESTAMP=$(date +%Y%m%d-%H%M%S)
    DESCRIPTIVE_TAG="frontend-latest"
    DATED_TAG="frontend-${TIMESTAMP}"
    
    podman tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DESCRIPTIVE_TAG}
    podman tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DATED_TAG}
    
    echo -e "${GREEN}📅 Tags creados: ${DESCRIPTIVE_TAG}, ${DATED_TAG}${NC}"
else
    # Para versiones específicas, crear tag descriptivo
    DESCRIPTIVE_TAG="frontend-${VERSION}"
    
    podman tag ${IMAGE_NAME}:${VERSION} ${ECR_REPO}:${DESCRIPTIVE_TAG}
    
    echo -e "${GREEN}🏷️  Tags creados: ${DESCRIPTIVE_TAG}${NC}"
fi

# Login a ECR público
echo -e "${YELLOW}🔐 Haciendo login a ECR público con el perfil '${AWS_PROFILE_NAME}'...${NC}"
aws --profile "${AWS_PROFILE_NAME}" ecr-public get-login-password --region us-east-1 | podman login --username AWS --password-stdin public.ecr.aws

# Push de las imágenes
echo -e "${YELLOW}⬆️  Subiendo imágenes a ECR...${NC}"
echo -e "${BLUE}   Esto puede tomar varios minutos dependiendo de tu conexión...${NC}"

if [ "${VERSION}" = "latest" ]; then
    # Subir todos los tags para latest
    podman push ${ECR_REPO}:${DESCRIPTIVE_TAG}
    podman push ${ECR_REPO}:${DATED_TAG}
    echo -e "${GREEN}✅ Imágenes subidas con tags: ${DESCRIPTIVE_TAG}, ${DATED_TAG}${NC}"
else
    # Subir tags para versiones específicas
    podman push ${ECR_REPO}:${DESCRIPTIVE_TAG}
    echo -e "${GREEN}✅ Imagen subida con tag: ${DESCRIPTIVE_TAG}${NC}"
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
echo -e "   podman run -d \\"
echo -e "     --name probability-front-central \\"
echo -e "     --restart unless-stopped \\"
echo -e "     --network app-network \\"
echo -e "     -p 8080:80 \\"
echo -e "     ${ECR_REPO}:${DESCRIPTIVE_TAG}"

echo ""
echo -e "${BLUE}📝 Configuración de la imagen:${NC}"
echo -e "   • Puerto interno:     80"
echo -e "   • Puerto expuesto:    8080"
echo -e "   • Cliente (SSE):      ${PUBLIC_API_URL}"
echo -e "   • Servidor (Actions): ${SERVER_API_URL}"

echo ""
echo -e "${BLUE}🌐 Repositorio ECR:${NC}"
echo -e "   https://gallery.ecr.aws/c1l9h7c9/probability"
echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}✨ ¡Listo para desplegar en tu servidor ARM64!${NC}"
