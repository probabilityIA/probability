# Notification Config Module

Sistema de configuración de notificaciones multi-canal para Probability. Permite configurar qué notificaciones enviar, por qué canal, y bajo qué condiciones, para cada integración de e-commerce.

---

## ¿Qué hace este módulo?

Este módulo permite a los negocios **configurar notificaciones automáticas** que se disparan cuando ocurren eventos específicos en sus órdenes (creación, cambio de estado, envío, cancelación, etc.).

### Problema que resuelve

En una plataforma multi-tenant como Probability, cada negocio:
- Tiene múltiples integraciones (Shopify, Amazon, MercadoLibre)
- Necesita notificar a sus clientes por diferentes canales (WhatsApp, Email, SMS)
- Quiere diferentes mensajes para diferentes eventos (pedido creado, enviado, entregado)
- Necesita filtrar cuándo enviar cada notificación (solo para ciertos estados, métodos de pago, etc.)

**Este módulo centraliza y hace configurable todo este sistema de notificaciones.**

---

## ¿Cómo funciona?

### Flujo Conceptual

```
┌─────────────────────────────────────────────────────────────────┐
│  1. EVENTO OCURRE                                               │
│  Una orden es creada en Shopify                                │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│  2. SISTEMA BUSCA CONFIGURACIONES                               │
│  ¿Hay configs activas para esta integración + evento?          │
│  → Business: "Mi Tienda"                                        │
│  → Integration: "Shopify Mi Tiendita"                          │
│  → Event: "order.created"                                       │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│  3. VALIDA CONDICIONES                                          │
│  ¿Cumple con los filtros configurados?                         │
│  → Estado de la orden: "created" ✓                             │
│  → Método de pago: "contra_entrega" ✓                          │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│  4. ENVÍA NOTIFICACIÓN                                          │
│  Por el canal configurado:                                      │
│  → WhatsApp: "Tu pedido #1234 ha sido confirmado"              │
│  → Email: "Confirmación de Pedido"                             │
└─────────────────────────────────────────────────────────────────┘
```

---

## Arquitectura de 3 Niveles

El módulo sigue una **jerarquía de tres niveles** que permite flexibilidad y reutilización:

### Nivel 1: Tipos de Notificación (Canales)

**¿Qué es?** Define los canales de comunicación disponibles.

**Tabla:** `notification_types`

**Ejemplos:**
- WhatsApp Business
- Email
- SMS
- SSE (Server-Sent Events - notificaciones en tiempo real en la web)

**Características:**
- Cada tipo tiene un código único (`whatsapp`, `email`, `sms`, `sse`)
- Puede estar activo o inactivo globalmente
- Define un esquema de configuración específico (ej: para WhatsApp se necesita API key, número, etc.)

