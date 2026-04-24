# Orders Module

Sistema central de gestión de pedidos multi-canal para Probability. Recibe, normaliza, valida y gestiona órdenes provenientes de múltiples plataformas de e-commerce (Shopify, Amazon, MercadoLibre, WhatsApp), unificándolas en un modelo canónico.

---

## 📌 ¿Qué hace este módulo?

El módulo `orders` es el **núcleo del sistema de gestión de pedidos** de Probability. Centraliza todas las órdenes de venta sin importar su origen, las normaliza a un formato estándar, y gestiona su ciclo de vida completo.

### Problema que resuelve

En una plataforma multi-tenant como Probability, cada negocio:
- Vende por múltiples canales (Shopify, WhatsApp, Amazon, MercadoLibre)
- Cada canal tiene su propio formato de orden
- Necesita ver todas sus órdenes en un solo lugar
- Requiere validación automática de clientes y productos
- Necesita calcular probabilidad de entrega exitosa
- Quiere enviar confirmaciones por WhatsApp

**Este módulo unifica y gestiona todas las órdenes en un solo sistema.**

---

## 🔄 ¿Cómo funciona?

### Flujo Conceptual

```
+----------------------------------------------------------------+
|  1. INTEGRACIÓN GENERA ORDEN                                   |
|  Shopify, WhatsApp, MercadoLibre -> Envía webhook              |
+----------------------------------------------------------------+
                            v
+----------------------------------------------------------------+
|  2. NORMALIZACIÓN A FORMATO CANÓNICO                           |
|  Webhook -> ProbabilityOrderDTO (formato unificado)             |
|  -> Cada integración mapea su formato al formato Probability    |
+----------------------------------------------------------------+
                            v
+----------------------------------------------------------------+
|  3. VALIDACIÓN Y ENRIQUECIMIENTO                               |
|  -> Verificar si orden ya existe (evitar duplicados)            |
|  -> Validar/Crear cliente (por email o DNI)                     |
|  -> Validar/Crear productos (por variant_id, mappings externos, SKU o barcode) |
|  -> Mapear estados específicos a estados Probability            |
+----------------------------------------------------------------+
                            v
+----------------------------------------------------------------+
|  4. PERSISTENCIA COMPLETA                                      |
|  -> Guardar orden principal (orders)                            |
|  -> Guardar items (order_items)                                 |
|  -> Guardar direcciones (addresses)                             |
|  -> Guardar pagos (payments)                                    |
|  -> Guardar envíos (shipments)                                  |
|  -> Guardar datos crudos originales (order_channel_metadata)    |
+----------------------------------------------------------------+
                            v
+----------------------------------------------------------------+
|  5. EVENTOS Y PROCESAMIENTO ASÍNCRONO                          |
|  -> Publicar evento al fanout RabbitMQ (orders.events)          |
|  -> 5 consumers: invoicing, whatsapp, score, inventory, events  |
|  -> Events consumer -> SSE, email, WhatsApp via EventDispatcher  |
|  -> Score consumer calcula probabilidad de entrega              |
+----------------------------------------------------------------+
                            v
+----------------------------------------------------------------+
|  6. CONFIRMACIÓN OPCIONAL POR WHATSAPP                         |
|  -> Enviar mensaje de confirmación al cliente                   |
|  -> Recibir respuesta SÍ/NO                                     |
|  -> Actualizar estado de confirmación                           |
+----------------------------------------------------------------+
```

---

## 🏗️ Responsabilidades del Módulo

### 1. Recepción y Normalización de Órdenes

**¿Qué hace?**
- Recibe órdenes de múltiples integraciones en formato canónico (`ProbabilityOrderDTO`)
- Normaliza datos de diferentes plataformas a un modelo unificado
- Valida estructura y campos requeridos

**Formato canónico**: `ProbabilityOrderDTO`
```go
type ProbabilityOrderDTO struct {
    // Identificadores
    ExternalID      string  // ID de la plataforma origen (ej: shopify_12345)
    OrderNumber     string  // Número de orden visible al cliente
    IntegrationID   uint    // ID de la integración en Probability
    IntegrationType string  // "shopify", "whatsapp", "meli"
    Platform        string  // "shopify", "whatsapp", "mercadolibre"

    // Cliente
    CustomerName  string
    CustomerEmail string
    CustomerPhone string
    CustomerDNI   string

    // Financiero
    Subtotal     float64
    Tax          float64
    Discount     float64
    ShippingCost float64
    TotalAmount  float64
    Currency     string

    // Items
    Items []OrderItemDTO

    // Direcciones
    ShippingAddress  AddressDTO
    BillingAddress   AddressDTO

    // Pago
    Payments []PaymentDTO

    // Envío
    Shipments []ShipmentDTO

    // Estados específicos de la plataforma
    FinancialStatus   string  // "paid", "pending", "refunded"
    FulfillmentStatus string  // "fulfilled", "unfulfilled", "partial"
    Status            string  // "open", "closed", "cancelled"

    // Metadata cruda
    RawData map[string]interface{}
}
```

---

### 2. Validación de Catálogo y Clientes

**Validación de Clientes:**
- Busca cliente existente por **email** o **DNI**
- Si no existe -> Crea nuevo cliente automáticamente
- Asigna `customer_id` a la orden

