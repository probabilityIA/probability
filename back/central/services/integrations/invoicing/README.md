# Integrations — Facturación Electrónica (Invoicing)

Módulo que agrupa todos los proveedores de facturación electrónica soportados por Probability. Cada proveedor consume solicitudes de facturación desde RabbitMQ, genera la factura en el sistema externo y publica el resultado de vuelta a la cola de respuestas.

---

## Proveedores

| Proveedor | type_id | Cola de entrada | Estado |
|-----------|---------|-----------------|--------|
| Softpymes | 5 | `invoicing.softpymes.requests` | Completo |
| Factus | 7 | `invoicing.factus.requests` | Completo |
| Siigo | 8 | `invoicing.siigo.requests` | Completo |
| Alegra | 9 | `invoicing.alegra.requests` | Esqueleto |
| World Office | 10 | `invoicing.world_office.requests` | Esqueleto |
| Helisa | 11 | `invoicing.helisa.requests` | Esqueleto |

---

## Arquitectura de colas

El módulo de órdenes (`modules/invoicing`) publica solicitudes de facturación a **una sola cola de entrada unificada**. El `router` las enruta al proveedor correcto según el campo `provider` del mensaje.

```
modules/invoicing
    │
    └── PUBLISH ──► invoicing.requests  (cola unificada)
                           │
                    invoicing/router
                           │
           ┌───────────────┼───────────────┬─────────────────┐
           ▼               ▼               ▼                 ▼
  invoicing.softpymes  invoicing.factus  invoicing.siigo  invoicing.alegra ...
  .requests            .requests         .requests        .requests
           │               │               │
    softpymes/consumer  factus/consumer  siigo/consumer
           │               │               │
    Softpymes API      Factus API       Siigo API
```

El `router` lee el campo `provider` del encabezado del mensaje y reenvía el payload completo (sin transformación) a la cola del proveedor correspondiente.

---

## Estructura

```
invoicing/
├── bundle.go               # Orquestador — inicializa todos los proveedores + router
├── router/
│   └── bundle.go           # Router centralizado: invoicing.requests → proveedor
├── softpymes/
│   ├── bundle.go
│   └── internal/
│       ├── domain/         # Entidades, ports, DTOs, errores
│       ├── app/            # Lógica de negocio (process_order, test_connection)
│       └── infra/
│           ├── primary/consumer/       # Consumer RabbitMQ
│           └── secondary/
│               ├── client/             # Cliente HTTP Softpymes API
│               ├── core/               # Adaptador → IIntegrationContract
│               └── queue/              # Publisher de respuestas
├── factus/   (misma estructura)
├── siigo/    (misma estructura)
├── alegra/   (misma estructura — esqueleto)
├── world_office/  (misma estructura — esqueleto)
└── helisa/   (misma estructura — esqueleto)
```

---

## Flujo de una solicitud de facturación

```
1. modules/invoicing publica en invoicing.requests:
   { "invoice_id": 42, "provider": "factus", "operation": "create", ... }

2. invoicing/router consume y reenvía a invoicing.factus.requests

3. factus/consumer recibe el mensaje:
   - Descifra credenciales del core (DecryptCredential)
   - Construye el request tipado para la API de Factus
   - Llama client.CreateInvoice()
   - Publica resultado a la cola de respuestas (invoicing.results)
```

---

## Autenticación por proveedor

Cada proveedor maneja su propio token con cache en memoria (`token_cache.go`).

### Factus — OAuth2 Password Grant
| Aspecto | Detalle |
|---------|---------|
| Endpoint auth | `POST /oauth/token` (form-data) |
| `grant_type` inicial | `password` |
| `grant_type` renovación | `refresh_token` |
| TTL access token | 10 min (600s) |
| TTL refresh token | 1h (3600s) |
| Estrategia | Cache → Refresh → Login completo |

**Credenciales requeridas:**

| Campo | Descripción |
|-------|-------------|
| `base_url` | URL base de la API de Factus |
| `client_id` | Client ID OAuth |
| `client_secret` | Client Secret OAuth |
| `username` | Usuario de la cuenta Factus |
| `password` | Contraseña de la cuenta Factus |

### Siigo — Bearer Token (sin refresh)
| Aspecto | Detalle |
|---------|---------|
| Endpoint auth | `POST /v1/auth` (JSON body) |
| TTL access token | 24h (86400s) |
| Estrategia | Cache → Login completo (no hay refresh) |
| Header especial | `Authorization: <account_id>`, `Partner-Id: <partner_id>` |

