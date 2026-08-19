# Plan: remision de salida de mercancia hacia Siigo

Fecha: 2026-08-19
Modulo: `back/central/services/integrations/invoicing/siigo` + `services/modules/invoicing`
Estado: investigacion cerrada con reserva, plan pendiente de decisiones del negocio

---

## 1. Punto de partida

Hoy NO existe remision en el modulo. El consumer de Siigo
(`internal/infra/primary/consumer/invoice_request_consumer.go:157`) despacha
estas operaciones y ninguna mas:

| Operacion | Endpoint Siigo |
|---|---|
| `create` / `retry` | `POST /v1/invoices` |
| `check_status` | `GET /v1/invoices/{id}` |
| `cancel` | `POST /v1/invoices/{id}/annul` |
| `cash_receipt` | `POST /v1/vouchers` |
| `credit_note` | `POST /v1/credit-notes` |
| `create_journal` | `POST /v1/journals` |
| `compare` / `list_items` / `inventory_sync` | `GET /v1/products`, `/v1/warehouses` |
| `list_bank_accounts` / `list_siigo_warehouses` | catalogos |

El `inventory_sync` es unidireccional Siigo -> Probability. Probability no
escribe ningun documento de salida en Siigo. El unico descargue de inventario
en Siigo ocurre por la factura de venta.

Nota: el modulo `services/modules/siigoreferrals` NO es remisiones, es un
formulario de captacion de leads.

---

## 2. Investigacion del endpoint

### 2.1 Limitacion de esta sesion

La politica de egress de la organizacion bloquea el dominio de Siigo. Verificado:

```
curl https://developers.siigo.com/reference   -> CONNECT tunnel failed, 403
curl https://api.siigo.com/v1/remissions      -> CONNECT tunnel failed, 403
WebFetch siigonube.portaldeclientes.siigo.com -> EGRESS_BLOCKED
WebFetch siigoapi.docs.apiary.io              -> EGRESS_BLOCKED
```

Segun `/root/.ccr/README.md` un 403 del proxy es denegacion de politica: no se
reintenta ni se rodea. Por eso la verificacion directa contra la documentacion
oficial queda pendiente y hay que hacerla desde una maquina sin la restriccion.

### 2.2 Evidencia indirecta recogida

**SDK oficial de Siigo** (`github.com/SiigoSAS/siigo_sdk_javascript`, ultima
actualizacion oct-2023, generado desde su OpenAPI). Expone 18 clases y ~46
operaciones:

```
account-groups, accounts-payable, cost-centers, credit-notes, customers,
available-documents, document-types, asset-groups, fixed-assets,
invoices (+ /annul /pdf /stamp /mail /stamp/errors), journals, payment-types,
price-lists, products, taxes, test-balance-report, users, vouchers, warehouses
```

**No hay recurso de remisiones.** Tampoco lo tienen los clientes de terceros
revisados (`saulmoralespa/siigo-api-php`, `srdorado/siigo-client-php`). El repo
oficial `siigodocs/siigo-api-docs` solo trae ejemplos de conexion, sin catalogo
de recursos.

### 2.3 Por que esto no es concluyente

El SDK esta desactualizado frente a la API real: nuestro propio cliente ya
consume `POST /v1/webhooks`, que no figura en ese listado. "No esta en el SDK"
es evidencia fuerte, no prueba.

### 2.4 Lo que si esta confirmado

La remision existe como documento del producto Siigo: Siigo Nube tiene
"Configurar remisiones" y "Elaborar remisiones"; Siigo Pyme la llama "documento
tipo S - Nota de remision". Siigo tambien soporta el cruce remision -> factura
desde la UI ("Elaborar factura electronica con cruce de remision").

Conclusion: es funcionalidad de producto. La duda es solo si esta expuesta
por API.

---

## 3. Paso 0 (bloqueante): confirmar el endpoint

Antes de escribir una linea de codigo, cualquiera de estas tres:

- **A.** Abrir `developers.siigo.com/reference` desde una maquina sin la
  restriccion de egress y buscar "remision" / "remission" en el listado de
  recursos.
- **B.** Probe autenticado contra la cuenta: `GET /v1/document-types` con el
  token real y ver si aparece un tipo de documento de remision y con que
  `type`. Es un GET, seguro con la cuenta prestada.
- **C.** Escribir a soporteapi@siigo.com preguntando explicitamente por (1)
  endpoint de remisiones y (2) si el `POST /v1/invoices` acepta referenciar
  una remision para el cruce.

Criterio: si no hay endpoint, se cae la Fase 2 en adelante y aplica el Plan B.

---

## 4. Plan A - la API soporta remisiones

Sigue el patron ya probado de `credit_note` y `create_journal`. Referencia
directa: `modules/invoicing/internal/app/create_credit_note.go`.

### Fase 1 - Dominio y persistencia

- `back/migration/shared/models/remission.go`, espejo de `credit_note.go`:
  `OrderID`, `BusinessID`, `IntegrationID`, `InternalNumber` (`RM-...`),
  `ExternalID`, `RemissionNumber`, `WarehouseID`, `Status`
  (`pending|issued|failed|cancelled|invoiced`), `IssuedAt`, `ProviderResponse`
  jsonb, `CreatedByID`.