**Validación de Productos:**
- Para cada item de la orden:
  - Busca producto por `external_variant_id`
  - Si no existe, intenta por `external_sku`
  - Si no existe, intenta por `external_barcode`
  - Si no existe, intenta por `SKU` interno
  - Si no existe, intenta por `barcode` interno
  - Si no existe, intenta por `external_product_id`
  - Si aun no existe -> Crea nuevo producto automáticamente y registra el mapping externo
  - Asigna `product_id` al item

**Beneficio**: Las ordenes nunca fallan por falta de cliente/producto y, ademas, evitan duplicar variantes cuando el canal envia un `variant_id` distinto al `SKU`.

---

### 3. Mapeo de Estados

Cada plataforma tiene sus propios estados. Probability normaliza a tres tipos de estados:

#### Estados de Orden (OrderStatus)
- `pending` -> Pendiente de procesamiento
- `processing` -> En preparación
- `shipped` -> Enviado
- `delivered` -> Entregado
- `completed` -> Completado
- `cancelled` -> Cancelado
- `refunded` -> Reembolsado
- `failed` -> Fallido
- `on_hold` -> En espera

#### Estados de Pago (PaymentStatus)
- `pending` -> Pendiente de pago
- `paid` -> Pagado
- `partially_paid` -> Pago parcial
- `refunded` -> Reembolsado
- `partially_refunded` -> Reembolso parcial
- `voided` -> Anulado
- `authorized` -> Autorizado
- `expired` -> Expirado

#### Estados de Fulfillment (FulfillmentStatus)
- `unfulfilled` -> Sin preparar
- `partial` -> Parcialmente preparado
- `fulfilled` -> Completamente preparado
- `restocked` -> Devuelto a inventario
- `on_hold` -> En espera

**Mapeo de Shopify -> Probability:**
```go
// Ejemplo
Shopify "pending" -> Probability pending (OrderStatus)
Shopify "paid" -> Probability paid (PaymentStatus)
Shopify "fulfilled" -> Probability fulfilled (FulfillmentStatus)
```

---

### 4. Persistencia Multi-Tabla

Una sola orden se guarda en **múltiples tablas relacionadas**:

#### Tabla Principal: `orders`
- Información general de la orden
- Referencias a cliente, integración, business
- Campos financieros agregados
- Estados mapeados
- Score de entrega

#### Tablas Relacionadas:
- **`order_items`** -> Productos de la orden (N items)
- **`addresses`** -> Direcciones de envío/facturación (1-2 direcciones)
- **`payments`** -> Pagos asociados (1-N pagos)
- **`shipments`** -> Envíos y tracking (1-N envíos)
- **`order_channel_metadata`** -> Datos crudos originales (webhook completo)

**Ventajas**:
- Estructura normalizada (evita duplicación)
- Fácil consultar items, pagos, envíos por separado
- Se preserva el webhook original en `order_channel_metadata`

---

### 5. Cálculo de Score de Entrega

**¿Qué es el Score?**
El **Delivery Probability Score** (0-100) estima la probabilidad de que una orden sea entregada exitosamente.

**Factores que afectan el score:**

| Factor | Impacto | Puntos |
|--------|---------|--------|
| Cliente nuevo (sin historial) | Negativo | -10 |
| Pago contra entrega (COD) | Negativo | -10 |
| Dirección incompleta | Negativo | -15 |
| Teléfono inválido | Negativo | -10 |
| Monto muy alto (>$500k) | Negativo | -5 |
| Cliente con historial exitoso | Positivo | +20 |
| Pago anticipado (prepaid) | Positivo | +10 |
| Dirección completa con coordenadas | Positivo | +10 |

**Cálculo**:
```
Score Base = 50

Score Final = Score Base
            + Puntos Historial Cliente
            + Puntos Método de Pago
            - Puntos Factores Negativos

Min: 0, Max: 100
```

**Ejemplo**:
```json
{
  "order_id": "123e4567-e89b-12d3-a456-426614174000",
  "delivery_probability": 65,
  "negative_factors": [
    "Pago contra entrega (COD)",
    "Cliente nuevo"
  ]
}
```

**¿Cuándo se calcula?**
- Asíncronamente después de crear/actualizar la orden
- Se publica evento al fanout RabbitMQ (`orders.events`)
- Score consumer (`orders.events.score`) procesa el evento y actualiza la orden

---

### 6. Sistema de Eventos (RabbitMQ Fanout)

El módulo Orders publica eventos a un **único exchange fanout** (`orders.events`). Cada consumer tiene su propia cola bindeada al fanout, recibiendo una copia de cada evento.

**Exchange:** `orders.events` (tipo fanout)

**Colas bindeadas (5):**

| Cola | Consumer | Propósito |
|------|----------|-----------|
| `orders.events.invoicing` | Invoicing | Facturación automática |
| `orders.events.whatsapp` | WhatsApp | Notificaciones WhatsApp |
| `orders.events.score` | Score | Cálculo de probabilidad de entrega |
| `orders.events.inventory` | Inventory | Actualización de inventario |
| `orders.events.events` | EventDispatcher | SSE, email, WhatsApp via eventos unificados |

