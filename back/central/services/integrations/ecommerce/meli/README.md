# MercadoLibre Integration Module

Integration type ID: **3** (`IntegrationTypeMercadoLibre`)

## Overview

Connects Probability with MercadoLibre to receive and synchronize orders. Supports real-time order ingestion via IPN (Instant Payment Notifications) and bulk synchronization via the Orders Search API.

All orders are mapped to the canonical `ProbabilityOrderDTO` format and published to the `probability.orders.canonical` RabbitMQ queue for downstream processing.

## Architecture

```
meli/
├── bundle.go                          # Wiring: client, service, queue, usecase, handler
└── internal/
    ├── domain/                        # Pure domain — zero external dependencies
    │   ├── entities.go                # MeliOrder, MeliBuyer, MeliPayment, MeliShipping, TokenResponse, etc.
    │   ├── ports.go                   # IMeliClient, IIntegrationService, OrderPublisher
    │   ├── errors.go                  # Domain errors (ErrTokenExpired, ErrRateLimited, etc.)
    │   └── query_params.go            # GetOrdersParams (offset/limit), GetOrdersResult
    ├── app/usecases/                  # Business logic
    │   ├── constructor.go             # IMeliUseCase interface + New()
    │   ├── test_connection.go         # Validates credentials via GET /users/me
    │   ├── process_notification.go    # IPN flow: fetch order → map → publish
    │   ├── sync_orders.go             # Bulk sync with pagination (offset/limit, max 50)
    │   ├── refresh_token.go           # OAuth token management (lazy refresh)
    │   └── mapper/
    │       └── order_mapper.go        # MapMeliOrderToProbability + status mappings
    └── infra/
        ├── primary/handlers/
        │   ├── constructor.go         # IHandler + RegisterRoutes
        │   └── handle_notification.go # POST /meli/notifications — respond 200, process async
        └── secondary/
            ├── client/                # HTTP client for MeLi REST API
            │   ├── constructor.go     # MeliClient + newAuthorizedRequest helper
            │   ├── get_order.go       # GET /orders/{id}
            │   ├── get_orders.go      # GET /orders/search?seller={id}
            │   ├── get_shipment.go    # GET /shipments/{id}
            │   ├── get_user.go        # GET /users/me
            │   ├── refresh_token.go   # POST /oauth/token
            │   └── response/          # JSON-tagged structs + ToDomain()
            │       ├── meli_order_response.go
            │       └── meli_shipping_response.go
            ├── core/                  # Adapter to integration core
            │   ├── core.go            # IIntegrationContract (TestConnection, Sync, GetWebhookURL)
            │   └── integration_service.go  # Adapter core → domain.IIntegrationService
            └── queue/                 # RabbitMQ publisher
                ├── rabbitmq_publisher.go
                ├── noop_publisher.go
                ├── mapper/
                │   └── canonical_order_mapper.go
                └── request/
                    └── canonical_order_dto.go
```

## MercadoLibre API Endpoints Used

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/users/me` | Test connection / get seller_id |
| GET | `/orders/{id}` | Fetch single order (from IPN notification) |
| GET | `/orders/search?seller={id}` | Bulk order sync with pagination |
| GET | `/shipments/{id}` | Shipping details (address, carrier, status) |
| POST | `/oauth/token` | Refresh access token |

Base URL: `https://api.mercadolibre.com`

## Authentication

MercadoLibre uses **OAuth 2.0** with short-lived tokens:

- **Access Token**: Valid for 6 hours
- **Refresh Token**: Used to obtain new access tokens
- **App ID + Client Secret**: Required for token refresh

Token management is **lazy** — `EnsureValidToken()` checks expiration before every API call and refreshes automatically if the token expires within 5 minutes.

### Credentials Storage

| Field | Storage | Description |
|-------|---------|-------------|
| `access_token` | Encrypted credentials | Current OAuth access token |
| `refresh_token` | Encrypted credentials | OAuth refresh token |
| `client_secret` | Encrypted credentials | App secret from MeLi developer portal |
| `app_id` | Integration config | App ID from MeLi developer portal |
| `seller_id` | Integration config + store_id | MeLi user ID of the seller |
| `token_expires_at` | Integration config | RFC3339 timestamp of token expiration |

## Notification Flow (IPN)

MercadoLibre sends IPN notifications to `POST /integrations/meli/notifications`:

```
MeLi IPN → Handler (respond 200) → goroutine →
  1. Parse notification body (topic + resource)
  2. Filter: only process "orders_v2" topic
  3. Extract order_id from resource ("/orders/123456789")
  4. Find integration by seller_id (notification.user_id → store_id)
  5. EnsureValidToken (refresh if needed)
  6. GET /orders/{id} → full order data
  7. GET /shipments/{id} → shipping details (optional, non-fatal)
  8. MapMeliOrderToProbability → canonical DTO
  9. Publish to RabbitMQ queue "probability.orders.canonical"
```

### Notification Body Example

```json
{
  "resource": "/orders/2000004456789",
  "user_id": 123456789,
  "topic": "orders_v2",
  "application_id": 1234567890,
  "attempts": 1,
  "sent": "2026-02-24T10:00:00.000-04:00",
  "received": "2026-02-24T10:00:00.100-04:00"
}
```

