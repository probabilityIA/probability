# Plan: Sincronizacion manual de estados Envioclick -> Probability

## Objetivo

Endpoint manual para que un admin (desde el frontend) dispare una resincronizacion
de estados de envios con el API de Envioclick. Reutiliza el flujo async ya construido
(`webhook_update`) — una vez disparado, cada shipment se actualiza por la cola y el
frontend ve los cambios via SSE.

## Decisiones confirmadas por el usuario

1. **Trigger:** endpoint manual `POST /api/v1/shipments/sync-status` (activable desde el modulo shipments del frontend)
2. **Scope:** solo Envioclick por ahora (disenar reutilizable para Enviame/MiPaquete cuando tengan API)
3. **Filtro:** status IN (`pending`, `in_transit`, `picked_up`, `out_for_delivery`, `on_hold`) + rango de fechas (`date_from`, `date_to`) + `carrier ILIKE 'envioclick'`
4. **Batch size:** 40 shipments por batch
5. **Backfill:** no automatico — el usuario lo dispara desde el frontend cuando este listo

---

## Arquitectura (reusa 100% lo existente)

```
Frontend
  POST /api/v1/shipments/sync-status?business_id=X&date_from=...&date_to=...
       |
       v
Shipments handler -> SyncShipmentsUseCase
  1. Query DB: SELECT id, tracking_number, metadata->envioclick_id_order
                FROM shipments
                WHERE business_id=X
                  AND status IN (...pendings...)
                  AND tracking_number IS NOT NULL
                  AND carrier ILIKE 'envioclick'
                  AND created_at BETWEEN date_from AND date_to
  2. Responde 202 Accepted con total + correlation_id (async)
  3. En goroutine: divide en batches de 40 y publica a transport.requests
       operation="sync_batch"
       payload={ tracking_numbers: [...], id_orders: [...], business_id }
       |
       v (RabbitMQ transport.requests)
Transport router -> transport.envioclick.requests
       |
       v
Envioclick request_consumer
  - Nuevo case "sync_batch":
     * Llama client.TrackByOrdersBatch (si hay id_orders) o track individual en paralelo
     * Para cada resultado:
         - Normaliza el API response a NormalizedWebhookUpdate
         - Publica TransportResponseMessage operation="webhook_update" (una por tracking)
       |
       v (RabbitMQ transport.responses)
Shipments response_consumer.handleWebhookUpdate  <-- YA IMPLEMENTADO
  - Actualiza shipment + order + delivered_at/shipped_at
  - Publica SSE PublishTrackingUpdated
       |
       v
Frontend recibe eventos SSE, refresca UI
```

### Por que NO llamar directo al API desde el handler

- El handler HTTP debe responder rapido (202 Accepted).
- El flujo async (rabbit) ya tiene retry, error handling, SSE, update de orders.
- Si el API de Envioclick tarda (40 requests × 500ms = 20s), no bloquea al usuario.
- Reutiliza handleWebhookUpdate tal cual — cero duplicacion.

---

## Estado del API de Envioclick (verificado con curl)

El API `/track` y `/track-by-orders` devuelven:

```json
{
  "data": {
    "status": "Entregado" | "En transito",
    "statusDetail": "Envio entregado" | "Salida a ruta" | "Ingresado A Bodega" | ...,
    "realPickupDate": "2026-04-14 10:00:00",
    "realDeliveryDate": "2026-04-15 15:46:00" | null,
    "arrivalDate": "2026-04-15 09:54:00"
  }
}
```

**Importante:** el API NO devuelve `statusStep`, `events[]`, ni `incidence`. Solo `status`
+ `statusDetail` + fechas. Por eso la sincronizacion usa el `status` como `statusStep`
(son compatibles: "Entregado" y "En transito" ya estan en el mapper).

### Tratamiento de `statusDetail` durante sync

Algunos `statusDetail` dan pistas de sub-estado:
- `"Salida a ruta"`, `"Ingresado A Bodega"`, `"Viajando En Ruta Regional"` -> son sub-estados de "En transito"
- `"En reparto"`, `"En distribucion"` -> sub-estados de "En distribucion" (out_for_delivery)

