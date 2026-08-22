# 2026-08-22 - Tiendanube: las ordenes entran, pero la pagada queda sin estado

## Resumen

Se creo un producto y dos ordenes de prueba en la tienda demo de Tiendanube
(`sebas y corotos`, store 8126740, integracion 264, business Demo 26). Todo llega
a Probability en ~1 segundo y mapea bien, salvo el estado: una orden que nace
pagada queda con `status = 'paid'` y `status_id = NULL`, o sea sin estado del
catalogo, y no entra al flujo de bodega.

## Lo que si funciono

| Dato | Resultado |
|---|---|
| Llegada por webhook | `order/created`, `order/paid`, `order/updated` procesados, 200 en 0 ms |
| Latencia | orden creada 21:11:19 -> guardada 21:11:20 |
| Cliente | creado/deduplicado con nombre, email, telefono normalizado a `+57...` y DNI |
| Direccion | calle+numero, barrio, ciudad, departamento, codigo postal, pais CO |
| Montos | subtotal 55.000, total 55.000, COP |
| Producto | creado solo en Probability (`PRD_niz-wtEHjZET`, SKU `TN-TEST-001`) |
| Mapping producto-canal | `product_business_integrations` 10351 con `external_product_id 362572707` y `external_variant_id 1581840507` |
| Idempotencia | la 2a orden reuso el mismo producto; eventos repetidos no duplicaron la orden |

Ordenes de prueba: `2051703380` (#100, nacio pagada) y `2051706151` (#101, nacio
con pago pendiente y luego se marco recibida).

## El problema

`MapOrderStatus` (`tiendanube/internal/app/usecases/mapper/status_mapper.go`)
devuelve `"paid"` cuando el pago esta hecho. Ese codigo **no existe** en el
catalogo `order_statuses` (ahi hay pending, processing, shipped, delivered,
completed, cancelled...) ni en la maquina de estados de `orders`.

Al resolver `status_id`, `mapOrderStatusID` intenta tres caminos y falla en los
tres:

1. `statusmapper.GetIntegrationTypeID("tiendanube")` devolvia **0** porque el
   switch no tenia el caso; sin type id no se consulta `order_status_mappings`.
2. Aun con el type id, no habia **ninguna fila** en `order_status_mappings` para
   `integration_type_id = 17`.
3. Fallback por codigo: `"paid"` no esta en `order_statuses`.

Resultado en la orden #100: `status = 'paid'`, `status_id = NULL`. En el listado
de ordenes del API la fila sale con el estado vacio.

La #101 se salvo de casualidad: nacio como `pending`, el consumidor de inventario
la reservo y la movio a `picking` (`status_id = 12`). Al marcarla pagada, el log
lo dice sin ambiguedad:

```
WRN Integración realizó salto de estado que no cumple flujo v2 — aceptado por ser
    fuente externa  from=picking to=paid integration_type=tiendanube
```

Es decir: cada `order/updated` del canal **retrocede** el estado operativo de
`picking` a `paid`. La orden sale del flujo de bodega.

## Correccion aplicada (sin desplegar)

- `statusmapper/integration_type.go`: se agrego `case "tiendanube": return 17`.
- `migration/internal/infra/repository/migrate_tiendanube_order_statuses.go`
  (nuevo, idempotente, registrado en `constructor.go`): siembra los
  `integration_channel_statuses` del canal y los `order_status_mappings` del tipo
  17 (`open`/`paid` -> 2 processing, `closed`/`completed` -> 5, `cancelled`/
  `voided`/`abandoned` -> 6, `pending` -> 1, `refunded` -> 7, `shipped`/
  `in_transit` -> 17, `delivered` -> 4).

Con esto la orden pagada deja de quedar sin estado: pasa a "En procesamiento",
igual que Shopify, Jumpseller y MercadoLibre, que ya tenian su semilla.

## Pendiente de decision (no se toco)

