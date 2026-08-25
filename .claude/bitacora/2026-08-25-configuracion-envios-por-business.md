# Modulo de configuracion de envios y el paquete que se declara al carrier

Fecha: 2026-08-25
Modulos: back/central `shippingconfig` (nuevo) + `shipments`, front/central `shipping-config` (nuevo)

## Que problema resolvia

La configuracion de empaque y transportadoras vivia en tres lugares que no se
hablaban entre si:

| Que | Donde vivia | Quien la usaba |
|---|---|---|
| Cajas de empaque | `warehouses.metadata` | solo checkout WooCommerce y creacion de orden |
| Transportadoras permitidas | `integrations.config.allowed_carriers_*` (por tienda) | solo checkout WooCommerce |
| Dimensiones al generar guia | el formulario del modal | nadie mas |

Consecuencia medida en produccion: la orden 14636 de Viga (business 46) cotizo en
el checkout con **3 kg / 30x40x30** (Caja Mediana, correcto) y la guia se genero
con **1 kg / 10x10x10**. Cuatro prendas talla 5XL declaradas como una caja de
10 cm. El cliente pago $33.433 de envio y la guia costo $20.292.

De 574 guias en 60 dias, 133 salieron con el default. Por negocio: Viga 85 de 88,
Moto Mello 13 de 13, LaPercha 9 de 9. La causa de fondo es que el catalogo no
tiene medidas (Viga: 0 de 2.266 productos; Moto Mello: 2 de 5.152), pero la caja
estandar de la bodega si estaba configurada y la generacion de guia no la leia.

## Que se construyo

**Tabla `shipping_configs`**: `business_id` obligatorio, `warehouse_id` nullable
(NULL = config del negocio, con valor = override de esa bodega), unica por
`(business_id, warehouse_id)`. Guarda `package_strategy`, `boxes` y `carriers`.

**Modulo `shippingconfig`** (hexagonal) con `GET/PUT /shipping-config`,
`PUT/DELETE /shipping-config/warehouses/:id` y
`PUT /shipping-config/warehouses/:id/default`.

**Transportadoras por negocio**: cada una con `enabled`, `allow_cod`,
`allow_prepaid` y un switch `direct` para convenio propio, con tres estados:
`unavailable` (sale por la cuenta de Probability), `pending` (el cliente activo
su convenio pero el conector aun no existe) y `active` (hay conector e
integracion enlazada). La compuerta es `DirectIntegrationAvailable()`, hoy
siempre false. **La configuracion de transportadoras de WooCommerce NO se movio**:
es del canal y sigue en su integracion.

**`integration_types.is_system_provider`**: marca las integraciones internas de la
plataforma (envioclick, enviame, mipaquete, shipit). Probability las opera para
todos los negocios y el cliente no las ve. Las integraciones directas de
transportadora, cuando existan, iran con el flag en false.

**Resolvedor unico** `shared/shippingpkg`: logica pura, sin base de datos, para
que shipments y shippingconfig la compartan sin violar el aislamiento de repos.
Cascada de dimensiones: caja estandar -> dimensiones del producto -> default
10x10x10. Cascada de peso: carrito -> catalogo -> 1 kg, tomando siempre el mayor
entre contenido y caja.

## Correcciones que trae el resolvedor

1. El peso ya no se pierde: si el carrito no lo trae, cae al `products.weight`.
2. El peso de la caja no aplasta al real: se usa `max(contenido, caja)`.
3. Dimensiones parciales ya no descartan todo: antes, si faltaba un lado, el
   `AND` estricto mandaba 10x10x10; ahora se usan los lados que existan.
4. Los cuatro flujos leen la misma configuracion: checkout Woo, checkout Shopify
   (que mandaba 10x10x10 fijo), cotizacion del panel y generacion de guia.

En generacion de guia y cotizacion del panel solo se sobreescribe el paquete si
el payload venia con el default; si el operador tecleo medidas, se respetan
(caso Mystic Rose, que escribe 33x15x16 a mano en cada guia).

## Verificacion

- Tests de `shared/shippingpkg`: caso Viga 14636, peso de catalogo, dimensiones
  parciales, descarte por tamano, y **equivalencia con el selector anterior** para
  cantidades 1..8 (misma caja elegida, para no cambiar tarifas en silencio).
- E2E local contra un mock del carrier: enviando el default `1 kg 10x10x10` sobre
  una orden de 4 items, el carrier recibio `3 kg 30x40x30`. Con medidas tecleadas
  (33x15x16), se respetaron.
- Navegador: icono en la barra de Envios y nodo "Envios" en Tus Integraciones,
  ambos abren el formulario; guardar persiste en `shipping_configs`.
- `go build ./...`, tests Go del modulo y `vitest` 739/739 en verde.

## Pendientes

1. La migracion de `warehouses.metadata` a `shipping_configs` esta escrita pero
   solo corrio en local (la copia Demo no tiene cajas, migro 0 filas). En
   produccion hay 5 bodegas con estrategia configurada.
2. Retirar la seccion de empaque de `WarehouseForm` recien despues de migrar, para
   no dejar dos fuentes editables.
3. `order_items` duplicados: la orden 14636 tiene 8 filas para 4 prendas (cuatro
   con `product_sku` vacio). Infla `totalQuantity` y puede cambiar la caja elegida.
4. Acumular volumen real en vez de tomar el lado maximo entre items.
