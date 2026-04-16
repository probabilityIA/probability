# Modulo E-commerce - Integraciones

## Descripcion

El modulo `ecommerce` agrupa todas las integraciones con plataformas de comercio electronico y marketplaces. Su responsabilidad principal es **recibir ordenes desde plataformas externas** (via webhooks, notificaciones IPN o sincronizacion manual), **transformarlas a un formato canonico unificado** (`ProbabilityOrderDTO`) y **publicarlas a la cola de RabbitMQ** `probability.orders.canonical` para que el modulo `orders` las persista en base de datos.

Cada plataforma se implementa como un sub-modulo independiente con arquitectura hexagonal propia.

---

## Proveedores Soportados

| Proveedor     | type_id | Paquete               | Estado       | Mecanismo de ingesta                      |
|---------------|---------|----------------------|--------------|-------------------------------------------|
| Shopify       | 1       | `ecommerce/shopify`       | Completo     | Webhooks + OAuth + Sync + Compliance GDPR |
| MercadoLibre  | 3       | `ecommerce/meli`          | Completo     | Notificaciones IPN + Sync + RefreshToken  |
| WooCommerce   | 4       | `ecommerce/woocommerce`   | Completo     | Webhook + Sync                            |
| VTEX          | 16      | `ecommerce/vtex`          | Completo     | Webhook Hook v1 + Sync                    |
| Tiendanube    | 17      | `ecommerce/tiendanube`    | Esqueleto    | Webhook (solo TestConnection)             |
| Magento       | 18      | `ecommerce/magento`       | Esqueleto    | Webhook (solo TestConnection)             |
| Amazon        | 19      | `ecommerce/amazon`        | Esqueleto    | Notificacion SQS/SNS (solo TestConnection)|
| Falabella     | 20      | `ecommerce/falabella`     | Esqueleto    | Webhook (solo TestConnection)             |
| Exito         | 21      | `ecommerce/exito`         | Esqueleto    | Webhook (solo TestConnection)             |

Los proveedores marcados como "Esqueleto" tienen la estructura hexagonal completa, publisher de RabbitMQ y endpoint de webhook, pero solo implementan `TestConnection`. La logica de sincronizacion y procesamiento real de ordenes queda pendiente.

---

## Flujo General de Datos

```
Plataforma externa (Shopify, MeLi, VTEX, WooCommerce, etc.)
        |
        | Webhook / IPN / Sync API
        v
  Handler HTTP (infra/primary/handlers)
        |
        | Parsea el payload, responde 200 inmediatamente
        | Procesa en goroutine aparte (async)
        v
  Use Case (app/usecases)
        |
        | Obtiene credenciales via IntegrationService (core)
        | Consulta orden completa a la API del proveedor
        | Mapea a formato canonico (ProbabilityOrderDTO)
        v
  OrderPublisher (infra/secondary/queue)
        |
        | Mapea dominio -> DTO serializable (con tags JSON)
        | Serializa a JSON y publica a RabbitMQ
        v
  Cola: probability.orders.canonical
        |
        | Consumer en modulo orders
        v
  Modulo orders -> Base de datos
```

---

## Formato Canonico (`canonical/order.go`)

Todos los proveedores mapean sus ordenes a `ProbabilityOrderDTO` antes de publicar. Este DTO es dominio puro (sin etiquetas JSON). La serializacion con tags ocurre en `infra/secondary/queue/request/` de cada proveedor.

### Entidades principales

| Entidad                        | Descripcion                                                        |
|--------------------------------|--------------------------------------------------------------------|
| `ProbabilityOrderDTO`          | Orden completa: financieros, cliente, estado, metadata, presentment|
| `ProbabilityOrderItemDTO`      | Linea de producto: SKU, precios, descuentos, impuestos, peso      |
| `ProbabilityAddressDTO`        | Direccion (shipping/billing) con coordenadas opcionales            |
| `ProbabilityPaymentDTO`        | Pago: metodo, gateway, estado, transaccion, reembolso              |
| `ProbabilityShipmentDTO`       | Envio: tracking, transportadora, guia, dimensiones, warehouse      |
| `ProbabilityChannelMetadataDTO`| Metadata del canal: datos crudos, estado de sincronizacion         |

Soporta precios en moneda presentment (moneda local del cliente) para plataformas multi-moneda como Shopify.

---

## Endpoints HTTP

Todas las rutas se registran bajo el prefijo del router (tipicamente `/api/v1`).

