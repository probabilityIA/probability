# Módulo Siigo - Facturación Electrónica

Proveedor colombiano de facturación electrónica con certificación DIAN.
**Integration type_id:** `8` | **Queue:** `invoicing.siigo.requests`

---

## ¿Qué se hizo?

Implementación completa del módulo Siigo como tercer proveedor de facturación electrónica, siguiendo el mismo patrón de Factus (type_id=6). El módulo:

- Se autentica con la API de Siigo mediante `POST /v1/auth` con JSON y headers especiales
- Busca o crea clientes en Siigo antes de facturar
- Crea facturas electrónicas vía `POST /v1/invoices`
- Consume mensajes de RabbitMQ desde la cola `invoicing.siigo.requests`
- Publica resultados (éxito o error) a `invoicing.responses`
- Se integra al routing centralizado de `invoicing.core`

### Archivos creados

```
siigo/
+-- bundle.go                                         # Factory + IIntegrationContract
+-- internal/
    +-- domain/
    |   +-- dtos/
    |   |   +-- invoice_types.go                      # Credentials, CustomerData, ItemData, CreateInvoiceRequest/Result, AuditData
    |   |   +-- customer_types.go                     # CustomerResult, CreateCustomerRequest, ListInvoicesParams/Result
    |   +-- entities/config.go                        # InvoicingConfig (réplica local)
    |   +-- errors/errors.go                          # Errores específicos de Siigo
    |   +-- ports/ports.go                            # ISiigoClient, IInvoiceUseCase + eventos
    +-- app/
    |   +-- constructor.go                            # Use case stub
    |   +-- process_order_for_invoicing.go            # Stub (procesamiento real en consumer)
    +-- infra/
        +-- primary/consumer/
        |   +-- invoice_request_consumer.go           # Consumer RabbitMQ
        +-- secondary/
            +-- client/
            |   +-- client.go                         # Client struct + New()
            |   +-- auth.go                           # loginWithCredentials() - POST /v1/auth
            |   +-- token_cache.go                    # Cache access_token 24h
            |   +-- create_invoice.go                 # POST /v1/invoices
            |   +-- get_customer.go                   # GET /v1/customers?identification=xxx
            |   +-- create_customer.go                # POST /v1/customers
            |   +-- list_invoices.go                  # GET /v1/invoices
            |   +-- request/invoice.go                # Structs del body de Siigo
            |   +-- response/
            |   |   +-- invoice.go                    # Structs de respuesta de factura
            |   |   +-- customer.go                   # Structs de respuesta de cliente
            |   +-- mappers/
            |       +-- invoice.go                    # dtos.CreateInvoiceRequest -> request.SiigoInvoice
            |       +-- customer.go                   # response.Customer -> dtos.CustomerResult
            +-- queue/
                +-- response_publisher.go             # Publisher a invoicing.responses
```

### Archivos modificados

| Archivo | Cambio |
|---------|--------|
| `integrations/core/public.go` | `IntegrationTypeSiigo = 8` + case `"siigo"` en `getIntegrationTypeCodeAsInt()` |
| `integrations/invoicing/core/bundle.go` | `QueueSiigoRequests = "invoicing.siigo.requests"` + case en `getProviderQueue()` |
| `integrations/bundle.go` | Import + inicialización y registro del bundle Siigo |
| `modules/invoicing/internal/app/create_invoice.go` | `case 8: return dtos.ProviderSiigo` en `resolveProvider()` |

---

## ¿Cómo funciona?

### Flujo completo

```
modules/invoicing
    +-- CreateInvoice()
        +-- resolveProvider(integrationID)
            +-- type_id=8 -> provider="siigo"
            +-- PublishInvoiceRequest -> invoicing.requests

invoicing.core (router)
    +-- handleInvoiceRequest()
        +-- provider="siigo" -> invoicing.siigo.requests

siigo.InvoiceRequestConsumer
    +-- handleInvoiceRequest()
        +-- processCreateInvoice()
            +-- DecryptCredential(username, access_key, account_id, partner_id)
            +-- siigoClient.CreateInvoice()
            |   +-- authenticate() -> POST /v1/auth
            |   +-- GetCustomerByIdentification() -> GET /v1/customers?identification=xxx
            |   |   +-- [si no existe] CreateCustomer() -> POST /v1/customers
            |   +-- POST /v1/invoices
            +-- responsePublisher.PublishResponse() -> invoicing.responses
```

