# Mecanismo de Panic/Restart

## Concepto

El mecanismo de panic/restart permite que los contenedores se auto-recuperen cuando sus dependencias no están disponibles. En lugar de iniciar en un estado roto, el contenedor verifica sus dependencias y **sale con error (panic)** si no están disponibles, permitiendo que Podman lo reinicie automáticamente gracias a `restart: always`.

## Implementación

### Frontend (front/central/docker/startup.sh)

El frontend verifica que el backend esté disponible antes de iniciar Next.js:

```bash
#!/bin/sh
set -e

# Extraer solo el host para health checks
BACKEND_HOST=$(echo "$BACKEND_URL" | sed 's|^\(https\?://[^/]*\).*|\1|')
MAX_RETRIES=10
RETRY_INTERVAL=5

echo "🚀 Starting frontend with backend verification..."
echo "🏥 Health check URL: $BACKEND_HOST/health"

# Function to check backend connectivity
check_backend() {
    wget -q -O- -T 3 "$BACKEND_HOST/health" >/dev/null 2>&1
    return $?
}

# Wait for backend with retries
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if check_backend; then
        echo "✅ Backend is reachable"
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))

    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        echo "❌ PANIC: Backend not available after $MAX_RETRIES attempts"
        echo "💥 Exiting with error code to trigger container restart..."
        exit 1  # ← Panic! Podman reiniciará el contenedor
    fi

    sleep $RETRY_INTERVAL
done

echo "🎯 Backend is healthy, starting Next.js server..."
exec node server.js
```

**Comportamiento:**
1. Intenta conectar al backend 10 veces (50 segundos total)
2. Si el backend no responde → `exit 1`
3. Podman detecta el exit code 1 → reinicia el contenedor
4. El contenedor vuelve a intentar conectar al backend
5. Cuando el backend esté disponible → inicia Next.js

### Nginx (infra/nginx/entrypoint.sh)

Nginx verifica que backend Y frontend estén disponibles:

```bash
#!/bin/sh
set -e

BACKEND_URL="http://back-central:3050"
FRONTEND_URL="http://front-central:3000"
MAX_RETRIES=10
RETRY_INTERVAL=5

echo "🚀 Starting nginx with upstream verification..."

# Function to check upstream connectivity
check_upstream() {
    URL=$1
    NAME=$2
    wget -q -O- -T 3 "$URL/health" >/dev/null 2>&1 || \
    wget -q -O- -T 3 "$URL" >/dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "✅ $NAME is reachable at $URL"
        return 0
    fi
    return 1
}

# Wait for backend
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if check_upstream "$BACKEND_URL" "Backend"; then
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))

    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        echo "❌ PANIC: Backend not available"
        exit 1  # ← Panic!
    fi

    sleep $RETRY_INTERVAL
done

# Wait for frontend (mismo patrón)
# ...

echo "🎯 All upstreams are healthy, starting nginx..."
envsubst '\\$DOMAIN \\$SSL_CERT_PATH \\$SSL_KEY_PATH' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf
exec nginx -g 'daemon off;'
```

**Comportamiento:**
1. Verifica backend (10 intentos)
2. Verifica frontend (10 intentos)
3. Si alguno falla → `exit 1` → Podman reinicia nginx
4. Cuando ambos estén disponibles → inicia nginx

## Configuración en Docker/Podman Compose

```yaml
services:
  front-central:
    image: probability-frontend:latest
    restart: always  # ← Crítico para panic/restart
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3000/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 15s

  nginx:
    image: probability-nginx:latest
    restart: always  # ← Crítico para panic/restart
```

**IMPORTANTE:** NO usar `depends_on`. Las dependencias se manejan en código (entrypoint/startup scripts), no en Docker Compose.

## Ventajas

### 1. Auto-Recovery
Si el backend se cae y luego vuelve, el frontend automáticamente se recupera sin intervención manual.

### 2. Orden de Inicio No Importa
No importa si nginx inicia antes que el frontend. Nginx esperará hasta que el frontend esté disponible.

