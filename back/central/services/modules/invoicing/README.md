# Módulo de Facturación (modules/invoicing)

## Propósito

Gestión centralizada de facturas electrónicas para TODOS los proveedores de facturación (Softpymes, Alegra, Siigo, etc.).

## Responsabilidades

### ✅ Que SÍ hace este módulo

- **CRUD de facturas (Invoice)**: Crear, listar, obtener, cancelar y reintentar facturas
- **CRUD de notas de crédito (CreditNote)**: Gestión completa de notas de crédito
- **Gestión de configuraciones (InvoicingConfig)**: Configuraciones por integración
- **Sincronización automática**: Consumidores de RabbitMQ para facturación automática
- **Listado general con filtros**: Buscar facturas por negocio, estado, integración, etc.
- **Reportes y estadísticas**: KPIs, tendencias y análisis de facturación

### ❌ Que NO hace este módulo

- **NO gestiona credenciales de proveedores** (ver `integrations/core`)
- **NO implementa lógica específica de proveedores** (ver `integrations/invoicing/*`)
- **NO registra nuevos tipos de proveedores** (ver `integrations/core`)

## Arquitectura

Este módulo sigue **Arquitectura Hexagonal (Ports & Adapters)**:

```
modules/invoicing/
+-- bundle.go              # Ensambla el módulo
+-- internal/
    +-- domain/            # Núcleo - Reglas de negocio
    |   +-- entities/      # Entidades PURAS (sin tags)
    |   +-- dtos/          # Data Transfer Objects
    |   +-- ports/         # Interfaces (contratos)
    |   +-- errors/        # Errores de dominio
    |   +-- constants/     # Constantes
    +-- app/               # Casos de uso
    |   +-- constructor.go
    |   +-- create_invoice.go
    |   +-- bulk_create_invoices_async.go # ✨ NUEVO - Procesamiento asíncrono
    |   +-- get_summary.go      # ✨ NUEVO - Resumen de KPIs
    |   +-- get_stats.go        # ✨ NUEVO - Estadísticas detalladas
    |   +-- get_trends.go       # ✨ NUEVO - Tendencias temporales
    |   +-- deprecated_providers.go  # Métodos deprecados (retornan error)
    +-- infra/
        +-- primary/       # Adaptadores de entrada
        |   +-- handlers/  # HTTP handlers (Gin)
        |   +-- queue/     # Consumers (RabbitMQ)
        |       +-- consumer/
        |           +-- retry_consumer.go
        |           +-- bulk_invoice_consumer.go
        +-- secondary/     # Adaptadores de salida
            +-- repository/ # Repositorios DB (GORM)
            +-- queue/      # Publishers (RabbitMQ)
            +-- redis/      # SSE Publisher (Redis Pub/Sub)
                +-- sse_publisher.go
```

## Relación con Integraciones

Este módulo **delega la ejecución real** a proveedores específicos mediante `integrations/core`:

```
modules/invoicing (lógica de negocio)
        v usa
integrations/core (orquestador)
        v delega a
integrations/invoicing/softpymes (proveedor específico)
integrations/invoicing/alegra (futuro)
integrations/invoicing/siigo (futuro)
```

## Endpoints HTTP

### Facturas

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/invoicing/invoices` | Crear factura manual |
| `GET` | `/invoicing/invoices` | Listar facturas con filtros |
| `GET` | `/invoicing/invoices/:id` | Obtener detalle de factura |
| `POST` | `/invoicing/invoices/:id/cancel` | Cancelar factura |
| `POST` | `/invoicing/invoices/:id/retry` | Reintentar emisión de factura |
| `POST` | `/invoicing/invoices/:id/credit-notes` | Crear nota de crédito |

#### Filtros disponibles para listado

```
GET /invoicing/invoices?business_id=1&status=issued&integration_id=2&invoicing_integration_id=3
```

- `business_id` (uint): Filtrar por negocio
- `status` (string): Estados: `pending`, `issued`, `failed`, `cancelled`
- `integration_id` (uint): Filtrar por integración origen (Shopify, MercadoLibre, etc.)
- `invoicing_integration_id` (uint): Filtrar por proveedor de facturación (Softpymes, Alegra, etc.)
- `order_id` (string): Buscar factura de una orden específica
- `created_after` (date): Facturas creadas después de esta fecha
- `created_before` (date): Facturas creadas antes de esta fecha

### Estadísticas y Resúmenes (✨ NUEVO)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/invoicing/summary` | Resumen general con KPIs principales |
| `GET` | `/invoicing/stats` | Estadísticas detalladas para dashboards |
| `GET` | `/invoicing/trends` | Tendencias temporales para gráficos |