### Autenticación

Siigo usa un mecanismo diferente a Factus y Softpymes:

| Aspecto | Siigo | Factus |
|---------|-------|--------|
| Endpoint | `POST /v1/auth` | `POST /oauth/token` |
| Body | JSON | form-data |
| Credenciales | `username` + `access_key` | `client_id` + `client_secret` + `username` + `password` |
| Headers especiales | `Authorization: <account_id>`, `Partner-Id: <partner_id>` | ninguno |
| TTL | 24h / 86400s (buffer 30 min) | 10 min access / 1h refresh |
| Refresh | No tiene | Sí (OAuth2) |

**Request de auth:**
```json
POST /v1/auth
Authorization: <account_id>
Partner-Id: <partner_id>
Content-Type: application/json

{
  "username": "user@siigo.com",
  "access_key": "1234567890"
}
```

**Response:**
```json
{
  "access_token": "eyJ...",
  "expires_in": 86400,
  "token_type": "Bearer"
}
```

### Credenciales requeridas

| Campo | Descripción | Header en auth |
|-------|-------------|----------------|
| `username` | Email del usuario Siigo | Body JSON |
| `access_key` | Clave de acceso API | Body JSON |
| `account_id` | Subscription key / Siigo Account ID | `Authorization: <account_id>` |
| `partner_id` | ID del partner integrador | `Partner-Id: <partner_id>` |
| `api_url` | URL base (opcional, default: `https://api.siigo.com`) | — |

### Configuración de la integración (`invoice_config`)

Estos campos deben configurarse en el `invoice_config` de la integración en Siigo:

| Campo | Tipo | Descripción | Ejemplo |
|-------|------|-------------|---------|
| `document_id` | int | ID del tipo de documento en Siigo (FV = Factura de Venta) | `24` |
| `payment_method_id` | int | ID del método de pago en Siigo | `5264` |
| `tax_id` | int | ID del impuesto IVA en Siigo | `1` |
| `customer_id_type` | string | Código del tipo de documento del cliente | `"13"` (CC), `"31"` (NIT) |
| `person_type` | string | Tipo de persona del cliente | `"Person"` o `"Company"` |

> **Nota:** Los IDs exactos (`document_id`, `payment_method_id`, `tax_id`) se obtienen consultando la API de Siigo con las credenciales del cliente. Varían por cuenta.

### Manejo de clientes

A diferencia de Factus y Softpymes, Siigo requiere gestión explícita de clientes:

1. **Buscar** cliente por `customer_dni` -> `GET /v1/customers?identification=<dni>`
2. Si **no existe**, crear el cliente -> `POST /v1/customers`
3. Usar la identificación del cliente en el body de la factura

---

## ¿Qué falta?

### Crítico (bloqueante para producción)

- [ ] **Probar con credenciales reales** — El mapeo del body de factura (`request/invoice.go`) es una aproximación basada en la documentación. Necesita validarse contra una cuenta Siigo real, en especial los campos `document_id`, `payment_method_id` y `tax_id`
- [ ] **Validar respuesta de la API** — Los campos de `CreateInvoiceResponse` (especialmente `id`, `name`, `Errors`) deben verificarse contra respuestas reales de Siigo, ya que la estructura puede variar
- [ ] **Obtener IDs de configuración** — Los valores de `document_id`, `payment_method_id` y `tax_id` son específicos de cada cuenta Siigo. Se deben consultar con:
  - `GET /v1/document-types` -> para `document_id`
  - `GET /v1/payment-types` -> para `payment_method_id`
  - `GET /v1/taxes` -> para `tax_id`

