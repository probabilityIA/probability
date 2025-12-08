#!/bin/zsh

# Script para actualizar servicios de Docker Compose
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🔄 Actualizando servicios de Docker Compose${NC}"

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.yaml" ]; then
  echo -e "${RED}❌ docker-compose.yaml no encontrado. Ejecuta desde el directorio correcto.${NC}"
  exit 1
fi

# Parar servicios
echo -e "${YELLOW}⏹️ Parando servicios...${NC}"
docker-compose down

# Limpiar imágenes no utilizadas
echo -e "${YELLOW}🧹 Limpiando imágenes no utilizadas...${NC}"
docker image prune -f

# Reconstruir y levantar servicios
echo -e "${YELLOW}🔨 Reconstruyendo y levantando servicios...${NC}"
docker-compose up -d --build

# Verificar estado de los servicios
echo -e "${YELLOW}📊 Verificando estado de los servicios...${NC}"
docker-compose ps

echo -e "${GREEN}✅ Servicios actualizados exitosamente${NC}"
echo -e "${YELLOW}🌐 Frontend disponible en: http://localhost/ (puerto 80)${NC}"
echo -e "${YELLOW}🔧 Backend disponible en: http://localhost:3050${NC}"
echo -e "${YELLOW}🗄️ Base de datos disponible en: localhost:5433${NC}"

