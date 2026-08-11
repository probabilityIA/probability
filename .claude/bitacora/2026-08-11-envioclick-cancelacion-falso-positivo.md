# EnvioClick: cancelaciones con falso positivo

EnvioClick cambio el comportamiento de su API de cancelacion entre el 7 y el 8 de
agosto de 2026: dejo de cancelar las guias pero siguio respondiendo `200 OK`.
Probability daba esa respuesta por exitosa y marcaba el envio como `cancelled`,
asi que nadie se enteraba de que la guia seguia viva en la transportadora.

## Sintoma

Reportan que las cancelaciones desde Probability dan falso positivo: la orden
queda cancelada en Probability pero no en EnvioClick.

Cuatro guias del business 36 (integracion 51, carrier ENVIA) canceladas el
2026-08-11 seguian activas en la transportadora horas despues:

| Guia | idOrder | Cancelada en Probability | Estado real en EnvioClick | Recoleccion |
|------|---------|--------------------------|---------------------------|-------------|
| 034058168864 | 4613110 | 09:10 | Pendiente de recoleccion | 2026-08-12 |
| 034058168817 | 4613100 | 09:11 | Pendiente de recoleccion | 2026-08-12 |
| 034058187592 | 4616473 | 10:31 | Pendiente de recoleccion | 2026-08-14 |
| 034058168475 | 4613038 | 10:38 | Pendiente de recoleccion | 2026-08-13 |

Riesgo concreto: guias que el negocio cree canceladas pueden ser recolectadas y
cobradas.

## Diagnostico

### La respuesta que llegaba

Las cuatro llamadas fueron `200 OK` con las tres listas vacias:

```
POST https://api.envioclickpro.com.co/api/v2/cancellation/batch/order
{"idOrders":[4613038]}

{"status": "OK", "status_codes": [200],
 "status_messages": [{"request": "Request processed."}],
 "data": {"not_valid_orders": [], "only_cancel_orders": [], "to_refund_orders": []}}
```

`"Request processed."` es solo un acuse de recibo. No confirma cancelacion.

### La evidencia que fecha el cambio

`shipment_sync_logs` guarda request y response de cada llamada (desde el commit
`dc42daf2`). Comparando el mismo endpoint y el mismo body a lo largo del tiempo:

| Fecha | `to_refund_orders` | Cancelo? |
|-------|--------------------|----------|
| 2026-07-06 .. 2026-08-07 03:23 | `[idOrder]` | Si, siempre |
| desde 2026-08-08 15:42 | `[]` | No, nunca |

Ultima buena: orden `4612012`, 2026-08-07 03:23, `"to_refund_orders": [4612012]`.
Primera rota: orden `4613136`, 2026-08-08 15:42, las tres listas vacias.

El codigo de cancelacion no se toca desde el 2026-04-29. El cambio es de
EnvioClick, no nuestro.

### Alcance real

Seis cancelaciones afectadas, todas desde el 08-08:

- `4613136` (08-08) y `4609403` (08-10): quedaron canceladas porque el equipo las
  cancelo a mano en la plataforma. **No** las cancelo la API.
- `4613110`, `4613100`, `4616473`, `4613038` (08-11): quedaron vivas.

Las 15 cancelaciones anteriores al 08-08 (ultimos 30 dias) si quedaron canceladas.

### Prueba directa

Cancelacion manual por API de la orden `4613038`, ya cancelada por Probability a
las 10:38 del mismo dia:

```
track   -> "Pendiente de recoleccion"
cancel  -> 200 OK, las tres listas vacias
track   -> "Pendiente de recoleccion"   (inmediato)
track   -> "Pendiente de recoleccion"   (+5 min)
```

EnvioClick acepta la cancelacion cuantas veces se le pida y no cancela.

### Hipotesis descartadas

- **`idOrder` invalido.** No. `POST /track-by-orders` con `{"orders":[4613038]}`
  lo reconoce y devuelve su estado.
