#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP="$ROOT/mobile_central"
MOCK="$ROOT/mock_server"
WEB_PORT="${WEB_PORT:-5100}"
MOCK_PORT="${MOCK_PORT:-5199}"
LOG_DIR="/tmp/probability-mobile"
CTL="$LOG_DIR/flutterctl"

mkdir -p "$LOG_DIR"

stop_all() {
  pkill -f "flutter_tools.snapshot run -d web-server" 2>/dev/null || true
  pkill -f "flutter.sh run -d web-server" 2>/dev/null || true
  pkill -f "mock_server/server.js" 2>/dev/null || true
  pkill -f "frontend_server_aot" 2>/dev/null || true
  sleep 2
}

start_mock() {
  MOCK_PORT="$MOCK_PORT" nohup node "$MOCK/server.js" > "$LOG_DIR/mock.log" 2>&1 &
  sleep 1
  curl -sf "http://localhost:$MOCK_PORT/api/v1/dashboard/stats" > /dev/null && echo "mock ok :$MOCK_PORT"
}

start_web() {
  rm -f "$CTL"
  mkfifo "$CTL"
  cd "$APP"
  nohup bash -c "exec 3<>'$CTL'; flutter run -d web-server --web-port $WEB_PORT --web-hostname 0.0.0.0 \
    --dart-define=APP_ENV=development \
    --dart-define=API_BASE_URL=http://localhost:$MOCK_PORT/api/v1 <&3" > "$LOG_DIR/web.log" 2>&1 &
  echo "flutter web arrancando en :$WEB_PORT (log $LOG_DIR/web.log)"
}

case "${1:-start}" in
  start)
    stop_all
    start_mock
    start_web
    ;;
  stop)
    stop_all
    echo "detenido"
    ;;
  restart-app)
    echo "R" > "$CTL"
    echo "hot restart enviado"
    ;;
  reload)
    echo "r" > "$CTL"
    echo "hot reload enviado"
    ;;
  status)
    ss -ltn 2>/dev/null | grep -E "$WEB_PORT|$MOCK_PORT" || echo "sin servicios"
    ;;
  logs)
    tail -n "${2:-40}" "$LOG_DIR/web.log"
    ;;
  mock-logs)
    tail -n "${2:-40}" "$LOG_DIR/mock.log"
    ;;
  *)
    echo "uso: dev.sh {start|stop|restart-app|reload|status|logs|mock-logs}"
    exit 1
    ;;
esac
