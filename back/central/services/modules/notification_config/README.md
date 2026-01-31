# Notification Config Module

## Descripción

Módulo de configuración de notificaciones con **arquitectura jerárquica de tres niveles**:

1. **Notification Types** (Tipos de Notificación) - Canales de comunicación: WhatsApp, SSE, Email, SMS
2. **Notification Event Types** (Eventos de Notificación) - Eventos específicos por tipo: order.created, order.shipped, invoice.created
3. **Business Notification Configs** (Configuraciones de Negocio) - Configuraciones específicas por negocio/integración

## Arquitectura Hexagonal

```
notification_config/
├── bundle.go                    # Ensamblador del módulo
└── internal/
    ├── domain/                  # 🔵 DOMINIO (núcleo puro)
    │   ├── entities/
    │   │   ├── notification_type.go              # Entidad NotificationType
    │   │   ├── notification_event_type.go        # Entidad NotificationEventType (NUEVA)
    │   │   ├── business_notification_config.go   # Entidad refactorizada
    │   │   └── notification_config.go            # IntegrationNotificationConfig (legacy)
    │   ├── dtos/
    │   ├── ports/
    │   │   ├── repository.go    # Interfaces de repositorios
    │   │   └── usecase.go       # Interfaces de casos de uso
    │   └── errors/
    │
    ├── app/                     # 🟢 APLICACIÓN (casos de uso)
    │   ├── constructor.go
    │   ├── create_notification_type.go           # CRUD NotificationType
    │   ├── get_notification_types.go
    │   ├── update_notification_type.go
    │   ├── delete_notification_type.go
    │   ├── create_notification_event_type.go     # CRUD NotificationEventType
    │   ├── get_notification_event_types.go
    │   ├── update_notification_event_type.go
    │   ├── delete_notification_event_type.go
    │   ├── create_notification_config.go         # CRUD BusinessNotificationConfig
    │   ├── list_notification_configs.go
    │   ├── update_notification_config.go
    │   └── delete_notification_config.go
    │
    └── infra/                   # 🔴 INFRAESTRUCTURA
        ├── primary/             # Adaptadores de entrada
        │   └── handlers/
        │       ├── notification_type/           # Handlers para NotificationType
        │       │   ├── constructor.go
        │       │   ├── routes.go
        │       │   ├── create_handler.go
        │       │   ├── list_handler.go
        │       │   ├── get_by_id_handler.go
        │       │   ├── update_handler.go
        │       │   ├── delete_handler.go
        │       │   ├── request/
        │       │   ├── response/
        │       │   └── mappers/
        │       │
        │       ├── notification_event_type/     # Handlers para NotificationEventType
        │       │   ├── constructor.go
        │       │   ├── routes.go
        │       │   ├── create_handler.go
        │       │   ├── list_handler.go
        │       │   ├── get_by_id_handler.go
        │       │   ├── update_handler.go
        │       │   ├── delete_handler.go
        │       │   ├── request/
        │       │   ├── response/
        │       │   └── mappers/
        │       │
        │       └── notification_config/         # Handlers para BusinessNotificationConfig
        │           ├── constructor.go
        │           ├── routes.go
        │           ├── create_handler.go
        │           ├── list_handler.go
        │           ├── get_by_id_handler.go
        │           ├── update_handler.go
        │           ├── delete_handler.go
        │           ├── request/
        │           ├── response/
        │           └── mappers/
        │
        └── secondary/           # Adaptadores de salida
            └── repository/
                ├── constructor.go
                ├── notification_type_repository.go        # Repositorio NotificationType
                ├── notification_event_type_repository.go  # Repositorio NotificationEventType
                ├── repository.go                          # Repositorio BusinessNotificationConfig
                └── mappers/
                    ├── notification_type_to_domain.go
                    ├── notification_type_to_model.go
                    ├── notification_event_type_to_domain.go
                    ├── notification_event_type_to_model.go
                    ├── to_domain.go
                    └── to_model.go
```

## Jerarquía de Datos

### 1. Notification Types (Nivel Superior)

**Tabla:** `notification_types`

Tipos de canales de notificación disponibles:

