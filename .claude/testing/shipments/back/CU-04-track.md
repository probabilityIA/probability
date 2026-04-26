# CU-04: Rastrear envio (Track)

## Endpoint
`POST /api/v1/shipments/tracking/{tracking_number}/track`

## Precondiciones
- CU-03 completado (guia generada con TRACKING_NUMBER_1)

## Caso 4.1: Tracking exitoso

**Request:**
```
POST /api/v1/shipments/tracking/{TRACKING_NUMBER_1}/track
Authorization: Bearer {SUPER_ADMIN_TOKEN}
```

**Expected Response (202):**
```json
{
  "success": true,
  "message": "Solicitud de tracking enviada.",
  "correlation_id": "{CORRELATION_ID}"
}
```

**Verificaciones:**
- [ ] Status code = 202
- [ ] success = true
- [ ] correlation_id presente

**Verificacion MCP (esperar ~5s):**
```sql
SELECT id, status, metadata FROM shipments WHERE id = {SHIPMENT_ID_1};
```
- [ ] metadata contiene 'tracking_events' con al menos 2 eventos
- [ ] status puede ser 'Pendiente de Recoleccion' o mapeado a 'pending'/'in_transit'

## Caso 4.2: Tracking publico por tracking number

**Request:**
```
GET /api/v1/tracking/search?tracking_number={TRACKING_NUMBER_1}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] Retorna datos del shipment (tracking_number, carrier, status)

## Caso 4.3: Error - Tracking number inexistente

**Request:**
```
POST /api/v1/shipments/tracking/FAKE-TRACKING-999/track
Authorization: Bearer {SUPER_ADMIN_TOKEN}
```

**Verificaciones:**
- [ ] Respuesta con error (shipment no encontrado o error en tracking)