**Formato del mensaje:** `OrderEventMessage`
```json
{
  "event_id": "uuid",
  "event_type": "order.created",
  "order_id": "uuid",
  "business_id": 1,
  "integration_id": 5,
  "timestamp": "2026-03-02T10:30:00Z",
  "order": {
    "id": "uuid",
    "order_number": "#1234",
    "internal_number": "ORD-2026-0001",
    "external_id": "shopify_12345",
    "total_amount": 150000,
    "currency": "COP",
    "customer_name": "Juan Pérez",
    "customer_email": "juan@example.com",
    "customer_phone": "+573001234567",
    "platform": "shopify",
    "integration_id": 5
  },
  "changes": {
    "current_status": "processing",
    "previous_status": "pending"
  },
  "metadata": {}
}
```

**Tipos de eventos:**
- `order.created` — Orden creada
- `order.updated` — Orden actualizada
- `order.cancelled` — Orden cancelada
- `order.status_changed` — Cambio de estado

---

### 7. Confirmación por WhatsApp

**Flujo**:

1. **Solicitar confirmación** (POST `/orders/:id/request-confirmation`)
   ```json
   {
     "order_id": "uuid"
   }
   ```

2. **Sistema publica mensaje a RabbitMQ**:
   ```json
   {
     "phone": "+573001234567",
     "message": "Hola Juan, confirma tu orden #1234 por $150.000. Responde SÍ para confirmar o NO para cancelar.",
     "order_id": "uuid",
     "type": "confirmation_request"
   }
   ```

3. **Módulo WhatsApp envía mensaje**

4. **Cliente responde "SÍ" o "NO"**

5. **WhatsApp publica respuesta a RabbitMQ**:
   ```json
   {
     "order_id": "uuid",
     "confirmation_status": "yes",
     "confirmed_at": "2026-01-31T10:35:00Z"
   }
   ```

6. **Consumer actualiza orden**:
   ```sql
   UPDATE orders
   SET is_confirmed = true,
       confirmation_status = 'yes'
   WHERE id = 'uuid';
   ```

---

## 📋 Entidades Principales

### 1. ProbabilityOrder (Orden Principal)

**Tabla**: `orders`

**Campos clave**:

```go
type ProbabilityOrder struct {
    // Identificadores
    ID             string    // UUID
    ExternalID     string    // ID de la plataforma (shopify_12345)
    OrderNumber    string    // Número visible (#1234)
    InternalNumber string    // Número interno (ORD-2026-0001)

    // Relaciones
    BusinessID         *uint  // Negocio dueño
    IntegrationID      uint   // Integración origen
    IntegrationType    string // "shopify", "whatsapp"
    Platform           string // Nombre de la plataforma
    CustomerID         *uint  // Cliente (FK)

    // Financiero
    Subtotal     float64
    Tax          float64
    Discount     float64
    ShippingCost float64
    TotalAmount  float64
    Currency     string
    CodTotal     *float64  // Monto contra entrega

    // Cliente (datos denormalizados)
    CustomerName  string
    CustomerEmail string
    CustomerPhone string
    CustomerDNI   string

    // Dirección de envío (denormalizada)
    ShippingStreet  string
    ShippingCity    string
    ShippingState   string
    ShippingCountry string
    ShippingLat     *float64
    ShippingLng     *float64

    // Pago
    PaymentMethodID *uint
    IsPaid          bool
    PaidAt          *time.Time

    // Logística
    TrackingNumber string
    GuideID        *uint
    WarehouseID    *uint
    DriverID       *uint

    // Estados
    Status              string  // Estado general
    StatusID            *uint   // FK a order_statuses
    PaymentStatusID     *uint   // FK a payment_statuses
    FulfillmentStatusID *uint   // FK a fulfillment_statuses

    // Score de entrega
    DeliveryProbability *float64  // 0-100
    NegativeFactors     []string  // Array de factores negativos

    // Confirmación WhatsApp
    IsConfirmed         *bool
    ConfirmationStatus  string
    ConfirmedAt         *time.Time

    // JSONB (datos estructurados)
    Items              []byte  // Items JSON
    Metadata           []byte  // Metadata adicional
    FinancialDetails   []byte  // Detalles financieros
    ShippingDetails    []byte  // Detalles de envío
    PaymentDetails     []byte  // Detalles de pago
    FulfillmentDetails []byte  // Detalles de fulfillment

    // Auditoría
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

---

### 2. ProbabilityOrderItem (Item de Orden)

**Tabla**: `order_items`

**Relación**: N:1 con `orders`

```go
type ProbabilityOrderItem struct {
    ID      uint
    OrderID string  // FK a orders

    // Producto
    ProductID   *uint   // FK a products
    ProductSKU  string
    ProductName string
    VariantID   string  // ID de variante en el canal, si existe

    // Cantidades y precios
    Quantity   int
    UnitPrice  float64
    TotalPrice float64
    Discount   float64
    Tax        float64
    TaxRate    float64

    // Información adicional
    ImageURL  string
    ProductURL string
    Weight    float64
    RequiresShipping bool

    // JSONB
    Properties []byte  // Propiedades adicionales
}
```

---

### 3. ProbabilityAddress (Dirección)

**Tabla**: `addresses`

**Relación**: N:1 con `orders`

```go
type ProbabilityAddress struct {
    ID      uint
    OrderID string  // FK a orders

    // Tipo
    Type string  // "shipping" o "billing"

    // Datos de contacto
    FirstName string
    LastName  string
    Company   string
    Phone     string

    // Ubicación
    Street     string
    Street2    string
    City       string
    State      string
    Country    string
    PostalCode string
    Latitude   *float64
    Longitude  *float64

    // Instrucciones
    Instructions string
}
```

---

### 4. ProbabilityPayment (Pago)

**Tabla**: `payments`

**Relación**: N:1 con `orders`

```go
type ProbabilityPayment struct {
    ID      uint
    OrderID string  // FK a orders

    // Método de pago
    PaymentMethodID *uint  // FK a payment_methods

    // Montos
    Amount       float64
    Currency     string
    ExchangeRate float64

    // Estado
    Status      string
    PaidAt      *time.Time
    ProcessedAt *time.Time

    // Transacción
    TransactionID    string
    PaymentReference string
    Gateway          string

    // Reembolsos
    RefundAmount float64
    RefundedAt   *time.Time
    FailureReason string
}
```

---

### 5. ProbabilityShipment (Envío)

**Tabla**: `shipments`

**Relación**: N:1 con `orders`

```go
type ProbabilityShipment struct {
    ID      uint
    OrderID string  // FK a orders

    // Tracking
    TrackingNumber string
    TrackingURL    string

    // Transportadora
    Carrier     string
    CarrierCode string

    // Guía
    GuideID  *uint   // FK a guides
    GuideURL string

    // Estado
    Status      string
    ShippedAt   *time.Time
    DeliveredAt *time.Time

    // Logística
    WarehouseID *uint  // FK a warehouses
    DriverID    *uint  // FK a drivers
    IsLastMile  bool

    // Dimensiones
    Weight float64
    Height float64
    Width  float64
    Length float64
}
```

---

### 6. ProbabilityOrderChannelMetadata (Datos Crudos)

**Tabla**: `order_channel_metadata`

**Relación**: N:1 con `orders`

**Propósito**: Preservar el webhook original de la plataforma

```go
type ProbabilityOrderChannelMetadata struct {
    ID      uint
    OrderID string  // FK a orders

    // Canal
    ChannelSource string  // "shopify", "whatsapp", "meli"

    // Datos crudos
    RawData []byte  // JSON completo del webhook

    // Versionado
    Version    string
    ReceivedAt time.Time
    ProcessedAt *time.Time

    // Sincronización
    IsLatest     bool
    LastSyncedAt *time.Time
    SyncStatus   string
}
```

**Beneficio**: Si necesitas los datos originales (ej: campos custom de Shopify), están aquí.

---

## 🚀 API Endpoints

### CRUD de Órdenes

#### Crear Orden Manual
```http
POST /api/orders

