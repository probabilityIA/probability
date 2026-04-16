#!/bin/bash

# Script para limpiar imágenes Docker/Podman antiguas
# Uso: ./cleanup_images.sh [--force]
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
NC='\033[0m'

FORCE=false

# Procesar argumentos
while [[ $# -gt 0 ]]; do
  case $1 in
    --force|-f)
      FORCE=true
      shift
      ;;
    *)
      echo -e "${RED}❌ Argumento desconocido: $1${NC}"
      echo "Uso: $0 [--force]"
      exit 1
      ;;
  esac
done

echo -e "${MAGENTA}🧹 Script de Limpieza de Imágenes Docker/Podman${NC}"
echo ""

# Detectar si usa Docker o Podman
if command -v docker &> /dev/null; then
    CONTAINER_CMD="docker"
    echo -e "${BLUE}🐳 Usando Docker${NC}"
elif command -v podman &> /dev/null; then
    CONTAINER_CMD="podman"
    echo -e "${BLUE}🦭 Usando Podman${NC}"
else
    echo -e "${RED}❌ No se encontró Docker ni Podman${NC}"
    exit 1
fi

echo ""

# Mostrar espacio ANTES
echo -e "${MAGENTA}💾 Espacio en disco ANTES de limpiar:${NC}"
df -h / | grep -E "Filesystem|/$"
echo ""

# Mostrar imágenes actuales
echo -e "${MAGENTA}📦 Imágenes actuales:${NC}"
$CONTAINER_CMD images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | head -15
echo ""

# Contar imágenes antes
IMAGES_BEFORE=$($CONTAINER_CMD images -q | wc -l)
CONTAINERS_RUNNING=$($CONTAINER_CMD ps -q | wc -l)

echo -e "${BLUE}📊 Estadísticas actuales:${NC}"
echo -e "   • Imágenes totales: ${IMAGES_BEFORE}"
echo -e "   • Contenedores corriendo: ${CONTAINERS_RUNNING}"
echo ""

# Confirmar antes de limpiar (si no se usa --force)
if [ "$FORCE" = false ]; then
  echo -e "${YELLOW}⚠️  Esta operación eliminará:${NC}"
  echo "   • Todas las imágenes no utilizadas por contenedores activos"
  echo "   • Imágenes dangling (sin tag)"
  echo "   • Imágenes antiguas con tags que ya no se usan"
  echo ""
  read -p "¿Continuar con la limpieza? (y/N): " -n 1 -r
  echo ""
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}❌ Operación cancelada${NC}"
    exit 0
  fi
  echo ""
fi

# Limpieza paso a paso
echo -e "${GREEN}🧹 Iniciando limpieza...${NC}"
echo ""

# 1. Limpiar contenedores detenidos
echo -e "${YELLOW}🗑️  Paso 1: Eliminando contenedores detenidos...${NC}"
STOPPED_CONTAINERS=$($CONTAINER_CMD container prune -f 2>&1 | grep -oP '(?<=Total reclaimed space: ).*' || echo "0B")
echo -e "${GREEN}  ✅ Espacio liberado: ${STOPPED_CONTAINERS}${NC}"
echo ""

# 2. Limpiar imágenes dangling
echo -e "${YELLOW}🗑️  Paso 2: Eliminando imágenes dangling (<none>)...${NC}"
DANGLING_SPACE=$($CONTAINER_CMD image prune -f 2>&1 | grep -oP '(?<=Total reclaimed space: ).*' || echo "0B")
echo -e "${GREEN}  ✅ Espacio liberado: ${DANGLING_SPACE}${NC}"
echo ""

# 3. Limpiar TODAS las imágenes no utilizadas
echo -e "${YELLOW}🗑️  Paso 3: Eliminando TODAS las imágenes no utilizadas...${NC}"
ALL_IMAGES_SPACE=$($CONTAINER_CMD image prune -a -f 2>&1 | grep -oP '(?<=Total reclaimed space: ).*' || echo "0B")
echo -e "${GREEN}  ✅ Espacio liberado: ${ALL_IMAGES_SPACE}${NC}"
echo ""

# 4. Limpiar volúmenes no utilizados (opcional)
echo -e "${YELLOW}🗑️  Paso 4: Eliminando volúmenes huérfanos...${NC}"
VOLUMES_SPACE=$($CONTAINER_CMD volume prune -f 2>&1 | grep -oP '(?<=Total reclaimed space: ).*' || echo "0B")
echo -e "${GREEN}  ✅ Espacio liberado: ${VOLUMES_SPACE}${NC}"
echo ""

# 5. Limpiar redes no utilizadas
echo -e "${YELLOW}🗑️  Paso 5: Eliminando redes no utilizadas...${NC}"
$CONTAINER_CMD network prune -f > /dev/null 2>&1
echo -e "${GREEN}  ✅ Redes limpiadas${NC}"
echo ""

# 6. Limpiar build cache (solo Docker)
if [ "$CONTAINER_CMD" = "docker" ]; then
  echo -e "${YELLOW}🗑️  Paso 6: Limpiando build cache...${NC}"
  BUILD_CACHE=$($CONTAINER_CMD builder prune -f 2>&1 | grep -oP '(?<=Total reclaimed space: ).*' || echo "0B")
  echo -e "${GREEN}  ✅ Espacio liberado: ${BUILD_CACHE}${NC}"
  echo ""
fi

# Contar imágenes después
IMAGES_AFTER=$($CONTAINER_CMD images -q | wc -l)
IMAGES_REMOVED=$((IMAGES_BEFORE - IMAGES_AFTER))

echo -e "${GREEN}✅ Limpieza completada exitosamente${NC}"
echo ""

echo -e "${BLUE}📊 Resumen:${NC}"
echo -e "   • Imágenes eliminadas: ${IMAGES_REMOVED}"
echo -e "   • Imágenes restantes: ${IMAGES_AFTER}"
echo -e "   • Contenedores corriendo: ${CONTAINERS_RUNNING} (sin cambios)"
echo ""

# Mostrar espacio DESPUÉS
echo -e "${MAGENTA}💾 Espacio en disco DESPUÉS de limpiar:${NC}"
df -h / | grep -E "Filesystem|/$"
echo ""

# Mostrar imágenes restantes
echo -e "${MAGENTA}📦 Imágenes restantes:${NC}"
$CONTAINER_CMD images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
echo ""

echo -e "${GREEN}🎉 ¡Limpieza completada!${NC}"
echo -e "${MAGENTA}💡 Tip: Ejecuta este script periódicamente para mantener el espacio limpio${NC}"
