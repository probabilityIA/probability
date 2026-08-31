# WooCommerce COD: la comision del carrier se declaraba dos veces

Resumen: cuando `cod_includes_shipping=true` (Woo cobra el envio en el
checkout), `AmountToCollect` volvia a sumarle la comision del carrier a un
`cod_total` que ya la traia incluida. El cliente final quedaba declarado a
pagar de mas al recibir la guia.

## Sintoma

Reclamo de Viga (business 46) por la orden 14668: la plataforma decia que se
recaudaria 115.017, pero el PDF de la guia (Coordinadora, tracking
84151729296) imprimia "Valor a recaudar: 121.133". Diferencia: 6.116, exacto
al `cod_carrier_fee` de la guia.

## Diagnostico

`orders.cod_total` para ordenes de WooCommerce con `cod_includes_shipping=true`
ya viene con la comision del carrier sumada, porque `mapQuoteRatesToWoo`
(`shipments/internal/infra/primary/handlers/woocommerce-shipping-rates.go:379`)
le suma `codCarrierFee + codProbabilityMargin` al costo que se muestra en el
checkout de la tienda: el cliente ya acepto pagar ese total sabiendo que ahi va
la comision.

Verificado con datos reales del pedido: `subtotal` 90.000 + `shipping_cost`
18.901 (flete+seguro+margen Probability) + `cod_carrier_fee` 6.116 =
`cod_total` 115.017. Coincide al peso.

Al generar la guia, `OrderCodBasis.NetTarget` (`shipments/internal/domain/ports.go`)
veia `CodIncludesShipping=true` y devolvia `CodTotal` tal cual, tratandolo como
si fuera el neto sin comision. `AmountToCollect`/`ResolveCODValue` le sumaban
la comision otra vez para calcular que declarar a EnvioClick:
115.017 + 6.116 = 121.133, el valor que termino impreso en el PDF y que
EnvioClick le cobraria al comprador.

## Causa raiz

`NetTarget`/`AmountToCollect` (`ports.go`) asumian siempre que `CodTotal` es
neto (sin comision) y le sumaban la comision del carrier. Correcto para
ordenes manuales (`CodIncludesShipping=false`, ver `orders/internal/app/usecasecreateorder`
y el mapper de WooCommerce cuando `cod_includes_shipping=false`: ahi si arma
`cod_total` sin comision, "que se agrega despues al generarla"). Incorrecto
para WooCommerce con `cod_includes_shipping=true`: ahi `cod_total` ya es bruto.

El test `order_cod_basis_test.go` tenia un caso llamado literalmente "total ya
incluye flete: solo suma comision" que codificaba la suposicion equivocada
(que `cod_includes_shipping` solo cubre el flete, no la comision).

## Correccion

`shipments/internal/domain/ports.go`: `NetTarget` recibe ahora tambien
`embeddedCarrierFee` y, cuando `CodIncludesShipping=true`, resta la comision
ya incluida en `CodTotal` antes de que `AmountToCollect` la vuelva a sumar.
Con la misma comision en ambos extremos (la que ya trae `cod_total` y la que
declara la guia), el resultado neto es `CodTotal` sin cambios; si la comision
real difiere de la cotizada, `AmountToCollect` ajusta la diferencia en vez de
duplicarla.

Caller actualizado: `generate-guide.go` `overrideCodValue` pasa `carrierFee` a
`NetTarget`. Test `order_cod_basis_test.go` corregido: el caso 3 ahora espera
176.494 (antes 182.610, que era el bug).

## Verificacion

`go build ./...` limpio. `go test ./services/modules/shipments/internal/domain/...`
en verde, incluido el caso corregido. Desplegado (push a `main`, commit
`ea633276`). El usuario probo en el Woo de pruebas (orden nueva + orden ya
sincronizada) y confirmo que el valor a recaudar de la guia ya no duplica la
comision.

## Auditoria de impacto

Solo 3 negocios tienen WooCommerce con `cod_includes_shipping=true`: 26 (Demo,
sin impacto real), 46 (Viga) y 52. Filtrando `orders.cod_includes_shipping=true`
(columna copiada en la orden, la fuente real que usa `GetOrderCodBasis`, no el
config de la integracion) y `integration_type='woocommerce'`:

- Business 52: 0 ordenes afectadas.
- Business 46: **solo 2 ordenes**, 14666 y 14668, ambas generadas el
  2026-08-29 y ambas siguen `pending` (guia creada, sin recoleccion/entrega
  todavia). `cod_carrier_fee` 6.116 cada una -> sobrecobro potencial 12.232 en
  total, **pero no se ha cobrado nada aun** porque no se han entregado.

El resto de guias COD de Viga con comision (85 registros) son de
`integration_type='platform'` (ordenes manuales), cuya integracion Plataforma
(id 81) tiene `cod_includes_shipping=false`: ese flujo nunca duplico la
comision, no esta afectado.

## Pendientes

- [ ] Deploy del fix.
- [ ] Decidir que hacer con las guias 84151729283 (orden 14666) y 84151729296
      (orden 14668): estan generadas con el valor declarado incorrecto
      (211.864 para 14666 / 121.133 para 14668, con la comision duplicada).
      Como siguen `pending`,
      se pueden cancelar y regenerar con el fix ya desplegado antes de que el
      mensajero recaude, evitando el sobrecobro real al cliente final.
- [ ] Confirmar si `cod_includes_shipping=true` lleva poco tiempo activo para
      Viga en WooCommerce (explicaria por que solo hay 2 ordenes con este
      patron pese a que la tienda tiene meses de historial) o si el flujo Woo
      simplemente se uso poco hasta ahora.
