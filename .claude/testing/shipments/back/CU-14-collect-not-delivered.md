# CU-14: Validacion - cobrar shipment no entregado

## Objetivo
Intentar cobrar SHIPMENT_ID_3 (cod=250.000) que sigue en `pending` debe rechazarse
con `ErrShipmentNotDelivered`.

## Precondiciones
- CU-10 ejecutado, SHIPMENT_ID_3 status=`pending` (NO ejecutado en CU-11).
- ORDER_ID_3 con `cod_total=250000`, `is_paid=false`.

## Endpoint
```
POST /api/v1/shipments/{SHIPMENT_ID_3}/collect-cod
Authorization: Bearer {DEMO_TOKEN}
Content-Type: application/json

{ "notes": "no debe permitir" }
```

## Verificaciones
- [ ] Status 400
- [ ] `success=false`
- [ ] `message` contiene "delivered" (mensaje del error `ErrShipmentNotDelivered`)

## Verificacion MCP postgres
```sql
SELECT is_paid FROM orders WHERE id = '{ORDER_ID_3}';
SELECT count(*) FROM payments WHERE order_id = '{ORDER_ID_3}';
```
- [ ] `orders.is_paid = false` (sin cambios)
- [ ] `count = 0` (no se creo payment)
