# 2026-09-14 - COD: el corte pagaba cod_total aunque la guia devolvia menos, y cambio de transportadora

**Ticket:** TKT-000088 | **Business:** 46 (Viga ropa deportiva) | **Transportadoras:** Interrapidisimo, Coordinadora, Envia (EnvioClick)

Continuacion de [2026-09-14-cod-total-mezclaba-comision-checkout-woocommerce.md](2026-09-14-cod-total-mezclaba-comision-checkout-woocommerce.md).

## Resumen

Al generar la guia, EnvioClick calibra el valor a recaudar y responde la comision
exacta para ese valor. El sistema usaba ese resultado para la guia pero **no lo
guardaba**: `orders.cod_total` se quedaba con el valor estimado al crear la orden
y el corte contra entrega paga con `cod_total`. En varias ordenes se le iba a
pagar al negocio mas de lo que devuelve EnvioClick. Las guias estaban bien; el
registro no.

## Regla vigente (decision del usuario)

- **Manda la guia.** `shipments.cod_collect_amount` es lo que imprime el PDF.
- Al negocio nunca se le paga mas de lo que devuelve la guia
  (`cod_collect_amount - cod_carrier_fee`). En ordenes sin corte, si `cod_total`
  es mayor se ajusta; si es menor (a favor de Probability) se deja.
- Guia generada con **otra transportadora** distinta a la cotizada en la tienda:
  cobra neto + comision real, la paga el cliente final, el negocio recibe su neto.
  El modal de guia avisa la diferencia antes de generar.
- Billetera (flete + seguro + margenes) no cambia con estos ajustes.
- Ordenes con corte confirmado no se tocan.

## Casos y numeros reales

| Orden | Causa | cod_total antes | Guia (PDF) | Comision | Devuelve EnvioClick | cod_total ahora |
|---|---|---|---|---|---|---|
| 15788 | Viga genero la guia con Coordinadora; la tienda cotizo Interrapidisimo (5.238) | 65.990 | 71.228 | 6.116 | 65.112 | 65.112 |
| 14685 | Plugin cotizo comision 11.501, real 13.528 (bug viejo de cotizacion) | 199.068 | 210.569 | 13.528 | 197.041 | 197.041 |
| 14689 | Igual que 14685 | 65.595 | 69.273 | 5.121 | 64.152 | 64.152 |
| 14670 | Igual que 14685 | 60.645 | 64.323 | 4.827 | 59.496 | 59.496 |
| 14687 | Igual que 14685 | 57.593 | 61.271 | 4.645 | 56.626 | 56.626 |
| VIG-0068 | Manual: cod_total con envio estimado; la guia uso el envio real | 159.300,90 | 169.627 | 11.092 | 158.535 | 158.535 |

Evidencia 14685: PDF "Valor a cobrar $ 210.569"; sondeos de EnvioClick
codValue 197.041 -> comision 12.723, 209.764 -> 13.480, 210.569 -> 13.528; guia
generada con 210.569.

Quedaron a favor de Probability y no se tocaron: VIG-0083 (+600), VIG-0080 (+559),
15770 y 15787 (+195), 15765 (+2) y 9 ordenes con +1 peso de redondeo.

**VIG-0122**: marcada contra entrega (`cod_total` 64.119) pero el metodo de pago es
transferencia y la guia salio sin recaudo (PDF "Valor a cobrar $0", sin `codValue`
en el log). Por decision del usuario se deja asi; hay que excluirla a mano del corte.

## Correccion

| Commit (main) | Que |
|---|---|
| `84a26e41` | El backend entrega el cobro COD calculado (`shared/cod.Summarize`); front y app solo muestran |
| `d0fdc3e4` | `shipments.cod_collect_amount`, guia con otra transportadora cobra neto + comision real, aviso en modal y nota en detalle |
| `5f4c190d` | Fix 500 en contra entrega: faltaba la columna en el `JOIN LATERAL` de codreport |
| `06d1876e` | El neto mostrado es `min(cod_total, cod_collect_amount - cod_carrier_fee)` |

