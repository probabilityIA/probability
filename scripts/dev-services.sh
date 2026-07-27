#!/bin/bash
# ═══════════════════════════════════════════════════════════════
# dev-services.sh — Gestor de servicios de desarrollo con tmux
#
# Compatible con bash y zsh
# ═══════════════════════════════════════════════════════════════

set -eo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMUX_SESSION="prob"
NEXT_HEAP_MB="${NEXT_HEAP_MB:-6144}"

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
        frontend) echo "cd $FRONTEND_DIR && NODE_OPTIONS='--max-old-space-size=$NEXT_HEAP_MB' pnpm dev" ;;
        testing) echo "cd $TESTING_DIR && go run cmd/main.go" ;;
        test-front) echo "cd $FRONTEND_TESTING_DIR && NODE_OPTIONS='--max-old-space-size=$NEXT_HEAP_MB' pnpm dev -p 3051" ;;
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
    tmux list-windows -t "$TMUX_SESSION" -F "#{window_name}" 2>/dev/null | grep -qx "$1"
}

is_running() {
    local svc="$1"
    local win
    win=$(get_svc_window "$svc")
    [ -n "$win" ] && window_exists "$win"
}

MEM_WARN_MB=7000

descendants() {
    local pid="$1" child
    for child in $(pgrep -P "$pid" 2>/dev/null || true); do
        descendants "$child"
        echo "$child"
    done
}

kill_tree() {
    local pid="$1" signal="${2:--TERM}" p
    [ -z "$pid" ] && return 0
    for p in $(descendants "$pid") "$pid"; do
        kill "$signal" "$p" 2>/dev/null || true
    done
}

port_listeners() {
    local port="$1" pids
    pids=$(lsof -ti "tcp:$port" -sTCP:LISTEN 2>/dev/null || true)
    if [ -z "$pids" ]; then
        pids=$(ss -ltnp 2>/dev/null | grep -E "[:.]$port[[:space:]]" | grep -oP 'pid=\K[0-9]+' | sort -u || true)
    fi
    echo "$pids"
}

kill_port() {
    local port="$1"
    local pids pid

    pids=$(port_listeners "$port")
    if [ -z "$pids" ]; then
        return 0
    fi

    echo -e "  ${YELLOW}Puerto $port ocupado por: $pids${NC}"
    for pid in $pids; do
        kill_tree "$pid" -TERM
    done
    sleep 2

    pids=$(port_listeners "$port")
    if [ -n "$pids" ]; then
        echo -e "  ${YELLOW}No cedieron, forzando: $pids${NC}"
        for pid in $pids; do
            kill_tree "$pid" -KILL
        done
        sleep 1
    fi

    pids=$(port_listeners "$port")
    if [ -n "$pids" ]; then
        fuser -k "$port/tcp" 2>/dev/null || true
        sleep 1
    fi

    if [ -n "$(port_listeners "$port")" ]; then
        echo -e "  ${RED}Puerto $port sigue ocupado por $(port_listeners "$port")${NC}"
        return 1
    fi
}

proc_mem_mb() {
    local pid="$1" total=0 rss p
    [ -z "$pid" ] && { echo 0; return; }
    for p in "$pid" $(descendants "$pid"); do
        rss=$(ps -o rss= -p "$p" 2>/dev/null | tr -d ' ' || true)
        [ -n "$rss" ] && total=$((total + rss))
    done
    echo $((total / 1024))
}

svc_pane_pid() {
    local win="$1"
    tmux list-panes -t "$TMUX_SESSION:$win" -F '#{pane_pid}' 2>/dev/null | head -1 || true
}

