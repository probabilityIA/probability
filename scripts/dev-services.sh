#!/bin/bash
# ═══════════════════════════════════════════════════════════════
# dev-services.sh — Gestor de servicios de desarrollo con tmux
#
# Compatible con bash y zsh
# ═══════════════════════════════════════════════════════════════

set -eo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMUX_SESSION="prob"

BACKEND_DIR="$PROJECT_ROOT/back/central"
FRONTEND_DIR="$PROJECT_ROOT/front/central"
TESTING_DIR="$PROJECT_ROOT/back/testing"
FRONTEND_TESTING_DIR="$PROJECT_ROOT/front/testing"
DOCKER_LOCAL="$PROJECT_ROOT/infra/compose-local"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

get_svc_window() {
    case "$1" in
        backend) echo "backend" ;;
        frontend) echo "frontend" ;;
        testing) echo "testing" ;;
        test-front) echo "test-front" ;;
        mobile-web) echo "mobile-web" ;;
        *) echo "" ;;
    esac
}

get_svc_cmd() {
    case "$1" in
        backend) echo "cd $BACKEND_DIR && mkdir -p log && go run cmd/main.go" ;;
        frontend) echo "cd $FRONTEND_DIR && pnpm dev" ;;
        testing) echo "cd $TESTING_DIR && go run cmd/main.go" ;;
        test-front) echo "cd $FRONTEND_TESTING_DIR && pnpm dev -p 3051" ;;
        mobile-web) echo "cd $PROJECT_ROOT/mobile/mobile_central && flutter run -d chrome --dart-define=APP_ENV=development --dart-define=API_BASE_URL=http://localhost:3050/api/v1" ;;
        *) echo "" ;;
    esac
}

get_svc_port() {
    case "$1" in
        backend) echo "3050" ;;
        frontend) echo "3000" ;;
        testing) echo "9092" ;;
        test-front) echo "3051" ;;
        *) echo "" ;;
    esac
}

ensure_session() {
    if ! tmux has-session -t "$TMUX_SESSION" 2>/dev/null; then
        tmux new-session -d -s "$TMUX_SESSION" -n "main"
        tmux send-keys -t "$TMUX_SESSION:main" "echo '=== Probability Dev Session ===" C-m
    fi
}

window_exists() {
    tmux list-windows -t "$TMUX_SESSION" 2>/dev/null | grep -q " $1 "
}

is_running() {
    local svc="$1"
    local win
    win=$(get_svc_window "$svc")
    [ -n "$win" ] && window_exists "$win"
}

