#!/usr/bin/env bash
# Corre back/migration directamente en el EC2 de produccion via SSM, sin tunel
# a la RDS desde la laptop. El tunel (aws-tunnel.sh) depende de un proceso vivo
# en la maquina local: si la terminal se cierra, se duerme o pierde la sesion
# SSM, el tunel muere a mitad de camino. El EC2 ya tiene ruta de red directa a
# la RDS (por eso el backend corre ahi sin tunel) y ya tiene, en
# infra/compose-prod/.env, las mismas variables que el binario de migracion
# necesita (mismo esquema que back/central/.env). Este script:
#
#   1. Compila back/migration para linux/arm64 (igual arquitectura del EC2).
#   2. Lo sube a S3 y lo ejecuta ahi con AWS-RunShellScript (igual patron que
#      .github/scripts/ssm-deploy.sh usa para los deploys).
#   3. Corre el binario parado en infra/compose-prod, para que tome ese .env.
#   4. Imprime la salida remota y limpia todo (binario temporal + S3).
#
# Uso:
#   ./scripts/run-migration-prod.sh
#
# Antes de correrlo: agregar la migracion nueva a Migrate() en
# back/migration/internal/infra/repository/constructor.go (ver MIGRACIONES.md).
# Despues de correrlo y verificar el efecto: dejar Migrate() en cero otra vez.
set -euo pipefail

export AWS_PROFILE="${AWS_PROFILE:-probability}"
export AWS_REGION="${AWS_REGION:-us-east-1}"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$(mktemp -d /tmp/probability-migration-build.XXXXXX)"
trap 'rm -rf "$BUILD_DIR"' EXIT

echo "Compilando back/migration para linux/arm64..."
mkdir -p "$BUILD_DIR/artifacts"
(
  cd "$ROOT/back/migration"
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o "$BUILD_DIR/artifacts/migration-binary" ./cmd
)

cat > "$BUILD_DIR/deploy.sh" <<'REMOTE'
set -e
DEST=/home/ubuntu/probability/infra/compose-prod
cp "$DEPLOY_DIR/artifacts/migration-binary" "$DEST/migration-tmp-binary"
chmod +x "$DEST/migration-tmp-binary"
cd "$DEST"
set -a
source .env
set +a
# El .env de compose-prod usa DB_PASSWORD/DB_SSLMODE (los nombres que docker-compose
# traduce a DB_PASS/PGSSLMODE dentro del contenedor). El binario nativo los lee
# directo del entorno, asi que hay que mapearlos aca. RELAX_ENV=1 evita que el
# binario aborte por variables que no aplican fuera del backend (USER_PASS_DEFAULT,
# EMAIL_USER_DEFAULT, JWT_SECRET, S3_*, SMTP_*): la migracion no las usa.
export DB_PASS="$DB_PASSWORD"
export PGSSLMODE="$DB_SSLMODE"
export RELAX_ENV=1
STATUS=0
./migration-tmp-binary || STATUS=$?
rm -f "$DEST/migration-tmp-binary"
exit "$STATUS"
REMOTE

echo "Ejecutando en el EC2 de produccion via SSM..."
(
  cd "$BUILD_DIR"
  "$ROOT/.github/scripts/ssm-deploy.sh" deploy.sh
)
