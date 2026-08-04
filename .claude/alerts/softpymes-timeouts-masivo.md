# SoftPymes reventando por timeouts en facturacion masiva (business 34 "sin intermediarios")

Fecha: 2026-07-23

## Contexto

El negocio 34 corrio 5 jobs masivos el 2026-07-23 (1000 + 1000 + 1000 + 297 + 533 ordenes).
Los primeros pasaron casi limpios; los de la noche colapsaron. Errores en
`invoice_sync_logs` de ese dia:

- 406x `invoice response has no info:` (SoftPymes responde 200 con body vacio bajo carga)
- 230x `context deadline exceeded` en POST `/app/integration/sales_invoice/` (timeout 30s)
- 193x timeout en POST `/oauth/integration/login/` (auth tambien saturado)

El job de 533 ordenes quedo con `failed=890` (contadores inflados por reintentos).

## Causas en codigo

1. `shared/httpclient` (resty) con `RetryCount: 2`: ante timeout reintenta el POST de
   creacion de factura 2 veces mas SIN check de idempotencia (el check solo corre con
   `operation == "retry"`). Triplica la carga justo cuando SoftPymes esta lento.
2. Consumer `invoicing.softpymes.requests` con 3 workers (default, `SOFTPYMES_INVOICE_WORKERS`
   no seteado en prod) pegandole sin pausa; cada factura son 3-8 llamadas HTTP
   (auth, ensureCustomer, POST, hasta 4 GetDocument, recibo de caja).
3. `RetryConsumer` cada 5 min re-postea hasta 50 fallidas; si `findExistingInvoiceByOrderID`
   falla (tambien por timeout) hace "proceeding with creation".
4. Sin rate limiting ni backoff/circuit breaker hacia SoftPymes.

## Resuelto 2026-07-23/24

- Retry automatico de resty eliminado, workers=1, timeout 90s, retry fail-closed,
  busqueda de existentes paginada (commit f6eec211, desplegado).
- Errores de proveedor-caido (timeouts de red, login inalcanzable, 5xx, aborto de
  verificacion) ya NO consumen retry_count: se reprograman cada 15 min hasta que
  SoftPymes vuelva (response_consumer, isProviderUnavailableError).

## Items

- [RESUELTO 2026-07-23] Auditoria de duplicados ejecutada contra la API de SoftPymes
  (3.337 documentos del dia, cruce por comment `order:<uuid>`): 0 fantasmas,
  0 inconsistencias, solo 7 ordenes con doble factura (pares identicos).
- [URGENTE - MANUAL] Anular en Pymes+ (web) los 7 documentos sobrantes; la API de
  integracion NO tiene endpoint de anulacion. Anular estos, conservando el par:
  0000006255 (queda 0000006259), 0000006256 (queda 0000006258),
  0000006257 (queda 0000006260), 0000006943 (queda 0000007176),
  0000007282 (queda 0000007287), 0000007283 (queda 0000007286),
  0000007284 (queda 0000007285). Total doble-facturado ~1.425.000 COP.
- [IMPORTANTE] Bug latente: `CancelInvoice` del cliente SoftPymes apunta a
  `/app/integration/sales_invoice/cancel/` que NO existe (404 de ruta) — la
  cancelacion de facturas SoftPymes nunca ha funcionado. Confirmar con SoftPymes
  si existe endpoint de anulacion; si no, quitar/deshabilitar esa opcion en UI.
- [IMPORTANTE] Seguridad: el cliente SoftPymes loguea apiKey/apiSecret en texto
  plano en consola (request_body en authenticate + SetDebug). Redactar esos logs.
- [IMPORTANTE] Relanzar las facturas que murieron en `cancelled` durante la caida
  (agotaron max_retries contra un proveedor caido, antes del fix de presupuesto).
- [DESEABLE] Arreglar contadores de `bulk_invoice_jobs` (failed > total_orders).

## Incidente 2026-08-03: RDS OOM congelo 449 facturas en "pending"

Bulk job `04aad4f7-d117-477e-a459-83587d8d1602` (business 34, 540 ordenes,
lanzado 15:02 Bogota). A las 15:33-15:37 Bogota el RDS `database-1` se quedo
sin memoria y se reinicio (evento AWS: "workload causing the system to run
critically low on memory"; RDS bajo shared_buffers de 23081 a 11295). Durante
la caida (~9.580 connection refused en 2 min) el consumer de SoftPymes drenó el
backlog de la cola: cada create fallo en ~4ms (no podia leer la integracion de
la DB) y la respuesta de error tampoco se pudo persistir.

Resultado: 449 invoices en `pending` con sync log `create` en `processing`
(response_status=0, sin body). NUNCA llegaron a SoftPymes: check_status barrio
8 dias / 1.417 documentos y no existen -> reintentar creacion es seguro, sin
riesgo de duplicados. El cron de reconciliacion solo hace check_status (query),
nunca re-crea: estas facturas NO van a avanzar solas ("Document not found yet —
DIAN still validating" en loop).

- [EN CURSO 2026-08-03] Relanzamiento ejecutado con autorizacion del usuario:
  449 invoices pasadas a status='failed', sus 449 sync logs create
  processing->failed con next_retry_at=now, y 5 logs query pending->cancelled
  (UPDATE directo en RDS, transaccional). El retry consumer las procesa en
  lotes de 50 cada 5 min via RetryInvoice (idempotencia fail-closed contra
  SoftPymes). Primer lote verificado OK en logs de prod (49/50 publicadas,
  facturas emitiendose). Drenaje estimado ~3h. Verificar al final que
  status='pending' o 'failed' quede en 0 para el burst 20:02:45-20:03:03 UTC.
- [URGENTE] Cerrar el bulk job 04aad4f7 (quedo en `processing`, 539/540,
  successful=100, failed=5, contadores nunca actualizados tras la caida).
- [IMPORTANTE] El response_consumer pierde la respuesta si la DB esta caida al
  procesarla (el mensaje se ACKea y el estado queda congelado). Falta nack/requeue
  o retry con backoff cuando el fallo es de DB, igual que isProviderUnavailableError.
- [IMPORTANTE] Capacity: el RDS hace OOM con jobs masivos grandes (ya bajo
  shared_buffers solo). Evaluar subir instancia o limitar tamano de bulk.

## Criterio de cierre

Retry no-idempotente eliminado + verificacion de duplicados hecha + throttle/backoff
implementado y probado con un job masivo real sin errores en cascada.
