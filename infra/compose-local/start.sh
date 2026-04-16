#!/bin/bash

# Script para gestionar servicios de desarrollo local

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

case "${1:-up}" in
  up|start)
    echo "🚀 Levantando servicios de desarrollo local..."
    docker-compose up -d
    echo ""
    echo "✅ Servicios iniciados:"
    echo "   📊 PostgreSQL:    localhost:5432"
    echo "   🔴 Redis:         localhost:6379"
    echo "   📦 MinIO:         http://localhost:9000 (API)"
    echo "   🖥️  MinIO Console: http://localhost:9001"
    echo "   🐰 RabbitMQ:      localhost:5672 (AMQP)"
    echo "   🎛️  RabbitMQ UI:   http://localhost:15672"
    echo ""
    echo "💡 Credenciales por defecto:"
    echo "   PostgreSQL: postgres/postgres"
    echo "   MinIO:      minioadmin/minioadmin"
    echo "   RabbitMQ:   admin/admin"
    echo ""
    echo "📋 Para ver logs: docker-compose logs -f"
    ;;
  
  down|stop)
    echo "🛑 Deteniendo servicios..."
    docker-compose down
    echo "✅ Servicios detenidos"
    ;;
  
  restart)
    echo "🔄 Reiniciando servicios..."
    docker-compose restart
    echo "✅ Servicios reiniciados"
    ;;
  
  logs)
    docker-compose logs -f "${2:-}"
    ;;
  
  status)
    echo "📊 Estado de los servicios:"
    docker-compose ps
    ;;
  
  clean)
    echo "🧹 Limpiando servicios y volúmenes (⚠️  esto borra todos los datos)..."
    read -p "¿Estás seguro? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
      docker-compose down -v
      echo "✅ Limpieza completada"
    else
      echo "❌ Operación cancelada"
    fi
    ;;
  
  *)
    echo "Uso: $0 {up|down|restart|logs|status|clean}"
    echo ""
    echo "Comandos:"
    echo "  up       - Levantar servicios (default)"
    echo "  down     - Detener servicios"
    echo "  restart  - Reiniciar servicios"
    echo "  logs     - Ver logs (opcional: nombre del servicio)"
    echo "  status   - Ver estado de servicios"
    echo "  clean    - Detener y eliminar volúmenes (⚠️  borra datos)"
    exit 1
    ;;
esac
