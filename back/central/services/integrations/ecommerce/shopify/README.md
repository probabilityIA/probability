# Shopify Integration

Módulo de integración con Shopify. Recibe eventos de órdenes vía webhooks, los mapea al formato canónico de Probability y los publica a RabbitMQ para procesamiento downstream.

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

---

## Arquitectura

```
bundle.go
│
├── infra/primary/handlers/          ← Adaptadores de entrada
│   ├── WebhookHandler               ← Recibe webhooks de Shopify (HMAC + async)
│   ├── OAuthHandlers                ← Flujo OAuth2
│   └── ComplianceWebhookHandler     ← GDPR/CCPA
│
├── internal/app/usecases/           ← Lógica de negocio
│   ├── CreateOrder
│   ├── ProcessOrderPaid/Updated/Cancelled/Fulfilled/PartiallyFulfilled
│   ├── SyncOrders
│   ├── CreateWebhook / ListWebhooks / DeleteWebhook / VerifyWebhooks
│   └── TestConnection
│
└── infra/secondary/                 ← Adaptadores de salida
    ├── client/         ← HTTP client Shopify API (Resty, 30s timeout)
    ├── queue/          ← Publisher RabbitMQ (canonical orders)
    └── core/           ← ShopifyCore — implementa core.IIntegrationContract
```

---

## Dependencias

### Entrada (lo que recibe)

| Fuente | Mecanismo | Descripción |
|--------|-----------|-------------|
| Shopify | `POST /integrations/shopify/webhook` | Eventos de órdenes vía webhook |
| Shopify | `GET /shopify/callback` | OAuth2 redirect |
| API interna | HTTP REST | Sync manual, gestión de webhooks |

### Salida (lo que produce)

| Destino | Mecanismo | Descripción |
|---------|-----------|-------------|
| RabbitMQ | `probability.orders.canonical` | Órdenes en formato canónico |
| Shopify API | HTTP (Resty) | Crear/listar/eliminar webhooks, fetch órdenes |

### Servicios externos requeridos

| Servicio | Interfaz | Para qué |
|----------|----------|----------|
| `core.IIntegrationService` | Inyectada | Buscar integración por store_id, descifrar credenciales, actualizar config |
| RabbitMQ | `rabbitmq.IQueue` | Publicar órdenes (opcional — hay no-op fallback) |
| Base de datos | `db.IDatabase` | Persistencia de órdenes |

---

## Flujo principal: Webhook → RabbitMQ

```
Shopify → POST /integrations/shopify/webhook
              │
              ├─ Extrae: X-Shopify-Topic, X-Shopify-Shop-Domain, X-Shopify-HMAC-SHA256
              ├─ Lee body como bytes
              ├─ Valida HMAC-SHA256 (secret por tienda, fallback a secret global)
              ├─ Responde 200 OK inmediato
              │
              └─ Goroutine async:
                    │
                    ├─ Parse JSON → ShopifyOrder
                    ├─ Dispatch por topic:
                    │     orders/create             → CreateOrder()
                    │     orders/paid               → ProcessOrderPaid()
                    │     orders/updated            → ProcessOrderUpdated()
                    │     orders/cancelled          → ProcessOrderCancelled()
                    │     orders/fulfilled          → ProcessOrderFulfilled()
                    │     orders/partially_fulfilled→ ProcessOrderPartiallyFulfilled()
                    │
                    └─ Cada handler:
                          1. GetIntegrationByExternalID(shopDomain, typeID=1)
                          2. Llenar ShopifyOrder: BusinessID, IntegrationID
                          3. MapShopifyOrderToProbability() → ProbabilityOrderDTO
                          4. EnrichOrderWithDetails() (payment, fulfillment, financial)
                          5. Agregar ChannelMetadata con JSON raw de Shopify
                          6. orderPublisher.Publish() → RabbitMQ
```

---

## Auto-creación de webhooks

Al crear una integración de tipo Shopify, `bundle.go` registra un observador que automáticamente crea los 6 webhooks necesarios en la tienda:

```
OnIntegrationCreated(ShopifyTypeID) →
    CreateWebhook(integrationID, baseURL) →
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
POST   /integrations/shopify/webhook                    ← Webhooks de órdenes (HMAC)
POST   /integrations/shopify/webhook/:integration_id    ← Ruta alternativa
POST   /integrations/shopify/webhooks/compliance        ← Compliance GDPR/CCPA
POST   /integrations/shopify/connect                    ← Iniciar OAuth (JWT)
POST   /integrations/shopify/connect/custom             ← OAuth custom (JWT)
GET    /shopify/callback                                 ← OAuth2 callback
POST   /integrations/shopify/auth/login                 ← Login con session token
POST   /integrations/shopify/config                     ← Obtener configuración
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
```

---

## Adapter de core

`infra/secondary/core/ShopifyCore` implementa `core.IIntegrationContract` y delega todo al use case.
Es lo que core invoca cuando opera sobre una integración de tipo Shopify (test connection, sync, webhooks).

```
core → ShopifyCore.TestConnection()           → useCase.TestConnection()
core → ShopifyCore.SyncOrdersByIntegrationID() → useCase.SyncOrders()
core → ShopifyCore.CreateWebhook()            → useCase.CreateWebhook()
core → ShopifyCore.GetWebhookURL()            → construye "{baseURL}/integrations/shopify/webhook"
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
