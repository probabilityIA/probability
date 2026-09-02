#!/usr/bin/env bash
# Deploy blue-green del backend en el EC2 de produccion. Lo ejecuta
# ssm-deploy.sh via SSM como usuario ubuntu. Recibe VERSION_FULL por entorno y
# DEPLOY_DIR con los artefactos descargados de S3.
#
# El color nuevo se levanta al lado del que esta sirviendo. Solo cuando responde
# /health se mueve el upstream de nginx y se recarga. Si no responde, se borra el
# contenedor nuevo y el viejo sigue atendiendo: el deploy falla sin caida.
set -e

DEST=/home/ubuntu/probability/infra/compose-prod
REPO_URL=476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend
mkdir -p "$DEST"

echo "📁 Copiando artefactos de deploy"
cp "$DEPLOY_DIR/artifacts/docker-compose.yaml" "$DEST/"
if [ -f "$DEPLOY_DIR/artifacts/prometheus.yml" ]; then
  cp "$DEPLOY_DIR/artifacts/prometheus.yml" "$DEST/"
  [ -d "$DEPLOY_DIR/artifacts/grafana" ] && cp -r "$DEPLOY_DIR/artifacts/grafana" "$DEST/"
  echo "✅ Observability config copiada"
fi

. "$DEPLOY_DIR/artifacts/bluegreen-lib.sh"

export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH

# Directorio de logs persistentes del backend (sobrevive al redeploy)
# El contenedor corre como appuser (UID 1000): sin este chown no puede escribir
sudo mkdir -p /home/ubuntu/probability/logs/back-central
sudo chown -R 1000:1000 /home/ubuntu/probability/logs/back-central

if ! command -v aws &> /dev/null || ! aws --version &> /dev/null; then
  echo "📦 Instalando AWS CLI para ARM64..."
  sudo rm -rf /usr/local/aws-cli /usr/local/bin/aws /usr/local/bin/aws_completer
  curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
  unzip -q awscliv2.zip
  sudo ./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli
  rm -rf aws awscliv2.zip
  /usr/local/bin/aws --version || { echo "❌ AWS CLI no se instalo"; exit 1; }
fi

cd "$DEST" || exit 1

echo "🔐 Login a ECR..."
/usr/local/bin/aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 476702565908.dkr.ecr.us-east-1.amazonaws.com

# Pull con reintentos: containerd falla de forma transitoria con
# "unable to lease content: lease does not exist". Si falla, no se toco nada
# todavia y la version en produccion sigue intacta.
pull_backend_image() {
  local attempt
  for attempt in 1 2 3; do
    if docker pull "$REPO_URL:$VERSION_FULL"; then
      return 0
    fi
    echo "⚠️  Pull fallo (intento $attempt/3), reintentando..."
    docker rmi "$REPO_URL:$VERSION_FULL" 2>/dev/null || true
    sleep $((attempt * 5))
  done
  echo "❌ No se pudo descargar la imagen $VERSION_FULL tras 3 intentos"
  echo "   El backend NO fue modificado y sigue corriendo la version anterior."
  return 1
}

echo "📦 Pulling $REPO_URL:$VERSION_FULL"
pull_backend_image
docker tag "$REPO_URL:$VERSION_FULL" "$REPO_URL:latest"

bg_init
bg_ensure_nginx

ACTIVE=$(bg_active_color back)
TARGET=$(bg_other_color "$ACTIVE")
FRONT_ACTIVE=$(bg_active_color front)
NEW_CONTAINER="central_reserve_prod_$TARGET"
OLD_CONTAINER="central_reserve_prod_$ACTIVE"

echo "📊 Version: $VERSION_FULL"
echo "🎨 Backend activo: $ACTIVE -> desplegando en: $TARGET"

# Restos de un intento anterior que haya quedado a medias
docker rm -f "$NEW_CONTAINER" >/dev/null 2>&1 || true

echo "🚀 Levantando back-central-$TARGET..."
docker compose -f docker-compose.yaml --profile green up -d --no-deps "back-central-$TARGET"

if ! bg_wait_healthy "$NEW_CONTAINER" "http://localhost:3050/health" 150; then
  echo "❌ El color nuevo no paso el health check. Se descarta y sigue sirviendo $ACTIVE."
  docker rm -f "$NEW_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

bg_write_upstreams "$TARGET" "$FRONT_ACTIVE"
if ! bg_reload_nginx; then
  echo "❌ nginx rechazo la configuracion. Volviendo a $ACTIVE."
  bg_write_upstreams "$ACTIVE" "$FRONT_ACTIVE"
  bg_reload_nginx || true
  docker rm -f "$NEW_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

echo "✅ Trafico servido por back-central-$TARGET"

bg_retire "$OLD_CONTAINER"
bg_drop_legacy central_reserve_prod

echo "🧹 Limpiando imagenes antiguas de backend..."
docker images --format "{{.Repository}}:{{.Tag}} {{.ID}}" | \
  grep "probability-backend" | \
  grep -v "$VERSION_FULL" | \
  grep -v "latest" | \
  awk '{print $2}' | \
  xargs -r docker rmi -f 2>/dev/null || true
docker image prune -f

if docker compose -f docker-compose.yaml config --services 2>/dev/null | grep -q "prometheus"; then
  echo "📊 Iniciando stack de observabilidad..."
  docker compose -f docker-compose.yaml up -d --no-recreate cadvisor prometheus grafana
fi

echo "✅ Backend desplegado en $TARGET con la version $VERSION_FULL"