Propuesta: combinar `status` + `statusDetail` y pasar el mas granular al mapper:
1. Si `statusDetail` matchea algun statusStep conocido, usar statusDetail
2. Si no, usar `status`
3. El mapper normaliza y mapea

---

## Archivos a crear/modificar

### Envioclick (agrega logica de sync)

```
envioclick/internal/domain/
├── sync_batch.go                       # NUEVO — SyncBatchRequest + BatchItem
└── status_mapper.go                    # modificado — helper ApiStatusToStep
envioclick/internal/app/
└── sync_usecase.go                     # NUEVO — orchestra track + normalize
envioclick/internal/infra/primary/consumer/
└── transport_request_consumer.go       # modificado — case "sync_batch"
```

### Shipments (agrega endpoint + usecase)

```
shipments/internal/domain/
└── sync.go                             # NUEVO — SyncShipmentsRequest + SyncShipmentsResult
shipments/internal/app/usecaseshipment/
└── sync.go                             # NUEVO — UC que query + publica a transport.requests
shipments/internal/infra/primary/handlers/
├── sync-shipments.go                   # NUEVO — HTTP handler
└── router.go                           # modificado — registrar ruta
shipments/internal/infra/secondary/repository/
└── sync_queries.go                     # NUEVO — query filtrando por carrier + fechas + status
shipments/internal/infra/secondary/queue/
└── transport_request_publisher.go      # modificado — agregar PublishSyncBatch
```

---

## Contrato del endpoint

### Request
```http
POST /api/v1/shipments/sync-status
Authorization: Bearer <JWT>
Content-Type: application/json

{
  "provider": "envioclick",           // obligatorio por ahora solo envioclick
  "date_from": "2026-04-01",          // opcional, default 30 dias atras
  "date_to": "2026-04-17",            // opcional, default hoy
  "statuses": ["pending","in_transit"] // opcional, default son los "activos"
}
```

Super admin: acepta `?business_id=X` como query param. Usuario normal: toma su propio business del JWT.

### Response (202 Accepted)
```json
{
  "success": true,
  "correlation_id": "sync-<uuid>",
  "total_shipments": 67,
  "batches": 2,
  "batch_size": 40,
  "estimated_duration_seconds": 45,
  "message": "Sincronizacion iniciada. Los envios se actualizaran progresivamente."
}
```

### Error 400 — sin shipments que sincronizar
```json
{
  "success": false,
  "total_shipments": 0,
  "message": "No hay envios de envioclick para sincronizar en el rango indicado"
}
```

---

## Contrato interno `sync_batch`

Publicado por shipments, consumido por envioclick. Reusa `TransportRequestMessage`
agregando `Operation="sync_batch"`:

```json
{
  "provider": "envioclick",
  "operation": "sync_batch",
  "business_id": 35,
  "integration_id": 55,
  "correlation_id": "sync-abc-batch-1",
  "payload": {
    "items": [
      { "shipment_id": 34063, "tracking_number": "240050336731", "envioclick_id_order": 4263492 },
      { "shipment_id": 34061, "tracking_number": "240050328647", "envioclick_id_order": null },
      ...40 items
    ]
  }
}
```

El consumer de envioclick:
1. Para items con `envioclick_id_order` -> llama `TrackByOrdersBatch` (1 sola request HTTP)
2. Para los que no tienen id_order -> llama `/track` individual en paralelo (max 10 concurrent)
3. Por cada resultado, construye `WebhookUpdateMessage` y publica a `transport.responses`

---

## Helper: `ApiStatusToStep` (en status_mapper.go)

```go
// ApiStatusToStep converts the /track API response "status" + "statusDetail"
// to a statusStep that MapStatusStepToProbability can consume.
func ApiStatusToStep(status, statusDetail string) string {
    detail := normalize(statusDetail)
    // Sub-states seen from API that map to a more granular step
    switch {
    case strings.Contains(detail, "distribucion"), strings.Contains(detail, "reparto"):
        return "En Distribucion"
    case strings.Contains(detail, "recolec"):
        return "Envio Recolectado"
    case strings.Contains(detail, "entregado"):
        return "Entregado"
    }
    // Fallback to top-level status ("En transito", "Entregado")
    return status
}
```