#### 1. Resumen general (Summary)

```bash
GET /invoicing/summary?business_id=1&period=month
```

**Query Parameters:**
- `business_id` (uint, requerido): ID del negocio
- `period` (string, opcional): Período a analizar
  - `today`: Hoy
  - `week`: Esta semana
  - `month`: Este mes (default)
  - `year`: Este año
  - `all`: Últimos 10 años

**Response:**
```json
{
  "period": {
    "start": "2026-01-01T00:00:00Z",
    "end": "2026-01-31T23:59:59Z",
    "label": "Enero 2026"
  },
  "totals": {
    "total_invoices": 150,
    "total_amount": 45000000,
    "issued_count": 120,
    "issued_amount": 42000000,
    "failed_count": 20,
    "pending_count": 10
  },
  "by_status": [
    { "status": "issued", "count": 120, "amount": 42000000, "percentage": 80 },
    { "status": "failed", "count": 20, "amount": 2000000, "percentage": 13.3 }
  ],
  "by_provider": [
    { "provider_id": 5, "provider_name": "Softpymes", "count": 100, "amount": 35000000 }
  ],
  "recent_failures": [
    {
      "invoice_id": 123,
      "order_id": "456",
      "amount": 100000,
      "error": "API timeout",
      "failed_at": "2026-01-31T10:00:00Z"
    }
  ]
}
```

#### 2. Estadísticas detalladas (Stats)

```bash
GET /invoicing/stats?business_id=1&start_date=2026-01-01&end_date=2026-01-31
```

**Query Parameters:**
- `business_id` (uint, requerido): ID del negocio
- `start_date` (string, opcional): Fecha de inicio (formato: `YYYY-MM-DD`)
- `end_date` (string, opcional): Fecha de fin (formato: `YYYY-MM-DD`)
- `integration_id` (uint, opcional): Filtrar por integración origen
- `invoicing_integration_id` (uint, opcional): Filtrar por proveedor de facturación

**Response:**
```json
{
  "summary": {
    "total_invoices": 500,
    "total_amount": 150000000,
    "avg_amount": 300000,
    "success_rate": 85.5
  },
  "top_customers": [
    { "customer_name": "Cliente A", "invoice_count": 50, "total_amount": 15000000 }
  ],
  "monthly_breakdown": [
    { "month": "2026-01", "count": 150, "amount": 45000000, "success_rate": 90 }
  ],
  "failure_analysis": {
    "total_failures": 72,
    "by_reason": [
      { "reason": "API timeout", "count": 30, "percentage": 41.7 }
    ]
  },
  "processing_times": {
    "avg_seconds": 2.5,
    "p50_seconds": 2.0,
    "p95_seconds": 5.0,
    "p99_seconds": 10.0
  }
}
```

#### 3. Tendencias temporales (Trends)

```bash
GET /invoicing/trends?business_id=1&start_date=2026-01-01&end_date=2026-01-31&granularity=day&metric=count
```

**Query Parameters:**
- `business_id` (uint, requerido): ID del negocio
- `start_date` (string, requerido): Fecha de inicio (formato: `YYYY-MM-DD`)
- `end_date` (string, requerido): Fecha de fin (formato: `YYYY-MM-DD`)
- `granularity` (string, opcional): Granularidad de datos
  - `day`: Por día (default)
  - `week`: Por semana
  - `month`: Por mes
- `metric` (string, opcional): Métrica a visualizar
  - `count`: Cantidad de facturas (default)
  - `amount`: Monto total facturado
  - `success_rate`: Tasa de éxito (%)

**Response:**
```json
{
  "metric": "count",
  "granularity": "day",
  "data_points": [
    { "date": "2026-01-01", "value": 10, "success_rate": 90 },
    { "date": "2026-01-02", "value": 12, "success_rate": 85 }
  ],
  "trend": {
    "direction": "up",
    "percentage_change": 15.5,
    "comparison_period": "previous_period"
  }
}
```

