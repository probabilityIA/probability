#!/usr/bin/env bash
# Tuneles a la infraestructura de AWS via SSM. No hay puertos publicos ni .pem:
# el RDS solo acepta conexiones desde el EC2, y el EC2 no tiene SSH.
#
#   ./scripts/aws-tunnel.sh start     abre el tunel al RDS en 127.0.0.1:5433
#   ./scripts/aws-tunnel.sh stop      lo cierra
#   ./scripts/aws-tunnel.sh status    dice si esta arriba y prueba la conexion
#   ./scripts/aws-tunnel.sh ensure    lo abre solo si no esta corriendo (idempotente)
#   ./scripts/aws-tunnel.sh shell     shell interactiva en el EC2 de produccion
set -euo pipefail

PROFILE="${AWS_PROFILE:-probability}"
REGION="${AWS_REGION:-us-east-1}"
INSTANCE="${SSM_INSTANCE_ID:-i-0f3284d2a87127e57}"
RDS_HOST="${RDS_HOST:-database-1.capmmoe4cw2e.us-east-1.rds.amazonaws.com}"
LOCAL_PORT="${RDS_LOCAL_PORT:-5433}"
PIDFILE="/tmp/probability-rds-tunnel.pid"
LOGFILE="/tmp/probability-rds-tunnel.log"

is_up() {
  ss -ltn 2>/dev/null | grep -q "127.0.0.1:${LOCAL_PORT}"
}

start() {
  if is_up; then
    echo "El tunel ya esta arriba en 127.0.0.1:${LOCAL_PORT}"
    return 0
  fi
  command -v session-manager-plugin >/dev/null || {
    echo "Falta session-manager-plugin. Instalar:"
    echo "  curl -fsSL https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb -o /tmp/smp.deb && sudo dpkg -i /tmp/smp.deb"
    return 1
  }
  nohup aws ssm start-session \
    --profile "$PROFILE" --region "$REGION" \
    --target "$INSTANCE" \
    --document-name AWS-StartPortForwardingSessionToRemoteHost \
    --parameters "{\"host\":[\"${RDS_HOST}\"],\"portNumber\":[\"5432\"],\"localPortNumber\":[\"${LOCAL_PORT}\"]}" \
    > "$LOGFILE" 2>&1 &
  echo $! > "$PIDFILE"

  for _ in $(seq 1 20); do
    sleep 1
    is_up && { echo "Tunel arriba: 127.0.0.1:${LOCAL_PORT} -> ${RDS_HOST}:5432"; return 0; }
  done
  echo "No abrio el tunel. Log: $LOGFILE"
  tail -5 "$LOGFILE"
  return 1
}

stop() {
  [ -f "$PIDFILE" ] && kill "$(cat "$PIDFILE")" 2>/dev/null || true
  pkill -f "StartPortForwardingSessionToRemoteHost" 2>/dev/null || true
  rm -f "$PIDFILE"
  echo "Tunel cerrado"
}

status() {
  if is_up; then
    echo "Tunel ARRIBA en 127.0.0.1:${LOCAL_PORT}"
    if command -v psql >/dev/null && [ -f back/central/.env ]; then
      local u p n
      u=$(grep -m1 '^DB_USER=' back/central/.env | cut -d= -f2-)
      p=$(grep -m1 '^DB_PASS=' back/central/.env | cut -d= -f2-)
      n=$(grep -m1 '^DB_NAME=' back/central/.env | cut -d= -f2-)
      PGPASSWORD="$p" psql -h 127.0.0.1 -p "$LOCAL_PORT" -U "$u" -d "$n" \
        -tAc "select 'conexion OK -> '||inet_server_addr()" 2>&1 | tail -1
    fi
  else
    echo "Tunel ABAJO. Levantar con: ./scripts/aws-tunnel.sh start"
    return 1
  fi
}

case "${1:-status}" in
  start)  start ;;
  stop)   stop ;;
  status) status ;;
  ensure) is_up || start ;;
  shell)  exec aws ssm start-session --profile "$PROFILE" --region "$REGION" --target "$INSTANCE" ;;
  *) echo "uso: $0 {start|stop|status|ensure|shell}"; exit 1 ;;
esac
