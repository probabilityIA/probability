#!/bin/bash

# Script para verificar el estado de la cola de RabbitMQ
# Uso: ./scripts/check_rabbitmq_queue.sh [queue_name]

set -e

# Configuración por defecto
RABBITMQ_HOST="${RABBITMQ_HOST:-localhost}"
RABBITMQ_PORT="${RABBITMQ_PORT:-15672}"
RABBITMQ_USER="${RABBITMQ_USER:-admin}"
RABBITMQ_PASS="${RABBITMQ_PASS:-admin}"
RABBITMQ_VHOST="${RABBITMQ_VHOST:-%2F}"  # %2F es la codificación URL de "/"
QUEUE_NAME="${1:-probability.orders.canonical}"

echo "🐰 Estado de la cola de RabbitMQ"
echo "   Host: $RABBITMQ_HOST:$RABBITMQ_PORT"
echo "   Usuario: $RABBITMQ_USER"
echo "   Cola: $QUEUE_NAME"
echo "   VHost: $RABBITMQ_VHOST"
echo ""

# URL de la API de RabbitMQ
STATUS_URL="http://${RABBITMQ_HOST}:${RABBITMQ_PORT}/api/queues/${RABBITMQ_VHOST}/${QUEUE_NAME}"

# Verificar el estado de la cola
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
CONSUMERS=$(echo "$STATUS_RESPONSE" | grep -o '"consumers":[0-9]*' | grep -o '[0-9]*' || echo "0")
MESSAGE_STATS_PUBLISH=$(echo "$STATUS_RESPONSE" | grep -o '"publish":[0-9]*' | grep -o '[0-9]*' | head -1 || echo "0")
MESSAGE_STATS_DELIVER=$(echo "$STATUS_RESPONSE" | grep -o '"deliver":[0-9]*' | grep -o '[0-9]*' | head -1 || echo "0")

echo "📊 Estado de la cola:"
echo "   ┌─────────────────────────────────────┐"
echo "   │ Mensajes totales:        $MESSAGES"
echo "   │ Mensajes listos:         $MESSAGES_READY"
echo "   │ Mensajes sin confirmar:  $MESSAGES_UNACKNOWLEDGED"
echo "   │ Consumidores activos:     $CONSUMERS"
echo "   │ Total publicados:         $MESSAGE_STATS_PUBLISH"
echo "   │ Total entregados:         $MESSAGE_STATS_DELIVER"
echo "   └─────────────────────────────────────┘"
echo ""

if [ "$MESSAGES" -gt 1000 ]; then
    echo "⚠️  ADVERTENCIA: La cola tiene más de 1000 mensajes. Considera purgar la cola."
    echo "   Ejecuta: ./scripts/purge_rabbitmq_queue.sh $QUEUE_NAME"
elif [ "$MESSAGES" -gt 0 ]; then
    echo "ℹ️  La cola tiene $MESSAGES mensajes pendientes"
else
    echo "✅ La cola está vacía"
fi

if [ "$CONSUMERS" -eq 0 ]; then
    echo "⚠️  ADVERTENCIA: No hay consumidores activos. Los mensajes no se están procesando."
fi