### Shopify (`/integrations/shopify`)

| Metodo | Ruta                                              | Auth         | Descripcion                        |
|--------|---------------------------------------------------|--------------|-------------------------------------|
| GET    | `/integrations/shopify/config`                    | Publica      | Configuracion de la app Shopify     |
| POST   | `/integrations/shopify/auth/login`                | Publica      | Login con session token de Shopify  |
| POST   | `/integrations/shopify/connect`                   | JWT          | Iniciar flujo OAuth                 |
| POST   | `/integrations/shopify/connect/custom`            | JWT          | Iniciar OAuth con app custom        |
| GET    | `/integrations/shopify/oauth/token`               | Token/Cookie | Obtener token OAuth                 |
| GET    | `/shopify/callback`                               | State+HMAC   | Callback OAuth de Shopify           |
| POST   | `/integrations/shopify/webhook`                   | HMAC         | Webhook de ordenes                  |
| POST   | `/integrations/shopify/webhook/:integration_id`   | HMAC         | Webhook con integration_id en path  |
| POST   | `/integrations/shopify/webhooks/compliance`       | HMAC         | Webhooks GDPR/CCPA unificado        |
| POST   | `/integrations/shopify/webhooks/customers/data_request` | HMAC   | Solicitud de datos del cliente      |
| POST   | `/integrations/shopify/webhooks/customers/redact`  | HMAC        | Borrado de datos del cliente        |
| POST   | `/integrations/shopify/webhooks/shop/redact`       | HMAC        | Borrado de datos de la tienda       |

### MercadoLibre (`/meli`)

| Metodo | Ruta                  | Auth    | Descripcion                          |
|--------|-----------------------|---------|--------------------------------------|
| POST   | `/meli/notifications` | Publica | Notificaciones IPN (orders_v2, payments) |

### WooCommerce (`/woocommerce`)

| Metodo | Ruta                    | Auth    | Descripcion                          |
|--------|-------------------------|---------|--------------------------------------|
| POST   | `/woocommerce/webhook`  | Publica | Webhook de ordenes (HMAC opcional)   |

### VTEX (`/vtex`)

| Metodo | Ruta             | Auth    | Descripcion                          |
|--------|------------------|---------|--------------------------------------|
| POST   | `/vtex/webhook`  | Publica | Webhook Hook v1 (cambios de estado)  |

### Tiendanube (`/tiendanube`)

| Metodo | Ruta                   | Auth    | Descripcion                         |
|--------|------------------------|---------|-------------------------------------|
| POST   | `/tiendanube/webhook`  | Publica | Webhook (esqueleto)                 |

### Magento (`/magento`)

| Metodo | Ruta               | Auth    | Descripcion                         |
|--------|--------------------|---------|-------------------------------------|
| POST   | `/magento/webhook`  | Publica | Webhook (esqueleto)                |

### Amazon (`/amazon`)

| Metodo | Ruta                    | Auth    | Descripcion                         |
|--------|-------------------------|---------|-------------------------------------|
| POST   | `/amazon/notification`  | Publica | Notificacion SQS/SNS (esqueleto)   |

### Falabella (`/falabella`)

| Metodo | Ruta                  | Auth    | Descripcion                         |
|--------|-----------------------|---------|-------------------------------------|
| POST   | `/falabella/webhook`  | Publica | Webhook (esqueleto)                |

### Exito (`/exito`)

| Metodo | Ruta              | Auth    | Descripcion                         |
|--------|-------------------|---------|-------------------------------------|
| POST   | `/exito/webhook`  | Publica | Webhook (esqueleto)                |

---

## Cola de RabbitMQ

Todos los proveedores publican a una unica cola canonica:

| Aspecto   | Valor                              |
|-----------|------------------------------------|
| Cola      | `probability.orders.canonical`     |
| Publisher | Cada sub-modulo de ecommerce       |
| Consumer  | Modulo `orders` (`services/modules/orders/internal/infra/primary/queue/consumer.go`) |

Adicionalmente, Shopify publica eventos de sincronizacion al exchange de eventos de RabbitMQ usando `EventEnvelope` con categoria `"integration"`.

---

## Integracion con Otros Modulos

### Core de Integraciones (`integrations/core`)

Cada proveedor se registra en el core mediante `integrationCore.RegisterIntegration(typeID, provider)`. El core proporciona:

