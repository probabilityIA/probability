#!/usr/bin/env bash
# Deploy blue-green del frontend en el EC2 de produccion. Lo ejecuta
# ssm-deploy.sh via SSM como usuario ubuntu. Recibe VERSION_FULL por entorno y
# DEPLOY_DIR con los artefactos descargados de S3.
#
# Mismo patron que el backend: se levanta el color nuevo al lado del activo y
# solo se mueve el upstream de nginx cuando el nuevo responde.
set -e

DEST=/home/ubuntu/probability/infra/compose-prod
REPO_URL=476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-frontend
mkdir -p "$DEST"
if [ -d "$DEPLOY_DIR/artifacts" ]; then
  cp -r "$DEPLOY_DIR/artifacts/." "$DEST/"
  echo "Artefactos copiados a $DEST"
fi

. "$DEPLOY_DIR/artifacts/bluegreen-lib.sh"

export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH

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

pull_with_retry() {
  local attempt
  for attempt in 1 2; do
    echo "📦 Pull intento $attempt/2: $REPO_URL:$VERSION_FULL"
    if docker pull "$REPO_URL:$VERSION_FULL"; then
      return 0
    fi
    echo "⚠️  Pull fallo, limpiando cache de Docker..."
    docker image prune -af 2>/dev/null || true
    docker system prune -f 2>/dev/null || true
  done
  echo "❌ Pull fallo. El frontend NO fue modificado."
  return 1
}

pull_with_retry
docker tag "$REPO_URL:$VERSION_FULL" "$REPO_URL:latest"

bg_init
bg_ensure_nginx

ACTIVE=$(bg_active_color front)
TARGET=$(bg_other_color "$ACTIVE")
BACK_ACTIVE=$(bg_active_color back)
NEW_CONTAINER="frontend_prod_$TARGET"
OLD_CONTAINER="frontend_prod_$ACTIVE"

echo "📊 Version: $VERSION_FULL"
echo "🎨 Frontend activo: $ACTIVE -> desplegando en: $TARGET"

docker rm -f "$NEW_CONTAINER" >/dev/null 2>&1 || true

echo "🚀 Levantando front-central-$TARGET..."
docker compose -f docker-compose.yaml --profile green up -d --no-deps "front-central-$TARGET"

if ! bg_wait_healthy "$NEW_CONTAINER" "http://127.0.0.1:3000/" 150; then
  echo "❌ El color nuevo no paso el health check. Se descarta y sigue sirviendo $ACTIVE."
  docker rm -f "$NEW_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

bg_write_upstreams "$BACK_ACTIVE" "$TARGET"
if ! bg_reload_nginx; then
  echo "❌ nginx rechazo la configuracion. Volviendo a $ACTIVE."
  bg_write_upstreams "$BACK_ACTIVE" "$ACTIVE"
  bg_reload_nginx || true
  docker rm -f "$NEW_CONTAINER" >/dev/null 2>&1 || true
  exit 1
fi

echo "✅ Trafico servido por front-central-$TARGET"

bg_retire "$OLD_CONTAINER"
bg_drop_legacy frontend_prod

echo "🧹 Limpiando imagenes antiguas de frontend..."
docker images --format "{{.Repository}}:{{.Tag}} {{.ID}}" | \
  grep "probability-frontend" | \
  grep -v "$VERSION_FULL" | \
  grep -v "latest" | \
  awk '{print $2}' | \
  xargs -r docker rmi -f 2>/dev/null || true
docker image prune -f

echo "✅ Frontend desplegado en $TARGET con la version $VERSION_FULL"
