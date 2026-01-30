# Troubleshooting - Sistema de Deployment

## Índice de Problemas Comunes

1. [Workflow falla en Build](#workflow-falla-en-build)
2. [Workflow falla en Deploy](#workflow-falla-en-deploy)
3. [Container en loop de restart](#container-en-loop-de-restart)
4. [Nginx da 502 Bad Gateway](#nginx-da-502-bad-gateway)
5. [Puerto ocupado durante deploy](#puerto-ocupado-durante-deploy)
6. [Imagen no se actualiza](#imagen-no-se-actualiza)
7. [Container queda en "Stopping"](#container-queda-en-stopping)
8. [Site down después de deploy](#site-down-después-de-deploy)

---

## Workflow falla en Build

### Error: "podman login failed"

**Síntomas:**
```
Error: error logging in to "476702565908.dkr.ecr.us-east-1.amazonaws.com"
```

**Causa:** Credenciales AWS inválidas o expiradas

**Solución:**
1. Verificar secretos en GitHub:
   ```bash
   # En GitHub repo → Settings → Secrets → Actions
   # Verificar: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY
   ```

2. Regenerar credenciales en AWS IAM si es necesario

3. Re-ejecutar el workflow

### Error: "podman build failed"

**Síntomas:**
```
Error: error building at STEP "RUN npm ci": exit status 1
```

**Causa:** Fallo en instalación de dependencias o build

**Solución:**
1. Revisar logs del workflow para el error específico

2. Probar build localmente:
   ```bash
   cd front/central  # O el directorio correspondiente
   podman build -f docker/Dockerfile .
   ```

3. Si el error es de dependencias:
   ```bash
   # Actualizar lockfiles
   npm install
   git add package-lock.json
   git commit -m "fix: update package-lock.json"
   git push
   ```

### Error: "ECR repository not found"

**Síntomas:**
```
Error: repository does not exist: probability-backend
```

**Causa:** Repositorio ECR no existe en AWS

**Solución:**
```bash
# Crear repositorio en ECR
aws ecr create-repository \
  --repository-name probability-backend \
  --region us-east-1
```

---

## Workflow falla en Deploy

### Error: "SSH connection failed"

**Síntomas:**
```
Permission denied (publickey)
```

**Causa:** SSH key inválida o servidor no alcanzable

**Solución:**
1. Verificar secret `EC2_SSH_KEY`:
   ```bash
   # La key debe estar en formato PEM, completa:
   -----BEGIN RSA PRIVATE KEY-----
   ...
   -----END RSA PRIVATE KEY-----
   ```

2. Verificar que el servidor esté corriendo:
   ```bash
   ping app.probabilityia.com.co
   ```

3. Verificar security groups en AWS (puerto 22 abierto)

### Error: "podman-compose command not found"

**Síntomas:**
```
bash: podman-compose: command not found
```

**Causa:** podman-compose no está instalado en el servidor

**Solución:**
```bash
# Conectarse al servidor
ssh ubuntu@server

# Instalar podman-compose
pip3 install podman-compose

# Verificar instalación
podman-compose --version
```

### Error: "Container failed to start"

**Síntomas:**
```
❌ El contenedor no está corriendo después de 3 intentos
Estado final: exited
```

**Causa:** El contenedor inicia pero sale inmediatamente

**Solución:**
1. Ver logs del contenedor:
   ```bash
   ssh ubuntu@server 'podman logs backend_prod'
   ```

2. Verificar variables de entorno en `.env`

3. Verificar que las dependencias (DB, Redis) estén disponibles

---

## Container en Loop de Restart

### Frontend restarting continuamente

**Síntomas:**
```bash
podman ps
# frontend_prod  Up 3 seconds (starting)  # Vuelve a reiniciar
```

**Causa:** Backend no disponible, panic/restart mechanism activo

**Diagnóstico:**
```bash
# 1. Ver logs del frontend
podman logs --tail 50 frontend_prod

# Buscar:
# ❌ PANIC: Backend not available after 10 attempts
```

**Solución:**
```bash
# 2. Verificar backend
podman ps | grep backend
curl http://localhost:3050/health

# 3. Si backend está caído, levantarlo
cd ~/probability/infra/compose-prod
podman-compose up -d back-central

# 4. El frontend se auto-recuperará en 1-2 minutos
```

### Nginx restarting continuamente

**Síntomas:**
Nginx reinicia constantemente

**Causa:** Backend o frontend no disponibles

**Diagnóstico:**
```bash
# Ver logs de nginx
podman logs --tail 50 nginx_prod

# Buscar:
# ⚠️  Attempt 10/10: Frontend not available
# ❌ PANIC: Backend not available
```

**Solución:**
```bash
# Verificar que backend y frontend estén corriendo
podman ps | grep -E "backend|frontend"

# Levantar los que falten
podman-compose up -d back-central front-central

# Nginx se auto-recuperará
```

---

## Nginx da 502 Bad Gateway

### Después de deploy de backend/frontend

**Síntomas:**
```bash
curl https://app.probabilityia.com.co
# 502 Bad Gateway
```

**Causa:** Nginx cachea las IPs de los upstreams y no se actualiza automáticamente

**Solución:**
```bash
# Opción 1: Reiniciar nginx (rápido)
ssh ubuntu@server 'podman restart nginx_prod'

# Opción 2: Verificar y reconectar manualmente
ssh ubuntu@server << 'EOF'
  # Verificar DNS
  podman exec nginx_prod nslookup back-central
  podman exec nginx_prod nslookup front-central

  # Reiniciar nginx
  podman restart nginx_prod
EOF
```

### Nginx no puede conectarse a upstream

**Síntomas:**
```
connect() failed (113: Host is unreachable) while connecting to upstream
```

**Causa:** Firewall/network rules bloqueando conexiones

**Solución:**
```bash
# 1. Verificar reglas de iptables
sudo iptables -L FORWARD -n | head -5

# Debe ver:
# Chain FORWARD (policy ACCEPT)
# ACCEPT  all -- 10.89.0.0/24  anywhere

# 2. Si la policy es DROP:
sudo iptables -P FORWARD ACCEPT

# 3. Agregar reglas si faltan:
sudo iptables -I FORWARD 1 -s 10.89.0.0/24 -j ACCEPT
sudo iptables -I FORWARD 2 -d 10.89.0.0/24 -j ACCEPT

# 4. Reiniciar nginx
podman restart nginx_prod
```

---

## Puerto ocupado durante deploy

### Error: "address already in use"

**Síntomas:**
```
Error: rootlessport listen tcp 0.0.0.0:8080: bind: address already in use
```

**Causa:** Proceso zombie de Podman (rootlessport) o contenedor viejo

**Solución:**
```bash
# 1. Identificar proceso ocupando el puerto
sudo lsof -i :8080

# 2. Matar proceso (el workflow ya debería hacer esto)
sudo fuser -k 8080/tcp
pkill -9 -f "rootlessport.*8080"

# 3. Verificar que no hay contenedores del mismo nombre
podman ps -a | grep frontend_prod

# 4. Remover contenedor con dependencias
podman rm -f --depend frontend_prod

# 5. Intentar levantar nuevamente
podman-compose up -d front-central
```

---

## Imagen no se actualiza

### Container usando imagen vieja después de deploy

**Síntomas:**
Workflow completa exitosamente pero el contenedor usa la imagen anterior

**Diagnóstico:**
```bash
# 1. Ver workflow logs (verificar VERSION_FULL)
gh run view --log

# 2. Verificar imagen en el servidor
ssh ubuntu@server << 'EOF'
  # Ver imágenes disponibles
  podman images | grep probability-backend

  # Ver qué imagen usa el contenedor
  podman inspect backend_prod --format '{{.Image}}'
EOF
```

**Solución:**
```bash
ssh ubuntu@server << 'EOF'
  cd ~/probability/infra/compose-prod

  # 1. Ver última versión en ECR
  aws ecr describe-images \
    --repository-name probability-backend \
    --region us-east-1 \
    --query 'sort_by(imageDetails,& imagePushedAt)[-1].imageTags[0]'

  # 2. Pull manual de la versión específica
  VERSION_FULL="2026.01.37.a3d554c"  # Usar la del workflow

  podman pull 476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:$VERSION_FULL

  # 3. Tag como latest
  podman tag \
    476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:$VERSION_FULL \
    476702565908.dkr.ecr.us-east-1.amazonaws.com/probability-backend:latest

  # 4. Forzar recreación del contenedor
  podman stop backend_prod
  podman rm -f backend_prod
  podman-compose up -d back-central
EOF
```

---

## Container queda en "Stopping"

### Container stuck en estado "Stopping"

**Síntomas:**
```bash
podman ps -a
# backend_prod  Stopping  # No termina de detenerse
```

**Causa:** Contenedor no responde a SIGTERM, esperando timeout

**Solución:**
```bash
# 1. Forzar kill
podman kill backend_prod

# 2. Si no funciona, eliminar por ID con dependencias
CONTAINER_ID=$(podman ps -a --filter name=backend_prod --format '{{.ID}}')
podman rm -f --depend $CONTAINER_ID

# 3. Limpiar procesos zombie
sudo fuser -k 3050/tcp
pkill -9 -f "rootlessport.*3050"

# 4. Levantar nuevamente
podman-compose up -d back-central
```

---

## Site down después de deploy

### Todo parece funcionar pero el sitio no carga

**Diagnóstico completo:**

```bash
# 1. Verificar estado de contenedores
ssh ubuntu@server 'podman ps'

# Todos deben estar "Up" y "healthy"

# 2. Verificar acceso local
ssh ubuntu@server << 'EOF'
  curl -I http://localhost:8080  # Frontend
  curl -I http://localhost:3050/health  # Backend
  curl -I http://localhost:80  # Nginx HTTP
  curl -Ik https://localhost:443  # Nginx HTTPS
EOF

# 3. Verificar desde Internet
curl -I https://app.probabilityia.com.co

# 4. Verificar Security Groups AWS
# - Puerto 80: 0.0.0.0/0
# - Puerto 443: 0.0.0.0/0

# 5. Verificar iptables (ver problema "Nginx da 502")
```

### Checklist completo:

- [ ] Contenedores corriendo: `podman ps`
- [ ] Health checks pasando: `podman ps` (columna STATUS)
- [ ] Logs sin errores: `podman logs <service>_prod`
- [ ] Puertos listening: `sudo ss -tlnp | grep -E "80|443|8080|3050"`
- [ ] iptables FORWARD=ACCEPT: `sudo iptables -L FORWARD -n | head -1`
- [ ] Security Groups abiertos (AWS Console)
- [ ] DNS resuelve: `nslookup app.probabilityia.com.co`

---

## Comandos Útiles de Diagnóstico

### Ver todo el estado del sistema

```bash
#!/bin/bash
# diagnostico-completo.sh

echo "=== CONTENEDORES ==="
podman ps -a --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'

echo -e "\n=== IMÁGENES ==="
podman images | grep probability

echo -e "\n=== HEALTH CHECKS ==="
curl -s -o /dev/null -w "Frontend: %{http_code}\n" http://localhost:8080
curl -s -o /dev/null -w "Backend: %{http_code}\n" http://localhost:3050/health
curl -s -o /dev/null -w "Nginx HTTP: %{http_code}\n" http://localhost:80
curl -sk -o /dev/null -w "Nginx HTTPS: %{http_code}\n" https://localhost:443

echo -e "\n=== PUERTOS ==="
sudo ss -tlnp | grep -E "80|443|8080|3050|8081"

echo -e "\n=== IPTABLES ==="
sudo iptables -L FORWARD -n | head -3

echo -e "\n=== DISK SPACE ==="
df -h / | tail -1

echo -e "\n=== MEMORY ==="
free -h

echo -e "\n=== ÚLTIMO DEPLOY ==="
podman images --format 'table {{.Repository}}\t{{.Tag}}\t{{.CreatedAt}}' | \
  grep probability | head -5
```

### Ver logs en tiempo real

```bash
# Todos los contenedores
podman ps --format '{{.Names}}' | \
  xargs -I {} sh -c 'echo "=== {} ===" && podman logs --tail 5 {}'

# Un contenedor específico con follow
podman logs -f backend_prod

# Filtrar solo errores
podman logs backend_prod 2>&1 | grep -i error
```

### Restart completo del sistema

```bash
#!/bin/bash
# restart-completo.sh
# ⚠️  USAR SOLO EN EMERGENCIA

cd ~/probability/infra/compose-prod

echo "Deteniendo todos los servicios..."
podman-compose down

echo "Limpiando procesos zombie..."
sudo fuser -k 80/tcp 443/tcp 8080/tcp 8081/tcp 3050/tcp
pkill -9 -f rootlessport

echo "Esperando 5 segundos..."
sleep 5

echo "Levantando servicios en orden..."
podman-compose up -d redis rabbitmq
sleep 10

podman-compose up -d back-central
sleep 15

podman-compose up -d front-central front-website
sleep 15

podman-compose up -d nginx

echo "Estado final:"
podman ps
```

---

## Contacto y Escalación

Si ninguna de estas soluciones funciona:

1. **Recopilar información:**
   ```bash
   # Ejecutar script de diagnóstico
   bash diagnostico-completo.sh > diagnostico.txt

   # Obtener logs de workflows
   gh run view --log > workflow.log
   ```

2. **Revisar documentación adicional:**
   - README.md
   - panic-restart-mechanism.md
   - workflow-structure.md

3. **Buscar en logs del sistema:**
   ```bash
   journalctl -u podman -n 100
   ```

4. **GitHub Issues:**
   Crear issue en el repositorio con:
   - Síntomas observados
   - Logs relevantes
   - Output de diagnóstico-completo.sh
   - Workflow run ID
