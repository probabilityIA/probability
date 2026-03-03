#!/bin/sh
# Startup script with backend connectivity verification
# Waits for testing backend before starting Next.js server

set -e

BACKEND_HOST="${TESTING_BACKEND_URL:-http://back-testing:9092}"
MAX_RETRIES=10
RETRY_INTERVAL=5

echo "Starting testing frontend with backend verification..."
echo "Backend URL: $BACKEND_HOST"
echo "Health check URL: $BACKEND_HOST/api/v1/health"

# Function to check backend connectivity
check_backend() {
    wget -q -O- -T 3 "$BACKEND_HOST/api/v1/health" >/dev/null 2>&1
    return $?
}

# Wait for backend with retries
echo "Checking backend connectivity..."
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if check_backend; then
        echo "Backend is reachable at $BACKEND_HOST"
        break
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Attempt $RETRY_COUNT/$MAX_RETRIES: Backend not available, retrying in ${RETRY_INTERVAL}s..."

    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        echo "WARN: Backend still not available after $MAX_RETRIES attempts, starting anyway..."
    fi

    sleep $RETRY_INTERVAL
done

echo "CENTRAL_API_URL=${CENTRAL_API_URL:-not set}"
echo "Starting Next.js server on port 3051..."
exec node server.js
