# Modulo de Facturacion Electronica (Invoicing)

Modulo que agrupa todos los proveedores de facturacion electronica soportados por Probability. Cada proveedor opera como un sub-modulo independiente que consume solicitudes de facturacion desde RabbitMQ, genera la factura en el sistema externo del proveedor y publica el resultado de vuelta a la cola de respuestas.

**Ruta:** `services/integrations/invoicing/`

---

## Proposito

- Facturar ordenes de forma automatica o manual contra proveedores de facturacion electronica colombianos.
- Abstraer las diferencias de API, autenticacion y formato entre proveedores detras de una interfaz unificada (`IIntegrationContract`).
- Enrutar dinamicamente las solicitudes al proveedor correcto segun el `integration_type_id` de la integracion configurada para cada negocio.

---

## Proveedores

| Proveedor    | type_id | Cola de entrada                    | Estado    | Operaciones implementadas                                          |
|--------------|---------|------------------------------------|-----------|--------------------------------------------------------------------|
| Softpymes    | 5       | `invoicing.softpymes.requests`     | Completo  | CreateInvoice, CreateCreditNote, GetDocument, ListDocuments, Compare |
| Factus       | 7       | `invoicing.factus.requests`        | Completo  | CreateInvoice, ListBills, GetBill, TestConnection                  |
| Siigo        | 8       | `invoicing.siigo.requests`         | Completo  | CreateInvoice, GetCustomer, CreateCustomer, ListInvoices, TestConnection |
| Alegra       | 9       | `invoicing.alegra.requests`        | Esqueleto | CreateInvoice (stub), TestConnection (stub)                        |
| World Office | 10      | `invoicing.world_office.requests`  | Esqueleto | CreateInvoice (stub), TestConnection (stub)                        |
| Helisa       | 11      | `invoicing.helisa.requests`        | Esqueleto | CreateInvoice (stub), TestConnection (stub)                        |

Los proveedores marcados como "Esqueleto" tienen la estructura hexagonal completa y los consumers de RabbitMQ activos, pero la logica de llamada a la API real esta pendiente de implementacion.

---

## Entidades y DTOs principales

Cada proveedor define sus propios DTOs internos (regla de aislamiento de modulos). Las estructuras comunes son:

- **OrderEventMessage / OrderSnapshot / OrderItemSnapshot** -- Replica local del payload de ordenes que llega por RabbitMQ. Contiene datos del cliente, items, totales, moneda e integracion.
- **InvoiceRequestMessage** -- Mensaje recibido desde `modules/invoicing`, con `invoice_id`, `provider`, `operation`, `invoice_data` y `correlation_id`.
- **InvoiceResponseMessage** -- Respuesta publicada de vuelta con estado (`success`/`error`), numero de factura, external_id, document_json, CUFE, audit data y tiempos de procesamiento.
- **CompareResponseMessage** -- Respuesta del flujo de comparacion: lista de documentos del proveedor en un rango de fechas.
- **InvoicingConfig / FilterConfig** -- Configuracion de facturacion automatica por integracion, con filtros por monto, estado de pago, metodos de pago, tipos de cliente, regiones, etc.
- **Credentials** -- Cada proveedor tiene su propio struct de credenciales (api_key/api_secret para Softpymes, client_id/client_secret/username/password para Factus, username/access_key/account_id/partner_id para Siigo).

---

## Colas RabbitMQ

| Cola                              | Direccion          | Descripcion                                      |
|-----------------------------------|--------------------|--------------------------------------------------|
| `invoicing.requests`              | Entrada (unificada)| Todas las solicitudes de facturacion              |
| `invoicing.responses`             | Salida (unificada) | Respuestas de todos los proveedores               |
| `invoicing.events`                | Salida             | Eventos de facturacion (creada, fallida, etc.)    |
| `invoicing.bulk.create`           | Entrada            | Trabajos de facturacion masiva                    |
| `invoicing.softpymes.requests`    | Entrada proveedor  | Consumer de Softpymes                             |
| `invoicing.factus.requests`       | Entrada proveedor  | Consumer de Factus                                |
| `invoicing.siigo.requests`        | Entrada proveedor  | Consumer de Siigo                                 |
| `invoicing.alegra.requests`       | Entrada proveedor  | Consumer de Alegra                                |
| `invoicing.world_office.requests` | Entrada proveedor  | Consumer de World Office                          |
| `invoicing.helisa.requests`       | Entrada proveedor  | Consumer de Helisa                                |

### Formato del mensaje de solicitud