- **Nombre de campo equivocado.** No. `idOrders` es el correcto para cancelacion:
  con un id inexistente responde `422 "orders not found"`, y con `{"orders":[...]}`
  responde `422 idOrders es obligatorio`. Ojo con la inconsistencia de EnvioClick:
  `track-by-orders` usa `orders` y `cancellation/batch/order` usa `idOrders`. El
  codigo lo tiene bien en ambos.
- **Guia en estado no cancelable.** No. El `Track` previo a cada cancelacion
  devolvio "Pendiente de recoleccion", que es el estado valido.
- **Error de red, 4xx o 5xx.** No. Los cuatro requests fueron `200` en 0.7-4.1s.
- **Procesamiento asincrono con retraso.** No. Cinco minutos despues seguia
  activa. Las del 08-08 y 08-10 que si terminaron canceladas lo fueron a mano,
  no por la API (confirmado con el equipo).
- **Antiguedad de la guia o carrier especifico.** No explica nada: una guia de 28
  dias si cancelo antes del 08-08, y otras dos de ENVIA tambien.

## Causa raiz

De EnvioClick: su API acepta la cancelacion, responde `200 OK` y no cancela. Si
no puede cancelar deberia devolver la orden en `not_valid_orders` o un `422`.

Amplificado por Probability: `operations.go` solo trataba como error el caso
`not_valid_orders` no vacio. Cualquier otra respuesta se daba por exito, se
marcaba el envio como `cancelled` y se emitia `shipment.cancelled`, dejando el
fallo invisible para el usuario.

## Correccion

Commit `0088c3ed` en `main`.

`envioclick/internal/app/operations.go`:

- Tras `CancelBatch` se re-consulta `Track` hasta 3 veces, con 2s entre intentos,
  y solo se da por cancelado si el carrier reporta estado `cancelad*`.
- Sin confirmacion del carrier se devuelve error: el envio **no** se marca
  `cancelled` y el front recibe el fallo por SSE (`PublishCancelFailed`), con el
  mensaje de cancelar manualmente en la plataforma de EnvioClick.
- `not_valid_orders` pasa a devolver un error real, en vez de un
  `CancelResponse{Status:"error"}` que el consumidor no distinguia del exito.

La verificacion corre en el consumidor de cola (el endpoint responde `202`), asi
que los ~7s del peor caso no afectan ningun request HTTP. El timeout de 30s del
cliente es por request, no por operacion.

Tests en `operations_test.go`: se actualizaron los tres que afirmaban el
comportamiento viejo y se agrego
`TestCancel_Error_CarrierAcceptsButShipmentStaysActive`, que reproduce el caso
exacto de hoy (EnvioClick acepta, la guia sigue activa, debe fallar).

## Verificacion

`go build ./...` y `go test ./...` en verde antes del push. Deploy por el
workflow Backend CI/CD (run `31525517938`).

## Pendientes

- **Cancelar en la plataforma de EnvioClick** las cuatro guias que quedaron
  vivas: `034058168864`, `034058168817` (recoleccion 08-12), `034058168475`
  (08-13), `034058187592` (08-14). El fix no las cancela, solo evita que el
  problema se repita en silencio.
- **Reclamo a EnvioClick.** El caso esta armado: misma URL, mismo body, misma
  version de API, codigo nuestro sin cambios desde el 29 de abril, y la respuesta
  cambiada entre el 7 y el 8 de agosto con los bodies de antes y despues.
  Mientras no lo arreglen, toda cancelacion por API va a fallar de forma visible.
- **Fallback singular muerto.** `operations.go` cae a
  `DELETE /shipment/{idShipment}` cuando no hay `idOrder`, pero le pasa el numero
  de guia, y EnvioClick nunca nos devuelve un `idShipment` al generar (solo
  `idOrder` e `idRate`). Ese camino no puede funcionar. Hoy no se dispara porque
  siempre hay `idOrder`.
