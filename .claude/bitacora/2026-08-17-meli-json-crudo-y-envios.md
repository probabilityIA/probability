# MercadoLibre: sin JSON crudo, sin direccion y con las notificaciones de envio rotas

Las ordenes de MercadoLibre entraban sin direccion de destino y sin JSON crudo.
Aparecieron tres problemas encadenados: el parseo de `/shipments/{id}/items`
reventaba, el numero de orden guardado es el `pack_id`, y el raw solo se guardaba
al crear la orden, nunca al actualizarla.

## Sintoma

Orden `2000014575744355` (business 46, integracion 254): en el detalle no salia
direccion ni mapa, y el modal de JSON crudo decia "esta orden no tiene datos
crudos guardados".

```sql
SELECT count(*) AS total,
       count(*) FILTER (WHERE COALESCE(shipping_street,'') <> '') AS con_direccion
FROM orders WHERE platform = 'mercadolibre';
-- 8 | 0     <- ninguna de las 8 ordenes de ML tenia direccion
```

`order_channel_metadata` por canal: shopify 46.343, woocommerce 1.827,
jumpseller 2, **mercadolibre 0**.

## Diagnostico

Se reproduce mandando una IPN sintetica al backend de produccion (la firma solo
se loguea, no bloquea):

```bash
curl -X POST http://localhost:3050/api/v1/meli/notification \
  -H 'Content-Type: application/json' \
  -d '{"topic":"orders_v2","resource":"/orders/2000014575744355","user_id":228589645,"attempts":1}'
```

Respuesta en logs: `"fetching order: meli: order not found"`.

Ese numero **no es un order_id, es un pack_id**: `pack.go:19` hace
`merged.ID = packID` al consolidar packs, asi que la orden queda guardada con el
id del pack y despues `/orders/{ese id}` da 404.

Reintentando por el envio (`topic: shipments`, `/shipments/47791295504`):

```
"resolving shipment orders: meli client: parsing shipment items:
 json: cannot unmarshal string into Go struct field .order_id of type int64"
```

`GET /shipments/{id}/items` devuelve `order_id` como **texto**, y el struct lo
esperaba `int64`. Con eso, **toda** notificacion de topic `shipments` de
MercadoLibre venia fallando desde siempre.

Sobre el JSON crudo: el mapper de ML si arma el `ChannelMetadata`
(`order_mapper.go:188`) y el publisher lo serializa. El hueco estaba en el
consumidor: `map_and_save_order.go:34` manda las ordenes que ya existen a
`UpdateOrder`, y ese use case **no guardaba `ChannelMetadata`** (solo el mock de
test lo implementaba). Como las ordenes de ML se actualizan varias veces, el raw
nunca se persistia. Afectaba a todos los canales, no solo a ML.

### Hipotesis descartadas

- **"El mapper de ML no captura la direccion"**: falsa. `order_mapper.go:107`
  arma la direccion completa (calle, ciudad, departamento, zip, lat, lng,
  telefono) desde `shippingDetail.ReceiverAddress`, y el struct de respuesta
  tiene los tags correctos.
- **"`GET /shipments/{id}` esta fallando"**: no hay evidencia de eso. Cuando
  falla se loguea "Failed to fetch shipment detail" y ese log no aparece; ademas
  el tracking y la URL de guia si llegan (`orders.guide_link` con el formato
  exacto del mapper), lo que indica que la llamada responde 200.

## Causa raiz

Tres cosas distintas:

1. `get_shipment_orders.go`: `order_id` parseado como numero cuando ML lo manda
   como texto.
2. `usecaseupdateorder`: no persistia el JSON crudo del canal.
3. La direccion del comprador **no viene en el JSON de la orden**. El raw de ML
   solo trae `"shipping": { "id": 47791295504 }`; la direccion vive en
   `/shipments/{id}` (`receiver_address`), que es otra llamada.

## Correccion

- `get_shipment_orders.go`: tipo `flexibleInt64` que acepta `order_id` como
  numero o como texto (commit `d01854f8`, desplegado).
- `usecaseupdateorder/update_order.go`: nueva `saveChannelMetadata`, llamada
  **antes** del early-return de "no hay cambios" (si no, una orden sin cambios de
  campos nunca guardaria su raw). Marca el anterior con `is_latest = false` e
  inserta el nuevo. Requirio `MarkChannelMetadataNotLatest` en el puerto y el
  repositorio. Aplica a todos los canales.

## Verificacion

Tras el deploy, misma IPN de envio:

```
INF MercadoLibre IPN notification received  resource=/shipments/47791295504
INF Order published to queue successfully   order_number=2000014575744355
INF MercadoLibre order published            order_id=2000017981110684
INF MercadoLibre order published            order_id=2000017981120028
```

Ya resuelve las dos ordenes hijas del pack (antes moria en el parseo), y el JSON
queda guardado:

```sql
SELECT channel_source, is_latest, length(raw_data::text)
FROM order_channel_metadata ocm JOIN orders o ON o.id = ocm.order_id
WHERE o.order_number = '2000014575744355';
-- mercadolibre | false | 3249
-- mercadolibre | true  | 3217
```

La direccion **sigue vacia**: el raw de la orden confirma que ML solo manda
`shipping.id`.

## Pendientes

- **Direccion de destino**: falta comprobar que devuelve `/shipments/{id}` para
  este vendedor. Sospecha principal: ML protege los datos del comprador
  (Mercado Envios Full) o falta permiso en la app. Para cerrarlo hace falta
  guardar tambien el JSON del shipment, no solo el de la orden.
- **`pack_id` como numero de orden**: agrupar el pack esta bien, pero se pierde
  el order_id real y no se puede volver a consultar la orden por API. Conviene
  conservarlo (en `external_id` o en metadata).
- **Token de ML**: durante la prueba salio `"meli: access token expired"` en
  `enrichBillingInfo`. Revisar el refresco del token de la integracion 254.