### 3. Zero-Config Deployments
Durante un deployment:
1. Backend nuevo se levanta
2. Frontend detecta que el backend viejo se cayó → panic
3. Frontend reinicia → conecta al backend nuevo
4. Todo funciona automáticamente

### 4. Fail Fast
Si hay un problema real (ej: backend no puede conectarse a la BD), el contenedor no inicia en un estado zombie. Sale con error inmediatamente y los logs muestran el problema claramente.

## Desventajas y Limitaciones

### 1. Tiempo de Inicio Más Largo
Si las dependencias no están listas, el contenedor puede tardar más en iniciar (hasta 50 segundos en el peor caso).

### 2. Restart Loops Visibles
Durante problemas, verás múltiples restarts en `podman ps`:
```
frontend_prod  Up 2 seconds (starting)  # Restart 1
frontend_prod  Up 5 seconds (starting)  # Restart 2
frontend_prod  Up 8 seconds (healthy)   # Finalmente inició
```

### 3. Logs Repetitivos
Los logs mostrarán múltiples intentos:
```
⚠️  Attempt 1/10: Backend not available, retrying in 5s...
⚠️  Attempt 2/10: Backend not available, retrying in 5s...
...
❌ PANIC: Backend not available after 10 attempts
💥 Exiting with error code to trigger container restart...
🚀 Starting frontend with backend verification...  # Nuevo intento
```

## Alternativas Consideradas

### 1. `depends_on` con `wait-for-it`
**Problema:** No maneja reconexiones después del inicio. Si el backend se cae después de iniciar, el frontend queda en estado roto.

### 2. Retry en la aplicación
**Problema:** La aplicación ya inició (ocupando puerto), pero no puede servir requests. Genera 502 errors.

### 3. Health checks de Docker/Podman
**Problema:** Solo marcan el contenedor como "unhealthy", pero no reinician automáticamente.

## Configuración de Timeouts

### Valores Actuales
- **MAX_RETRIES:** 10 intentos
- **RETRY_INTERVAL:** 5 segundos
- **Timeout total:** 50 segundos

### Ajustar según necesidad

Para servicios críticos que deben iniciar rápido:
```bash
MAX_RETRIES=5
RETRY_INTERVAL=3
# Timeout total: 15 segundos
```

Para servicios que pueden esperar más:
```bash
MAX_RETRIES=20
RETRY_INTERVAL=10
# Timeout total: 200 segundos (3.3 minutos)
```

## Monitoreo

### Ver logs de panic/restart
```bash
# Frontend
podman logs frontend_prod | grep -E "PANIC|Starting|Ready"

# Nginx
podman logs nginx_prod | grep -E "PANIC|Starting|upstreams"
```

### Contar restarts
```bash
# Ver cuántas veces se reinició un contenedor
podman inspect frontend_prod | jq '.RestartCount'
```

### Ver estado de health checks
```bash
podman ps --format 'table {{.Names}}\t{{.Status}}'
```

## Troubleshooting

### Problema: Frontend en loop infinito de restarts

**Verificar:**
1. ¿El backend está corriendo?
   ```bash
   podman ps | grep backend
   ```

2. ¿El backend responde en /health?
   ```bash
   curl http://localhost:3050/health
   ```

3. ¿El frontend puede resolver el DNS "back-central"?
   ```bash
   podman exec frontend_prod ping -c 2 back-central
   ```

### Problema: Nginx en loop de restarts

**Verificar:**
1. ¿Backend y frontend están corriendo?
   ```bash
   podman ps | grep -E "backend|frontend"
   ```

2. ¿Nginx puede resolver los DNS?
   ```bash
   podman exec nginx_prod nslookup back-central
   podman exec nginx_prod nslookup front-central
   ```

## Mejoras Futuras

1. **Metrics/Prometheus:** Exponer métricas de restart count
2. **Alerting:** Alertar si un contenedor reinicia más de N veces en X minutos
3. **Circuit breaker:** Dejar de intentar después de muchos fallos y requerir intervención manual
4. **Exponential backoff:** Aumentar el intervalo entre reintentos progresivamente
