# Alerta: MercadoLibre - stock de publicaciones con variantes

Fecha: 2026-08-12
Modulo: `back/central/services/integrations/ecommerce/meli`
Reportado por: negocio Viga ropa deportiva (business_id 46, integracion 254)

## Contexto

El modulo de ML solo leia y escribia stock a nivel de publicacion. Las
publicaciones con variantes (talla/color) quedaban fuera del sync: no se
mapeaban, no recibian push, y su stock nunca se actualizaba.

Medicion en la cuenta de Viga (seller 228589645) al 2026-08-12:

- 300 publicaciones: 262 sin variantes, 38 con variantes.
- Mapeadas: 252, TODAS sin variantes.
- Las 38 con variantes suman ~1.039 variantes y 6.278 unidades.
- 34 variantes activas publicaban stock que aqui esta en 0 (84 unidades
  fantasma) -> ventas que se caen. Caso real: SKU SHD239-2XL, orden ML
  #2000014491203861 del 12-ago 05:52, cancelada por falta de stock.
- 40 variantes activas estaban en 0 en ML teniendo 256 unidades aqui.

## Como quedo resuelto

El commit `7156dbec` (SanCam04) implemento el push por variante contra
`PUT /items/{id}/variations/{variation_id}`. Ese mecanismo se conservo: es mas
simple que mandar el array completo a `PUT /items/{id}` y no tiene el modo de
fallo catastrofico de esa via (**las variantes ausentes del array se borran**,
confirmado en la doc de ML).

Encima se corrigieron y agregaron:

1. **El descubrimiento no traia los SKU.** `multigetItems` pedia
   `attributes=...,variations`. Con esa proyeccion ML devuelve las variaciones
   SIN su array `attributes` y con `seller_custom_field` en null, asi que
   `extractSKUFrom` devolvia siempre "" y `itemToProducts` las descartaba todas:
   las 38 publicaciones con variantes producian CERO filas. Corregido usando
   `include_attributes=all`. Verificado contra la cuenta real: 1.300 filas
   (262 simples + 1.038 variantes).

2. **No mandar 0 cuando no hay registro de inventario.** `GetStockForProducts`
   devuelve un mapa y lo ausente salia 0, apagando la publicacion en ML. Ahora
   se salta y se conserva el stock que ML ya tiene. Afecta a 230 items de Viga.

3. **Grupo "Sin SKU en el canal" en el comparativo.** Las variantes sin
   `SELLER_SKU` no se pueden emparejar y desaparecian en silencio. Ahora salen
   como grupo propio con banner explicando que el match es por SKU. En Viga son
   320 filas en 27 publicaciones.

4. **Agrupacion visual por publicacion.** `parent_ref`, `parent_label` y
   `variant_label` viajan hasta el front, que pinta una barra de color y un chip
   con el nombre del aviso. Sin eso el comparativo pasa de 300 a 1.300 filas
   sueltas e ilegibles.

5. **Indice unico.** Era `(product_id, integration_id)`, o sea un producto solo
   podia guardar una direccion por integracion; los 35 SKUs de Viga que estan en
   dos publicaciones (prenda suelta + combo) solo mapeaban una. Migrado a
   `(product_id, integration_id, COALESCE(external_variant_id, ''))`. El COALESCE
   es necesario: con NULL crudo Postgres los trata como distintos y se perdia la
   unicidad para los canales sin variantes. Se quitaron los tags `uniqueIndex`
   del modelo. `UpsertProductIntegrationMapping` tambien busca por variante.

## Fulfillment

Probado contra la cuenta real el 2026-08-12: el endpoint por variante funciona.
`PUT /items/MCO1419041847/variations/180397877404 {"available_quantity":1}` ->
200, cambio solo esa variante, las otras 7 intactas, ninguna borrada.
Restaurado a su valor original, tambien 200.

En publicaciones de fulfillment ML lo rechaza:

    variations.available_quantity.not_modifiable
    "Available quantity is not modifiable in fulfillment items"

Es correcto: ahi el stock lo administra ML en su bodega. Por eso se guarda
`external_logistic_type` en `product_business_integrations` al descubrir los
productos, y el sync salta esas publicaciones sin gastar una llamada para
fallar. El valor se refresca solo: `processItemNotification` (topics `items`,
`items_prices`, `stock-locations`) ya hacia un `GetItem`, y ahora actualiza el
tipo logistico ahi mismo. Si el vendedor saca una publicacion de fulfillment,
el webhook lo trae y vuelve a entrar al sync sin intervencion.

## Deteccion de SKU cambiado

El id de la variante nunca cambia, asi que si el vendedor le edita el SKU a una
talla el mapeo sigue "funcionando": apunta a la misma variante, pero esa variante
ya es otro producto. Le seguiriamos empujando el inventario del producto viejo,
sin ningun error.

Se detecta en el comparativo comparando `external_sku` (la copia al momento de
mapear) contra el SKU que la variante tiene hoy. Sale como grupo propio
`sku_changed`, con el SKU viejo y el nuevo, para que el negocio lo corrija. NO se
desasocia ni se re-empareja solo: es una decision del cliente, no una que
podamos adivinar.

