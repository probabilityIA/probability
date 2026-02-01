# RabbitMQ Events - Orders Module

## 📋 Resumen

El módulo Orders publica eventos de órdenes en **RabbitMQ** para garantizar el procesamiento confiable por parte de otros módulos (Invoicing, Notifications, Events).

Estos eventos se publican **simultáneamente en Redis** (para tiempo real/scoring) y **RabbitMQ** (para garantía de entrega).

---

## 🔄 Sistema Dual de Publicaciones

```
ORDER EVENT
    ↓
┌───────────────────────┐
│  Event Publisher      │
└───────────────────────┘
           ↓
    ┌──────┴──────┐
    ↓             ↓
┌─────────┐   ┌──────────┐
│  REDIS  │   │ RABBITMQ │
│ Pub/Sub │   │  Queues  │
└─────────┘   └──────────┘
    ↓             ↓
Scoring       Invoicing, Notifications
```

**Características:**
- **Redis:** Pub/Sub para eventos en tiempo real (scoring, dashboard SSE)
- **RabbitMQ:** Garantía de entrega at-least-once (facturas, notificaciones críticas)
- **Asíncrono:** Publicaciones en goroutines para no bloquear respuesta HTTP
- **Tolerante a fallos:** Redis falla silenciosamente, RabbitMQ registra errores

---

## 📡 Queues Publicadas

### 1. `orders.events.created`

**Evento:** Nueva orden creada

**Consumidores esperados:** Invoicing (facturas automáticas), Notifications (notificación al cliente), Events (webhooks)

**Payload:**
```json
{
  "event_id": "20260131101500-abc123",
  "event_type": "order.created",
  "order_id": "uuid-de-la-orden",
  "business_id": 123,
  "integration_id": 456,
  "timestamp": "2026-01-31T10:15:00Z",
  "order": {
    "id": "uuid-de-la-orden",
    "order_number": "ORD-001",
    "internal_number": "INT-001",
    "external_id": "shopify-12345",
    "total_amount": 499.99,
    "currency": "COP",
    "customer_name": "Juan Pérez",
    "customer_email": "juan@example.com",
    "customer_phone": "+573001234567",
    "platform": "shopify",
    "created_at": "2026-01-31T10:15:00Z",
    "updated_at": "2026-01-31T10:15:00Z"
  }
}
```

**Casos de uso:**
- Generar factura electrónica automáticamente
- Enviar email de confirmación al cliente
- Notificar al vendedor sobre nueva orden
- Disparar webhook a sistemas externos

---

### 2. `orders.events.updated`

**Evento:** Orden actualizada

**Consumidores esperados:** Notifications (notificar cambios), Events (webhooks de actualización)

**Payload:**
```json
{
  "event_id": "20260131102000-def456",
  "event_type": "order.updated",
  "order_id": "uuid-de-la-orden",
  "business_id": 123,
  "integration_id": 456,
  "timestamp": "2026-01-31T10:20:00Z",
  "order": {
    "id": "uuid-de-la-orden",
    "order_number": "ORD-001",
    "total_amount": 549.99,
    "currency": "COP",
    "platform": "shopify",
    "updated_at": "2026-01-31T10:20:00Z"
  }
}
```

**Casos de uso:**
- Notificar al cliente sobre cambios en su orden
- Actualizar sistemas de inventario
- Actualizar dashboard en tiempo real

---

### 3. `orders.events.cancelled`

**Evento:** Orden cancelada

**Consumidores esperados:** Invoicing (generar nota de crédito), Notifications (notificar cancelación)

**Payload:**
```json
{
  "event_id": "20260131103000-ghi789",
  "event_type": "order.cancelled",
  "order_id": "uuid-de-la-orden",
  "business_id": 123,
  "integration_id": 456,
  "timestamp": "2026-01-31T10:30:00Z",
  "order": {
    "id": "uuid-de-la-orden",
    "order_number": "ORD-001",
    "status": "cancelled",
    "platform": "shopify",
    "updated_at": "2026-01-31T10:30:00Z"
  }
}
```

**Casos de uso:**
- Generar nota de crédito automática (si ya se facturó)
- Liberar inventario reservado
- Notificar al cliente sobre cancelación
- Revertir pagos si aplica

---

### 4. `orders.events.status_changed`

**Evento:** Cambio de estado de orden

**Consumidores esperados:** Notifications (notificar cambio de estado), Dashboard (actualización en tiempo real)

**Payload:**
```json
{
  "event_id": "20260131104000-jkl012",
  "event_type": "order.status_changed",
  "order_id": "uuid-de-la-orden",
  "business_id": 123,
  "integration_id": 456,
  "timestamp": "2026-01-31T10:40:00Z",
  "changes": {
    "previous_status": "pending",
    "current_status": "shipped"
  },
  "order": {
    "id": "uuid-de-la-orden",
    "order_number": "ORD-001",
    "status": "shipped",
    "platform": "shopify",
    "updated_at": "2026-01-31T10:40:00Z"
  }
}
```

