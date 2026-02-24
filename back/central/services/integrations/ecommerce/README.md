# Integrations — E-Commerce

Módulo que agrupa todos los proveedores de comercio electrónico soportados por Probability. Cada proveedor implementa el contrato `IIntegrationContract` del core y publica sus órdenes al canal canónico de RabbitMQ.

---

## Proveedores

| Proveedor | type_id | Estado | Paquete |
|-----------|---------|--------|---------|
| Shopify | 1 | Completo | `ecommerce/shopify` |
| MercadoLibre | 3 | Esqueleto | `ecommerce/meli` |
| WooCommerce | 4 | Funcional | `ecommerce/woocommerce` |
| VTEX | 16 | Esqueleto | `ecommerce/vtex` |
| Tiendanube | 17 | Esqueleto | `ecommerce/tiendanube` |
| Magento | 18 | Esqueleto | `ecommerce/magento` |
| Amazon | 19 | Esqueleto | `ecommerce/amazon` |
| Falabella | 20 | Esqueleto | `ecommerce/falabella` |
| Éxito | 21 | Esqueleto | `ecommerce/exito` |

---

## Estructura

```
ecommerce/
├── bundle.go               # Orquestador — inicializa y registra todos los proveedores
├── canonical/
│   └── order.go            # ProbabilityOrderDTO — formato canónico compartido
├── shopify/
│   ├── bundle.go
│   └── internal/
│       ├── domain/         # Entidades, ports, DTOs (sin tags JSON)
│       ├── app/usecases/   # Lógica de negocio
│       └── infra/
│           ├── primary/handlers/   # Endpoints HTTP (webhooks, OAuth)
│           └── secondary/
│               ├── client/         # Cliente HTTP Shopify API
│               ├── core/           # Adaptador → IIntegrationContract
│               └── queue/          # Publicador RabbitMQ
├── meli/
│   └── internal/  (misma estructura)
├── woocommerce/
│   └── internal/  (misma estructura)
├── vtex/
│   └── internal/  (misma estructura)
├── tiendanube/
│   └── internal/  (misma estructura)
├── magento/
│   └── internal/  (misma estructura)
├── amazon/
│   └── internal/  (misma estructura)
├── falabella/
│   └── internal/  (misma estructura)
└── exito/
    └── internal/  (misma estructura)
```

---

## Patrón de inicialización

El `ecommerce/bundle.go` es el único punto de entrada. Cada proveedor expone su `New()` que retorna un `IIntegrationContract`; el bundle padre hace el `RegisterIntegration`:

```go
// ecommerce/bundle.go
func New(router, logger, config, rabbitMQ, database, integrationCore) {
    // Shopify (type_id=1) — se auto-registra (incluye OnIntegrationCreated para webhooks automáticos)
    shopify.New(router, logger, config, integrationCore, rabbitMQ, database)

    // MercadoLibre (type_id=3)
    meliProvider := meli.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeMercadoLibre, meliProvider)

    // WooCommerce (type_id=4)
    wooProvider := woocommerce.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeWoocommerce, wooProvider)

    // VTEX (type_id=16)
    vtexProvider := vtex.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeVTEX, vtexProvider)

    // Tiendanube (type_id=17)
    tiendanubeProvider := tiendanube.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeTiendanube, tiendanubeProvider)

    // Magento (type_id=18)
    magentoProvider := magento.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeMagento, magentoProvider)

    // Amazon (type_id=19)
    amazonProvider := amazon.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeAmazon, amazonProvider)

    // Falabella (type_id=20)
    falabellaProvider := falabella.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeFalabella, falabellaProvider)

    // Éxito (type_id=21)
    exitoProvider := exito.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeExito, exitoProvider)
}
```

> Shopify maneja su propio registro internamente porque también configura un observer `OnIntegrationCreated` que crea webhooks automáticamente al activar una integración.

---

## Contrato de proveedor

Todo proveedor debe implementar `IIntegrationContract` (definido en `integrations/core/internal/domain/provider_contract.go`). Los métodos no soportados se heredan via `BaseIntegration` que retorna `ErrNotSupported`.

```go
type IIntegrationContract interface {
    TestConnection(ctx, config, credentials) error

    // Sincronización de órdenes
    SyncOrdersByIntegrationID(ctx, integrationID) error
    SyncOrdersByIntegrationIDWithParams(ctx, integrationID, params) error

    // Webhooks
    GetWebhookURL(ctx, baseURL, integrationID) (*WebhookInfo, error)
    ListWebhooks(ctx, integrationID) ([]interface{}, error)
    DeleteWebhook(ctx, integrationID, webhookID) error
    VerifyWebhooksByURL(ctx, integrationID, baseURL) ([]interface{}, error)
    CreateWebhook(ctx, integrationID, baseURL) (interface{}, error)
}
```

