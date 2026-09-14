# 2026-09-14 - cod_total mezclaba la comision del carrier en ordenes del plugin de WooCommerce

**Ticket:** TKT-000088 | **Business:** 46 (Viga ropa deportiva) y 26 (Demo) | **Canal:** WooCommerce con plugin de cotizacion

## Resumen

`orders.cod_total` significaba dos cosas: en ordenes manuales producto + envio
sin comision del carrier; en ordenes cotizadas por el plugin, el total del
checkout con la comision adentro. Quien no miraba `cod_includes_shipping`
sumaba la comision otra vez. Se separo la comision en su propia columna.

## Sintoma

Orden 15789 (Interrapidisimo, guia 240061196214):

| Donde | Mostraba | Real |
|---|---|---|
| WhatsApp de guia "Valor a recaudar" | 78.732 | 73.367 |
| Detalle "Cobro al cliente final" | 78.732 | 73.367 |
| Detalle "Total para el negocio, sin comision" | 73.367 | 68.002 |
| Modulo contra entrega "Total cliente" | 78.732 | 73.367 |
| PDF de Interrapidisimo "Valor a cobrar" | 73.367 | 73.367 |

La confirmacion de pedido salio bien (73.367) solo porque se envio 5 s antes
de que existiera el envio con la comision.

## Diagnostico

- PDF descargado del `guide_url`: "Valor a cobrar: $ 73.367".
- `shipment_sync_logs`: sondeos 68.002 -> 73.048 -> 73.367, codValue final 73.367.
- Datos: `cod_total` 73.367 = 45.000 + 23.001 + 5.366 (la comision del checkout).
- Billetera: 23.001 debitado = flete 16.500 + seguro 1.101 + margen 3.900 + margen COD 1.500. Correcto.
- 27 guias del plugin de Viga comparadas contra su PDF: 25 cuadran con la tienda.

### Hipotesis descartadas

1. **"A Viga se le paga mal en los cortes"**: falso. Los cortes 55 y 56 tienen
   21 ordenes, todas manuales, pagadas exacto a `cod_total` neto. Ninguna orden
   del plugin habia entrado a un corte.
2. **"Se le cobro de mas en la billetera"**: falso. 27 de 27 guias con un solo
   debito igual a `total_cost`; la comision no pasa por la billetera.
3. **"Es solo visual"**: casi. Ademas de pantalla y WhatsApp, las guias
   14666 y 14668 (del 29/08, antes de `ea633276`) declararon la comision doble.

## Causa raiz

1. `cod_total` sin un significado unico; el flag que lo distingue
   (`cod_includes_shipping`) se hereda de la configuracion de la tienda y vale
   `true` por defecto (`GetIntegrationCodIncludesShipping`), no de cada orden.
2. `ea633276` (31/08) arreglo solo la guia; `33d5ae43` (28/07) hacia que
   WhatsApp sumara la comision y `54392d58` (15/07) lo mismo en la pantalla,
   cuando todas las ordenes eran manuales.

Por que puede haber diferencia entre comisiones: el checkout cotiza la comision
y la guia la vuelve a medir despues (15788 cotizada con Interrapidisimo 5.238,
guia con Coordinadora 6.116). Regla vigente: el comprador paga lo que vio en la
tienda y el negocio asume la diferencia.

## Correccion

- `7bdb67ba` (parche): pantallas y WhatsApp miran `cod_includes_shipping`.
- `8c547183` (fondo): columna `orders.cod_checkout_carrier_fee`. El mapper de
  WooCommerce la llena solo si la linea de envio trae `quote_id` y `cod = 1`, y
  deja `cod_total` neto. Guia, pago COD, manifiesto, etiquetas, WhatsApp y front
  cobran `cod_total + cod_checkout_carrier_fee`.
- Data: `FixCodTotalCheckoutWoo`, idempotente, corrida en produccion el
  2026-09-14 despues de agregar la columna y desplegar. 63 ordenes (45 Viga, 18
  Demo): 61 restan la comision, 14679 y 15765 ya estaban netas y solo guardan la
  comision. Fuera: 14642, 150 y 151 (tarifa sin contra entrega), 140-142 (no
  cuadran con la tienda) y cualquier orden dentro de un corte.
- Orden de despliegue: DDL -> codigo -> DML. El codigo nuevo necesita la
  columna, y con datos viejos se comporta como el parche, asi que ningun
  momento cobro mal.

## Verificacion

- Local: orden simulada por webhook (copia de la 15789): 68.002 + 5.365;
  detalle 68.002 / 5.365 / 73.367; modulo contra entrega 68.002 / 73.367;
  `order.updated` conserva valores; migracion sobre formato viejo y segunda
  corrida sin cambios.
- Produccion: 62 de 63 ordenes cumplen `cod_total + comision = total de la
  tienda`; 14679 no porque fue editada a mano (productos 94.900 contra 45.000
  en la tienda). API de la 15789: `cod_total` 68.002, `cod_checkout_carrier_fee` 5.365.
- Respaldo de los valores previos de 69 ordenes: guardado fuera del repo.

No verificado de punta a punta: generar una guia real, marcar pago COD y el
envio real de WhatsApp (cubiertos por pruebas unitarias con 15789 y 15788).

## Pendientes

- [ ] Guias 14666 y 14668: declaradas con la comision doble (211.864 y 121.133), siguen sin recoger. Cancelar y regenerar.
- [ ] 14679: valores editados a mano, no coinciden con la tienda. Revisar con Viga.
- [ ] Shopify y ordenes sin plugin quedan con `cod_includes_shipping = true` por defecto: si generan guia, el negocio asume la comision sin saberlo. Decidir la regla.
- [ ] Manifiesto y etiquetas propias imprimen `cod_total` sin comision en ordenes manuales (anterior a este caso).
- [ ] El resumen "Recaudado" del modulo contra entrega ahora suma `cod_total` neto tambien para las ordenes del plugin.
