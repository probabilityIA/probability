# Backend Central - Arquitectura Hexagonal por Módulos

## Descripción

API REST principal del proyecto Probability. Gestiona autenticación, módulos de negocio e integraciones con plataformas de e-commerce.

## Stack Tecnológico

- **Go 1.23**
- **Gin** - Framework HTTP
- **GORM** - ORM para PostgreSQL
- **Redis** - Cache y pub/sub
- **RabbitMQ** - Cola de mensajes
- **MinIO/S3** - Almacenamiento de archivos
- **JWT** - Autenticación
- **Zerolog** - Logging estructurado

## Arquitectura Hexagonal

Cada módulo sigue el patrón de **Arquitectura Hexagonal (Ports & Adapters)**:

```
module/
+-- bundle.go              # Punto de entrada - ensambla el módulo
+-- internal/
    +-- domain/            # NÚCLEO - Reglas de negocio
    |   +-- entities.go    # Entidades del dominio
    |   +-- ports.go       # Interfaces (contratos)
    |   +-- dtos.go        # Data Transfer Objects
    |   +-- errors.go      # Errores del dominio
    +-- app/               # APLICACIÓN - Casos de uso
    |   +-- constructor.go # Factory del usecase
    |   +-- *.go           # Implementación de casos de uso
    +-- infra/             # INFRAESTRUCTURA - Adaptadores
        +-- primary/       # Adaptadores de entrada (drivers)
        |   +-- handlers/  # HTTP handlers (Gin)
        +-- secondary/     # Adaptadores de salida (driven)
            +-- repository/  # Implementación DB (GORM)
            +-- client/      # Clientes HTTP externos
            +-- queue/       # Publicadores de mensajes
```

## Flujo de Dependencias

```
Handler (HTTP) -> UseCase (App) -> Repository (Infra)
      v              v                 v
   primary/      domain/ports      secondary/

Las dependencias SIEMPRE apuntan hacia el dominio (inversión de dependencias)
```

## Estructura de Servicios

### Auth (`services/auth/`)

Sistema de autenticación y autorización RBAC:

| Módulo | Descripción |
|--------|-------------|
| `login/` | Autenticación JWT, API keys, cambio de contraseña |
| `users/` | CRUD de usuarios, asignación de roles |
| `roles/` | Gestión de roles y permisos |
| `permissions/` | CRUD de permisos |
| `resources/` | Recursos protegidos |
| `actions/` | Acciones permitidas sobre recursos |
| `bussines/` | Gestión de negocios (multi-tenant) |
| `middleware/` | Middleware de autenticación |

### Modules (`services/modules/`)

Módulos de negocio core:

| Módulo | Descripción |
|--------|-------------|
| `orders/` | Gestión de pedidos, mapeo desde integraciones |
| `products/` | Catálogo de productos |
| `payments/` | Procesamiento de pagos |
| `shipments/` | Gestión de envíos |
| `events/` | Sistema de eventos y notificaciones |
| `dashboard/` | Métricas y estadísticas |
| `orderstatus/` | Estados de pedidos |
| `paymentstatus/` | Estados de pagos |
| `fulfillmentstatus/` | Estados de fulfillment |
| `notification_config/` | Configuración de notificaciones |
| `ai/` | Funcionalidades de IA |

### Integrations (`services/integrations/`)

Conectores con plataformas externas:

| Módulo | Descripción |
|--------|-------------|
| `shopify/` | Integración con Shopify (webhooks, sync) |
| `amazon/` | Marketplace Amazon |
| `meli/` | MercadoLibre |
| `whatsApp/` | Notificaciones WhatsApp |
| `core/` | Lógica compartida de integraciones |
| `events/` | Eventos de integraciones |

## Patrones de Código

### Bundle (Composición del módulo)

```go
// bundle.go - Ensambla todas las capas del módulo
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, cfg env.IConfig) {
    // 1. Infraestructura secundaria (adaptadores de salida)
    repo := repository.New(db, logger)

    // 2. Capa de aplicación (casos de uso)
    useCase := app.New(repo, logger, cfg)

    // 3. Infraestructura primaria (adaptadores de entrada)
    handler := handlers.New(useCase, logger)

    // 4. Registrar rutas HTTP
    handler.RegisterRoutes(router)
}
```