## Sync Flow

Bulk synchronization via `SyncOrders` (last 30 days) or `SyncOrdersWithParams`:

```
Core.SyncOrdersByIntegrationID → UseCase.SyncOrders →
  1. Get integration + seller_id
  2. EnsureValidToken
  3. Launch async goroutine:
     - Paginate: GET /orders/search (offset 0, 50, 100...)
     - For each order: GET /shipments/{id} → map → publish
     - Rate limit: 1s between pages
     - Auto-retry on token expiration
```

### Pagination

MeLi uses `offset/limit` (not page/per_page):
- Max 50 orders per request
- Response includes `paging.total` for total count
- Loop until `offset + limit >= total`

## Status Mappings

### Order Status

| MeLi Status | Probability Status |
|-------------|-------------------|
| confirmed | pending |
| payment_required | pending |
| payment_in_process | pending |
| paid | paid |
| partially_paid | paid |
| cancelled | cancelled |

### Payment Status

| MeLi Status | Probability Status |
|-------------|-------------------|
| approved | paid |
| pending / in_process / in_mediation | pending |
| rejected | failed |
| refunded | refunded |
| cancelled | cancelled |

### Shipping Status

| MeLi Status | Probability Status |
|-------------|-------------------|
| ready_to_ship | pending |
| shipped | shipped |
| delivered | delivered |
| not_delivered | failed |
| cancelled | cancelled |

## Integration Setup

To create a MercadoLibre integration:

1. Create an app at [MercadoLibre Developers](https://developers.mercadolibre.com/)
2. Complete the OAuth flow to get access + refresh tokens
3. Create integration via the core API:
   - **type_id**: 3
   - **store_id**: seller_id (MeLi user ID)
   - **config**: `{ "app_id": "...", "seller_id": 123456, "token_expires_at": "..." }`
   - **credentials**: `{ "access_token": "APP_USR-...", "refresh_token": "TG-...", "client_secret": "..." }`
4. Configure notification URL in MeLi developer portal:
   - URL: `https://api.yourdomain.com/integrations/meli/notifications`
   - Topics: `orders_v2`

## Testing

```bash
# Test connection
POST /api/integrations/{id}/test

# Manual notification simulation
POST /api/integrations/meli/notifications
Content-Type: application/json

{
  "resource": "/orders/123456789",
  "user_id": 123456,
  "topic": "orders_v2",
  "application_id": 12345,
  "attempts": 1,
  "sent": "2026-02-24T10:00:00.000-04:00",
  "received": "2026-02-24T10:00:00.100-04:00"
}

# Trigger bulk sync
POST /api/integrations/{id}/sync
```

## Key Differences vs WooCommerce

| Aspect | WooCommerce | MercadoLibre |
|--------|-------------|--------------|
| Auth | Basic Auth (key/secret) | OAuth 2.0 (token refresh every 6h) |
| Webhook payload | Full order in body | Only topic + resource URL |
| Pagination | page/per_page (headers) | offset/limit (JSON body) |
| Shipping data | Inline in order | Separate endpoint GET /shipments/{id} |
| Rate limit | 500ms between pages | 1s between pages |
| Webhook signature | HMAC-SHA256 (base64) | x-signature: ts={ts},v0={hash} |

## API Documentation References

### Official MercadoLibre Developer Docs

- **Developer Portal**: https://developers.mercadolibre.com.co/
- **Authentication (OAuth 2.0)**: https://developers.mercadolibre.com.co/es_co/autenticacion-y-autorizacion
- **Token Refresh**: https://developers.mercadolibre.com.co/es_co/registra-tu-aplicacion (sección OAuth)
- **Notifications (IPN)**: https://developers.mercadolibre.com.co/es_co/recibir-notificaciones
- **Orders API**: https://developers.mercadolibre.com.co/es_co/gestiona-ventas
- **Orders Search**: https://developers.mercadolibre.com.co/es_co/gestiona-ventas#Buscar-ordenes
- **Shipments API**: https://developers.mercadolibre.com.co/es_co/gestiona-envios
- **Items API**: https://developers.mercadolibre.com.co/es_co/publica-productos

### API Reference (Endpoints)

- **Auth**: https://api.mercadolibre.com/oauth/token
- **User**: https://api.mercadolibre.com/users/me
- **Orders**: https://api.mercadolibre.com/orders/{id}
- **Orders Search**: https://api.mercadolibre.com/orders/search?seller={id}
- **Shipments**: https://api.mercadolibre.com/shipments/{id}

### Guides & Resources

- **Crear aplicacion**: https://developers.mercadolibre.com.co/devcenter
- **IPN Topics disponibles**: `orders_v2`, `payments`, `items`, `shipments`, `questions`
- **Rate Limits**: https://developers.mercadolibre.com.co/es_co/api-docs (varía por endpoint, ~10k requests/hora)
- **Status codes de ordenes**: `confirmed`, `payment_required`, `payment_in_process`, `paid`, `partially_paid`, `cancelled`
- **Status codes de pagos**: `approved`, `pending`, `in_process`, `in_mediation`, `rejected`, `refunded`, `cancelled`
- **Status codes de envios**: `ready_to_ship`, `shipped`, `delivered`, `not_delivered`, `cancelled`