Body:
{
  "business_id": 1,
  "integration_id": 5,
  "integration_type": "manual",
  "platform": "admin",
  "order_number": "ORD-001",
  "customer_name": "Juan Pérez",
  "customer_email": "juan@example.com",
  "customer_phone": "+573001234567",
  "total_amount": 150000,
  "currency": "COP",
  "items": [
    {
      "product_sku": "ABC-123",
      "product_name": "Producto 1",
      "quantity": 2,
      "unit_price": 75000
    }
  ]
}

Response: 201 Created
{
  "id": "uuid",
  "order_number": "ORD-001",
  "status": "pending",
  "total_amount": 150000,
  "created_at": "2026-01-31T10:30:00Z"
}
```

#### Mapear y Guardar Orden (Desde Integración)
```http
POST /api/orders/map

Body: ProbabilityOrderDTO (formato canónico)
{
  "external_id": "shopify_12345",
  "integration_id": 5,
  "integration_type": "shopify",
  "platform": "shopify",
  "order_number": "#1234",
  "customer_name": "Juan Pérez",
  "customer_email": "juan@example.com",
  "total_amount": 150000,
  "currency": "COP",
  "items": [...],
  "shipping_address": {...},
  "payments": [...],
  "raw_data": {...}  # Webhook completo
}

Response: 201 Created / 200 OK (si ya existía)
{
  "id": "uuid",
  "external_id": "shopify_12345",
  "order_number": "#1234",
  "status": "pending",
  "delivery_probability": 65,
  "negative_factors": ["Pago contra entrega (COD)"],
  "created_at": "2026-01-31T10:30:00Z"
}
```

#### Listar Órdenes
```http
GET /api/orders?page=1&page_size=20&business_id=1&status=pending

Query Parameters:
- page: int (default: 1)
- page_size: int (default: 20, max: 100)
- business_id: uint
- integration_id: uint
- status_id: uint
- payment_status_id: uint
- fulfillment_status_id: uint
- customer_email: string
- customer_phone: string
- order_number: string
- internal_number: string
- is_paid: bool
- is_cod: bool
- warehouse_id: uint
- driver_id: uint
- start_date: datetime (YYYY-MM-DD)
- end_date: datetime (YYYY-MM-DD)
- sort_by: string (default: "created_at")
- sort_order: string ("asc" / "desc", default: "desc")

