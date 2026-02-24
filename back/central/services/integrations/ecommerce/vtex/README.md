# Integración VTEX

Módulo de integración con la plataforma de e-commerce **VTEX** para sincronizar órdenes en tiempo real hacia Probability.

**type_id:** `16` (`IntegrationTypeVTEX`)

---

## Documentación VTEX (Links)

| Recurso | URL |
|---------|-----|
| API de Órdenes (OMS) | https://developers.vtex.com/docs/api-reference/orders-api |
| Listar órdenes | https://developers.vtex.com/docs/api-reference/orders-api#get-/api/oms/pvt/orders |
| Detalle de orden | https://developers.vtex.com/docs/api-reference/orders-api#get-/api/oms/pvt/orders/-orderId- |
| Flujo de estados | https://help.vtex.com/en/tutorial/order-flow-and-status--tutorials_196 |
| Feed y Hook de órdenes | https://developers.vtex.com/docs/guides/orders-feed |
| Hook (webhook push) | https://developers.vtex.com/docs/guides/orders-feed#hook |
| Autenticación (API Keys) | https://developers.vtex.com/docs/guides/authentication |
| Catálogo API | https://developers.vtex.com/docs/api-reference/catalog-api |

---

## Autenticación

VTEX usa **API Key + API Token** como headers en cada request. No hay flujo OAuth ni refresh token.

```
X-VTEX-API-AppKey: {appKey}
X-VTEX-API-AppToken: {appToken}
Content-Type: application/json
```

**Base URL:** `https://{accountName}.vtexcommercestable.com.br`

### Credenciales almacenadas

| Campo | Ubicación | Descripción |
|-------|-----------|-------------|
| `store_url` | `config` (texto plano) | URL base de la tienda VTEX |
| `api_key` | `credentials` (encriptado) | X-VTEX-API-AppKey |
| `api_token` | `credentials` (encriptado) | X-VTEX-API-AppToken |

---

## Funcionalidades

### 1. Test Connection

Verifica credenciales con una llamada ligera al catálogo:

```
GET /api/catalog/pvt/product?_offset=0&_limit=1
```

Si retorna `200 OK`, las credenciales son válidas. Si retorna `401/403`, son inválidas.

### 2. Sincronización de órdenes (SyncOrders)

Sincroniza órdenes de los últimos 30 días (configurable). Se ejecuta en **background** (goroutine).

**Flujo:**
1. Obtener integración del core por ID
2. Desencriptar credenciales (`api_key`, `api_token`)
3. Llamar `GET /api/oms/pvt/orders` con paginación (`page`, `per_page`)
4. Por cada orden del listado: `GET /api/oms/pvt/orders/{orderId}` para detalle completo
5. Mapear a `canonical.ProbabilityOrderDTO`
6. Publicar a RabbitMQ (`probability.orders.canonical`)

**Paginación:** VTEX usa `page` + `per_page` (default 15 por página).

**Rate limiting:** 500ms entre páginas para no saturar la API.

**Filtros soportados:**
- `created_at_min` / `created_at_max` → se convierten a `f_creationDate` de VTEX
- `status` → se convierte a `f_status`

### 3. Webhook (Hook v1)

VTEX envía un POST cuando cambia el estado de una orden.

**Endpoint:** `POST /integrations/vtex/webhook`

**Payload ejemplo:**
```json
{
  "Domain": "Fulfillment",
  "OrderId": "v1234567-01",
  "State": "payment-approved",
  "LastState": "payment-pending",
  "LastChange": "2026-02-24T10:30:00.0000000+00:00",
  "CurrentChange": "2026-02-24T10:31:00.0000000+00:00",
  "Origin": {
    "Account": "mitienda",
    "Key": "vtexappkey-mitienda-ABCDEF"
  }
}
```

**Flujo:**
1. Parsear JSON del body
2. Responder `200 OK` inmediatamente
3. En background: buscar integración, desencriptar credenciales
4. Llamar `GET /api/oms/pvt/orders/{orderId}` para detalle completo
5. Mapear a canonical DTO y publicar a RabbitMQ

> **Nota:** El webhook usa `Origin.Account` para identificar la integración. Requiere que el account name de VTEX esté almacenado como `store_id` de la integración.

---

## Mapeo de estados

### Orden (VTEX → Probability)

| Estado VTEX | Estado Probability |
|-------------|-------------------|
| `order-created`, `waiting-for-sellers-confirmation` | `pending` |
| `payment-pending`, `waiting-for-authorization`, `approve-payment` | `pending` |
| `payment-approved`, `authorize-fulfillment`, `window-to-cancel` | `paid` |
| `ready-for-handling`, `start-handling`, `handling` | `processing` |
| `invoice`, `invoiced` | `invoiced` |
| `canceled`, `cancellation-requested`, `cancel` | `cancelled` |

