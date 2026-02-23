# Integración Softpymes - Facturación Electrónica

Módulo de integración con **Softpymes** para emisión de facturas electrónicas en Colombia. Opera de forma **completamente asíncrona y sin estado local**, apoyándose en RabbitMQ, Redis e IntegrationCore.

---

## Estructura del Módulo

```
softpymes/
├── bundle.go                                     # Punto de entrada: inicializa consumers y cliente
└── internal/
    ├── domain/
    │   ├── entities/
    │   │   ├── catalog.go                        # Constantes DIAN (tipos de doc, monedas, etc.)
    │   │   └── config.go                         # InvoicingConfig (réplica local)
    │   ├── dtos/
    │   │   └── invoice_types.go                  # DTOs tipados (CreateInvoiceRequest, Result, etc.)
    │   ├── ports/
    │   │   └── ports.go                          # Interfaces: ISoftpymesClient, IInvoiceUseCase
    │   └── errors/
    │       └── errors.go                         # Errores de dominio
    ├── app/
    │   ├── constructor.go                        # Constructor del use case
    │   └── process_order_for_invoicing.go        # Use case: facturación automática por evento
    └── infra/
        ├── primary/
        │   └── consumer/
        │       ├── invoice_request_consumer.go   # Escucha invoicing.softpymes.requests
        │       └── order_consumer.go             # Escucha orders.events.invoicing
        └── secondary/
            ├── client/
            │   ├── client.go                     # Cliente HTTP base
            │   ├── auth.go                       # Autenticación y test de conexión
            │   ├── invoice.go                    # CreateInvoice()
            │   ├── credit_note.go                # CreateCreditNote()
            │   ├── get_document.go               # GetDocumentByNumber()
            │   ├── list_documents.go             # ListDocuments()
            │   ├── customer.go                   # ensureCustomerExists(), createCustomer()
            │   └── token_cache.go                # Cache en memoria del Bearer token
            ├── cache/
            │   └── config_cache.go               # Redis cache de InvoicingConfig
            ├── queue/
            │   └── response_publisher.go         # Publica respuestas a invoicing.responses
            └── integration_cache/
                └── client.go                     # Lee metadata y credenciales de IntegrationCore
```

---

## Flujos Principales

### Flujo A — Facturación bajo demanda (desde Invoicing Module)

```
Invoicing Module
    │
    └─▶ [invoicing.softpymes.requests]
            │
            ▼
    InvoiceRequestConsumer.Start()
            │
            ├─ Obtiene integración desde IntegrationCore
            ├─ Desencripta api_key / api_secret
            ├─ Client.CreateInvoice()
            │       ├─ authenticate()               ← POST /oauth/integration/login/
            │       ├─ ensureCustomerExists()        ← GET/POST /app/integration/customer
            │       └─ POST /app/integration/sales_invoice/
            │
            ├─ Espera 3s para procesamiento DIAN
            ├─ GetDocumentByNumber()                ← POST /app/integration/search/documents/
            │
            └─▶ [invoicing.responses]               ← InvoiceResponseMessage
```

### Flujo B — Facturación automática (desde Orders Module)

```
Orders Module
    │
    └─▶ [orders.events.invoicing]
            │
            ▼
    OrderConsumer.Start()
            │
            ▼
    ProcessOrderForInvoicing() [Use Case]
            │
            ├─ ConfigCache (Redis) → fallback IntegrationCore
            ├─ Validar filtros (monto, pago, estado)
            ├─ Verificar duplicado en Redis Hash
            ├─ Obtener credenciales desde integration_cache
            ├─ Client.CreateInvoice()               ← mismo flujo que Flujo A
            └─ Marcar como procesado en Redis (30 días)
```

---

## Autenticación con Softpymes

| Aspecto | Detalle |
|---------|---------|
| Endpoint | `POST /oauth/integration/login/` |
| Body | `{"apiKey": "...", "apiSecret": "..."}` |
| Header requerido | `Referer: <URL del cliente>` |
| Response | `{"accessToken": "...", "expiresInMin": 60, "tokenType": "Bearer"}` |
| Cache | En memoria (TokenCache). Se invalida 5 min antes de vencer o al recibir 401 |

---

## Endpoints de Softpymes Consumidos

| Endpoint | Método | Propósito |
|----------|--------|-----------|
| `/oauth/integration/login/` | POST | Obtener Bearer token |
| `/app/integration/customer` | GET | Buscar cliente por NIT |
| `/app/integration/customer` | POST | Crear cliente (tercero) |
| `/app/integration/sales_invoice/` | POST | Crear factura electrónica |
| `/app/integration/search/documents/` | POST | Listar / buscar documentos |
| `/search/documents/notes/` | POST | Crear nota de crédito |

---

## Queues de RabbitMQ