Response: 200 OK
{
  "data": [
    {
      "id": "uuid",
      "order_number": "#1234",
      "customer_name": "Juan Pérez",
      "total_amount": 150000,
      "currency": "COP",
      "status": "pending",
      "created_at": "2026-01-31T10:30:00Z"
    }
  ],
  "total": 150,
  "page": 1,
  "page_size": 20,
  "total_pages": 8
}
```

#### Obtener Orden por ID
```http
GET /api/orders/:id

Response: 200 OK
{
  "id": "uuid",
  "external_id": "shopify_12345",
  "order_number": "#1234",
  "internal_number": "ORD-2026-0001",
  "business_id": 1,
  "integration_id": 5,
  "integration_type": "shopify",
  "platform": "shopify",
  "customer": {
    "id": 100,
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "phone": "+573001234567"
  },
  "subtotal": 130000,
  "tax": 20000,
  "shipping_cost": 0,
  "discount": 0,
  "total_amount": 150000,
  "currency": "COP",
  "status": "pending",
  "payment_status": "pending",
  "fulfillment_status": "unfulfilled",
  "delivery_probability": 65,
  "negative_factors": ["Pago contra entrega (COD)"],
  "items": [...],
  "addresses": [...],
  "payments": [...],
  "shipments": [...],
  "created_at": "2026-01-31T10:30:00Z",
  "updated_at": "2026-01-31T10:30:00Z"
}
```

#### Actualizar Orden
```http
PATCH /api/orders/:id

Body:
{
  "status": "processing",
  "tracking_number": "ABC123456",
  "delivery_probability": 80
}

Response: 200 OK
{
  "id": "uuid",
  "order_number": "#1234",
  "status": "processing",
  "tracking_number": "ABC123456",
  "delivery_probability": 80,
  "updated_at": "2026-01-31T11:00:00Z"
}
```

#### Eliminar Orden
```http
DELETE /api/orders/:id

Response: 204 No Content
```

---

### Endpoints Especiales

#### Obtener Datos Crudos de la Orden
```http
GET /api/orders/:id/raw

Response: 200 OK
{
  "order_id": "uuid",
  "channel_source": "shopify",
  "raw_data": {
    "id": 12345,
    "email": "juan@example.com",
    "created_at": "2026-01-31T10:25:00Z",
    "line_items": [...],
    "shipping_address": {...},
    "customer": {...},
    # ... Webhook completo de Shopify
  },
  "received_at": "2026-01-31T10:25:05Z",
  "processed_at": "2026-01-31T10:25:10Z"
}
```

#### Solicitar Confirmación por WhatsApp
```http
POST /api/orders/:id/request-confirmation

