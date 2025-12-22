#!/bin/bash

# Script para purgar la cola de RabbitMQ
# Uso: ./scripts/purge_rabbitmq_queue.sh [queue_name]

set -e

# Configuración por defecto
RABBITMQ_HOST="${RABBITMQ_HOST:-localhost}"
RABBITMQ_PORT="${RABBITMQ_PORT:-15672}"
RABBITMQ_USER="${RABBITMQ_USER:-admin}"
RABBITMQ_PASS="${RABBITMQ_PASS:-admin}"
RABBITMQ_VHOST="${RABBITMQ_VHOST:-%2F}"  # %2F es la codificación URL de "/"
QUEUE_NAME="${1:-probability.orders.canonical}"

echo "🐰 Purgando cola de RabbitMQ..."
echo "   Host: $RABBITMQ_HOST:$RABBITMQ_PORT"
echo "   Usuario: $RABBITMQ_USER"
echo "   Cola: $QUEUE_NAME"
echo "   VHost: $RABBITMQ_VHOST"
echo ""

# URL de la API de RabbitMQ para purgar la cola
PURGE_URL="http://${RABBITMQ_HOST}:${RABBITMQ_PORT}/api/queues/${RABBITMQ_VHOST}/${QUEUE_NAME}/contents"

# Primero, verificar el estado de la cola
echo "📊 Verificando estado de la cola..."
STATUS_URL="http://${RABBITMQ_HOST}:${RABBITMQ_PORT}/api/queues/${RABBITMQ_VHOST}/${QUEUE_NAME}"

STATUS_RESPONSE=$(curl -s -u "${RABBITMQ_USER}:${RABBITMQ_PASS}" "$STATUS_URL" || echo "{}")

if echo "$STATUS_RESPONSE" | grep -q "error"; then
    echo "❌ Error: La cola '$QUEUE_NAME' no existe o no se puede acceder"
    echo "   Respuesta: $STATUS_RESPONSE"
    exit 1
fi

# Extraer información de la cola
MESSAGES=$(echo "$STATUS_RESPONSE" | grep -o '"messages":[0-9]*' | grep -o '[0-9]*' || echo "0")
MESSAGES_READY=$(echo "$STATUS_RESPONSE" | grep -o '"messages_ready":[0-9]*' | grep -o '[0-9]*' || echo "0")
MESSAGES_UNACKNOWLEDGED=$(echo "$STATUS_RESPONSE" | grep -o '"messages_unacknowledged":[0-9]*' | grep -o '[0-9]*' || echo "0")

echo "   Mensajes totales: $MESSAGES"
echo "   Mensajes listos: $MESSAGES_READY"
echo "   Mensajes sin confirmar: $MESSAGES_UNACKNOWLEDGED"
echo ""

if [ "$MESSAGES" -eq 0 ] && [ "$MESSAGES_READY" -eq 0 ]; then
    echo "✅ La cola ya está vacía"
    exit 0
fi

# Confirmar antes de purgar
read -p "⚠️  ¿Estás seguro de que quieres purgar la cola? (s/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[SsYy]$ ]]; then
    echo "❌ Operación cancelada"
    exit 0
fi

# Purgar la cola usando DELETE (solo elimina mensajes ready)
echo "🗑️  Purgando mensajes listos (ready)..."
PURGE_RESPONSE=$(curl -s -X DELETE -u "${RABBITMQ_USER}:${RABBITMQ_PASS}" "$PURGE_URL" || echo "{}")

if echo "$PURGE_RESPONSE" | grep -q "error"; then
    echo "⚠️  No se pudieron purgar los mensajes listos: $PURGE_RESPONSE"
else
    echo "✅ Mensajes listos purgados"
fi

# Para purgar mensajes sin confirmar, necesitamos usar rabbitmqctl o reiniciar el consumidor
# Intentar usar la API de purga completa (requiere permisos de administrador)
PURGE_COMPLETE_URL="http://${RABBITMQ_HOST}:${RABBITMQ_PORT}/api/queues/${RABBITMQ_VHOST}/${QUEUE_NAME}/contents"
echo ""
echo "🗑️  Intentando purgar todos los mensajes (incluyendo sin confirmar)..."

# Método alternativo: usar rabbitmqctl si está disponible
if command -v rabbitmqctl &> /dev/null; then
    echo "   Usando rabbitmqctl..."
    docker exec rabbitmq_local rabbitmqctl purge_queue "$QUEUE_NAME" 2>/dev/null || \
    rabbitmqctl purge_queue "$QUEUE_NAME" 2>/dev/null || \
    echo "   rabbitmqctl no disponible, intentando con API..."
fi

# Si hay mensajes sin confirmar, sugerir reiniciar el consumidor
if [ "$MESSAGES_UNACKNOWLEDGED" -gt 0 ]; then
    echo ""
    echo "⚠️  IMPORTANTE: Hay $MESSAGES_UNACKNOWLEDGED mensajes sin confirmar."
    echo "   Estos mensajes no se pueden purgar mientras el consumidor esté activo."
    echo "   Opciones:"
    echo "   1. Detener el servicio que consume la cola temporalmente"
    echo "   2. Reiniciar el servicio para que los mensajes vuelvan a la cola"
    echo "   3. Esperar a que se procesen (puede tomar mucho tiempo)"
fi

echo "✅ Proceso de purga completado"

# Verificar el estado después de purgar
echo ""
echo "📊 Verificando estado después de la purga..."
sleep 1
STATUS_RESPONSE_AFTER=$(curl -s -u "${RABBITMQ_USER}:${RABBITMQ_PASS}" "$STATUS_URL" || echo "{}")
MESSAGES_AFTER=$(echo "$STATUS_RESPONSE_AFTER" | grep -o '"messages":[0-9]*' | grep -o '[0-9]*' || echo "0")
MESSAGES_READY_AFTER=$(echo "$STATUS_RESPONSE_AFTER" | grep -o '"messages_ready":[0-9]*' | grep -o '[0-9]*' || echo "0")

echo "   Mensajes totales: $MESSAGES_AFTER"
echo "   Mensajes listos: $MESSAGES_READY_AFTER"
echo ""

if [ "$MESSAGES_AFTER" -eq 0 ] && [ "$MESSAGES_READY_AFTER" -eq 0 ]; then
    echo "✅ La cola está completamente vacía"
else
    echo "⚠️  Aún hay mensajes en la cola (pueden estar siendo procesados)"
fi