### Ports (Interfaces del dominio)

```go
// domain/ports.go - Contratos que la infraestructura debe implementar
type IRepository interface {
    Create(ctx context.Context, entity *Entity) error
    GetByID(ctx context.Context, id string) (*Entity, error)
    List(ctx context.Context, filters map[string]interface{}) ([]Entity, error)
    Update(ctx context.Context, entity *Entity) error
    Delete(ctx context.Context, id string) error
}
```

### UseCase (Casos de uso)

```go
// app/constructor.go
type Iapp interface {
    DoSomething(ctx context.Context, dto domain.SomeDTO) (*domain.Response, error)
}

type UseCase struct {
    repository domain.IRepository
    log        log.ILogger
}

func New(repo domain.IRepository, log log.ILogger) Iapp {
    return &UseCase{repository: repo, log: log}
}
```

### Handler (HTTP)

```go
// infra/primary/handlers/handler.go
type Handler struct {
    useCase app.Iapp
    log     log.ILogger
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
    router.GET("/items", h.List)
    router.POST("/items", h.Create)
    router.GET("/items/:id", h.GetByID)
}
```

## Shared (`shared/`)

Utilidades compartidas entre módulos:

| Paquete | Descripción |
|---------|-------------|
| `db/` | Conexión PostgreSQL (GORM) |
| `redis/` | Cliente Redis |
| `rabbitmq/` | Cliente RabbitMQ |
| `storage/` | Cliente S3/MinIO |
| `jwt/` | Generación y validación JWT |
| `email/` | Servicio de emails |
| `log/` | Logger estructurado (zerolog) |
| `env/` | Configuración de entorno |
| `errs/` | Errores personalizados |
| `httpclient/` | Cliente HTTP (resty) |

## Entry Point

```
cmd/
+-- main.go                    # Punto de entrada
+-- internal/
    +-- server/
    |   +-- init.go            # Inicialización del servidor
    +-- routes/
    |   +-- router.go          # Router principal
    |   +-- api_routes.go      # Rutas de API
```

## Comandos

```bash
# Ejecutar servidor
go run cmd/main.go

# Build
go build -o main cmd/main.go

# Tests
go test ./...
```

## Variables de Entorno

```env
# Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=probability
PGSSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASS=admin

# S3/MinIO
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=probability

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Server
PORT=8080
GIN_MODE=debug
```

## Convenciones

1. **Interfaces**: Prefijo `I` (ej: `IRepository`, `ILogger`)
2. **DTOs**: Sufijo `DTO` (ej: `CreateUserDTO`, `UserResponseDTO`)
3. **Constructores**: Función `New()` que retorna la interfaz
4. **Contexto**: Siempre pasar `context.Context` como primer parámetro
5. **Errores**: Usar errores tipados del dominio
6. **Logging**: Usar el logger estructurado de `shared/log`
7. **Validación**: Usar `go-playground/validator` con tags en structs

## Debugging

**El debugging en este proyecto se hace mediante logging focalizado, NO con debuggers externos.**

### Sistema de Logging Dual

- **Logs normales** -> Solo consola (`.Info()`, `.Warn()`, `.Error()`)
- **Logs de debugging** -> Solo archivo (`.DebugToFile()`)

### Activación

En `.env`:
```bash
ENABLE_DEBUG_FILE_LOGGING=true
LOG_DIRECTORY=/home/cam/Desktop/probability/back/central/log
LOG_LEVEL=debug
```

### Uso

```go
// Log normal - solo consola
uc.log.Info(ctx).Str("order_id", id).Msg("Processing order")

// Log de debugging - solo archivo (JSON)
uc.log.DebugToFile(ctx).
    Str("order_id", id).
    Interface("config", config).
    Interface("filters", filters).
    Msg("📋 Detailed configuration")
```

### Análisis

```bash
# Ver logs de una orden
grep '"order_id":"ABC123"' log/app-2026-02-09.log | jq .

# Validadores fallidos
grep "DETAILED: Validation failed" log/app-2026-02-09.log | jq .
```

**Ver:** `/back/central/LOG_DEBUGGING.md` para documentación completa
**Reglas:** `/.claude/rules/debugging-logging.md` para convenciones
