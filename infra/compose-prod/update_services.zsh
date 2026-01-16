#!/bin/zsh

# Script para actualizar servicios de Podman Compose desde ECR
# Actualiza solo: font-central, font-website y back-central
# NO actualiza: redis, rabbitmq, nginx
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}🔄 Actualizando servicios de aplicación desde ECR${NC}"

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.yaml" ]; then
  echo -e "${RED}❌ docker-compose.yaml no encontrado. Ejecuta desde el directorio correcto.${NC}"
  exit 1
fi

# Servicios a actualizar (solo frontends y backend)
SERVICES_TO_UPDATE=("font-central" "font-website" "back-central")

# Descargar las imágenes más recientes desde ECR (solo para los servicios especificados)
echo -e "${BLUE}📥 Descargando imágenes más recientes desde ECR...${NC}"
for service in "${SERVICES_TO_UPDATE[@]}"; do
  echo -e "${YELLOW}  → Descargando imagen para: ${service}${NC}"
  podman-compose pull "$service"
done

# Limpiar imágenes antiguas no utilizadas (dangling images) ANTES de actualizar
echo -e "${YELLOW}🧹 Limpiando imágenes antiguas no utilizadas...${NC}"
podman image prune -f

# Usar podman-compose up -d para actualizar los servicios
# Esto maneja mejor las dependencias y el orden de inicio
echo -e "${YELLOW}🚀 Actualizando servicios con las nuevas imágenes...${NC}"
echo -e "${YELLOW}   (esto recreará los contenedores con las nuevas imágenes)${NC}"
podman-compose up -d --no-deps "${SERVICES_TO_UPDATE[@]}"

# Esperar a que los servicios estén listos
echo -e "${YELLOW}⏳ Esperando a que los servicios estén listos...${NC}"
sleep 5

# Verificar estado de los servicios
echo -e "${YELLOW}📊 Verificando estado de los servicios...${NC}"
podman-compose ps

# Verificar que los servicios estén saludables
echo -e "${YELLOW}🏥 Verificando salud de los servicios...${NC}"
for service in "${SERVICES_TO_UPDATE[@]}"; do
  if podman-compose ps "$service" | grep -q "Up"; then
    echo -e "${GREEN}  ✅ $service está corriendo${NC}"
  else
    echo -e "${RED}  ❌ $service NO está corriendo${NC}"
  fi
done

# Recargar nginx para que detecte los servicios actualizados
echo -e "${YELLOW}🔄 Recargando nginx para detectar servicios actualizados...${NC}"
if podman-compose exec -T nginx nginx -s reload > /dev/null 2>&1; then
  echo -e "${GREEN}  ✅ Nginx recargado exitosamente${NC}"
else
  echo -e "${YELLOW}  ⚠️  Nginx no pudo recargarse (puede que no esté corriendo o no sea necesario)${NC}"
fi

echo -e "${GREEN}✅ Servicios actualizados exitosamente${NC}"
echo -e "${BLUE}📋 Servicios actualizados:${NC}"
echo -e "   • font-central (Frontend Central)"
echo -e "   • font-website (Frontend Website)"
echo -e "   • back-central (Backend Central)"
echo -e "${BLUE}📋 Servicios NO modificados (siguen corriendo):${NC}"
echo -e "   • redis (Redis Cache)"
echo -e "   • rabbitmq (RabbitMQ Message Queue)"
echo -e "${YELLOW}🌐 Frontend disponible en: http://localhost/ (puerto 80)${NC}"
echo -e "${YELLOW}🔧 Backend disponible en: http://localhost:3050${NC}"
