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

# El entrypoint de nginx corre como root y puede crear active.conf; sin este
# chown el deploy (usuario ubuntu) no puede reescribirlo y aborta a mitad de
# camino, dejando el color nuevo arriba pero sin trafico. Paso exactamente eso
# el 2026-09-02 y tumbo el sitio.
bg_init() {
  sudo mkdir -p "$UPSTREAMS_DIR"
  sudo chown -R ubuntu:ubuntu "$UPSTREAMS_DIR"

  if [ ! -f "$ACTIVE_FILE" ]; then
    echo "🎨 Sin archivo de upstreams, tomando el color que este corriendo"
    bg_write_upstreams "$(bg_running_color back)" "$(bg_running_color front)"
    return
  fi

  # Si el archivo nombra un color que no existe, nginx entra en bucle de panico.
  # Se corrige solo apuntando al que si esta corriendo.
  back=$(bg_file_color back)
  front=$(bg_file_color front)
  fixed=0
  if ! bg_color_running back "$back"; then
    back=$(bg_running_color back); fixed=1
  fi
  if ! bg_color_running front "$front"; then
    front=$(bg_running_color front); fixed=1
  fi
  if [ "$fixed" = "1" ]; then
    echo "⚠️  active.conf apuntaba a un color sin contenedor, corrigiendo"
    bg_write_upstreams "$back" "$front"
  fi
}

# bg_container_prefix back|front
bg_container_prefix() {
  if [ "$1" = "back" ]; then echo central_reserve_prod; else echo frontend_prod; fi
}

# bg_color_running back|front <color>
bg_color_running() {
  docker ps --format '{{.Names}}' | grep -qx "$(bg_container_prefix "$1")_$2"
}

# bg_running_color back|front -> color con contenedor arriba (blue si ninguno)
bg_running_color() {
  for c in blue green; do
    if bg_color_running "$1" "$c"; then echo "$c"; return; fi
  done
  echo blue
}

# bg_file_color back|front -> lo que dice el archivo, sin validar
bg_file_color() {
  color=$(grep -o "$1-central-[a-z]*" "$ACTIVE_FILE" 2>/dev/null | head -1 | sed "s/$1-central-//")
  case "$color" in
    blue|green) echo "$color" ;;
    *) echo blue ;;
  esac
}

# bg_active_color back|front -> el color que realmente esta sirviendo
bg_active_color() {
  color=$(bg_file_color "$1")
  if bg_color_running "$1" "$color"; then
    echo "$color"
  else
    bg_running_color "$1"
  fi
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
