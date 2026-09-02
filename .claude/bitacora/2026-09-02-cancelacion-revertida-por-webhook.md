# Cancelacion revertida por un webhook de tracking atrasado

Un envio se cancelo bien en EnvioClick y se marco `cancelled` en Probability,
pero medio segundo despues un update de tracking lo devolvio a `pending`. Para
el usuario parecia que la cancelacion no habia funcionado.

## Sintoma

Orden MYS-0852, shipment 48290, guia `034058470201`, carrier ENVIA.
El usuario cancela, EnvioClick queda cancelada, pero la lista de envios sigue
mostrando el envio como Pendiente y al reintentar "da error".

## Evidencia

`shipment_sync_logs` del shipment 48290, 2026-09-02:

| Hora (UTC) | Llamada | Respuesta |
|---|---|---|
| 20:17:33.665 | `POST /v2/track` | `Pendiente de Recoleccion` |
| 20:17:33.723 | `POST /v2/cancellation/batch/order` `{"idOrders":[4670410]}` | `to_refund_orders: [4670410]` |
| 20:17:33.729 | `POST /v2/track` | `Cancelado` |

O sea: EnvioClick si cancelo y la verificacion del commit `0088c3ed` lo confirmo.

`shipments.metadata.tracking_events` del mismo envio:

```
15:17:33 -05  status=cancelled  raw_status=Cancelado           <- la cancelacion
00:00:00      status=pending    raw_status=Pendiente de recoleccion  <- la piso
```

El segundo evento trae `has_incidence` y una fecha de carrier a medianoche: es un
`webhook_update`, no la cancelacion. `shipments.updated_at` quedo en
20:17:34.280, un segundo despues de la cancelacion, con `status = pending` y
`carrier_status = 'Pendiente de recoleccion'`.

## Causa raiz

`response_consumer.go` pisaba el estado sin mirar el que ya tenia:

- `handleWebhookUpdate` hacia `shipment.Status = probabilityStatus` siempre.
- `handleTrackResponse` hacia `shipment.Status = status` siempre.

El webhook de EnvioClick venia con la foto anterior a la cancelacion (evento del
carrier fechado a medianoche) y llego despues. Gana el ultimo que escribe, no el
mas nuevo.

`cancelled` es terminal en Probability: ningun tracking posterior deberia sacarlo
de ahi.

## Hipotesis descartadas

- **La cancelacion fallo en EnvioClick.** No. `to_refund_orders: [4670410]` y el
  track posterior devolvio `Cancelado`. Este NO es el falso positivo de
  `2026-08-11-envioclick-cancelacion-falso-positivo.md`.
- **La logica de "ya esta cancelada alla" no existia.** Existe:
  `envioclick/internal/app/operations.go` linea 63 devuelve exito cuando el track
  previo dice `cancelad*`. El reintento del usuario habria funcionado.
- **Error del backend en el reintento.** El popup que vio el usuario decia
  `Server Action "4039921f8f880d72653f0ef5f6fef675bbc8ff3074" was not found on
  the server`: es la pestana con el bundle de un deploy viejo, no el backend.
  Se resuelve recargando.
- **`envioclick_order_id` nulo en la fila.** Lo esta, pero el id llega igual por
  `metadata.envioclick_id_order` (4670410), asi que no rompio nada aca.

## Correccion

`revivesCancelledShipment(current, incoming)` en `response_consumer.go`: si el
envio ya esta `cancelled` y llega cualquier otro estado, se conserva `cancelled`
y `carrier_status`, se deja el evento en el historial y se escribe un `Warn`
(pasa el filtro de logs de CloudWatch, asi que un "entregado despues de
cancelado" queda visible). Aplicado en `handleWebhookUpdate` y en
`handleTrackResponse`.

Tests en `response_consumer_test.go`:
`TestHandleWebhookUpdate_CancelledShipment_KeepsCancelled`,
`TestHandleWebhookUpdate_NonCancelledShipment_AppliesStatus`,
`TestRevivesCancelledShipment`.

## Alcance

Cuatro envios historicos con un evento `cancelled` en el historial y estado
distinto de `cancelled`: 48290 (este), 38712, 34247, 34071.

## Pendiente

- El shipment 48290 sigue en `pending` en produccion. Con el fix desplegado,
  volver a cancelar desde la UI lo deja bien: el track dice `Cancelado`, se marca
  `cancelled` y ya nada lo revierte. No hace falta tocar la base a mano.
