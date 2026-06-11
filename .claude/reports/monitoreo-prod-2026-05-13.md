# Monitoreo Produccion - 2026-05-13

Ventana: 2 horas, cada 10 min. Cron job: `bd3e0961` (4,14,24,34,44,54 * * * *).

Comando base:
```
ssh -i .../probability.pem ubuntu@ec2-3-224-189-33.compute-1.amazonaws.com
cd /home/ubuntu/probability/infra/compose-prod
docker ps + docker compose logs --since 10m back-central nginx front-central rabbitmq
grep -iE 'panic|error|fatal|5[0-9]{2}|reject|retry|oom|killed'
```

---

## Iter 1 - 15:58 local (20:58 UTC)

**Estado:** OK con 1 nota persistente.

- Containers todos `Up`. `website_prod` aparece **unhealthy** hace 2 semanas (estado preexistente, no es regresion).
- back-central: solo INF. Webhooks Shopify (`sin-intermediarios-co.myshopify.com`, topic `orders/updated`) procesados OK, eventos SSE despachados a business_id=34, integration_id=57. Sin panic / error / 5xx / reject.
- nginx: solo 200s, trafico desde `186.81.100.188` navegando /orders, /customers, /products, etc.
- front-central: sin output en filtro de errores.
- rabbitmq: sin errores ni warnings.

**Bugs / cosas raras:** ninguno nuevo. Pendiente investigar healthcheck de `website_prod` (no urgente).

---

## Iter 2 - 16:14 local (21:14 UTC)

**Estado:** WARN - 1 error recurrente detectado.

- Containers sin cambios.
- **ERR back-central 16:14:17:** `Failed to confirm sale` error=`"no se encontro bodega por defecto para el negocio"` function=`handleConfirmSale` module=`inventory.consumer` order_id=`19b4fa9a-b6ad-4432-b6c3-8c176e1b2c31` (business_id=34, integration_id=57, Shopify external_id=5373773971665).
  - Disparado por el evento derivado `order.delivered` (status_changed -> shipped -> delivered).
  - Causa: business 34 no tiene bodega por defecto configurada. El consumer de inventario no puede descontar stock al confirmar venta.
  - Impacto: el inventario no se actualiza al entregar la orden. La orden se procesa, pero el modulo inventory queda inconsistente.
  - Accion sugerida: configurar bodega por defecto para business 34, o ajustar el handler para fallar mas suavemente / re-encolar.
- Resto del back: webhooks Shopify procesando normal, eventos SSE OK.
- nginx: solo 200s. front-central y rabbitmq sin errores.

---

## Iter 3 - 16:24 local (21:24 UTC)

**Estado:** OK. Sin novedad - solo logs INF de procesamiento normal de webhooks Shopify y eventos SSE.

---

## Iter 4 - 16:34 local (21:34 UTC)

**Estado:** WARN - 1 error en invoicing retry consumer + 1 warning de salto de estado.

- Containers sin cambios.
- **ERR back-central 16:34:23:** `Failed to process invoice` error=`"reintento no permitido"` invoice_id=`24386` sync_status=`failed` module=`invoicing.retry_consumer`.
  - El retry consumer encontro 1 retry pendiente, intento procesar invoice 24386, y el usecase rechazo el reintento. `failed=1 success=0 total=1`.
  - Causa probable: invoice ya esta en estado terminal o supero max reintentos pero sigue encolado.
  - Impacto: bucle ruidoso. La factura no se va a sincronizar. Revisar tabla de invoices/retries y limpiar o desactivar este retry.
- **WRN back-central 16:34:14:** `Integracion realizo salto de estado que no cumple flujo v2 - aceptado por ser fuente externa` from=`confirmed` to=`in_transit` order_id=`bbeb33cb-...` Shopify. Comportamiento esperado para fuentes externas, no es bug.
- Resto del trafico normal. nginx 200s. front-central/rabbitmq sin errores.

---

## Iter 5 - 16:44 local (21:44 UTC)

**Estado:** OK. Sin novedad - solo INF de webhooks Shopify y eventos SSE. El retry de invoicing no aparecio esta ventana.

