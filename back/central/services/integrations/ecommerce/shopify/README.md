# Shopify Integration

Módulo de integración con Shopify. Recibe eventos de órdenes vía webhooks y sincroniza órdenes bajo demanda. Mapea al formato canónico de Probability y publica a RabbitMQ para procesamiento downstream.

---

## Resumen

| Aspecto | Valor |
|---------|-------|
| Integration Type ID | `1` |
| API Shopify | `2024-10` |
| Queue de salida | `probability.orders.canonical` |
| Validación de webhooks | HMAC-SHA256 por tienda |
| Procesamiento | Asíncrono (goroutine) — responde 200 OK inmediato |
| Auto-creación de webhooks | Sí, al crear una integración |
| Sync por lotes | Rangos > 14 días se dividen en chunks de 7 días vía RabbitMQ |

---

## Arquitectura

```
bundle.go
|
+-- infra/primary/handlers/          <- Adaptadores de entrada
|   +-- WebhookHandler               <- Recibe webhooks de Shopify (HMAC + async)
|   +-- OAuthHandlers                <- Flujo OAuth2
|   +-- ComplianceWebhookHandler     <- GDPR/CCPA
|
+-- internal/app/usecases/           <- Lógica de negocio
|   +-- CreateOrder
|   +-- ProcessOrderPaid/Updated/Cancelled/Fulfilled/PartiallyFulfilled
|   +-- SyncOrders / SyncOrdersWithParams / GetOrders
|   +-- CreateWebhook / ListWebhooks / DeleteWebhook / VerifyWebhooks
|   +-- TestConnection
|
+-- infra/secondary/                 <- Adaptadores de salida
    +-- client/         <- HTTP client Shopify API (Resty, 30s timeout)
    +-- queue/          <- Publisher RabbitMQ (canonical orders)
    +-- core/           <- ShopifyCore — implementa core.IIntegrationContract
```

---

## Dependencias

### Entrada (lo que recibe)

| Fuente | Mecanismo | Descripción |
|--------|-----------|-------------|
| Shopify | `POST /integrations/shopify/webhook` | Eventos de órdenes vía webhook |
| Shopify | `GET /shopify/callback` | OAuth2 redirect |
| API interna | HTTP REST | Sync manual, gestión de webhooks |
| RabbitMQ | `integration.sync.batches` | Lotes de sincronización (via módulo core) |

### Salida (lo que produce)

| Destino | Mecanismo | Descripción |
|---------|-----------|-------------|
| RabbitMQ | `probability.orders.canonical` | Órdenes en formato canónico |
| RabbitMQ | `events.exchange` | Eventos SSE: `integration.sync.completed`, `integration.sync.failed` |
| Shopify API | HTTP (Resty) | Crear/listar/eliminar webhooks, fetch órdenes |

### Servicios externos requeridos

| Servicio | Interfaz | Para qué |
|----------|----------|----------|
| `core.IIntegrationService` | Inyectada | Buscar integración por store_id, descifrar credenciales, actualizar config |
| RabbitMQ | `rabbitmq.IQueue` | Publicar órdenes y eventos SSE |
| Base de datos | `db.IDatabase` | Persistencia de órdenes |

---

## Flujo de sincronización: Sync -> RabbitMQ -> SSE

La sincronización de órdenes tiene dos modos: **directo** (rango <= 14 días) y **por lotes** (rango > 14 días). Ambos convergen en el mismo pipeline de eventos.

### Diagrama completo

