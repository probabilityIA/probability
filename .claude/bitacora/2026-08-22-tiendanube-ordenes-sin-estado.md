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

## Push-back de estado y guia hacia Tiendanube (implementado el mismo dia)

Faltaba la vuelta completa: el cliente de ordenes de Tiendanube solo leia
(`GetOrder`, `GetOrders`). Cuando Probability generaba la guia, el comprador
seguia viendo su pedido "por empaquetar" en la tienda y el correo de seguimiento
que Tiendanube manda solo nunca salia.

La mitad ya existia: el modulo de ordenes publica cada cambio de estado al
exchange fanout `orders.events`, y **Jumpseller y MercadoLibre ya se
suscribian**. Tiendanube no.

Lo implementado:

- `shared/rabbitmq`: cola `orders.events.tiendanube`.
- `tiendanube/internal/infra/primary/queue/order_status_consumer.go`: consume
  `order.status_changed`, filtra por plataforma y llama al usecase.
- `client/fulfillment.go`: `ListFulfillmentOrders`, `UpdateFulfillmentOrder`
  (PATCH), `CreateTrackingEvent` (POST) y `CancelOrder`.
- `usecases/update_status.go`: homologacion de estados y progresion.

Endpoints usados (doc: https://tiendanube.github.io/api-documentation/resources/fulfillment-order):

```
GET   /orders/{id}/fulfillment-orders
PATCH /orders/{id}/fulfillment-orders/{fo_id}     {"status":"DISPATCHED","tracking_info":{"code","url","notify_customer"}}
POST  /orders/{id}/fulfillment-orders/{fo_id}/tracking-events
POST  /orders/{id}/cancel
```

Homologacion:

| Probability | Fulfillment | Evento de seguimiento |
|---|---|---|
| picking, packing, ready_to_ship | PACKED | - |
| assigned_to_driver, picked_up, shipped | DISPATCHED | dispatched |
| in_transit | DISPATCHED | in_transit |
| out_for_delivery | DISPATCHED | out_for_delivery |
| delivered | DELIVERED | delivered |
| delivery_novelty, delivery_failed | DISPATCHED | delivery_attempt_failed |
| rejected, return_in_transit, returned | DISPATCHED | returned_to_sender |
| cancelled | - | `POST /orders/{id}/cancel` |

Detalles que importan:

- **La API no deja saltar pasos**, asi que `pasosHasta` aplica los intermedios en
  orden: una orden en UNPACKED que pasa a `in_transit` recibe PACKED y luego
  DISPATCHED. Nunca retrocede: si ya esta DISPATCHED y llega una novedad, solo se
  registra el evento.
- El `tracking_info` viaja en el paso **DISPATCHED** con `notify_customer: true`,
  que es lo que dispara el aviso de Tiendanube al comprador. Sin numero de guia no
  se manda el objeto.
- El link de rastreo no venia en el evento: se agrego `tracking_url` al
  `OrderSnapshot` del modulo de ordenes, tomando `tracking_link` y cayendo a
  `guide_link`. Eso tambien le sirve a Jumpseller y MercadoLibre.
- Apagado por defecto: requiere `status_sync_enabled` en el config de la
  integracion, con interruptor nuevo en el formulario de edicion.
- Un 404 al listar fulfillment orders se trata como "no hay nada que actualizar",
  no como error.

Tests en `update_status_test.go`: progresion de pasos, no retroceso, sync
apagado, estado sin homologar, transito con evento, entrega, novedad sin patch,
cancelacion, orden sin guia y orden sin id externo.

**Sin probar contra la API real**: hace falta desplegar, activar el interruptor
en la integracion 264 y mover una orden de la tienda de prueba por picking ->
in_transit -> delivered, comprobando en el admin de Tiendanube que el pedido
avanza, muestra la guia y lista los eventos de seguimiento.

## Dos interruptores por integracion para el estado de las ordenes

El primer intento fue una regla global: una orden que ya estaba en `picking` o
mas adelante ignoraba lo que dijera el canal. No sirve, porque cada negocio opera
distinto:

- **sin intermediarios** solo usa Probability para facturar; el estado siempre lo
  manda el canal.
- **Viga ropa deportiva** centraliza la operacion aca.
- **MercadoLibre** es quien notifica el avance del envio.
- **WooCommerce** se notifica desde Probability.

Asi que la direccion del flujo es configuracion, no regla fija. Cada integracion
tiene ahora dos interruptores, **ambos encendidos por defecto**:

| Config | Sentido | Efecto |
|---|---|---|
| `status_inbound_enabled` | canal -> Probability | apagado: el canal no mueve el estado operativo |
| `status_sync_enabled` | Probability -> canal | apagado: no se escribe nada en el canal |

Con la entrada apagada, `updateOrderStatus` no toca `status` ni `status_id`, pero
el `original_status` se sigue guardando y el estado de pago sigue su camino. Las
**cancelaciones y los reembolsos entran siempre**, incluso con la entrada
apagada: si no, se despacha una orden que el comprador ya cancelo en la tienda.

Detalles de implementacion:

- El flag se lee de `integrations.config` con un SELECT replicado en el repo de
  ordenes (`IsChannelStatusInboundEnabled`). Si falta la llave o falla la
  consulta, se asume encendido: nunca se pierde un estado por un error de lectura.
- La salida (`status_sync_enabled`) tambien pasa a encendida por defecto en
  Tiendanube y Jumpseller, que antes exigian activarla a mano.
- UI: componente compartido `ChannelStatusSyncSection` con los dos switches, ya
  montado en Tiendanube, WooCommerce, MercadoLibre, Jumpseller y Shopify. VTEX
  conserva sus `ToggleRow` y solo se le agrego el de entrada.


## Secciones del pedido original en todos los canales

El JSON original **si** se guardaba: vive en `order_channel_metadata` y lo tienen
todos los canales (Shopify 17.677/17.677, Woo 1.001/1.001, Tiendanube 9/9,
Jumpseller 2/2, MELI 120/124, VTEX 0/1 y esa unica es del mock).

Lo que estaba vacio eran las columnas `orders.financial_details`,
`shipping_details`, `payment_details` y `fulfillment_details`, que **solo Shopify
llenaba** con extractores propios (`shopify_details_extractor.go`). No son una
copia del pedido: son las secciones sueltas, y el backend las lee — el mapeo de
`payment_status_id` saca el `financial_status` de `payment_details`, y por eso
estaba hardcodeado a Shopify.

Ahora hay un extractor compartido en `canonical/raw_sections.go`: recibe el JSON
del canal y un perfil de llaves, y arma las cuatro secciones. `Root` cubre a los
canales que envuelven la orden en un objeto (Jumpseller la manda dentro de
`order`). Un payload ilegible devuelve secciones vacias: nunca descarta la orden.

**Las llaves de cada perfil salieron del JSON real guardado en produccion**, no
de la documentacion:

```sql
SELECT o.platform, string_agg(DISTINCT k, ', ' ORDER BY k)
FROM orders o
JOIN order_channel_metadata m ON m.order_id = o.id
CROSS JOIN LATERAL jsonb_object_keys(m.raw_data::jsonb) AS k
GROUP BY 1;
```

Perfiles: `TiendanubeSections`, `JumpsellerSections`, `MeliSections` y
`WooCommerceSections`, conectados en el mapper de cada canal donde ya se armaba
el `ChannelMetadata`. Shopify conserva sus extractores, que son mas detallados.

VTEX tambien quedo conectado, pero es la excepcion: su mapper si arma el
`ChannelMetadata`, lo que pasa es que la unica orden VTEX en produccion es del
mock y entro sin raw. Como no hay un pedido real de donde sacar las llaves, el
perfil `VTEXSections` salio de los tags del response del cliente
(`vtex_order_response.go`) y conviene contrastarlo con el primer pedido real.

Con esto se puede generalizar el mapeo de estado de pago a los demas canales, que
es lo que hoy solo funciona para Shopify.
