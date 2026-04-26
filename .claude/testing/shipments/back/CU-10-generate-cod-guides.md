# CU-10: Generar guias EnvioClick contra entrega via mock

## Objetivo
Para cada una de las ordenes COD creadas en CU-09 (N=1..4), generar una guia EnvioClick
usando el mock (puerto 9091) y verificar que el shipment queda creado y vinculado.

## Precondiciones
- CU-09 ejecutado (ORDER_ID_1..4 disponibles).
- Mock EnvioClick corriendo en `http://localhost:9091/api/v2`.
- Integracion EnvioClick del business 26 con `is_testing=true` (la imagen del usuario confirma esto).

## Paso 10.1: Cotizar (POST /shipments/quote)

Por cada orden:

```json
{
  "order_uuid": "{ORDER_ID_N}",
  "packages": [{"weight": 1, "height": 10, "width": 15, "length": 20}],
  "description": "Test COD {N}",
  "contentValue": {ORDER_TOTAL},
  "codValue": {COD_TOTAL},
  "codPaymentMethod": "cash",
  "includeGuideCost": false,
  "insurance": false,
  "origin":      {"daneCode": "11001000", "address": "Calle 100 #15-20"},
  "destination": {"daneCode": "08001000", "address": "Carrera 54 #72-80"}
}
```

- [ ] Status 200/202
- [ ] Respuesta contiene `correlation_id`
- [ ] (Eventualmente via SSE/poll) `rates` con al menos una opcion `cod=true` -> guardar `idRate`, `carrier`, `flete` -> ID_RATE_N, CARRIER_N, FLETE_N

> Nota: si el flujo de quote es asincrono via SSE y el runner no tiene SSE,
> el runner puede saltar este paso y usar `idRate=1, carrier="COORDINADORA", flete=15000`
> (valores aceptados por el mock para test).

## Paso 10.2: Generar guia (POST /shipments/generate)

```json
{
  "idRate": {ID_RATE_N},
  "carrier": "{CARRIER_N}",
  "order_uuid": "{ORDER_ID_N}",
  "myShipmentReference": "COD {ORDER_NUMBER_N}",
  "external_order_id": "{ORDER_NUMBER_N}",
  "requestPickup": false,
  "pickupDate": "2026-04-26",
  "insurance": false,
  "description": "COD shipment",
  "contentValue": {ORDER_TOTAL},
  "codValue": {COD_TOTAL},
  "codPaymentMethod": "cash",
  "includeGuideCost": false,
  "totalCost": {FLETE_N},
  "packages": [{"weight":1,"height":10,"width":15,"length":20}],
  "origin":      {"daneCode":"11001000","address":"Calle 100 #15-20","company":"Demo","firstName":"Demo","lastName":"Test","email":"demo@probability.com","phone":"3001234567","suburb":"Usaquen","crossStreet":"Calle 100","reference":"Of 301"},
  "destination": {"daneCode":"08001000","address":"Carrera 54 #72-80","company":"COD Test {N}","firstName":"Cliente","lastName":"COD","email":"cod{N}@test.com","phone":"300100000{N}","suburb":"Norte","crossStreet":"Carrera 54","reference":"Apto 1"}
}
```

- [ ] Status 202
- [ ] `shipment_id` > 0  -> guardar SHIPMENT_ID_N
- [ ] (esperar ~5s) `GET /shipments/{SHIPMENT_ID_N}` retorna:
  - [ ] `tracking_number` no vacio
  - [ ] `carrier` no vacio
  - [ ] `is_test = true`
  - [ ] `order_id == ORDER_ID_N`

## Paso 10.3: Verificar que aparece en /shipments/cod

```
GET /api/v1/shipments/cod?page_size=50
```

- [ ] Cada SHIPMENT_ID_N esta en la lista
- [ ] `cod_total` correcto, `is_paid=false`
