#!/bin/bash
set -e

# Script para actualizar contenedores en producción de forma segura
# Sin dependencias entre contenedores - cada servicio maneja reconexiones

COMPOSE_DIR="/home/ubuntu/probability/infra/compose-prod"
PROJECT_NAME="compose-prod"

echo "============================================"
echo "  Actualización Segura de Producción"
echo "============================================"
echo ""

echo "🔍 Verificando contenedores actuales..."
sudo podman ps -a --filter "label=io.podman.compose.project=${PROJECT_NAME}" --format "table {{.Names}}\t{{.Status}}"

echo ""
echo "🛑 Deteniendo servicios..."
cd "$COMPOSE_DIR"
sudo podman-compose down 2>/dev/null || true

echo ""
echo "🧹 Limpiando contenedores residuales con dependencias..."
# Usar --depend para eliminar contenedores con dependencias
sudo podman ps -a --filter "label=io.podman.compose.project=${PROJECT_NAME}" --format "{{.ID}}" | \
while read container_id; do
    if [ -n "$container_id" ]; then
        echo "  Eliminando contenedor: $container_id"
        sudo podman rm -f --depend "$container_id" 2>/dev/null || true
    fi
done

echo ""
echo "✅ Verificando limpieza..."
remaining=$(sudo podman ps -a --filter "label=io.podman.compose.project=${PROJECT_NAME}" -q | wc -l)
if [ "$remaining" -eq 0 ]; then
    echo "  Todos los contenedores eliminados correctamente"
else
    echo "  ⚠️  Quedan $remaining contenedores"
    sudo podman ps -a --filter "label=io.podman.compose.project=${PROJECT_NAME}"
fi

echo ""
echo "🔌 Verificando puertos liberados..."
PORTS=("80" "443" "8080" "8081" "3050" "15672")
all_free=true
for port in "${PORTS[@]}"; do
    if sudo ss -tlnp | grep -q ":${port}\s"; then
        echo "  ⚠️  Puerto $port ocupado"
        all_free=false
    else
        echo "  ✅ Puerto $port libre"
    fi
done

if [ "$all_free" = false ]; then
    echo ""
    echo "⚠️  Hay puertos ocupados. ¿Desea continuar? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Actualización cancelada"
        exit 1
    fi
fi

echo ""
echo "🚀 Levantando servicios (sin dependencias)..."
sudo podman-compose up -d

echo ""
echo "⏳ Esperando 20 segundos para que los servicios inicien..."
sleep 20

echo ""
echo "📊 Estado de los contenedores:"
sudo podman ps --filter "label=io.podman.compose.project=${PROJECT_NAME}" --format "table {{.Names}}\t{{.Status}}"

echo ""
echo "🏥 Verificando salud de los servicios..."
echo ""
for service in website_prod redis_prod rabbitmq_prod backend_prod frontend_prod nginx_prod; do
    status=$(sudo podman ps --filter "name=^${service}$" --format "{{.Status}}" 2>/dev/null || echo "NOT FOUND")
    if echo "$status" | grep -q "Up"; then
        echo "  ✅ $service: $status"
    else
        echo "  ❌ $service: $status"
    fi
done

echo ""
echo "🌐 Verificando conectividad..."
echo "  Desde localhost:"
if curl -s -I http://localhost | grep -q "HTTP"; then
    echo "    ✅ http://localhost responde"
else
    echo "    ❌ http://localhost no responde"
fi

echo ""
echo "  Desde Internet (puede tardar):"
if curl -s -I https://app.probabilityia.com.co --max-time 5 | grep -q "HTTP"; then
    echo "    ✅ https://app.probabilityia.com.co responde"
else
    echo "    ⚠️  https://app.probabilityia.com.co no responde (puede tomar tiempo)"
fi

echo ""
echo "============================================"
echo "  ✅ Actualización completada"
echo "============================================"
echo ""
echo "Para verificar logs de un servicio:"
echo "  sudo podman logs -f backend_prod"
echo ""
echo "Para verificar desde Internet:"
echo "  curl -I https://app.probabilityia.com.co"
echo ""