---

## Iter 6 - 16:54 local (21:54 UTC)

**Estado:** OK. Sin novedad - 3 webhooks Shopify procesados (orders 5375611470033, 5375238963409, 5375167070417), todos INF.

---

## Iter 7 - 17:04 local (22:04 UTC)

**Estado:** WARN - reincidencia del bucle de invoice 24386.

- Containers sin cambios.
- **ERR back-central 17:02:18:** mismo `Failed to process invoice` invoice_id=`24386` error=`"reintento no permitido"` sync_status=`failed`. **Patron confirmado:** se dispara aprox cada ~28 min (iter 4 = 16:34, iter 7 = 17:02). El retry consumer requeue indefinido sin progresar.
- Nuevo evento sano: webhook `orders/paid` recibido 17:01:14 (external_id=5375626346705, order_id=8d69fddc-...) procesado OK.
- Resto INF normal: webhooks `orders/updated`, eventos SSE.
- nginx 200s. front/rabbit limpios.

**Accion sugerida (post-monitoreo):** marcar invoice 24386 como dead-letter o ajustar handler para no requeue cuando sync_status=`failed` y el usecase rechace.

---

## Iter 8 - 17:14 local (22:14 UTC)

**Estado:** WARN - 2 `Failed to confirm sale` (mismo bug de iter 2).