```
+-------------------------------------------------------------------------+
|  FRONTEND (IntegrationList.tsx)                                         |
|                                                                         |
|  Usuario presiona "Iniciar Sincronización"                              |
|       |                                                                 |
|       ▼                                                                 |
|  POST /integrations/:id/sync -> HTTP 202 Accepted                       |
|       |                                                                 |
|       ▼                                                                 |
|  Escucha SSE (módulo events) esperando eventos de progreso              |
+---------------------+---------------------------------------------------+
                      |
+---------------------▼---------------------------------------------------+
|  MÓDULO INTEGRATIONS/CORE (handler sync-orders.go)                      |
|                                                                         |
|  ¿Rango > 14 días?                                                      |
|       |                                                                 |
|       +- SÍ (batch) ------------------+                                 |
|       |   SplitDateRange(7 días)      |                                 |
|       |   Publica N msgs a cola       |                                 |
|       |   integration.sync.batches    |                                 |
|       |        |                      |                                 |
|       |        ▼                      |                                 |
|       |   SSE: batched.started        |                                 |
|       |   { total_batches, ... }      |                                 |
|       |                               |                                 |
|       +- NO (directo) --- go func()   |                                 |
|       |   Llama SyncOrders()          |                                 |
|       |   directamente                |                                 |
|       |                               |                                 |
+-------+-------------------------------+---------------------------------+
        |                               |
        |    +--------------------------+|
        |    |  BATCH CONSUMER          ||
        |    |  (sync_batch_consumer)   ||
        |    |                          ||
        |    |  Para cada lote:         ||
        |    |  Resuelve provider ->     ||
        |    |  provider.SyncOrders()   ||
        |    |       |                  ||
        |    |       ▼                  ||
        |    |  SSE: batch.completed    ||
        |    |  o batch.failed          ||
        |    +------+-------------------+|
        |           |                    |
        ▼           ▼                    |
+-----------------------------------------------------------------------+
|  MÓDULO SHOPIFY (usecases/sync_orders.go)                             |
|                                                                       |
|  SyncOrdersWithParams():                                              |
|       |                                                               |
|       ▼                                                               |
|  GetOrders() <- llama Shopify API paginado                            |
|       |                                                               |
|       +- Por cada página de órdenes:                                  |
|       |   MapShopifyOrderToProbability()                              |
|       |   orderPublisher.Publish() ----------+                        |
|       |                                      |                        |
|       ▼                                      |                        |
|  Retorna (totalFetched, error)               |                        |
|       |                                      |                        |
|       ▼                                      |                        |
|  SSE: integration.sync.completed             |                        |
|  { total_fetched, duration }                 |                        |
|  --- o ---                                   |                        |
|  SSE: integration.sync.failed                |                        |
|  { error }                                   |                        |
+----------------------------------------------+------------------------+
                                               |
                                               ▼
+----------------------------------------------------------------------+
|  RABBITMQ: probability.orders.canonical                              |
|                                                                      |
|  Cola con órdenes en formato ProbabilityOrderDTO                     |
|  (una orden = un mensaje)                                            |
+--------------------------+-------------------------------------------+
                           |
                           ▼
+----------------------------------------------------------------------+
|  MÓDULO ORDERS (queue/consumer.go) <- OrderConsumer                   |
|                                                                      |
|  handleMessage():                                                    |
|       |                                                              |
|       ▼                                                              |
|  Deserializa ProbabilityOrderDTO                                     |
|       |                                                              |
|       ▼                                                              |
|  createUC.MapAndSaveOrder()                                          |
|       |                                                              |
|       +- Orden nueva    -> guarda en DB -> SSE: order.created          |
|       +- Orden existente-> actualiza DB -> SSE: order.updated          |
|       +- Error          -> guarda en order_errors                     |
|                           -> SSE: order.rejected { reason }           |
|                                                                      |
|  GARANTÍA: cada orden publicada a la cola genera                     |
|  EXACTAMENTE 1 evento SSE (created | updated | rejected)            |
|  Esto permite que el frontend complete la barra de progreso:         |
|  created + updated + rejected === totalFetched                       |
+--------------------------+-------------------------------------------+
                           |
                           ▼
+----------------------------------------------------------------------+
|  RABBITMQ: events.exchange (topic exchange)                          |
|                                                                      |
|  Todos los eventos SSE se publican aquí como EventEnvelope:          |
|  { type, category, business_id, integration_id, data }               |
|                                                                      |
|  Routing key = event type (e.g. "integration.sync.order.created")    |
+--------------------------+-------------------------------------------+
                           |
                           ▼
+----------------------------------------------------------------------+
|  MÓDULO EVENTS (consumer/rabbitmq_consumer.go)                       |
|                                                                      |
|  Consume de "events.unified" (bindeada a events.exchange con "#")    |
|       |                                                              |
|       ▼                                                              |
|  EventDispatcher.HandleEvent()                                       |
|       |                                                              |
|       +- Busca notification_configs en Redis cache                   |
|       +- Rutea por canal: SSE | WhatsApp | Email                     |
|       +- Si no hay configs -> broadcast SSE por defecto               |
|              |                                                       |
|              ▼                                                       |
|  ssePublisher.PublishEvent() -> SSE push al frontend                  |
+--------------------------+-------------------------------------------+
                           |
                           ▼
+----------------------------------------------------------------------+
|  FRONTEND (IntegrationList.tsx)                                      |
|                                                                      |
|  useSSE() recibe eventos y actualiza estado:                         |
|                                                                      |
|  Eventos de sincronización:                                          |
|  +- integration.sync.started        -> inicializa contadores          |
|  +- integration.sync.completed      -> guarda totalFetched            |
|  +- integration.sync.failed         -> muestra error                  |
|  |                                                                   |
|  Eventos por lote (rango > 14 días):                                 |
|  +- integration.sync.batched.started -> inicializa N lotes            |
|  +- integration.sync.batch.completed -> marca lote verde              |
|  +- integration.sync.batch.failed    -> marca lote rojo               |
|  |                                                                   |
|  Eventos por orden (alimentan barra de progreso):                    |
|  +- integration.sync.order.created   -> +1 creada (verde)            |
|  +- integration.sync.order.updated   -> +1 actualizada (amarillo)    |
|  +- integration.sync.order.rejected  -> +1 rechazada (rojo)          |
|                                                                      |
|  Barra: created + updated + rejected / totalFetched = progreso       |
|  Cuando created + updated + rejected >= totalFetched -> 100%          |
+----------------------------------------------------------------------+
```

