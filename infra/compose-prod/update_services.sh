#!/bin/bash

# Script para actualizar servicios de Podman Compose desde ECR
# Actualiza solo: front-central, front-website y back-central
# NO actualiza: redis, rabbitmq, nginx

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}Actualizando servicios de aplicacion desde ECR${NC}"

# Verificar que estamos en el directorio correcto
if [ ! -f "podman-compose.yaml" ]; then
  echo -e "${RED}podman-compose.yaml no encontrado. Ejecuta desde el directorio correcto.${NC}"
  exit 1
fi

# Verificar que Podman este disponible
if ! command -v podman &> /dev/null; then
  echo -e "${RED}Podman no esta instalado${NC}"
  exit 1
fi

if ! command -v podman-compose &> /dev/null; then
  echo -e "${RED}podman-compose no esta instalado${NC}"
  exit 1
fi

# Login a ECR publico
echo -e "${YELLOW}Login a ECR publico...${NC}"
aws ecr-public get-login-password --region us-east-1 --profile probability | \
  podman login --username AWS --password-stdin public.ecr.aws

# Servicios a actualizar (solo frontends y backend)
SERVICES_TO_UPDATE=("front-central" "front-website" "back-central")

# Descargar las imagenes mas recientes desde ECR
echo -e "${BLUE}Descargando imagenes mas recientes desde ECR...${NC}"
for service in "${SERVICES_TO_UPDATE[@]}"; do
  echo -e "${YELLOW}  -> Descargando imagen para: ${service}${NC}"
  podman-compose -f podman-compose.yaml pull "$service" || true
done

# Limpiar imagenes antiguas no utilizadas ANTES de actualizar
echo -e "${YELLOW}Limpiando imagenes antiguas no utilizadas...${NC}"
podman image prune -f

# Actualizar los servicios con las nuevas imagenes
echo -e "${YELLOW}Actualizando servicios con las nuevas imagenes...${NC}"
for service in "${SERVICES_TO_UPDATE[@]}"; do
  echo -e "${YELLOW}  -> Recreando: ${service}${NC}"
  podman-compose -f podman-compose.yaml up -d --force-recreate "$service" || true
done

# Esperar a que los servicios esten listos
echo -e "${YELLOW}Esperando a que los servicios esten listos...${NC}"
sleep 5

# Verificar estado de los servicios
echo -e "${YELLOW}Verificando estado de los servicios...${NC}"
podman-compose -f podman-compose.yaml ps

# Verificar que los servicios esten corriendo
echo -e "${YELLOW}Verificando salud de los servicios...${NC}"
for service in "${SERVICES_TO_UPDATE[@]}"; do
  container_name=$(podman-compose -f podman-compose.yaml ps -q "$service" 2>/dev/null)
  if [ -n "$container_name" ]; then
    status=$(podman inspect --format '{{.State.Status}}' "$container_name" 2>/dev/null || echo "unknown")
    if [ "$status" = "running" ]; then
      echo -e "${GREEN}  $service esta corriendo${NC}"
    else
      echo -e "${RED}  $service NO esta corriendo (status: $status)${NC}"
    fi
  else
    echo -e "${RED}  $service NO encontrado${NC}"
  fi
done

# Recargar nginx para que detecte los servicios actualizados
echo -e "${YELLOW}Recargando nginx...${NC}"
if podman exec nginx_prod nginx -s reload > /dev/null 2>&1; then
  echo -e "${GREEN}  Nginx recargado exitosamente${NC}"
else
  echo -e "${YELLOW}  Nginx no pudo recargarse (puede que no este corriendo)${NC}"
fi

echo -e "${GREEN}Servicios actualizados exitosamente${NC}"
echo ""
echo -e "${BLUE}Servicios actualizados:${NC}"
echo "   - front-central (Frontend Central)"
echo "   - front-website (Frontend Website)"
echo "   - back-central (Backend Central)"
echo ""
echo -e "${BLUE}Servicios NO modificados (siguen corriendo):${NC}"
echo "   - redis (Redis Cache)"
echo "   - rabbitmq (RabbitMQ Message Queue)"
echo "   - nginx (Reverse Proxy)"
echo ""
echo -e "${YELLOW}Frontend disponible en: https://app.probabilityia.com.co${NC}"
echo -e "${YELLOW}Backend disponible en: http://localhost:3050${NC}"
