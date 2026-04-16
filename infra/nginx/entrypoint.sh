#!/bin/sh
# Nginx entrypoint with upstream connectivity verification
# If upstreams are not available after retries, panic and let Podman restart

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
    wget -q -O- -T 3 "$URL/health" >/dev/null 2>&1 || wget -q -O- -T 3 "$URL" >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ $NAME is reachable at $URL"
        return 0
    fi
    return 1
}

# Wait for backend
echo "🔍 Checking backend connectivity..."
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if check_upstream "$BACKEND_URL" "Backend"; then
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "⚠️  Attempt $RETRY_COUNT/$MAX_RETRIES: Backend not available, retrying in ${RETRY_INTERVAL}s..."

    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        echo "❌ PANIC: Backend still not available after $MAX_RETRIES attempts"
        echo "💥 Exiting with error code to trigger container restart..."
        exit 1
    fi

    sleep $RETRY_INTERVAL
done

# Wait for frontend
echo "🔍 Checking frontend connectivity..."
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if check_upstream "$FRONTEND_URL" "Frontend"; then
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "⚠️  Attempt $RETRY_COUNT/$MAX_RETRIES: Frontend not available, retrying in ${RETRY_INTERVAL}s..."

    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        echo "❌ PANIC: Frontend still not available after $MAX_RETRIES attempts"
        echo "💥 Exiting with error code to trigger container restart..."
        exit 1
    fi

    sleep $RETRY_INTERVAL
done

echo "🎯 All upstreams are healthy, starting nginx..."
envsubst '\$DOMAIN \$SSL_CERT_PATH \$SSL_KEY_PATH' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf
exec nginx -g 'daemon off;'
# Deploy trigger 2026-01-29_17:18:31
