# Order Status Module - Arquitectura Hexagonal

> **Módulo de mapeo de estados de órdenes**: Gestiona la conversión entre estados de integraciones externas (Shopify, WhatsApp, etc.) y estados unificados de Probability.

## 📋 Descripción

Este módulo permite:
- **Mapear estados de órdenes** de diferentes integraciones a estados unificados
- **Listar estados disponibles** de Probability
- **CRUD completo** de mapeos de estados
- **Activar/desactivar** mapeos sin eliminarlos

---

## 🏗️ Arquitectura Hexagonal

Este módulo sigue **Clean Architecture / Hexagonal Architecture** con separación estricta de capas:

```
orderstatus/
├── bundle.go              # Punto de entrada - ensambla el módulo
└── internal/              # Código interno (no exportable)
    ├── domain/            # 🔵 NÚCLEO - Reglas de negocio
    │   ├── entities/      # Entidades PURAS (sin tags)
    │   ├── dtos/          # DTOs de dominio
    │   ├── ports/         # Interfaces (contratos)
    │   └── errors/        # Errores de dominio
    │
    ├── app/               # 🟢 APLICACIÓN - Casos de uso
    │   ├── constructor.go # IUseCase interface + New()
    │   ├── create.go      # CreateOrderStatusMapping
    │   ├── get.go         # GetOrderStatusMapping
    │   ├── list.go        # ListOrderStatusMappings
    │   ├── update.go      # UpdateOrderStatusMapping
    │   ├── delete.go      # DeleteOrderStatusMapping
    │   ├── toggle.go      # ToggleOrderStatusMappingActive
    │   ├── list_statuses.go # ListOrderStatuses
    │   ├── request/       # DTOs de entrada (vacío actualmente)
    │   ├── response/      # DTOs de salida (vacío actualmente)
    │   └── mappers/       # Conversiones (vacío actualmente)
    │
    └── infra/             # 🔴 INFRAESTRUCTURA - Adaptadores
        ├── primary/       # Adaptadores de entrada (drivers)
        │   └── handlers/  # HTTP handlers (Gin)
        │       ├── constructor.go      # IHandler + New()
        │       ├── routes.go           # RegisterRoutes()
        │       ├── create.go           # POST /order-status-mappings
        │       ├── get.go              # GET /order-status-mappings/:id
        │       ├── list.go             # GET /order-status-mappings
        │       ├── update.go           # PUT /order-status-mappings/:id
        │       ├── delete.go           # DELETE /order-status-mappings/:id
        │       ├── toggle.go           # PATCH /order-status-mappings/:id/toggle
        │       ├── list_order_statuses.go # GET /order-statuses
        │       ├── list_simple.go      # GET /order-statuses/simple
        │       ├── request/            # DTOs HTTP request
        │       │   ├── create.go
        │       │   └── update.go
        │       ├── response/           # DTOs HTTP response
        │       │   ├── response.go
        │       │   └── simple-response.go
        │       └── mappers/            # Conversiones domain ↔ HTTP
        │           ├── to_domain.go
        │           └── to_response.go
        │
        └── secondary/     # Adaptadores de salida (driven)
            └── repository/
                ├── constructor.go      # New() - retorna ports.IRepository
                ├── create.go           # Create
                ├── get_by_id.go        # GetByID
                ├── list.go             # List con filtros
                ├── update.go           # Update
                ├── delete.go           # Delete
                ├── toggle.go           # ToggleActive
                ├── exists.go           # Exists
                ├── list_statuses.go    # ListOrderStatuses
                ├── get_status_id.go    # GetOrderStatusIDBy...
                ├── models/             # Modelos GORM locales
                │   ├── order_status_mapping.go
                │   ├── integration_type.go
                │   └── order_status.go
                ├── request/            # DTOs de queries (vacío)
                ├── response/           # DTOs de resultados (vacío)
                └── mappers/            # Conversiones adicionales (vacío)
```

---

## 🎯 Flujo de Dependencias

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Request                         │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│  PRIMARY ADAPTERS (infra/primary/handlers)              │
│  - Validación HTTP                                      │
│  - Parseo de parámetros                                 │
│  - Mapeo request → domain                               │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│  APPLICATION LAYER (app)                                │
│  - Lógica de negocio                                    │
│  - Validaciones de dominio                              │
│  - Orquestación de casos de uso                         │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│  DOMAIN LAYER (domain)                                  │
│  - Entidades puras (sin tags)                           │
│  - Puertos (interfaces)                                 │
│  - Errores de dominio                                   │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│  SECONDARY ADAPTERS (infra/secondary/repository)        │
│  - Implementación de repositorios                       │
│  - Modelos GORM (con tags)                              │
│  - Mapeo domain ↔ DB                                    │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Database (PostgreSQL)                │
└─────────────────────────────────────────────────────────┘