- **IIntegrationCore**: Registro de proveedores, gestion de integraciones, observadores de eventos.
- **IIntegrationContract**: Interfaz que cada proveedor implementa (`TestConnection`, `SyncOrdersByIntegrationID`, `GetWebhookURL`).
- **BaseIntegration**: Implementacion base con metodos opcionales (retornan `ErrNotSupported`).
- Lectura y descifrado de credenciales (`DecryptCredential`).
- Actualizacion de configuracion (`UpdateIntegrationConfig`).

Shopify ademas usa `OnIntegrationCreated` para crear webhooks automaticamente al crear una integracion.

### Modulo Orders (`modules/orders`)

Consume la cola `probability.orders.canonical` y persiste las ordenes. No hay dependencia directa de codigo; la comunicacion es exclusivamente via RabbitMQ.

### Modulo Invoicing (`integrations/invoicing`)

Las ordenes procesadas por `orders` pueden disparar facturacion automatica a traves de la cola `probability.orders.to_invoicing`. Este flujo es downstream y no tiene relacion directa con el modulo ecommerce.

---

## Contrato de Proveedor

Todo proveedor implementa `IIntegrationContract` (definido en `integrations/core`). Los metodos no soportados se heredan via `BaseIntegration`.

### Estado de implementacion por proveedor

| Metodo                           | Shopify | MeLi | WooCommerce | VTEX | Tiendanube | Magento | Amazon | Falabella | Exito |
|----------------------------------|---------|------|-------------|------|------------|---------|--------|-----------|-------|
| `TestConnection`                 | SI      | SI   | SI          | SI   | SI         | SI      | TODO   | SI        | SI    |
| `SyncOrdersByIntegrationID`      | SI      | SI   | SI          | SI   | TODO       | TODO    | TODO   | TODO      | TODO  |
| `SyncOrdersByIntegrationIDWithParams` | -  | SI   | SI          | SI   | TODO       | TODO    | TODO   | TODO      | TODO  |
| `GetWebhookURL`                  | SI      | SI   | base        | base | base       | base    | base   | base      | base  |
| `HandleWebhook / HandleNotification` | SI | SI   | SI          | SI   | TODO       | TODO    | TODO   | TODO      | TODO  |
| `ListWebhooks`                   | SI      | N/A  | TODO        | TODO | TODO       | TODO    | TODO   | TODO      | TODO  |
| `CreateWebhook`                  | SI      | N/A  | TODO        | TODO | TODO       | TODO    | TODO   | TODO      | TODO  |
| `DeleteWebhook`                  | SI      | N/A  | TODO        | TODO | TODO       | TODO    | TODO   | TODO      | TODO  |

---

## Arquitectura Hexagonal por Sub-modulo

Cada proveedor sigue la misma estructura:

```
proveedor/
  bundle.go                                  # Punto de entrada, wiring de dependencias
  internal/
    domain/
      entities.go                            # Entidades de dominio (sin tags)
      ports.go                               # Interfaces: Client, IntegrationService, OrderPublisher
      errors.go                              # Errores de dominio
      constants.go                           # Constantes de estados (solo algunos)
      dtos.go                                # DTOs de API externa (solo algunos)
      query_params.go                        # Parametros de consulta (solo algunos)
    app/
      usecases/
        constructor.go                       # Interfaz IXxxUseCase y constructor New()
        test_connection.go                   # Verificacion de credenciales
        sync_orders.go                       # Sincronizacion masiva (proveedores completos)
        process_webhook.go                   # Procesamiento de webhook en tiempo real
        mapper/
          order_mapper.go                    # Entidad proveedor -> ProbabilityOrderDTO
    infra/
      primary/
        handlers/
          constructor.go                     # Interfaz IHandler y constructor
          handle_webhook.go                  # Handler HTTP del webhook/notificacion
          router.go                          # Registro de rutas (Shopify tiene rutas mas extensas)
      secondary/
        client/
          constructor.go                     # Cliente HTTP hacia la API del proveedor
          get_order.go / get_orders.go       # Llamadas a la API externa
          response/                          # DTOs de respuesta con tags JSON
        core/
          core.go                            # Adaptador IIntegrationContract
          integration_service.go             # Adaptador IIntegrationService -> core
        queue/
          rabbitmq_publisher.go              # Publisher a probability.orders.canonical
          noop_publisher.go                  # Publisher no-op (fallback sin RabbitMQ)
          mapper/
            canonical_order_mapper.go        # Dominio -> serializable (con tags JSON)
          request/
            canonical_order_dto.go           # DTO serializable con etiquetas JSON
```

### Capas