- Containers sin cambios.
- **ERR back-central 17:14:10:** `Failed to confirm sale` order_id=`71f8f838-...` (#92534) - business 34 sin bodega default.
- **ERR back-central 17:14:12:** mismo error, order_id=`2254805a-...` (#92458). Disparado tras evento derivado `order.delivered`.
- **WRN 17:14:12:** salto `in_transit -> delivered` en order 2254805a (Shopify externa, esperado).
- Frontend muestra lineas tipo `code: 'failed'`, `category: 'failed'`, `delivery_failed` en logs - son **payloads de datos siendo serializados** (no errores reales del proceso). Falso positivo del grep.
- nginx: solo 200s, navegacion en `/invoicing/invoices`.

**Confirmado patron iter 2 + iter 8:** cada orden Shopify que llega a `delivered` falla el confirm sale del inventory consumer en business 34. Acumulando inconsistencia de stock.

---

## Iter 9 - 17:24 local (22:24 UTC)

**Estado:** WARN - 2 errores: invoice 24386 (recurrente) + bug nuevo de announcement con super admin.

- Containers sin cambios.
- **ERR back-central 17:19:03:** invoice 24386 retry rechazado nuevamente (`reintento no permitido`). Patron confirmado: 16:34, 17:02, 17:19 - intervalos de 28 min y luego 17 min (no es exactamente periodico).
- **BUG NUEVO 17:22:51:** `query failed` FK violation en `announcement_views`:
  ```
  INSERT INTO announcement_views (announcement_id=14, user_id=6, business_id=0, action='closed', ...)
  fk_announcement_views_business (SQLSTATE 23503)
  ```
  - Devuelve HTTP 400 al frontend (POST `/api/v1/announcements/14/view`).
  - Causa: usuario es **super admin** (business_id=0 en JWT), pero la tabla `announcement_views` tiene FK a `businesses.id` que no acepta 0.
  - Frontend tambien lo logea como ⨯ Error.
  - Impacto: super admins no pueden marcar anuncios como vistos/cerrados.
  - Fix: o aceptar `business_id NULL` para super admin, o omitir el insert cuando business_id=0.
- SSE long-poll con latencias ~54s (normal para SSE).
- nginx 200s. rabbit limpio.

---

## Iter 10 - 17:34 local (22:34 UTC)

**Estado:** OK. Sin errores nuevos.

- Containers sin cambios.
- WRN salto `paid -> confirmed` en order f2efbae2 (#92875), webhook `orders/fulfilled` Shopify externa - esperado.
- Webhook procesado OK + score calculado.
- No retry de invoice 24386 esta ventana, no FK violation.
- Front-central muestra solo dumps de datos `code: 'failed'`/`delivery_failed` (falsos positivos del grep).
- nginx 200s, rabbit limpio.

---

## Iter 11 - 17:44 local (22:44 UTC)

**Estado:** WARN - patrones recurrentes ya conocidos.

- Containers sin cambios.
- **ERR 17:39:44:** `Failed to confirm sale` order `be9b0293-...` (otro delivered de Shopify, sin bodega default en biz 34). Acumulado.
- **ERR 17:41:23:** invoice 24386 retry rechazado nuevamente (`reintento no permitido`). 4ta repeticion en el monitoreo (16:34, 17:02, 17:19, 17:41).
- SSE long-poll con duraciones ~348s y ~352s (~6 min) - mas largo que iter 9 pero dentro de comportamiento esperado.
- Sin bugs nuevos esta ventana.

---

## Iter 12 - 17:54 local (22:54 UTC) - FINAL

**Estado:** WARN - patrones recurrentes ya conocidos.

- Containers sin cambios.
- **ERR 17:52:47:** `Failed to confirm sale` order `cd685374-...` (#92444), bodega default biz 34.
- **ERR 17:54:03:** `Failed to confirm sale` order `1df635b3-...` (#92591), mismo bug.
- Webhooks Shopify procesando OK (orders/updated, derivados a shipped/delivered).
- nginx 200s + busqueda activa de places-search (oficinas Servientrega, Inter Rapidisimo, Coordinadora, etc.) - usuario interactivo.
- Sin bugs nuevos. Cron `0dd1c6fc` detenido tras esta iteracion.

---

## Resumen de la ventana 2h (15:58 - 17:54 local, 2026-05-13)

**Estado general:** Servicio funcionando. Trafico normal. Sin caidas, sin panics, sin OOM, sin restarts, sin webhooks rechazados, sin 5xx en nginx.

**Bugs detectados (priorizar fix):**

1. **CRITICO ops (recurrente):** business 34 sin **bodega por defecto**. Cada orden Shopify que llega a `delivered` falla `handleConfirmSale` en `inventory.consumer`. Vistos: iters 2, 8 (x2), 11, 12 (x2) = **6 ocurrencias** en 2h. El stock NO se descuenta para esas ordenes.
   - Fix: configurar bodega default para business 34, o degradar handler para fallar suavemente / continuar.

2. **CRITICO bucle de logs:** invoice **24386** queda atascada en retry. El consumer encuentra el item pendiente, llama `RetryInvoice`, el usecase rechaza con `"reintento no permitido"` (sync_status=`failed`). Vistos: iters 4, 7, 9, 11 = 4 disparos (intervalo irregular ~17-28 min).
   - Fix: marcar 24386 como dead-letter/abandoned, o ajustar query de `Found pending retries` para excluir items con sync_status=`failed` que el usecase rechaza.

3. **BUG super admin:** super admin (`business_id=0` en JWT) no puede registrar vistas de anuncios. `INSERT INTO announcement_views` viola FK `fk_announcement_views_business`. Endpoint `POST /api/v1/announcements/14/view` devuelve 400. Visto en iter 9 (1 vez).
   - Fix: aceptar `business_id NULL` para super admin en la tabla, o saltarse el insert cuando business_id=0, o usar `business_id` desde query param (patron del resto del modulo).

**Comportamiento esperado (no son bugs):**
- WRN `Integracion realizo salto de estado que no cumple flujo v2` para Shopify (paid->confirmed, in_transit->delivered, etc.) - aceptados por ser fuente externa.
- SSE long-poll con duraciones largas (54s - 6 min) - normal.
- `website_prod unhealthy` desde hace 2 semanas - preexistente, no es regresion.

**Trafico observado:**
- Webhooks Shopify (orders/updated, orders/paid, orders/fulfilled) procesando con normalidad para `sin-intermediarios-co.myshopify.com` (integration_id=57, business_id=34).
- Eventos derivados (status_changed -> shipped -> delivered) y SSE broadcast OK.
- Frontend: navegacion normal, busqueda de places-search activa.
- Bold/Meta WhatsApp: sin webhooks recibidos en esta ventana.

---
