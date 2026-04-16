#!/bin/sh
# Startup script with backend connectivity verification
# If backend is not available after retries, panic and let Podman restart

set -e

# API_BASE_URL puede tener path (/api/v1), extraer solo el host para health check
BACKEND_URL="${API_BASE_URL:-http://back-central:3050}"
# Extraer solo el host (sin path) para el health check
BACKEND_HOST=$(echo "$BACKEND_URL" | sed 's|^\(https\?://[^/]*\).*|\1|')
MAX_RETRIES=10
RETRY_INTERVAL=5

echo "🚀 Starting frontend with backend verification..."
echo "📡 Backend URL: $BACKEND_URL"
echo "🏥 Health check URL: $BACKEND_HOST/health"

# Function to check backend connectivity
check_backend() {
    wget -q -O- -T 3 "$BACKEND_HOST/health" >/dev/null 2>&1
    return $?
}

# Wait for backend with retries
echo "🔍 Checking backend connectivity..."
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if check_backend; then
        echo "✅ Backend is reachable at $BACKEND_HOST"
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

echo "🎯 Backend is healthy, starting Next.js server..."
exec node server.js
# Deploy trigger 2026-01-29_17:18:31
