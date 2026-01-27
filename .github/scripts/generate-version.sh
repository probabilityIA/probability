#!/bin/bash
set -e

# Argumentos
REPO_NAME=$1  # ej: probability-backend
AWS_REGION=${2:-us-east-1}
AWS_ACCOUNT_ID=${3:-476702565908}

# Calcular año-mes actual
YEAR_MONTH=$(date +%Y.%m)

# Consultar tags existentes en ECR
echo "Consultando versiones existentes en ECR para $REPO_NAME..."
EXISTING_TAGS=$(aws ecr describe-images \
  --repository-name "$REPO_NAME" \
  --region "$AWS_REGION" \
  --query 'imageDetails[].imageTags[]' \
  --output text 2>/dev/null || echo "")

# Encontrar el número de build más alto para el mes actual
MAX_BUILD=0
for tag in $EXISTING_TAGS; do
  # Buscar tags con formato YYYY.MM.N o YYYY.MM.N-SHA
  if [[ $tag =~ ^${YEAR_MONTH}\.([0-9]+) ]]; then
    BUILD_NUM=${BASH_REMATCH[1]}
    if [ "$BUILD_NUM" -gt "$MAX_BUILD" ]; then
      MAX_BUILD=$BUILD_NUM
    fi
  fi
done

# Incrementar al siguiente número de build
NEXT_BUILD=$((MAX_BUILD + 1))

# Generar SHA corto
SHORT_SHA=$(echo "${GITHUB_SHA}" | cut -c1-7)

# Generar versiones
VERSION_FULL="${YEAR_MONTH}.${NEXT_BUILD}-${SHORT_SHA}"
VERSION_SHORT="${YEAR_MONTH}.${NEXT_BUILD}"
VERSION_CALVER="${YEAR_MONTH}.${NEXT_BUILD}"

# Escribir a GITHUB_OUTPUT
echo "VERSION_FULL=${VERSION_FULL}" >> "$GITHUB_OUTPUT"
echo "VERSION_SHORT=${VERSION_SHORT}" >> "$GITHUB_OUTPUT"
echo "VERSION_CALVER=${VERSION_CALVER}" >> "$GITHUB_OUTPUT"
echo "SHORT_SHA=${SHORT_SHA}" >> "$GITHUB_OUTPUT"

# Logging
echo "📦 Versión generada: ${VERSION_FULL}"
echo "   CalVer: ${VERSION_CALVER}"
echo "   Git SHA: ${SHORT_SHA}"
