#!/usr/bin/env sh
# Funciones compartidas de blue-green para los deploys de back-central y
# front-central. Se ejecuta dentro del EC2 (dash, no bash: nada de pipefail
# ni arrays).
#
# La fuente de verdad del color activo es el archivo de upstreams que lee
# nginx. No hay estado en ningun otro lado.

COMPOSE_DIR=/home/ubuntu/probability/infra/compose-prod
UPSTREAMS_DIR="$COMPOSE_DIR/nginx-upstreams"
ACTIVE_FILE="$UPSTREAMS_DIR/active.conf"

bg_init() {
  mkdir -p "$UPSTREAMS_DIR"
  if [ ! -f "$ACTIVE_FILE" ]; then
    echo "🎨 Sin archivo de upstreams, arrancando en blue"
    bg_write_upstreams blue blue
  fi
}

# bg_active_color back|front -> blue|green
bg_active_color() {
  service="$1"
  color=$(grep -o "${service}-central-[a-z]*" "$ACTIVE_FILE" 2>/dev/null | head -1 | sed "s/${service}-central-//")
  case "$color" in
    blue|green) echo "$color" ;;
    *) echo blue ;;
  esac
}

bg_other_color() {
  if [ "$1" = "blue" ]; then echo green; else echo blue; fi
}

# bg_write_upstreams <color_back> <color_front>
bg_write_upstreams() {
  cat > "$ACTIVE_FILE" <<EOF
upstream probability_backend {
    server back-central-$1:3050 max_fails=3 fail_timeout=30s;
}

upstream probability_frontend {
    server front-central-$2:3000 max_fails=3 fail_timeout=30s;
}
EOF
  echo "🎨 Upstreams: backend=$1 frontend=$2"
}

# Recrea nginx si todavia no tiene el montaje de upstreams (primer deploy tras
# migrar a blue-green, o si alguien lo levanto con el compose viejo).
bg_ensure_nginx() {
  if ! docker exec nginx_prod test -d /etc/nginx/upstreams 2>/dev/null; then
    echo "🔧 nginx sin el montaje de upstreams, recreando..."
    cd "$COMPOSE_DIR" || return 1
    docker compose -f docker-compose.yaml up -d --force-recreate nginx
    sleep 5
  fi
}

# bg_wait_healthy <container> <url interna> <timeout segundos>
bg_wait_healthy() {
  container="$1"
  url="$2"
  timeout="${3:-90}"
  waited=0

  echo "🏥 Esperando a que $container responda en $url (max ${timeout}s)..."
  while [ "$waited" -lt "$timeout" ]; do
    state=$(docker inspect "$container" --format '{{.State.Status}}' 2>/dev/null || echo missing)
    if [ "$state" != "running" ] && [ "$state" != "created" ]; then
      echo "❌ $container esta en estado '$state'"
      docker logs --tail 50 "$container" 2>/dev/null || true
      return 1
    fi
    if docker exec "$container" wget -q -O- -T 3 "$url" >/dev/null 2>&1; then
      echo "✅ $container responde despues de ${waited}s"
      return 0
    fi
    sleep 3
    waited=$((waited + 3))
  done

  echo "❌ $container no respondio en ${timeout}s"
  docker logs --tail 50 "$container" 2>/dev/null || true
  return 1
}

bg_reload_nginx() {
  if ! docker exec nginx_prod nginx -t 2>&1; then
    echo "❌ La configuracion de nginx no valida, no se recarga"
    return 1
  fi
  docker exec nginx_prod nginx -s reload
  echo "🔁 nginx recargado (sin cortar conexiones)"
}

# bg_retire <container> - saca de servicio el color viejo tras el switch
bg_retire() {
  old="$1"
  if docker ps -a --format '{{.Names}}' | grep -qx "$old"; then
    echo "⏳ Drenando $old 15s antes de apagarlo..."
    sleep 15
    docker stop -t 15 "$old" >/dev/null 2>&1 || true
    docker rm -f "$old" >/dev/null 2>&1 || true
    echo "🗑️  $old retirado"
  fi
}

# bg_drop_legacy <container> - contenedor con el nombre de antes de blue-green
bg_drop_legacy() {
  legacy="$1"
  if docker ps -a --format '{{.Names}}' | grep -qx "$legacy"; then
    echo "🧹 Eliminando contenedor previo a blue-green: $legacy"
    docker rm -f "$legacy" >/dev/null 2>&1 || true
  fi
}
