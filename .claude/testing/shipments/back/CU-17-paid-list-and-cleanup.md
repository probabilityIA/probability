# CU-17: Validar lista filtrada de cobradas + verificacion final

## Objetivo
Despues de los cobros de CU-12, CU-13 y CU-15, validar que el listado COD refleja
correctamente los estados.

## Precondiciones
- CU-12, CU-13, CU-15 ejecutados (3 ordenes pagadas).
- CU-14 ejecutado (1 orden COD aun por cobrar - pending).
- CU-16 (orden no COD: NO debe aparecer en este listado).

## Paso 17.1: GET solo recaudados
```
GET /api/v1/shipments/cod?is_paid=true&page_size=50
Authorization: Bearer {DEMO_TOKEN}
```
- [ ] Incluye SHIPMENT_ID_1, SHIPMENT_ID_2, SHIPMENT_ID_4
- [ ] No incluye SHIPMENT_ID_3 ni SHIPMENT_ID_5
- [ ] Cada item tiene `is_paid=true` y `paid_at` no nulo

## Paso 17.2: GET solo por cobrar entregados
```
GET /api/v1/shipments/cod?is_paid=false&status=delivered&page_size=50
```
- [ ] No incluye SHIPMENT_ID_1, _2, _4 (ya pagados)
- [ ] Si SHIPMENT_ID_3 estuviera delivered apareceria; como esta pending, no debe estar aqui

## Paso 17.3: Total general
```
GET /api/v1/shipments/cod?page_size=1
```
- [ ] `total >= 4` (al menos los SHIPMENT_ID_1..4 creados)

## Paso 17.4: Suma de cobros (verificacion contable)
```sql
SELECT COALESCE(SUM(amount),0) AS total_cobrado
FROM payments
WHERE order_id IN ('{ORDER_ID_1}','{ORDER_ID_2}','{ORDER_ID_4}')
  AND status='completed';
```
- [ ] `total_cobrado = 50000 + 1500000 + 320000 = 1870000`