### Eventos SSE completos

| Evento | Origen | Datos | Cuándo |
|--------|--------|-------|--------|
| `integration.sync.started` | Shopify use case | `integration_id` | Al iniciar fetch desde Shopify API |
| `integration.sync.completed` | Shopify use case | `total_fetched`, `duration` | Todas las órdenes publicadas a cola |
| `integration.sync.failed` | Shopify use case | `error` | Error fatal en el fetch |
| `integration.sync.batched.started` | Core handler | `total_batches`, `date_from`, `date_to`, `chunk_days` | Rango > 14 días, se crean N lotes |
| `integration.sync.batch.completed` | Core batch consumer | `batch_index`, `duration`, `created_at_min`, `created_at_max` | Un lote procesado OK |
| `integration.sync.batch.failed` | Core batch consumer | `batch_index`, `error`, `duration` | Un lote falló |
| `integration.sync.order.created` | Orders create UC | `order_number`, `created_at`, `status` | Orden nueva guardada en DB |
| `integration.sync.order.updated` | Orders update UC | `order_number`, `updated_at`, `status` | Orden existente actualizada |
| `integration.sync.order.rejected` | Orders consumer | `order_number`, `external_id`, `reason` | Orden falló al procesarse |

### Garantía de completitud

Cada orden publicada a `probability.orders.canonical` genera **exactamente 1** evento SSE:

- `order.created` si es nueva
- `order.updated` si ya existía
- `order.rejected` si falló (duplicada, FK violation, error de mapeo, etc.)

Esto garantiza que `created + updated + rejected === totalFetched`, permitiendo que la barra de progreso del frontend llegue a 100% sin necesidad de timeouts.

---

## Flujo principal: Webhook -> RabbitMQ

```
Shopify -> POST /integrations/shopify/webhook
              |
              +- Extrae: X-Shopify-Topic, X-Shopify-Shop-Domain, X-Shopify-HMAC-SHA256
              +- Lee body como bytes
              +- Valida HMAC-SHA256 (secret por tienda, fallback a secret global)
              +- Responde 200 OK inmediato
              |
              +- Goroutine async:
                    |
                    +- Parse JSON -> ShopifyOrder
                    +- Dispatch por topic:
                    |     orders/create             -> CreateOrder()
                    |     orders/paid               -> ProcessOrderPaid()
                    |     orders/updated            -> ProcessOrderUpdated()
                    |     orders/cancelled          -> ProcessOrderCancelled()
                    |     orders/fulfilled          -> ProcessOrderFulfilled()
                    |     orders/partially_fulfilled-> ProcessOrderPartiallyFulfilled()
                    |
                    +- Cada handler:
                          1. GetIntegrationByExternalID(shopDomain, typeID=1)
                          2. Llenar ShopifyOrder: BusinessID, IntegrationID
                          3. MapShopifyOrderToProbability() -> ProbabilityOrderDTO
                          4. EnrichOrderWithDetails() (payment, fulfillment, financial)
                          5. Agregar ChannelMetadata con JSON raw de Shopify
                          6. orderPublisher.Publish() -> RabbitMQ
```

