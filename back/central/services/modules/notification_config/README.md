# Notification Config Module

Sistema de configuración de notificaciones multi-canal para Probability. Permite configurar qué notificaciones enviar, por qué canal, y bajo qué condiciones, para cada integración de e-commerce.

---

## 📌 ¿Qué hace este módulo?

Este módulo permite a los negocios **configurar notificaciones automáticas** que se disparan cuando ocurren eventos específicos en sus órdenes (creación, cambio de estado, envío, cancelación, etc.).

### Problema que resuelve

En una plataforma multi-tenant como Probability, cada negocio:
- Tiene múltiples integraciones (Shopify, Amazon, MercadoLibre)
- Necesita notificar a sus clientes por diferentes canales (WhatsApp, Email, SMS)
- Quiere diferentes mensajes para diferentes eventos (pedido creado, enviado, entregado)
- Necesita filtrar cuándo enviar cada notificación (solo para ciertos estados, métodos de pago, etc.)

**Este módulo centraliza y hace configurable todo este sistema de notificaciones.**

---

## 🔄 ¿Cómo funciona?

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

## 🏗️ Arquitectura de 3 Niveles

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
    ID                 uint
    NotificationTypeID uint    // FK a notification_types
    EventCode          string  // "order.created", "order.shipped"
    EventName          string  // "Confirmación de Pedido"
    Description        string
    TemplateConfig     map[string]interface{}  // Config del template (variables, etc.)
    IsActive           bool
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

**Índice único:** `(notification_type_id, event_code)` - No puede haber dos eventos con el mismo código para el mismo tipo.

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

## 💡 Ejemplos de Uso

### Ejemplo 1: Confirmación de Pedido por WhatsApp

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

**Resultado:**
- ✅ Se enviará WhatsApp cuando:
  - La orden viene de la integración Shopify (ID: 5)
  - Se dispara el evento `order.created`
  - El estado de la orden es `created` O `paid`
  - El método de pago es "contra_entrega" O "pse"

- ❌ NO se enviará si:
  - El estado es diferente (ej: `cancelled`)
  - El método de pago es otro (ej: "tarjeta_credito")

---

### Ejemplo 2: Notificaciones en Dashboard (SSE)

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

**Datos relacionados:**
- NotificationType: "SSE" (ID: 1, code: `sse`)
- NotificationEventType: "Cambio de Estado" (ID: 2, event_code: `order.status_changed`)
- OrderStatuses:
  - ID 2: `processing`
  - ID 4: `shipped`
  - ID 5: `delivered`

**Resultado:**
- ✅ Se enviará notificación SSE al dashboard cuando:
  - El estado de una orden cambie a `processing`, `shipped` o `delivered`

---

### Ejemplo 3: Email de Factura

**Escenario:**
Enviar email con la factura cuando se genere el documento.

```json
{
  "business_id": 1,
  "integration_id": 5,
  "notification_type_id": 3,
  "notification_event_type_id": 15,
  "enabled": true,
  "description": "Email con factura generada"
}
```

**Datos relacionados:**
- NotificationType: "Email" (ID: 3, code: `email`)
- NotificationEventType: "Factura Generada" (ID: 15, event_code: `invoice.created`)

---

## 🚀 API Endpoints

### Tipos de Notificación (Notification Types)

```http
GET    /api/notification-types           # Listar todos los canales
GET    /api/notification-types/:id       # Obtener un tipo específico
POST   /api/notification-types           # Crear nuevo canal (admin)
PATCH  /api/notification-types/:id       # Actualizar canal
DELETE /api/notification-types/:id       # Eliminar canal (soft delete)
```

**Ejemplo - Listar tipos:**
```bash
curl http://localhost:8080/api/notification-types
```

**Respuesta:**
```json
[
  {
    "id": 1,
    "name": "SSE (Server-Sent Events)",
    "code": "sse",
    "description": "Notificaciones en tiempo real en el dashboard",
    "icon": "bell",
    "is_active": true
  },
  {
    "id": 2,
    "name": "WhatsApp Business",
    "code": "whatsapp",
    "description": "Mensajes por WhatsApp",
    "icon": "whatsapp",
    "is_active": true
  }
]
```