```go
type NotificationType struct {
    ID           uint
    Name         string  // "WhatsApp Business"
    Code         string  // "whatsapp" (unique)
    Description  string
    Icon         string
    IsActive     bool
    ConfigSchema map[string]interface{}  // Esquema JSON de configuración
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

---

### Nivel 2: Tipos de Evento (Qué pasó)

**¿Qué es?** Define los eventos específicos que pueden ocurrir en cada canal.

**Tabla:** `notification_event_types`

**Relación:** Cada evento pertenece a UN tipo de notificación.

**Ejemplos para WhatsApp:**
- `order.created` → "Confirmación de Pedido"
- `order.shipped` → "Tu pedido ha sido enviado"
- `order.delivered` → "Tu pedido ha sido entregado"
- `order.canceled` → "Pedido cancelado"
- `invoice.created` → "Factura disponible"

**Ejemplos para SSE (notificaciones web):**
- `order.created` → "Nueva Orden en el Dashboard"
- `order.status_changed` → "Estado de Orden Actualizado"

```go
type NotificationEventType struct {
    ID                   uint
    NotificationTypeID   uint    // FK a notification_types
    EventCode            string  // "order.created", "order.shipped"
    EventName            string  // "Confirmación de Pedido"
    Description          string
    TemplateConfig       map[string]interface{}  // Config del template
    IsActive             bool
    AllowedOrderStatusIDs []uint // Estados permitidos (vacío = todos)
    CreatedAt            time.Time
    UpdatedAt            time.Time
}
```

**Índice único:** `(notification_type_id, event_code)` - No puede haber dos eventos con el mismo código para el mismo tipo.

**Relación M2M con Order Statuses (AllowedOrderStatuses):**
- Tabla pivote: `notification_event_type_allowed_statuses`
- Define qué estados de orden son válidos para cada tipo de evento
- Si está vacío → se permiten todos los estados
- Ejemplo: `order.created` → solo `[pending, processing]`, `order.shipped` → solo `[shipped, delivered]`
- Se usa en el frontend para filtrar los toggles de estados en el formulario de reglas

---

### Nivel 3: Configuraciones de Negocio (Cuándo y cómo enviar)

**¿Qué es?** Configura qué notificaciones enviar para cada integración de un negocio.

**Tabla:** `business_notification_configs`

**Relación:** Conecta una integración con un tipo de notificación y un evento.

```go
type BusinessNotificationConfig struct {
    ID                      uint
    BusinessID              uint  // FK a businesses (el negocio dueño)
    IntegrationID           uint  // FK a integrations (de dónde viene el evento)
    NotificationTypeID      uint  // FK a notification_types (por dónde enviar)
    NotificationEventTypeID uint  // FK a notification_event_types (qué evento)
    Enabled                 bool  // ¿Está activa esta config?
    Filters                 map[string]interface{}  // Filtros adicionales (JSON)
    Description             string
    CreatedAt               time.Time
    UpdatedAt               time.Time
    DeletedAt               *time.Time  // Soft delete

    // Relaciones Many-to-Many
    OrderStatusIDs []uint  // Estados de orden que disparan la notificación
}
```

**Índice único:** `(integration_id, notification_type_id, notification_event_type_id)` - Una integración no puede tener dos configs iguales.

**Relación M2M con Order Statuses:**
- Tabla intermedia: `business_notification_config_order_statuses`
- Permite filtrar: "Solo enviar WhatsApp cuando el estado sea 'created' o 'paid'"
- Estados disponibles: `pending`, `processing`, `shipped`, `delivered`, `completed`, `cancelled`, `refunded`, `failed`, `on_hold`

---

## Flujo Integration-Centric (Batch Sync)

A partir de v3.0.0, la gestión de configs se centra en la **integración**: en lugar de crear/editar configs individualmente, se gestionan N reglas por integración y se sincronizan de una vez.

### Flujo UI

```
┌──────────────────────────────┐
│  1. LISTA AGRUPADA           │
│  Muestra integraciones con   │
│  sus reglas (count, canales) │
│  [Configurar] [+ Agregar]   │
└──────────────────────────────┘
         ↓ Click "Agregar"
┌──────────────────────────────┐
│  2. INTEGRATION PICKER       │
│  Seleccionar integración     │
│  ecommerce del negocio       │
└──────────────────────────────┘
         ↓ Selecciona una
┌──────────────────────────────┐
│  3. INTEGRATION RULES FORM   │
│  Gestionar N reglas:         │
│  ┌──────────────────────┐    │
│  │ Regla 1: WhatsApp +  │    │
│  │ order.created +      │    │
│  │ [pending, processing]│    │
│  └──────────────────────┘    │
│  ┌──────────────────────┐    │
│  │ Regla 2: SSE +       │    │
│  │ order.status_changed │    │
│  └──────────────────────┘    │
│  [+ Agregar regla]           │
│  [Guardar]                   │
└──────────────────────────────┘
         ↓ Guardar
