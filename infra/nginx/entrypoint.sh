#!/bin/sh
# Nginx entrypoint. Resuelve el color activo de blue-green y verifica que el
# backend responda antes de arrancar. Si no responde, sale con error y docker
# reinicia el contenedor.
#
# El frontend NO se espera a proposito: nginx tiene que poder arrancar aunque el
# frontend este abajo, porque el frontend le pega al backend a traves de nginx
# (puerto interno 8088). Esperarlo seria un ciclo.

set -e

UPSTREAMS_DIR=/etc/nginx/upstreams
ACTIVE_FILE="$UPSTREAMS_DIR/active.conf"
MAX_RETRIES=10
RETRY_INTERVAL=5

mkdir -p "$UPSTREAMS_DIR"

if [ ! -f "$ACTIVE_FILE" ]; then
    echo "⚠️  No hay $ACTIVE_FILE, escribiendo el color por defecto (blue)"
    cat > "$ACTIVE_FILE" <<EOF
upstream probability_backend {
    server back-central-blue:3050 max_fails=3 fail_timeout=30s;
}

upstream probability_frontend {
    server front-central-blue:3000 max_fails=3 fail_timeout=30s;
}
EOF
fi

BACKEND_TARGET=$(grep -o '[a-z-]*back-central-[a-z]*:[0-9]*' "$ACTIVE_FILE" | head -1)
BACKEND_TARGET=${BACKEND_TARGET:-back-central-blue:3050}
BACKEND_URL="http://$BACKEND_TARGET"

echo "🚀 Starting nginx"
echo "🎨 Backend activo: $BACKEND_TARGET"

RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if wget -q -O- -T 3 "$BACKEND_URL/health" >/dev/null 2>&1; then
        echo "✅ Backend alcanzable en $BACKEND_URL"
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "⚠️  Intento $RETRY_COUNT/$MAX_RETRIES: backend no disponible, reintentando en ${RETRY_INTERVAL}s..."

    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        echo "❌ PANIC: backend sigue sin responder tras $MAX_RETRIES intentos"
        exit 1
    fi

    sleep $RETRY_INTERVAL
done

echo "🎯 Arrancando nginx..."
envsubst '\$DOMAIN \$SSL_CERT_PATH \$SSL_KEY_PATH' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf
exec nginx -g 'daemon off;'