- Tabla hija `remission_items` o reuso del patron de `invoice_items`.
- Indice unico parcial: una sola remision vigente por orden.
- Migracion `XXX_remissions.go` en `back/migration`, AutoMigrate solo del
  modelo nuevo.

### Fase 2 - Cliente Siigo

- `client/create_remission.go`, `get_remission.go`, `list_remissions.go`
- `client/request/remission.go`, `client/response/remission.go`,
  `client/mappers/remission.go`
- `client/find_existing_remission.go`: mismo mecanismo de idempotencia que
  `find_existing_invoice.go` (marcador `order:<orderID> | #<orderNumber>` en
  `observations`, busqueda por `customer_identification` + rango de fechas
  antes de cada POST).
- `domain/ports/ports.go`: metodos nuevos en `ISiigoClient`.
- `domain/dtos`: `CreateRemissionRequest` / `CreateRemissionResult`.

### Fase 3 - Consumer

- `consumer/process_remission.go`.
- `case "remission"` en el switch de `invoice_request_consumer.go:157`.
- El router (`invoicing/router/bundle.go`) rutea por provider, no por
  operacion: no requiere cambio, pero hay que verificarlo.

### Fase 4 - Modulo invoicing

- `modules/invoicing/internal/domain/dtos/operation_remission.go` con
  `OperationRemission = "remission"`.
- `modules/invoicing/internal/app/create_remission.go`: valida la orden,
  resuelve provider (solo Siigo; el resto -> error de operacion no soportada),
  persiste la remision en `pending`, crea `InvoiceSyncLog`, publica a
  `invoicing.requests`.
- `handle_remission_response.go` en el consumer de respuestas.
- Handler + rutas: `POST /api/v1/orders/:id/remissions`, `GET` de listado con
  paginacion obligatoria.
- Multi-tenant (`.claude/rules/multi-tenant-security.md`): `business_id` del
  token; super admin por query param; validar que la orden y la integracion
  pertenezcan a ese business.

### Fase 5 - Disparador

- Manual primero (accion en el detalle de la orden).
- Automatico despues: config `remission_on_status` en `invoice_config`
  (vacio = apagado), consumida en el `order_consumer.go` de invoicing, que ya
  escucha la cola `orders.events.invoicing`. Estado natural:
  `ready_to_ship`.
- Idempotencia en dos capas: indice unico en BD + `find_existing_remission`
  en el cliente.

### Fase 6 - El cruce remision -> factura (lo critico)

Si la remision descarga inventario en Siigo y despues la factura lo descarga
otra vez, el stock queda mal, y como nuestro `inventory_sync` es Siigo ->
Probability, el error se propaga hacia aca y hacia los canales.

Hay que resolver con Siigo si el `POST /v1/invoices` acepta referenciar la
remision. **Si por API no se puede cruzar, no se debe emitir remision + factura
para la misma orden de forma automatica: o una o la otra.** Esta decision
condiciona la Fase 5.

### Fase 7 - Frontend

- `front/central/src/services/integrations/invoicing/siigo/domain/types.ts`:
  hoy `SiigoConfig` esta vacio (`{}`) y toda la config de Siigo se edita a mano
  en BD. Agregar `remission_document_id`, `remission_warehouse_id`,
  `remission_on_status` arrastra el pendiente 2 de
  `.claude/alerts/siigo-pendientes.md`.
- Accion en el detalle de la orden y estado visible de la remision.

### Fase 8 - Testing E2E

- `.claude/testing/siigo/back/CU-NN-remision-salida.md` segun
  `.claude/rules/testing.md`.
- **Bloqueante conocido:** la cuenta Siigo de pruebas es prestada y de SOLO
  LECTURA (`.claude/alerts/siigo-pendientes.md`, punto 3). Sin una cuenta
  propia no se puede validar ningun POST, ni el de remision ni el cruce.

---

## 5. Plan B - la API no soporta remisiones

- **B1. Comprobante contable.** Representar la salida con `POST /v1/journals`,
  que ya esta implementado. Deja el asiento, pero no genera el PDF de remision
  ni mueve el kardex igual que el documento nativo.
- **B2. Remision solo en Probability.** Modelo propio + PDF propio, sin
  enviarla a Siigo. Cubre la necesidad operativa (el papel que acompana la
  mercancia) sin tocar el ERP ni arriesgar doble descargue.
- **B3.** No hacerlo y esperar a que Siigo lo exponga.

Recomendacion: **B2**, y en paralelo pedirle el endpoint a Siigo.

---

## 6. Decisiones pendientes del negocio

1. **Que problema resuelve la remision aqui:** el papel que acompana la
   mercancia, o el descargue de inventario en Siigo antes de facturar?
   Cambia el plan completo (B2 basta para lo primero; lo segundo exige el
   endpoint real).
2. **Se emite ademas de la factura, o en vez de?** Define si hace falta el
   cruce de la Fase 6.
3. **Manual o automatico por estado de orden?**

---

## 7. Riesgos

- Doble descargue de inventario en Siigo si no hay cruce (Fase 6).
- No se puede aceptar en E2E sin cuenta Siigo con permiso de escritura.
- Config de Siigo no editable desde la UI: cualquier campo nuevo nace
  requiriendo edicion manual en BD.

Alertas relacionadas: `.claude/alerts/siigo-pendientes.md`,
`.claude/alerts/siigo-stock-sin-nivel-de-inventario.md`,
`.claude/alerts/inventario-saliente-por-canal.md`.
