# Cotizaciones mostraba "Guia generada" para guias que nunca salieron

VIG-0095 aparecia sin guia en el modulo de guias pero con "Guia generada" en el
modulo de cotizaciones. La guia habia fallado con 422 del carrier y la
cotizacion nunca se entero.

## Sintoma

- Orden `VIG-0095` (`f7e14413-b4a5-4f25-8464-4460549f3272`), business 46 (Viga
  ropa deportiva), manual, COD, $172.500. Sin `tracking_number` ni `guide_link`.
- Shipment 43412: `status = failed`, sin tracking ni `guide_url`.
- Cotizacion 6542 (`shipping_quotes`): `status = guide_generated`.
- Ordenes y guias decian "fallida". Cotizaciones decia "Guia generada".

## Diagnostico

`shipment_sync_logs` id 13543, 2026-08-17 20:30:09:

```
POST /api/v2/shipment  ->  422
{"error":["Unprocessed Entity.",
  {"destination":{"email":["Dato invalido, el campo debe ser de tipo email"]}}]}
```

El `request_payload` de la cotizacion tenia el destino de la cotizacion de
prueba, no el de la orden:

```json
"destination": { "email": "", "phone": "", "firstName": "Destinatario",
                 "lastName": "Default", "address": "carrera 46 # 50 -21c" }
```

El destinatario real era Ferney Gonzalez, tel 3155707677, carrera 11#17-27
barrio Juan mayor, y sin email.

Linea de tiempo: cotizacion creada 20:25:06 desde el panel, orden creada
20:30:01, generacion disparada y rechazada 20:30:09.

### Hipotesis descartadas

- **La guia se genero y se perdio el tracking.** No: no hay registro de exito en
  `shipment_sync_logs` ni cobro en `transaction`. La guia nunca existio.
- **El modulo de guias filtra mal y la esconde.** No: el shipment esta en
  `failed`, que es el estado correcto. El que mentia era cotizaciones.
- **Doble generacion como en la alerta de guias duplicadas.** No: un solo
  shipment, un solo intento.

## Causa raiz

Dos defectos independientes.

1. **El payload de generacion reusa el destino de la cotizacion.**
   `order_created_autogen.go` clona `quote.RequestPayload` y solo sobreescribe
   `idRate`, `carrier`, `order_uuid`, `external_order_id` y
   `myShipmentReference`. El bloque `destination` nunca se refresca con los datos
   reales de la orden, asi que se envio el destinatario placeholder con el email
   vacio.

2. **La cotizacion se marcaba generada antes de saber el resultado.** El mismo
   archivo ponia `guide_generated` justo despues de publicar a RabbitMQ. Cuando
   llegaba el 422, `handleGenerateResponse` marcaba el shipment `failed` y
   emitia el SSE, pero no tocaba la `shipping_quote`: quedaba en
   `guide_generated` para siempre.

## Correccion

Codigo:

- Estado nuevo `generating`: la cotizacion pasa a `guide_generated` solo cuando
  el carrier confirma, y a `failed` con motivo cuando rechaza. Aplica tanto a la
  autogeneracion como a la asociacion manual con `guide_requested`.
- `response_consumer.go`: `markQuoteFailed` / `markQuoteGuideGenerated` cierran
  el ciclo desde la respuesta del carrier, buscando la cotizacion por
  `order_uuid`.
- `domain/quote_error.go`: sanitizador de mensajes. Traduce el error del
  proveedor a un motivo accionable en español y **nunca deja pasar el nombre del
  broker**. Sin coincidencia conocida devuelve un mensaje generico, no el texto
  crudo. `MentionsProvider` es la red de seguridad.
- El guard de idempotencia de `order_created_consumer.go` ahora cubre tambien
  `generating`, para no reabrir la puerta a guias duplicadas.
- Front: badge "Generando guia", filtro nuevo, columna "Motivo del error" y
  boton "Reintentar" en las cotizaciones fallidas, que abre el modal de
  generacion de guia con la orden cargada. Ese modal arma el destino con los
  datos reales y exige email valido por schema, asi que obliga a corregir el
  dato que causo el 422 antes de reintentar. El titulo "Generar Guia Envioclick"
  paso a "Generar Guia de Envio".

Datos (produccion, `fixVig0095CotizacionFallida`):

- Cotizacion 6542 (VIG-0095) y 6441 (DEM-0040, business 26) pasaron de
  `guide_generated` a `failed` con su motivo real. Eran las dos unicas
  inconsistentes.

## Verificacion

```sql
SELECT count(*) FROM shipping_quotes q
JOIN shipments s ON s.order_id = q.order_uuid
WHERE q.status = 'guide_generated' AND s.status = 'failed';
-- 0
```

## Pendientes

El defecto 1 sigue abierto: la autogeneracion continua enviando el destino de la
cotizacion. Mientras no se refresque con los datos de la orden, toda cotizacion
del panel hecha con destinatario placeholder va a fallar en 422. Ver
`.claude/alerts/autogen-guia-usa-destino-de-la-cotizacion.md`.