kill_port() {
    local port="$1"
    local pids
    pids=$(lsof -ti ":$port" 2>/dev/null || true)
    if [ -n "$pids" ]; then
        echo -e "  ${YELLOW}Matando procesos en puerto $port: $pids${NC}"
        echo "$pids" | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

start_infra() {
    echo -e "${BLUE}[infra]${NC} Verificando servicios Docker..."
    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "redis_local"; then
        echo -e "${GREEN}[infra]${NC} Ya corriendo"
        return 0
    fi
    echo -e "${BLUE}[infra]${NC} Iniciando PostgreSQL, Redis, RabbitMQ, MinIO..."
    docker-compose -f "$DOCKER_LOCAL/docker-compose.yaml" up -d 2>&1 | tail -5
    echo -e "${BLUE}[infra]${NC} Esperando inicialización (3s)..."
    sleep 3
    echo -e "${GREEN}[infra]${NC} Listo"
}

stop_infra() {
    echo -e "${BLUE}[infra]${NC} Deteniendo servicios Docker..."
    docker-compose -f "$DOCKER_LOCAL/docker-compose.yaml" down 2>&1 | tail -3
    echo -e "${GREEN}[infra]${NC} Detenido"
}

start_service() {
    local svc="$1"
    local win cmd port

    win=$(get_svc_window "$svc")
    cmd=$(get_svc_cmd "$svc")
    port=$(get_svc_port "$svc")

    if [ -z "$win" ] || [ -z "$cmd" ]; then
        echo -e "${RED}Servicio desconocido: $svc${NC}"
        return 1
    fi

    if is_running "$svc"; then
        echo -e "${YELLOW}[$svc]${NC} Ya está corriendo en ventana tmux '$win'"
        return 0
    fi

    if [ -n "$port" ]; then
        kill_port "$port"
    fi

    ensure_session

    echo -e "${BLUE}[$svc]${NC} Iniciando en ventana tmux '$win'..."
    tmux new-window -t "$TMUX_SESSION" -n "$win"
    tmux send-keys -t "$TMUX_SESSION:$win" "$cmd" C-m

    echo -e "${GREEN}[$svc]${NC} Iniciado${port:+ (puerto $port)}"
}

stop_service() {
    local svc="$1"
    local win port

    win=$(get_svc_window "$svc")
    port=$(get_svc_port "$svc")

    if [ -z "$win" ]; then
        echo -e "${RED}Servicio desconocido: $svc${NC}"
        return 1
    fi

    if ! is_running "$svc"; then
        echo -e "${YELLOW}[$svc]${NC} No está corriendo"
        if [ -n "$port" ]; then
            kill_port "$port"
        fi
        return 0
    fi

    echo -e "${BLUE}[$svc]${NC} Deteniendo..."

    tmux send-keys -t "$TMUX_SESSION:$win" C-c 2>/dev/null || true
    sleep 2

    tmux kill-window -t "$TMUX_SESSION:$win" 2>/dev/null || true

    if [ -n "$port" ]; then
        kill_port "$port"
    fi

    echo -e "${GREEN}[$svc]${NC} Detenido"
}

restart_service() {
    local svc="$1"
    stop_service "$svc"
    sleep 1
    start_service "$svc"
}

read_logs() {
    local svc="$1"
    local lines="${2:-80}"
    local win

    win=$(get_svc_window "$svc")

    if [ -z "$win" ]; then
        echo -e "${RED}Servicio desconocido: $svc${NC}"
        return 1
    fi

    if ! is_running "$svc"; then
        echo -e "${YELLOW}[$svc]${NC} No está corriendo, no hay logs activos"
        return 1
    fi

    tmux capture-pane -t "$TMUX_SESSION:$win" -p -S "-$lines" 2>/dev/null
}

show_status() {
    echo ""
    echo -e "${CYAN}═══ Probability Dev Services ═══${NC}"
    echo ""

    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "redis_local"; then
        echo -e "  ${GREEN}●${NC} infra        Docker (PG:5433 Redis:6379 RMQ:5672 MinIO:9000)"
    else
        echo -e "  ${RED}○${NC} infra        Docker (detenido)"
    fi

    for svc in backend frontend testing test-front mobile-web; do
        local port port_info
        port=$(get_svc_port "$svc")
        port_info=""
        [ -n "$port" ] && port_info=" :$port"

        if is_running "$svc"; then
            if [ -n "$port" ] && lsof -i ":$port" &>/dev/null; then
                echo -e "  ${GREEN}●${NC} ${svc}$(printf '%*s' $((13 - ${#svc})) '')Corriendo${port_info}"
            else
                echo -e "  ${YELLOW}◐${NC} ${svc}$(printf '%*s' $((13 - ${#svc})) '')Iniciando...${port_info}"
            fi
        else
            echo -e "  ${RED}○${NC} ${svc}$(printf '%*s' $((13 - ${#svc})) '')Detenido${port_info}"
        fi
    done

    echo ""

    if tmux has-session -t "$TMUX_SESSION" 2>/dev/null; then
        echo -e "  ${CYAN}tmux session:${NC} $TMUX_SESSION"
        echo -e "  ${CYAN}ventanas:${NC} $(tmux list-windows -t "$TMUX_SESSION" -F '#{window_name}' 2>/dev/null | tr '\n' ' ')"
    else
        echo -e "  ${YELLOW}Sin sesión tmux activa${NC}"
    fi
    echo ""
}

kill_zombies() {
    echo -e "${BLUE}Buscando procesos zombie...${NC}"
    local found=0

    local go_pids
    go_pids=$(pgrep -f "go run cmd/main.go" 2>/dev/null || true)
    if [ -n "$go_pids" ]; then
        echo -e "  ${YELLOW}Go run:${NC} $go_pids"
        echo "$go_pids" | xargs kill -9 2>/dev/null || true
        found=1
    fi

    local next_pids
    next_pids=$(pgrep -f "next dev" 2>/dev/null || true)
    if [ -n "$next_pids" ]; then
        echo -e "  ${YELLOW}Next.js:${NC} $next_pids"
        echo "$next_pids" | xargs kill -9 2>/dev/null || true
        found=1
    fi

    for port in 3000 3050 3051 9090 9091 9092; do
        local port_pids
        port_pids=$(lsof -ti ":$port" 2>/dev/null || true)
        if [ -n "$port_pids" ]; then
            echo -e "  ${YELLOW}Puerto $port:${NC} $port_pids"
            echo "$port_pids" | xargs kill -9 2>/dev/null || true
            found=1
        fi
    done

    if [ "$found" -eq 0 ]; then
        echo -e "${GREEN}No hay procesos zombie${NC}"
    else
        echo -e "${GREEN}Limpieza completa${NC}"
    fi
}

show_ports() {
    echo ""
    echo -e "${CYAN}═══ Puertos del proyecto ═══${NC}"
    echo ""

    for port in 3000 3050 3051 5433 5672 6379 9000 9001 9090 9091 9092 15672; do
        local pid proc
        pid=$(lsof -ti ":$port" 2>/dev/null | head -1 || true)
        if [ -n "$pid" ]; then
            proc=$(ps -p "$pid" -o comm= 2>/dev/null || echo "?")
            echo -e "  ${GREEN}●${NC} :$port  $proc (pid $pid)"
        else
            echo -e "  ${RED}○${NC} :$port  libre"
        fi
    done
    echo ""
}

CMD="${1:-}"
SVC="${2:-}"
ARG3="${3:-}"

case "$CMD" in
    start)
        if [ "$SVC" = "all" ]; then
            start_infra
            start_service "backend"
            start_service "frontend"
        elif [ "$SVC" = "infra" ]; then
            start_infra
        elif [ -n "$SVC" ]; then
            start_service "$SVC"
        else
            echo "Uso: $0 start <servicio|all>"
            echo "Servicios: infra, backend, frontend, testing, test-front, mobile-web, all"
            exit 1
        fi
        ;;
    stop)
        if [ "$SVC" = "all" ]; then
            for s in mobile-web test-front testing frontend backend; do
                stop_service "$s" 2>/dev/null || true
            done
            stop_infra
            tmux kill-session -t "$TMUX_SESSION" 2>/dev/null || true
        elif [ "$SVC" = "infra" ]; then
            stop_infra
        elif [ -n "$SVC" ]; then
            stop_service "$SVC"
        else
            echo "Uso: $0 stop <servicio|all>"
            exit 1
        fi
        ;;
    restart)
        if [ -z "$SVC" ]; then
            echo "Uso: $0 restart <servicio>"
            exit 1
        fi
        if [ "$SVC" = "infra" ]; then
            stop_infra
            start_infra
        else
            restart_service "$SVC"
        fi
        ;;
    status)
        show_status
        ;;
    logs)
        if [ -z "$SVC" ]; then
            echo "Uso: $0 logs <servicio> [lineas]"
            exit 1
        fi
        read_logs "$SVC" "${ARG3:-80}"
        ;;
    tail)
        if [ -z "$SVC" ]; then
            echo "Uso: $0 tail <servicio>"
            exit 1
        fi
        read_logs "$SVC" 40
        ;;
    kill-zombies)
        kill_zombies
        ;;
    ports)
        show_ports
        ;;
    help|--help|-h|"")
        echo ""
        echo -e "${CYAN}dev-services.sh${NC} — Gestor de servicios de desarrollo"
        echo ""
        echo "Comandos:"
        echo "  start <svc|all>     Iniciar servicio(s)"
        echo "  stop <svc|all>      Detener servicio(s)"
        echo "  restart <svc>       Reiniciar servicio"
        echo "  status              Estado de todos los servicios"
        echo "  logs <svc> [N]      Leer últimas N líneas (default: 80)"
        echo "  tail <svc>          Log resumido (40 líneas)"
        echo "  kill-zombies        Matar procesos huérfanos"
        echo "  ports               Ver puertos en uso"
        echo ""
        echo "Servicios:"
        echo "  infra       Docker (PostgreSQL, Redis, RabbitMQ, MinIO)"
        echo "  backend     Go API (puerto 3050)"
        echo "  frontend    Next.js dashboard (puerto 3000)"
        echo "  testing     Go testing server (puertos 9090-9092)"
        echo "  test-front  Next.js testing UI (puerto 3051)"
        echo "  mobile-web  Flutter web en Chrome"
        echo "  all         infra + backend + frontend"
        echo ""
        ;;
    *)
        echo -e "${RED}Comando desconocido: $CMD${NC}"
        echo "Usa '$0 help' para ver los comandos disponibles"
        exit 1
        ;;
esac
