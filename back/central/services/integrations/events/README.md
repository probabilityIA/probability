# integrations/events

Módulo publicador de eventos de sincronización. Su responsabilidad es notificar al resto del sistema cuando una integración procesa órdenes desde plataformas externas (Shopify, WooCommerce, etc.).

---

## Qué hace

Expone funciones globales (`PublishSyncOrderCreated`, `PublishSyncCompleted`, etc.) que cualquier use case de integración puede llamar sin conocer la infraestructura. Internamente serializa el evento a JSON y lo publica al canal Redis correspondiente.

También expone dos endpoints HTTP propios:
- `GET /integrations/events/sse/:businessID` — conexión SSE directa al EventManager interno
- `GET /integrations/events/sync-status/:integrationID` — consulta si una sincronización está en curso

---

## Estructura

```
integrations/events/
├── bundle.go                               # Inicialización + instancia global + funciones Publish*
└── internal/
    ├── app/
    │   └── event_service.go                # IntegrationEventService (implementa IIntegrationEventService)
    ├── domain/
    │   ├── ports.go                        # IIntegrationEventService, IIntegrationEventPublisher
    │   ├── dtos.go                         # SyncOrderCreatedEvent, SyncCompletedEvent, SyncParams, etc.
    │   ├── sync_state.go                   # SyncState, SyncStatus
    │   └── connection.go                   # IntegrationSSEFilter, ConnectionManager
    └── infra/
        ├── primary/
        │   ├── handlers/
        │   │   └── sse_handler.go          # HandleSSE(), GetSyncStatus()
        │   └── routes.go
        └── secondary/
            ├── events/                     # EventManager SSE in-memory (fallback)
            └── redis/
                └── publisher.go            # IntegrationEventRedisPublisher
```

---

## Flujo de publicación

```
Use case de integración
    │
    └── events.PublishSyncOrderCreated(ctx, integrationID, businessID, data)
                │
                ├── redisPublisher != nil
                │       └── IntegrationEventRedisPublisher.Publish()
                │               └── json.Marshal → client.Publish(canal Redis)
                │
                └── fallback (Redis no disponible)
                        └── eventServiceInstance.PublishSync* → EventManager SSE local
```

El fallback in-memory garantiza que el módulo funcione aunque Redis no esté disponible, con la limitación de que el evento solo llega a clientes conectados a `/integrations/events/sse`.

---

## Tipos de evento

| Tipo | Cuándo se publica |
|------|-------------------|
| `integration.sync.started` | Al comenzar una sincronización |
| `integration.sync.order.created` | Orden nueva creada exitosamente |
| `integration.sync.order.updated` | Orden existente actualizada |
| `integration.sync.order.rejected` | Orden descartada (duplicada, inválida, etc.) |
| `integration.sync.completed` | Sincronización finalizada con éxito |
| `integration.sync.failed` | Sincronización fallida con error |

---

## Canal Redis

| Canal | Constante |
|-------|-----------|
| `probability:integrations:orders:sync:events` | `redisclient.ChannelIntegrationsSyncOrders` |

---

## Uso desde un use case

```go
import "github.com/secamc93/probability/back/central/services/integrations/events"

// Al sincronizar órdenes de Shopify:
events.PublishSyncStarted(ctx, integrationID, &businessID, events.SyncStartedEvent{
    IntegrationID:   integrationID,
    IntegrationType: "shopify",
    Params:          syncParams,
    StartedAt:       time.Now(),
})

for _, order := range orders {
    events.PublishSyncOrderCreated(ctx, integrationID, &businessID, events.SyncOrderCreatedEvent{
        OrderID:     order.ID,
        OrderNumber: order.Number,
        ExternalID:  order.ExternalID,
        Platform:    "shopify",
        Status:      order.Status,
        CreatedAt:   order.CreatedAt,
        SyncedAt:    time.Now(),
    })
}

events.PublishSyncCompleted(ctx, integrationID, &businessID, events.SyncCompletedEvent{
    TotalOrders:    total,
    CreatedOrders:  created,
    UpdatedOrders:  updated,
    RejectedOrders: rejected,
    Duration:       time.Since(start),
    CompletedAt:    time.Now(),
})
```

No es necesario instanciar nada. `New()` en `bundle.go` inicializa la instancia global al arrancar el servidor.

---

## Relación con modules/events

Este módulo **publica**. `modules/events` **consume**.

La comunicación es unidireccional a través de Redis: este módulo no conoce ni importa `modules/events`. La entrega final al frontend (SSE) es responsabilidad exclusiva de `modules/events`.

Ver [modules/events/README.md](../../../modules/events/README.md) para entender cómo se procesa el evento una vez publicado.
