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

## Cierre (2026-08-18)

Se guardo tambien el JSON de `GET /shipments/{id}` (adjunto al raw de la orden
bajo `shipment_detail`) y con eso quedo a la vista la causa real de la falta de
direccion: **ML no devuelve `receiver_address`, devuelve `destination`**.

```json
"destination": {
  "receiver_name": "Jose de Jesus Mejia Rodriguez",
  "receiver_phone": "XXXXXXX",
  "shipping_address": {
    "address_line": "CALLE 6 15-03", "street_name": "CALLE 6", "street_number": "15-03",
    "city": { "name": "Aguachica" }, "state": { "id": "CO-CES", "name": "Cesar" },
    "zip_code": "205010", "latitude": 8.308608, "longitude": -73.6192217,
    "neighborhood": { "name": "Olaya Herrera" }
  }
}
```

El struct esperaba el formato viejo (`receiver_address`), asi que
`ReceiverAddress` quedaba siempre nil y la orden entraba sin destino. La
hipotesis de que ML ocultaba la direccion por privacidad era **falsa**: lo unico
enmascarado es el telefono (`"XXXXXXX"`, ahora se descarta en vez de guardarse).

Corregido leyendo `destination.shipping_address` con fallback al formato viejo,
con test sobre el payload real.

### Packs duplicados

Al sincronizar aparecio el segundo problema: un carrito de MeLi es **una orden
por producto** agrupadas por `pack_id`, y las dos rutas lo guardaban distinto.

| Ruta | external_id | Resultado |
|---|---|---|
| Webhook | `pack_id` (consolidado) | 1 orden con N items |
| Sync | `order_id` de cada hija | N ordenes de 1 item |

Dos claves distintas para la misma compra, asi que quedaba duplicada. En Viga
Sport: 2 packs + sus 5 hijas, **$287.900 contados dos veces**.

Se unifico: el sync consolida igual que el webhook y publica el pack una sola
vez; los ids reales de las hijas quedan en el raw (`pack_order_ids`); y la orden
lleva `channel_pack_id` con badge **PACK** en el listado y el detalle.

### Estados sin mapear

`order_status_mappings` no tenia **ninguna** fila para MercadoLibre
(integration_type_id 3), mientras Shopify tenia 20 y WooCommerce 10. Sin mapeo,
el fallback busca "paid" en `order_statuses` (que no existe) y la orden quedaba
con `status_id` NULL: la lista mostraba el string crudo en gris. Se sembraron
los 8 estados con el criterio de los demas canales.

### Etiqueta de envio

El `guide_url` que guardamos es la URL de la API de ML, que exige el token del
vendedor: abrirla en el navegador da `{"status":401,"message":"Invalid token"}`.
El modulo de envios ya lo resolvia con `guideHref` -> `/internal/meli-label/{id}`,
pero el listado de ordenes abria la URL cruda. Se movio el helper a
`shared/utils` y ahora ambos lo usan.

### Limpieza de datos (produccion)

- Se borraron 333 ordenes de ML del business 46 que importo un sync historico
  (dejando las 10 recientes) y luego las 5 hijas de los 2 packs.
- Quedaron 5 ordenes, todas con direccion, coordenadas y estado mapeado.
- El borrado fue **fisico**: `orders.DeletedAt` es `*time.Time` y no
  `gorm.DeletedAt`, asi que GORM no filtra por `deleted_at` y el soft delete no
  ocultaba nada en la UI.

## Pendientes

- **`orders.DeletedAt` es `*time.Time`** (`migration/shared/models/order.go:17`):
  el soft delete no oculta nada porque ninguna consulta agrega
  `deleted_at IS NULL`. Cambiarlo a `gorm.DeletedAt` afecta a todas las
  consultas de ordenes, hay que hacerlo con pruebas.
- **Token de ML**: durante las pruebas salio `"meli: access token expired"` en
  `enrichBillingInfo`. Revisar el refresco del token de la integracion 254.