┌──────────────────────────────┐
│  4. BATCH SYNC               │
│  PUT /notification-configs/  │
│  sync?business_id=X          │
│  → Crea nuevas               │
│  → Actualiza existentes      │
│  → Elimina removidas         │
│  → Todo en UNA transacción   │
└──────────────────────────────┘
```

### Componentes Frontend

| Componente | Descripción |
|-----------|-------------|
| `ConfigListTable` | Lista agrupada por integración (logo, nombre, count, canales, botón Configurar) |
| `IntegrationPicker` | Modal para seleccionar integración ecommerce |
| `IntegrationRulesForm` | Form principal: carga existentes, gestiona N reglas, sync batch |
| `RuleCard` | Card compacta para una regla (canal, evento, estados filtrados, toggle) |
| `NotificationEventTypeForm` | CRUD de event types con selector de allowed statuses |

---

## Ejemplos de Uso

### Ejemplo 1: Batch Sync - Configurar reglas para Shopify

**Escenario:**
"Mi Tienda" quiere configurar 3 reglas de notificación para su integración Shopify de una vez.

**Request:**

```bash
curl -X PUT http://localhost:8080/api/v1/notification-configs/sync?business_id=1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "integration_id": 5,
    "rules": [
      {
        "notification_type_id": 2,
        "notification_event_type_id": 3,
        "enabled": true,
        "description": "WhatsApp al crear orden",
        "order_status_ids": [1, 2]
      },
      {
        "notification_type_id": 2,
        "notification_event_type_id": 4,
        "enabled": true,
        "description": "WhatsApp al enviar",
        "order_status_ids": [3, 4]
      },
      {
        "id": 7,
        "notification_type_id": 1,
        "notification_event_type_id": 2,
        "enabled": true,
        "description": "SSE cambio de estado (existente, actualizar)",
        "order_status_ids": []
      }
    ]
  }'
```

**Response:**

```json
{
  "created": 2,
  "updated": 1,
  "deleted": 0,
  "configs": [
    {
      "id": 15,
      "business_id": 1,
      "integration_id": 5,
      "notification_type_id": 2,
      "notification_event_type_id": 3,
      "enabled": true,
      "description": "WhatsApp al crear orden",
      "order_status_ids": [1, 2]
    },
    {
      "id": 16,
      "business_id": 1,
      "integration_id": 5,
      "notification_type_id": 2,
      "notification_event_type_id": 4,
      "enabled": true,
      "description": "WhatsApp al enviar",
      "order_status_ids": [3, 4]
    },
    {
      "id": 7,
      "business_id": 1,
      "integration_id": 5,
      "notification_type_id": 1,
      "notification_event_type_id": 2,
      "enabled": true,
      "description": "SSE cambio de estado (existente, actualizar)",
      "order_status_ids": []
    }
  ]
}
```

**Lógica del sync:**
- Rules sin `id` → se crean como nuevas
- Rules con `id` → se actualizan
- Configs existentes cuyo `id` no aparece en el request → se eliminan (soft delete)
- Validación: no permite duplicados `(notification_type_id, notification_event_type_id)` en el mismo request

---

### Ejemplo 2: Confirmación de Pedido por WhatsApp (Config individual)

**Escenario:**
"Mi Tienda" quiere enviar un mensaje de WhatsApp cuando se crea una orden en su tienda Shopify, solo si el pago es contra entrega o PSE.

**Configuración:**

```json
{
  "business_id": 1,
  "integration_id": 5,
  "notification_type_id": 2,
  "notification_event_type_id": 10,
  "enabled": true,
  "order_status_ids": [1, 3],
  "filters": {
    "payment_methods": ["contra_entrega", "pse"]
  },
  "description": "Confirmación por WhatsApp para órdenes de Shopify"
}
```

**Datos relacionados:**
- Business: "Mi Tienda" (ID: 1)
- Integration: "Shopify - Mi Tiendita" (ID: 5)
- NotificationType: "WhatsApp" (ID: 2, code: `whatsapp`)
- NotificationEventType: "Confirmación de Pedido" (ID: 10, event_code: `order.created`)
- OrderStatuses:
  - ID 1: `created`
  - ID 3: `paid`

---

### Ejemplo 3: Notificaciones en Dashboard (SSE)

**Escenario:**
"Mi Tienda" quiere mostrar notificaciones en tiempo real en el dashboard cuando cambia el estado de una orden.

**Configuración:**

```json
{
  "business_id": 1,
  "integration_id": 5,
  "notification_type_id": 1,
  "notification_event_type_id": 2,
  "enabled": true,
  "order_status_ids": [2, 4, 5],
  "description": "Notificaciones en tiempo real en el dashboard"
}
```

---

## API Endpoints

### Tipos de Notificación (Notification Types)

```http
GET    /api/v1/notification-types           # Listar todos los canales
GET    /api/v1/notification-types/:id       # Obtener un tipo específico
POST   /api/v1/notification-types           # Crear nuevo canal (admin)
PATCH  /api/v1/notification-types/:id       # Actualizar canal
DELETE /api/v1/notification-types/:id       # Eliminar canal (soft delete)
```

---

### Tipos de Evento (Notification Event Types)

```http
GET    /api/v1/notification-event-types?notification_type_id=2  # Listar eventos (filtrable)
GET    /api/v1/notification-event-types/:id                     # Obtener evento
POST   /api/v1/notification-event-types                         # Crear evento
PATCH  /api/v1/notification-event-types/:id                     # Actualizar evento
DELETE /api/v1/notification-event-types/:id                     # Eliminar evento
```

**Respuesta incluye `allowed_order_status_ids`:**

```json
{
  "id": 3,
  "notification_type_id": 2,
  "event_code": "order.created",
  "event_name": "Confirmación de Pedido",
  "description": "Se envía cuando se crea una nueva orden",
  "is_active": true,
  "allowed_order_status_ids": [1, 2]
}
```

**Crear/Actualizar con allowed statuses:**

```bash
POST /api/v1/notification-event-types
{
  "notification_type_id": 2,
  "event_code": "order.ready_for_pickup",
  "event_name": "Pedido Listo para Recoger",
  "is_active": true,
  "allowed_order_status_ids": [4, 5]
}

