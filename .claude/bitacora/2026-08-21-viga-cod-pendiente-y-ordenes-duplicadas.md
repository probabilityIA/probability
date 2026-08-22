# Viga: 3 envios COD sin pagar y ordenes duplicadas

Fecha: 2026-08-21
Business: 46 (Viga ropa deportiva). Reclamo por correo del 19/08/2026.

## Reclamo

Tres envios de mas de un mes marcados PENDIENTE por el negocio:

| Fecha | Cliente | Guia | Valor | Orden |
|---|---|---|---|---|
| 09/07 | Adriana del Pilar Rojas | 84151675104 (Coordinadora) | 77.000 | VIG-0014 |
| 14/07 | Maria Katherine Zabala | 84151679648 (Coordinadora) | 73.050,40 | VIG-0022 |
| 14/07 | Octavio Calderon | 240056906880 (Interrapidisimo) | 106.163,60 | VIG-0025 |

Total: 256.214,00

## Verificacion

Las tres estan `delivered` con `delivered_at` 14/07, 20/07 y 21/07 y cod_total
igual al valor reclamado. NINGUNA tiene fila en `cod_payment_cut_order`, o sea
que en `/shipments/cod` salen como PENDIENTE. El reclamo es correcto: la
plataforma coincide con el negocio, no estan pagadas.

`paidExpr` en `codreport/internal/infra/secondary/repository/orders_repo.go`
define pagada = existe `cod_payment_cut_order` con corte `confirmed`. No hay
proceso automatico: `SelectableCutOrders` lista las entregadas del rango que no
esten ya vinculadas y el admin elige cuales meter al corte. Estas quedaron
disponibles y nunca se seleccionaron. Cortes del negocio: 40, 41, 42, 43, 44,
46, 50, 51 confirmados y 52 en `draft`; sus rangos cubren julio, asi que la
omision fue humana al armar el corte, no un filtro del query.

## Duplicados

Patron: guia generada -> cancelada o fallida -> en vez de reintentar sobre la
misma orden, se crea una orden NUEVA con el mismo cliente y monto.

- Adriana: VIG-0013 (09/07 14:32, status 12 "Seleccionando productos", sin guia)
  y VIG-0014 (09/07 16:24, guia real, entregada).
- Octavio: VIG-0024 (14/07 22:19, guia 240056906035 CANCELADA) y VIG-0025
  (14/07 22:28, guia 240056906880 entregada).
- Ervin (mismo patron, no reclamado): VIG-0011 con 3 shipments
  (84151671323 cancelada + 2 `failed` sin tracking) y VIG-0012 con la guia
  buena. Esto es lo que se ve como "triplicado": una orden con varios shipments
  aparece varias veces en el listado de envios aunque solo una guia sea real.

Solo la duplicada de Octavio tiene guia real duplicada.

## Plata: doble cobro de wallet en Octavio

`transaction` USAGE/GUIDE del business 46:

- 47481d95-dfa6-49aa-b8c1-6b298ede1d99 - shipment 36296 (VIG-0024,
  240056906035, CANCELADA) - 16.123,90 - 2026-07-14 22:22
- 818e5dde-72ce-4e44-a4fd-98c7a6e32397 - shipment 36297 (VIG-0025,
  240056906880, entregada) - 16.163,60 - 2026-07-14 22:30

La guia cancelada NO se reembolso. Es el item pendiente de
`.claude/alerts/guias-duplicadas-doble-cobro.md`: cancelar no acredita el wallet.
Falta reintegrar 16.123,90.

## Acciones

1. Incluir VIG-0014, VIG-0022 y VIG-0025 en el proximo corte COD y pagar
   256.214,00.
2. Reintegrar 16.123,90 al wallet de Viga por la guia cancelada de VIG-0024.
3. Depurar VIG-0013 y VIG-0024 (cancelarlas / marcarlas duplicadas) para que no
   sigan inflando el listado de ordenes.
