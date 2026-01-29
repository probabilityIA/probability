# Estructura de Workflows CI/CD

## Patrón Común

Todos los workflows siguen la misma estructura de 2 jobs:

```yaml
jobs:
  build-and-push:
    # Build de imagen y push a ECR
    runs-on: ubuntu-24.04-arm
    outputs:
      version_full: ${{ steps.version.outputs.VERSION_FULL }}
      version_short: ${{ steps.version.outputs.VERSION_SHORT }}

  deploy:
    # Deploy al servidor EC2
    needs: build-and-push
    runs-on: ubuntu-latest
```

## Job 1: Build and Push

### Steps

1. **Checkout code**
   ```yaml
   - uses: actions/checkout@v4
   ```

2. **Setup del lenguaje** (Node.js / Go)
   ```yaml
   - uses: actions/setup-node@v4  # Para frontend/website
   # O
   - uses: actions/setup-go@v5    # Para backend
   ```

3. **Install Podman**
   ```yaml
   - run: |
       sudo apt-get update
       sudo apt-get install -y podman
   ```

4. **Configure AWS credentials**
   ```yaml
   - uses: aws-actions/configure-aws-credentials@v4
     with:
       aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
       aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
       aws-region: us-east-1
   ```

5. **Login to ECR**
   ```bash
   aws ecr get-login-password | \
     podman login --username AWS --password-stdin \
     476702565908.dkr.ecr.us-east-1.amazonaws.com
   ```

6. **Generate version**
   ```bash
   .github/scripts/generate-version.sh probability-<service> us-east-1 476702565908
   ```

   Genera:
   - `VERSION_FULL`: `2026.01.37.a3d554c` (año.día.contador.sha)
   - `VERSION_SHORT`: `2026.01.37`
   - `SHORT_SHA`: `a3d554c`

7. **Build image with Podman**
   ```bash
   podman build \
     --platform linux/arm64 \
     --no-cache \
     -t ${REPO_URL}:latest \
     -t ${REPO_URL}:${VERSION_FULL} \
     -t ${REPO_URL}:${VERSION_SHORT} \
     -t ${REPO_URL}:${SHORT_SHA} \
     .
   ```

8. **Push image to ECR**
   ```bash
   podman push ${REPO_URL}:latest
   podman push ${REPO_URL}:${VERSION_FULL}
   podman push ${REPO_URL}:${VERSION_SHORT}
   podman push ${REPO_URL}:${SHORT_SHA}
   ```

9. **Create Git Tag**
   ```bash
   git tag -a "v${VERSION_SHORT}" -m "Release v${VERSION_SHORT}"
   git push origin "v${VERSION_SHORT}" --force
   ```

## Job 2: Deploy

### Steps

1. **Checkout code**
   ```yaml
   - uses: actions/checkout@v4
   ```

2. **Setup SSH**
   ```yaml
   - uses: webfactory/ssh-agent@v0.9.0
     with:
       ssh-private-key: ${{ secrets.EC2_SSH_KEY }}
   ```

3. **Upload compose file**
   ```bash
   scp infra/compose-prod/podman-compose.yaml \
     ubuntu@server:~/probability/infra/compose-prod/
   ```

4. **Deploy to server**

   Script ejecutado via SSH con heredoc:

   ```bash
   ssh ubuntu@server "VERSION_FULL='$VERSION_FULL' bash -s" << 'EOF'
     # 1. Setup
     export PATH=/usr/local/bin:/usr/bin:/bin:$HOME/.local/bin:$PATH
     cd ~/probability/infra/compose-prod

     # 2. Login to ECR
     aws ecr get-login-password | podman login ...

     # 3. Pull y tag de nueva imagen
     OLD_IMAGE_ID=$(podman images --format "{{.ID}}" .../:latest)
     podman rmi .../:latest
     podman rmi .../:$VERSION_FULL
     podman pull .../:$VERSION_FULL
     podman tag .../:$VERSION_FULL .../:latest
     NEW_IMAGE_ID=$(podman images --format "{{.ID}}" .../:latest)

     # 4. Detener y eliminar contenedor viejo
     podman stop -t 10 <service>_prod
     podman rm -f --depend <service>_prod

     # 5. Verificar eliminación completa
     IDS=$(podman ps -a --filter name=<service>_prod --format '{{.ID}}')
     if [ -n "$IDS" ]; then
       echo "$IDS" | xargs -r podman rm -f --depend
     fi

     # 6. Liberar puerto
     sudo fuser -k <PORT>/tcp
     pkill -9 -f "rootlessport.*<PORT>"
     sleep 2

     # 7. Pull nueva imagen (otra vez, para asegurar)
     podman pull .../:$VERSION_FULL
     podman tag .../:$VERSION_FULL .../:latest

     # 8. Levantar servicio
     podman-compose -f podman-compose.yaml up -d <service>

     # 9. Esperar inicio
     sleep 15

     # 10. Verificar estado con reintentos
     MAX_RETRIES=3
     RETRY_COUNT=0
     while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
       CONTAINER_STATUS=$(podman inspect <service>_prod --format '{{.State.Status}}')
       if [ "$CONTAINER_STATUS" = "running" ]; then
         break
       fi
       RETRY_COUNT=$((RETRY_COUNT + 1))
       sleep 5
     done

     # 11. Verificar que usa la imagen correcta
     CONTAINER_IMAGE_ID=$(podman inspect <service>_prod --format '{{.Image}}')
     # Comparar IDs (primeros 12 caracteres)

     # 12. Limpiar imágenes antiguas
     podman images | grep "probability-<service>" | \
       grep -v "$VERSION_FULL" | grep -v "latest" | \
       awk '{print $2}' | xargs -r podman rmi -f

     podman image prune -f
   EOF
   ```