Datos en produccion (ver `back/migration/MIGRACIONES.md`):

1. DDL `shipments.cod_collect_amount`.
2. Relleno desde `shipment_sync_logs` (ultimo `/api/v2/shipment` exitoso con `codValue`): 91 envios.
3. Relleno leyendo el PDF de las guias de Viga sin log: 63 envios.
   Etiquetas: Coordinadora "Valor a recaudar: en TARJETA o Efectivo N", Envia
   "Valor a recaudar ... $N", Interrapidisimo "Valor a cobrar ... $ N".
4. `cod_total` de 15788 y de las 5 ordenes de la tabla.

## Verificacion

- Produccion por API: 170 ordenes COD de Viga, detalle / editar / contra entrega
  coinciden, 0 inconsistencias. 0 ordenes sin corte pagando mas que su guia.
- Playwright en produccion: contra entrega (11 paginas, 3 meses), detalle y editar
  de 14685, VIG-0068 y 15788: neto + comision = cobro de la guia.
- Local contra produccion antes del deploy: 629 ordenes comparadas contra la regla.

## Hipotesis descartadas

1. **"La guia estaba mal"**: falso. Cada PDF coincide con el `codValue` enviado y
   EnvioClick informo la comision para ese valor.
2. **"EnvioClick suma la comision encima del valor declarado"**: falso. El PDF
   imprime exactamente el `codValue` que enviamos; la comision se descuenta.
3. **"La comision del carrier estaba mal registrada"**: falso. `cod_carrier_fee`
   coincide con la cotizacion de EnvioClick para el valor final.
4. **Reemplazar `cod_collect_amount` para cobrarle a Viga**: no, esa columna es la
   verdad del PDF. Lo que paga el corte es `cod_total`.

## Si algo falla en el futuro

```sql
-- Ordenes sin corte que pagarian mas de lo que devuelve la guia
SELECT o.order_number, o.cod_total, s.cod_collect_amount, s.cod_carrier_fee,
       s.cod_collect_amount - s.cod_carrier_fee AS devuelve
FROM orders o
JOIN LATERAL (SELECT * FROM shipments s WHERE s.order_id = o.id AND s.deleted_at IS NULL
              ORDER BY s.created_at DESC LIMIT 1) s ON true
WHERE o.business_id = 46 AND o.deleted_at IS NULL AND o.cod_total > 0
  AND s.cod_collect_amount IS NOT NULL AND s.cod_carrier_fee > 0
  AND s.status NOT IN ('cancelled','returned')
  AND NOT EXISTS (SELECT 1 FROM cod_payment_cut_order c WHERE c.order_id = o.id AND c.deleted_at IS NULL)
  AND o.cod_total > s.cod_collect_amount - s.cod_carrier_fee + 0.5;
```

- Guia nueva sin `cod_collect_amount`: revisar que EnvioClick devuelva `codValue`
  en la respuesta (`transport_request_consumer.go`) y que `response_consumer.go`
  lo guarde.
- Si una pantalla y la guia no cuadran, comparar contra el PDF del `guide_url`,
  no contra la formula.
- Cualquier `JOIN LATERAL` o SELECT que elija columnas de `shipments` a mano
  necesita `cod_collect_amount` (asi nacio el 500 de `5f4c190d`).

## Pendiente

- Aviso del modal de guia no visto en pantalla (requiere cotizar con EnvioClick).
- Modal de guias masivas sin aviso de cambio de transportadora.
- Otros negocios: guias viejas sin `cod_collect_amount` usan la formula.
- El corte sigue pagando `cod_total`; un caso nuevo como VIG-0068 (envio real
  distinto) requiere ajuste manual o cambiar el corte para pagar el neto de la guia.
- Confirmar contra la liquidacion real de EnvioClick.
- App movil requiere release para leer `cod_customer_charge`.