```json
{
  "invoice_id": 42,
  "provider": "factus",
  "operation": "create",
  "correlation_id": "uuid-v4",
  "timestamp": "2026-02-23T10:00:00Z",
  "invoice_data": {
    "integration_id": 5,
    "customer": { "name": "...", "email": "...", "dni": "..." },
    "items": [{ "sku": "...", "name": "...", "quantity": 1, "unit_price": 50000 }],
    "total": 150000,
    "subtotal": 126050,
    "tax": 23950,
    "config": { "referer": "..." }
  }
}
```

Operaciones soportadas en el campo `operation`: `create`, `retry`, `compare`.

---

## Enrutamiento dinamico de proveedores (resolveProvider)

El modulo `modules/invoicing` (no este modulo) es quien determina a que proveedor enviar una factura. El flujo es:

1. `modules/invoicing` recibe la solicitud de facturar una orden.
2. Llama a `resolveProvider(ctx, integrationID)` que consulta `integrations.integration_type_id` en la base de datos.
3. Mapea el `type_id` a un string de proveedor:

| type_id | Proveedor   |
|---------|-------------|
| 5       | `softpymes` |
| 7       | `factus`    |
| 8       | `siigo`     |
| default | `softpymes` (fallback) |

4. Publica el mensaje a la cola unificada `invoicing.requests` con el campo `provider` ya resuelto.
5. El **router** de este modulo (`invoicing/router/bundle.go`) consume `invoicing.requests`, lee el campo `provider` del encabezado y reenvia el mensaje integro (sin transformacion) a la cola del proveedor correspondiente (`invoicing.<provider>.requests`).
6. El consumer del proveedor recibe el mensaje, descifra credenciales via `integrationCore.DecryptCredential()`, llama a la API externa y publica la respuesta a `invoicing.responses`.

```
modules/invoicing
    |
    |  resolveProvider() --> type_id --> provider string
    |
    +-- PUBLISH --> invoicing.requests
                        |
                   invoicing/router
                        |
           +------------+------------+--- ...
           v            v            v
  invoicing.softpymes  invoicing.factus  invoicing.siigo
  .requests            .requests         .requests
           |            |            |
    consumer        consumer      consumer
           |            |            |
    Softpymes API   Factus API    Siigo API
           |            |            |
           +----------- + -----------+
                        |
                   invoicing.responses
                        |
                   modules/invoicing
                   (response consumer)
```

---

## Autenticacion por proveedor

Cada proveedor gestiona su propia autenticacion con cache en memoria (`token_cache.go`).

### Softpymes -- Bearer Token
- Credenciales: `api_key`, `api_secret`, `referer`
- URL dinamica desde `integration_types.base_url` / `base_url_test`
- Cache de token en memoria

### Factus -- OAuth2 Password Grant
- Endpoint auth: `POST /oauth/token` (form-data)
- `grant_type=password` para login inicial, `grant_type=refresh_token` para renovar
- TTL access token: 10 min (buffer 2 min)
- TTL refresh token: 1h (buffer 5 min)
- Credenciales: `client_id`, `client_secret`, `username`, `password`, `base_url`
- Estrategia: Cache -> Refresh -> Login completo

### Siigo -- Bearer Token (sin refresh)
- Endpoint auth: `POST /v1/auth` (JSON body)
- TTL access token: 24h (buffer 30 min)
- Headers especiales: `Authorization: <account_id>`, `Partner-Id: <partner_id>` para auth; `Authorization: Bearer <token>` para requests post-auth
- Credenciales: `username`, `access_key`, `account_id`, `partner_id`, `base_url`
- Estrategia: Cache -> Login completo (no tiene refresh token)

### Alegra, World Office, Helisa
- Pendiente de implementacion (esqueletos con stubs).

---

## Estructura de directorios

```
invoicing/
+-- bundle.go                 # Orquestador: inicializa todos los proveedores + router
+-- router/
|   +-- bundle.go             # Router centralizado: invoicing.requests -> proveedor
+-- softpymes/
|   +-- bundle.go
|   +-- internal/
|       +-- domain/
|       |   +-- dtos/         # Credentials, CreateInvoiceRequest, AuditData
|       |   +-- entities/     # InvoicingConfig, FilterConfig, catalogos DIAN
|       |   +-- ports/        # ISoftpymesClient, IInvoiceUseCase, OrderEventMessage
|       |   +-- errors/
|       +-- app/
|       |   +-- constructor.go
|       |   +-- process_order_for_invoicing.go
|       |   +-- test_connection.go
|       +-- infra/
|           +-- primary/consumer/          # Consumer RabbitMQ (create, retry, compare)
|           +-- secondary/
|               +-- client/                # HTTP client (auth, invoice, credit_note, get/list docs)
|               +-- core/                  # Adaptador -> IIntegrationContract
|               +-- queue/                 # Publisher de responses
+-- factus/    (misma estructura)
+-- siigo/     (misma estructura, incluye create_customer y get_customer)
+-- alegra/    (misma estructura -- esqueleto)
+-- world_office/  (misma estructura -- esqueleto)
+-- helisa/    (misma estructura -- esqueleto)
```