### Configuraciones

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/invoicing/configs` | Crear configuración |
| `GET` | `/invoicing/configs` | Listar configuraciones |
| `GET` | `/invoicing/configs/:id` | Obtener configuración |
| `PUT` | `/invoicing/configs/:id` | Actualizar configuración |
| `DELETE` | `/invoicing/configs/:id` | Eliminar configuración |

#### Sistema de Filtros Avanzados

Las configuraciones de facturación incluyen un sistema de **filtros avanzados** que permite controlar qué órdenes se facturan automáticamente.

**Ejemplo de configuración con filtros:**

```json
{
  "business_id": 1,
  "integration_id": 5,
  "invoicing_provider_id": 10,
  "enabled": true,
  "auto_invoice": true,
  "filters": {
    "min_amount": 100000,
    "payment_status": "paid",
    "order_types": ["delivery"],
    "exclude_products": ["GIFT-CARD-001"],
    "shipping_regions": ["Bogotá", "Medellín", "Cali"]
  }
}
```

**Filtros disponibles:**

| Categoría | Filtro | Tipo | Descripción |
|-----------|--------|------|-------------|
| **Monto** | `min_amount` | `float64` | Monto mínimo para facturar |
| | `max_amount` | `float64` | Monto máximo para facturar |
| **Pago** | `payment_status` | `string` | Estado de pago requerido (`"paid"`) |
| | `payment_methods` | `[]uint` | IDs de métodos de pago permitidos |
| **Orden** | `order_types` | `[]string` | Tipos de orden permitidos |
| | `exclude_statuses` | `[]string` | Estados de orden a excluir |
| **Productos** | `exclude_products` | `[]string` | SKUs a excluir |
| | `include_products_only` | `[]string` | Solo estos SKUs |
| | `min_items_count` | `int` | Mínimo de items en la orden |
| | `max_items_count` | `int` | Máximo de items en la orden |
| **Cliente** | `customer_types` | `[]string` | Tipos de cliente permitidos |
| | `exclude_customer_ids` | `[]string` | IDs de clientes a excluir |
| **Ubicación** | `shipping_regions` | `[]string` | Regiones/departamentos permitidos |
| **Fecha** | `date_range` | `object` | Rango de fechas permitido |

**Ejemplos de uso:**

1. **Ecommerce básico**: Solo facturar órdenes pagadas mayores a $100.000
```json
{
  "filters": {
    "min_amount": 100000,
    "payment_status": "paid"
  }
}
```

2. **Marketplace B2B**: Solo clientes empresariales, órdenes grandes con mínimo 5 productos
```json
{
  "filters": {
    "min_amount": 500000,
    "customer_types": ["juridica"],
    "min_items_count": 5,
    "exclude_statuses": ["cancelled", "refunded"]
  }
}
```

3. **Tienda regional**: Solo delivery en ciudades específicas, sin gift cards
```json
{
  "filters": {
    "order_types": ["delivery"],
    "exclude_products": ["GIFT-CARD-001"],
    "shipping_regions": ["Bogotá", "Medellín", "Cali"]
  }
}
```

**Nota:** Los filtros se evalúan en modo AND (todos deben cumplirse). Si algún filtro falla, la orden no se factura automáticamente.

### Notas de Crédito

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/invoicing/credit-notes` | Listar notas de crédito |
| `GET` | `/invoicing/credit-notes/:id` | Obtener detalle de nota de crédito |

### Proveedores (⚠️ DEPRECADO)

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|--------|
| `POST` | `/invoicing/providers` | Crear proveedor | ⚠️ DEPRECATED |
| `GET` | `/invoicing/providers` | Listar proveedores | ⚠️ DEPRECATED |
| `GET` | `/invoicing/providers/:id` | Obtener proveedor | ⚠️ DEPRECATED |
| `PUT` | `/invoicing/providers/:id` | Actualizar proveedor | ⚠️ DEPRECATED |
| `POST` | `/invoicing/providers/:id/test` | Probar conexión | ⚠️ DEPRECATED |

**⚠️ NOTA**: Estos endpoints están deprecados. Usar `integrations/core` para gestión de proveedores de facturación.

## Migración: Gestión de Proveedores

### ❌ Antes (Deprecado)

```bash
POST /invoicing/providers
GET /invoicing/providers
```

### ✅ Ahora (Usar integrations/core)

```bash
GET /integrations?category=invoicing&business_id=1
POST /integrations  # Con category_id=invoicing
```

## Estados de Factura