start_infra() {
    echo -e "${BLUE}[infra]${NC} Verificando servicios Docker..."
    if docker ps --format '{{.Names}}' 2>/dev/null | grep -q "redis_local"; then
        echo -e "${GREEN}[infra]${NC} Ya corriendo"
        return 0
    fi
    echo -e "${BLUE}[infra]${NC} Iniciando PostgreSQL, Redis, RabbitMQ..."
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

    local pane_pid
    pane_pid=$(svc_pane_pid "$win")

    tmux send-keys -t "$TMUX_SESSION:$win" C-c 2>/dev/null || true
    sleep 2

    if [ -n "$pane_pid" ]; then
        kill_tree "$pane_pid" -TERM
        sleep 1
        kill_tree "$pane_pid" -KILL
    fi

    tmux kill-window -t "$TMUX_SESSION:$win" 2>/dev/null || true

    if [ -n "$port" ]; then
        kill_port "$port" || true
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
        echo -e "  ${GREEN}●${NC} infra        Docker (PG:5433 Redis:6379 RMQ:5672)"
    else
        echo -e "  ${RED}○${NC} infra        Docker (detenido)"
    fi

    for svc in backend frontend testing test-front mobile-web; do
        local port port_info
        port=$(get_svc_port "$svc")
        port_info=""
        [ -n "$port" ] && port_info=" :$port"

        if is_running "$svc"; then
            local win mem mem_info
            win=$(get_svc_window "$svc")
            mem=$(proc_mem_mb "$(svc_pane_pid "$win")")
            mem_info=""
            if [ "$mem" -gt 0 ]; then
                if [ "$mem" -ge "$MEM_WARN_MB" ]; then
                    mem_info="  ${RED}${mem} MB (reinicialo)${NC}"
                else
                    mem_info="  ${CYAN}${mem} MB${NC}"
                fi
            fi

            if [ -n "$port" ] && [ -n "$(port_listeners "$port")" ]; then
                echo -e "  ${GREEN}●${NC} ${svc}$(printf '%*s' $((13 - ${#svc})) '')Corriendo${port_info}${mem_info}"
            else
                echo -e "  ${YELLOW}◐${NC} ${svc}$(printf '%*s' $((13 - ${#svc})) '')Iniciando...${port_info}${mem_info}"
            fi
        else
            echo -e "  ${RED}○${NC} ${svc}$(printf '%*s' $((13 - ${#svc})) '')Detenido${port_info}"
        fi
    done

    local huerfanos
    huerfanos=$(pgrep -f "next-server|next dev|go run cmd/main.go" 2>/dev/null | wc -l)
    if [ "$huerfanos" -gt 4 ]; then
        echo ""
        echo -e "  ${YELLOW}Hay $huerfanos procesos next/go vivos: revisa con '$0 mem'${NC}"
    fi

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
    next_pids=$(pgrep -f "next-server|next dev|next start" 2>/dev/null || true)
    if [ -n "$next_pids" ]; then
        echo -e "  ${YELLOW}Next.js:${NC} $next_pids"
        for pid in $next_pids; do
            kill_tree "$pid" -KILL
        done
        found=1
    fi

    for port in 3000 3050 3051 9090 9091 9092; do
        local port_pids
        port_pids=$(port_listeners "$port")
        if [ -n "$port_pids" ]; then
            echo -e "  ${YELLOW}Puerto $port:${NC} $port_pids"
            for pid in $port_pids; do
                kill_tree "$pid" -KILL
            done
            found=1
        fi
    done

    if [ "$found" -eq 0 ]; then
        echo -e "${GREEN}No hay procesos zombie${NC}"
    else
        echo -e "${GREEN}Limpieza completa${NC}"
    fi
}

guard_services() {
    local limit="${1:-$MEM_WARN_MB}"
    local interval="${2:-60}"

    echo -e "${CYAN}Vigilando frontend: limite ${limit} MB, revision cada ${interval}s${NC}"
    echo -e "${YELLOW}Ctrl-C para detener${NC}"

    while true; do
        if is_running "frontend"; then
            local mem
            mem=$(proc_mem_mb "$(svc_pane_pid frontend)")
            if [ "$mem" -ge "$limit" ]; then
                echo -e "$(date +%T) ${RED}frontend en ${mem} MB, reiniciando${NC}"
                restart_service "frontend"
                sleep 30
            else
                echo -e "$(date +%T) frontend ${mem} MB"
            fi
        else
            echo -e "$(date +%T) frontend detenido"
        fi
        sleep "$interval"
    done
}

show_mem() {
    echo ""
    echo -e "${CYAN}═══ Memoria de los procesos del proyecto ═══${NC}"
    echo ""
    printf "  %-8s %-10s %-9s %s\n" "PID" "RSS(MB)" "TIEMPO" "COMANDO"

    local pids pid rss etime args mb total=0
    pids=$(pgrep -f "next-server|next dev|go run cmd/main.go|pnpm dev|/exe/main" 2>/dev/null || true)
    for pid in $pids; do
        rss=$(ps -o rss= -p "$pid" 2>/dev/null | tr -d ' ' || true)
        [ -z "$rss" ] && continue
        mb=$((rss / 1024))
        total=$((total + mb))
        etime=$(ps -o etime= -p "$pid" 2>/dev/null | tr -d ' ' || echo "?")
        args=$(ps -o args= -p "$pid" 2>/dev/null | cut -c1-45 || echo "?")
        if [ "$mb" -ge "$MEM_WARN_MB" ]; then
            printf "  ${RED}%-8s %-10s %-9s %s${NC}\n" "$pid" "$mb" "$etime" "$args"
        else
            printf "  %-8s %-10s %-9s %s\n" "$pid" "$mb" "$etime" "$args"
        fi
    done

    echo ""
    echo -e "  Total procesos del proyecto: ${CYAN}${total} MB${NC}"
    free -h 2>/dev/null | awk 'NR==2 {printf "  RAM del equipo: usada %s de %s, disponible %s\n", $3, $2, $7}'
    echo ""
    echo -e "  ${YELLOW}Un next-server sobre ${MEM_WARN_MB} MB o con horas encima conviene reiniciarlo:${NC}"
    echo -e "    $0 restart frontend"
    echo ""
}

show_ports() {
    echo ""
    echo -e "${CYAN}═══ Puertos del proyecto ═══${NC}"
    echo ""

    for port in 3000 3050 3051 5433 5672 6379 9090 9091 9092 15672; do
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
    mem|memoria)
        show_mem
        ;;
    guard)
        guard_services "${SVC:-}" "${ARG3:-}"
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
        echo "  mem                 Ver memoria de los procesos del proyecto"
        echo "  guard [MB] [seg]    Vigila el frontend y lo reinicia si se pasa de MB"
        echo ""
        echo "Servicios:"
        echo "  infra       Docker (PostgreSQL, Redis, RabbitMQ)"
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
