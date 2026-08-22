#!/usr/bin/env bash
# Deploy del backend en el EC2 de produccion. Lo ejecuta ssm-deploy.sh via SSM
# como usuario ubuntu. Recibe VERSION_FULL por entorno y DEPLOY_DIR con los
# artefactos descargados de S3 (docker-compose.yaml, prometheus.yml, grafana/).
set -e

DEST=/home/ubuntu/probability/infra/compose-prod
mkdir -p "$DEST"

echo "📁 Copiando artefactos de deploy"
cp "$DEPLOY_DIR/artifacts/docker-compose.yaml" "$DEST/"
if [ -f "$DEPLOY_DIR/artifacts/prometheus.yml" ]; then
  cp "$DEPLOY_DIR/artifacts/prometheus.yml" "$DEST/"
  [ -d "$DEPLOY_DIR/artifacts/grafana" ] && cp -r "$DEPLOY_DIR/artifacts/grafana" "$DEST/"
  echo "✅ Observability config copiada"
fi

export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH

# Crear directorio si no existe
mkdir -p ~/probability/infra/compose-prod

# Directorio de logs persistentes del backend (sobrevive al redeploy)
# El contenedor corre como appuser (UID 1000): sin este chown no puede escribir
sudo mkdir -p /home/ubuntu/probability/logs/back-central
sudo chown -R 1000:1000 /home/ubuntu/probability/logs/back-central

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

# Obtener el ID de la imagen actual antes del pull
echo "📦 Pulling image with version tag: $VERSION_FULL"
OLD_IMAGE_ID=$(docker images --format "{{.ID}}" 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest 2>/dev/null || echo "")

# Eliminar imágenes locales para forzar nuevo pull
docker rmi 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest 2>/dev/null || true
docker rmi 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:$VERSION_FULL 2>/dev/null || true

# Pull image con versión específica (garantiza que es la imagen recién construida)
# Con reintentos: containerd falla de forma transitoria con
# "unable to lease content: lease does not exist" y sin retry el deploy
# aborta dejando el backend caído.
pull_backend_image() {
  local attempt
  for attempt in 1 2 3; do
    if docker pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:$VERSION_FULL; then
      return 0
    fi
    echo "⚠️  Pull falló (intento $attempt/3), reintentando..."
    docker rmi 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:$VERSION_FULL 2>/dev/null || true
    sleep $((attempt * 5))
  done
  echo "❌ Error: no se pudo descargar la imagen $VERSION_FULL tras 3 intentos"
  echo "   El backend NO fue modificado y sigue corriendo la versión anterior."
  return 1
}
pull_backend_image

# Tag la imagen con versión como latest localmente
docker tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:$VERSION_FULL \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest

# Obtener el ID de la nueva imagen
NEW_IMAGE_ID=$(docker images --format "{{.ID}}" 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest)

echo "📊 Old image ID: ${OLD_IMAGE_ID:-none}"
echo "📊 New image ID: $NEW_IMAGE_ID"
echo "📊 Version: $VERSION_FULL"

# Nueva imagen detectada, actualizar solo backend
echo "✅ Nueva imagen detectada (Version: $VERSION_FULL), actualizando backend..."

# 1. Detener y eliminar solo el contenedor backend
echo "🛑 Deteniendo backend..."
docker stop -t 10 central_reserve_prod || true

echo "🗑️  Eliminando contenedor backend y sus dependientes..."
docker rm -f central_reserve_prod || true

# Verificar y eliminar cualquier contenedor central_reserve_prod que quede
echo "🔍 Verificando eliminación completa..."
BACKEND_IDS=$(docker ps -a --filter name=central_reserve_prod --format '{{.ID}}' 2>/dev/null || true)
if [ -n "$BACKEND_IDS" ]; then
  echo "⚠️  Contenedores backend encontrados, eliminando por ID con dependientes..."
  echo "$BACKEND_IDS" | xargs -r docker rm -f
fi
echo "✅ Contenedor backend eliminado completamente"

# 2. Liberar puerto 3050 (matar proceso conmon/rootlessport)
echo "🔫 Liberando puerto 3050..."
sudo fuser -k 3050/tcp || true
pkill -9 -f "rootlessport.*3050" || true
sleep 2

# 3. Levantar solo el servicio back-central (contenedor: central_reserve_prod)
# La imagen ya fue descargada y tageada como latest antes de tocar el
# contenedor: acá no se hace red, solo arrancar.
echo "🚀 Levantando backend con nueva imagen..."
docker compose -f docker-compose.yaml up -d back-central

# Nota: Nginx detectará automáticamente la reconexión con el backend
# No es necesario recargar nginx ya que usa reintentos automáticos
sleep 3

echo "⏳ Esperando a que el contenedor inicie (20 segundos)..."
sleep 20
  
# Verificar que el contenedor está corriendo (con reintentos)
MAX_RETRIES=3
RETRY_COUNT=0
CONTAINER_STATUS=""

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  CONTAINER_STATUS=$(docker inspect central_reserve_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
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
docker logs --tail 30 central_reserve_prod 2>/dev/null || true

# Verificar estado final
CONTAINER_STATUS=$(docker inspect central_reserve_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
if [ "$CONTAINER_STATUS" != "running" ]; then
  echo "❌ El contenedor no está corriendo después de $MAX_RETRIES intentos"
  echo "📋 Estado final: $CONTAINER_STATUS"
  echo "📋 Revisando logs de error:"
  docker logs --tail 50 central_reserve_prod 2>/dev/null || true

  # Intentar reiniciar manualmente
  echo "🔄 Intentando reiniciar manualmente..."
  docker restart central_reserve_prod 2>/dev/null || docker start central_reserve_prod 2>/dev/null || true
  sleep 10

  # Verificar una vez más
  CONTAINER_STATUS=$(docker inspect central_reserve_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
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
CONTAINER_IMAGE_ID=$(docker inspect central_reserve_prod --format '{{.Image}}' 2>/dev/null || echo "")
EXPECTED_IMAGE_ID=$(docker inspect 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:$VERSION_FULL --format '{{.Id}}' 2>/dev/null || echo "")

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

# Limpiar imágenes antiguas de backend (mantener solo las últimas 2 versiones)
echo "🧹 Limpiando imágenes antiguas de backend..."

# Eliminar todas las imágenes de backend excepto la actual
docker images --format "{{.Repository}}:{{.Tag}} {{.ID}}" | \
  grep "probability-backend" | \
  grep -v "$VERSION_FULL" | \
  grep -v "latest" | \
  awk '{print $2}' | \
  xargs -r docker rmi -f 2>/dev/null || true

# Limpiar imágenes sin usar (dangling)
docker image prune -f

# Reiniciar nginx para refrescar DNS cache
echo "🔄 Reiniciando nginx para refrescar DNS cache..."
docker restart nginx_prod
sleep 3
echo "✅ Nginx reiniciado correctamente"

# Arrancar stack de observabilidad si está definido (idempotente: no reinicia lo que ya corre)
if docker compose -f docker-compose.yaml config --services 2>/dev/null | grep -q "prometheus"; then
  echo "📊 Iniciando stack de observabilidad (cAdvisor + Prometheus + Grafana)..."
  docker compose -f docker-compose.yaml up -d --no-recreate cadvisor prometheus grafana
  echo "✅ Stack de observabilidad iniciado"
fi

echo "✅ Backend deployed successfully with version $VERSION_FULL"
