# Sistema de Deployment Automático - Probability

## Arquitectura de Deployment

```
┌─────────────────────────────────────────┐
│  GitHub Push to main                    │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  GitHub Actions (concurrency-controlled)│
│  - Build ARM64 image (Podman)           │
│  - Push to ECR (version-tagged)         │
│  - SSH to EC2 server                    │
│  - Deploy con verificación              │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  Server Deployment Steps                │
│  1. Pull new image (version-tagged)     │
│  2. Stop old container (-t 10)          │
│  3. Remove with --depend flag           │
│  4. Kill port processes (fuser -k)      │
│  5. Pull specific version again         │
│  6. Tag as latest locally               │
│  7. Start new container                 │
│  8. Verify container status (3 retries) │
│  9. Verify image ID match               │
│  10. Clean old images                   │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  Panic/Restart Mechanism                │
│  - Frontend: Verifica backend health    │
│    Si falla → exit 1 → Podman restart   │
│  - Nginx: Verifica backend + frontend   │
│    Si falla → exit 1 → Podman restart   │
│  - restart: always → auto-recovery      │
└─────────────────────────────────────────┘
```

## Workflows Disponibles

### 1. Backend CI/CD
- **Trigger:** Push a `back/central/**` o `back/migration/**`
- **Concurrency group:** `backend-deploy`
- **Puerto:** 3050 (interno), expuesto solo a localhost
- **Tiempo estimado:** ~3-4 minutos

### 2. Frontend CI/CD
- **Trigger:** Push a `front/central/**`
- **Concurrency group:** `frontend-deploy`
- **Puerto:** 8080 (expuesto externamente)
- **Tiempo estimado:** ~3-4 minutos

### 3. Website CI/CD
- **Trigger:** Push a `front/website/**`
- **Concurrency group:** `website-deploy`
- **Puerto:** 8081 (expuesto externamente)
- **Tiempo estimado:** ~2-3 minutos

### 4. Nginx CI/CD
- **Trigger:** Push a `infra/nginx/**`
- **Concurrency group:** `nginx-deploy`
- **Puertos:** 80, 443 (expuestos externamente)
- **Tiempo estimado:** ~2 minutos

## Características Clave

### Concurrency Control
Cada workflow tiene su propio grupo de concurrencia con `cancel-in-progress: true`:
- Previene deployments simultáneos del mismo servicio
- Cancela workflows en progreso cuando llega un nuevo push
- Elimina race conditions y conflictos de puertos

### Version Tagging
Cada imagen se tagea con 4 tags:
- `latest` - Última versión deployada
- `VERSION_FULL` - Ejemplo: `2026.01.37.a3d554c`
- `VERSION_SHORT` - Ejemplo: `2026.01.37`
- `SHORT_SHA` - Ejemplo: `a3d554c`

### Verificación de Deployment
Después de deployar, el workflow verifica:
1. Container status = "running"
2. Image ID match (nueva imagen vs contenedor)
3. Container logs (últimas 30 líneas)
4. Reintentos automáticos si falla (máx 3)

### Auto-Recovery con Panic/Restart
Ver: `panic-restart-mechanism.md`

## Estado de Contenedores en Producción

```bash
# Ver estado
ssh ubuntu@app.probabilityia.com.co 'podman ps'

# Logs de un contenedor
ssh ubuntu@app.probabilityia.com.co 'podman logs --tail 50 backend_prod'

# Reiniciar un servicio manualmente
ssh ubuntu@app.probabilityia.com.co 'cd ~/probability/infra/compose-prod && podman-compose restart backend'
```

## Problemas Comunes y Soluciones

### Problema: Nginx da 502 después de deploy de backend/frontend
**Causa:** Nginx cachea las IPs de los upstreams

**Solución:**
```bash
ssh ubuntu@app.probabilityia.com.co 'podman restart nginx_prod'
```

**Mejora futura:** Automatizar el restart de nginx después de deployments de upstreams.

### Problema: Puerto ocupado durante deployment
**Causa:** Proceso zombie de Podman (rootlessport)

**Solución:** El workflow ya incluye:
```bash
sudo fuser -k <PORT>/tcp
pkill -9 -f "rootlessport.*<PORT>"
```

### Problema: Container en estado "Stopping" permanente
**Causa:** Contenedor tiene dependencias (depends_on fue removido, pero puede haber referencias)

**Solución:**
```bash
podman rm -f --depend <container_name>
```

### Problema: Frontend/Nginx en loop de restart
**Causa:** Backend no disponible, panic/restart mechanism activo

**Verificar:**
```bash
# Ver logs del backend
podman logs backend_prod

# Ver health del backend
curl http://localhost:3050/health
```

## Monitoreo de Deployments

### Via GitHub API
```bash
gh run list --limit 5
gh run watch <run-id>
```

### Via Logs del Servidor
```bash
ssh ubuntu@app.probabilityia.com.co 'journalctl -u podman -f'
```

## Rollback Manual

Si un deployment falla y necesitas hacer rollback:

```bash
# 1. Conectarse al servidor
ssh ubuntu@app.probabilityia.com.co

# 2. Ir al directorio de compose
cd ~/probability/infra/compose-prod

# 3. Listar imágenes disponibles
podman images | grep probability-backend

# 4. Tagear versión anterior como latest
podman tag 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:2026.01.36 \
           476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest

# 5. Reiniciar el servicio
podman-compose restart backend
```

## Variables de Entorno

Las variables de entorno están en:
- **Producción:** `~/probability/infra/compose-prod/.env` (en el servidor)
- **Workflows:** GitHub Secrets

### Secrets Necesarios en GitHub
- `AWS_ACCESS_KEY_ID` - Para ECR
- `AWS_SECRET_ACCESS_KEY` - Para ECR
- `EC2_SSH_KEY` - Private key para SSH al servidor
- `EC2_USER` - Usuario SSH (ubuntu)
- `EC2_HOST` - IP o dominio del servidor
- `NEXT_PUBLIC_API_BASE_URL` - URL pública del API (frontend)
- `API_BASE_URL` - URL interna del API (frontend SSR)

## Mejoras Futuras

1. **Auto-restart de nginx:** Agregar step en workflows de backend/frontend para reiniciar nginx
2. **Health checks más robustos:** Verificar endpoints específicos además de `/health`
3. **Notificaciones:** Slack/Discord notifications para deployments exitosos/fallidos
4. **Métricas:** Prometheus/Grafana para monitorear uptime durante deployments
5. **Staging environment:** Ambiente de staging para probar deployments antes de producción
6. **Blue-Green deployment:** Para zero-downtime garantizado
7. **Automated rollback:** Rollback automático si health checks fallan después del deploy

## Contacto

Para problemas con el sistema de deployment, contactar al equipo de DevOps o revisar:
- GitHub Actions logs
- Logs del servidor: `/var/log/` y `podman logs`
- CLAUDE.md en la raíz del proyecto