PATCH /api/v1/notification-event-types/3
{
  "event_name": "Confirmación de Pedido Actualizada",
  "allowed_order_status_ids": [1, 2, 3]
}
```

**Nota sobre update:** Enviar `allowed_order_status_ids: []` limpia todos los estados permitidos. No enviar el campo mantiene los estados existentes sin cambio.

---

### Configuraciones de Negocio (Business Notification Configs)

```http
GET    /api/v1/notification-configs?business_id=1&integration_id=5  # Listar configs
GET    /api/v1/notification-configs/:id                             # Obtener config
POST   /api/v1/notification-configs                                 # Crear config
PATCH  /api/v1/notification-configs/:id                             # Actualizar config
DELETE /api/v1/notification-configs/:id                             # Eliminar config
PUT    /api/v1/notification-configs/sync?business_id=X              # Batch sync (NUEVO v3.0)
```

---

### Batch Sync (NUEVO v3.0)

```http
PUT /api/v1/notification-configs/sync?business_id=X
```

**Propósito:** Sincronizar todas las reglas de notificación de una integración en una sola operación transaccional.

**Request Body:**

```json
{
  "integration_id": 5,
  "rules": [
    {
      "notification_type_id": 2,
      "notification_event_type_id": 3,
      "enabled": true,
      "description": "Nueva regla",
      "order_status_ids": [1, 2]
    },
    {
      "id": 7,
      "notification_type_id": 1,
      "notification_event_type_id": 2,
      "enabled": true,
      "description": "Regla existente actualizada",
      "order_status_ids": []
    }
  ]
}
```

**Lógica:**

| Campo `id` en rule | Acción |
|--------------------|--------|
| No presente / nil | Crea nueva config |
| Presente (ej: `7`) | Actualiza config existente |
| Config existente NO en rules | Se elimina (soft delete) |

**Validaciones:**
- `integration_id` es obligatorio
- No permite duplicados `(notification_type_id, notification_event_type_id)` en el mismo request
- Cada rule debe tener `notification_type_id` y `notification_event_type_id`

**Response:**

```json
{
  "created": 1,
  "updated": 1,
  "deleted": 2,
  "configs": [ ... ]
}
```

**Efectos secundarios:**
- Invalida cache de configs por integración
- Opera dentro de una transacción DB (todo o nada)

---

## Arquitectura Técnica

### Estructura de Carpetas (Arquitectura Hexagonal)

```
notification_config/
├── bundle.go                    # Ensamblador del módulo
└── internal/
    ├── domain/                  # DOMINIO (núcleo puro)
    │   ├── entities/            # Entidades sin tags
    │   │   ├── notification_type.go
    │   │   ├── notification_event_type.go
    │   │   └── business_notification_config.go
    │   ├── dtos/                # DTOs de dominio
    │   │   ├── filter.go
    │   │   └── sync.go          # DTOs para batch sync
    │   ├── ports/               # Interfaces
    │   │   ├── repository.go
    │   │   └── usecase.go
    │   └── errors/              # Errores de dominio
    │
    ├── app/                     # APLICACIÓN (casos de uso)
    │   ├── constructor.go
    │   ├── create*.go
    │   ├── update*.go
    │   ├── delete*.go
    │   ├── get*.go
    │   ├── list*.go
    │   ├── sync.go              # Caso de uso batch sync
    │   ├── request/
    │   ├── response/
    │   └── mappers/
    │
    ├── infra/                   # INFRAESTRUCTURA
    │   ├── primary/             # Adaptadores de entrada
    │   │   └── handlers/
    │   │       ├── notification_type/
    │   │       ├── notification_event_type/
    │   │       └── notification_config/
    │   │           ├── constructor.go
    │   │           ├── routes.go
    │   │           ├── create_handler.go
    │   │           ├── list_handler.go
    │   │           ├── sync_handler.go       # Handler batch sync
    │   │           └── request/
    │   │               └── sync_request.go   # Request DTO HTTP
    │   │
    │   └── secondary/           # Adaptadores de salida
    │       ├── repository/
    │       │   ├── repository.go
    │       │   ├── notification_type_repository.go
    │       │   ├── notification_event_type_repository.go
    │       │   ├── sync_configs.go            # Repo transaccional sync
    │       │   ├── order_status_queries.go    # Queries replicadas (aislamiento)
    │       │   └── mappers/
    │       └── cache/
    │           ├── constructor.go
    │           ├── warmup_cache.go
    │           ├── invalidate_configs_by_integration.go
    │           └── ...
    │
    └── mocks/                   # Mocks para testing
