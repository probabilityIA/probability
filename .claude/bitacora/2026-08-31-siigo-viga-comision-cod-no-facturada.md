# 2026-08-31 - Viga: la comision COD del carrier no se factura en Siigo

**Business:** 46 (Viga ropa deportiva) | **Integracion facturacion:** 228 (Siigo) |
**Canal:** WooCommerce 221 (`cod_includes_shipping: true`)

## Reporte del cliente

1. "La comision del carrier no se facturo".
2. "La orden cuando llega de WooCommerce factura doble envio".

## Que se verifico

Se descargaron los PDF reales de Siigo por
`GET /api/v1/siigo/invoices/<id>/pdf?business_id=46` (FV-2-1, FV-2-4, FV-2-5,
FV-2-6, FV-2-7, FV-2-14) y se comparo contra `orders`, `shipments` y
`shipping_details` del canal.

## 1. Comision COD no facturada - CONFIRMADO

`orders.shipping_cost` es el **costo de la guia** (flete + seguro minimo +
seguro extra + margen COD), deliberadamente **sin** la comision del carrier
(`cod_carrier_fee`), segun `order_mapper.go` y
`.claude/rules/guias-contra-entrega.md`.

`create_invoice.go` copiaba ese valor a `invoices.shipping_cost`, y el mapper
de Siigo (`.../siigo/.../mappers/invoice.go`) arma con el la linea "Envio".
Resultado: en contra entrega se factura menos de lo que el cliente paga en la
puerta, por exactamente la comision del carrier.

| Factura | Orden | Envio facturado | Envio cobrado en Woo | Faltante |
|---|---|---|---|---|
| FV-2-7  | 14665 | 12.006 | 18.122 | 6.116 |
| FV-2-8  | 14666 | 23.132 | 29.248 | 6.116 |
| FV-2-9  | 14667 | 20.789 | 25.234 | 4.445 |
| FV-2-10 | 14668 | 18.901 | 25.017 | 6.116 |
| FV-2-12 | 14670 | 15.645 | 19.322 | 3.677 |
| FV-2-13 | 14672 | 28.382 | 34.498 | 6.116 |
| FV-2-14 | 14673 | 19.031 | 25.147 | 6.116 |

Total sin facturar: **$38.702** en 7 facturas.

Ejemplo cerrado, orden 14673: Woo cobra 70.147 (45.000 producto + 25.147 envio,
donde 25.147 = flete 16.501 + seguro min 580 + seguro extra 450 + margen
Probability 1.500 + comision carrier 6.116). `cod_total` = 70.147, que es lo que
recauda el mensajero. La FV-2-14 salio por **64.031**. Diferencia = 6.116.

Solo afecta ordenes COD. En prepago (FV-2-4, FV-2-5, FV-2-6, FV-2-11) no hay
comision y `shipping_cost` coincide con lo cobrado: faltante 0.

## 2. Doble envio - NO se reproduce

Los 6 PDF de Siigo traen **una sola** linea "Envio" y el total es
productos + envio. Ninguna factura de la base tiene dos lineas de envio.

El doble envio real fue el que corrigio `b3611857` (2026-08-27 21:01):
`total_amount` de las ordenes de Woo traia el total del canal (productos +
envio) y el modal de Editar Orden le volvia a sumar `shipping_cost`
(`OrderForm.tsx:567`). Desde ese commit `total_amount` es solo productos.

**Queda dato viejo:** 14 ordenes de Woo del business 46 creadas entre el
2026-08-13 y el 2026-08-27 conservan `total_amount = productos + envio`. En esas
el modal sigue mostrando el envio dos veces. Facturas emitidas de ese lote
(FV-2-1, FV-2-4, FV-2-5) estan bien en el PDF: el pago se calcula sumando
`items[]`, no `req.Total` (fix `8c67a376`).

## Hipotesis descartadas

- Doble linea "Envio" en el request a Siigo: no, el mapper la agrega una sola vez.
- Item de envio duplicado desde `order_items`: no existe ningun item de envio
  en `order_items` del business 46.
- `include_shipping` de `filter_config.go`: esta declarado pero nunca se usa en
  el camino de Siigo; no es la causa.

## Fix aplicado (2026-08-31, sin desplegar)

Se decidio **una sola linea "Envio"** con envio y comision sumados, para no
depender de que el negocio cree un servicio nuevo en Siigo (el catalogo de Viga
solo tiene el codigo `Envio`).

