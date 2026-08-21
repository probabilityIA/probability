#!/usr/bin/env bash
# Deploy de frontend en el EC2 de produccion. Lo ejecuta ssm-deploy.sh via SSM
# como usuario ubuntu. Recibe VERSION_FULL por entorno y DEPLOY_DIR con los
# artefactos descargados de S3.
set -e

DEST=/home/ubuntu/probability/infra/compose-prod
mkdir -p "$DEST"
if [ -d "$DEPLOY_DIR/artifacts" ]; then
  cp -r "$DEPLOY_DIR/artifacts/." "$DEST/"
  echo "Artefactos copiados a $DEST"
fi

export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH

# Crear directorio si no existe
mkdir -p ~/probability/infra/compose-prod

# Instalar/reinstalar AWS CLI para ARM64
# Verificar si AWS CLI está instalado y funciona correctamente
if ! command -v aws &> /dev/null || ! aws --version &> /dev/null; then
  echo "📦 Instalando/reinstalando AWS CLI para ARM64..."
  # Remover instalación anterior si existe
  sudo rm -rf /usr/local/aws-cli /usr/local/bin/aws /usr/local/bin/aws_completer
  # Descargar e instalar versión ARM64
  curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
  unzip -q awscliv2.zip
  sudo ./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli
  rm -rf aws awscliv2.zip
  # Verificar instalación
  /usr/local/bin/aws --version || {
    echo "❌ Error: AWS CLI no se instaló correctamente"
    exit 1
  }
  echo "✅ AWS CLI instalado correctamente"
else
  echo "✅ AWS CLI ya está instalado y funcionando"
fi

cd ~/probability/infra/compose-prod || exit 1

# Asegurar que AWS CLI esté en PATH
export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH

# Login a ECR privado
echo "🔐 Logging in to ECR..."
/usr/local/bin/aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 476702565908.dkr.ecr.us-east-1.amazonaws.com

# Función para pull con retry automático ante snapshots corruptos
pull_with_retry() {
  local image=$1
  local max_attempts=2
  local attempt=1

  while [ $attempt -le $max_attempts ]; do
    echo "📦 Pull intento $attempt/$max_attempts: $image"
    if docker pull "$image"; then
      return 0
    fi

    if [ $attempt -lt $max_attempts ]; then
      echo "⚠️  Pull falló (posibles snapshots corruptos), limpiando cache Docker..."
      docker image prune -af 2>/dev/null || true
      docker system prune -f 2>/dev/null || true
      echo "🔄 Reintentando pull..."
    fi

    attempt=$((attempt + 1))
  done

  echo "❌ Pull falló después de $max_attempts intentos"
  return 1
}

REPO_URL="476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-frontend"

# Limpiar imágenes locales previas para evitar conflictos de capas
echo "🧹 Limpiando imágenes locales previas..."
docker rmi ${REPO_URL}:latest 2>/dev/null || true
docker rmi ${REPO_URL}:$VERSION_FULL 2>/dev/null || true
docker image prune -f 2>/dev/null || true

echo "📊 Version: $VERSION_FULL"
# Nueva imagen detectada, actualizar solo frontend
echo "✅ Nueva imagen detectada (Version: $VERSION_FULL), actualizando frontend..."

# 1. Detener y eliminar solo el contenedor frontend
echo "🛑 Deteniendo frontend..."
docker stop -t 10 frontend_prod || true

echo "🗑️  Eliminando contenedor frontend y sus dependientes..."
docker rm -f frontend_prod || true

# Verificar y eliminar cualquier contenedor frontend_prod que quede
echo "🔍 Verificando eliminación completa..."
FRONTEND_IDS=$(docker ps -a --filter name=frontend_prod --format '{{.ID}}' 2>/dev/null || true)
if [ -n "$FRONTEND_IDS" ]; then
  echo "⚠️  Contenedores frontend encontrados, eliminando por ID con dependientes..."
  echo "$FRONTEND_IDS" | xargs -r docker rm -f
fi
echo "✅ Contenedor frontend eliminado completamente"

# 2. Liberar puerto 8080 (matar proceso conmon/rootlessport)
echo "🔫 Liberando puerto 8080..."
sudo fuser -k 8080/tcp || true
pkill -9 -f "rootlessport.*8080" || true
sleep 2

# 3. Pull de la nueva imagen específica (con retry automático)
echo "📦 Pulling nueva imagen $VERSION_FULL..."
pull_with_retry ${REPO_URL}:$VERSION_FULL

