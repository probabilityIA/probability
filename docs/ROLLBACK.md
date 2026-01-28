# Procedimiento de Rollback - CalVer Versioning

Este documento describe el procedimiento para hacer rollback a una versión anterior de los contenedores en producción.

## Contexto

Con el sistema de versionamiento CalVer implementado, cada despliegue genera 4 tags por imagen:
- `:latest` - Siempre apunta a la versión más reciente
- `:YYYY.MM.N-SHA` - Versión completa (ej: `2026.01.5-18c9e87`)
- `:YYYY.MM.N` - Versión corta CalVer (ej: `2026.01.5`)
- `:SHA` - Git commit SHA corto (ej: `18c9e87`)

ECR retiene las últimas 5 imágenes para permitir rollback.

## Procedimiento de Rollback

### 1. Conectar al servidor

```bash
ssh user@server
```

### 2. Login a ECR

```bash
aws ecr get-login-password --region us-east-1 | \
  podman login --username AWS --password-stdin 476702565908.dkr.ecr.us-east-1.amazonaws.com
```

### 3. Listar versiones disponibles

Para ver las últimas versiones disponibles en ECR:

```bash
# Backend
aws ecr describe-images \
  --repository-name probability-backend \
  --region us-east-1 \
  --query 'sort_by(imageDetails,&imagePushedAt)[-5:].[imageTags[0],imagePushedAt]' \
  --output table

# Frontend
aws ecr describe-images \
  --repository-name probability-frontend \
  --region us-east-1 \
  --query 'sort_by(imageDetails,&imagePushedAt)[-5:].[imageTags[0],imagePushedAt]' \
  --output table

# Website
aws ecr describe-images \
  --repository-name probability-website \
  --region us-east-1 \
  --query 'sort_by(imageDetails,&imagePushedAt)[-5:].[imageTags[0],imagePushedAt]' \
  --output table

# Nginx
aws ecr describe-images \
  --repository-name probability-nginx \
  --region us-east-1 \
  --query 'sort_by(imageDetails,&imagePushedAt)[-5:].[imageTags[0],imagePushedAt]' \
  --output table
```

### 4. Hacer rollback a versión específica

Elige el servicio que necesitas revertir:

#### Backend

```bash
# Establecer la versión a la que quieres hacer rollback
VERSION_ROLLBACK="2026.01.3-abc1234"  # Cambiar por la versión deseada

# Pull de la versión específica
podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:${VERSION_ROLLBACK}

# Retag como latest localmente
podman tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:${VERSION_ROLLBACK} \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest

# Navegar al directorio de compose
cd ~/probability/infra/compose-prod

# Recrear el contenedor
podman stop backend_prod
podman rm backend_prod
podman-compose -f podman-compose.yaml up -d --no-deps back-central

# Reiniciar nginx para reconectar
podman restart nginx_prod

# Verificar logs
podman logs --tail 50 backend_prod
```

#### Frontend

```bash
VERSION_ROLLBACK="2026.01.3-abc1234"  # Cambiar por la versión deseada

podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-frontend:${VERSION_ROLLBACK}
podman tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-frontend:${VERSION_ROLLBACK} \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-frontend:latest

cd ~/probability/infra/compose-prod
podman stop frontend_prod
podman rm frontend_prod
podman-compose -f podman-compose.yaml up -d --no-deps front-central
podman restart nginx_prod

podman logs --tail 50 frontend_prod
```

#### Website

```bash
VERSION_ROLLBACK="2026.01.3-abc1234"  # Cambiar por la versión deseada

podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:${VERSION_ROLLBACK}
podman tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:${VERSION_ROLLBACK} \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-website:latest

cd ~/probability/infra/compose-prod
podman stop website_prod
podman rm website_prod
podman-compose -f podman-compose.yaml up -d --no-deps front-website
podman restart nginx_prod

podman logs --tail 50 website_prod
```

#### Nginx