| Estado | Descripción |
|--------|-------------|
| `pending` | Factura pendiente de emisión |
| `issued` | Factura emitida exitosamente |
| `failed` | Error al emitir la factura |
| `cancelled` | Factura cancelada |

## Casos de Uso Principales

### 1. Crear factura para una orden

```go
invoice, err := useCase.CreateInvoice(ctx, &dtos.CreateInvoiceDTO{
    OrderID: "uuid-de-la-orden",
    InvoicingIntegrationID: 5, // ID de integración de Softpymes
})
```

### 2. Listar facturas con filtros

```go
invoices, err := useCase.ListInvoices(ctx, map[string]interface{}{
    "business_id": 1,
    "status": "issued",
    "integration_id": 2,
})
```

### 3. Obtener resumen de facturación

```go
summary, err := useCase.GetSummary(ctx, businessID, "month")
```

### 4. Obtener estadísticas detalladas

```go
stats, err := useCase.GetDetailedStats(ctx, businessID, map[string]interface{}{
    "start_date": "2026-01-01",
    "end_date": "2026-01-31",
})
```

### 5. Obtener tendencias temporales

```go
trends, err := useCase.GetTrends(ctx, businessID, "2026-01-01", "2026-01-31", "day", "count")
```

## Eventos Publicados

### RabbitMQ (Procesamiento Asíncrono)

| Evento | Descripción |
|--------|-------------|
| `invoice.created` | Factura creada exitosamente |
| `invoice.failed` | Error al crear factura |
| `invoice.cancelled` | Factura cancelada |
| `credit_note.created` | Nota de crédito creada |

### Redis Pub/Sub -> SSE (Notificaciones en Tiempo Real)

Además de RabbitMQ, el módulo publica eventos a **Redis Pub/Sub** para que el frontend reciba actualizaciones en tiempo real via **Server-Sent Events (SSE)**.

| Evento | Descripción |
|--------|-------------|
| `invoice.created` | Factura emitida exitosamente |
| `invoice.failed` | Error al emitir factura (con mensaje de error) |
| `invoice.cancelled` | Factura cancelada |
| `credit_note.created` | Nota de crédito creada |
| `bulk_job.progress` | Progreso de job de facturación masiva |
| `bulk_job.completed` | Job de facturación masiva finalizado |

**Canal Redis:** `probability:invoicing:events` (configurable via env)

## Flujo SSE (Notificaciones en Tiempo Real)

```
+-----------------------------------------------------------------+
|                    FLUJO DE EVENTOS SSE                         |
+-----------------------------------------------------------------+
|                                                                 |
|  CreateInvoice / CancelInvoice / BulkJob                        |
|       |                                                         |
|       +--> RabbitMQ (procesamiento asíncrono, sin cambios)       |
|       |                                                         |
|       +--> Redis Pub/Sub ---> Events Module suscribe              |
|               |                  |                              |
|               |                  +--> InvoiceEventSubscriber     |
|               |                  |     (lee del canal Redis)    |
|               |                  |                              |
|               |                  +--> InvoiceEventConsumer       |
|               |                  |     (convierte a Event)      |
|               |                  |                              |
|               |                  +--> EventManager broadcast     |
|               |                        |                        |
|               |                        +--> SSE connections      |
|               |                              |                  |
|               |                              +--> Frontend       |
|               |                                   useInvoiceSSE |
|               |                                                 |
|  Canal: probability:invoicing:events                            |
+-----------------------------------------------------------------+
```

### Arquitectura del Publisher SSE

El módulo utiliza `IInvoiceSSEPublisher` (definido en `domain/ports/ports.go`) con dos implementaciones:

- **`SSEPublisher`** (`infra/secondary/redis/sse_publisher.go`): Publica a Redis Pub/Sub de forma no-bloqueante (goroutine). Si falla, solo registra un log de error sin afectar el flujo principal.
- **`noopSSEPublisher`**: Implementación vacía para cuando Redis no está disponible. El módulo funciona normalmente sin SSE.

### Formato del Evento JSON

```json
{
  "id": "20260208143025-a7b3c9d2",
  "event_type": "invoice.created",
  "business_id": 5,
  "timestamp": "2026-02-08T14:30:25Z",
  "data": {
    "invoice_id": 123,
    "order_id": "ord-456",
    "invoice_number": "FV-001",
    "total_amount": 150.50,
    "currency": "COP",
    "status": "issued",
    "customer_name": "Juan Pérez",
    "external_url": "https://..."
  }
}
```

