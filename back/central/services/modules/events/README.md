# modules/events

Centro neurálgico de notificaciones en tiempo real. Consume eventos desde tres canales Redis, los filtra según la configuración de notificaciones de cada negocio, y los entrega a clientes frontend vía Server-Sent Events (SSE).

---

## Qué hace

Actúa como dispatcher centralizado entre los sistemas que generan eventos (órdenes, facturación, integraciones) y los clientes web que deben recibirlos en tiempo real. Ningún módulo publicador conoce este módulo: la comunicación es enteramente a través de Redis Pub/Sub.

---

## Estructura

```
modules/events/
├── bundle.go                                       # Inicialización del módulo
└── internal/
    ├── app/
    │   ├── constructor.go                          # OrderEventConsumer (órdenes internas)
    │   ├── invoice_consumer.go                     # InvoiceEventConsumer (facturación)
    │   └── integration_consumer.go                 # IntegrationEventConsumer (sync de integraciones)
    ├── domain/
    │   ├── ports.go                                # IEventPublisher, INotificationConfigRepository
    │   ├── dtos.go                                 # Event (genérico), EventType, SSEConnectionFilter
    │   ├── order_events.go                         # OrderEvent, OrderEventType, OrderEventData
    │   ├── invoice_events.go                       # InvoiceEvent, InvoiceEventType
    │   ├── integration_events.go                   # IntegrationEvent, IntegrationEventType
    │   ├── notification_config.go                  # NotificationConfig (configuración por negocio)
    │   └── connection.go                           # SSEConnection, SSEConnectionFilter
    └── infra/
        ├── primary/
        │   ├── handlers/
        │   │   └── sse_gin.go                      # HTTP Handler SSE
        │   └── routes.go                           # GET /events/sse/:businessID
        └── secondary/
            ├── redis/
            │   ├── subscriber.go                   # OrderEventSubscriber
            │   ├── invoice_subscriber.go           # InvoiceEventSubscriber
            │   └── integration_subscriber.go       # IntegrationEventSubscriber
            ├── events/                             # EventManager: broadcast SSE + caché
            └── repository/
                └── notification_config_repository.go
```

---

## Canales Redis consumidos

| Canal | Constante | Publicador | Tipos de evento |
|-------|-----------|-----------|-----------------|
| `probability:orders:state:events` | `ChannelOrdersEvents` | `modules/orders` | `order.created`, `order.status_changed`, `order.cancelled`, etc. |
| `probability:invoicing:state:events` | `ChannelInvoicingEvents` | `modules/invoicing` | `invoice.created`, `invoice.failed`, `bulk_job.completed`, etc. |
| `probability:integrations:orders:sync:events` | `ChannelIntegrationsSyncOrders` | `integrations/events` | `integration.sync.order.created`, `integration.sync.completed`, etc. |

---

## Flujo de datos

```
Redis Pub/Sub (3 canales)
        │
        ├── OrderEventSubscriber
        │       └── OrderEventConsumer
        │               ├── shouldNotifyEvent() — filtra por NotificationConfig + estado permitido
        │               └── publishOrderEvent() — convierte a Event genérico
        │
        ├── InvoiceEventSubscriber
        │       └── InvoiceEventConsumer
        │               └── publishInvoiceEvent() — convierte a Event genérico
        │
        └── IntegrationEventSubscriber
                └── IntegrationEventConsumer
                        ├── shouldNotifyEvent() — filtra por NotificationConfig
                        └── publishIntegrationEvent() — convierte a Event genérico
                                │
                                ▼
                        EventManager.PublishEvent()
                                │
                                ▼
                        broadcastToBusinesses()
                          (filtra conexiones SSE activas por business_id)
                                │
                                ▼
                        Cliente HTTP recibe evento SSE
```

Cada subscriber corre en su propia goroutine con un buffer de 100 eventos. El EventManager mantiene caché de los últimos 2000 eventos por negocio para rehidratación al reconectar.

---

## Filtrado con NotificationConfig

Antes de enviar un evento al EventManager, los consumers consultan la tabla `notification_configs` en base de datos:

```
¿El evento tiene business_id?
    NO  → notificar siempre
    SÍ  → GetByBusinessAndEventType(business_id, event_type)
              │
              ├── No existe config → notificar (comportamiento por defecto)
              ├── config.Enabled == false → descartar
              └── config.Enabled == true
                      │
                      └── Para order.status_changed:
                              ¿config.OrderStatusCodes vacío?
                                  SÍ  → notificar todos los estados
                                  NO  → solo si currentStatus está en la lista
```

Esto permite que cada negocio configure exactamente qué eventos quiere recibir sin modificar código.

---

## Endpoint HTTP

```
GET /events/sse/:businessID
```

| Query param | Tipo | Descripción |
|-------------|------|-------------|
| `integration_id` | uint | Filtrar por integración específica |
| `event_types` | string (comas) | Solo recibir estos tipos de evento |

**Formato SSE de respuesta:**
```
event: integration.sync.order.created
data: {"id":"uuid","type":"integration.sync.order.created","business_id":"123","integration_id":13,"timestamp":"...","data":{...}}

event: order.status_changed
data: {"id":"uuid","type":"order.status_changed","business_id":"123","data":{"current_status":"fulfilled",...}}
```

Al conectar, el EventManager rehidrata automáticamente los eventos recientes que el cliente pudo haber perdido.

---

## Inicialización

```go
// modules/bundle.go
if redisClient != nil {
    events.New(router, database, logger, redisClient)
}
```

El módulo no se inicializa si Redis no está disponible, dado que su función central depende de él.

---

## Relación con integrations/events

`integrations/events` genera eventos durante sincronizaciones y los publica a Redis. Este módulo los consume desde el canal `ChannelIntegrationsSyncOrders`, aplica los filtros de `NotificationConfig`, y los entrega por SSE.

La dependencia es **unidireccional**: este módulo no importa ni conoce `integrations/events`. Solo comparte el nombre del canal, definido como constante en `shared/redis/channels.go`.

```
integrations/events ──(Redis)──▶ modules/events ──(SSE)──▶ Frontend
```

Ver [integrations/events/README.md](../../integrations/events/README.md) para entender qué eventos se publican y desde dónde.
