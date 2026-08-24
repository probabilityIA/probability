#!/usr/bin/env bash
# Deploy de testing-frontend en el EC2 de produccion. Lo ejecuta ssm-deploy.sh via SSM
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

REPO_URL="476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-testing-frontend"

echo "Cleaning local images..."
docker rmi ${REPO_URL}:latest 2>/dev/null || true
docker rmi ${REPO_URL}:$VERSION_FULL 2>/dev/null || true
docker image prune -f 2>/dev/null || true

echo "Pulling image $VERSION_FULL..."
docker pull ${REPO_URL}:$VERSION_FULL
docker tag ${REPO_URL}:$VERSION_FULL ${REPO_URL}:latest

echo "Stopping testing frontend..."
docker stop -t 10 testing_frontend_prod || true
docker rm -f testing_frontend_prod || true

echo "Liberando puerto 8082..."
sudo fuser -k 8082/tcp || true
sleep 2

echo "Starting testing frontend with new image..."
docker compose -f docker-compose.yaml up -d --no-deps front-testing
sleep 3

echo "Waiting for container to start (60 seconds max)..."
WAIT_COUNT=0
MAX_WAIT=12
while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
  CONTAINER_STATUS=$(docker inspect testing_frontend_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
  if [ "$CONTAINER_STATUS" = "running" ]; then
    echo "Container is running after $((WAIT_COUNT * 5 + 3)) seconds"
    break
  fi
  if [ "$CONTAINER_STATUS" = "exited" ] || [ "$CONTAINER_STATUS" = "dead" ]; then
    echo "Container exited prematurely (status: $CONTAINER_STATUS)"
    echo "=== Container logs (stdout + stderr) ==="
    docker logs --tail 50 testing_frontend_prod 2>&1 || true
    echo "=== Container inspect ==="
    docker inspect testing_frontend_prod --format '{{.State.ExitCode}} {{.State.Error}}' 2>&1 || true
    echo "=== OOMKilled check ==="
    docker inspect testing_frontend_prod --format '{{.State.OOMKilled}}' 2>&1 || true
    echo "Attempting restart..."
    docker restart testing_frontend_prod 2>/dev/null || true
    sleep 15
    CONTAINER_STATUS=$(docker inspect testing_frontend_prod --format '{{.State.Status}}' 2>/dev/null || echo "not_found")
    if [ "$CONTAINER_STATUS" != "running" ]; then
      echo "=== Logs after restart ==="
      docker logs --tail 30 testing_frontend_prod 2>&1 || true
      echo "Error: Could not start container"
      exit 1
    fi
    break
  fi
  WAIT_COUNT=$((WAIT_COUNT + 1))
  sleep 5
done

if [ "$CONTAINER_STATUS" != "running" ]; then
  echo "Container did not start within 60 seconds (status: $CONTAINER_STATUS)"
  echo "=== Container logs ==="
  docker logs --tail 50 testing_frontend_prod 2>&1 || true
  exit 1
fi

echo "Container running correctly"

# Clean old images
docker images --format "{{.Repository}}:{{.Tag}} {{.ID}}" | \
  grep "probability-testing-frontend" | \
  grep -v "$VERSION_FULL" | grep -v "latest" | \
  awk '{print $2}' | xargs -r docker rmi -f 2>/dev/null || true
docker image prune -f

# Restart nginx to pick up new upstream
echo "Restarting nginx..."
docker restart nginx_prod
sleep 3

echo "Testing frontend deployed successfully with version $VERSION_FULL"