El retroceso de estado sigue existiendo: con el mapeo puesto, un `order/updated`
de Tiendanube pisara `picking` con `processing`. El codigo ya detecta el salto y
lo acepta "por ser fuente externa". Lo correcto seria que una actualizacion del
canal no retroceda una orden que ya avanzo a fase de bodega o ultima milla, salvo
cancelacion o entrega. Afecta a **todos** los canales, no solo a Tiendanube.

## Gaps menores encontrados

- `order/paid` sobre una orden ya existente no resuelve `payment_status_id`
  (quedo NULL en la #101; en la #100, que nacio pagada, quedo en 3).
- Tiendanube no guarda los JSON crudos (`metadata`, `shipping_details`,
  `payment_details`, `fulfillment_details` todos NULL). Solo Shopify los guarda
  completos y WooCommerce guarda `shipping_details`. MercadoLibre y Jumpseller
  tampoco.
- `order_items.product_sku` y `product_name` van vacios. No es de Tiendanube: en
  los ultimos 60 dias estan vacios en los 26.000 items de todos los canales.
- El producto creado desde la orden queda con `track_inventory = false` aunque en
  Tiendanube tenga stock limitado (25 unidades).
- `product/updated` se ignora a proposito en el webhook ("no es un evento de
  orden"): crear o editar un producto en Tiendanube no lo trae a Probability. El
  producto llega solo cuando entra una orden que lo contiene, o por el
  reconcile/apply manual.

## Hipotesis descartadas

- "La orden no llego": llego en 1 segundo y quedo completa; lo que falta es el
  estado.
- "El consumidor de inventario no corre": si corre, movio la #101 a `picking`. La
  #100 no avanzo porque su estado `paid` no es un nodo valido de la maquina de
  estados y no puede transicionar a `picking`.
- "Falta un webhook": los 7 estan registrados y los tres eventos de orden
  llegaron. El de "empaquetado" llega como `order/updated` y no cambia nada,
  porque el shipping status no se refleja en la orden.

## Firma del webhook (corregido el mismo dia)

El webhook de ordenes aceptaba cualquier POST: leia el body, tomaba el
`integration_id` del query string y procesaba. El header
`x-linkedstore-hmac-sha256` que Tiendanube envia en cada request solo se
registraba en los webhooks de privacidad, nunca se verificaba. Con los
`integration_id` siendo secuenciales, cualquiera podia inyectar ordenes en el
negocio de un cliente.

Tiendanube firma con **HMAC-SHA256 en hex sobre el body crudo, usando el
`client_secret` de la app** (mismo esquema que Shopify). Documentacion:
https://tiendanube.github.io/api-documentation/resources/webhook

Implementado en `handlers/webhook_signature.go`:

- `firmaValida` calcula el HMAC sobre los bytes tal cual llegan, antes de
  cualquier `json.Unmarshal`, y compara con `hmac.Equal` (tiempo constante).
- `verificarWebhook` resuelve la integracion, toma el `client_secret` que
  corresponda segun `is_testing` (soporta `test_client_secret`) y ademas valida
  que el `store_id` del payload sea el de esa integracion. Sin esa segunda
  validacion, alguien con su propia tienda y una firma legitima podria apuntar el
  webhook al `integration_id` de otro negocio.
- Si algo no cuadra: `Warn` de una linea y **401**, sin procesar.
- En los webhooks de privacidad, que no escriben nada, la firma se verifica pero
  solo se registra (`firma_valida=true/false`) y se sigue respondiendo 200, para
  no arriesgar el App Review.

Tests en `webhook_signature_test.go`: firma correcta, secret distinto, body
alterado, firma vacia, secret vacio, firma no hexadecimal y espacios en el header.

**Por que el riesgo de cortar ordenes al desplegar es bajo:** el `client_secret`
guardado es el mismo que se usa en el intercambio `code -> access_token`, y ese
intercambio funciono cinco veces seguidas hoy en los ciclos de OAuth. Si
estuviera mal, el OAuth habria fallado. Aun asi, al desplegar conviene crear una
orden de prueba y confirmar en el log que no aparece
"Webhook de Tiendanube rechazado".

El mock de `back/testing` no emite webhooks, asi que no se ve afectado.
