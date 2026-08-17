# Guias contra entrega - Verificacion obligatoria

Checklist para validar cualquier cambio que toque tarifas, guias o recaudo COD.
Nacio del bug de precios del cotizador expres (2026-08-16), donde cada pantalla
sumaba los componentes de la tarifa con una formula distinta.

## 1. La formula unica del costo

```
guia  = flete + seguro minimo + (seguro extra si asegura)
                + (margen COD si hay contra entrega Y la tarifa la soporta)
total = guia + (comision del carrier si hay contra entrega)
```

Vive en `front/central/src/shared/utils/rate-pricing.ts`
(`rateGuideCost`, `rateCarrierFee`, `rateTotalCost`). **Nadie mas debe sumar
tarifas a mano.** Si una pantalla nueva muestra precios de envio, usa esas
funciones; si necesita otra combinacion, se agrega ahi, no en el componente.

Referencias que fijan la formula: `woocommerce-shipping-rates.go` (mapQuoteRatesToWoo)
y `ShippingForm.tsx` (handleGenerate).

Errores clasicos que evita:
- Omitir el seguro minimo. En Interrapidisimo llega a $6.530: el usuario ve un
  precio y se le cobra otro.
- Sumar el margen COD a un envio que no es contra entrega.
- Olvidar que el margen COD solo aplica si `rate.cod === true`.

## 2. Que significa cada monto

| Campo | Que es | Quien lo calcula |
|---|---|---|
| `orders.shipping_cost` | costo de la guia para el negocio | front al crear la orden; lo pisa el costo real al generar guia |
| `orders.cod_total` | **neto** que le queda al negocio | producto + envio, SIN comision |
| `codValue` declarado al carrier | lo que el mensajero cobra al cliente | `cod_total` + comision |
| `shipments.cod_carrier_fee` | comision real del carrier | se reescribe tras la calibracion COD |
| `shipments.cod_probability_margin` | margen de Probability sobre la comision | monto fijo por negocio |

`OrderCodBasis.AmountToCollect` (`shipments/internal/domain/ports.go`) es la
fuente de verdad: `NetTarget + comision`. Que `cod_total` no incluya la comision
**no es un bug**, es la convencion.

## 3. Prueba E2E obligatoria

Antes de dar por bueno un cambio en este terreno, correr el flujo completo y
comparar los cinco puntos. Deben cuadrar TODOS:

1. **Cotizador** - precio en tarjeta y total con comision.
2. **Modal de retomar** - misma tarifa, mismo numero. Si difieren, hay bug.
3. **Orden en BD** - `shipping_cost`, `cod_total`, `is_cod`.
4. **PDF de la guia** - el "Valor a recaudar" impreso por la transportadora debe
   ser `producto + envio + comision`. Descargar el `guide_url` y leerlo, no
   confiar en la respuesta del API.
5. **Modulo contra entrega** (`/shipments/cod` o `GET /cod-report/orders`) -
   `cod_total`, `net`, `cod_carrier_fee`, `has_guide`, `guide_number`.

Ademas: `carrier_cost + applied_margin == total_cost`.

Probar con **Interrapidisimo** ademas de Coordinadora: es la que tiene el seguro
minimo mas alto y variable, donde los errores de formula se ven mas grandes.

## 4. Payload de generacion de guia

EnvioClick valida y rechaza con 422. Limites que ya nos mordieron:

- `external_order_id` y `myShipmentReference`: **maximo 28 caracteres**. Mandar
  el `order_number`, NUNCA el UUID de la orden (36 caracteres).
- `origin.email` valido y `origin.phone` de 7 a 10 digitos: salen de la bodega.
  Una bodega sin contacto rompe la guia aunque la cotizacion haya funcionado.
- `suburb` de origen y destino: entre 2 y 30 caracteres, no puede ir vacio.

El diagnostico esta siempre en `shipment_sync_logs` (`request_payload`,
`response_body`, `error_message`), no en el toast del front.

## 5. Cuidado al probar en local

`back/central/.env` apunta a **RDS de produccion**. Una prueba local crea
cotizaciones, ordenes y guias reales, y la guia le programa recoleccion a la
transportadora. Cancelar siempre al terminar y avisar al usuario antes.