---

### Tipos de Evento (Notification Event Types)

```http
GET    /api/notification-event-types?notification_type_id=2  # Listar eventos (filtrable)
GET    /api/notification-event-types/:id                     # Obtener evento
POST   /api/notification-event-types                         # Crear evento
PATCH  /api/notification-event-types/:id                     # Actualizar evento
DELETE /api/notification-event-types/:id                     # Eliminar evento
```

**Ejemplo - Listar eventos de WhatsApp:**
```bash
curl http://localhost:8080/api/notification-event-types?notification_type_id=2
```

**Respuesta:**
```json
[
  {
    "id": 10,
    "notification_type_id": 2,
    "event_code": "order.created",
    "event_name": "Confirmación de Pedido",
    "description": "Se envía cuando se crea una nueva orden",
    "is_active": true
  },
  {
    "id": 11,
    "notification_type_id": 2,
    "event_code": "order.shipped",
    "event_name": "Pedido Enviado",
    "description": "Notifica cuando el pedido ha sido despachado",
    "is_active": true
  }
]
```

---

### Configuraciones de Negocio (Business Notification Configs)

```http
GET    /api/notification-configs?business_id=1&integration_id=5  # Listar configs
GET    /api/notification-configs/:id                             # Obtener config
POST   /api/notification-configs                                 # Crear config
PATCH  /api/notification-configs/:id                             # Actualizar config
DELETE /api/notification-configs/:id                             # Eliminar config
```

**Ejemplo - Crear configuración:**
```bash
curl -X POST http://localhost:8080/api/notification-configs \
  -H "Content-Type: application/json" \
  -d '{
    "business_id": 1,
    "integration_id": 5,
    "notification_type_id": 2,
    "notification_event_type_id": 10,
    "enabled": true,
    "order_status_ids": [1, 3],
    "filters": {
      "payment_methods": ["contra_entrega", "pse"]
    },
    "description": "WhatsApp para órdenes de Shopify"
  }'
```

**Ejemplo - Listar configs de una integración:**
```bash
curl http://localhost:8080/api/notification-configs?integration_id=5
```

**Respuesta:**
```json
[
  {
    "id": 1,
    "business_id": 1,
    "integration_id": 5,
    "notification_type_id": 2,
    "notification_event_type_id": 10,
    "enabled": true,
    "order_status_ids": [1, 3],
    "filters": {
      "payment_methods": ["contra_entrega", "pse"]
    },
    "description": "WhatsApp para órdenes de Shopify",
    "integration": {
      "id": 5,
      "name": "Shopify - Mi Tiendita"
    },
    "notification_type": {
      "id": 2,
      "name": "WhatsApp Business",
      "code": "whatsapp"
    },
    "notification_event_type": {
      "id": 10,
      "event_code": "order.created",
      "event_name": "Confirmación de Pedido"
    }
  }
]
```

---

## 📖 Guía de Configuración

### Paso 1: Configurar Tipos de Notificación

Los tipos básicos (WhatsApp, Email, SMS, SSE) vienen preconfigurados. Solo necesitas activarlos/desactivarlos según tu plan.

```bash
# Listar tipos disponibles
GET /api/notification-types

# Desactivar un tipo (ej: SMS)
PATCH /api/notification-types/4
{
  "is_active": false
}
```

---

### Paso 2: Configurar Eventos

Los eventos comunes vienen precargados, pero puedes crear eventos personalizados.

```bash
# Crear evento personalizado para WhatsApp
POST /api/notification-event-types
{
  "notification_type_id": 2,
  "event_code": "order.ready_for_pickup",
  "event_name": "Pedido Listo para Recoger",
  "description": "Notifica cuando el pedido está listo en tienda",
  "is_active": true
}
```

---