```go
type NotificationType struct {
    ID           uint
    Name         string  // "WhatsApp", "SSE", "Email", "SMS"
    Code         string  // "whatsapp", "sse", "email", "sms" (unique)
    Description  string
    Icon         string
    IsActive     bool
    ConfigSchema map[string]interface{}  // Esquema de configuración específico
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

**Datos iniciales:**
- SSE (Server-Sent Events)
- WhatsApp Business
- Email
- SMS

### 2. Notification Event Types (Nivel Medio)

**Tabla:** `notification_event_types`

Eventos específicos por tipo de notificación:

```go
type NotificationEventType struct {
    ID                 uint
    NotificationTypeID uint    // FK a notification_types
    EventCode          string  // "order.created", "order.shipped", etc.
    EventName          string  // "Pedido Creado", "Pedido Enviado"
    Description        string
    TemplateConfig     map[string]interface{}  // Configuración de template
    IsActive           bool
    CreatedAt          time.Time
    UpdatedAt          time.Time
    NotificationType   *NotificationType  // Relación
}
```

**Índice único:** `(notification_type_id, event_code)`

**Ejemplos:**

**WhatsApp:**
- `order.created` → Confirmación de Pedido
- `order.shipped` → Pedido Enviado
- `order.delivered` → Pedido Entregado
- `order.canceled` → Pedido Cancelado
- `invoice.created` → Factura Generada

**SSE:**
- `order.created` → Nueva Orden
- `order.status_changed` → Cambio de Estado

### 3. Business Notification Configs (Nivel Inferior)

**Tabla:** `business_notification_configs`

Configuraciones específicas por negocio/integración:

```go
type BusinessNotificationConfig struct {
    ID                      uint
    BusinessID              uint  // FK a businesses
    IntegrationID           uint  // FK a integrations (origen del evento)
    NotificationTypeID      uint  // FK a notification_types (canal de salida)
    NotificationEventTypeID uint  // FK a notification_event_types (tipo de evento)
    Enabled                 bool
    Filters                 map[string]interface{}  // Filtros adicionales
    Description             string
    CreatedAt               time.Time
    UpdatedAt               time.Time
    DeletedAt               *time.Time

    // Relaciones
    Integration           *Integration
    NotificationType      *NotificationType
    NotificationEventType *NotificationEventType
    OrderStatusIDs        []uint  // M2M con order_statuses
}
```

**Índice único:** `(integration_id, notification_type_id, notification_event_type_id)`

**Relación M2M:** `business_notification_config_order_statuses`
- Permite configurar en qué estados de orden disparar la notificación
- Estados disponibles: pending, processing, shipped, delivered, completed, cancelled, refunded, failed, on_hold

## Flujo de Uso

### Ejemplo Completo

**Configuración:**
```
Business: "Mi Tienda" (ID: 1)
Integration: "Shopify - Mi Tiendita" (ID: 5, type: shopify)
NotificationType: "WhatsApp" (ID: 2, code: "whatsapp")
NotificationEventType: "Confirmación de Pedido" (ID: 10, event_code: "order.created")
OrderStatuses: [created (ID: 1), paid (ID: 3)]
```

**Resultado:**
- Cuando una orden de la integración Shopify (ID: 5) genera el evento `order.created`
- Y el estado de la orden es `created` O `paid`
- → Se envía una notificación por WhatsApp

## API Endpoints

### Notification Types

```http
GET    /api/notification-types           # Listar todos los tipos
GET    /api/notification-types/:id       # Obtener por ID
POST   /api/notification-types           # Crear nuevo tipo
PATCH  /api/notification-types/:id       # Actualizar tipo
DELETE /api/notification-types/:id       # Eliminar tipo (soft delete)
```

### Notification Event Types

```http
GET    /api/notification-event-types?notification_type_id=2  # Listar eventos (filtrable por tipo)
GET    /api/notification-event-types/:id                     # Obtener por ID
POST   /api/notification-event-types                         # Crear nuevo evento
PATCH  /api/notification-event-types/:id                     # Actualizar evento
DELETE /api/notification-event-types/:id                     # Eliminar evento (soft delete)
```

### Business Notification Configs

```http
GET    /api/notification-configs?business_id=1&integration_id=5  # Listar configs (filtrable)
GET    /api/notification-configs/:id                             # Obtener por ID
POST   /api/notification-configs                                 # Crear nueva config
PATCH  /api/notification-configs/:id                             # Actualizar config
DELETE /api/notification-configs/:id                             # Eliminar config (soft delete)
```

## Modelos GORM

Los modelos GORM con tags están centralizados en:

**`/back/migration/shared/models/`**
- `notification_type.go` - Modelo con tags GORM
- `notification_event_type.go` - Modelo con tags GORM
- `notification_config.go` - Modelo con tags GORM (refactorizado)

**Migración:**
- Script SQL: `/back/migration/shared/sql/migrate_notification_system_refactor.sql`
- Incluye creación de tablas, datos iniciales y migración de configs existentes

## Reglas de Arquitectura Hexagonal

### ✅ Domain (Entidades Puras)

```go
// ✅ CORRECTO - Sin tags
type NotificationType struct {
    ID          uint
    Name        string
    Code        string
    IsActive    bool
}
```

### ❌ Domain (NO hacer esto)

```go
// ❌ INCORRECTO - Con tags (esto va en models de migration)
type NotificationType struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"size:100;not null"`
    IsActive bool   `gorm:"default:true"`
}
```

### ✅ Repository (Usa modelos de migration)

