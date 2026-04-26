# CU-15: Validacion - doble cobro

## Objetivo
Cobrar dos veces el mismo shipment (SHIPMENT_ID_4, cod=320.000): la segunda llamada
debe retornar `ErrOrderAlreadyPaid`.

## Precondiciones
- CU-11 ejecutado: SHIPMENT_ID_4 con status=delivered.
- ORDER_ID_4 con `cod_total=320000`, `is_paid=false`.

## Paso 15.1: Primer cobro (debe ser exitoso)
```
POST /api/v1/shipments/{SHIPMENT_ID_4}/collect-cod
Authorization: Bearer {DEMO_TOKEN}

{ "notes": "primer cobro CU-15" }
```
- [ ] Status 200, `is_paid=true`

## Paso 15.2: Segundo cobro (debe fallar)
```
POST /api/v1/shipments/{SHIPMENT_ID_4}/collect-cod
{ "notes": "intento de doble cobro" }
```
- [ ] Status 400
- [ ] `message` contiene "already paid"

## Verificaciones MCP postgres
```sql
SELECT count(*) FROM payments WHERE order_id = '{ORDER_ID_4}';
```
- [ ] `count = 1` (no se creo segundo payment)