## Script de Versioning

**`.github/scripts/generate-version.sh`**

```bash
#!/bin/bash
set -e

REPO_NAME=$1
AWS_REGION=$2
AWS_ACCOUNT_ID=$3

# Obtener último contador del día desde ECR
YEAR=$(date +%Y)
DAY_OF_YEAR=$(date +%j)
DATE_PREFIX="${YEAR}.${DAY_OF_YEAR}"

LATEST_TAG=$(aws ecr describe-images \
  --repository-name "$REPO_NAME" \
  --region "$AWS_REGION" \
  --query "sort_by(imageDetails,& imagePushedAt)[-1].imageTags[0]" \
  --output text 2>/dev/null || echo "")

if [[ "$LATEST_TAG" =~ ^${DATE_PREFIX}\.([0-9]+)\. ]]; then
  COUNTER=$((${BASH_REMATCH[1]} + 1))
else
  COUNTER=1
fi

SHORT_SHA=$(git rev-parse --short=7 HEAD)

VERSION_FULL="${DATE_PREFIX}.${COUNTER}.${SHORT_SHA}"
VERSION_SHORT="${DATE_PREFIX}.${COUNTER}"

# Output para GitHub Actions
echo "VERSION_FULL=$VERSION_FULL" >> $GITHUB_OUTPUT
echo "VERSION_SHORT=$VERSION_SHORT" >> $GITHUB_OUTPUT
echo "SHORT_SHA=$SHORT_SHA" >> $GITHUB_OUTPUT
```

**Formato de versión:**
- `YYYY.DDD.N.XXXXXXX`
  - `YYYY` - Año (2026)
  - `DDD` - Día del año (001-365)
  - `N` - Contador del día (incrementa por cada deploy)
  - `XXXXXXX` - Short SHA del commit

**Ejemplos:**
- `2026.029.1.a3d554c` - Primer deploy del día 29 de 2026
- `2026.029.2.f8b2c91` - Segundo deploy del mismo día
- `2026.030.1.c4e9a12` - Primer deploy del día siguiente

## Concurrency Control

Agregado a todos los workflows:

```yaml
concurrency:
  group: <service>-deploy  # backend-deploy, frontend-deploy, etc.
  cancel-in-progress: true
```

**Comportamiento:**
- Solo un workflow del mismo tipo puede ejecutarse a la vez
- Si llega un nuevo push mientras uno está corriendo, el viejo se cancela
- Previene race conditions y conflictos de puertos

## Triggers

### Backend
```yaml
on:
  push:
    branches: [main]
    paths:
      - 'back/central/**'
      - 'back/migration/**'
      - '.github/workflows/backend-ci-cd.yml'
  workflow_dispatch:  # Manual trigger
```

### Frontend
```yaml
on:
  push:
    branches: [main]
    paths:
      - 'front/central/**'
      - '.github/workflows/frontend-ci-cd.yml'
  workflow_dispatch:
```

### Website
```yaml
on:
  push:
    branches: [main]
    paths:
      - 'front/website/**'
      - '.github/workflows/website-ci-cd.yml'
  workflow_dispatch:
```

### Nginx
```yaml
on:
  push:
    branches: [main]
    paths:
      - 'infra/nginx/**'
      - '.github/workflows/nginx-ci-cd.yml'
  workflow_dispatch:
```

## Tiempos de Ejecución

| Workflow | Build | Deploy | Total |
|----------|-------|--------|-------|
| Backend  | ~2min | ~1min  | ~3min |
| Frontend | ~2min | ~1min  | ~3min |
| Website  | ~1min | ~1min  | ~2min |
| Nginx    | ~30s  | ~1min  | ~2min |

## Optimizaciones Aplicadas

### 1. ARM64 Native
- Runners: `ubuntu-24.04-arm`
- Build más rápido (nativo, sin emulación)
- Imágenes más pequeñas

### 2. Cache de Dependencias
```yaml
- uses: actions/setup-node@v4
  with:
    node-version: '20'
    cache: 'npm'  # ← Cache automático
    cache-dependency-path: front/central/package-lock.json
```

### 3. No-Cache Build
```bash
podman build --no-cache ...
```
- Asegura que siempre se usa la última versión de dependencias
- Previene bugs por cache stale

### 4. Specific Version Pull
```bash
podman pull .../:$VERSION_FULL  # No :latest
```
- Garantiza que se descarga la imagen correcta
- Previene race conditions con tags

## Secretos Requeridos

### AWS
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

### SSH
- `EC2_SSH_KEY` - Private key (formato PEM)
- `EC2_USER` - Usuario SSH (ubuntu)
- `EC2_HOST` - IP o dominio del servidor

### Build Args (Frontend)
- `NEXT_PUBLIC_API_BASE_URL` - URL pública: `https://app.probabilityia.com.co/api/v1`
- `API_BASE_URL` - URL interna: `http://back-central:3050/api/v1`

## Logs y Debugging

### Ver logs de un workflow
```bash
gh run view <run-id> --log
```

### Ver logs de un job específico
```bash
gh run view <run-id> --log --job <job-id>
```

### Download logs
```bash
gh run download <run-id>
```

### Ver workflows en progreso
```bash
gh run watch
```

## Mejoras Futuras

1. **Cache de imágenes Docker:** Usar GitHub Container Registry cache
2. **Matrix builds:** Build para múltiples arquitecturas (ARM64 + AMD64)
3. **Deploy strategies:** Blue-green, canary deployments
4. **Rollback automático:** Si health checks fallan, rollback a versión anterior
5. **Notificaciones:** Slack/Discord para deployments
6. **Métricas:** Tiempo de deploy, éxito/fallo rate