**Casos de uso:**
- Notificar al cliente que su orden fue enviada
- Actualizar dashboard con estado en tiempo real
- Disparar automáticamente acciones según estado (ej: enviar tracking al pasar a "shipped")

---

### 5. `orders.confirmation.requested`

**Evento:** Solicitud de confirmación WhatsApp (YA EXISTENTE)

**Consumidores esperados:** WhatsApp Service

**Payload:**
```json
{
  "event_type": "order.confirmation_requested",
  "order_id": "uuid-de-la-orden",
  "order_number": "ORD-001",
  "business_id": 123,
  "customer_name": "Juan Pérez",
  "customer_phone": "+573001234567",
  "customer_email": "juan@example.com",
  "total_amount": 499.99,
  "currency": "COP",
  "items_summary": "2x Producto A, 1x Producto B",
  "shipping_address": "Calle 123 #45-67, Bogotá, Cundinamarca",
  "payment_method_id": 1,
  "integration_id": 456,
  "platform": "shopify",
  "timestamp": 1738334100
}
```

**Casos de uso:**
- Enviar mensaje de WhatsApp solicitando confirmación de orden

---

## 🛠️ Consumir Eventos

### Estructura de Consumidor Básico (Go)

```go
package consumer

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type OrderEventMessage struct {
    EventID       string                 `json:"event_id"`
    EventType     string                 `json:"event_type"`
    OrderID       string                 `json:"order_id"`
    BusinessID    *uint                  `json:"business_id"`
    IntegrationID *uint                  `json:"integration_id"`
    Timestamp     time.Time              `json:"timestamp"`
    Order         *OrderSnapshot         `json:"order,omitempty"`
    Changes       map[string]interface{} `json:"changes,omitempty"`
}

type OrderSnapshot struct {
    ID             string    `json:"id"`
    OrderNumber    string    `json:"order_number"`
    TotalAmount    float64   `json:"total_amount"`
    Currency       string    `json:"currency"`
    CustomerEmail  string    `json:"customer_email,omitempty"`
    Platform       string    `json:"platform"`
}

func StartConsumer(rabbit rabbitmq.IQueue, queueName string) {
    rabbit.Consume(context.Background(), queueName, func(ctx context.Context, body []byte) error {
        var event OrderEventMessage
        if err := json.Unmarshal(body, &event); err != nil {
            return fmt.Errorf("error unmarshaling event: %w", err)
        }

        // Procesar según tipo de evento
        switch event.EventType {
        case "order.created":
            return handleOrderCreated(ctx, &event)
        case "order.updated":
            return handleOrderUpdated(ctx, &event)
        case "order.cancelled":
            return handleOrderCancelled(ctx, &event)
        case "order.status_changed":
            return handleStatusChanged(ctx, &event)
        default:
            // Ignorar eventos desconocidos
            return nil
        }
    })
}

func handleOrderCreated(ctx context.Context, event *OrderEventMessage) error {
    // Ejemplo: Generar factura
    fmt.Printf("Generando factura para orden %s\n", event.OrderID)
    return nil
}
```

---

## 📊 Monitoreo

### Logs Estructurados

Todas las publicaciones generan logs estructurados con zerolog:

```json
{
  "level": "info",
  "order_id": "uuid-de-la-orden",
  "event_type": "order.created",
  "queue": "orders.events.created",
  "message": "✅ Order event published to RabbitMQ"
}
```

### Verificar Queue en RabbitMQ UI

1. Acceder a RabbitMQ Management UI: http://localhost:15672
2. Usuario: admin / Contraseña: admin
3. Ir a **Queues** > Buscar queues `orders.events.*`
4. Ver mensajes pendientes, consumidores activos, rate de publicación

---

## 🔍 Troubleshooting

### Evento no llega a consumidor

**Verificar:**
1. ✅ Queue existe en RabbitMQ
2. ✅ Consumidor está conectado (ver "Consumers" en RabbitMQ UI)
3. ✅ Exchange correcto (si usas exchanges en lugar de queues directas)
4. ✅ Logs de publicación muestran éxito

### Errores de deserialización

**Solución:** Asegurarse que el consumidor use la misma estructura `OrderEventMessage` o sea flexible con campos opcionales (`omitempty`).

### Duplicados

**Causa:** RabbitMQ garantiza at-least-once delivery. Los consumidores DEBEN ser idempotentes.

**Solución:** Usar `event_id` como clave de deduplicación en el consumidor.

---

## 📝 Changelog

### 2026-01-31
- ✅ Implementado sistema dual de publicaciones (Redis + RabbitMQ)
- ✅ Agregados eventos: created, updated, cancelled, status_changed
- ✅ Documentación completa de queues y payloads

---

**Última actualización:** 2026-01-31
**Autor:** Orders Module Team
