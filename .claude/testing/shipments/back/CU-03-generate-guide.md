# CU-03: Generar guia (Generate)

## Endpoint
`POST /api/v1/shipments/generate`

## Precondiciones
- CU-02 completado (al menos una cotizacion exitosa)
- Guardar un idRate de CU-02.1

## Caso 3.1: Generar guia exitosamente

**Request:**
```
POST /api/v1/shipments/generate
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "idRate": {ID_RATE_FROM_QUOTE},
  "carrier": "{CARRIER_FROM_QUOTE}",
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "myShipmentReference": "Orden prob-0016",
  "external_order_id": "prob-0016",
  "requestPickup": false,
  "pickupDate": "2026-04-15",
  "insurance": true,
  "description": "Test shipment",
  "contentValue": 300000,
  "includeGuideCost": false,
  "codPaymentMethod": "",
  "totalCost": {FLETE_FROM_QUOTE},
  "packages": [{"weight": 2, "height": 15, "width": 20, "length": 25}],
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
    "daneCode": "05001000",
    "address": "Carrera 43A #1-50",
    "company": "Cliente",
    "firstName": "Carlos",
    "lastName": "Arrieta",
    "email": "carlos@test.com",
    "phone": "3009876543",
    "suburb": "El Poblado",
    "crossStreet": "Carrera 43A #1-50",
    "reference": "Casa"
  }
}
```

**Expected Response (202):**
```json
{
  "success": true,
  "message": "Solicitud de generacion de guia enviada. Sera procesada en breve.",
  "correlation_id": "{CORRELATION_ID}",
  "shipment_id": 123
}
```

**Verificaciones:**
- [ ] Status code = 202
- [ ] success = true
- [ ] shipment_id > 0 -> guardar como SHIPMENT_ID_1
- [ ] correlation_id presente

**Verificacion MCP (esperar ~5s para que el async procese):**
```sql
SELECT id, order_id, tracking_number, carrier, guide_url, status, is_test, metadata
FROM shipments WHERE id = {SHIPMENT_ID_1};
```
- [ ] tracking_number no es NULL (asignado por mock)
- [ ] guide_url no es NULL (URL del PDF)
- [ ] carrier no es NULL
- [ ] status = 'pending'
- [ ] is_test = true
- [ ] metadata contiene 'envioclick_id_order'

**Verificacion MCP - orden actualizada:**
```sql
SELECT tracking_number, carrier, guide_link FROM orders
WHERE id = 'a9a03251-f01a-409f-b848-9de8598e84de';
```
- [ ] tracking_number no es NULL (sincronizado desde shipment)
- [ ] guide_link no es NULL

-> Guardar TRACKING_NUMBER_1 y ENVIOCLICK_ID_ORDER del metadata

## Caso 3.2: Error - Sin origin

**Request:**
```
POST /api/v1/shipments/generate
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "idRate": 1001,
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "description": "Test",
  "contentValue": 50000,
  "packages": [{"weight": 1, "height": 10, "width": 10, "length": 10}],
  "destination": {"daneCode": "05001000", "address": "Test"}
}
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] Error indica "origin es requerido"

## Caso 3.3: Error - Sin destination

**Request:**
```
POST /api/v1/shipments/generate
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "idRate": 1001,
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "description": "Test",
  "contentValue": 50000,
  "packages": [{"weight": 1, "height": 10, "width": 10, "length": 10}],
  "origin": {"daneCode": "11001000", "address": "Test"}
}
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] Error indica "destination es requerido"
