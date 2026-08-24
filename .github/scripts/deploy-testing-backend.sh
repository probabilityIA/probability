#!/usr/bin/env bash
# Deploy de testing-backend en el EC2 de produccion. Lo ejecuta ssm-deploy.sh via SSM
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

mkdir -p ~/probability/infra/compose-prod

if ! command -v aws &> /dev/null || ! aws --version &> /dev/null; then
  echo "Installing AWS CLI for ARM64..."
  sudo rm -rf /usr/local/aws-cli /usr/local/bin/aws /usr/local/bin/aws_completer
  curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
  unzip -q awscliv2.zip
  sudo ./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli
  rm -rf aws awscliv2.zip
  /usr/local/bin/aws --version || { echo "Error: AWS CLI install failed"; exit 1; }
fi

cd ~/probability/infra/compose-prod || exit 1
export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH

echo "Logging in to ECR..."
/usr/local/bin/aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 476702565908.dkr.ecr.us-east-1.amazonaws.com

REPO_URL="476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-testing-backend"

echo "Pulling image with version tag: $VERSION_FULL"
docker rmi ${REPO_URL}:latest 2>/dev/null || true
docker rmi ${REPO_URL}:$VERSION_FULL 2>/dev/null || true

docker pull ${REPO_URL}:$VERSION_FULL
docker tag ${REPO_URL}:$VERSION_FULL ${REPO_URL}:latest

echo "Stopping testing backend..."
docker stop -t 10 testing_backend_prod || true
docker rm -f testing_backend_prod || true

echo "Liberando puerto 9092..."
sudo fuser -k 9092/tcp || true
sleep 2

echo "Starting testing backend with new image..."
docker compose -f docker-compose.yaml up -d back-testing
sleep 5

echo "Waiting for container to start (15 seconds)..."
sleep 15

CONTAINER_STATUS=$(docker inspect testing_backend_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
if [ "$CONTAINER_STATUS" != "running" ]; then
  echo "Container not running, attempting restart..."
  docker logs --tail 30 testing_backend_prod 2>/dev/null || true
  docker restart testing_backend_prod 2>/dev/null || true
  sleep 10
  CONTAINER_STATUS=$(docker inspect testing_backend_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
  if [ "$CONTAINER_STATUS" != "running" ]; then
    echo "Error: Could not start container"
    exit 1
  fi
fi

echo "Container running correctly"

# Clean old images
docker images --format "{{.Repository}}:{{.Tag}} {{.ID}}" | \
  grep "probability-testing-backend" | \
  grep -v "$VERSION_FULL" | grep -v "latest" | \
  awk '{print $2}' | xargs -r docker rmi -f 2>/dev/null || true
docker image prune -f

echo "Testing backend deployed successfully with version $VERSION_FULL"
