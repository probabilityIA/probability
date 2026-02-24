# WooCommerce Integration

Integración con la REST API v3 de WooCommerce para sincronizar órdenes desde tiendas WordPress/WooCommerce hacia Probability.

**type_id:** 4 (`IntegrationTypeWoocommerce`)

---

## Funcionalidades

| Funcionalidad | Estado |
|---------------|--------|
| TestConnection | ✅ |
| SyncOrders (paginado, async) | ✅ |
| Webhook (HMAC, async) | ✅ |
| Order Mapper → Canonical | ✅ |
| Frontend Config Form | ✅ |
| ListWebhooks / CreateWebhook | ⏳ TODO |

---

## Autenticación

HTTP Basic Auth con `consumer_key` y `consumer_secret` generados desde el panel de WordPress.

```
WordPress Admin → WooCommerce → Ajustes → Avanzado → REST API → Agregar clave
```

| Campo | Tipo | Ejemplo |
|-------|------|---------|
| `store_url` | config | `https://mitienda.com` |
| `consumer_key` | credential | `ck_xxxxxxxxxxxxxxxxxxxx` |
| `consumer_secret` | credential | `cs_xxxxxxxxxxxxxxxxxxxx` |

---

## API REST v3

**Base URL:** `{store_url}/wp-json/wc/v3`

### Endpoints utilizados

| Método | Endpoint | Uso |
|--------|----------|-----|
| `GET` | `/system_status` | TestConnection |
| `GET` | `/orders` | SyncOrders (paginado) |
| `GET` | `/orders/{id}` | GetOrder individual |

### Paginación

WooCommerce usa paginación basada en página (no cursor):

- Query params: `page`, `per_page` (max 100), `after`, `before` (ISO 8601), `status`, `orderby`, `order`
- Response headers: `X-WP-Total` (total de registros), `X-WP-TotalPages`

---

## Webhook

### Configuración en WooCommerce

```
WordPress Admin → WooCommerce → Ajustes → Avanzado → Webhooks → Agregar webhook
```

- **URL de entrega:** `{BASE_URL}/integrations/woocommerce/webhook`
- **Método:** POST
- **Eventos:** `order.created`, `order.updated`, `order.deleted`, `order.restored`
- **Secret:** Configurar en WooCommerce y en la variable `WOOCOMMERCE_WEBHOOK_SECRET`

### Headers que envía WooCommerce

| Header | Descripción |
|--------|-------------|
| `X-WC-Webhook-Topic` | Evento (ej: `order.created`) |
| `X-WC-Webhook-Source` | URL de la tienda origen |
| `X-WC-Webhook-Signature` | `base64(HMAC-SHA256(secret, body))` |
| `X-WC-Webhook-ID` | ID del webhook en WooCommerce |
| `X-WC-Webhook-Delivery-ID` | ID único de la entrega |

### Validación HMAC

Opcional — solo se valida si `WOOCOMMERCE_WEBHOOK_SECRET` está configurada como variable de entorno. El handler responde `200 OK` inmediatamente y procesa la orden en background.

---

## Mapeo de estados

| WooCommerce | Probability | Descripción |
|-------------|-------------|-------------|
| `pending` | `pending` | Pago pendiente |
| `processing` | `paid` | Pago recibido, procesando |
| `on-hold` | `on_hold` | En espera de confirmación |
| `completed` | `fulfilled` | Orden completada y entregada |
| `cancelled` | `cancelled` | Cancelada |
| `refunded` | `refunded` | Reembolsada |
| `failed` | `failed` | Pago fallido |

---

## Variables de entorno

| Variable | Requerida | Descripción |
|----------|-----------|-------------|
| `WOOCOMMERCE_WEBHOOK_SECRET` | No | Secret para validar firma HMAC de webhooks |
| `RABBITMQ_ORDERS_CANONICAL_QUEUE` | No | Cola de destino (default: `probability.orders.canonical`) |

---

## Estructura del módulo

```
woocommerce/
├── README.md
├── bundle.go                          # Ensambla el módulo, retorna IIntegrationContract
└── internal/
    ├── domain/
    │   ├── entities.go                # Integration, WooCommerceOrder, Billing, Shipping, LineItem, etc.
    │   ├── errors.go                  # ErrMissingStoreURL, ErrWebhookInvalidSignature, etc.
    │   ├── ports.go                   # IWooCommerceClient, IIntegrationService, OrderPublisher
    │   └── query_params.go            # GetOrdersParams, GetOrdersResult
    ├── app/usecases/
    │   ├── constructor.go             # IWooCommerceUseCase (4 métodos)
    │   ├── test_connection.go         # Valida credenciales via GET /system_status
    │   ├── sync_orders.go             # Sync paginado async (30 días default, 500ms rate limit)
    │   ├── process_webhook.go         # Deserializa → mapea → publica a RabbitMQ
    │   └── mapper/
    │       └── order_mapper.go        # WooCommerceOrder → ProbabilityOrderDTO
    └── infra/
        ├── primary/handlers/
        │   ├── constructor.go         # IHandler + RegisterRoutes
        │   └── handle_webhook.go      # HMAC validation + async processing
        └── secondary/
            ├── client/
            │   ├── constructor.go     # HTTP client con Basic Auth
            │   ├── get_orders.go      # GetOrders (paginado), GetOrder (individual)
            │   └── response/
            │       └── woo_order_response.go  # JSON structs + ToDomain()
            ├── core/
            │   ├── core.go            # IIntegrationContract adapter
            │   └── integration_service.go
            └── queue/
                ├── mapper/canonical_order_mapper.go
                ├── rabbitmq_publisher.go
                ├── noop_publisher.go
                └── request/canonical_order_dto.go
```

---

## Documentación oficial de WooCommerce

- [REST API v3 — Referencia completa](https://woocommerce.github.io/woocommerce-rest-api-docs/)
- [REST API — Orders](https://woocommerce.github.io/woocommerce-rest-api-docs/#orders)
- [REST API — Autenticación](https://woocommerce.github.io/woocommerce-rest-api-docs/#authentication)
- [Webhooks — Guía](https://woocommerce.com/document/webhooks/)
- [Webhooks — REST API](https://woocommerce.github.io/woocommerce-rest-api-docs/#webhooks)
- [Order Statuses](https://woocommerce.com/document/managing-orders/#order-statuses)
- [Generar claves REST API](https://woocommerce.com/document/woocommerce-rest-api/)
