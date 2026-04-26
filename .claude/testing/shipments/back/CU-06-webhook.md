# CU-06: Webhook EnvioClick

## Endpoint
`POST /api/v1/webhooks/envioclick`

## Precondiciones
- Generar una nueva guia para tener un shipment activo (repetir flujo CU-02 + CU-03 con nueva orden o crear shipment manual)
- Guardar TRACKING_NUMBER_2 y SHIPMENT_ID_2

## Caso 6.1: Webhook con evento in_transit (ultimo evento)

El fix toma `Events[len-1]` como el evento mas reciente (orden cronologico).

**Request:**
```
POST /api/v1/webhooks/envioclick
Content-Type: application/json

{
  "trackingCode": "{TRACKING_NUMBER_2}",
  "myShipmentReference": "Orden test",
  "idOrder": 12345,
  "realPickupDate": "2026-04-14 08:00:00",
  "realDeliveryDate": "",
  "events": [
    {
      "timestamp": "2026-04-14T08:00:00",
      "statusStep": "Pendiente de Recoleccion",
      "status": "pending",
      "statusDetail": "Esperando recoleccion",
      "incidence": false,
      "description": "Guia creada"
    },
    {
      "timestamp": "2026-04-14T10:00:00",
      "statusStep": "Envio Recolectado",
      "status": "picked_up",
      "statusDetail": "Paquete recogido",
      "incidence": false,
      "description": "Recolectado en origen"
    },
    {
      "timestamp": "2026-04-14T14:00:00",
      "statusStep": "En transito",
      "status": "in_transit",
      "statusDetail": "En ruta",
      "incidence": false,
      "description": "Paquete en camino"
    }
  ]
}
```

**Expected Response (200):**
```json
{
  "success": true,
  "message": "Estado actualizado correctamente",
  "tracking": "{TRACKING_NUMBER_2}",
  "new_status": "in_transit",
  "has_incidence": false
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] new_status = "in_transit" (tomo el ULTIMO evento "En transito", no el primero "Pendiente")
- [ ] has_incidence = false

**Verificacion MCP:**
```sql
SELECT id, status, shipped_at FROM shipments WHERE id = {SHIPMENT_ID_2};
```
- [ ] status = 'in_transit'
- [ ] shipped_at no es NULL (parseado de realPickupDate)

## Caso 6.2: Webhook con evento delivered

**Request:**
```
POST /api/v1/webhooks/envioclick
Content-Type: application/json

{
  "trackingCode": "{TRACKING_NUMBER_2}",
  "myShipmentReference": "Orden test",
  "idOrder": 12345,
  "realPickupDate": "2026-04-14 08:00:00",
  "realDeliveryDate": "2026-04-15 11:30:00",
  "events": [
    {
      "timestamp": "2026-04-14T08:00:00",
      "statusStep": "Pendiente de Recoleccion",
      "status": "pending",
      "statusDetail": "",
      "incidence": false,
      "description": "Guia creada"
    },
    {
      "timestamp": "2026-04-14T14:00:00",
      "statusStep": "En transito",
      "status": "in_transit",
      "statusDetail": "",
      "incidence": false,
      "description": "En camino"
    },
    {
      "timestamp": "2026-04-15T11:30:00",
      "statusStep": "Entregado",
      "status": "delivered",
      "statusDetail": "Recibido por portero",
      "incidence": false,
      "description": "Entregado exitosamente"
    }
  ]
}
```

**Expected Response (200):**
```json
{
  "success": true,
  "new_status": "delivered",
  "has_incidence": false
}
```

**Verificaciones:**
- [ ] new_status = "delivered" (ultimo evento)

**Verificacion MCP:**
```sql
SELECT id, status, delivered_at FROM shipments WHERE id = {SHIPMENT_ID_2};
```
- [ ] status = 'delivered'
- [ ] delivered_at = '2026-04-15 11:30:00' (parseado de realDeliveryDate)

## Caso 6.3: Webhook con incidencia

**Request:**
```
POST /api/v1/webhooks/envioclick
Content-Type: application/json

{
  "trackingCode": "{TRACKING_NUMBER_2}",
  "myShipmentReference": "Orden test",
  "idOrder": 12345,
  "events": [
    {
      "timestamp": "2026-04-14T08:00:00",
      "statusStep": "En transito",
      "status": "in_transit",
      "incidence": false,
      "description": "En camino"
    },
    {
      "timestamp": "2026-04-15T16:00:00",
      "statusStep": "Novedad",
      "status": "incident",
      "incidence": true,
      "incidenceType": "direccion_erronea",
      "description": "No se encontro la direccion"
    }
  ]
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] new_status = "failed" (incidence=true mapea a failed)
- [ ] has_incidence = true

## Caso 6.4: Webhook con tracking inexistente

**Request:**
```
POST /api/v1/webhooks/envioclick
Content-Type: application/json

{
  "trackingCode": "FAKE-NONEXISTENT-123",
  "events": [
    {
      "timestamp": "2026-04-14T08:00:00",
      "statusStep": "En transito",
      "incidence": false
    }
  ]
}
```

**Expected Response (200):**
```json
{
  "success": true,
  "message": "Envio no encontrado en el sistema, ignorado"
}
```

**Verificaciones:**
- [ ] Status code = 200 (no retorna error para evitar reintentos de EnvioClick)
- [ ] message indica "no encontrado" o "ignorado"

## Caso 6.5: Webhook sin eventos

**Request:**
```
POST /api/v1/webhooks/envioclick
Content-Type: application/json

{
  "trackingCode": "{TRACKING_NUMBER_2}",
  "events": []
}
```

**Expected Response (200):**
```json
{
  "success": true,
  "message": "Sin eventos en el payload"
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] No se actualiza el shipment