Response: 200 OK
{
  "message": "Confirmation request sent successfully",
  "order_id": "uuid",
  "phone": "+573001234567",
  "sent_at": "2026-01-31T10:35:00Z"
}
```

---

## 📊 Casos de Uso Principales

### Caso 1: Crear Orden desde Shopify

**Flujo completo**:

1. **Shopify envía webhook** cuando se crea una orden
2. **Módulo Shopify recibe webhook** y lo convierte a `ProbabilityOrderDTO`
3. **Módulo Shopify llama** `POST /orders/map` con el DTO canónico
4. **Módulo Orders**:
   - Verifica si orden ya existe (`external_id` + `integration_id`)
   - Si existe -> Actualiza
   - Si no existe -> Crea nueva
5. **Valida/Crea Cliente**:
   - Busca por email: `juan@example.com`
   - Si no existe -> Crea cliente nuevo
6. **Valida/Crea Productos**:
   - Para cada item, busca primero por `variant_id`
   - Si no existe, intenta por mappings externos, `SKU` y `barcode`
   - Si no existe -> Crea producto nuevo y registra el mapping externo
7. **Mapea Estados**:
   - `financial_status: "paid"` -> `payment_status_id: 2` (paid)
   - `fulfillment_status: "unfulfilled"` -> `fulfillment_status_id: 1` (unfulfilled)
8. **Guarda en BD**:
   - `orders` -> Orden principal
   - `order_items` -> 3 items
   - `addresses` -> 1 dirección de envío
   - `payments` -> 1 pago
   - `order_channel_metadata` -> Webhook original
9. **Publica Eventos**:
   - `order.created` -> RabbitMQ fanout (`orders.events`)
   - Todos los consumers reciben copia: invoicing, whatsapp, score, inventory, events
10. **Consumer calcula score** (asíncrono):
    - Cliente nuevo -> -10 pts
    - Pago anticipado -> +10 pts
    - Score final: 60
11. **Actualiza orden** con score y factores negativos

---

### Caso 2: Calcular Score de Entrega

**Actor**: RabbitMQ Score Consumer (`orders.events.score`)
**Trigger**: Evento `order.created` o `order.updated` en fanout

**Flujo**:

1. **Consumer recibe evento** con `order_id`
2. **Obtiene orden completa** con relaciones
3. **Obtiene historial del cliente**:
   ```sql
   SELECT COUNT(*) as total_orders,
          SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END) as delivered
   FROM orders
   WHERE customer_id = ?
   ```
4. **Calcula score base**: 50 puntos
5. **Suma/resta factores**:
   - Cliente nuevo (0 órdenes previas): -10
   - Pago contra entrega (`is_cod = true`): -10
   - Dirección incompleta (sin `shipping_lat`): -15
   - Teléfono inválido (menos de 10 dígitos): -10
   - Monto alto (>$500k): -5
   - Cliente con historial exitoso (>80% entregadas): +20
   - Pago anticipado (`is_paid = true`): +10
6. **Score final**: `max(0, min(100, score_base + factores))`
7. **Identifica factores negativos**:
   ```json
   [
     "Pago contra entrega (COD)",
     "Dirección incompleta"
   ]
   ```
8. **Actualiza orden**:
   ```sql
   UPDATE orders
   SET delivery_probability = 55,
       negative_factors = '["Pago contra entrega (COD)", "Dirección incompleta"]'
   WHERE id = 'uuid';
   ```
9. **Publica evento** `order.updated`

---

### Caso 3: Listar Órdenes con Filtros

**Usuario**: Admin del negocio
**Objetivo**: Ver órdenes pendientes de su negocio

**Request**:
```http
GET /api/orders?business_id=1&status_id=1&is_cod=true&page=1&page_size=20
```

**Flujo**:

1. **Handler parsea query params**
2. **Construye filtros**:
   ```go
   filters := map[string]interface{}{
       "business_id": 1,
       "status_id": 1,
       "is_cod": true,
   }
   ```
3. **Repository construye query**:
   ```sql
   SELECT o.*,
          c.name as customer_name,
          os.name as status_name
   FROM orders o
   LEFT JOIN clients c ON o.customer_id = c.id
   LEFT JOIN order_statuses os ON o.status_id = os.id
   WHERE o.business_id = 1
     AND o.status_id = 1
     AND o.is_cod = true
     AND o.deleted_at IS NULL
   ORDER BY o.created_at DESC
   LIMIT 20 OFFSET 0;
   ```
4. **Cuenta total** (para paginación):
   ```sql
   SELECT COUNT(*) FROM orders
   WHERE business_id = 1
     AND status_id = 1
     AND is_cod = true
     AND deleted_at IS NULL;
   ```
5. **Mapea a DTOs** (`OrderSummary`)
6. **Retorna respuesta paginada**:
   ```json
   {
     "data": [...],
     "total": 45,
     "page": 1,
     "page_size": 20,
     "total_pages": 3
   }
   ```

---

## 🏛️ Arquitectura Técnica

### Estructura de Carpetas

```
orders/
+-- bundle.go                    # Ensamblador del módulo
+-- internal/
    +-- domain/                  # 🔵 CAPA DE DOMINIO
    |   +-- entities/            # Entidades principales
    |   +-- dtos/                # DTOs de dominio
    |   +-- ports/               # Interfaces (contratos)
    |   +-- errors/              # Errores de dominio
    |
    +-- app/                     # 🟢 CAPA DE APLICACIÓN
    |   +-- helpers/             # Helpers compartidos
    |   |   +-- statusmapper/    # Mapeo de estados por plataforma
    |   +-- usecaseorder/        # Casos de uso CRUD (get, list, delete)
    |   +-- usecasecreateorder/  # Crear orden + publicar eventos
    |   +-- usecaseupdateorder/  # Actualizar orden + publicar eventos
    |   +-- usecaseorderscore/   # Cálculo de score de entrega
    |
    +-- mocks/                   # Mocks para testing
    |
    +-- infra/                   # 🔴 CAPA DE INFRAESTRUCTURA
        +-- primary/             # Adaptadores de entrada
        |   +-- handlers/        # HTTP handlers (Gin)
        |   |   +-- constructor.go
        |   |   +-- router.go
        |   |   +-- request/     # DTOs HTTP de entrada (CON tags)
        |   |   +-- response/    # DTOs HTTP de salida (CON tags)
        |   |   +-- mappers/     # Conversión Domain ↔ HTTP
        |   |   +-- *.go         # Un handler por archivo
        |   +-- queue/           # RabbitMQ consumers
        |       +-- consumer.go          # Score consumer
        |       +-- whatsapp_consumer.go # WhatsApp response consumer
        |
        +-- secondary/           # Adaptadores de salida
            +-- repository/      # PostgreSQL (GORM)
            |   +-- constructor.go
            |   +-- repository.go
            |   +-- status_queries.go  # Consultas réplica de estados
            |   +-- mappers/
            +-- eventpublisher/  # Publisher de integración (Shopify events)
            |   +-- constructor.go
            +-- queue/           # RabbitMQ fanout publisher
                +-- order_publisher.go   # Publica al fanout orders.events
                +-- response/            # Structs de mensaje
                +-- mappers/             # Mapeo Order -> OrderEventMessage
