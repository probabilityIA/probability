#!/usr/bin/env bash
# ═══════════════════════════════════════════════════════════════
# dev-services.sh — Gestor de servicios de desarrollo con tmux
#
# Uso:
#   ./scripts/dev-services.sh <comando> [servicio]
#
# Comandos:
#   start <servicio|all>     Iniciar servicio(s) en tmux
#   stop <servicio|all>      Detener servicio(s)
#   restart <servicio>       Reiniciar un servicio
#   status                   Ver estado de todos los servicios
#   logs <servicio> [N]      Leer últimas N líneas de log (default: 80)
#   tail <servicio>          Capturar log en vivo (últimas líneas del panel tmux)
#   kill-zombies             Matar procesos Go/Next.js huérfanos
#   ports                    Ver puertos en uso
#
# Servicios:
#   infra       Docker (PostgreSQL, Redis, RabbitMQ, MinIO)
#   backend     Go API central (puerto 3050)
#   frontend    Next.js dashboard (puerto 3000)
#   testing     Go testing server (puertos 9090-9092)
#   test-front  Next.js testing frontend (puerto 3051)
#   mobile-web  Flutter web en Chrome
#   all         Todos los servicios principales (infra+backend+frontend)
# ═══════════════════════════════════════════════════════════════

set -euo pipefail

PROJECT_ROOT="/home/cam/Desktop/probability"
TMUX_SESSION="prob"

# Directorios
BACKEND_DIR="$PROJECT_ROOT/back/central"
FRONTEND_DIR="$PROJECT_ROOT/front/central"
TESTING_DIR="$PROJECT_ROOT/back/testing"
FRONTEND_TESTING_DIR="$PROJECT_ROOT/front/testing"
MOBILE_DIR="$PROJECT_ROOT/mobile/mobile_central"
DOCKER_LOCAL="$PROJECT_ROOT/infra/compose-local"

# Mapeo servicio → ventana tmux
declare -A SVC_WINDOW=(
    [backend]="backend"
    [frontend]="frontend"
    [testing]="testing"
    [test-front]="test-front"
    [mobile-web]="mobile-web"
)

# Mapeo servicio → comando
declare -A SVC_CMD=(
    [backend]="cd $BACKEND_DIR && go run cmd/main.go"
    [frontend]="cd $FRONTEND_DIR && pnpm dev"
    [testing]="cd $TESTING_DIR && go run cmd/main.go"
    [test-front]="cd $FRONTEND_TESTING_DIR && pnpm dev -p 3051"
    [mobile-web]="cd $MOBILE_DIR && flutter run -d chrome --dart-define=APP_ENV=development --dart-define=API_BASE_URL=http://localhost:3050/api/v1 --dart-define=DEV_EMAIL=\$(grep DEV_EMAIL .env.dev 2>/dev/null | cut -d= -f2) --dart-define=DEV_PASSWORD=\$(grep DEV_PASSWORD .env.dev 2>/dev/null | cut -d= -f2)"
)

# Mapeo servicio → puerto
declare -A SVC_PORT=(
    [backend]="3050"
    [frontend]="3000"
    [testing]="9092"
    [test-front]="3051"
)

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# ─── Helpers ────────────────────────────────────────────

ensure_session() {
    if ! tmux has-session -t "$TMUX_SESSION" 2>/dev/null; then
        tmux new-session -d -s "$TMUX_SESSION" -n "main"
        tmux send-keys -t "$TMUX_SESSION:main" "echo '=== Probability Dev Session ==='" Enter
    fi
}

window_exists() {
    tmux list-windows -t "$TMUX_SESSION" 2>/dev/null | grep -q " $1 "
    return $?
}

