# CU-09: Crear ordenes COD para el negocio Demo

## Objetivo
Crear 4 ordenes COD via API (que la UI invoca al usar el modal manual de orden) para
poder ejercitar el resto de los casos: listado, generacion de guia, cobro y validaciones.

## Precondiciones
- DEMO_TOKEN obtenido (CU-08).
- Integracion `manual` (id=35) activa en business 26.

## Endpoint
```
POST /api/v1/orders
Authorization: Bearer {DEMO_TOKEN}
Content-Type: application/json
```

## Body de cada orden

```json
{
  "business_id": 26,
  "integration_id": 35,
  "integration_type": "platform",
  "platform": "manual",
  "external_id": "cod-test-{TIMESTAMP_NS}-{N}",
  "subtotal": {AMOUNT},
  "total_amount": {AMOUNT},
  "cod_total": {AMOUNT},
  "currency": "COP",
  "customer_name": "COD Test {N}",
  "customer_email": "codtest{N}@test.com",
  "customer_phone": "300100000{N}",
  "customer_dni": "100000000{N}",
  "status": "pending",
  "payment_method_id": 1,
  "occurred_at": "2026-04-25T12:00:00Z",
  "imported_at": "2026-04-25T12:00:00Z",
  "order_items": [{
    "product_sku": "SKU-COD-{N}",
    "product_name": "Producto COD {N}",
    "quantity": 1,
    "unit_price": {AMOUNT},
    "total_price": {AMOUNT},
    "currency": "COP"
  }]
}
```

## Variantes a crear

| N | AMOUNT (COD) | Uso |
|---|---|---|
| 1 | 50000   | Happy path cobro pequeno (CU-12) |
| 2 | 1500000 | Happy path cobro alto (CU-13) |
| 3 | 250000  | Validacion no-delivered (CU-14) |
| 4 | 320000  | Validacion doble cobro (CU-15) |

Ademas se crea una orden NO COD (cod_total ausente) para CU-16:

| 5 | total=99000, **sin** cod_total | Validacion ErrOrderNotCOD |

## Verificaciones por orden

- [ ] Status 200, `success=true`
- [ ] `data.ID` (UUID) generado -> guardar como ORDER_ID_{N}
- [ ] `data.OrderNumber` no vacio (ej: prob-XXXX)
- [ ] `data.CodTotal` == AMOUNT (o null para N=5)
- [ ] `data.BusinessID` == 26