```bash
VERSION_ROLLBACK="2026.01.3-abc1234"  # Cambiar por la versión deseada

podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-nginx:${VERSION_ROLLBACK}
podman tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-nginx:${VERSION_ROLLBACK} \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-nginx:latest

cd ~/probability/infra/compose-prod
podman-compose -f podman-compose.yaml up -d --force-recreate --no-deps nginx

podman logs --tail 50 nginx_prod
```

### 5. Verificar el rollback

```bash
# Verificar que el contenedor está corriendo
podman ps | grep <service>_prod

# Verificar la imagen que está usando
podman inspect <service>_prod --format '{{.Image}}'

# Verificar logs en tiempo real
podman logs -f <service>_prod
```

## Verificación de versiones desplegadas

Para ver qué versión está corriendo actualmente en producción:

```bash
# Ver todas las imágenes locales con sus tags
podman images | grep probability

# Ver imagen específica del contenedor
podman inspect backend_prod --format '{{.Config.Image}}'
podman inspect frontend_prod --format '{{.Config.Image}}'
podman inspect website_prod --format '{{.Config.Image}}'
podman inspect nginx_prod --format '{{.Config.Image}}'
```

## Git Tags

Cada versión desplegada también tiene un tag en git. Para ver los tags:

```bash
# Listar tags recientes
git tag -l "v2026.*" --sort=-version:refname | head -10

# Ver detalles de un tag específico
git show v2026.01.5
```

## Notas importantes

1. **Retención**: ECR retiene las últimas 5 imágenes. Si necesitas hacer rollback a una versión más antigua, verifica que aún esté disponible en ECR.

2. **Tag latest**: El tag `:latest` siempre está protegido en ECR y nunca se elimina automáticamente.

3. **Rollback de emergencia**: Si necesitas hacer rollback urgente y tienes problemas con AWS CLI, puedes usar el SHA corto directamente:
   ```bash
   podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:<SHA>
   ```

4. **Múltiples servicios**: Si necesitas hacer rollback de múltiples servicios a la vez, hazlo uno por uno y verifica que cada uno funciona antes de continuar con el siguiente.

5. **Nginx**: Siempre reinicia nginx después de hacer rollback de backend, frontend o website para que resuelva las nuevas IPs de los contenedores.

## Ejemplo completo de rollback

Escenario: El backend deployment v2026.01.5 tiene un bug crítico y necesitas volver a v2026.01.4.

```bash
# 1. SSH al servidor
ssh user@server

# 2. Login a ECR
aws ecr get-login-password --region us-east-1 | \
  podman login --username AWS --password-stdin 476702565908.dkr.ecr.us-east-1.amazonaws.com

# 3. Listar versiones para confirmar que v2026.01.4 existe
aws ecr describe-images \
  --repository-name probability-backend \
  --region us-east-1 \
  --query 'sort_by(imageDetails,&imagePushedAt)[-5:].[imageTags[0],imagePushedAt]' \
  --output table

# 4. Pull versión anterior
VERSION_ROLLBACK="2026.01.4-f214c2d"
podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:${VERSION_ROLLBACK}

# 5. Retag como latest
podman tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:${VERSION_ROLLBACK} \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest

# 6. Recrear contenedor
cd ~/probability/infra/compose-prod
podman stop backend_prod
podman rm backend_prod
podman-compose -f podman-compose.yaml up -d --no-deps back-central

# 7. Reiniciar nginx
podman restart nginx_prod

# 8. Verificar
echo "Esperando 20 segundos..."
sleep 20
podman logs --tail 50 backend_prod

# 9. Confirmar que está usando la versión correcta
podman inspect backend_prod --format '{{.Config.Image}}'

# 10. Verificar que el servicio responde
curl -I https://app.probabilityia.com.co/api/v1/health || echo "Health check endpoint"
```

## Contacto de emergencia

Si tienes problemas durante el rollback, contacta al equipo de DevOps o consulta los logs detallados:

```bash
# Logs completos
podman logs <service>_prod > /tmp/rollback-debug.log

# Estado del sistema
podman ps -a
podman images
```