Uso:
```go
step := ApiStatusToStep(apiResp.Status, apiResp.StatusDetail)
probStatus, unknown := MapStatusStepToProbability(step, false)
```

---

## Fases de implementacion

### Fase A — Dominio y sync usecase en envioclick
1. Crear `envioclick/internal/domain/sync_batch.go` con `SyncBatchItem`, `SyncBatchRequest`.
2. Extender `status_mapper.go` con `ApiStatusToStep(status, detail)`.
3. Crear `envioclick/internal/app/sync_usecase.go`:
   - `SyncBatch(ctx, items []SyncBatchItem, apiKey, baseURL) []NormalizedWebhookUpdate`
   - Llama TrackByOrdersBatch si hay id_orders, si no /track en paralelo
   - Normaliza cada resultado

### Fase B — Consumer envioclick: case "sync_batch"
1. En `transport_request_consumer.go`, agregar `case "sync_batch"`:
   - Resuelve apiKey + baseURL (ya existe)
   - Llama SyncBatch
   - Por cada NormalizedUpdate, usa el existente `WebhookResponsePublisher.PublishWebhookUpdate`
   - Asegura que cada mensaje lleva `ShipmentID` (del item original) y `BusinessID` (del request)

### Fase C — Shipments: query, publisher y usecase
1. `shipments/.../repository/sync_queries.go`:
   - `ListShipmentsForSync(ctx, filter) ([]SyncShipmentRow, error)`
   - SELECT s.id, s.tracking_number, s.metadata->>'envioclick_id_order' FROM shipments WHERE ...
2. `shipments/.../queue/transport_request_publisher.go`:
   - `PublishSyncBatch(ctx, provider, businessID, integrationID, correlationID, items)`
3. `shipments/.../app/usecaseshipment/sync.go`:
   - `SyncShipments(ctx, req) (*SyncResult, error)`
   - Query shipments, divide en batches de 40, publica cada batch con correlation_id derivado

### Fase D — Handler HTTP
1. `shipments/.../handlers/sync-shipments.go`:
   - Parsea request (provider, date_from, date_to, statuses)
   - Resuelve business_id (super admin vs normal)
   - Resuelve integration_id de envioclick para el business
   - Llama UC + responde 202 con totals
2. `router.go`: `shipments.POST("/sync-status", h.SyncShipmentStatus)` (con JWT)

### Fase E — Tests
1. Unit tests para `ApiStatusToStep` y `SyncBatch`.
2. Integration test: mock API envioclick + verify publish.

### Fase F — Verificacion end-to-end
1. Disparar con curl desde local:
   ```bash
   curl -X POST http://localhost:3050/api/v1/shipments/sync-status \
     -H "Authorization: Bearer $JWT" \
     -H "Content-Type: application/json" \
     -d '{"provider":"envioclick","date_from":"2026-04-01"}'
   ```
2. Verificar en DB: shipments actualizados por `handleWebhookUpdate`.
3. Verificar webhook_logs: `source='envioclick'`, `event_type='sync'`.
4. Verificar orders sincronizadas.

---

## Preguntas abiertas (no bloqueantes)

1. **Frontend:** `¿boton "Sincronizar"` en el modulo de shipments + filtros de fecha? Lo hace el usuario en otro PR.
2. **Event type en webhook_logs:** cuando viene por sync, ¿usamos `event_type="sync"` o seguimos con `"tracking_update"` para no fragmentar? Voto: `"sync"` (distinguir origenes).
3. **Rate limit al API de Envioclick:** 40 req en paralelo puede saturar. Agregar un semaforo para max 10 concurrent dentro del batch.
4. **Shipments sin `envioclick_id_order`:** los mas viejos no tienen ese metadata. Caen al path de `/track` individual. Funciona pero es mas lento.

---

## Fuera de scope (PRs futuros)

- Sync automatico (cron) cada N minutos.
- Sync para Enviame/MiPaquete cuando tengan API accesible.
- Retry automatico de webhooks fallidos (ya estan en `webhook_logs` con `status=failed`).
- Frontend con boton de sincronizacion manual + indicador de progreso.
