# CU-16: Validacion - cobrar orden que no es COD

## Objetivo
La orden 5 fue creada sin `cod_total`. Intentar cobrarla por su shipment debe
retornar `ErrOrderNotCOD`.

## Precondiciones
- CU-09 caso N=5 ejecutado: ORDER_ID_5 con `cod_total=null`.
- Crear shipment para esa orden via POST /shipments con status delivered.

## Paso 16.1: Crear shipment manual delivered
```
POST /api/v1/shipments
Authorization: Bearer {DEMO_TOKEN}

{
  "order_id": "{ORDER_ID_5}",
  "client_name": "No COD",
  "destination_address": "Calle 99 #1-2",
  "tracking_number": "TEST-NO-COD-{TS}",
  "carrier": "TEST",
  "status": "delivered",
  "delivered_at": "2026-04-25T18:00:00Z"
}
```
- [ ] Status 200/201 -> guardar SHIPMENT_ID_5

## Paso 16.2: Intentar cobro
```
POST /api/v1/shipments/{SHIPMENT_ID_5}/collect-cod
{ "notes": "no debe permitir" }
```
- [ ] Status 400
- [ ] `message` contiene "not a cash on delivery" (o equivalente espanol)

## Verificaciones MCP postgres
```sql
SELECT count(*) FROM payments WHERE order_id = '{ORDER_ID_5}';
```
- [ ] `count = 0`
