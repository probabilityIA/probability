# Módulo Factus - Facturación Electrónica

Proveedor colombiano de facturación electrónica con certificación DIAN.
**Integration type_id:** `7` | **Queue:** `invoicing.factus.requests`

> ⚠️ **Corrección aplicada:** El código original tenía `IntegrationTypeFactus = 6` pero en la DB el type_id real de Factus es **7** (id=6 es "Plataforma"). Esto fue corregido en `core/public.go` y en `create_invoice.go`.

---

## ¿Qué se implementó?

Módulo completo de facturación electrónica para Factus con OAuth2 Password Grant + Refresh Token. El módulo:

- Autentica con `POST /oauth/token` usando form-data (OAuth2)
- Renueva el access token automáticamente con refresh token
- Crea facturas electrónicas vía `POST /v1/bills/validate`
- Consulta facturas paginadas vía `GET /v1/bills`
- Obtiene detalle de una factura vía `GET /v1/bills/show/:number`
- Consume mensajes de RabbitMQ desde `invoicing.factus.requests`
- Publica resultados (éxito o error) a `invoicing.responses`
- Cachea configuraciones en Redis (TTL: 1h)
- Lee credenciales desde el cache de IntegrationCore (Redis)

### Estructura de archivos

```
factus/
+-- bundle.go                                                # Factory + IIntegrationContract
+-- internal/
    +-- domain/
    |   +-- dtos/
    |   |   +-- invoice_types.go                            # Credentials, CustomerData, ItemData, CreateInvoiceRequest/Result, AuditData
    |   |   +-- bill_types.go                              # ListBillsParams, Bill, BillDetail, ListBillsResult, BillsPagination
    |   +-- entities/config.go                              # InvoicingConfig + FilterConfig (réplica local)
    |   +-- errors/errors.go                               # ErrAuthFailed, ErrInvoiceCreationFailed, ErrMissingCredentials, ErrTokenExpired
    |   +-- ports/ports.go                                 # IFactusClient, IInvoiceUseCase + estructuras de eventos
    +-- app/
    |   +-- constructor.go                                  # Use case stub
    |   +-- process_order_for_invoicing.go                 # Stub (procesamiento real en consumer)
    +-- infra/
        +-- primary/consumer/
        |   +-- invoice_request_consumer.go                # Consumer RabbitMQ
        +-- secondary/
            +-- client/
            |   +-- client.go                              # Client struct + New() + endpointURL()
            |   +-- auth.go                                # authenticate() / loginWithPassword() / refreshAccessToken() / TestAuthentication()
            |   +-- token_cache.go                         # Cache dual: access_token (10 min) + refresh_token (1h)
            |   +-- create_invoice.go                      # POST /v1/bills/validate
            |   +-- list_bills.go                          # GET /v1/bills (filtros + paginación)
            |   +-- get_bill.go                            # GET /v1/bills/show/:number
            |   +-- request/invoice.go                     # CreateBillBody, CreateBillCustomer, CreateBillItem
            |   +-- response/
            |   |   +-- invoice.go                         # CreateBill, CreatedBill, CreateBillData
            |   |   +-- list_bills.go                      # Bills, Bill, BillsData, BillsPagination
            |   |   +-- get_bill.go                        # GetBillDetail
            |   +-- mappers/
            |       +-- invoice.go                         # BuildCreateBillRequest() + helpers GetConfigString/GetConfigInt
            |       +-- list_bills.go                      # BillsToListResult()
            |       +-- get_bill.go                        # GetBillToDetail()
            |       +-- config.go                          # GetConfigString(), GetConfigInt()
            +-- cache/
            |   +-- config_cache.go                        # Redis cache para InvoicingConfig (TTL 1h)
            +-- integration_cache/
                +-- client.go                              # Lectura del Redis cache de IntegrationCore (meta + credenciales)
```

---

## ¿Cómo funciona?

### Flujo completo