### Mejoras deseables

- [ ] **Endpoints de consulta de catálogos** — Implementar `GetDocumentTypes()`, `GetPaymentTypes()`, `GetTaxes()` para ayudar al usuario a configurar la integración desde el frontend
- [ ] **Manejo de notas crédito** — `POST /v1/credit-notes` para anulaciones de facturas
- [ ] **Webhook de estado DIAN** — Siigo notifica cuando la DIAN acepta/rechaza la factura; implementar endpoint receptor
- [ ] **CUFE y QR en respuesta** — Verificar que `metadata.cufe` y `metadata.qr` vengan correctamente en la respuesta de creación
- [ ] **Factura en PDF** — Siigo puede retornar URL de PDF; agregar `pdf_url` al resultado
- [ ] **Paginación real en `ListInvoices()`** — Verificar que los query params `page` y `page_size` sean los correctos en la API de Siigo

---

## Documentación de referencia

### Donde esta la documentacion real

La fuente de verdad de la API es el **API Blueprint** de Siigo. Es un solo
archivo con TODOS los recursos, endpoints, payloads, enums y codigos de error:

- **Apiary (render oficial):** https://siigoapi.docs.apiary.io/
- **Blueprint crudo (`siigoapi.apib`, ~6.000 lineas), grepeable:**
  https://github.com/jdlar1/siigo-mcp/blob/master/siigoapi.apib

El crudo es mas util para responder "existe el endpoint X?" en 10 segundos:

```bash
git clone --depth 1 https://github.com/jdlar1/siigo-mcp /tmp/siigo-doc
grep -n "^## .*\[\(GET\|POST\|PUT\|DELETE\)" /tmp/siigo-doc/siigoapi.apib  # todos los endpoints
grep -in "<lo que busques>" /tmp/siigo-doc/siigoapi.apib
```

Es un mirror de terceros: verificar la fecha del ultimo commit antes de
confiar en el para algo nuevo, y contrastar con Apiary si hay duda.

- **Portal de desarrolladores:** https://developers.siigo.com
- **API Reference:** https://developers.siigo.com/reference
- **Autenticación:** https://developers.siigo.com/docs/autenticacion
- **Facturas:** https://developers.siigo.com/docs/crear-factura-de-venta
- **Clientes:** https://developers.siigo.com/docs/crear-cliente
- **Soporte tecnico API:** soporteapi@siigo.com

### Recursos que expone la API (blueprint del 2026-05-05)

`/products` `/customers` `/quotations` `/invoices` `/purchases`
`/credit-notes` `/vouchers` `/payment-receipts` `/journals`
`/purchase-support-documents` + catalogos (`/document-types`, `/taxes`,
`/payment-types`, `/cost-centers`, `/users`, `/warehouses`, `/price-lists`,
`/fixed-assets`) + reportes + `/webhooks`.

Enum `DocumentType`: `FV` (factura de venta), `RC` (recibo de caja),
`NC` (nota credito), `FC` (factura de compra), `CC` (comprobante contable).
Ademas `DS`, `RP` y `C` como valores de `/v1/document-types?type=`.

### NO existe: remisiones

**La API de Siigo no expone remisiones / salida de mercancia.** Verificado el
2026-08-19 contra el blueprint completo: `grep -ic remis` devuelve 0, y no hay
tipo de documento de remision en el enum. La remision si existe en el producto
Siigo Nube, pero solo por UI.

No volver a investigarlo sin revisar antes
`.claude/plan/siigo-remision-salida.md`, que tiene la evidencia y las
alternativas.

---

## Variables de entorno

```bash
# URL base de la API de Siigo (opcional, default: https://api.siigo.com)
SIIGO_API_URL=https://api.siigo.com
```

## Notas de integración en BD

```sql
-- El tipo de integración ya existe en la DB:
SELECT id, code, name FROM integration_types WHERE id = 8;
-- id=8, code="siigo", name="Siigo"
```