| Queue | Dirección | Propósito |
|-------|-----------|-----------|
| `invoicing.softpymes.requests` | Entrada | Solicitudes de facturación desde Invoicing Module |
| `orders.events.invoicing` | Entrada | Eventos de órdenes nuevas para facturación automática |
| `invoicing.responses` | Salida | Resultado de facturación para Invoicing Module |

---

## Cacheado en Redis

| Key Pattern | Tipo | TTL | Propósito |
|-------------|------|-----|-----------|
| `probability:invoicing:config:{integration_id}` | String (JSON) | 1 hora | Config de facturación automática |
| `probability:invoices:processed:{order_id}` | Hash | 30 días | Prevenir facturación duplicada |
| `integration:meta:{integration_id}` | String | — | Metadata de la integración (IntegrationCore) |
| `integration:creds:{integration_id}` | String | — | Credenciales encriptadas (IntegrationCore) |

---

## Gestión de Clientes (Terceros)

Antes de crear una factura, el cliente verifica si el tercero existe en Softpymes:

1. `GET /app/integration/customer?identification={nit}` → si existe, usa su `branchCode`
2. Si no existe → `POST /app/integration/customer` con:
   - Tipo: Persona Natural (`thirdType = "N"`)
   - Identificación: Cédula de Ciudadanía (`13`)
   - Nombre dividido en `firstName` / `lastName`
   - Campos obligatorios fijos: `maidenName="."`, `otherName="."`
   - Defaults: `email="noreply@probability.com"`, `cityCode="001"`, `departmentCode="11"`

---

## Configuración Requerida

En `integration.config` de la integración Softpymes:

```json
{
    "referer": "https://empresa.softpymes.com.co",
    "resolution_id": 18000123,
    "branch_code": "001",
    "customer_branch_code": "001",
    "seller_nit": "123456789",
    "default_customer_nit": "999999999",
    "company_nit": "123456789"
}
```

| Campo | Requerido | Descripción |
|-------|-----------|-------------|
| `referer` | ✅ Sí | URL del cliente en Softpymes (header de auth) |
| `resolution_id` | ✅ Sí | ID de resolución DIAN para numeración |
| `branch_code` | No | Sucursal del documento (default: `"001"`) |
| `customer_branch_code` | No | Sucursal del cliente (default: `"001"`) |
| `seller_nit` | No | NIT del vendedor |
| `default_customer_nit` | No | NIT por defecto si el cliente no tiene |
| `company_nit` | No | NIT de la empresa para crear clientes |

---

## Filtros de Facturación Automática

El use case `ProcessOrderForInvoicing` valida los siguientes criterios antes de facturar:

1. Config habilitada (`enabled = true`)
2. Facturación automática activa (`auto_invoice = true`)
3. Orden no procesada previamente (Redis Hash)
4. Monto `>=` `min_amount` (si está configurado)
5. Método de pago en lista permitida (si está configurado)
6. Estado de pago coincide con el requerido (si está configurado)

---

## Mapeo de Datos

### Moneda

| Entrada | Código Softpymes |
|---------|-----------------|
| `COP` / `cop` | `"P"` (Peso) |
| `USD` / `usd` | `"D"` (Dólar) |
| Cualquier otro | `"P"` (default) |

### Items

- `unitCode`: `"UNI"` (Unidades — estándar DIAN)
- `unitValue`: formato string `"%.2f"`
- `quantity`: `float64`
- Fecha del documento: zona horaria Bogotá (UTC-5), formato `YYYY-MM-DD`

---

## Trazabilidad (AuditData)

Cada llamada a la API de Softpymes genera un `AuditData` que se incluye siempre en el resultado (incluso en error):

```go
type AuditData struct {
    RequestURL     string      // Ej: "/app/integration/sales_invoice/"
    RequestPayload interface{} // Payload enviado
    ResponseStatus int         // HTTP status code
    ResponseBody   string      // Body crudo de la respuesta
}
```

---

## Principios de Arquitectura

| Principio | Implementación |
|-----------|----------------|
| Sin base de datos propia | Solo HTTP, RabbitMQ y Redis |
| Aislamiento de módulos | DTOs replicados localmente; no importa de otros módulos |
| Dual Read Pattern | ConfigCache (Redis) primero, fallback a IntegrationCore |
| Idempotencia | Redis Hash previene facturación duplicada |
| Audit trail | AuditData en cada request HTTP |
| Token resilience | TokenCache en memoria con invalidación automática por 401 |

---

## Dependencias Internas

| Módulo | Propósito |
|--------|-----------|
| `services/integrations/core` | Desencriptar credenciales de la integración |
| `shared/httpclient` | Cliente HTTP con reintentos |
| `shared/rabbitmq` | Consumo y publicación de mensajes |
| `shared/redis` | Cache de config e idempotencia |
| `shared/log` | Logger centralizado (zerolog) |
