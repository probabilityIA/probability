# Una orden con direccion en Zarzal se geocodifico en Guadalajara de Buga y la guia se devolvio

Mystic Rose reclamo que la guia de la orden MYS-0796 se devolvio por "direccion
incorrecta", y que el mapa que veian en su app (Envia) mostraba Buga, no Zarzal
como decia la direccion. Sospechaban que el sistema dejo escribir dos ciudades
en el mismo campo.

## Sintoma

Orden MYS-0796, shipment 45995, guia `034058358010`, carrier ENVIA. Novedad
`Direccion incorrecta/insuficiente`. El texto de la direccion decia
`Cra 9 # 11-36 piso 2 | Barrio: Gonzalo Echeverry | Zarzal (Valle del Cauca)`.

## Diagnostico

`shipment_sync_logs.request_payload` del `generate` mostraba `destination.suburb:
"Zarzal (Valle del Cauca)"` (correcto) pero `destination.daneCode: "76111000"`
— ese codigo DANE es de **Guadalajara de Buga**, no de Zarzal (`76895`).

La orden (`orders.geozone_city_id = 1044`, geozona "GUADALAJARA DE BUGA") tenia
`shipping_lat/lng = 3.90217, -76.295571` — a ~5 km del centro de Buga y a ~55 km
del centro real de Zarzal. `orders.shipping_city` y `addresses.city` siempre
dijeron "ZARZAL" correctamente; el problema nunca fue el texto.

`ResolveOrderGeozone` (`orders/internal/infra/secondary/repository/geozone_queries.go`)
resuelve la geozona por `ST_Contains` sobre esas coordenadas ANTES de intentar
un fallback por texto (que si habria emparejado bien "ZARZAL"), y ese fallback
solo corre si el primero no dejo nada — por eso nunca se ejecuto.

El origen de las coordenadas malas: la orden se creo manualmente (`Orden creada
manualmente` por la usuaria del negocio, `order_history`). En el formulario, el
campo de Direccion (`AddressAutocomplete`, backed por Google Geocoding) fija
`addressCoords` con la sugerencia que se haya clickeado; el campo separado de
"Ciudad y Departamento" (`handleCitySelect`) solo corrige el texto
(`shipping_city`/`shipping_state`) y nunca vuelve a geocodificar. La usuaria
reporto exactamente eso: eligio una sugerencia de direccion que resulto en Buga,
lo noto, y corrigio manualmente la ciudad a Zarzal — pero las coordenadas de
Buga quedaron pegadas.

## Hipotesis descartadas

- **El sistema dejo escribir dos ciudades en el campo de texto.** No: el texto
  final siempre decia "Zarzal (Valle del Cauca)", nunca aparecio "Buga" en
  ningun campo de texto guardado.
- **Bug de EnvioClick/el carrier.** No, el carrier hizo exactamente lo que se le
  pidio: el `daneCode` que le llego (Buga) fue el que nuestro propio sistema le
  mando.

## Correccion

`handleCitySelect` (`OrderForm.tsx`) ahora vuelve a geocodificar
(`GET /geocode?address=...&city=...`) cuando el usuario corrige la ciudad a
mano, y reemplaza `addressCoords` con el resultado nuevo (o las limpia si la
geocodificacion falla, en vez de dejar coordenadas de otra ciudad).

## Alcance / Pendiente

- No se corrigio el dato historico de la orden MYS-0796 (la guia ya se devolvio;
  el reenvio con la direccion correcta se coordina aparte con el cliente).
- No se audito si hay otras ordenes con el mismo patron (ciudad corregida
  despues de elegir una sugerencia de direccion). Si vuelve a aparecer un caso
  similar, buscar por `orders.geozone_city_id` que no corresponda al
  `shipping_city` de texto.
- Se le explico al cliente que el error fue del sistema, no de ellos.
