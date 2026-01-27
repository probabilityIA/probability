# Sistema de Versionamiento CalVer

Este documento describe el sistema de versionamiento automático CalVer (Calendar Versioning) implementado para los contenedores Docker/Podman del proyecto Probability.

## Formato de versiones

**Formato**: `YYYY.MM.BUILD-SHA`

Donde:
- `YYYY`: Año (ej: 2026)
- `MM`: Mes (ej: 01 para enero)
- `BUILD`: Número incremental del build del mes (se resetea cada mes)
- `SHA`: Git commit SHA corto (7 caracteres)

### Ejemplos

- `2026.01.1-18c9e87` - Primer build de enero 2026
- `2026.01.2-f214c2d` - Segundo build de enero 2026
- `2026.02.1-abc1234` - Primer build de febrero 2026 (contador reseteado)

## Tags de imagen

Cada imagen desplegada recibe **4 tags** en ECR:

1. `:latest` - Siempre apunta a la versión más reciente
2. `:YYYY.MM.N-SHA` - Versión completa con SHA (ej: `2026.01.5-18c9e87`)
3. `:YYYY.MM.N` - Versión CalVer sin SHA (ej: `2026.01.5`)
4. `:SHA` - Solo el SHA del commit (ej: `18c9e87`)

### Ejemplo en ECR

```
476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest
476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:2026.01.5-18c9e87
476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:2026.01.5
476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:18c9e87
```

## Git Tags automáticos

Cada versión desplegada también crea un tag en git:

- Formato: `v<VERSION_SHORT>`
- Ejemplo: `v2026.01.5`

El tag incluye metadata con:
- Nombre de la imagen completa
- SHA del commit
- Servicio (Backend, Frontend, Website, Nginx)

### Ver tags en git

```bash
# Listar tags recientes
git tag -l "v2026.*" --sort=-version:refname | head -10

# Ver detalles de un tag
git show v2026.01.5
```

## Política de retención en ECR

Las políticas de lifecycle en ECR:

1. **Tag `:latest`**: Protegido, nunca se elimina
2. **Imágenes tagueadas**: Se mantienen las últimas 5
3. **Imágenes sin tag**: Se eliminan después de 7 días

Esto permite:
- Rollback a cualquiera de las últimas ~5 versiones
- Historial de ~1-2 semanas de deploys
- Limpieza automática de imágenes antiguas

## Workflow de despliegue

### 1. Generación de versión

El script `.github/scripts/generate-version.sh`:
1. Calcula `YYYY.MM` de la fecha actual
2. Consulta ECR para ver builds existentes del mes
3. Encuentra el número más alto y lo incrementa
4. Genera las variables de versión

### 2. Build de imagen

Se construye la imagen con los 4 tags simultáneamente:

```yaml
podman build \
  -t ${REPO_URL}:latest \
  -t ${REPO_URL}:${VERSION_FULL} \
  -t ${REPO_URL}:${VERSION_SHORT} \
  -t ${REPO_URL}:${SHORT_SHA} \
  .
```

### 3. Push a ECR

Se suben los 4 tags:

```yaml
podman push ${REPO_URL}:latest
podman push ${REPO_URL}:${VERSION_FULL}
podman push ${REPO_URL}:${VERSION_SHORT}
podman push ${REPO_URL}:${SHORT_SHA}
```

### 4. Creación de Git Tag

Se crea un tag anotado en git:

```yaml
git tag -a "v${VERSION_SHORT}" -m "Release v${VERSION_SHORT} - Service

Image: service:${VERSION_FULL}
Commit: ${GITHUB_SHA}"

git push origin "v${VERSION_SHORT}"
```

### 5. Deployment

El servidor:
1. Hace pull de la imagen con tag `VERSION_FULL` específico
2. La retaguea como `:latest` localmente
3. Podman Compose usa siempre `:latest`

## Servicios versionados

Los siguientes servicios tienen versionamiento CalVer:

1. **Backend** (`probability-backend`)
2. **Frontend** (`probability-frontend`)
3. **Website** (`probability-website`)
4. **Nginx** (`probability-nginx`)

## Workflows afectados