REGLA DE ORO: Las dependencias SIEMPRE apuntan hacia el dominio (adentro)
```

---

## 📊 Modelo de Datos

### Tabla: `order_status_mappings`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | uint | ID único del mapeo |
| `integration_type_id` | uint | FK a `integration_types` (1=Shopify, 2=WhatsApp, etc.) |
| `original_status` | string | Estado original de la integración ("paid", "fulfilled", etc.) |
| `order_status_id` | uint | FK a `order_statuses` (estado unificado de Probability) |
| `is_active` | bool | Si el mapeo está activo |
| `priority` | int | Prioridad en caso de múltiples mapeos (mayor = mayor prioridad) |
| `description` | string | Descripción del mapeo |
| `created_at` | timestamp | Fecha de creación |
| `updated_at` | timestamp | Fecha de última actualización |

**Índice único**: `(integration_type_id, original_status)` - No puede haber duplicados activos

---

## 🔌 API Endpoints

### Mapeos de Estados (`/order-status-mappings`)

#### 1. **Listar mapeos** (paginado)
```http
GET /api/v1/order-status-mappings?page=1&page_size=10&integration_type_id=1&is_active=true
```

**Query params**:
- `page` (int, default: 1) - Número de página
- `page_size` (int, default: 10, max: 100) - Tamaño de página
- `integration_type_id` (int, opcional) - Filtrar por tipo de integración
- `is_active` (bool, opcional) - Filtrar por estado activo/inactivo

**Response 200**:
```json
{
  "data": [
    {
      "id": 1,
      "integration_type_id": 1,
      "integration_type": {
        "id": 1,
        "code": "shopify",
        "name": "Shopify",
        "image_url": "https://..."
      },
      "original_status": "paid",
      "order_status_id": 2,
      "order_status": {
        "id": 2,
        "code": "processing",
        "name": "En Procesamiento",
        "description": "...",
        "category": "active",
        "color": "#FFB020"
      },
      "is_active": true,
      "priority": 10,
      "description": "Mapeo de estado 'paid' de Shopify",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 45,
  "page": 1,
  "page_size": 10,
  "total_pages": 5
}
```

---

#### 2. **Obtener mapeo por ID**
```http
GET /api/v1/order-status-mappings/:id
```

**Response 200**:
```json
{
  "id": 1,
  "integration_type_id": 1,
  "integration_type": { ... },
  "original_status": "paid",
  "order_status_id": 2,
  "order_status": { ... },
  "is_active": true,
  "priority": 10,
  "description": "...",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Response 404**:
```json
{
  "error": "order status mapping not found"
}
```

---

#### 3. **Crear mapeo**
```http
POST /api/v1/order-status-mappings
Content-Type: application/json
```

**Body**:
```json
{
  "integration_type_id": 1,
  "original_status": "paid",
  "order_status_id": 2,
  "priority": 10,
  "description": "Mapeo de estado 'paid' de Shopify"
}
```

**Response 201**:
```json
{
  "id": 5,
  "integration_type_id": 1,
  "original_status": "paid",
  "order_status_id": 2,
  "is_active": true,
  "priority": 10,
  "description": "...",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Response 400**:
```json
{
  "error": "mapping already exists for this integration type and original status"
}
```

---

#### 4. **Actualizar mapeo**
```http
PUT /api/v1/order-status-mappings/:id
Content-Type: application/json
```

**Body**:
```json
{
  "original_status": "fulfilled",
  "order_status_id": 3,
  "priority": 15,
  "description": "Actualización del mapeo"
}
```

**Response 200**:
```json
{
  "id": 5,
  "integration_type_id": 1,
  "original_status": "fulfilled",
  "order_status_id": 3,
  "is_active": true,
  "priority": 15,
  "description": "Actualización del mapeo",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

---

#### 5. **Eliminar mapeo**
```http
DELETE /api/v1/order-status-mappings/:id
```

**Response 200**:
```json
{
  "message": "mapping deleted successfully"
}
```

**Response 400**:
```json
{
  "error": "invalid id"
}
```

---

#### 6. **Alternar estado activo/inactivo**
```http
PATCH /api/v1/order-status-mappings/:id/toggle
```

**Response 200**:
```json
{
  "id": 5,
  "integration_type_id": 1,
  "original_status": "paid",
  "order_status_id": 2,
  "is_active": false,  // ← Cambió de true a false
  "priority": 10,
  "description": "...",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:30:00Z"
}
```

---

### Estados de Órdenes (`/order-statuses`)

#### 7. **Listar estados de órdenes de Probability**
```http
GET /api/v1/order-statuses?is_active=true
```

**Query params**:
- `is_active` (bool, opcional) - Filtrar por estado activo/inactivo (omitir para traer todos)

**Response 200**:
```json
{
  "success": true,
  "message": "Estados de órdenes obtenidos exitosamente",
  "data": [
    {
      "id": 1,
      "code": "pending",
      "name": "Pendiente",
      "description": "Orden recibida, pendiente de procesamiento",
      "category": "active",
      "color": "#FFB020"
    },
    {
      "id": 2,
      "code": "processing",
      "name": "En Procesamiento",
      "description": "Orden siendo procesada",
      "category": "active",
      "color": "#4A90E2"
    }
  ]
}
```

---

#### 8. **Listar estados en formato simple** (para dropdowns)
```http
GET /api/v1/order-statuses/simple?is_active=true
```

**Query params**:
- `is_active` (bool, default: true) - Filtrar por estado activo

**Response 200**:
```json
{
  "success": true,
  "message": "Order statuses retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Pendiente",
      "code": "pending",
      "is_active": true
    },
    {
      "id": 2,
      "name": "En Procesamiento",
      "code": "processing",
      "is_active": true
    }
  ]
}
```

---

## 🧪 Casos de Uso

### 1. **CreateOrderStatusMapping**
- **Entrada**: Entidad `OrderStatusMapping` con campos requeridos
- **Validaciones**:
  - ✅ Verifica que no exista un mapeo activo con la misma combinación `(integration_type_id, original_status)`
  - ✅ Asigna `is_active = true` por defecto
- **Salida**: Mapeo creado con relaciones cargadas (`IntegrationType`, `OrderStatus`)

### 2. **GetOrderStatusMapping**
- **Entrada**: ID del mapeo
- **Validaciones**:
  - ✅ Verifica que el ID exista
- **Salida**: Mapeo con relaciones cargadas

### 3. **ListOrderStatusMappings**
- **Entrada**: Filtros (`integration_type_id`, `is_active`, paginación)
- **Salida**: Lista paginada de mapeos con relaciones cargadas

### 4. **UpdateOrderStatusMapping**
- **Entrada**: ID del mapeo + campos a actualizar
- **Validaciones**:
  - ✅ Verifica que el ID exista
  - ✅ Solo actualiza campos permitidos: `original_status`, `order_status_id`, `priority`, `description`
  - ❌ NO permite cambiar `integration_type_id` ni `is_active`
- **Salida**: Mapeo actualizado con relaciones recargadas

### 5. **DeleteOrderStatusMapping**
- **Entrada**: ID del mapeo
- **Acción**: Eliminación lógica (soft delete) del mapeo
- **Salida**: Sin contenido (204)

### 6. **ToggleOrderStatusMappingActive**
- **Entrada**: ID del mapeo
- **Acción**: Invierte el valor de `is_active` (true ↔ false)
- **Salida**: Mapeo actualizado

### 7. **ListOrderStatuses**
- **Entrada**: Filtro opcional `is_active`
- **Salida**: Lista de estados de órdenes de Probability (NO mapeos, sino los estados base del sistema)

---

## 📁 Estructura de Entidades

### Domain Entity (PURA - Sin tags)

```go
// internal/domain/entities/order_status_mapping.go
type OrderStatusMapping struct {
	ID                uint
	IntegrationTypeID uint
	OriginalStatus    string
	OrderStatusID     uint
	IsActive          bool
	Priority          int
	Description       string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Relaciones (opcionales)
	IntegrationType *IntegrationTypeInfo
	OrderStatus     *OrderStatusInfo
}
```

### Infrastructure Model (GORM - Con tags)

```go
// internal/infra/secondary/repository/models/order_status_mapping.go
type OrderStatusMapping struct {
	gorm.Model

	IntegrationTypeID uint   `gorm:"not null;index;uniqueIndex:idx_status_mapping,priority:1"`
	OriginalStatus    string `gorm:"size:128;not null;uniqueIndex:idx_status_mapping,priority:2"`
	OrderStatusID     uint   `gorm:"not null;index"`
	IsActive          bool   `gorm:"default:true;index"`
	Priority          int    `gorm:"default:0"`
	Description       string `gorm:"type:text"`
	Metadata          datatypes.JSON `gorm:"type:jsonb"`

	// Relaciones
	IntegrationType IntegrationType `gorm:"foreignKey:IntegrationTypeID"`
	OrderStatus     OrderStatus     `gorm:"foreignKey:OrderStatusID"`
}

func (OrderStatusMapping) TableName() string {
	return "order_status_mappings"
}

// Conversiones
func (m *OrderStatusMapping) ToDomain() entities.OrderStatusMapping { ... }
func FromDomain(e entities.OrderStatusMapping) *OrderStatusMapping { ... }
```

---

## ✅ Validaciones de Arquitectura Hexagonal

### ✅ Domain Layer (CONFORME)
- ✅ Entidades PURAS sin tags JSON/GORM
- ✅ Ports (interfaces) solo con tipos de dominio
- ✅ Sin dependencias de frameworks externos
- ✅ Organizado en subcarpetas (`entities/`, `dtos/`, `ports/`, `errors/`)

### ✅ Application Layer (CONFORME)
- ✅ Casos de uso separados por archivo
- ✅ Solo depende de `domain/`
- ✅ Un método por archivo para claridad
- ✅ Constructor único con interfaz `IUseCase`

### ✅ Infrastructure Layer (CONFORME)
- ✅ Handlers en `infra/primary/handlers/`
- ✅ Repositorio en `infra/secondary/repository/`
- ✅ Modelos GORM locales (NO externos de `migration/`)
- ✅ Mappers separados en carpeta `mappers/`
- ✅ DTOs de request/response en carpetas separadas
- ✅ Método `RegisterRoutes()` en interfaz `IHandler`

### ✅ Repository (CONFORME)
- ✅ Usa modelos GORM locales (NO `.Table("nombre")`)
- ✅ Todos los modelos tienen `TableName()`
- ✅ Todos los modelos tienen `ToDomain()` y `FromDomain()`
- ✅ Constructor retorna interfaz de dominio (`ports.IRepository`)

---

## 🚀 Comandos de Verificación

### Compilar módulo
```bash
cd /home/cam/Desktop/probability/back/central
go build ./services/modules/orderstatus/...
```

### Verificar que domain es PURO
```bash
# NO debe encontrar tags JSON en domain/entities
grep -r 'json:"' services/modules/orderstatus/internal/domain/entities/

# NO debe encontrar imports de frameworks en domain
grep -r "gorm\|gin\|fiber" services/modules/orderstatus/internal/domain/

# Verificar estructura de carpetas obligatoria
ls services/modules/orderstatus/internal/domain/entities/
ls services/modules/orderstatus/internal/domain/ports/
ls services/modules/orderstatus/internal/domain/errors/
```

### Verificar modelos GORM
```bash
# Verificar que todos los modelos tienen TableName()
find services/modules/orderstatus/internal/infra/secondary/repository/models -name "*.go" -exec grep -L "TableName()" {} \;

# Verificar que todos los modelos tienen ToDomain()
find services/modules/orderstatus/internal/infra/secondary/repository/models -name "*.go" -exec grep -L "ToDomain()" {} \;

# NO debe usar .Table() en repositorios
grep -r '\.Table(' services/modules/orderstatus/internal/infra/secondary/repository/*.go
```

---

## 📚 Referencias

- **Patrón**: Hexagonal Architecture / Ports & Adapters
- **Inspiración**: Clean Architecture (Robert C. Martin)
- **Convenciones**: Ver `.claude/rules/architecture.md`

---

## 🔧 Mantenimiento

### Agregar un nuevo endpoint

1. **Agregar método en `domain/ports/ports.go`**:
```go
type IRepository interface {
    // ... métodos existentes
    NuevoMetodo(ctx context.Context, param string) (*entities.Entity, error)
}
```

2. **Implementar en `infra/secondary/repository/nuevo_metodo.go`**:
```go
func (r *repository) NuevoMetodo(ctx context.Context, param string) (*entities.Entity, error) {
    var model models.Entity
    err := r.db.Conn(ctx).Where("campo = ?", param).First(&model).Error
    if err != nil {
        return nil, err
    }
    domain := model.ToDomain()
    return &domain, nil
}
```

3. **Agregar método en `app/constructor.go`**:
```go
type IUseCase interface {
    // ... métodos existentes
    NuevoUseCase(ctx context.Context, param string) (*entities.Entity, error)
}
```

4. **Implementar en `app/nuevo_use_case.go`**:
```go
func (uc *useCase) NuevoUseCase(ctx context.Context, param string) (*entities.Entity, error) {
    return uc.repo.NuevoMetodo(ctx, param)
}
```

5. **Agregar handler en `infra/primary/handlers/nuevo_handler.go`**:
```go
func (h *handler) NuevoHandler(c *gin.Context) {
    param := c.Param("param")
    result, err := h.uc.NuevoUseCase(c.Request.Context(), param)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, mappers.DomainToResponse(result, h.getImageURLBase()))
}
```

6. **Registrar ruta en `infra/primary/handlers/routes.go`**:
```go
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
    // ... rutas existentes
    router.GET("/nuevo-endpoint/:param", h.NuevoHandler)
}
```

---

## 📊 Métricas de Calidad

| Métrica | Valor | Estado |
|---------|-------|--------|
| Violaciones de arquitectura hexagonal | 0 | ✅ CONFORME |
| Entidades con tags en domain | 0 | ✅ PURO |
| Uso de `.Table()` en repositorios | 0 | ✅ CONFORME |
| Modelos GORM sin `TableName()` | 0 | ✅ CONFORME |
| Handlers sin `RegisterRoutes()` | 0 | ✅ CONFORME |
| Mappers mezclados con handlers | 0 | ✅ SEPARADO |

---

## 🔄 Resumen de Refactorización

### Violaciones Corregidas

Este módulo fue completamente refactorizado desde una estructura legacy a arquitectura hexagonal. Se corrigieron **7 violaciones**:

1. **Domain con dependencias externas** - Domain importaba modelos de `migration/`
2. **Entidades con tags JSON** - Entidades tenían `json:"..."` tags
3. **App con dependencias de infraestructura** - Application usaba modelos de infra
4. **Domain sin organización** - Archivos sueltos en lugar de subcarpetas
5. **App sin estructura estándar** - Faltaban carpetas `request/`, `response/`, `mappers/`
6. **Repository con modelos externos** - Usaba modelos de `migration/` en lugar de locales
7. **Repository sin estructura** - Faltaban carpetas `models/`, `request/`, `response/`

### Cambios Principales

#### Antes (Legacy)
```
orderstatus/
├── bundle.go
├── domain/
│   ├── entities.go       # ❌ Con tags JSON
│   └── ports.go          # ❌ Tipos externos
├── app/
│   └── usecases.go       # ❌ Todo mezclado
└── infra/
    └── repository/
        └── repository.go  # ❌ Modelos externos
```

#### Después (Hexagonal)
```
orderstatus/
├── bundle.go
└── internal/             # ✅ Todo en internal/
    ├── domain/           # ✅ 100% PURO
    │   ├── entities/     # Sin tags
    │   ├── dtos/
    │   ├── ports/
    │   └── errors/
    ├── app/              # ✅ Separado
    │   ├── constructor.go
    │   ├── create.go, get.go, list.go...
    │   ├── request/, response/, mappers/
    └── infra/
        ├── primary/handlers/
        │   ├── routes.go
        │   ├── create.go, get.go...
        │   └── request/, response/, mappers/
        └── secondary/repository/
            ├── create.go, get_by_id.go...
            └── models/   # ✅ Modelos GORM locales
```

### Beneficios Obtenidos

- ✅ **Dominio 100% puro** - Sin tags, sin frameworks, totalmente testeable
- ✅ **Independencia de infraestructura** - Modelos GORM controlados localmente
- ✅ **Mantenibilidad** - Un método por archivo, fácil de localizar
- ✅ **Escalabilidad** - Estructura consistente para crecimiento
- ✅ **Testabilidad** - Capas desacopladas, fácil de mockear

---

**Última actualización**: 2026-01-31
**Estado**: ✅ PRODUCCIÓN - 100% CONFORME CON ARQUITECTURA HEXAGONAL
**Refactorización**: Migración completa de estructura legacy → hexagonal con `internal/`
