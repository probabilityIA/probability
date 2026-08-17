# ALERTA: la autogeneracion de guia envia el destino de la cotizacion, no el de la orden

Fecha: 2026-08-17

## Contexto

`maybeAutoGenerate` (`shipments/internal/infra/primary/queue/consumer/order_created_autogen.go`)
clona `quote.RequestPayload` y solo sobreescribe `idRate`, `carrier`,
`order_uuid`, `external_order_id` y `myShipmentReference`. El bloque
`destination` viaja tal como quedo en la cotizacion.

Cuando la cotizacion se hizo desde el panel con datos de prueba, la guia se
pide con ese destinatario. En VIG-0095 (business 46) eso significo mandar
`email: ""`, `phone: ""` y `firstName: "Destinatario" / lastName: "Default"` con
una direccion distinta a la de la orden. El carrier respondio 422 y la guia no
salio. Diagnostico completo en
`.claude/bitacora/2026-08-17-cotizacion-guia-generada-fantasma.md`.

## Mitigacion parcial (2026-08-17)

`POST /shipments/quotes/:id/retry-guide` reintenta la MISMA cotizacion (misma
tarifa y mismo idRate) pero reconstruyendo el bloque `destination` con los datos
reales de la orden, y valida email/telefono/direccion antes de gastar el viaje al
carrier. `buildRetryPayload` en `handlers/retry-quote-guide.go` es la referencia
de como debe quedar el payload.

**La autogeneracion sigue sin ese refresco**: publica el destino de la
cotizacion tal cual. El arreglo de fondo es que `maybeAutoGenerate` use la misma
construccion que el reintento.

## Items

### Urgente

- Refrescar `payload["destination"]` con los datos reales de la orden
  (destinatario, email, telefono, direccion, ciudad, barrio, dane) antes de
  publicar la peticion de generacion. Hoy la cotizacion manda datos que pueden
  no tener nada que ver con la orden: ademas de fallar, si pasara la validacion
  generaria una guia dirigida a la persona equivocada.

### Importante

- Validar el destinatario antes de publicar: si falta email o telefono, fallar
  la cotizacion con motivo claro (`SanitizeGuideError` ya tiene los mensajes) en
  vez de gastar el viaje al carrier y aterrizar en un 422 generico.
- Definir el fallback de email. El carrier lo exige y muchas ordenes manuales no
  lo tienen; hoy eso condena la guia. Decidir entre un email de la tienda o
  bloquear la generacion pidiendo el dato.

### Deseable

- Cubrir con test el caso "cotizacion con destino placeholder + orden con
  destinatario real": el payload publicado debe llevar el de la orden.

## Criterio para cerrar

El payload publicado toma el destino de la orden, hay validacion previa con
motivo legible, y una orden manual sin email de cliente no puede terminar en un
422 del carrier.
