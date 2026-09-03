# 2026-09-03 - WooCommerce COD: la comision se cotizaba sobre el producto, no sobre el recaudo

**Ticket:** TKT-000071 | **Business:** 46 (Viga ropa deportiva) | **Canal:** WooCommerce 221 (`cod_includes_shipping: true`)
**Solo afecta** ordenes que llegan de WooCommerce con el precio fijado en el checkout.

## Reporte del cliente

"La comision del carrier cambio en los ultimos dias, ahora la orden se genera con
comision menor y perdemos plata. Eso ya lo habiamos solucionado."

## No es una regresion

Los dos fixes del 2026-08-06 siguen en el codigo y funcionando
(`probe.ContentValue = declared` en `cod_value.go`, `FindCODFee` por `rateID`).
La calibracion de la guia corrio bien en la orden reclamada: los sondeos de
`shipment_sync_logs` para el shipment 48509 fueron 56.625 -> 60.994 -> 61.271.

Lo que cambio es quien absorbe el faltante. El fix `ea633276` (2026-08-31, "no
duplicar la comision COD") dejo de cobrarle la diferencia al comprador; desde
entonces la absorbe el negocio, y eso destapo un bug distinto que estaba tapado.

## Causa raiz

`buildWooQuotePayload` (`shipments/.../handlers/woocommerce-shipping-rates.go`)
mandaba `codValue = contentValue = solo el valor de los productos`. EnvioClick
cobra la comision COD sobre **lo que recauda el mensajero**: producto + flete +
seguros + margen + la propia comision. Ademas el seguro se cobra sobre el
`contentValue` declarado, asi que tambien sube.

Es circular y por eso nadie lo habia resuelto en el checkout: la comision depende
del declarado y el declarado incluye la comision. **No diverge**: con
Interrapidisimo cada vuelta suma 5,95% de la anterior, o sea una serie geometrica
con solucion cerrada.

Solo se ve con **Interrapidisimo** (5% + IVA, proporcional). Coordinadora es 2%
con minimo $6.116 y a estos montos siempre pega el minimo, que no cambia con el
declarado: por eso cuadra al peso y el bug quedo invisible durante semanas.

## Numeros reales

Orden 14687 (guia 240060477972, factura FV-2-27 / FEVD 27):

| Concepto | Checkout | Real (guia) |
|---|---|---|
| Comision carrier | 3.677 | 4.645 |
| Seguro | 675 | 920 |
| Costo de la guia | 12.593 | 12.838 |
| Total envio | 16.270 | 17.483 |

Recaudo 61.270 - comision 4.645 - guia 12.838 = **43.787** por un producto de
45.000. Faltan **1.213** (968 de comision + 245 de seguro sobre el mayor
declarado).

Ordenes afectadas desde el 25/08 (las de Coordinadora cierran en 0):

| Orden | Comision checkout | Comision real | Perdida |
|---|---|---|---|
| 14670 | 3.677 | 4.827 | 1.150 |
| 14685 | 11.501 | 13.528 | 2.027 |
| 14687 | 3.677 | 4.645 | 968 |
| 14689 | 3.677 | 5.121 | 1.444 |

El flujo manual (`platform`, `cod_includes_shipping=false`) **no esta afectado**:
39 ordenes revisadas en el mismo periodo, desfase promedio $16 (redondeo del
`math.Ceil`). Ahi la calibracion decide el declarado al generar la guia y puede
subirlo libremente.

## La factura no era el problema

FV-2-27 se timbro el 03/09 a las 02:08, en el mismo instante en que entro la
orden; la guia se genero 14 horas despues, a las 11:05. Cuando se facturo, la
comision real no existia. La linea "Envio" (16.270) cuadra al peso con lo que el
mensajero le cobra al comprador, que es lo correcto: no se puede facturar 17.483
contra un pago de 61.270 (esa validacion `sum(payments) == total` es la que
revento facturas reales en `8c67a376`). Corregir el checkout arregla la factura
sola.

## Correccion

1. **Calibracion en el checkout** (`woocommerce-cod-calibration.go`, nuevo):
   segunda cotizacion de sondeo con `codValue = contentValue = V2`, ajuste afin
   por tarifa de la comision y del costo (flete + seguros + margen COD) contra el
   declarado, y solucion del punto fijo
   `D = (productos + a_fee + a_costo) / (1 - b_fee - b_costo)`.
   No se quema ninguna formula de transportadora: se **mide**, igual que
   `ResolveCODValue`. Sirve para tarifa proporcional, plana y con minimo.
   Si el sondeo falla o la tarifa no aparece en el, esa tarifa queda como hoy.
   Para la 14687 el checkout habria cobrado envio 17.580 / comision 4.723.
2. **Switch "Facturar solo con guia generada"** (filtro `require_guide`), pedido
   por el usuario como red de seguridad: la factura no se emite al entrar la
   orden sino cuando la guia existe. Al generarse la guia, el consumer de
   respuestas de shipments publica `order.guide_generated` en
   `orders.events.invoicing` y el `OrderConsumer` de facturacion redispara.
   Es seguro repetir: `InvoiceExistsForOrder` corta el duplicado.

Tests: `woocommerce_cod_calibration_test.go` (Interrapidisimo y Coordinadora en
4 montos, el caso 14687 y el sondeo vacio) y `require_guide_test.go`.

## Hipotesis descartadas

- **"Volvio el bug del 06/08"**: no. El codigo esta intacto y los sondeos de la
  guia se ven correctos en `shipment_sync_logs`.
- **"EnvioClick cambio su tarifa"**: no. `floor((840 + 0,05*D) * 1,19)` reproduce
  al peso las comisiones de todas las ordenes revisadas.
- **"Es un bug de facturacion / del mapper de Siigo"**: no. La factura es fiel a
  lo recaudado; lo que esta mal es lo recaudado.
- **"Es un problema de tiempos, facturamos antes de la guia"**: solo explica por
  que la factura no podia saber la comision real. Facturar despues no recupera un
  peso, porque el monto se lo prometio la tienda al comprador en el checkout.

## Verificacion en produccion (2026-09-03, business 26 Demo)

Desplegado (commit `3312a469`, Backend y Frontend CI/CD en verde). Probado con la
tienda de pruebas `woo.probabilityia.com.co` contra produccion, con guias
**reales** a EnvioClick, y la facturacion Siigo apuntando al mock
(`back-testing:9095`, integracion 198 con `is_testing: true`).

Orden #150, Interrapidisimo, guia real 240060488806:

| | |
|---|---|
| Checkout | flete 6.518 + seguro 1.384 + margen 1.500 + comision 4.505 = 13.907 -> cobra 58.907 |
| Sondeos de la guia | 54.402 -> 58.638 -> 58.906 |
| Guia | declarado 58.906, comision real 4.504, costo 9.402 |
| Neto | 58.906 - 4.504 = 54.402 = 45.000 + 9.402 (exacto) |
| PDF Interrapidisimo | "Valor a cobrar: $58.906" |
| Factura mock | FE-24446-1001, Envio 13.907, pago 58.907 |

La comision cotizada en el checkout (4.505) y la cobrada por el carrier (4.504)
difieren en **1 peso**. Con el codigo anterior habrian sido 3.677 contra 4.504.

Orden #151, Coordinadora, guia real 84151734575: checkout 15.261 -> cobra 60.261;
sondeos 54.145 -> 60.261 (pendiente 0, pega el minimo); neto 60.261 - 6.116 =
54.145 = 45.000 + 9.145 exacto; PDF "Valor a recaudar: 60,261"; factura
FE-24446-1002.

Wallet debitado 18.547 = 9.402 + 9.145, exacto. Las dos guias quedaron canceladas
en EnvioClick.

**Switch `require_guide`:** #150 entro antes de activarlo y se facturo a los 3
segundos sin guia; #151 entro con el switch activo, quedo sin factura, y la
factura salio sola **38 segundos despues de generarse la guia**, disparada por
`order.guide_generated`.

Latencia del checkout: 8 s en contra entrega (dos cotizaciones) contra 5,5 s en
prepago. El sondeo agrega ~2,5 s.

### Hallazgos secundarios de la prueba

- **`retry-guide` no calibra.** `POST /shipments/quotes/:id/retry-guide` ->
  `buildRetryPayload` no setea `codNetTarget`, asi que `applyCODCalibration`
  no corre y se declara el `codValue` crudo de la cotizacion (el valor de los
  productos). Ademas manda `external_order_id` con el UUID de 36 caracteres y
  EnvioClick rechaza mas de 28. Ese camino esta roto para contra entrega.
- **Cancelar una guia no devuelve el saldo al wallet** (los 18.547 siguen
  descontados de Demo). Relacionado con
  `.claude/alerts/wallet-cobro-guias-no-atomico.md`.
- **`pickupDate` es obligatorio** en el payload de generacion; sin el EnvioClick
  responde 422 "el campo debe ser de tipo fecha AAAA-MM-DD".

### Configuracion que quedo cambiada en Demo

- Integracion 197 (Woo): `INTERRAPIDISIMO` agregado a `allowed_carriers_cod`.
- Integracion 198 (Siigo mock): `seller_id = 629` (faltaba, y por eso fallaba
  la primera factura).
- Config de facturacion 17: `require_guide: true`, `final_customer_when_no_id: true`.
- EC2 `woo-store` (`i-0e41ea3a2f1747cd3`) encendida.

## Los otros dos puntos de entrada del mismo error (misma fecha)

El usuario probo desde el panel y encontro que **Interrapidisimo no aparece en
las opciones contra entrega**. Es el mismo error de fondo en otro camino, uno
que la primera verificacion no cubrio: se probo el checkout y la generacion de
guia, no el cotizador del panel.

**Cotizador del panel** (`POST /shipments/quote`): mandaba
`codValue = cod_total` y `contentValue = subtotal`. Reproducido contra el API
real, cambiando una sola variable:

| codValue | contentValue | Interrapidisimo contra entrega |
|---|---|---|
| 45.000 | 45.000 | si |
| 60.261 | 45.000 | **no** |
| 60.261 | 60.261 | si |

Es la causa raiz del 2026-08-06 aplicada al panel. Corregido con
`alignCODContentValue`, que iguala `contentValue` a `codValue` en contra
entrega, que ademas es lo que la guia termina declarando.

**Retomar cotizacion** (`POST /shipments/quotes/:id/retry-guide`):
`buildRetryPayload` no seteaba `codNetTarget` ni pasaba por `overrideCodValue`,
asi que `applyCODCalibration` no corria y se declaraba el `codValue` crudo de la
cotizacion. Tampoco llevaba `totalCost` ni `codCarrierFee`, que salen de la
tarifa elegida. Y mandaba `external_order_id` con el UUID de 36 caracteres
cuando EnvioClick rechaza mas de 28. Todo corregido, con tests.

### Falsa alarma verificada: el valor que "no coincidia"

En la orden 151 el panel mostraba Coordinadora en 10.700 contra los 9.145 del
checkout. No es error de calculo: **el flete cambio de 6.391 a 8.550 en hora y
media**. Comprobado cotizando con codValue de 45.000, 60.261 y con los dos
iguales: el flete da 8.550 en los tres casos, asi que no depende del monto
declarado. Es la deriva de tarifa del proveedor, que asume quien fija el precio
en el checkout.

## Prueba end to end desde el navegador (2026-09-03)

Compra real en `woo.probabilityia.com.co` con Playwright, replicando las
condiciones de la orden 14687 de Viga: producto de 45.000, Envigado,
Calle39Dsur #25-50, contra entrega, Interrapidisimo.

| | |
|---|---|
| Checkout mostro | Subtotal 45.000, Entrega 13.907, Total 58.907 |
| Orden 158 en Probability | cod_total 58.907, comision 4.505, envio 9.402 |
| Facturas al entrar | 0 (filtro require_guide) |
| Sondeos de la guia | 54.402 -> 58.638 -> 58.906 |
| Guia real | 240060500174, comision real 4.504 |
| Neto al negocio | 58.906 - 4.504 = 54.402 = 45.000 + 9.402 (exacto) |
| PDF Interrapidisimo | "Valor a cobrar: $58.906" |
| Factura | FE-24446-1003, emitida 0,1 s despues de la guia, envio 13.907 |

Con el codigo anterior esa orden habria cobrado 3.677 de comision y el negocio
habria perdido unos 830 pesos.

Detalles que costaron tiempo en montar la prueba:

- Los productos de la tienda estaban en `outofstock` y WooCommerce los pinta
  como "Leer mas": no se pueden agregar al carrito. Hay que reponer stock por
  la REST API y devolverlo despues.
- El plugin cotiza contra entrega **solo cuando cambia el metodo de pago**. Al
  entrar al checkout cotiza prepago; hay que alternar a otro metodo y volver a
  contra entrega para que dispare la cotizacion COD.
- El bloque de checkout limpia los campos de envio cuando se re-renderiza por
  el cambio de departamento: hay que rellenarlos despues.

## Pendiente

- [x] Desplegar (`3312a469`).
- [ ] Las 3 ordenes del 03/09 (14685, 14687, 14689) siguen sin recaudar. Se
      pueden cancelar y regenerar despues del deploy para no perder los $4.439.
- [ ] FV-2-27 ya esta ante la DIAN: nota credito o factura adicional por $1.213.
- [ ] Revisar el margen COD de Viga en WooCommerce ($1.500) contra el del flujo
      manual ($5.400): es el que deberia absorber la deriva de tarifa entre el
      checkout y la guia.
- [ ] Decidir si se activa `require_guide` en Viga (business 46).
- [x] `retry-guide` sin calibracion COD ni `external_order_id` acotado (corregido
      y desplegado el 2026-09-03, commit `39c7f007`).