```

---

### Tablas de Base de Datos

| Tabla | Descripción |
|-------|-------------|
| `notification_types` | Canales (WhatsApp, Email, SMS, SSE) |
| `notification_event_types` | Eventos por canal (order.created, order.shipped) |
| `notification_event_type_allowed_statuses` | M2M: estados permitidos por tipo de evento |
| `business_notification_configs` | Configs de negocio (integración + canal + evento) |
| `business_notification_config_order_statuses` | M2M: estados que disparan la notificación |

### Modelos GORM (fuente de verdad)

Centralizados en `/back/migration/shared/models/`:
- `notification_type.go`
- `notification_event_type.go` (incluye `AllowedOrderStatuses` M2M)
- `notification_config.go`

---

### Reglas de Arquitectura Hexagonal

#### Domain (Entidades Puras)

```go
// CORRECTO - Sin tags, solo tipos nativos
type NotificationType struct {
    ID          uint
    Name        string
    Code        string
    IsActive    bool
}

type NotificationEventType struct {
    ID                    uint
    NotificationTypeID    uint
    EventCode             string
    EventName             string
    IsActive              bool
    AllowedOrderStatusIDs []uint  // Estados de orden permitidos
}
```

#### Repository (Usa modelos GORM externos)

```go
import "github.com/secamc93/probability/back/migration/shared/models"

var model models.NotificationType
db.Find(&model)

// Preload de relaciones M2M
db.Preload("AllowedOrderStatuses").Find(&eventTypes)
```

---

## Guía de Configuración

### Paso 1: Configurar Tipos de Notificación

Los tipos básicos (WhatsApp, Email, SMS, SSE) vienen preconfigurados. Solo necesitas activarlos/desactivarlos según tu plan.

```bash
GET /api/v1/notification-types

