# Guias duplicadas en EnvioClick: la ventana entre el envio y la respuesta

Fecha: 2026-08-25
Modulo: back/central shipments + integrations/transport/envioclick, front/central modal de guia

## Sintoma reportado

"En Probability se ve una sola guia pero en EnvioClick aparecen dos por la misma
orden." La segunda queda huerfana: con recoleccion programada, sin tracking en el
shipment y sin cobro al negocio, o sea la paga Probability.

## Diagnostico

El guard de idempotencia (`shipmentHasActiveGuide`) preguntaba a la BD si la
orden ya tenia tracking o guide_url. Pero la generacion es asincrona: el handler
responde 202, publica a RabbitMQ, y el tracking solo se persiste cuando vuelve la
respuesta. Entre esos dos momentos el guard estaba ciego y cualquier segundo
disparo creaba OTRA guia real.

El disparador es el front: `shipment-guide-modal.tsx` cortaba a los 45 s con
"Tiempo de espera agotado", hacia `setLoading(false)` y rehabilitaba el boton.
El usuario volvia a darle.

### Datos que sostienen el diagnostico

Duracion real de `POST /api/v2/shipment` en produccion, 952 llamadas de 90 dias:

| p50 | p90 | p99 | max | >30 s | >45 s |
|---|---|---|---|---|---|
| 4,9 s | 6,8 s | 18,3 s | 27,4 s | 0 | 0 |

O sea el carrier NO es lento. Lo que falla es el camino de vuelta
(cola de respuesta -> SSE al modal).

Caso MYS-0718 (shipment 43991, 2026-08-18), dos correlation_id distintos, ambos
`triggered_by=user`, usuario 20:

- 19:30:49 -> idOrder 4632833, tracker 2284436982, la llamada duro 5,6 s
- la respuesta nunca se persistio: ni tracking ni transaccion de wallet
- ~19:31:34 el modal muestra "Tiempo de espera agotado" y libera el boton
- 19:32:33 -> idOrder 4632845, tracker 2284436983, 5,9 s

Guias huerfanas confirmadas en 90 dias (log success cuyo tracker no quedo en
`shipments`): 2284436982 (18/08), 014160988633 (08/07), 84151641204 y
84151641209 (06/06).

### Hipotesis descartadas

- **El retry HTTP del cliente duplicaba.** No: `RetryCount: 2` esta anulado por
  la condicion `AddRetryCondition(r.StatusCode() == 429)`, que devuelve false
  ante un error de red, asi que un timeout no reintenta.
- **EnvioClick tarda y por eso el usuario espera.** No: p99 de 18 s contra un
  timeout de front de 45 s. El problema es que el resultado no llega a pantalla.
- **Es el mismo bug de 2026-06-18 (dos shipments).** Parcial: aquel se arreglo,
  este ocurre sobre el MISMO shipment.

## Segundo agujero, encontrado al probar

El `httpclient` tiene timeout de 30 s. Al simular un carrier lento con un mock,
la llamada se corto a los 30 s, el shipment quedo `failed` y el sistema permitia
reintentar de inmediato... con la guia ya creada del otro lado. Es la via que
explica las huerfanas sin log de exito.

## Correccion

Backend:

- Estado `generating` en `shipments`, reservado con `MarkShipmentGenerating`
  (UPDATE condicional) ANTES de publicar a la cola; se libera si el publish falla.
- Seccion critica bajo lock por orden: `pg_advisory_xact_lock(hashtext(order_id))`
  en `WithOrderGuideLock`, aplicado en generacion manual, reintento de cotizacion
  y autogeneracion. Sin esto, N clics simultaneos creaban N shipments distintos y
  cada uno bloqueaba solo el suyo (medido: 3 de 5 pasaban).
- Estado `needs_verification` cuando el carrier no confirma respuesta
  (`ErrTransportUnreachable` -> `error_kind: unreachable`): no se permite
  regenerar hasta verificar.
- `generating` caduca a los 15 min (`GuideGenerationStaleAfter`) para no trabar
  ordenes si se pierde la respuesta.

Front:

- El modal ya no rehabilita el boton a los 45 s: consulta el shipment cada 5 s
  (hasta 2 min) y cierra con exito, error o aviso de verificacion.
- Boton bloqueado con "Verifica la guia antes de reintentar" cuando el shipment
  quedo en `needs_verification`.

## Verificacion (local, mock de EnvioClick en :9099)

| Escenario | Antes | Despues |
|---|---|---|
| 2do clic con guia en vuelo | 2da guia | 409 `in_flight`, 1 llamada |
| 5 clics simultaneos | 3 guias | 1x202 + 4x409, 1 llamada, 1 shipment |
| Timeout de 30 s | `failed`, reintento duplicaba | `needs_verification`, 409 |
| Generacion normal | ok | ok, reintento da 409 "guia activa" |

E2E en navegador con ordenes DEM-0043 y DEM-0044 (business 26): el modal muestra
el aviso, el boton queda deshabilitado y el carrier recibio una sola llamada por
escenario.

## Pendiente

Reconciliar contra EnvioClick por `myShipmentReference` antes de crear, para
resolver solo los `needs_verification` en vez de pedir revision manual.
