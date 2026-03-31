# Monitoring Service

Servicio independiente de monitoreo de la infraestructura Docker de Probability. Permite visualizar el estado de todos los contenedores, ver logs en tiempo real, ejecutar acciones (restart/stop/start), y monitorear recursos del servidor (CPU, RAM, disco).

Funciona **fuera de Nginx** con puertos propios expuestos directamente al host, para seguir operativo incluso si el resto del stack se cae.

## Arquitectura

```
Browser ──► Next.js (:3002) ──SSR fetch──► Go API (:3070) ──► Docker Socket
                │                              │                 /var/run/docker.sock
                │◄──────── SSE (logs) ─────────│
                │                              │──► PostgreSQL (auth)
                │                              │──► /proc (system stats)
```

### Stack

| Componente | Tecnologia | Puerto |
|------------|-----------|--------|
| **monitoring-api** | Go 1.25 + Gin + Docker SDK | 3070 |
| **monitoring-web** | Next.js 16.1 + React 19 + TailwindCSS 4 | 3002 |

### Acceso en produccion

```
http://ec2-3-224-189-33.compute-1.amazonaws.com:3002
```

Puertos restringidos por Security Group a IPs autorizadas.

---

## Backend (Go API)

### Arquitectura Hexagonal

```
back/monitoring/
├── cmd/main.go                          # Entry point
├── Dockerfile                           # Multi-stage build (golang → alpine)
├── .env                                 # Config local
└── internal/
    ├── domain/
    │   ├── entities/
    │   │   ├── container.go             # Container, ContainerStats, SystemStats, LogEntry, ComposeService
    │   │   └── user.go                  # MonitoringUser (auth via DB)
    │   ├── dtos/
    │   │   ├── auth.go                  # LoginRequest, LoginResponse
    │   │   └── container.go             # ContainerActionRequest, LogStreamRequest
    │   ├── ports/ports.go               # IDockerService, IUserRepository, IUseCase
    │   └── errors/errors.go             # Custom error types
    ├── app/
    │   ├── constructor.go               # UseCase constructor
    │   ├── login.go                     # Auth: bcrypt + JWT (24h)
    │   ├── list_containers.go           # Listar contenedores filtrados por compose project
    │   ├── get_container.go             # Detalle de un contenedor
    │   ├── get_stats.go                 # Stats CPU/RAM de un contenedor
    │   ├── get_system_stats.go          # Stats del servidor host (CPU/RAM/Disk)
    │   ├── container_action.go          # restart, stop, start
    │   ├── stream_logs.go              # Streaming de logs via Docker SDK
    │   └── get_compose_services.go      # Listar servicios del compose
    └── infra/
        ├── primary/handlers/
        │   ├── constructor.go           # Handler + IHandler interface
        │   ├── routes.go               # Registro de rutas
        │   ├── middleware.go            # JWT middleware (HS256)
        │   ├── login_handler.go         # POST /api/v1/auth/login
        │   ├── list_containers_handler.go
        │   ├── get_container_handler.go
        │   ├── get_stats_handler.go
        │   ├── get_system_stats_handler.go
        │   ├── get_logs_handler.go
        │   ├── stream_logs_handler.go   # SSE endpoint
        │   ├── container_action_handler.go
        │   ├── get_compose_services_handler.go
        │   ├── health_handler.go
        │   ├── request/                 # Request DTOs con tags
        │   ├── response/               # Response DTOs con tags
        │   └── mappers/                # Domain ↔ HTTP mappers
        └── secondary/
            ├── docker/
            │   ├── constructor.go       # Docker SDK client
            │   ├── containers.go        # List, Inspect, filtrado por label
            │   ├── actions.go           # Restart, Stop, Start
            │   ├── logs.go              # Log streaming via Docker API
            │   ├── stats.go             # Container CPU/RAM/Network stats
            │   └── system_stats.go      # Host stats via /proc y syscall
            └── repository/
                ├── constructor.go       # GORM PostgreSQL
                └── user.go              # GetUserByEmail para auth
```

### Endpoints

