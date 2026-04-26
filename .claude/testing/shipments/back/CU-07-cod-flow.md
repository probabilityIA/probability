# CU-07: Flujo contra entrega (COD) completo

## Precondiciones
- CU-01 completado
- Crear una orden COD en la base de datos para el negocio Demo

## Caso 7.1: Crear orden COD de prueba

**Accion MCP:** Crear orden con cod_total para testing
```sql
UPDATE orders SET cod_total = 150000
WHERE id = 'a9a03251-f01a-409f-b848-9de8598e84de' AND business_id = 26;
```

**Verificacion MCP:**
```sql
SELECT id, order_number, cod_total, total_amount FROM orders
WHERE id = 'a9a03251-f01a-409f-b848-9de8598e84de';
```
- [ ] cod_total = 150000

## Caso 7.2: Cotizar envio COD

**Request:**
```
POST /api/v1/shipments/quote
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "packages": [{"weight": 1.5, "height": 12, "width": 18, "length": 22}],
  "description": "Orden COD test",
  "contentValue": 300000,
  "codValue": 150000,
  "codPaymentMethod": "cash",
  "includeGuideCost": false,
  "insurance": true,
  "origin": {
    "daneCode": "11001000",
    "address": "Calle 100 #15-20"
  },
  "destination": {
    "daneCode": "08001000",
    "address": "Carrera 54 #72-80"
  }
}
```

**Verificaciones:**
- [ ] Status 200 o 202
- [ ] rates contiene opciones con cod=true
- [ ] Al menos las primeras 2 rates tienen cod=true (garantizado por el mock cuando codValue > 0)
- [ ] Guardar un idRate donde cod=true -> ID_RATE_COD

## Caso 7.3: Generar guia COD

**Request:**
```
POST /api/v1/shipments/generate
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "idRate": {ID_RATE_COD},
  "carrier": "{CARRIER_COD}",
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "myShipmentReference": "Orden COD prob-0016",
  "external_order_id": "prob-0016",
  "requestPickup": false,
  "pickupDate": "2026-04-15",
  "insurance": true,
  "description": "COD shipment",
  "contentValue": 300000,
  "codValue": 150000,
  "codPaymentMethod": "cash",
  "includeGuideCost": false,
  "totalCost": {FLETE_COD},
  "packages": [{"weight": 1.5, "height": 12, "width": 18, "length": 22}],
  "origin": {
    "daneCode": "11001000",
    "address": "Calle 100 #15-20",
    "company": "ProbabilityIA",
    "firstName": "Admin",
    "lastName": "Test",
    "email": "test@probability.com",
    "phone": "3001234567",
    "suburb": "Usaquen",
    "crossStreet": "Calle 100 #15-20",
    "reference": "Oficina 301"
  },
  "destination": {
    "daneCode": "08001000",
    "address": "Carrera 54 #72-80",
    "company": "Cliente COD",
    "firstName": "Juan",
    "lastName": "Perez",
    "email": "juan@test.com",
    "phone": "3109876543",
    "suburb": "Norte Centro Historico",
    "crossStreet": "Carrera 54 #72-80",
    "reference": "Apto 502"
  }
}
```

**Verificaciones:**
- [ ] Status 202
- [ ] shipment_id > 0 -> guardar como SHIPMENT_ID_COD

**Verificacion MCP (esperar ~5s):**
```sql
SELECT id, tracking_number, carrier, status, is_test FROM shipments
WHERE id = {SHIPMENT_ID_COD};
```
- [ ] tracking_number asignado
- [ ] status = 'pending'
- [ ] is_test = true

## Caso 7.4: Cancelar guia COD y verificar sincronizacion

**Request:**
```
POST /api/v1/shipments/{TRACKING_NUMBER_COD}/cancel
Authorization: Bearer {SUPER_ADMIN_TOKEN}
```

**Verificaciones:**
- [ ] Status 202

**Verificacion MCP (esperar ~10s):**
```sql
SELECT id, status FROM shipments WHERE id = {SHIPMENT_ID_COD};
```
- [ ] status = 'cancelled'

```sql
SELECT id, order_number, status FROM orders
WHERE id = 'a9a03251-f01a-409f-b848-9de8598e84de';
```
- [ ] status = 'cancelled' (sincronizado por fix en response_consumer)

## Caso 7.5: Limpiar - restaurar orden de prueba

**Accion MCP:**
```sql
UPDATE orders SET cod_total = NULL, status = 'open',
  tracking_number = NULL, carrier = NULL, guide_link = NULL
WHERE id = 'a9a03251-f01a-409f-b848-9de8598e84de';
```

**Verificacion:**
- [ ] Orden restaurada a estado original