### Estado de implementación por proveedor

| Método | Shopify | MercadoLibre | WooCommerce | VTEX | Tiendanube | Magento | Amazon | Falabella | Éxito |
|--------|---------|--------------|-------------|------|------------|---------|--------|-----------|-------|
| `TestConnection` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ⏳ TODO | ✅ | ✅ |
| `SyncOrdersByIntegrationID` | ✅ | ⏳ TODO | ✅ | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO |
| `GetWebhookURL` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `HandleWebhook` | ✅ | ⏳ TODO | ✅ | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO |
| `ListWebhooks` | ✅ | ⬜ N/A | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO |
| `DeleteWebhook` | ✅ | ⬜ N/A | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO |
| `VerifyWebhooksByURL` | ✅ | ⬜ N/A | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO |
| `CreateWebhook` | ✅ | ⬜ N/A | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO | ⏳ TODO |

---

## WooCommerce — Detalle de implementación

### API REST v3

- **Autenticación:** HTTP Basic Auth (`consumer_key:consumer_secret`)
- **Base URL:** `{store_url}/wp-json/wc/v3`
- **Test endpoint:** `GET /system_status`
- **Órdenes:** `GET /orders` — paginación via headers `X-WP-Total`, `X-WP-TotalPages`

### Webhook

El handler en `POST /integrations/woocommerce/webhook` procesa los headers de WooCommerce:

| Header | Descripción |
|--------|-------------|
| `X-WC-Webhook-Topic` | Evento (`order.created`, `order.updated`, `order.deleted`, `order.restored`) |
| `X-WC-Webhook-Source` | URL de la tienda origen |
| `X-WC-Webhook-Signature` | HMAC-SHA256 del body codificado en base64 |

La validación HMAC es opcional — solo se activa si la variable de entorno `WOOCOMMERCE_WEBHOOK_SECRET` está configurada. El handler responde `200 OK` inmediatamente y procesa la orden de forma asíncrona.

### SyncOrders

- Por defecto sincroniza los últimos 30 días
- Soporta parámetros: `created_at_min`, `created_at_max`, `status`
- Paginación automática (hasta 100 órdenes por página)
- Rate limiting: 500ms entre páginas
- Ejecución asíncrona (retorna inmediatamente)

### Mapeo de estados

| WooCommerce | Probability |
|-------------|-------------|
| `pending` | `pending` |
| `processing` | `paid` |
| `on-hold` | `on_hold` |
| `completed` | `fulfilled` |
| `cancelled` | `cancelled` |
| `refunded` | `refunded` |
| `failed` | `failed` |

### Estructura de archivos

```
woocommerce/
├── bundle.go
└── internal/
    ├── domain/
    │   ├── entities.go        # Integration + WooCommerceOrder y sub-entidades
    │   ├── errors.go          # Errores específicos + HMAC/sync errors
    │   ├── ports.go           # IWooCommerceClient (TestConnection, GetOrders, GetOrder)
    │   └── query_params.go    # GetOrdersParams con ToQueryString()
    ├── app/usecases/
    │   ├── constructor.go     # IWooCommerceUseCase (4 métodos)
    │   ├── test_connection.go
    │   ├── sync_orders.go     # Sync con paginación async
    │   ├── process_webhook.go # Deserializar → mapear → publicar
    │   └── mapper/
    │       └── order_mapper.go  # WooCommerce → ProbabilityOrderDTO
    └── infra/
        ├── primary/handlers/
        │   ├── constructor.go
        │   └── handle_webhook.go  # HMAC + async processing
        └── secondary/
            ├── client/
            │   ├── constructor.go       # TestConnection (Basic Auth)
            │   ├── get_orders.go        # GetOrders, GetOrder
            │   └── response/
            │       └── woo_order_response.go  # JSON structs + ToDomain()
            ├── core/
            │   ├── core.go              # IIntegrationContract (Test, Sync, Webhook)
            │   └── integration_service.go
            └── queue/
                ├── mapper/              # Domain → Serializable
                ├── rabbitmq_publisher.go
                ├── noop_publisher.go
                └── request/             # Serializable DTOs con JSON tags
```

### Frontend

El formulario de configuración está en `front/central/src/services/integrations/ecommerce/woocommerce/ui/`. Campos:

