#!/bin/sh
# Startup script with backend connectivity verification
# If backend is not available after retries, panic and let Podman restart

set -e

BACKEND_URL="${API_BASE_URL:-http://back-central:3050}"
MAX_RETRIES=10
RETRY_INTERVAL=5

echo "🚀 Starting frontend with backend verification..."
echo "📡 Backend URL: $BACKEND_URL"

# Function to check backend connectivity
check_backend() {
    wget -q -O- -T 3 "$BACKEND_URL/health" >/dev/null 2>&1
    return $?
}

# Wait for backend with retries
echo "🔍 Checking backend connectivity..."
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if check_backend; then
        echo "✅ Backend is reachable at $BACKEND_URL"
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
