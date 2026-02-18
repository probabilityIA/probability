#!/bin/zsh

# Script para actualizar servicios de Docker Compose desde ECR
# Actualiza solo: font-central, font-website y back-central
# NO actualiza: redis, rabbitmq, nginx
# INCLUYE: Limpieza automática de imágenes antiguas
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
NC='\033[0m'

echo -e "${GREEN}🔄 Actualizando servicios de aplicación desde ECR${NC}"

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.yaml" ]; then
  echo -e "${RED}❌ docker-compose.yaml no encontrado. Ejecuta desde el directorio correcto.${NC}"
  exit 1
fi

# Servicios a actualizar (solo frontends y backend)
SERVICES_TO_UPDATE=("font-central" "font-website" "back-central")

# Mostrar espacio en disco ANTES de la actualización
echo -e "${MAGENTA}💾 Espacio en disco ANTES de actualizar:${NC}"
df -h / | grep -E "Filesystem|/$"
echo ""

# Mostrar imágenes actuales y su tamaño
echo -e "${MAGENTA}📦 Imágenes Docker actuales:${NC}"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | head -10
echo ""

# Descargar las imágenes más recientes desde ECR (solo para los servicios especificados)
echo -e "${BLUE}📥 Descargando imágenes más recientes desde ECR...${NC}"
for service in "${SERVICES_TO_UPDATE[@]}"; do
  echo -e "${YELLOW}  → Descargando imagen para: ${service}${NC}"
  docker compose pull "$service"
done

echo ""
echo -e "${YELLOW}🚀 Actualizando servicios con las nuevas imágenes...${NC}"
echo -e "${YELLOW}   (esto recreará los contenedores con las nuevas imágenes)${NC}"
docker compose up -d --no-deps "${SERVICES_TO_UPDATE[@]}"

# Esperar a que los servicios estén listos
echo -e "${YELLOW}⏳ Esperando a que los servicios estén listos...${NC}"
sleep 5

# Verificar estado de los servicios
echo -e "${YELLOW}📊 Verificando estado de los servicios...${NC}"
docker compose ps

# Verificar que los servicios estén saludables
echo -e "${YELLOW}🏥 Verificando salud de los servicios...${NC}"
ALL_SERVICES_OK=true
for service in "${SERVICES_TO_UPDATE[@]}"; do
  if docker compose ps "$service" | grep -q "Up"; then
    echo -e "${GREEN}  ✅ $service está corriendo${NC}"
  else
    echo -e "${RED}  ❌ $service NO está corriendo${NC}"
    ALL_SERVICES_OK=false
  fi
done

if [ "$ALL_SERVICES_OK" = true ]; then
  echo ""
  echo -e "${GREEN}✅ Todos los servicios están corriendo correctamente${NC}"
  echo -e "${GREEN}🧹 Procediendo con la limpieza de imágenes antiguas...${NC}"
  echo ""

  # LIMPIEZA AGRESIVA: Eliminar TODAS las imágenes no utilizadas
  # Esto incluye imágenes antiguas con tags que ya no están en uso
  echo -e "${YELLOW}🗑️  Eliminando imágenes antiguas no utilizadas (incluyendo las con tags)...${NC}"
  echo -e "${YELLOW}   Esta operación puede tardar unos segundos...${NC}"

  # Contar imágenes antes de limpiar
  IMAGES_BEFORE=$(docker images -q | wc -l)

  # Eliminar imágenes no utilizadas (sin preguntar)
  # -a: elimina TODAS las imágenes no usadas, no solo dangling
  # -f: forzar sin confirmación
  docker image prune -a -f

  # Contar imágenes después de limpiar
  IMAGES_AFTER=$(docker images -q | wc -l)
  IMAGES_REMOVED=$((IMAGES_BEFORE - IMAGES_AFTER))

  echo -e "${GREEN}  ✅ Limpieza completada${NC}"
  echo -e "${BLUE}  📊 Imágenes eliminadas: ${IMAGES_REMOVED}${NC}"
  echo -e "${BLUE}  📊 Imágenes restantes: ${IMAGES_AFTER}${NC}"
  echo ""

  # Opcional: Limpiar también volúmenes huérfanos (descomentala si quieres)
  # echo -e "${YELLOW}🗑️  Limpiando volúmenes huérfanos...${NC}"
  # docker volume prune -f

  # Mostrar espacio liberado
  echo -e "${MAGENTA}💾 Espacio en disco DESPUÉS de limpiar:${NC}"
  df -h / | grep -E "Filesystem|/$"
  echo ""

  # Mostrar imágenes restantes
  echo -e "${MAGENTA}📦 Imágenes Docker restantes:${NC}"
  docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
  echo ""

else
  echo ""
  echo -e "${RED}⚠️  ADVERTENCIA: Algunos servicios no están corriendo correctamente${NC}"
  echo -e "${YELLOW}⚠️  Se omitió la limpieza de imágenes para evitar problemas${NC}"
  echo -e "${YELLOW}💡 Revisa los logs con: docker compose logs <nombre-servicio>${NC}"
  exit 1
fi

# Recargar nginx para que detecte los servicios actualizados
echo -e "${YELLOW}🔄 Recargando nginx para detectar servicios actualizados...${NC}"
if docker compose exec -T nginx nginx -s reload > /dev/null 2>&1; then
  echo -e "${GREEN}  ✅ Nginx recargado exitosamente${NC}"
else
  echo -e "${YELLOW}  ⚠️  Nginx no pudo recargarse (puede que no esté corriendo o no sea necesario)${NC}"
fi

echo ""
echo -e "${GREEN}🎉 ¡Actualización completada exitosamente!${NC}"
echo ""
echo -e "${BLUE}📋 Servicios actualizados:${NC}"
echo -e "   • font-central (Frontend Central)"
echo -e "   • font-website (Frontend Website)"
echo -e "   • back-central (Backend Central)"
echo ""
echo -e "${BLUE}📋 Servicios NO modificados (siguen corriendo):${NC}"
echo -e "   • redis (Redis Cache)"
echo -e "   • rabbitmq (RabbitMQ Message Queue)"
echo ""
echo -e "${YELLOW}🌐 URLs de la aplicación:${NC}"
echo -e "${GREEN}   ✓ Frontend: https://www.probabilityia.com.co${NC}"
echo -e "${GREEN}   ✓ Backend:  https://www.probabilityia.com.co/api/v1${NC}"
echo -e "${GREEN}   ✓ Swagger:  https://www.probabilityia.com.co/swagger/${NC}"
echo ""
echo -e "${MAGENTA}💡 Tip: Las imágenes antiguas se eliminaron automáticamente para ahorrar espacio${NC}"
