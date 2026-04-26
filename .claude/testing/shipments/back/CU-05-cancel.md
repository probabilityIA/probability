# CU-05: Cancelar envio (Cancel)

## Endpoint
`POST /api/v1/shipments/{id}/cancel`

## Precondiciones
- CU-03 completado (guia generada con SHIPMENT_ID_1, TRACKING_NUMBER_1, ENVIOCLICK_ID_ORDER)

## Caso 5.1: Cancelar envio exitosamente

El cancel flow es asincrono:
1. Handler envia request a RabbitMQ
2. Consumer de envioclick verifica status via /track (debe ser "Pendiente...")
3. Consumer cancela via batch API (/v2cancellation/batch/order)
4. Response consumer actualiza shipment y orden a "cancelled"

**Request:**
```
POST /api/v1/shipments/{TRACKING_NUMBER_1}/cancel
Authorization: Bearer {SUPER_ADMIN_TOKEN}
```

**Expected Response (202):**
```json
{
  "success": true,
  "message": "Solicitud de cancelacion enviada. Sera procesada en breve.",
  "correlation_id": "{CORRELATION_ID}"
}
```

**Verificaciones:**
- [ ] Status code = 202
- [ ] success = true

**Verificacion MCP (esperar ~5-10s para procesamiento async):**
```sql
SELECT id, status, tracking_number FROM shipments WHERE id = {SHIPMENT_ID_1};
```
- [ ] status = 'cancelled'

**Verificacion MCP - orden sincronizada:**
```sql
SELECT id, order_number, status FROM orders
WHERE id = 'a9a03251-f01a-409f-b848-9de8598e84de';
```
- [ ] status = 'cancelled' (fix implementado en response_consumer.go)

## Caso 5.2: Verificar que el cancel usa batch API

El cancel verifica status primero con /track. El mock devuelve "Pendiente de Recoleccion" para shipments recien creados. Luego usa `/v2cancellation/batch/order` con el idOrder del metadata.

**Verificacion en logs del backend:**
- [ ] Log contiene "Proceeding to cancel via Batch API v2"
- [ ] Log NO contiene "falling back to singular DELETE"

## Caso 5.3: Error - Cancelar envio ya cancelado

Para este test, necesitamos generar una nueva guia primero (repetir CU-03.1 con otra orden o reutilizar).

**Request:**
```
POST /api/v1/shipments/{TRACKING_NUMBER_1}/cancel
Authorization: Bearer {SUPER_ADMIN_TOKEN}
```

**Verificaciones:**
- [ ] La solicitud se envia (202) pero el async falla
- [ ] SSE event shipment.cancel_failed se emite
- [ ] Mensaje de error indica que el envio no esta en estado cancelable

## Caso 5.4: Verificar status check case-insensitive

El fix en operations.go usa `strings.Contains(strings.ToLower(status), "pendiente")`.
Esto permite variaciones como:
- "Pendiente de Recoleccion"
- "PENDIENTE DE RECOLECCION"
- "pendiente"
- "Pendiente"

**Verificacion:** El caso 5.1 ya valida esto indirectamente, ya que el mock devuelve "Pendiente de Recoleccion" (sin tilde) y el cancel debe funcionar.