```go
import "github.com/secamc93/probability/back/migration/shared/models"

var model models.NotificationType
db.Find(&model)
```

## Migraciones

### 1. Ejecutar AutoMigrate

```bash
cd /back/central
go run cmd/main.go migrate
```

### 2. Ejecutar Script SQL

```bash
psql -U postgres -d probability_db -f /back/migration/shared/sql/migrate_notification_system_refactor.sql
```

### 3. Verificar Datos

```sql
-- Ver tipos de notificación
SELECT * FROM notification_types;

-- Ver eventos de notificación
SELECT
    net.id,
    nt.name as tipo,
    net.event_name,
    net.event_code,
    net.is_active
FROM notification_event_types net
JOIN notification_types nt ON net.notification_type_id = nt.id
ORDER BY nt.id, net.id;

-- Ver configuraciones de negocio
SELECT
    bnc.id,
    bnc.business_id,
    i.name as integration,
    nt.name as tipo_notificacion,
    net.event_name,
    bnc.enabled
FROM business_notification_configs bnc
JOIN integrations i ON bnc.integration_id = i.id
JOIN notification_types nt ON bnc.notification_type_id = nt.id
JOIN notification_event_types net ON bnc.notification_event_type_id = net.id;

-- Ver estados de orden asociados a una config
SELECT
    bnc.id,
    os.name as estado,
    os.code
FROM business_notification_configs bnc
JOIN business_notification_config_order_statuses bcs ON bnc.id = bcs.business_notification_config_id
JOIN order_statuses os ON bcs.order_status_id = os.id
WHERE bnc.id = 1;
```

## Testing

```bash
# Compilar
go build ./...

# Tests
go test ./...

# Test específico
go test ./internal/app/...
go test ./internal/infra/secondary/repository/...
```

## Convenciones

1. **Entidades de dominio:** Sin tags, solo tipos nativos de Go
2. **Modelos GORM:** Centralizados en `/back/migration/shared/models/`
3. **Repositorios:** Usan modelos de migration, retornan entidades de dominio
4. **Handlers:** Cada handler en su propio archivo
5. **Rutas:** Registradas en `routes.go` dentro de cada grupo de handlers
6. **Mappers:** Obligatorios en `request/`, `response/`, `mappers/` para cada handler

## Campos Deprecados (Migración)

Durante la migración, se mantienen campos deprecados para compatibilidad temporal:

```go
EventTypeDeprecated string  // Antiguo event_type (antes de refactorización)
// Se eliminará en versión futura
```

## Dependencias

- **GORM:** ORM para PostgreSQL
- **Gin:** Framework HTTP
- **datatypes.JSON:** Soporte para campos JSONB
- **Zerolog:** Logging estructurado

## Notas Importantes

1. **Unique constraints:** Evitan duplicados en combinaciones clave
2. **Soft deletes:** Implementados con `gorm.DeletedAt`
3. **Preload:** Usar `.Preload()` para cargar relaciones
4. **Validaciones:** Implementadas en capa de aplicación (use cases)
5. **Errores de dominio:** Tipados y centralizados en `domain/errors/`

## Changelog

### v2.0.0 - Refactorización Arquitectura Jerárquica (2026-01-31)

**BREAKING CHANGES:**
- Nueva estructura de tres niveles (NotificationType → NotificationEventType → BusinessNotificationConfig)
- Campo `channels` eliminado, reemplazado por `notification_type_id`
- Campo `event_type` deprecado, reemplazado por `notification_event_type_id`
- Agregado FK `integration_id` (integración que genera el evento)

**Nuevas Features:**
- CRUD completo de NotificationTypes
- CRUD completo de NotificationEventTypes
- Relación M2M con OrderStatuses para filtrar estados que disparan notificaciones
- Script de migración de datos existentes

**Arquitectura:**
- Handlers organizados en carpetas (`notification_type/`, `notification_event_type/`, `notification_config/`)
- Modelos GORM centralizados en `/back/migration/shared/models/`
- Mappers actualizados para usar modelos de migration

---

## Testing

### ✅ Estado de Tests

**Estado**: ✅ Todos los tests pasando
**Fecha**: 2026-01-31
**Arquitectura**: 100% Hexagonal (validado)

### 📊 Cobertura de Tests

#### Resumen Global

```
Capa de Aplicación (app/):              29.8% (5 casos de uso principales)
Capa de Handlers (notification_config): 88.4%
Total de tests:                         40 tests (20 app + 20 handlers)
Total pasando:                          ✅ 40/40 (100%)
```

#### Casos de Uso Testeados

| Caso de Uso | Cobertura | Tests | Estado |
|-------------|-----------|-------|--------|
| Create      | 100%      | 5     | ✅     |
| Update      | 100%      | 5     | ✅     |
| GetByID     | 100%      | 3     | ✅     |
| List        | 100%      | 4     | ✅     |
| Delete      | 100%      | 3     | ✅     |