PATCH /api/v1/notification-types/4
{
  "is_active": false
}
```

### Paso 2: Configurar Eventos con Estados Permitidos

Los eventos comunes vienen precargados con sus estados permitidos. Puedes crear eventos personalizados y definir qué estados son válidos.

```bash
POST /api/v1/notification-event-types
{
  "notification_type_id": 2,
  "event_code": "order.ready_for_pickup",
  "event_name": "Pedido Listo para Recoger",
  "is_active": true,
  "allowed_order_status_ids": [4, 5]
}
```

**Estados permitidos por defecto (seed):**

| Evento | Estados Permitidos |
|--------|-------------------|
| SSE order.created | pending, processing |
| SSE order.status_changed | todos |
| WA order.created | pending, processing |
| WA order.shipped | shipped, delivered |
| WA order.delivered | delivered, completed |
| WA order.canceled | cancelled, refunded |
| WA invoice.created | todos |

### Paso 3: Configurar Reglas por Integración (Batch Sync)

Usa el endpoint de sync para configurar todas las reglas de una integración de una vez.

```bash
PUT /api/v1/notification-configs/sync?business_id=1
{
  "integration_id": 5,
  "rules": [
    {
      "notification_type_id": 2,
      "notification_event_type_id": 3,
      "enabled": true,
      "order_status_ids": [1, 2],
      "description": "WhatsApp confirmación"
    },
    {
      "notification_type_id": 1,
      "notification_event_type_id": 2,
      "enabled": true,
      "order_status_ids": [],
      "description": "SSE cambio de estado"
    }
  ]
}
```

---

## Verificar Datos en BD

```sql
-- Ver tipos de notificación
SELECT * FROM notification_types;

-- Ver eventos con estados permitidos
SELECT
    net.id,
    nt.name as tipo,
    net.event_name,
    net.event_code,
    net.is_active,
    ARRAY_AGG(os.name) as allowed_statuses
FROM notification_event_types net
JOIN notification_types nt ON net.notification_type_id = nt.id
LEFT JOIN notification_event_type_allowed_statuses netas
    ON net.id = netas.notification_event_type_id
LEFT JOIN order_statuses os ON netas.order_status_id = os.id
WHERE net.deleted_at IS NULL
GROUP BY net.id, nt.name, net.event_name, net.event_code, net.is_active
ORDER BY nt.id, net.id;

-- Ver configuraciones agrupadas por integración
SELECT
    i.name as integration,
    COUNT(bnc.id) as total_rules,
    COUNT(bnc.id) FILTER (WHERE bnc.enabled) as active_rules,
    ARRAY_AGG(DISTINCT nt.name) as channels
FROM business_notification_configs bnc
JOIN integrations i ON bnc.integration_id = i.id
JOIN notification_types nt ON bnc.notification_type_id = nt.id
WHERE bnc.deleted_at IS NULL
GROUP BY i.name
ORDER BY i.name;

-- Ver estados asociados a una config
SELECT
    bnc.id,
    os.name as estado,
    os.code
FROM business_notification_configs bnc
JOIN business_notification_config_order_statuses bcs
  ON bnc.id = bcs.business_notification_config_id
