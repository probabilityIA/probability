#!/bin/bash

# Script para forzar la purga completa de la cola de RabbitMQ
# Este script elimina y recrea la cola para limpiar todos los mensajes
# ⚠️  ADVERTENCIA: Esto eliminará TODOS los mensajes, incluso los sin confirmar

set -e

# Configuración por defecto
RABBITMQ_HOST="${RABBITMQ_HOST:-localhost}"
RABBITMQ_PORT="${RABBITMQ_PORT:-15672}"
RABBITMQ_USER="${RABBITMQ_USER:-admin}"
RABBITMQ_PASS="${RABBITMQ_PASS:-admin}"
RABBITMQ_VHOST="${RABBITMQ_VHOST:-%2F}"  # %2F es la codificación URL de "/"
QUEUE_NAME="${1:-probability.orders.canonical}"

echo "⚠️  ADVERTENCIA: Este script ELIMINARÁ y RECREARÁ la cola"
echo "   Esto eliminará TODOS los mensajes, incluso los sin confirmar"
echo "   Host: $RABBITMQ_HOST:$RABBITMQ_PORT"
echo "   Usuario: $RABBITMQ_USER"
echo "   Cola: $QUEUE_NAME"
echo "   VHost: $RABBITMQ_VHOST"
echo ""

# Confirmar antes de proceder
read -p "⚠️  ¿Estás ABSOLUTAMENTE seguro? Esto es IRREVERSIBLE (escribe 'SI' para confirmar): " -r
if [[ ! $REPLY == "SI" ]]; then
    echo "❌ Operación cancelada"
    exit 0
fi

# URL de la API de RabbitMQ
QUEUE_URL="http://${RABBITMQ_HOST}:${RABBITMQ_PORT}/api/queues/${RABBITMQ_VHOST}/${QUEUE_NAME}"

# 1. Obtener información de la cola antes de eliminarla
echo "📊 Obteniendo información de la cola..."
QUEUE_INFO=$(curl -s -u "${RABBITMQ_USER}:${RABBITMQ_PASS}" "$QUEUE_URL" || echo "{}")

if echo "$QUEUE_INFO" | grep -q "error"; then
    echo "❌ Error: La cola '$QUEUE_NAME' no existe"
    exit 1
fi

# Extraer propiedades de la cola
DURABLE=$(echo "$QUEUE_INFO" | grep -o '"durable":[^,]*' | grep -o '[^:]*$' | tr -d ' ' || echo "true")
AUTO_DELETE=$(echo "$QUEUE_INFO" | grep -o '"auto_delete":[^,]*' | grep -o '[^:]*$' | tr -d ' ' || echo "false")
ARGUMENTS=$(echo "$QUEUE_INFO" | grep -o '"arguments":{[^}]*}' || echo "{}")

echo "   Durable: $DURABLE"
echo "   Auto-delete: $AUTO_DELETE"
echo ""

# 2. Eliminar la cola
echo "🗑️  Eliminando la cola..."
DELETE_RESPONSE=$(curl -s -X DELETE -u "${RABBITMQ_USER}:${RABBITMQ_PASS}" "$QUEUE_URL" || echo "{}")

if echo "$DELETE_RESPONSE" | grep -q "error"; then
    echo "❌ Error al eliminar la cola: $DELETE_RESPONSE"
    exit 1
fi

echo "✅ Cola eliminada exitosamente"
echo ""

# 3. Esperar un momento
sleep 2

# 4. Recrear la cola con las mismas propiedades
echo "🔄 Recreando la cola..."
# Usar rabbitmqctl si está disponible, o la API
if docker ps | grep -q rabbitmq; then
    echo "   Usando docker exec para recrear la cola..."
    docker exec rabbitmq_local rabbitmqctl declare queue name="$QUEUE_NAME" durable="$DURABLE" auto_delete="$AUTO_DELETE" 2>/dev/null || \
    echo "   Cola se recreará automáticamente cuando el consumidor se conecte"
else
    echo "   La cola se recreará automáticamente cuando el consumidor se conecte"
fi

echo ""
echo "✅ Proceso completado"
echo "   La cola '$QUEUE_NAME' ha sido eliminada y se recreará cuando el consumidor se conecte"
echo "   Todos los mensajes han sido eliminados"