```
POST   /api/v1/auth/login              # Login (email + password)
GET    /api/v1/auth/verify              # Verificar token JWT

GET    /api/v1/containers               # Listar contenedores
GET    /api/v1/containers/:id           # Detalle de un contenedor
GET    /api/v1/containers/:id/stats     # CPU/RAM/Network del contenedor
GET    /api/v1/containers/:id/logs      # Logs historicos (?tail=100)
GET    /api/v1/containers/:id/logs/stream  # SSE - logs en tiempo real

POST   /api/v1/containers/:id/restart   # Reiniciar contenedor
POST   /api/v1/containers/:id/stop      # Detener contenedor
POST   /api/v1/containers/:id/start     # Iniciar contenedor

GET    /api/v1/compose/services          # Listar servicios del compose
GET    /api/v1/system/stats              # CPU/RAM/Disk del servidor host

GET    /health                           # Health check
```

### Autenticacion

- Login contra tabla `users` de la DB compartida de Probability
- Solo usuarios con `scope_id = 1` (platform admin) pueden acceder
- JWT propio con secret independiente (`MONITORING_JWT_SECRET`)
- Token expira en 24 horas
- Middleware JWT en todas las rutas excepto `/health` y `/api/v1/auth/login`

### Docker SDK

- Se conecta via `/var/run/docker.sock` (montado como volume read-only)
- Filtra contenedores por label `com.docker.compose.project` (configurable via `COMPOSE_PROJECT`)
- Stats de contenedores via Docker Stats API (CPU, RAM, Network I/O)
- Stats del host via `/proc/stat`, `/proc/meminfo`, y `syscall.Statfs`

---

## Frontend (Next.js)

### Arquitectura Hexagonal adaptada

```
front/monitoring/
├── Dockerfile                           # Multi-stage (node → standalone)
├── middleware.ts                        # Auth middleware (redirect to /login)
├── src/
│   ├── app/
│   │   ├── layout.tsx                   # Root layout (dark theme)
│   │   ├── page.tsx                     # Redirect → /dashboard
│   │   ├── globals.css                  # Cyberpunk theme + animations
│   │   ├── login/page.tsx               # Client Component - login form
│   │   ├── dashboard/
│   │   │   ├── page.tsx                 # SSR - architecture view
│   │   │   └── [id]/page.tsx            # SSR - container detail + logs
│   │   └── api/
│   │       ├── logs/[id]/route.ts       # SSE proxy → Go API
│   │       ├── stats/[id]/route.ts      # Stats proxy → Go API
│   │       └── system/route.ts          # System stats proxy → Go API
│   ├── services/monitoring/
│   │   ├── domain/
│   │   │   ├── types.ts                 # Container, Stats, SystemStats, etc.
│   │   │   └── ports.ts                 # IMonitoringRepository
│   │   ├── infra/
│   │   │   ├── repository/api-repository.ts  # Fetch wrapper al Go API
│   │   │   └── actions/index.ts         # Server Actions (login, logout, restart, stop, start)
│   │   └── ui/
│   │       ├── components/
│   │       │   ├── ArchitectureView.tsx  # Vista por grupos de servicios con conexiones
│   │       │   ├── ContainerCard.tsx     # Card de un contenedor (estado, uptime, puertos)
│   │       │   ├── ContainerGrid.tsx     # Grid simple de cards
│   │       │   ├── ContainerDetail.tsx   # Info detallada de un contenedor
│   │       │   ├── ActionButtons.tsx     # Restart/Stop/Start con feedback
│   │       │   ├── StatsBar.tsx          # CPU/RAM bars de un contenedor (polling 5s)
│   │       │   ├── SystemStatsBar.tsx    # CPU/RAM/Disk del servidor (polling 5s)
│   │       │   ├── LogViewer.tsx         # Terminal con SSE, scanline, colorized logs
│   │       │   └── Header.tsx           # Nav con logo + logout
│   │       └── hooks/
│   │           ├── useLogStream.ts       # EventSource hook para SSE
│   │           ├── useContainerStats.ts  # Polling stats de un contenedor
│   │           └── useSystemStats.ts     # Polling stats del servidor
│   └── shared/
│       ├── auth/middleware.ts            # JWT cookie check + redirect
│       └── lib/api.ts                   # Fetch wrapper con token
```

### Paginas

| Ruta | Tipo | Descripcion |
|------|------|-------------|
| `/login` | Client Component | Login con email + password |
| `/dashboard` | SSR (force-dynamic) | Vista de arquitectura con todos los servicios agrupados |
| `/dashboard/[id]` | SSR (force-dynamic) | Detalle del contenedor + stats + logs en vivo |

### Dashboard - Vista de Arquitectura

