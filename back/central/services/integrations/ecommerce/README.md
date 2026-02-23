# Integrations — E-Commerce

Módulo que agrupa todos los proveedores de comercio electrónico soportados por Probability. Cada proveedor implementa el contrato `IIntegrationContract` del core y publica sus órdenes al canal canónico de RabbitMQ.

---

## Proveedores

| Proveedor | type_id | Estado | Paquete |
|-----------|---------|--------|---------|
| Shopify | 1 | Completo | `ecommerce/shopify` |
| MercadoLibre | 3 | Esqueleto | `ecommerce/meli` |
| WooCommerce | 4 | Esqueleto | `ecommerce/woocommerce` |

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
└── woocommerce/
    └── internal/  (misma estructura)
```

---

## Patrón de inicialización

El `ecommerce/bundle.go` es el único punto de entrada. Cada proveedor expone su `New()` que retorna un `IIntegrationContract`; el bundle padre hace el `RegisterIntegration`:

```go
// ecommerce/bundle.go
func New(router, logger, config, rabbitMQ, database, integrationCore) {
    // Shopify — se auto-registra (incluye OnIntegrationCreated para webhooks automáticos)
    shopify.New(router, logger, config, integrationCore, rabbitMQ, database)

    // MercadoLibre
    meliProvider := meli.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeMercadoLibre, meliProvider)

    // WooCommerce
    wooProvider := woocommerce.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeWoocommerce, wooProvider)
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

| Método | Shopify | MercadoLibre | WooCommerce |
|--------|---------|--------------|-------------|
| `TestConnection` | ✅ | ✅ | ✅ |
| `SyncOrdersByIntegrationID` | ✅ | ⏳ TODO | ⏳ TODO |
| `GetWebhookURL` | ✅ | ✅ | ✅ |
| `ListWebhooks` | ✅ | ⬜ N/A | ⏳ TODO |
| `DeleteWebhook` | ✅ | ⬜ N/A | ⏳ TODO |
| `VerifyWebhooksByURL` | ✅ | ⬜ N/A | ⏳ TODO |
| `CreateWebhook` | ✅ | ⬜ N/A | ⏳ TODO |

---

## Formato canónico de órdenes (`canonical/`)

Todos los proveedores mapean sus órdenes a `ProbabilityOrderDTO` antes de publicar a RabbitMQ. Este DTO es la única fuente de verdad del formato de orden entre el módulo de integraciones y el módulo de órdenes.

```
Shopify Order  ──┐
Meli Order     ──┼──► ProbabilityOrderDTO ──► probability.orders.canonical (RabbitMQ)
WooCommerce    ──┘
```

**Regla:** `canonical/order.go` **no tiene tags JSON**. Es dominio puro. La serialización con tags ocurre en `infra/secondary/queue/request/` de cada proveedor.

### Campos principales

| Campo | Descripción |
|-------|-------------|
| `IntegrationID` | ID de la integración origen |
| `IntegrationType` | Código del proveedor (`shopify`, `meli`, `woocommerce`) |
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

### Shopify
| Campo | Descripción |
|-------|-------------|
| `shop_domain` | Dominio de la tienda (`mitienda.myshopify.com`) |
| `access_token` | Token de acceso OAuth |
| `client_secret` | Secret para validar webhooks HMAC |

### MercadoLibre
| Campo | Descripción |
|-------|-------------|
| `access_token` | Token OAuth de la cuenta del vendedor |

### WooCommerce
| Campo | Descripción |
|-------|-------------|
| `store_url` | URL base de la tienda (`https://mitienda.com`) |
| `consumer_key` | Consumer Key de la REST API de WooCommerce |
| `consumer_secret` | Consumer Secret de la REST API |

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