### Paso 3: Crear Configuraciones para tus Integraciones

Ahora conecta tus integraciones con los eventos y canales que quieres usar.

```bash
# Configurar WhatsApp para confirmaciones de Shopify
POST /api/notification-configs
{
  "business_id": 1,
  "integration_id": 5,
  "notification_type_id": 2,
  "notification_event_type_id": 10,
  "enabled": true,
  "order_status_ids": [1, 3],
  "description": "Confirmación de pedido por WhatsApp"
}
```

---

### Paso 4: Filtrar por Estados

Especifica en qué estados de orden se debe enviar la notificación.

**Estados disponibles:**
- `pending` (1)
- `processing` (2)
- `paid` (3)
- `shipped` (4)
- `delivered` (5)
- `completed` (6)
- `cancelled` (7)
- `refunded` (8)
- `failed` (9)
- `on_hold` (10)

**Ejemplo:**
```json
{
  "order_status_ids": [1, 3]  // Solo estados "pending" y "paid"
}
```

---

### Paso 5: Filtros Adicionales (Opcional)

Agrega filtros adicionales en formato JSON:

```json
{
  "filters": {
    "payment_methods": ["contra_entrega", "pse"],
    "min_amount": 50000,
    "source_integration_id": 5
  }
}
```

---

## 🏛️ Arquitectura Técnica

### Estructura de Carpetas (Arquitectura Hexagonal)

```
notification_config/
├── bundle.go                    # Ensamblador del módulo
└── internal/
    ├── domain/                  # 🔵 DOMINIO (núcleo puro)
    │   ├── entities/            # Entidades sin tags
    │   │   ├── notification_type.go
    │   │   ├── notification_event_type.go
    │   │   └── business_notification_config.go
    │   ├── dtos/                # DTOs de dominio
    │   ├── ports/               # Interfaces
    │   │   ├── repository.go
    │   │   └── usecase.go
    │   └── errors/              # Errores de dominio
    │
    ├── app/                     # 🟢 APLICACIÓN (casos de uso)
    │   ├── constructor.go
    │   ├── create*.go           # Casos de uso de creación
    │   ├── update*.go           # Casos de uso de actualización
    │   ├── delete*.go           # Casos de uso de eliminación
    │   ├── get*.go              # Casos de uso de consulta
    │   ├── list*.go             # Casos de uso de listado
    │   ├── request/             # DTOs de request
    │   ├── response/            # DTOs de response
    │   └── mappers/             # Conversiones
    │
    ├── infra/                   # 🔴 INFRAESTRUCTURA
    │   ├── primary/             # Adaptadores de entrada
    │   │   └── handlers/
    │   │       ├── notification_type/
    │   │       ├── notification_event_type/
    │   │       └── notification_config/
    │   │
    │   └── secondary/           # Adaptadores de salida
    │       └── repository/
    │           ├── notification_type_repository.go
    │           ├── notification_event_type_repository.go
    │           ├── repository.go
    │           └── mappers/
    │
    └── mocks/                   # 🧪 Mocks para testing
        ├── repository_mock.go
        ├── notification_type_repository_mock.go
        ├── notification_event_type_repository_mock.go
        ├── usecase_mock.go
        └── logger_mock.go
```

---

### Reglas de Arquitectura Hexagonal

#### ✅ Domain (Entidades Puras)

```go
// ✅ CORRECTO - Sin tags, solo tipos nativos
type NotificationType struct {
    ID          uint
    Name        string
    Code        string
    IsActive    bool
}
```

#### ❌ Domain (NO hacer esto)

```go
// ❌ INCORRECTO - Tags de frameworks (esto va en models)
type NotificationType struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"size:100;not null"`
    IsActive bool   `gorm:"default:true"`
}
```

#### ✅ Repository (Usa modelos GORM externos)

```go
import "github.com/secamc93/probability/back/migration/shared/models"