```
modules/invoicing
    +-- CreateInvoice()
        +-- resolveProvider(integrationID)
            +-- type_id=7 -> provider="factus"
            +-- PublishInvoiceRequest -> invoicing.requests

invoicing.core (router)
    +-- handleInvoiceRequest()
        +-- provider="factus" -> invoicing.factus.requests

factus.InvoiceRequestConsumer
    +-- handleInvoiceRequest()
        +-- processCreateInvoice()
            +-- integrationCore.GetIntegrationByID()
            +-- integrationCore.DecryptCredential(client_id, client_secret, username, password)
            +-- factusClient.CreateInvoice()
            |   +-- authenticate()
            |   |   +-- [cache hit]  -> retorna token cacheado
            |   |   +-- [access exp] -> refreshAccessToken() con refresh_token
            |   |   +-- [ambos exp]  -> loginWithPassword() con credenciales
            |   +-- mappers.BuildCreateBillRequest() -> CreateBillBody
            |   +-- POST /v1/bills/validate
            +-- responsePublisher.PublishResponse() -> invoicing.responses
```

### Autenticación OAuth2

| Aspecto | Detalle |
|---------|---------|
| Endpoint | `POST /oauth/token` |
| Formato | **form-data** (no JSON) |
| Grant inicial | `grant_type=password` |
| Grant renovación | `grant_type=refresh_token` |
| Access token TTL | 10 min (buffer 2 min -> efectivo **8 min**) |
| Refresh token TTL | 1h (buffer 5 min -> efectivo **55 min**) |
| Estrategia cache | Dual token: intenta access -> intenta refresh -> login completo |

**Login inicial:**
```
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=password&client_id=...&client_secret=...&username=...&password=...
```

**Renovación:**
```
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&client_id=...&client_secret=...&refresh_token=...
```

**Response:**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "def50...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### Credenciales requeridas

| Campo | Descripción |
|-------|-------------|
| `client_id` | ID del cliente OAuth2 en Factus |
| `client_secret` | Secreto del cliente OAuth2 |
| `username` | Email del usuario Factus |
| `password` | Contraseña del usuario Factus |
| `api_url` | URL base (opcional, default: `https://api.factus.com.co`) |

### Configuración de la integración (`invoice_config`)

| Campo | Tipo | Descripción | Default |
|-------|------|-------------|---------|
| `numbering_range_id` | int | ID del rango de numeración DIAN | **requerido** |
| `reference_code` | string | Código de referencia (usa `order_id` si está vacío) | `order_id` |
| `payment_form` | string | Forma de pago (`"1"`=contado, `"2"`=crédito) | `"1"` |
| `payment_method_code` | string | Código método de pago (`"10"`=efectivo) | `"10"` |
| `operation_type` | int | Tipo de operación (`10`=estándar) | `10` |
| `document` | string | Tipo de documento (`"01"`=factura) | `"01"` |
| `legal_organization_id` | string | Tipo de organización (`"1"`=persona natural, `"2"`=jurídica) | `"2"` |
| `tribute_id` | string | Régimen tributario (`"21"`=no responsable IVA, `"22"`=responsable) | `"21"` |
| `identification_document_id` | string | Tipo de documento del cliente (`"3"`=NIT, `"13"`=CC) | `"3"` |
| `municipality_id` | string | Código DANE del municipio del cliente | `""` |
| `item_scheme_id` | string | Esquema de items | `"1"` |
| `unit_measure_id` | int | Unidad de medida (`70`=unidad) | `70` |
| `standard_code_id` | int | Código estándar del item | `1` |
| `item_tribute_id` | int | ID del tributo del item (IVA) | `1` |
| `default_tax_rate` | string | Tasa de IVA por defecto | `"19.00"` |

### Creación de factura — POST /v1/bills/validate

**Endpoint:** `POST /v1/bills/validate`
**Auth:** `Authorization: Bearer <access_token>`

Body principal:
```json
{
  "numbering_range_id": 123,
  "reference_code": "ORD-456",
  "payment_form": "1",
  "payment_method_code": "10",
  "operation_type": 10,
  "send_email": false,
  "document": "01",
  "customer": {
    "identification": "900123456",
    "names": "Cliente Ejemplo S.A.S.",
    "email": "cliente@ejemplo.com",
    "phone": "3001234567",
    "legal_organization_id": "2",
    "tribute_id": "21",
    "identification_document_id": "3",
    "municipality_id": "11001"
  },
  "items": [
    {
      "scheme_id": "1",
      "code_reference": "SKU-001",
      "name": "Producto Ejemplo",
      "quantity": 2,
      "price": 45000.00,
      "tax_rate": "19.00",
      "unit_measure_id": 70,
      "standard_code_id": 1,
      "is_excluded": 0,
      "tribute_id": 1
    }
  ]
}
```