## Pendiente

### IMPORTANTE

1. **Normalizar espacios interiores en el SKU.** 38 variantes activas de Viga
   traen `BH313 - 2XL` y aqui es `BH313-2XL`. `productmatch.Normalize()`
   (`shared/productmatch/productmatch.go:45`) solo hace `ToLower` + `TrimSpace`.
   Usar un normalizador propio de SKU, no tocar el compartido: lo usan Woo,
   Shopify y Siigo.

2. **El push en tiempo real gasta una llamada por variante.**
   `pushEcommerceStock` publica un mensaje por fila de mapeo. Si un sync de
   Siigo cambia 400 variantes son 400 PUT en serie; con el rate limit del
   cliente (100 req/min, `newRateLimiter(100)`) son ~4 minutos con la cola
   tapada. Agrupar por publicacion en una ventana corta lo reduce mucho.

3. **El sync masivo tampoco compara antes de empujar.** Manda un PUT por cada
   item mapeado aunque la cantidad ya coincida. En Viga eso son ~920 PUT por
   corrida, ~9 minutos. Saltar los que ya coinciden requiere conocer la cantidad
   actual en ML, que el reconcile ya trae.

### DESEABLE

4. **Concepto de combo/pack.** `MCO1434191343` "Pack X 3" vende 3 unidades por
   venta; empujarle el stock crudo del SKU siempre va a estar mal (12 unidades
   de JH348-4XL son 4 packs, no 12). Mientras no exista el concepto, considerar
   dejar los packs fuera del sync.

## Fulfillment

Probado contra la cuenta real el 2026-08-12: el endpoint por variante funciona.
`PUT /items/MCO1419041847/variations/180397877404 {"available_quantity":1}` ->
200, cambio solo esa variante, las otras 7 intactas, ninguna borrada.
Restaurado a su valor original, tambien 200.

En publicaciones de fulfillment ML lo rechaza:

    variations.available_quantity.not_modifiable
    "Available quantity is not modifiable in fulfillment items"

Es correcto: ahi el stock lo administra ML en su bodega. Por eso se guarda
`external_logistic_type` en `product_business_integrations` al descubrir los
productos, y el sync salta esas publicaciones sin gastar una llamada para
fallar. El valor se refresca solo: `processItemNotification` (topics `items`,
`items_prices`, `stock-locations`) ya hacia un `GetItem`, y ahora actualiza el
tipo logistico ahi mismo. Si el vendedor saca una publicacion de fulfillment,
el webhook lo trae y vuelve a entrar al sync sin intervencion.

## Deteccion de SKU cambiado

El id de la variante nunca cambia, asi que si el vendedor le edita el SKU a una
talla el mapeo sigue "funcionando": apunta a la misma variante, pero esa variante
ya es otro producto. Le seguiriamos empujando el inventario del producto viejo,
sin ningun error.

Se detecta en el comparativo comparando `external_sku` (la copia al momento de
mapear) contra el SKU que la variante tiene hoy. Sale como grupo propio
`sku_changed`, con el SKU viejo y el nuevo, para que el negocio lo corrija. NO se
desasocia ni se re-empareja solo: es una decision del cliente, no una que
podamos adivinar.

## Pendientes de datos del negocio (no son de codigo)

- Viga debe ponerle SKU a las 44 variantes de `MCO2359713856` "Buso Cuello
  Redondo" (361 unidades) y a las 5 de `MCO1493700221` "Boxer Licrados X 3".
- 21 SKUs que ML publica no existen como producto aqui (18 unidades):
  `BT11-3XL`, `BT18-5XL`, `BT116-5XL`, `SHD250-4XL`, `SHD239-6XL`, y uno
  truncado, `FS241-`.
- `MCO2348160328` tiene la variante `BT116-3XL` duplicada (var 180381210164 en
  0 y var 180381210166 en 3). Se corrige en ML.

## Riesgo aparte, mitigado

15 productos `SDH*` creados por el product sync el 2026-08-12 17:09 con stock 0
y ya mapeados. El SKU real en Siigo es `SHD*`; el vendedor escribio las letras
al reves en ML (`SELLER_SKU = SDH22-S` en `MCO3957168204`).

MITIGADO: los 15 no tienen ninguna fila en `inventory_levels`, asi que con el
punto 2 el sync los salta en vez de mandarles 0. Ya no apagan las ~133 unidades
reales. Igual hay que limpiarlos: ensucian el catalogo y el comparativo. Si
alguien les crea un nivel en 0, vuelve el riesgo.

## Orden para destapar a Viga

1. Correr el comparativo y asociar las variantes (sin mapeos el push no tiene a
   que apuntarle).
2. Recien ahi, el sync de inventario.

## Criterio para cerrar

- Una publicacion con variantes queda mapeada variante por variante y su stock
  coincide con Probability tras un sync, verificado en produccion.
- El push no gasta una llamada por variante.