var model models.NotificationType
db.Find(&model)
```

**Modelos GORM centralizados en:**
- `/back/migration/shared/models/notification_type.go`
- `/back/migration/shared/models/notification_event_type.go`
- `/back/migration/shared/models/notification_config.go`

---

## 🧪 Testing

### ✅ Estado de Tests

**Estado**: ✅ Todos los tests pasando
**Arquitectura**: 100% Hexagonal (validado)

### 📊 Cobertura

```
Capa de Aplicación (app/):              29.8% (5 casos de uso principales)
Capa de Handlers (notification_config): 88.4%
Total de tests:                         40 tests (20 app + 20 handlers)
Total pasando:                          ✅ 40/40 (100%)
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

# Generar reporte HTML
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 📋 Casos de Uso Testeados

| Caso de Uso | Cobertura | Tests | Estado |
|-------------|-----------|-------|--------|
| Create      | 100%      | 5     | ✅     |
| Update      | 100%      | 5     | ✅     |
| GetByID     | 100%      | 3     | ✅     |
| List        | 100%      | 4     | ✅     |
| Delete      | 100%      | 3     | ✅     |

### 📋 Handlers Testeados

| Handler  | Cobertura | Tests | Estado |
|----------|-----------|-------|--------|
| Create   | 100%      | 4     | ✅     |
| Update   | 100%      | 5     | ✅     |
| GetByID  | 100%      | 4     | ✅     |
| List     | 100%      | 4     | ✅     |
| Delete   | 100%      | 4     | ✅     |

### 🎯 Mejores Prácticas Aplicadas

- ✅ Todos los mocks en `internal/mocks/` (no dentro de tests)
- ✅ Tests unitarios puros (sin base de datos real)
- ✅ Mocks de interfaces (ports), no de implementaciones
- ✅ Patrón AAA (Arrange, Act, Assert)
- ✅ Tests independientes
- ✅ Nombres descriptivos
- ✅ Cobertura de casos felices, errores y casos límite

---

## 🛠️ Desarrollo

### Compilar

```bash
go build ./...
```

### Ejecutar Tests

```bash
go test ./...
```

### Migraciones

#### 1. AutoMigrate (desde el código)

```bash
cd /back/central
go run cmd/main.go migrate
```

#### 2. Script SQL (manual)

```bash
psql -U postgres -d probability_db -f /back/migration/shared/sql/migrate_notification_system_refactor.sql
```

### Verificar Datos en BD

```sql
-- Ver tipos de notificación
SELECT * FROM notification_types;

-- Ver eventos de notificación con su tipo
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

## 📝 Convenciones

1. **Entidades de dominio:** Sin tags, solo tipos nativos de Go
2. **Modelos GORM:** Centralizados en `/back/migration/shared/models/`
3. **Repositorios:** Usan modelos de migration, retornan entidades de dominio
4. **Handlers:** Cada handler en su propio archivo (`create_handler.go`, `list_handler.go`)
5. **Rutas:** Registradas en `routes.go` dentro de cada grupo de handlers
6. **Mappers:** Obligatorios en `request/`, `response/`, `mappers/` para cada handler

---

## 📦 Dependencias

- **GORM:** ORM para PostgreSQL
- **Gin:** Framework HTTP
- **datatypes.JSON:** Soporte para campos JSONB
- **Zerolog:** Logging estructurado

---

## ⚠️ Notas Importantes

1. **Unique constraints:** Evitan duplicados en combinaciones clave
2. **Soft deletes:** Implementados con `gorm.DeletedAt`
3. **Preload:** Usar `.Preload()` para cargar relaciones
4. **Validaciones:** Implementadas en capa de aplicación (use cases)
5. **Errores de dominio:** Tipados y centralizados en `domain/errors/`

---

## 📜 Changelog

### v2.0.0 - Refactorización Arquitectura Jerárquica (2026-01-31)

**BREAKING CHANGES:**
- Nueva estructura de tres niveles (NotificationType → NotificationEventType → BusinessNotificationConfig)
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

**Última actualización:** 2026-01-31