**Response exitosa:**
```json
{
  "status": "OK",
  "message": "Factura creada exitosamente",
  "data": {
    "bill": {
      "id": 789,
      "number": "SETP990000203",
      "cufe": "abc123...",
      "qr": "https://...",
      "total": "107100.00",
      "status": 1,
      "validated": "2026-02-22T10:30:00Z"
    }
  }
}
```

**Resultado mapeado al dominio:**
- `InvoiceNumber` <- `bill.number` (`"SETP990000203"`)
- `ExternalID` <- `bill.id` (como string)
- `CUFE` <- `bill.cufe`
- `QRCode` <- `bill.qr`
- `Total` <- `bill.total`
- `IssuedAt` <- `bill.validated`

### Cache de configuraciones (Redis)

El módulo tiene un `ConfigCache` que guarda `InvoicingConfig` en Redis:
- **Key:** `probability:invoicing:config:<integration_id>` (prefijo configurable vía `REDIS_INVOICING_CONFIG_PREFIX`)
- **TTL:** 1 hora

> **Nota:** La lectura del cache de configuraciones está implementada pero **no está conectada al consumer**. El consumer actual usa `integrationCore.GetIntegrationByID()` directamente.

---

## ¿Qué falta?

### Crítico (bugs / comportamiento incorrecto)

- [x] ~~`IntegrationTypeFactus = 6`~~ — **Corregido:** el DB tiene type_id=7, el código tenía 6. Ya fue corregido a `IntegrationTypeFactus = 7` y `case 7` en `resolveProvider()`.

### Pendiente de implementar

- [ ] **Notas crédito** — La API de Factus soporta `POST /v1/credit-notes` para anulaciones. No está implementado en el módulo
- [ ] **Descarga de PDF/XML** — `GET /v1/bills/{id}/download-pdf` y `GET /v1/bills/{id}/download-xml`. Útil para reenviar al cliente
- [ ] **Reenvío por email** — `POST /v1/bills/{id}/send-email` para enviar la factura al cliente
- [ ] **Conectar ConfigCache al consumer** — El `ConfigCache` (Redis) está implementado pero el consumer no lo usa; siempre va a DB
- [ ] **Conectar IntegrationCacheClient al consumer** — El cliente de cache de IntegrationCore (`integration_cache/client.go`) existe pero tampoco se usa en el consumer actual
- [ ] **Webhook de validación DIAN** — Factus puede notificar vía webhook cuando la DIAN acepta/rechaza; no hay endpoint receptor

### Mejoras deseables

- [ ] **Retry automático por token expirado (401)** — Si el proveedor retorna 401 durante `CreateInvoice`, limpiar cache y reintentar una vez
- [ ] **Métricas de facturación** — Contador de facturas exitosas/fallidas por integración

---

## Documentación de referencia

- **Portal Factus:** https://factus.com.co
- **API Reference:** https://docs.factus.com.co
- **Autenticación OAuth2:** https://docs.factus.com.co/#autenticacion
- **Crear factura:** https://docs.factus.com.co/#crear-factura-de-venta
- **Listar facturas:** https://docs.factus.com.co/#listar-facturas
- **Notas crédito:** https://docs.factus.com.co/#notas-credito
- **Rangos de numeración:** https://docs.factus.com.co/#rangos-de-numeracion

---

## Variables de entorno

```bash
# URL base de la API de Factus (opcional, default: https://api.factus.com.co)
FACTUS_API_URL=https://api.factus.com.co

# Prefijo de la key en Redis para config cache (opcional)
REDIS_INVOICING_CONFIG_PREFIX=probability:invoicing:config
```

## Nota de DB

```sql
-- Tipo de integración en la BD:
SELECT id, code, name FROM integration_types WHERE code = 'factus';
-- id=7, code="factus", name="Factus"

-- Integraciones Factus configuradas:
SELECT id, name, integration_type_id FROM integrations WHERE integration_type_id = 7;
```

## Tipos de integración en DB (referencia completa)

| id | code | name |
|----|------|------|
| 1 | Shopify | Shopify |
| 2 | Whastap | Whatsapp |
| 3 | Mercado Libre | Mercado Libre |
| 4 | Woocormerce | WooCommerce |
| 5 | softpymes | Softpymes Facturación |
| 6 | platform | Plataforma |
| 7 | factus | Factus |
| 8 | siigo | Siigo |
| 9 | alegra | Alegra |
| 10 | world_office | World Office |
| 11 | helisa | Helisa |