`dtos.InvoiceShippingCost` / `dtos.InvoiceTotalAmount`
(`invoicing/internal/domain/dtos/invoice_data.go`):

    si no es COD, o el canal no registro comision  -> igual que hoy
    si es COD y hay comision                       -> envio  += comision
                                                      total  += comision

La comision NO se deduce: se lee explicita del meta `cod_carrier_fee` que el
plugin de checkout de WooCommerce guarda en `orders.shipping_details`
(`codCarrierFeeFromOrder` en el repo de ordenes del modulo). Ese es el valor que
la transportadora le cobro al cliente en el checkout.

Se usa en `create_invoice.go`, `create_journal.go` y `register_manual_invoice.go`.
`retry_invoice.go` no cambia: lee lo ya guardado en la factura.

### Intento descartado: deducir la comision de cod_total

La primera version calculaba `envio = cod_total - subtotal + discount`. Cuadra
en WooCommerce, pero **`cod_total` no significa lo mismo en Shopify**. Medido
contra produccion, negocio 34 (sin intermediarios, 350 facturas COD emitidas via
Softpymes):

- orden #114510: envio habria pasado de 11.000 a **76.911** (le sumaba el
  descuento de 76.911).
- orden #113372: el total habria pasado de 211.780 a **243.360** (le sumaba un
  IVA que ya venia dentro del subtotal).

Es decir, habria roto 350 facturas de otro cliente para arreglar 7. Por eso la
version final lee un campo explicito en vez de despejar una ecuacion.

### Alcance verificado en produccion

Ordenes COD desde 2026-06-01 con meta `cod_carrier_fee`:

| Plataforma | Ordenes COD | Con meta |
|---|---|---|
| shopify | 1.541 | **0** |
| woocommerce | 243 | 16 |
| manual | 115 | **0** |
| jumpseller | 2 | **0** |

Las 16 son de dos negocios: 46 Viga (8, 7 facturadas) y 26 Demo (8, 0
facturadas). **Sin intermediarios (34) y su facturacion en Softpymes no se ven
afectados: cero ordenes.** El usuario pidio explicitamente no tocar ese flujo.

No se puso candado por plataforma a proposito: la condicion correcta de fondo es
"solo si el canal registro comision COD". Un `platform == woocommerce` haria que
el dia que el checkout de Shopify cobre comision se vuelva a facturar de menos
en silencio.

### Lo que quedo sin arreglar: invoices.total_amount

`invoices.total_amount` sale de `orders.total_amount`, que desde `b3611857` es
solo productos. La FV-2-14 queda guardada en 51.116 mientras el PDF de Siigo
dice 70.147. Se intento recalcularlo como `subtotal + tax - discount + envio` y
**eso fue justo lo que rompia a Shopify**: cada plataforma llena
`subtotal`/`tax`/`discount`/`total_amount` con semantica distinta (Shopify trae
el IVA dentro del subtotal, Woo no). Arreglarlo bien necesita una auditoria
plataforma por plataforma. Se dejo intacto.

### Prueba

- Unitaria (`invoice_shipping_test.go`): los 7 valores reales de Viga
  (18.122 / 29.248 / 25.234 / 25.017 / 19.322 / 34.498 / 25.147), mas casos
  Shopify con IVA incluido y con descuento que deben quedar **sin cambio**.
- Simulacion local contra el mock de Siigo (`back/testing`, :9095), base local
  5434, business 26, integracion 198, replicando la orden real
  `6033f6f8-3991-491a-a1a9-1529b4360799` (#14673) con su producto y su
  `shipping_details` de produccion:
  - con el codigo de produccion se reprodujo la FEVD 14 identica: `Envio`
    19.031, total **64.031**.
  - con el fix: `Envio` **25.147**, total **70.147**.
  - orden Shopify COD con IVA: `Envio` 14.000, total 211.780, **sin cambio**.
- Para reproducirla en local hizo falta agregar a la config lo que produccion ya
  tiene: `final_customer_when_no_id` (la orden viene sin cedula) y `seller_id`.
- Al mock se le agrego el contrato real que le faltaba: descuento por linea y
  validacion `sum(payments) == total de items` (la regla que hizo fallar
  facturas reales en `8c67a376`).

## Pendiente

Las 7 facturas ya timbradas ante la DIAN no se editan: se corrigen con nota
credito o con una factura adicional por la diferencia ($38.702).