- **Domain**: Entidades puras sin tags, interfaces (ports), errores. Sin dependencias externas (solo `context`, `time`, `uuid`).
- **Application (app/usecases)**: Logica de negocio. Depende solo de domain. Contiene los mappers de orden proveedor a DTO canonico.
- **Infrastructure (infra)**: Implementaciones concretas. Handlers HTTP (primary), clientes HTTP, adaptadores de core, publishers de RabbitMQ (secondary).

---

## Credenciales por Proveedor

Las credenciales se almacenan cifradas en la tabla `integrations` y se acceden via `DecryptCredential` del core.

| Proveedor    | Campos config            | Campos credential                                   |
|-------------|--------------------------|------------------------------------------------------|
| Shopify (1) | -                        | `shop_domain`, `access_token`, `client_secret`       |
| MeLi (3)    | -                        | `access_token`, `refresh_token`, `app_id`, `client_secret` |
| WooCommerce (4) | `store_url`           | `consumer_key`, `consumer_secret`                    |
| VTEX (16)   | `store_url`              | `api_key`, `api_token`                               |
| Tiendanube (17) | `store_url`           | `access_token`                                       |
| Magento (18) | `store_url`             | `access_token`                                       |
| Amazon (19) | `seller_id`              | `refresh_token`, `client_id`, `client_secret`        |
| Falabella (20) | `user_id`              | `api_key`                                            |
| Exito (21)  | `seller_id`              | `api_key`                                            |

MercadoLibre implementa renovacion automatica de tokens (`EnsureValidToken` / `RefreshToken`) ya que el access_token tiene expiracion.

---

## Patron de Inicializacion

El `ecommerce/bundle.go` es el unico punto de entrada. Cada proveedor expone su `New()` que retorna un `IIntegrationContract`; el bundle padre hace el `RegisterIntegration`:

```go
func New(router, logger, config, rabbitMQ, database, integrationCore) {
    // Shopify (type_id=1) -- se auto-registra (incluye OnIntegrationCreated)
    shopify.New(router, logger, config, integrationCore, rabbitMQ)

    // MercadoLibre (type_id=3)
    meliProvider := meli.New(router, logger, config, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeMercadoLibre, meliProvider)

    // WooCommerce (type_id=4)
    wooProvider := woocommerce.New(...)
    integrationCore.RegisterIntegration(core.IntegrationTypeWoocommerce, wooProvider)

    // ... mismo patron para VTEX, Tiendanube, Magento, Amazon, Falabella, Exito
}
```

Shopify maneja su propio registro internamente porque tambien configura un observer `OnIntegrationCreated` que crea webhooks automaticamente al activar una integracion.

---

## Notas Tecnicas

- Los webhooks responden `200 OK` inmediatamente y procesan la orden en una goroutine aparte. Esto es necesario porque las plataformas externas tienen timeouts cortos y reenvian la notificacion si no reciben respuesta rapida.
- Cada publisher tiene un fallback `NoOpPublisher` que se usa cuando RabbitMQ no esta disponible, evitando panics.
- Las entidades de dominio no tienen etiquetas JSON/GORM. La serializacion se hace en la capa de infraestructura con DTOs dedicados y mappers explicitos.
- Shopify es el proveedor mas maduro: OAuth completo, gestion automatica de webhooks, compliance GDPR/CCPA, y procesamiento diferenciado por tipo de evento (`orders/create`, `orders/paid`, `orders/updated`, `orders/cancelled`, `orders/fulfilled`, `orders/partially_fulfilled`).
- MercadoLibre usa IPN (Instant Payment Notification) en lugar de webhooks clasicos, y requiere renovacion automatica de tokens OAuth.

---

## Agregar un Nuevo Proveedor

1. Crear carpeta `ecommerce/<proveedor>/` con la estructura hexagonal estandar (copiar de un esqueleto existente como `tiendanube/`).
2. Implementar `IIntegrationContract` (embeber `BaseIntegration` y sobrescribir los metodos soportados).
3. Implementar el mapper a `canonical.ProbabilityOrderDTO` en `infra/secondary/queue/mapper/`.
4. Definir el `type_id` en `integrations/core/internal/domain/type_codes.go`.
5. Registrar en `ecommerce/bundle.go`:
   ```go
   miProvider := miproveedor.New(router, logger, config, rabbitMQ, integrationCore)
   integrationCore.RegisterIntegration(core.IntegrationTypeMiProveedor, miProvider)
   ```
