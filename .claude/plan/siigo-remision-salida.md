# Plan: remision de salida de mercancia hacia Siigo

Fecha: 2026-08-19
Modulo: `back/central/services/integrations/invoicing/siigo` + `services/modules/invoicing`
Estado: **RESUELTO. La API de Siigo no expone remisiones.** Plan A descartado,
aplica Plan B. Pendiente decision del negocio.

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

## 2. Investigacion del endpoint: NO EXISTE

### 2.0 Conclusion (2026-08-19)

**La API de Siigo no tiene endpoint de remisiones.** No hay que elegir un path
ni un shape de payload: el recurso no esta expuesto.

Fuente: el API Blueprint oficial completo de Siigo (`FORMAT: 1A`,
`HOST: https://api.siigo.com`, 6.029 lineas), vendorizado como `siigoapi.apib`
en `github.com/jdlar1/siigo-mcp` (commit del 2026-05-05). Es el mismo documento
que sirve `siigoapi.docs.apiary.io`, que desde esta sesion esta bloqueado.

Reproducir:

```
git clone --depth 1 https://github.com/jdlar1/siigo-mcp
grep -ic remis siigo-mcp/siigoapi.apib      # -> 0
grep -ic remission siigo-mcp/siigoapi.apib  # -> 0
```

La tabla de recursos del blueprint, completa:

| Recurso | Endpoint |
|---|---|
| Productos o servicios | `/products` |
| Clientes | `/customers` |
| Cotizaciones | `/quotations` |
| Facturas de venta | `/invoices` |
| Facturas de compra | `/purchases` |
| Notas credito | `/credit-notes` |
| Recibos de caja | `/vouchers` |
| Recibos de pago/egreso | `/payment-receipts` |
| Comprobantes contables | `/journals` |
| Documento soporte | `/purchase-support-documents` |
| Reportes financieros | varios |

Y el enum `DocumentType` del blueprint es exactamente:

```
FV - Factura de Venta
RC - Recibo de Caja
NC - Nota Credito
FC - Factura de Compra
CC - Comprobante Contable
```

(mas `DS`, `RP` y `C` como valores de consulta en `/v1/document-types?type=`).
No hay tipo de remision.

El blueprint esta al dia: incluye los cambios del sector salud vigentes desde
el 2025-07-22, `/v1/invoices/batch`, `/v1/payment-receipts` y
`/v1/purchase-support-documents`, todos posteriores al SDK oficial de 2023. Es
decir, Siigo si viene agregando documentos nuevos a la API - remisiones no es
uno de ellos.

Residuo de incertidumbre: el blueprint tiene 3,5 meses. Si Siigo publicara
remisiones entre mayo y hoy, no aparece aca. Es poco probable pero no
imposible; se confirma con un correo a soporteapi@siigo.com.

### 2.1 Limitacion de esta sesion (por que costo encontrarlo)

La politica de egress de la organizacion bloquea el dominio de Siigo. Verificado:

```
curl https://developers.siigo.com/reference   -> CONNECT tunnel failed, 403
curl https://api.siigo.com/v1/remissions      -> CONNECT tunnel failed, 403
WebFetch siigonube.portaldeclientes.siigo.com -> EGRESS_BLOCKED
WebFetch siigoapi.docs.apiary.io              -> EGRESS_BLOCKED
```

Segun `/root/.ccr/README.md` un 403 del proxy es denegacion de politica: no se
reintenta ni se rodea. GitHub, PyPI y npm si pasan, y por ahi salio la doc.

### 2.2 Callejones sin salida (para no repetirlos)

- **El SDK oficial no sirve como referencia.**
  `github.com/SiigoSAS/siigo_sdk_javascript` lista 18 clases y ~46 operaciones,
  sin remisiones, pero es de **oct-2023** y esta desactualizado: no incluye
  `/v1/webhooks`, que nuestro propio cliente ya consume. Ausencia ahi no
  probaba nada. Igual con `saulmoralespa/siigo-api-php` y
  `srdorado/siigo-client-php`.
- **`siigodocs/siigo-api-docs`** (repo oficial) solo trae ejemplos de conexion,
  1 commit, sin catalogo de recursos.
- **Buscar en Google es contraproducente.** Todo lo que devuelve son paginas del
  portal de CLIENTES sobre como elaborar una remision *desde la UI* de Siigo
  Nube, o "documento tipo S - Nota de remision" de Siigo Pyme. Hace parecer que
  la funcionalidad esta en la API cuando esta solo en el producto.

### 2.3 Lo que si existe (y confunde)

La remision existe como documento del producto Siigo: "Configurar remisiones" y
"Elaborar remisiones" en Siigo Nube, "documento tipo S" en Siigo Pyme, e incluso
el cruce remision -> factura ("Elaborar factura electronica con cruce de
remision"). Todo eso es UI. Nada de eso esta expuesto por API.

---

## 3. Paso 0: CERRADO

Era: confirmar si la API expone remisiones. Resuelto en 2.0 con el blueprint
oficial. **No las expone.** Queda un solo pendiente opcional: correo a
soporteapi@siigo.com para (1) confirmar que no hay nada posterior a mayo 2026 y
(2) pedirlo como feature. No bloquea el Plan B.

---

## 4. Plan A - DESCARTADO (la API no soporta remisiones)

Se conserva como referencia por si Siigo publica el recurso mas adelante. Hoy
no es ejecutable: no hay endpoint contra el cual programar.

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