El dashboard muestra los servicios agrupados por su topologia real:

```
┌─────────────────────────────────────────────┐
│  GATEWAY (nginx)                            │
│  Reverse proxy → routes traffic             │
└──────────┬─────────────────┬────────────────┘
           │ HTTPS → :8080   │ HTTPS → :8081
┌──────────▼──────────┐  ┌───▼──────────────┐
│  CENTRAL APP        │  │  WEBSITE          │
│  front-central      │  │  front-website    │
│     ↕ HTTP :3050    │  │  (static, no back)│
│  back-central       │  └──────────────────┘
└──────────┬──────────┘
           │ TCP :6379 / AMQP :5672
┌──────────▼──────────┐  ┌──────────────────┐
│  DATA & MESSAGING   │  │  TESTING          │
│  redis + rabbitmq   │  │  front-testing    │
└─────────────────────┘  │     ↕ HTTP :9092  │
                         │  back-testing     │
┌─────────────────────┐  └──────────────────┘
│  MONITORING          │
│  monitoring-web      │
│     ↕ HTTP :3070     │
│  monitoring-api      │
└──────────────────────┘
```

Cada grupo tiene:
- **Color propio** para identificacion visual rapida
- **Conexiones animadas internas** (front ↔ back) con protocolo y puerto
- **Conexiones animadas entre grupos** mostrando dependencias reales

### Diseno Visual

- **Tema:** Dark cyberpunk (#0a0a0f fondo, neon accents)
- **Animaciones:** Pulsing dots para status, flow animations para conexiones, scanline overlay en logs
- **LogViewer:** Terminal-style con colorizado automatico (error=rojo, warn=amber, info=cyan)
- **Stats:** Barras de progreso con gradientes neon y glow effects
- **Responsive:** Mobile-friendly con breakpoints sm/lg

---

## Deployment

### Docker Compose (produccion)

Ambos servicios estan en `infra/compose-prod/docker-compose.yaml`:

```yaml
monitoring-api:
    image: monitoring-api:latest        # Build local en el servidor ARM64
    ports: ["3070:3070"]
    volumes: ["/var/run/docker.sock:/var/run/docker.sock:ro"]
    environment:
        COMPOSE_PROJECT: "compose-prod"  # Filtra contenedores por este label
        JWT_SECRET: "${MONITORING_JWT_SECRET}"
        DB_HOST/DB_USER/DB_PASS...      # Auth contra DB compartida

monitoring-web:
    image: monitoring-web:latest
    ports: ["3002:3002"]
    environment:
        MONITORING_API_URL: "http://monitoring-api:3070"
    depends_on:
        monitoring-api: { condition: service_healthy }
```

### Build en servidor

```bash
# Desde el servidor EC2 (ARM64)
cd /home/ubuntu/probability-src/back/monitoring
docker build -t monitoring-api:latest .

cd /home/ubuntu/probability-src/front/monitoring
docker build -t monitoring-web:latest .

# Recrear contenedores
cd /home/ubuntu/probability
docker compose up -d monitoring-api monitoring-web
```

### Security Groups

Puertos 3002 y 3070 restringidos a IPs autorizadas en el Security Group `sg-03816f3607edc744b`.

---

## Variables de Entorno

### Backend (.env)

| Variable | Descripcion | Default |
|----------|-------------|---------|
| `HTTP_PORT` | Puerto del servidor | `3070` |
| `JWT_SECRET` | Secret para firmar JWTs | (requerido) |
| `COMPOSE_PROJECT` | Label de Docker Compose para filtrar | `probability` |
| `GIN_MODE` | Modo de Gin (debug/release) | `debug` |
| `LOG_LEVEL` | Nivel de log | `debug` |
| `DB_HOST` | Host de PostgreSQL | `localhost` |
| `DB_PORT` | Puerto de PostgreSQL | `5433` |
| `DB_NAME` | Base de datos | `postgres` |
| `DB_USER` | Usuario | `postgres` |
| `DB_PASS` | Password | `postgres` |
| `PGSSLMODE` | SSL mode | `disable` |

### Frontend (.env)

| Variable | Descripcion | Default |
|----------|-------------|---------|
| `MONITORING_API_URL` | URL interna del Go API (server-side) | `http://localhost:3070` |
| `NEXT_PUBLIC_MONITORING_API_URL` | URL del Go API (client-side) | `http://localhost:3070` |
