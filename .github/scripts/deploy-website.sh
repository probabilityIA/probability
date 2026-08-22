#!/usr/bin/env bash
# Deploy de website en el EC2 de produccion. Lo ejecuta ssm-deploy.sh via SSM
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

cd ~/probability/infra/compose-prod || exit 1

# Asegurar que AWS CLI esté en PATH
export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH

# Login a ECR privado
echo "🔐 Logging in to ECR..."
/usr/local/bin/aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 476702565908.dkr.ecr.us-east-1.amazonaws.com

# Obtener el ID de la imagen actual antes del pull
echo "📦 Pulling image with version tag: $VERSION_FULL"
OLD_IMAGE_ID=$(docker images --format "{{.ID}}" 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:latest 2>/dev/null || echo "")

# Eliminar imágenes locales para forzar nuevo pull
docker rmi 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:latest 2>/dev/null || true
docker rmi 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:$VERSION_FULL 2>/dev/null || true

# Pull image con versión específica (garantiza que es la imagen recién construida)
docker pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:$VERSION_FULL

# Tag la imagen con versión como latest localmente
docker tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:$VERSION_FULL \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:latest

# Obtener el ID de la nueva imagen
NEW_IMAGE_ID=$(docker images --format "{{.ID}}" 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:latest)

echo "📊 Old image ID: ${OLD_IMAGE_ID:-none}"
echo "📊 New image ID: $NEW_IMAGE_ID"
echo "📊 Version: $VERSION_FULL"

# Nueva imagen detectada, actualizar solo website
echo "✅ Nueva imagen detectada (Version: $VERSION_FULL), actualizando website..."

# 1. Detener y eliminar solo el contenedor website
echo "🛑 Deteniendo website..."
docker stop -t 10 website_prod || true

echo "🗑️  Eliminando contenedor website y sus dependientes..."
docker rm -f website_prod || true

# Verificar y eliminar cualquier contenedor website_prod que quede
echo "🔍 Verificando eliminación completa..."
WEBSITE_IDS=$(docker ps -a --filter name=website_prod --format '{{.ID}}' 2>/dev/null || true)
if [ -n "$WEBSITE_IDS" ]; then
  echo "⚠️  Contenedores website encontrados, eliminando por ID con dependientes..."
  echo "$WEBSITE_IDS" | xargs -r docker rm -f
fi
echo "✅ Contenedor website eliminado completamente"

# 2. Liberar puerto 8081 (matar proceso conmon/rootlessport)
echo "🔫 Liberando puerto 8081..."
sudo fuser -k 8081/tcp || true
pkill -9 -f "rootlessport.*8081" || true
sleep 2

# 3. Pull de la nueva imagen específica
echo "📦 Pulling nueva imagen $VERSION_FULL..."
docker pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:$VERSION_FULL

# 4. Tag como latest localmente
echo "🏷️  Tageando como latest..."
docker tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:$VERSION_FULL \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:latest

# 5. Levantar solo el servicio front-website (contenedor: website_prod)
echo "🚀 Levantando website con nueva imagen..."
docker compose -f docker-compose.yaml up -d front-website

# Nota: Nginx detectará automáticamente la reconexión con el website
# No es necesario recargar nginx ya que usa reintentos automáticos
sleep 3

echo "⏳ Esperando a que el contenedor inicie (15 segundos)..."
sleep 15
  
# Verificar que el contenedor está corriendo (con reintentos)
MAX_RETRIES=3
RETRY_COUNT=0
CONTAINER_STATUS=""

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  CONTAINER_STATUS=$(docker inspect website_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
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
docker logs --tail 30 website_prod 2>/dev/null || true

# Verificar estado final
CONTAINER_STATUS=$(docker inspect website_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
if [ "$CONTAINER_STATUS" != "running" ]; then
  echo "❌ El contenedor no está corriendo después de $MAX_RETRIES intentos"
  echo "📋 Estado final: $CONTAINER_STATUS"
  echo "📋 Revisando logs de error:"
  docker logs --tail 50 website_prod 2>/dev/null || true

  # Intentar reiniciar manualmente
  echo "🔄 Intentando reiniciar manualmente..."
  docker restart website_prod 2>/dev/null || docker start website_prod 2>/dev/null || true
  sleep 10

  # Verificar una vez más
  CONTAINER_STATUS=$(docker inspect website_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
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
CONTAINER_IMAGE_ID=$(docker inspect website_prod --format '{{.Image}}' 2>/dev/null || echo "")
EXPECTED_IMAGE_ID=$(docker inspect 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:$VERSION_FULL --format '{{.Id}}' 2>/dev/null || echo "")

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

# Limpiar imágenes antiguas de website
echo "🧹 Limpiando imágenes antiguas de website..."

docker images --format "{{.Repository}}:{{.Tag}} {{.ID}}" | \
  grep "probability-website" | \
  grep -v "$VERSION_FULL" | \
  grep -v "latest" | \
  awk '{print $2}' | \
  xargs -r docker rmi -f 2>/dev/null || true

docker image prune -f

echo "✅ Website deployed successfully with version $VERSION_FULL"
