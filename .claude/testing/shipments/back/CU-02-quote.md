# CU-02: Cotizar envio (Quote)

## Endpoint
`POST /api/v1/shipments/quote`

## Precondiciones
- CU-01 completado (tokens obtenidos, test mode activo)
- Mock server en localhost:9091

## Caso 2.1: Cotizacion exitosa (super admin con business_id)

**Request:**
```
POST /api/v1/shipments/quote
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "packages": [{"weight": 2, "height": 15, "width": 20, "length": 25}],
  "description": "Test shipment",
  "contentValue": 300000,
  "includeGuideCost": false,
  "codPaymentMethod": "",
  "insurance": false,
  "origin": {
    "daneCode": "11001000",
    "address": "Calle 100 #15-20"
  },
  "destination": {
    "daneCode": "05001000",
    "address": "Carrera 43A #1-50"
  }
}
```

**Expected Response (200 o 202):**
```json
{
  "success": true,
  "correlation_id": "{CORRELATION_ID_QUOTE}",
  "data": {
    "data": {
      "rates": [
        {
          "idRate": 1001,
          "carrier": "Servientrega",
          "product": "Mercancia Premier",
          "flete": 15000,
          "minimumInsurance": 300,
          "extraInsurance": 1500,
          "deliveryDays": 3,
          "cod": true
        }
      ]
    }
  }
}
```

**Verificaciones:**
- [ ] Status code = 200 o 202
- [ ] success = true
- [ ] correlation_id presente -> guardar como CORRELATION_ID_QUOTE
- [ ] Si 200: data.data.rates es array con 3-5 elementos
- [ ] Cada rate tiene: idRate, carrier, product, flete, deliveryDays, cod (boolean)
- [ ] flete > 0 para todos los rates
- [ ] Algunos rates tienen cod=true, otros cod=false

## Caso 2.2: Cotizacion con valor COD

**Request:**
```
POST /api/v1/shipments/quote
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "packages": [{"weight": 2, "height": 15, "width": 20, "length": 25}],
  "description": "COD shipment",
  "contentValue": 300000,
  "codValue": 300000,
  "codPaymentMethod": "cash",
  "includeGuideCost": false,
  "insurance": false,
  "origin": {
    "daneCode": "11001000",
    "address": "Calle 100 #15-20"
  },
  "destination": {
    "daneCode": "76001000",
    "address": "Avenida 6N #25-30"
  }
}
```

**Verificaciones:**
- [ ] Status code = 200 o 202
- [ ] success = true
- [ ] rates contiene al menos 1-2 opciones con cod=true
- [ ] El mock garantiza que las primeras 2 rates tienen cod=true cuando codValue > 0

## Caso 2.3: Error - DANE code invalido

**Request:**
```
POST /api/v1/shipments/quote
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "packages": [{"weight": 2, "height": 15, "width": 20, "length": 25}],
  "description": "Test",
  "contentValue": 50000,
  "includeGuideCost": false,
  "codPaymentMethod": "",
  "origin": {
    "daneCode": "INVALID",
    "address": "Test"
  },
  "destination": {
    "daneCode": "05001000",
    "address": "Test"
  }
}
```

**Verificaciones:**
- [ ] Respuesta con error (status error en SSE o response directa)
- [ ] Mensaje contiene "dane" o "invalido"

## Caso 2.4: Error - Sin paquetes

**Request:**
```
POST /api/v1/shipments/quote
Authorization: Bearer {SUPER_ADMIN_TOKEN}
Content-Type: application/json

{
  "order_uuid": "a9a03251-f01a-409f-b848-9de8598e84de",
  "description": "Test",
  "contentValue": 50000,
  "includeGuideCost": false,
  "codPaymentMethod": "",
  "origin": {"daneCode": "11001000", "address": "Test"},
  "destination": {"daneCode": "05001000", "address": "Test"}
}
```

**Verificaciones:**
- [ ] Status 400 (falta packages en validacion del handler)
- [ ] Error indica que packages es requerido