---

## Integracion con otros modulos

### modules/invoicing (modulo de negocio)
- **Publicador**: `modules/invoicing` crea la entidad `Invoice` en BD, resuelve el proveedor y publica a `invoicing.requests`.
- **Consumidor de respuestas**: `modules/invoicing` consume `invoicing.responses` para actualizar el estado de la factura (success/failed), guardar el `document_json`, `invoice_number`, `external_id`, CUFE, audit data, y actualizar el `invoice_sync_log`.

### integrations/core
- Provee `IIntegrationCore` / `IIntegrationService` para:
  - `GetIntegrationByID()` -- obtener metadata de la integracion (config, base_url, is_testing).
  - `DecryptCredential()` -- descifrar credenciales almacenadas (AES).
  - `RegisterIntegration()` -- registrar cada proveedor como implementacion de `IIntegrationContract`.
  - `TestConnection()` -- probar conectividad contra la API del proveedor.

### Contrato global (IIntegrationContract)
- Cada proveedor implementa `IIntegrationContract` a traves de un adaptador en `infra/secondary/core/core.go`.
- Embede `BaseIntegration` (que retorna `ErrNotSupported` para operaciones no soportadas como SyncOrders, Webhooks, UpdateInventory).
- Solo sobrescribe `TestConnection()`.

---

## Arquitectura (capas hexagonales)

Cada sub-modulo de proveedor sigue la arquitectura hexagonal del proyecto:

```
Domain (nucleo)
  +-- entities/    Entidades puras sin tags ni dependencias externas
  +-- dtos/        Structs de datos tipados (sin tags en domain)
  +-- ports/       Interfaces: ISoftpymesClient, IInvoiceUseCase, etc.
  +-- errors/      Errores de dominio

Application
  +-- app/         Use cases: TestConnection, ProcessOrderForInvoicing
                   Dependen solo de ports (interfaces)

Infrastructure
  +-- primary/
  |   +-- consumer/   Consumer RabbitMQ (adaptador primario, driving)
  +-- secondary/
      +-- client/     Cliente HTTP al API del proveedor (adaptador secundario, driven)
      +-- core/       Adaptador -> IIntegrationContract (puente con integrations/core)
      +-- queue/      Publisher de responses a invoicing.responses
```

La unica dependencia externa permitida entre sub-modulos es hacia `integrations/core` para acceso a credenciales y metadata de integraciones. Los sub-modulos de proveedores NO se importan entre si.

---

## Agregar un nuevo proveedor

1. Crear carpeta `invoicing/<proveedor>/` copiando la estructura de un proveedor esqueleto (ej: `alegra/`).
2. Definir el `type_id` en `core/internal/domain/type_codes.go` y re-exportar en `core/bundle.go`.
3. Declarar la cola en `shared/rabbitmq/queues.go` (`QueueInvoicing<Proveedor>Requests`).
4. Agregar la cola y el case en `invoicing/router/bundle.go` (`getProviderQueue`).
5. Agregar el case en `modules/invoicing/internal/app/create_invoice.go` (`resolveProvider`).
6. Registrar en `invoicing/bundle.go`:
   ```go
   miBundle := mi_proveedor.New(logger, rabbitMQ, integrationCore)
   integrationCore.RegisterIntegration(core.IntegrationTypeMiProveedor, miBundle)
   ```
7. Implementar el cliente HTTP (`client/`) con autenticacion, token cache y las operaciones de la API.
8. Implementar el use case (`app/`) con `TestConnection` y `ProcessOrderForInvoicing`.
9. Implementar el consumer (`consumer/`) para deserializar el mensaje, obtener credenciales, llamar al use case/client y publicar la respuesta.

El router se inicializa al final del `bundle.go` principal para garantizar que todas las colas de proveedores ya esten declaradas cuando comience a consumir.

---

Ultima actualizacion: 2026-03-01