JOIN order_statuses os ON bcs.order_status_id = os.id
WHERE bnc.id = 1;
```

---

## Desarrollo

### Compilar

```bash
go build ./...
```

### Ejecutar Tests

```bash
go test ./...
```

### Migraciones

```bash
cd /back/migration
go run cmd/main.go
```

Esto ejecuta:
1. `AutoMigrate(&models.NotificationEventType{})` - Crea/actualiza tabla + pivot M2M
2. `seedAllowedOrderStatusesByEventType()` - Inserta estados permitidos por evento

---

## Convenciones

1. **Entidades de dominio:** Sin tags, solo tipos nativos de Go
2. **Modelos GORM:** Centralizados en `/back/migration/shared/models/`
3. **Repositorios:** Usan modelos de migration, retornan entidades de dominio
4. **Handlers:** Cada handler en su propio archivo (`create_handler.go`, `sync_handler.go`)
5. **Rutas:** Registradas en `routes.go` dentro de cada grupo de handlers
6. **Mappers:** Obligatorios en `request/`, `response/`, `mappers/` para cada handler
7. **Aislamiento de repos:** Order status queries replicadas localmente en `order_status_queries.go`
8. **Cache:** Invalidación automática tras sync por integración

---

## Dependencias

- **GORM:** ORM para PostgreSQL
- **Gin:** Framework HTTP
- **datatypes.JSON:** Soporte para campos JSONB
- **Zerolog:** Logging estructurado
- **Redis:** Cache de configs por integración

---

## Notas Importantes

1. **Unique constraints:** Evitan duplicados en combinaciones clave
2. **Soft deletes:** Implementados con `gorm.DeletedAt`
3. **Preload:** Usar `.Preload("AllowedOrderStatuses")` para cargar estados permitidos
4. **Validaciones:** Implementadas en capa de aplicación (use cases)
5. **Errores de dominio:** Tipados y centralizados en `domain/errors/`
6. **Batch sync:** Ruta `/sync` registrada ANTES de `/:id` para evitar que Gin la interprete como ID
7. **Super admin:** `business_id` se resuelve vía `resolveBusinessID()` (query param para super admin, JWT para usuarios normales)

---

## Changelog

### v3.0.0 - Batch Sync + Allowed Statuses (2026-02-28)

**Nuevas Features:**
- **Batch Sync endpoint** (`PUT /notification-configs/sync`): crear, actualizar y eliminar configs de una integración en una sola transacción
- **AllowedOrderStatuses en Event Types**: relación M2M que define qué estados de orden son válidos por tipo de evento
- **Frontend integration-centric**: formulario centrado en integración con multi-reglas
- **ConfigListTable agrupada**: vista agrupada por integración (logo, conteo, canales)
- **IntegrationPicker**: modal de selección de integración ecommerce
- **RuleCard**: card compacta para gestionar una regla (canal, evento, estados filtrados)
- **CRUD Event Types con allowed statuses**: selector de estados permitidos en el formulario de tipos de evento

**Backend:**
- Nuevo use case `SyncByIntegration` con validación de duplicados y transacción
- Nuevo repo `SyncConfigs` con operaciones batch en `db.Transaction`
- Preload de `AllowedOrderStatuses` en todas las queries de event types
- Soporte de `allowed_order_status_ids` en create/update de event types
- Migración + seed de tabla pivote `notification_event_type_allowed_statuses`
- Cache invalidation tras sync

**Frontend:**
- `IntegrationRulesForm` - formulario multi-reglas por integración
- `RuleCard` - card con canal (botones color), evento (dropdown), estados (toggles filtrados)
- `IntegrationPicker` - selector de integraciones ecommerce
- `ConfigListTable` reescrita: agrupada por integración
- `NotificationEventTypeForm` con multi-select de allowed statuses
- `page.tsx` usa flujo IntegrationPicker -> IntegrationRulesForm

### v2.0.0 - Refactorización Arquitectura Jerárquica (2026-01-31)

**BREAKING CHANGES:**
- Nueva estructura de tres niveles (NotificationType -> NotificationEventType -> BusinessNotificationConfig)
- Campo `channels` eliminado, reemplazado por `notification_type_id`
- Campo `event_type` deprecado, reemplazado por `notification_event_type_id`
- Agregado FK `integration_id` (integración que genera el evento)

**Nuevas Features:**
- CRUD completo de NotificationTypes
- CRUD completo de NotificationEventTypes
- Relación M2M con OrderStatuses para filtrar estados
- Script de migración de datos existentes

**Arquitectura:**
- Handlers organizados en carpetas separadas
- Modelos GORM centralizados en `/back/migration/shared/models/`
- Tests completos (40 tests, 100% pasando)
- Arquitectura hexagonal 100% validada

---

**Última actualización:** 2026-02-28