```

---

### ✅ Estado de Arquitectura Hexagonal

**Cumplimiento**: ✅ **CONFORME 100%** (refactorizado 2026-01-31)

#### ✅ Aspectos Conformes

**Separación de Capas:**
- ✅ Separación de capas (Domain, App, Infra)
- ✅ Flujo de dependencias correcto (Infra -> App -> Domain)
- ✅ Domain 100% puro (SIN tags, SIN GORM, SIN Gin)
- ✅ Domain organizado en subcarpetas (`entities/`, `dtos/`, `ports/`, `errors/`)

**Domain Layer (Pureza Absoluta):**
- ✅ Entidades en `domain/entities/` sin tags (ni json, ni gorm, ni validate)
- ✅ DTOs en `domain/dtos/` sin tags (ni json, ni binding)
- ✅ Domain usa `[]byte` para datos JSON (NO usa `datatypes.JSON`)
- ✅ Ports en `domain/ports/` con interfaces puras
- ✅ Errors en `domain/errors/` con errores de dominio
- ✅ Domain NO importa GORM
- ✅ Domain NO importa Gin

**Infrastructure Layer (HTTP):**
- ✅ Handlers tienen DTOs HTTP separados en `request/` (CON tags)
- ✅ Handlers tienen DTOs HTTP separados en `response/` (CON tags)
- ✅ Handlers tienen mappers en `mappers/` (conversión Domain ↔ HTTP)
- ✅ DTOs HTTP usan `datatypes.JSON` (permitido en infra)
- ✅ Mappers convierten `datatypes.JSON` ↔ `[]byte`

**Otros Aspectos:**
- ✅ Repositorio usa modelos GORM correctamente (NO usa `.Table()`)
- ✅ Constructor único en `bundle.go`
- ✅ Handlers tienen `router.go` con `RegisterRoutes()`
- ✅ Un método por archivo en handlers

#### 📁 Estructura Refactorizada

```
orders/internal/
+-- domain/                      # 🔵 CAPA DE DOMINIO (PURA)
|   +-- entities/               # ✅ 11 archivos SIN tags
|   |   +-- order.go
|   |   +-- order_item.go
|   |   +-- address.go
|   |   +-- payment.go
|   |   +-- shipment.go
|   |   +-- channel_metadata.go
|   |   +-- product.go
|   |   +-- client.go
|   |   +-- order_status.go
|   |   +-- order_event.go
|   |   +-- order_error.go
|   +-- dtos/                   # ✅ 5 archivos SIN tags, usa []byte
|   |   +-- create_order_request.go
|   |   +-- update_order_request.go
|   |   +-- order_response.go
|   |   +-- order_summary.go
|   |   +-- probability_order_dto.go
|   +-- ports/                  # ✅ Interfaces
|   |   +-- ports.go
|   +-- errors/                 # ✅ Errores de dominio
|       +-- errors.go
|
+-- infra/primary/handlers/     # 🔴 CAPA DE INFRAESTRUCTURA
    +-- request/                # ✅ 4 archivos CON tags + datatypes.JSON
    |   +-- create_order.go
    |   +-- update_order.go
    |   +-- map_order.go
    |   +-- list_orders_filters.go
    +-- response/               # ✅ 5 archivos CON tags + datatypes.JSON
    |   +-- order.go
    |   +-- order_summary.go
    |   +-- order_raw.go
    |   +-- orders_list.go
    |   +-- error.go
    +-- mappers/                # ✅ 2 archivos de conversión
    |   +-- to_domain.go        # HTTP -> Domain (datatypes.JSON -> []byte)
    |   +-- to_response.go      # Domain -> HTTP ([]byte -> datatypes.JSON)
    +-- *.go                    # Handlers actualizados

```

#### 🔄 Flujo de Datos (Domain ↔ HTTP)

```
HTTP Request (CON tags + datatypes.JSON)
    v
[mappers.MapOrderRequestToDomain()]
    v
Domain DTO (SIN tags + []byte)
    v
[UseCase procesa]
    v
Domain Response (SIN tags + []byte)
    v
[mappers.OrderToResponse()]
    v
HTTP Response (CON tags + datatypes.JSON)
```

#### ✅ Verificaciones Pasadas

```bash
# ✅ Sin tags en domain/entities
grep -r 'json:\|gorm:\|validate:' internal/domain/entities/
# Resultado: vacío

# ✅ Sin tags en domain/dtos
grep -r 'json:\|gorm:\|validate:' internal/domain/dtos/
# Resultado: vacío

# ✅ Sin imports de GORM en domain
grep -r "^import\|^\t\"" internal/domain/ | grep gorm
# Resultado: vacío

# ✅ Sin imports de Gin en domain
grep -r "^import\|^\t\"" internal/domain/ | grep -E "gin|fiber|echo"
# Resultado: vacío

# ✅ Compilación exitosa
go build ./services/modules/orders/...
# Resultado: OK (sin errores en handlers)
```

---

## 🧪 Testing

### Estado Actual
- ✅ Tests unitarios implementados para casos de uso principales

### Tests Implementados

#### `usecasecreateorder` (app/usecasecreateorder/)
- Creación exitosa de orden
- Creación con cliente existente
- Manejo de errores de repositorio
- Publicación de eventos al fanout

#### `usecaseupdateorder` (app/usecaseupdateorder/)
- Actualización exitosa de campos
- Cambio de estado con publicación de eventos
- Manejo de orden no encontrada
- Publicación de eventos al fanout

#### `usecaseorder` (app/usecaseorder/)
- Get orden por ID
- List órdenes con filtros y paginación

#### `usecaseorderscore` (app/usecaseorderscore/)
- Cálculo de score de entrega

#### `statusmapper` (app/helpers/statusmapper/)
- Mapeo de estados por tipo de integración

### Ejecutar Tests

```bash
go test ./services/modules/orders/...
```

---

## 🛠️ Desarrollo

### Compilar

```bash
cd /back/central
go build ./services/modules/orders/...
```

### Ejecutar Tests

```bash
go test ./services/modules/orders/...
```

### Verificar Arquitectura

```bash
# Verificar tags en domain
grep -r 'json:"\|gorm:"' services/modules/orders/internal/domain/