**Credenciales requeridas:**

| Campo | Descripción |
|-------|-------------|
| `base_url` | URL base de la API de Siigo |
| `username` | Usuario de la cuenta |
| `access_key` | Access key (no es la contraseña) |
| `account_id` | ID de la empresa en Siigo |
| `partner_id` | Partner ID asignado por Siigo |

### Softpymes — Bearer Token
| Aspecto | Detalle |
|---------|---------|
| Estrategia | Cache → Login completo |

**Credenciales requeridas:**

| Campo | Descripción |
|-------|-------------|
| `base_url` | URL base de la API de Softpymes |
| `client_id` | Client ID |
| `client_secret` | Client Secret |
| `username` | Usuario |
| `password` | Contraseña |

---

## Colas RabbitMQ

| Cola | Dirección | Descripción |
|------|-----------|-------------|
| `invoicing.requests` | Entrada | Cola unificada — todos los proveedores |
| `invoicing.softpymes.requests` | Entrada proveedor | Consumer de Softpymes |
| `invoicing.factus.requests` | Entrada proveedor | Consumer de Factus |
| `invoicing.siigo.requests` | Entrada proveedor | Consumer de Siigo |
| `invoicing.alegra.requests` | Entrada proveedor | Consumer de Alegra |
| `invoicing.world_office.requests` | Entrada proveedor | Consumer de World Office |
| `invoicing.helisa.requests` | Entrada proveedor | Consumer de Helisa |

### Formato del mensaje de solicitud

```json
{
  "invoice_id": 42,
  "provider": "factus",
  "operation": "create",
  "correlation_id": "uuid-v4",
  "timestamp": "2026-02-23T10:00:00Z",
  "order": {
    "id": "123",
    "order_number": "ORD-001",
    "total_amount": 150000,
    "customer_name": "Juan Pérez",
    "customer_dni": "1234567890",
    "items": [...]
  }
}
```

El campo `provider` determina el enrutamiento. Valores válidos: `softpymes`, `factus`, `siigo`, `alegra`, `world_office`, `helisa`.

---

## Patrón de inicialización

El `invoicing/bundle.go` centraliza el registro de todos los proveedores:

```go
func New(config, logger, rabbitMQ, integrationCore) {
    // Cada New() retorna IIntegrationContract — el bundle hace el RegisterIntegration
    softpymesBundle := softpymes.New(config, logger, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeInvoicing, softpymesBundle)

    factusBundle := factus.New(logger, rabbitMQ, integrationCore)
    integrationCore.RegisterIntegration(core.IntegrationTypeFactus, factusBundle)

    // ... siigo, alegra, world_office, helisa ...

    // Router al final — las colas de proveedores deben estar declaradas primero
    router.New(logger, rabbitMQ)
}
```

> El router se inicializa **al final** para garantizar que las colas de cada proveedor ya están declaradas cuando comienza a consumir `invoicing.requests`.

---

## Agregar un nuevo proveedor

1. Crear carpeta `invoicing/<proveedor>/` con estructura hexagonal:
   ```
   <proveedor>/
   ├── bundle.go
   └── internal/
       ├── domain/
       │   ├── dtos/invoice_types.go
       │   ├── entities/config.go
       │   ├── errors/errors.go
       │   └── ports/ports.go
       ├── app/
       │   ├── constructor.go
       │   ├── process_order_for_invoicing.go
       │   └── test_connection.go
       └── infra/
           ├── primary/consumer/invoice_request_consumer.go
           └── secondary/
               ├── client/
               ├── core/core.go
               └── queue/response_publisher.go
   ```

2. Definir el `type_id` en `integrations/core/internal/domain/type_codes.go`

3. Agregar la cola en `invoicing/router/bundle.go`:
   ```go
   const QueueMiProveedorRequests = "invoicing.mi_proveedor.requests"
   // y en getProviderQueue():
   case "mi_proveedor":
       return QueueMiProveedorRequests
   ```

4. Registrar en `invoicing/bundle.go`:
   ```go
   miBundle := mi_proveedor.New(logger, rabbitMQ, integrationCore)
   integrationCore.RegisterIntegration(core.IntegrationTypeMiProveedor, miBundle)
   ```