- **Nombre de la Integración** (texto)
- **URL de la Tienda** (URL, requerida)
- **Consumer Key** (password, requerido)
- **Consumer Secret** (password, requerido)
- **Probar Conexión** → `testConnectionRawAction('woocommerce', config, credentials)`
- **Crear Integración** → `createIntegrationAction({ integration_type_id: 4, ... })`

---

## Formato canónico de órdenes (`canonical/`)

Todos los proveedores mapean sus órdenes a `ProbabilityOrderDTO` antes de publicar a RabbitMQ. Este DTO es la única fuente de verdad del formato de orden entre el módulo de integraciones y el módulo de órdenes.

```
Shopify Order      ──┐
Meli Order         ──┤
WooCommerce Order  ──┤
VTEX Order         ──┤
Tiendanube Order   ──┼──► ProbabilityOrderDTO ──► probability.orders.canonical (RabbitMQ)
Magento Order      ──┤
Amazon Order       ──┤
Falabella Order    ──┤
Éxito Order        ──┘
```

**Regla:** `canonical/order.go` **no tiene tags JSON**. Es dominio puro. La serialización con tags ocurre en `infra/secondary/queue/request/` de cada proveedor.

### Campos principales

| Campo | Descripción |
|-------|-------------|
| `IntegrationID` | ID de la integración origen |
| `IntegrationType` | Código del proveedor (`shopify`, `meli`, `woocommerce`, etc.) |
| `ExternalID` | ID de la orden en la plataforma externa |
| `OrderNumber` | Número de orden legible |
| `TotalAmount` | Total en moneda de la tienda |
| `TotalAmountPresentment` | Total en moneda local del cliente |
| `OrderItems` | Líneas de producto |
| `Payments` | Pagos asociados |
| `Shipments` | Envíos asociados |
| `Addresses` | Direcciones (billing, shipping) |
| `ChannelMetadata` | Datos crudos del canal y estado de sincronización |

---

## Credenciales por proveedor

Las credenciales se almacenan cifradas en la tabla `integrations` y se acceden vía `DecryptCredential` del core.

### Shopify (type_id=1)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `shop_domain` | credential | Dominio de la tienda (`mitienda.myshopify.com`) |
| `access_token` | credential | Token de acceso OAuth |
| `client_secret` | credential | Secret para validar webhooks HMAC |

### MercadoLibre (type_id=3)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `access_token` | credential | Token OAuth de la cuenta del vendedor |

### WooCommerce (type_id=4)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `store_url` | config | URL base de la tienda (`https://mitienda.com`) |
| `consumer_key` | credential | Consumer Key de la REST API |
| `consumer_secret` | credential | Consumer Secret de la REST API |

### VTEX (type_id=16)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `store_url` | config | URL de la tienda VTEX (`https://cuenta.vtexcommercestable.com.br`) |
| `api_key` | credential | App Key (`X-VTEX-API-AppKey`) |
| `api_token` | credential | App Token (`X-VTEX-API-AppToken`) |

### Tiendanube (type_id=17)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `store_url` | config | URL de la tienda |
| `access_token` | credential | Token OAuth de la aplicación |

### Magento (type_id=18)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `store_url` | config | URL base de Magento/Adobe Commerce |
| `access_token` | credential | Integration Access Token (Bearer) |

### Amazon (type_id=19)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `seller_id` | config | ID del vendedor en el marketplace |
| `refresh_token` | credential | Refresh token para SP-API OAuth |
| `client_id` | credential | Client ID de la aplicación SP-API |
| `client_secret` | credential | Client Secret de la aplicación SP-API |

> **Nota:** Amazon requiere un flujo OAuth completo (SP-API). `TestConnection` aún no está implementado.

### Falabella (type_id=20)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `user_id` | config | ID del vendedor en Falabella Seller Center |
| `api_key` | credential | API Key del vendedor |

### Éxito (type_id=21)
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `seller_id` | config | ID del vendedor en el marketplace Éxito |
| `api_key` | credential | API Key del vendedor |

---

## Agregar un nuevo proveedor

1. Crear carpeta `ecommerce/<proveedor>/` con la estructura hexagonal estándar
2. Implementar `IIntegrationContract` (embeber `BaseIntegration` y sobrescribir los métodos soportados)
3. Implementar el mapper `canonical.ProbabilityOrderDTO` en `infra/secondary/queue/mapper/`
4. Definir el `type_id` en `integrations/core/internal/domain/type_codes.go`
5. Registrar en `ecommerce/bundle.go`:
   ```go
   miProvider := mi_proveedor.New(router, logger, config, rabbitMQ, integrationCore)
   integrationCore.RegisterIntegration(core.IntegrationTypeMiProveedor, miProvider)
   ```
