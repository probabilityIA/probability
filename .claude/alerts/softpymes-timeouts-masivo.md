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

- [URGENTE] Riesgo de facturas DUPLICADAS en SoftPymes/DIAN del 2026-07-23: los
  406 "response has no info" y 230 timeouts de creacion pudieron crear la factura
  igual y el codigo viejo la re-creo. Cuando SoftPymes vuelva, buscar comments
  `order:<uuid>` repetidos de ese dia.
- [IMPORTANTE] Relanzar las facturas que murieron en `cancelled` durante la caida
  (agotaron max_retries contra un proveedor caido, antes del fix de presupuesto).
- [DESEABLE] Arreglar contadores de `bulk_invoice_jobs` (failed > total_orders).

## Criterio de cierre

Retry no-idempotente eliminado + verificacion de duplicados hecha + throttle/backoff
implementado y probado con un job masivo real sin errores en cascada.