Referencia completa: [Order flow and status](https://help.vtex.com/en/tutorial/order-flow-and-status--tutorials_196)

### Pago

| Contexto orden VTEX | Estado pago Probability |
|---------------------|------------------------|
| `payment-approved` hasta `invoiced` | `paid` |
| `payment-pending`, `waiting-for-authorization` | `pending` |
| `canceled`, `cancel` | `cancelled` |

### Envío

| Contexto orden VTEX | Estado envío Probability |
|---------------------|-------------------------|
| `invoiced` | `shipped` |
| `ready-for-handling`, `handling` | `pending` |
| `canceled` | `cancelled` |

---

## Mapeo de campos (VTEXOrder → ProbabilityOrderDTO)

| Campo VTEX | Campo Canonical | Nota |
|------------|-----------------|------|
| `OrderId` | `ExternalID` | ID del pedido en VTEX |
| `Sequence` | `OrderNumber` | Número secuencial legible |
| `Value / 100` | `TotalAmount` | VTEX usa centavos |
| `TotalItems / 100` | `Subtotal` | Centavos |
| `TotalDiscount / 100` | `Discount` | Se convierte a positivo |
| `TotalFreight / 100` | `ShippingCost` | Centavos |
| `Status` | `OriginalStatus` | Estado original |
| `ClientProfileData.Email` | `CustomerEmail` | |
| `ClientProfileData.FirstName + LastName` | `CustomerName` | Si es corporativo usa `CorporateName` |
| `ClientProfileData.Phone` | `CustomerPhone` | |
| `ClientProfileData.Document` | `CustomerDNI` | CPF/CNPJ/cédula |
| `Items[]` | `OrderItems[]` | Precios en centavos |
| `ShippingData.Address` | `Addresses[]` | GeoCoordinates: `[lng, lat]` |
| `PaymentData.Transactions[].Payments[]` | `Payments[]` | |
| `PackageAttachment.Packages[]` | `Shipments[]` | Incluye tracking |
| `CreationDate` | `OccurredAt` | |

**Conversión de centavos:** Todos los valores monetarios de VTEX son `int` en centavos. Se dividen por 100 para obtener `float64`.

---

## Arquitectura del módulo

```
vtex/
├── bundle.go                          # Ensamblaje del módulo
└── internal/
    ├── domain/
    │   ├── entities.go                # VTEXOrder, VTEXWebhookPayload, etc. (sin tags)
    │   ├── ports.go                   # IVTEXClient, IIntegrationService, OrderPublisher
    │   └── errors.go                  # Errores del dominio
    ├── app/usecases/
    │   ├── constructor.go             # IVTEXUseCase interface + New()
    │   ├── test_connection.go         # Verificación de credenciales
    │   ├── sync_orders.go             # Sincronización paginada en background
    │   ├── process_webhook.go         # Procesamiento de webhooks
    │   └── mapper/
    │       └── order_mapper.go        # VTEXOrder → ProbabilityOrderDTO
    └── infra/
        ├── primary/handlers/
        │   ├── constructor.go         # Handler HTTP + RegisterRoutes
        │   └── handle_webhook.go      # POST /integrations/vtex/webhook
        └── secondary/
            ├── client/
            │   ├── constructor.go     # HTTP client con headers VTEX
            │   ├── get_orders.go      # GET /api/oms/pvt/orders (lista)
            │   ├── get_order_by_id.go # GET /api/oms/pvt/orders/{id} (detalle)
            │   └── response/
            │       └── vtex_order_response.go  # Structs JSON + ToDomain()
            ├── core/
            │   ├── core.go            # IIntegrationContract adapter
            │   └── integration_service.go  # Adapter core → domain
            └── queue/
                ├── rabbitmq_publisher.go    # Publica a RabbitMQ
                ├── noop_publisher.go        # Fallback sin RabbitMQ
                ├── mapper/
                │   └── canonical_order_mapper.go  # Domain → Serializable
                └── request/
                    └── canonical_order_dto.go     # Structs con JSON tags
```

---

## Registro en el sistema

El módulo se registra en `ecommerce/bundle.go`:

```go
vtexProvider := vtex.New(router, logger, config, rabbitMQ, integrationCore)
integrationCore.RegisterIntegration(core.IntegrationTypeVTEX, vtexProvider)
```

Implementa `IIntegrationContract` con:
- `TestConnection` - verificar credenciales
- `GetWebhookURL` - retornar URL del webhook
- `SyncOrdersByIntegrationID` - sync últimos 30 días
- `SyncOrdersByIntegrationIDWithParams` - sync con filtros

---

## Configuración para el usuario

Al crear una integración VTEX desde el frontend:

1. **store_url**: `https://{accountName}.vtexcommercestable.com.br`
2. **api_key**: Obtener desde VTEX Admin > Account Settings > Application Keys
3. **api_token**: Se genera junto con el API Key

### Cómo configurar el webhook en VTEX

1. Ir a VTEX Admin > Orders > Settings > Feed and Hook
2. Seleccionar **Hook**
3. Configurar la URL: `{baseURL}/integrations/vtex/webhook`
4. Seleccionar eventos: orden creada, actualizada, cancelada
5. Guardar

---

## Pendiente / TODO

- [ ] **Webhook lookup por account_name**: Implementar `GetIntegrationByStoreID` en el `IIntegrationService` de VTEX (como tiene MeLi) para buscar la integración por el `Origin.Account` del webhook
- [ ] **Feed v3 polling**: Alternativa al webhook push para entornos donde no se puede configurar webhook
- [ ] **Sincronización de productos**: Actualmente solo sincroniza órdenes
- [ ] **Cancelación de órdenes**: Implementar cancel order vía API de VTEX
- [ ] **Tests de integración**: Tests con credenciales de prueba