---

## Auto-creación de webhooks

Al crear una integración de tipo Shopify, `bundle.go` registra un observador que automáticamente crea los 6 webhooks necesarios en la tienda:

```
OnIntegrationCreated(ShopifyTypeID) ->
    CreateWebhook(integrationID, baseURL) ->
        Crea en Shopify:
            orders/create
            orders/paid
            orders/updated
            orders/cancelled
            orders/fulfilled
            orders/partially_fulfilled
        Guarda webhook IDs en integration.config
```

URL del webhook: `{WEBHOOK_BASE_URL}/integrations/shopify/webhook`

---

## Rutas HTTP

```
POST   /integrations/shopify/webhook                    <- Webhooks de órdenes (HMAC)
POST   /integrations/shopify/webhook/:integration_id    <- Ruta alternativa
POST   /integrations/shopify/webhooks/compliance        <- Compliance GDPR/CCPA
POST   /integrations/shopify/connect                    <- Iniciar OAuth (JWT)
POST   /integrations/shopify/connect/custom             <- OAuth custom (JWT)
GET    /shopify/callback                                 <- OAuth2 callback
POST   /integrations/shopify/auth/login                 <- Login con session token
POST   /integrations/shopify/config                     <- Obtener configuración
```

---

## Interfaces del dominio

```go
// Puertos secundarios — lo que el use case necesita
type IIntegrationService interface {
    GetIntegrationByID(ctx, integrationID string) (*Integration, error)
    GetIntegrationByExternalID(ctx, externalID string, integrationType int) (*Integration, error)
    DecryptCredential(ctx, integrationID, fieldName string) (string, error)
    UpdateIntegrationConfig(ctx, integrationID string, config map[string]interface{}) error
}

type ShopifyClient interface {
    ValidateToken(ctx, storeName, accessToken string) (bool, map[string]interface{}, error)
    GetOrders(ctx, storeName, accessToken string, params *GetOrdersParams) ([]ShopifyOrder, string, error)
    GetOrder(ctx, storeName, accessToken, orderID string) (*ShopifyOrder, error)
    CreateWebhook(ctx, storeName, accessToken, webhookURL, event string) (string, error)
    ListWebhooks(ctx, storeName, accessToken string) ([]WebhookInfo, error)
    DeleteWebhook(ctx, storeName, accessToken, webhookID string) error
}

type OrderPublisher interface {
    Publish(ctx context.Context, order *ProbabilityOrderDTO) error
}

type ISyncEventPublisher interface {
    PublishSyncEvent(ctx, integrationID uint, businessID uint, eventType string, data map[string]interface{})
}
```

---

## Adapter de core

`infra/secondary/core/ShopifyCore` implementa `core.IIntegrationContract` y delega todo al use case.
Es lo que core invoca cuando opera sobre una integración de tipo Shopify (test connection, sync, webhooks).

```
core -> ShopifyCore.TestConnection()           -> useCase.TestConnection()
core -> ShopifyCore.SyncOrdersByIntegrationID() -> useCase.SyncOrders()
core -> ShopifyCore.CreateWebhook()            -> useCase.CreateWebhook()
core -> ShopifyCore.GetWebhookURL()            -> construye "{baseURL}/integrations/shopify/webhook"
```

---

## Variables de entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `WEBHOOK_BASE_URL` | URL base para webhooks | — |
| `URL_BASE_SWAGGER` | Fallback si no hay WEBHOOK_BASE_URL | — |
| `SHOPIFY_WEBHOOK_SECRET` | Secret global HMAC (fallback) | — |
| `SHOPIFY_DEBUG` | Habilitar logs HTTP del cliente | `false` |
| `RABBITMQ_ORDERS_CANONICAL_QUEUE` | Nombre de la queue de salida | `probability.orders.canonical` |