Para eventos de bulk job (`bulk_job.progress` / `bulk_job.completed`):

```json
{
  "id": "20260208143030-b8c4d0e3",
  "event_type": "bulk_job.progress",
  "business_id": 5,
  "timestamp": "2026-02-08T14:30:30Z",
  "data": {
    "job_id": 42,
    "total_orders": 100,
    "processed": 45,
    "successful": 40,
    "failed": 5,
    "progress": 45,
    "status": "processing"
  }
}
```

### Integración Frontend

El frontend usa el hook `useInvoiceSSE` que se conecta al endpoint SSE existente del módulo Events:

```
GET /api/v1/notify/sse/order-notify?business_id=1&event_types=invoice.created,invoice.failed,bulk_job.progress,bulk_job.completed
```

- **InvoiceList**: Escucha `invoice.created`, `invoice.failed`, `invoice.cancelled` para refrescar la lista y mostrar toasts.
- **BulkCreateInvoiceModal**: Escucha `bulk_job.progress` y `bulk_job.completed` para mostrar una barra de progreso en tiempo real.

## Variables de Entorno

```env
# Base de datos (compartida con otros módulos)
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=probability

# RabbitMQ (para eventos)
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASS=admin
# Redis SSE (notificaciones en tiempo real)
REDIS_INVOICE_EVENTS_CHANNEL=probability:invoicing:events
```

## Testing

```bash
# Ejecutar tests del módulo
go test ./services/modules/invoicing/...

# Ejecutar tests con cobertura
go test -cover ./services/modules/invoicing/...
```

## Ejemplos de Uso

### Ejemplo 1: Obtener resumen del mes actual

```bash
curl "http://localhost:8080/api/v1/invoicing/summary?business_id=1&period=month"
```

### Ejemplo 2: Listar facturas emitidas

```bash
curl "http://localhost:8080/api/v1/invoicing/invoices?business_id=1&status=issued"
```

### Ejemplo 3: Obtener tendencias de los últimos 30 días

```bash
curl "http://localhost:8080/api/v1/invoicing/trends?business_id=1&start_date=2026-01-01&end_date=2026-01-31&granularity=day&metric=count"
```

## Troubleshooting

### Error: "Provider not configured"

**Causa**: No existe una configuración de facturación para la integración de la orden.

**Solución**:
1. Crear configuración de facturación:
```bash
POST /invoicing/configs
{
  "integration_id": 2,
  "invoicing_integration_id": 5,
  "enabled": true
}
```

### Error: "Gestión de proveedores deprecada"

**Causa**: Intentando usar endpoints deprecados de `/invoicing/providers`.

**Solución**: Migrar a `integrations/core`:
```bash
# En lugar de:
POST /invoicing/providers

# Usar:
POST /integrations
{
  "business_id": 1,
  "integration_type_id": 10,  # Tipo "Softpymes"
  "category_id": 2,            # Categoría "Invoicing"
  "credentials": { ... }
}
```

## Roadmap

### ✅ Completado

- [x] Sincronización automática vía RabbitMQ
- [x] Notificaciones en tiempo real via SSE (Redis Pub/Sub -> Events Module -> Frontend)
- [x] Facturación masiva asíncrona con progreso en tiempo real
- [x] Creación masiva de facturas desde órdenes (bulk)

### 🚧 En Progreso

- [ ] Dashboard interactivo de facturación (frontend)
- [ ] Exportación de reportes (PDF, Excel)

### 📋 Planificado

- [ ] Soporte para facturación internacional
- [ ] Integración con más proveedores (Alegra, Siigo, etc.)
- [ ] Facturación recurrente/suscripciones
- [ ] Webhooks para notificaciones en tiempo real

## Contribuir

Al modificar este módulo, asegurarse de:

1. Seguir arquitectura hexagonal (no mezclar capas)
2. Mantener entidades del dominio sin tags de infraestructura
3. Actualizar documentación si se agregan endpoints
4. Escribir tests unitarios para nuevos casos de uso
5. No agregar lógica específica de proveedores aquí (usar `integrations/invoicing/*`)

## Última Actualización

**Fecha**: 2026-02-08

**Cambios recientes**:
- Notificaciones en tiempo real via SSE (Redis Pub/Sub)
- Facturación masiva asíncrona con barra de progreso en frontend
- Hook `useInvoiceSSE` para integración frontend
- Noop publisher para degradación elegante sin Redis