# Verifica si un servicio está corriendo revisando si la ventana tmux existe
is_running() {
    local svc="$1"
    local win="${SVC_WINDOW[$svc]:-}"
    [ -z "$win" ] && return 1
    tmux list-windows -t "$TMUX_SESSION" -F '#{window_name}' 2>/dev/null | grep -q "^${win}$"
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

# ─── Infraestructura (Docker) ──────────────────────────

start_infra() {
    echo -e "${BLUE}[infra]${NC} Verificando servicios Docker..."
    if docker ps --format '{{.Names}}' | grep -q "redis_local"; then
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

# ─── Servicios en tmux ──────────────────────────────────

start_service() {
    local svc="$1"
    local win="${SVC_WINDOW[$svc]:-}"
    local cmd="${SVC_CMD[$svc]:-}"
    local port="${SVC_PORT[$svc]:-}"

    if [ -z "$win" ] || [ -z "$cmd" ]; then
        echo -e "${RED}Servicio desconocido: $svc${NC}"
        echo "Servicios válidos: backend, frontend, testing, test-front, mobile-web"
        return 1
    fi

    # Si ya está corriendo, avisar
    if is_running "$svc"; then
        echo -e "${YELLOW}[$svc]${NC} Ya está corriendo en ventana tmux '$win'"
        return 0
    fi

    # Limpiar puerto si está ocupado
    if [ -n "$port" ]; then
        kill_port "$port"
    fi

    ensure_session

    echo -e "${BLUE}[$svc]${NC} Iniciando en ventana tmux '$win'..."
    tmux new-window -t "$TMUX_SESSION" -n "$win"
    tmux send-keys -t "$TMUX_SESSION:$win" "$cmd" Enter

    echo -e "${GREEN}[$svc]${NC} Iniciado${port:+ (puerto $port)}"
}

stop_service() {
    local svc="$1"
    local win="${SVC_WINDOW[$svc]:-}"
    local port="${SVC_PORT[$svc]:-}"

    if [ -z "$win" ]; then
        echo -e "${RED}Servicio desconocido: $svc${NC}"
        return 1
    fi

    if ! is_running "$svc"; then
        echo -e "${YELLOW}[$svc]${NC} No está corriendo"
        # Igualmente limpiar puerto por si quedó zombie
        if [ -n "$port" ]; then
            kill_port "$port"
        fi
        return 0
    fi

    echo -e "${BLUE}[$svc]${NC} Deteniendo..."

    # Enviar Ctrl+C primero para shutdown graceful
    tmux send-keys -t "$TMUX_SESSION:$win" C-c 2>/dev/null || true
    sleep 2

    # Cerrar la ventana tmux
    tmux kill-window -t "$TMUX_SESSION:$win" 2>/dev/null || true

    # Limpiar puerto por si quedó algo
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

# ─── Leer logs ──────────────────────────────────────────

read_logs() {
    local svc="$1"
    local lines="${2:-80}"
    local win="${SVC_WINDOW[$svc]:-}"

    if [ -z "$win" ]; then
        echo -e "${RED}Servicio desconocido: $svc${NC}"
        return 1
    fi

    if ! is_running "$svc"; then
        echo -e "${YELLOW}[$svc]${NC} No está corriendo, no hay logs activos"
        return 1
    fi

    # Capturar el contenido del panel tmux
    tmux capture-pane -t "$TMUX_SESSION:$win" -p -S "-$lines" 2>/dev/null
}

# ─── Estado ─────────────────────────────────────────────

show_status() {
    echo ""
    echo -e "${CYAN}═══ Probability Dev Services ═══${NC}"
    echo ""

    # Infra
    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "redis_local"; then
        echo -e "  ${GREEN}●${NC} infra        Docker (PG:5433 Redis:6379 RMQ:5672 MinIO:9000)"
    else
        echo -e "  ${RED}○${NC} infra        Docker (detenido)"
    fi

    # Servicios
    for svc in backend frontend testing test-front mobile-web; do
        local port="${SVC_PORT[$svc]:-}"
        local port_info=""
        [ -n "$port" ] && port_info=" :$port"

        if is_running "$svc"; then
            # Verificar si el puerto realmente responde
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

    # Sesión tmux
    if tmux has-session -t "$TMUX_SESSION" 2>/dev/null; then
        echo -e "  ${CYAN}tmux session:${NC} $TMUX_SESSION"
        echo -e "  ${CYAN}ventanas:${NC} $(tmux list-windows -t "$TMUX_SESSION" -F '#{window_name}' 2>/dev/null | tr '\n' ' ')"
    else
        echo -e "  ${YELLOW}Sin sesión tmux activa${NC}"
    fi
    echo ""
}

# ─── Kill zombies ───────────────────────────────────────

kill_zombies() {
    echo -e "${BLUE}Buscando procesos zombie...${NC}"
    local found=0

    # Go processes
    local go_pids
    go_pids=$(pgrep -f "go run cmd/main.go" 2>/dev/null || true)
    if [ -n "$go_pids" ]; then
        echo -e "  ${YELLOW}Go run:${NC} $go_pids"
        echo "$go_pids" | xargs kill -9 2>/dev/null || true
        found=1
    fi

    # Compiled Go binaries from tmp
    local go_tmp
    go_tmp=$(pgrep -f "/tmp/go-build.*main" 2>/dev/null || true)
    if [ -n "$go_tmp" ]; then
        echo -e "  ${YELLOW}Go tmp:${NC} $go_tmp"
        echo "$go_tmp" | xargs kill -9 2>/dev/null || true
        found=1
    fi

    # Next.js dev server
    local next_pids
    next_pids=$(pgrep -f "next dev" 2>/dev/null || true)
    if [ -n "$next_pids" ]; then
        echo -e "  ${YELLOW}Next.js:${NC} $next_pids"
        echo "$next_pids" | xargs kill -9 2>/dev/null || true
        found=1
    fi

    # Node processes on our ports
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

# ─── Puertos ────────────────────────────────────────────

show_ports() {
    echo ""
    echo -e "${CYAN}═══ Puertos del proyecto ═══${NC}"
    echo ""

    # Puertos de Docker (containers)
    local docker_ports
    docker_ports=$(docker ps --format '{{.Names}}:{{.Ports}}' 2>/dev/null || true)

    for port in 3000 3050 3051 5433 5672 6379 9000 9001 9090 9091 9092 15672; do
        local pid
        pid=$(lsof -ti ":$port" 2>/dev/null | head -1 || true)
        if [ -n "$pid" ]; then
            local proc
            proc=$(ps -p "$pid" -o comm= 2>/dev/null || echo "?")
            echo -e "  ${GREEN}●${NC} :$port  $proc (pid $pid)"
        elif echo "$docker_ports" | grep -qE "0.0.0.0:($port|[0-9]+-[0-9]+)->.*$port"; then
            local container
            container=$(echo "$docker_ports" | grep -E "0.0.0.0:($port|[0-9]+-[0-9]+)->.*$port" | head -1 | cut -d: -f1)
            echo -e "  ${GREEN}●${NC} :$port  docker ($container)"
        else
            echo -e "  ${RED}○${NC} :$port  libre"
        fi
    done
    echo ""
}

# ─── Main ───────────────────────────────────────────────

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
            # Limpiar sesión tmux si quedó vacía
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