#### Handlers Testeados

| Handler  | Cobertura | Tests | Estado |
|----------|-----------|-------|--------|
| Create   | 100%      | 4     | ✅     |
| Update   | 100%      | 5     | ✅     |
| GetByID  | 100%      | 4     | ✅     |
| List     | 100%      | 4     | ✅     |
| Delete   | 100%      | 4     | ✅     |

### 🗂️ Estructura de Tests

```
internal/
├── mocks/                                    # Todos los mocks centralizados
│   ├── repository_mock.go
│   ├── notification_type_repository_mock.go
│   ├── notification_event_type_repository_mock.go
│   ├── usecase_mock.go
│   └── logger_mock.go
│
├── app/                                      # Tests de Casos de Uso
│   ├── create_test.go                        # 5 tests
│   ├── update_test.go                        # 5 tests
│   ├── get_test.go                           # 3 tests
│   ├── list_test.go                          # 4 tests
│   └── delete_test.go                        # 3 tests
│
└── infra/primary/handlers/notification_config/  # Tests de Handlers
    ├── create_handler_test.go                # 4 tests
    ├── update_handler_test.go                # 5 tests
    ├── get_by_id_handler_test.go             # 4 tests
    ├── list_handler_test.go                  # 4 tests
    └── delete_handler_test.go                # 4 tests
```

### 🚀 Comandos de Testing

```bash
# Ejecutar todos los tests
go test ./internal/... -v

# Ejecutar solo tests de aplicación
go test ./internal/app -v

# Ejecutar solo tests de handlers
go test ./internal/infra/primary/handlers/notification_config -v

# Ver cobertura
go test ./internal/... -cover

# Generar reporte de cobertura HTML
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 🎯 Mejores Prácticas Aplicadas

#### Arquitectura Hexagonal
- ✅ Todos los mocks en `internal/mocks/` (no dentro de tests)
- ✅ Tests unitarios puros (sin base de datos real)
- ✅ Mocks de interfaces (ports), no de implementaciones
- ✅ Inversión de dependencias respetada

#### Testing Best Practices
- ✅ Tests independientes (sin estado compartido)
- ✅ Nombres descriptivos (documentan el comportamiento)
- ✅ Cobertura de casos felices, errores y casos límite
- ✅ Sin dependencias externas (DB, HTTP, filesystem)
- ✅ Tests rápidos (<50ms total)
- ✅ Patrón AAA (Arrange, Act, Assert)

#### Go Testing Conventions
- ✅ Package testing estándar (sin frameworks pesados)
- ✅ Funciones `Test*` siguiendo convención Go
- ✅ gin.TestMode para handlers HTTP

### 📋 Escenarios de Test por Caso de Uso

#### Create
- ✅ Creación exitosa
- ✅ Detecta duplicados (ErrDuplicateConfig)
- ✅ Error en validación de duplicados
- ✅ Error en persistencia
- ✅ Permite configs con condiciones diferentes

#### Update
- ✅ Actualización completa exitosa
- ✅ Actualización parcial
- ✅ Configuración no existe
- ✅ Error en persistencia
- ✅ Error al recuperar config actualizada

#### GetByID
- ✅ Obtención exitosa
- ✅ Config no encontrada
- ✅ Error de conexión a BD

#### List
- ✅ Listar todas las configs
- ✅ Listar con filtros
- ✅ Resultado vacío
- ✅ Error de conexión

#### Delete
- ✅ Eliminación exitosa
- ✅ Config no existe
- ✅ Error en eliminación

### 📋 Escenarios de Test por Handler

#### CreateHandler
- ✅ HTTP 201 Created
- ✅ HTTP 400 Bad Request (validación)
- ✅ HTTP 409 Conflict (duplicado)
- ✅ HTTP 500 Internal Server Error

#### UpdateHandler
- ✅ HTTP 200 OK
- ✅ HTTP 400 Bad Request (ID inválido)
- ✅ HTTP 400 Bad Request (body inválido)
- ✅ HTTP 404 Not Found
- ✅ HTTP 500 Internal Server Error

#### GetByIDHandler
- ✅ HTTP 200 OK
- ✅ HTTP 400 Bad Request
- ✅ HTTP 404 Not Found
- ✅ HTTP 500 Internal Server Error

#### ListHandler
- ✅ HTTP 200 OK (lista completa)
- ✅ HTTP 200 OK (con filtros)
- ✅ HTTP 200 OK (array vacío)
- ✅ HTTP 500 Internal Server Error

#### DeleteHandler
- ✅ HTTP 204 No Content
- ✅ HTTP 400 Bad Request
- ✅ HTTP 404 Not Found
- ✅ HTTP 500 Internal Server Error

---

**Última actualización:** 2026-01-31