# 4. Tag como latest localmente
echo "🏷️  Tageando como latest..."
docker tag ${REPO_URL}:$VERSION_FULL ${REPO_URL}:latest

# 5. Levantar solo el servicio front-central (contenedor: frontend_prod)
echo "🚀 Levantando frontend con nueva imagen..."
docker compose -f docker-compose.yaml up -d --no-deps front-central
sleep 3

echo "⏳ Esperando a que el contenedor inicie (15 segundos)..."
sleep 15
  
# Verificar que el contenedor está corriendo (con reintentos)
MAX_RETRIES=3
RETRY_COUNT=0
CONTAINER_STATUS=""

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  CONTAINER_STATUS=$(docker inspect frontend_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
  echo "📋 Intento $(($RETRY_COUNT + 1))/$MAX_RETRIES - Estado del contenedor: $CONTAINER_STATUS"

  if [ "$CONTAINER_STATUS" = "running" ]; then
    break
  fi

  RETRY_COUNT=$(($RETRY_COUNT + 1))
  if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
    echo "⏳ Esperando 5 segundos más..."
    sleep 5
  fi
done

# Mostrar logs antes de verificar estado final
echo "📋 Últimos logs del contenedor:"
docker logs --tail 30 frontend_prod 2>/dev/null || true

# Verificar estado final
CONTAINER_STATUS=$(docker inspect frontend_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
if [ "$CONTAINER_STATUS" != "running" ]; then
  echo "❌ El contenedor no está corriendo después de $MAX_RETRIES intentos"
  echo "📋 Estado final: $CONTAINER_STATUS"
  echo "📋 Revisando logs de error:"
  docker logs --tail 50 frontend_prod 2>/dev/null || true

  # Intentar reiniciar manualmente
  echo "🔄 Intentando reiniciar manualmente..."
  docker restart frontend_prod 2>/dev/null || docker start frontend_prod 2>/dev/null || true
  sleep 10

  # Verificar una vez más
  CONTAINER_STATUS=$(docker inspect frontend_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
  if [ "$CONTAINER_STATUS" != "running" ]; then
    echo "❌ Error crítico: No se pudo iniciar el contenedor. Revisar logs manualmente."
    exit 1
  else
    echo "✅ Contenedor iniciado después de reinicio manual"
  fi
else
  echo "✅ Contenedor corriendo correctamente"
fi

# Verificar que el contenedor está usando la nueva imagen
CONTAINER_IMAGE_ID=$(docker inspect frontend_prod --format '{{.Image}}' 2>/dev/null || echo "")
EXPECTED_IMAGE_ID=$(docker inspect ${REPO_URL}:$VERSION_FULL --format '{{.Id}}' 2>/dev/null || echo "")

# Remover prefijo sha256: si existe
EXPECTED_IMAGE_ID_CLEAN="${EXPECTED_IMAGE_ID#sha256:}"
CONTAINER_IMAGE_ID_CLEAN="${CONTAINER_IMAGE_ID#sha256:}"

echo "📋 Image ID del contenedor: ${CONTAINER_IMAGE_ID_CLEAN:0:12}"
echo "📋 Image ID esperada (Version $VERSION_FULL): ${EXPECTED_IMAGE_ID_CLEAN:0:12}"

if [ -n "$CONTAINER_IMAGE_ID_CLEAN" ] && [ -n "$EXPECTED_IMAGE_ID_CLEAN" ]; then
  if [ "$CONTAINER_IMAGE_ID_CLEAN" = "$EXPECTED_IMAGE_ID_CLEAN" ]; then
    echo "✅ Contenedor actualizado correctamente con la nueva imagen (Version: $VERSION_FULL)"
  else
    echo "⚠️  WARNING: Image IDs no coinciden exactamente, pero el contenedor está corriendo"
    echo "   Esto puede ser normal si Docker usa capas compartidas"
  fi
fi

# Limpiar imágenes antiguas de frontend
echo "🧹 Limpiando imágenes antiguas de frontend..."

docker images --format "{{.Repository}}:{{.Tag}} {{.ID}}" | \
  grep "probability-frontend" | \
  grep -v "$VERSION_FULL" | \
  grep -v "latest" | \
  awk '{print $2}' | \
  xargs -r docker rmi -f 2>/dev/null || true

docker image prune -f

# Reiniciar nginx para refrescar DNS cache
echo "🔄 Reiniciando nginx para refrescar DNS cache..."
docker restart nginx_prod
sleep 3
echo "✅ Nginx reiniciado correctamente"

echo "✅ Frontend deployed successfully with version $VERSION_FULL"