# Verificar imports prohibidos
grep -r "gorm\|gin" services/modules/orders/internal/domain/

# Verificar uso de .Table()
grep -r '\.Table(' services/modules/orders/internal/infra/secondary/repository/
```

---

## 📝 Convenciones

1. **Formato canónico**: Todas las integraciones deben convertir sus webhooks a `ProbabilityOrderDTO`
2. **Idempotencia**: Usar `external_id` + `integration_id` para evitar duplicados
3. **Validación automática**: Crear clientes/productos si no existen
4. **Eventos asíncronos**: Cálculo de score y notificaciones se procesan en background
5. **Soft deletes**: Usar `deleted_at` en lugar de eliminar físicamente
6. **Preservar datos crudos**: Siempre guardar webhook original en `order_channel_metadata`

---

## 📦 Dependencias

- **GORM**: ORM para PostgreSQL
- **Gin**: Framework HTTP
- **RabbitMQ**: Fanout de eventos + colas de integración
- **Zerolog**: Logging estructurado
- **UUID**: Generación de IDs únicos

---

## 📡 Sistema de Publicaciones

El módulo Orders publica eventos a un **único exchange fanout** de RabbitMQ. Cada consumer recibe una copia de cada evento.

### RabbitMQ Fanout
- **Exchange:** `orders.events` (tipo fanout)
- **Garantía:** At-least-once delivery
- **Publisher:** `IOrderRabbitPublisher` (infra/secondary/queue/)

### Colas Bindeadas (5)

| Cola | Consumer | Módulo | Propósito |
|------|----------|--------|-----------|
| `orders.events.invoicing` | Invoicing | invoicing | Facturación automática |
| `orders.events.whatsapp` | WhatsApp | whatsapp | Notificaciones WhatsApp |
| `orders.events.score` | Score | orders (interno) | Cálculo de probabilidad |
| `orders.events.inventory` | Inventory | inventory | Actualización de inventario |
| `orders.events.events` | EventDispatcher | events | SSE, email, WhatsApp unificado |

### Tipos de Eventos Publicados

| Evento | Descripción |
|--------|-------------|
| `order.created` | Orden creada |
| `order.updated` | Orden actualizada |
| `order.cancelled` | Orden cancelada |
| `order.status_changed` | Cambio de estado |

Todos los eventos se publican al fanout con formato `OrderEventMessage` que incluye un `OrderSnapshot` completo y un mapa de `Changes`.

---

## 🔗 Integraciones

### Módulos que consumen Orders

- **Dashboard**: Visualización de órdenes
- **Analytics**: Reportes y métricas
- **Invoicing**: Generación de facturas
- **Notifications**: Envío de notificaciones
- **Shipments**: Gestión de envíos

### Módulos que crean Orders

- **Shopify Integration**: Órdenes de Shopify
- **WhatsApp Integration**: Órdenes por WhatsApp
- **MercadoLibre Integration**: Órdenes de MercadoLibre
- **Amazon Integration**: Órdenes de Amazon

---

## 📜 Changelog

### v2.0.0 (2026-03-02) - Consolidación de Eventos: Fanout Único

**Eliminado:**
- `IEventsExchangePublisher` — ya no se publica al topic exchange `events.exchange`
- Publicación dual (Redis + topic exchange) eliminada completamente
- `events_exchange_publisher.go` eliminado

**Nuevo:**
- Cola `orders.events.events` bindeada al fanout existente
- Módulo `events` consume del fanout y transforma `OrderEventMessage` -> `entities.Event`
- EventDispatcher (SSE, email, WhatsApp) ahora recibe eventos via fanout, no via topic exchange

**Refactorizado:**
- Use cases `usecasecreateorder` y `usecaseupdateorder` solo publican al fanout
- WhatsApp consumer usa `IOrderRabbitPublisher` en vez de `IEventsExchangePublisher`
- Tests unitarios actualizados (mock de events publisher eliminado)

### v1.1.0 (2026-01-31) - Sistema Unificado de Publicaciones

**Features:**
- Sistema de publicaciones RabbitMQ con fanout exchange
- Eventos para todos los tipos: created, updated, cancelled, status_changed
- Estructura organizada con response/mappers en queue publisher

### v1.0.0

**Features:**
- CRUD completo de órdenes
- Mapeo desde formato canónico (`ProbabilityOrderDTO`)
- Validación automática de clientes y productos
- Mapeo de estados específicos de plataformas
- Cálculo de score de entrega
- Confirmación por WhatsApp (RabbitMQ)
- Preservación de datos crudos
- Arquitectura hexagonal (Domain, App, Infra)

---

**Última actualización:** 2026-03-02
