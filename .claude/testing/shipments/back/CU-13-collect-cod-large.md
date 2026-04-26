# CU-13: Cobro contra entrega (happy path - monto alto + verificacion DB)

## Objetivo
Cobrar la orden COD de 1.500.000 (SHIPMENT_ID_2) y validar a fondo el registro en
`payments` (amount, payment_method_id heredado, status, reference con notas).

## Precondiciones
- CU-11 ejecutado: SHIPMENT_ID_2 con status=delivered.
- ORDER_ID_2 con `cod_total=1500000`, `is_paid=false`, `payment_method_id=1`.

## Endpoint
```
POST /api/v1/shipments/{SHIPMENT_ID_2}/collect-cod
Authorization: Bearer {DEMO_TOKEN}
Content-Type: application/json

{ "notes": "Comprobante #INV-2026-001 - test CU-13" }
```

## Verificaciones HTTP
- [ ] Status 200
- [ ] `data.is_paid == true`, `data.cod_total == 1500000`

## Verificaciones MCP postgres
```sql
SELECT is_paid, paid_at FROM orders WHERE id = '{ORDER_ID_2}';
SELECT id, amount, payment_method_id, status, payment_reference
FROM payments WHERE order_id = '{ORDER_ID_2}' ORDER BY id DESC LIMIT 1;
```
- [ ] `orders.is_paid = true`
- [ ] Payment con `amount=1500000`, `payment_method_id=1`, `status='completed'`
- [ ] `payment_reference` contiene "INV-2026-001"