- `.github/workflows/backend-ci-cd.yml`
- `.github/workflows/frontend-ci-cd.yml`
- `.github/workflows/website-ci-cd.yml`
- `.github/workflows/nginx-ci-cd.yml`

## Consultar versiones

### Versiones en ECR

```bash
# Listar todas las versiones de un servicio
aws ecr describe-images \
  --repository-name probability-backend \
  --region us-east-1 \
  --query 'sort_by(imageDetails,&imagePushedAt)[*].[imageTags[0],imagePushedAt]' \
  --output table

# Solo últimas 5 versiones
aws ecr describe-images \
  --repository-name probability-backend \
  --region us-east-1 \
  --query 'sort_by(imageDetails,&imagePushedAt)[-5:].[imageTags[0],imagePushedAt]' \
  --output table
```

### Versiones desplegadas

```bash
# En el servidor
ssh user@server

# Ver imágenes locales
podman images | grep probability

# Ver versión del contenedor corriendo
podman inspect backend_prod --format '{{.Config.Image}}'
```

## Ventajas del sistema CalVer

1. **Trazabilidad**: Cada versión se puede rastrear al commit exacto de git
2. **Rollback fácil**: 5 versiones disponibles en ECR para rollback inmediato
3. **Identificación clara**: El formato `YYYY.MM.N` muestra cuándo fue desplegado
4. **Automatización**: Todo el proceso es automático en cada push a `main`
5. **Múltiples referencias**: 4 tags permiten referenciar la misma imagen de diferentes formas
6. **Git tags**: Correspondencia directa entre tags de imagen y tags de git

## Rollback

Ver la documentación completa de rollback en `docs/ROLLBACK.md`.

Resumen rápido:

```bash
# 1. Listar versiones disponibles
aws ecr describe-images --repository-name probability-backend --region us-east-1

# 2. Pull versión anterior
podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:2026.01.4-abc1234

# 3. Retag como latest
podman tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:2026.01.4-abc1234 \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest

# 4. Recrear contenedor
cd ~/probability/infra/compose-prod
podman stop backend_prod && podman rm backend_prod
podman-compose -f podman-compose.yaml up -d --no-deps back-central
```

## Terraform

Las políticas de lifecycle están definidas en:

```
infra/terraform/ecr.tf
```

Para aplicar cambios:

```bash
cd infra/terraform
terraform plan
terraform apply
```

## Mantenimiento

### Ajustar retención

Para cambiar el número de imágenes retenidas, edita `infra/terraform/ecr.tf`:

```terraform
{
  rulePriority = 2
  description  = "Mantener últimas X imágenes tagueadas"
  selection = {
    tagStatus   = "tagged"
    countType   = "imageCountMoreThan"
    countNumber = 5  # Cambiar este número
  }
  action = {
    type = "expire"
  }
}
```

### Limpiar imágenes antiguas manualmente

```bash
# En el servidor, limpiar imágenes locales no usadas
podman image prune -f

# Ver espacio liberado
podman system df
```

## Costos

- **ECR Storage**: ~$0.10/GB/mes
- **5 imágenes** (~200MB cada una) = ~$0.10/mes
- **Costo total estimado**: ~$1/mes para los 4 servicios

## Troubleshooting

### El script de versiones falla

```bash
# Verificar permisos de AWS
aws ecr describe-images --repository-name probability-backend --region us-east-1

# Verificar que GITHUB_OUTPUT esté definido
echo $GITHUB_OUTPUT
```

### Tags duplicados

Si por alguna razón se crea un tag duplicado, ECR simplemente reasignará el tag a la nueva imagen. Los tags no son únicos por imagen, sino por nombre.

### Build numbers incorrectos

El script consulta ECR para encontrar el número más alto del mes. Si ECR está vacío o no tiene imágenes del mes actual, empezará desde 1.

## Referencias

- CalVer Specification: https://calver.org/
- AWS ECR Lifecycle Policies: https://docs.aws.amazon.com/AmazonECR/latest/userguide/LifecyclePolicies.html
- Podman tagging: https://docs.podman.io/en/latest/markdown/podman-tag.1.html
