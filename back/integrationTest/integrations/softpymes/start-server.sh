#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "🚀 Starting SoftPymes Mock HTTP Server..."

# Puerto configurado
PORT=${SOFTPYMES_MOCK_PORT:-8082}

# Verificar si el puerto está en uso
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo "⚠️  Port $PORT is already in use"
    echo "🛑 Stopping existing process..."
    lsof -ti:$PORT | xargs kill -9 2>/dev/null || true
    sleep 1
fi

# Ejecutar servidor
echo "▶️  Starting server on port $PORT..."
go run server/main.go
