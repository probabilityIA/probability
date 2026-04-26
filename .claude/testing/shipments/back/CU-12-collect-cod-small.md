# CU-12: Cobro contra entrega (happy path - monto pequeno)

## Objetivo
Marcar como recaudada la orden COD de 50.000 (SHIPMENT_ID_1) y verificar que se
crea el registro en `payments` y la orden queda `is_paid=true`.

## Precondiciones
- CU-11 ejecutado: SHIPMENT_ID_1 con status=delivered.
- ORDER_ID_1 con `cod_total=50000`, `is_paid=false`.

## Endpoint
```
POST /api/v1/shipments/{SHIPMENT_ID_1}/collect-cod
Authorization: Bearer {DEMO_TOKEN}
Content-Type: application/json

{ "notes": "Recibido en efectivo - test CU-12" }
```

## Verificaciones HTTP
- [ ] Status 200, `success=true`
- [ ] `message` contiene "exitosamente"
- [ ] `data.is_paid == true`
- [ ] `data.paid_at` poblado

## Verificaciones MCP postgres
```sql
SELECT id, is_paid, paid_at FROM orders WHERE id = '{ORDER_ID_1}';
SELECT id, order_id, amount, status, paid_at FROM payments WHERE order_id = '{ORDER_ID_1}' ORDER BY id DESC LIMIT 1;
```
- [ ] `orders.is_paid = true` y `paid_at` no nulo
- [ ] Existe payment con `amount=50000`, `status='completed'`, `paid_at` no nulo
