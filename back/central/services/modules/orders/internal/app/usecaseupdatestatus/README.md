# Use Case: Update Status (Cambio de Estado de Ordenes)

**Endpoint:** `PUT /api/v1/orders/:id/status`
**Paquete:** `usecaseupdatestatus`
**Fecha de creacion:** 2026-03-29

---

## Proposito

Este caso de uso es el **unico punto de entrada autorizado** para cambiar el estado de una orden desde operacion interna de Probability (UI, API manual). Implementa validacion estricta de transiciones, strategy pattern por estado, historial de cambios, y publicacion de eventos a RabbitMQ.

---

## Regla Fundamental: Quien Controla el Estado

```
Integraciones externas (Shopify, Amazon, MercadoLibre)
  -> Escriben estado SIN restriccion (son fuente de verdad de su plataforma)
  -> Pasan por: MapAndSaveOrder -> usecaseupdateorder -> updateOrderStatus()
  -> NO usan este caso de uso

Operacion interna (UI de Probability, API manual)
  -> DEBE pasar por este caso de uso (PUT /orders/:id/status)
  -> Validacion estricta de transiciones
  -> Registra historial en order_history
  -> Publica eventos a RabbitMQ
```

### Bloqueo del endpoint general

El handler `PUT /orders/:id` (update general) **no mapea los campos Status, OriginalStatus ni StatusID** desde el request HTTP. Esto impide que el frontend cambie el estado por esa via. Ver `handlers/mappers/to_domain.go`.

---

## Flujo de Ejecucion (change_status.go)

```
1. Validar que el status destino es valido (IsValid)
2. Obtener la orden de la BD (GetOrderByID)
3. Validar que el estado actual NO es terminal (IsTerminal)
4. Validar la transicion (CanTransitionTo)
5. Guardar estado anterior
6. Ejecutar strategy del estado destino (executeStrategy)
7. Resolver StatusID desde el codigo (GetOrderStatusIDByCode)
8. Persistir cambios (UpdateOrder)
9. Registrar historial (CreateOrderHistory -> tabla order_history)
10. Publicar eventos a RabbitMQ (OrderEventTypeUpdated + OrderEventTypeStatusChanged)
```

---

## Mapa de Transiciones Permitidas

```
                          +----------+
                     +---->| on_hold  |----> pending
                     |    +----------+       |
                     |                       v
                     |    +---------------------+
                     +----|      picking        |<<---- inventory_issue
                     |    +---------------------+
                     |              |
                     |              v
                     +----+---------------------+
                     |    |      packing        |
                     |    +---------------------+
                     |              |
                     |              v
                          +---------------------+
                          |   ready_to_ship     |
                          +---------------------+
                                    |
                                    v
                          +---------------------+
                          | assigned_to_driver  |<<---- delivery_novelty
                          +---------------------+
                                    |
                                    v
                          +---------------------+
                          |     picked_up       |
                          +---------------------+
                                    |
                                    v
                          +---------------------+
                          |     in_transit      |
                          +---------------------+
                                    |
                                    v
                          +---------------------+
                          |  out_for_delivery   |
                          +---------------------+
                           |    |    |    |
                +----------+    |    |    +--------------+
                v               v    v                  v
          +----------+  +----------+ +----------+ +----------+
          |delivered |  | novelty  | | rejected | | d_failed |
          +----------+  +----------+ +----------+ +----------+
           |    |                         |              |
           v    v                         v              v
      completed |                   return_in_transit <<--+
           |    |                         |
           v    v                         v
       refunded                      returned -> refunded

    ★ cancelled: accesible desde CUALQUIER estado no-terminal
    ★ Terminales (sin salida): cancelled, refunded
```

### Transiciones Detalladas

| Desde | Puede ir a |
|-------|-----------|
| `pending` | picking, on_hold, **cancelled** |
| `on_hold` | pending, picking, **cancelled** |
| `picking` | packing, inventory_issue, on_hold, **cancelled** |
| `inventory_issue` | picking, **cancelled** |
| `packing` | ready_to_ship, on_hold, **cancelled** |
| `ready_to_ship` | assigned_to_driver, on_hold, **cancelled** |
| `assigned_to_driver` | picked_up, **cancelled** |
| `picked_up` | in_transit, **cancelled** |
| `in_transit` | out_for_delivery, **cancelled** |
| `out_for_delivery` | delivered, delivery_novelty, rejected, delivery_failed, **cancelled** |
| `delivered` | completed, refunded, return_in_transit, **cancelled** |
| `delivery_novelty` | assigned_to_driver, out_for_delivery, delivery_failed, return_in_transit, **cancelled** |
| `delivery_failed` | return_in_transit, **cancelled** |
| `rejected` | return_in_transit, **cancelled** |
| `return_in_transit` | returned, **cancelled** |
| `returned` | refunded, **cancelled** |
| `completed` | refunded, **cancelled** |
| `cancelled` | (terminal — sin salida) |
| `refunded` | (terminal — sin salida) |

La logica esta en `domain/entities/order_status.go` -> `validTransitions` map + `CanTransitionTo()`.

---

## Strategy Pattern (un archivo por estado destino)

Cada estado destino tiene su propio archivo `to_<status>.go` que aplica la logica especifica:

| Archivo | Estado destino | Metadata que acepta | Campo que modifica |
|---------|---------------|--------------------|--------------------|
| `to_picking.go` | picking | — | Status |
| `to_packing.go` | packing | — | Status |
| `to_ready_to_ship.go` | ready_to_ship | — | Status |
| `to_assigned_to_driver.go` | assigned_to_driver | `driver_id` (float64), `driver_name` (string) | Status, DriverID, DriverName |
| `to_picked_up.go` | picked_up | — | Status |
| `to_in_transit.go` | in_transit | `tracking_number`, `tracking_link` (string) | Status, TrackingNumber, TrackingLink |
| `to_out_for_delivery.go` | out_for_delivery | — | Status |
| `to_delivered.go` | delivered | — | Status, DeliveredAt (auto = now) |
| `to_delivery_novelty.go` | delivery_novelty | `reason` (string) | Status, Novelty |
| `to_delivery_failed.go` | delivery_failed | `reason` (string) | Status, Notes |
| `to_rejected.go` | rejected | `reason` (string) | Status, Notes |
| `to_return_in_transit.go` | return_in_transit | `tracking_number` (string) | Status, TrackingNumber |
| `to_returned.go` | returned | — | Status |
| `to_inventory_issue.go` | inventory_issue | `notes` (string) | Status, Novelty |
| `to_on_hold.go` | on_hold | `reason` (string) | Status, Notes |
| `to_cancelled.go` | cancelled | `reason` (string) | Status, Notes |
| `to_completed.go` | completed | — | Status |
| `to_refunded.go` | refunded | — | Status |
| `to_failed.go` | failed | — | Status |

---

## Request y Response

### Request (PUT /orders/:id/status)

```json
{
  "status": "assigned_to_driver",
  "metadata": {
    "driver_id": 1,
    "driver_name": "Carlos Ramirez"
  }
}
```

- `status` (string, required): codigo del estado destino
- `metadata` (object, optional): datos adicionales segun el estado

### Response (HTTP 200)

La orden completa actualizada en formato snake_case.

### Errores

| HTTP | Error | Cuando |
|------|-------|--------|
| 400 | `invalid status code: X` | Status no existe en el sistema |
| 400 | `invalid request` | Body vacio o mal formado |
| 404 | `order not found` | UUID no existe |
| 422 | `invalid status transition: cannot transition from X to Y` | Transicion no permitida |
| 422 | `order is in a terminal state: current status is X` | Orden en cancelled/refunded |

---

## Historial (order_history)

Cada cambio de estado registra una fila en `order_history`:

| Campo | Descripcion |
|-------|-------------|
| `order_id` | UUID de la orden |
| `previous_status` | Estado antes del cambio |
| `new_status` | Estado despues del cambio |
| `changed_by` | User ID del JWT (puede ser null) |
| `changed_by_name` | Nombre del usuario (desnormalizado) |
| `reason` | Extraido de `metadata.reason` si existe |
| `metadata` | JSON completo del metadata enviado |

---

## Eventos RabbitMQ

Al cambiar estado exitosamente se publican 2 eventos al exchange fanout:

1. **`order.updated`** — notifica que la orden cambio (consumers: invoicing, inventory, score, whatsapp, events)
2. **`order.status_changed`** — incluye `previous_status` y `current_status` (consumers: notificaciones, tracking)

Ambos se publican en goroutines (fire-and-forget con log de error).

---

## Dependencias

```go
type UseCaseUpdateStatus struct {
    repo                 ports.IRepository           // CRUD + GetOrderStatusIDByCode + CreateOrderHistory
    logger               log.ILogger                 // Logging
    rabbitEventPublisher ports.IOrderRabbitPublisher  // Eventos RabbitMQ
}
```

---

## Relacion con Otros Use Cases

| Use Case | Proposito | Cambia status | Valida transicion |
|----------|-----------|--------------|-------------------|
| **usecaseupdatestatus** (este) | Cambio de estado manual/UI | SI | SI (estricto) |
| usecaseupdateorder | Update general de orden | SI (solo integraciones) | NO (fuente externa) |
| usecasecreateorder | Crear orden nueva | SI (status inicial) | NO (es creacion) |

### Por que usecaseupdateorder no valida transiciones

Las integraciones externas (Shopify, Amazon) son fuente de verdad de su plataforma. Si Shopify dice que una orden paso de `pending` a `delivered` directamente, Probability debe reflejar eso. La validacion de transiciones es para operacion interna solamente.

Cuando usecaseupdateorder detecta un salto que violaria el flujo v2, registra un **warning en log** pero acepta el cambio.

---

## Testing

Tests unitarios en `change_status_test.go`. Plan de pruebas E2E en `/.claude/test/order-status-v2-testing.md`.

### Flujos probados (E2E)

1. Happy path completo (pending -> ... -> completed)
2. Cancelacion desde picking
3. Novedad de entrega + reintento con nuevo piloto
4. Rechazo + devolucion + reembolso
5. Entrega fallida + devolucion
6. On hold + reanudacion
7. Problema de inventario + recuperacion
8. Reembolso post-entrega
9. Reembolso post-completado
10. Transicion invalida (expect 422)
11. Estado terminal bloqueado (expect 422)

---

## Archivos del Paquete

```
usecaseupdatestatus/
+-- README.md                  # Este archivo
+-- constructor.go             # New() + struct
+-- change_status.go           # Orquestador principal + executeStrategy + saveOrderHistory
+-- change_status_test.go      # Tests unitarios
+-- publish_events.go          # Publicacion a RabbitMQ
+-- to_picking.go              # Strategy: -> picking
+-- to_packing.go              # Strategy: -> packing
+-- to_ready_to_ship.go        # Strategy: -> ready_to_ship
+-- to_assigned_to_driver.go   # Strategy: -> assigned_to_driver (metadata: driver_id, driver_name)
+-- to_picked_up.go            # Strategy: -> picked_up
+-- to_in_transit.go           # Strategy: -> in_transit (metadata: tracking_number, tracking_link)
+-- to_out_for_delivery.go     # Strategy: -> out_for_delivery
+-- to_delivered.go            # Strategy: -> delivered (auto: delivered_at = now)
+-- to_delivery_novelty.go     # Strategy: -> delivery_novelty (metadata: reason -> novelty)
+-- to_delivery_failed.go      # Strategy: -> delivery_failed (metadata: reason -> notes)
+-- to_rejected.go             # Strategy: -> rejected (metadata: reason -> notes)
+-- to_return_in_transit.go    # Strategy: -> return_in_transit (metadata: tracking_number)
+-- to_returned.go             # Strategy: -> returned
+-- to_inventory_issue.go      # Strategy: -> inventory_issue (metadata: notes -> novelty)
+-- to_on_hold.go              # Strategy: -> on_hold (metadata: reason -> notes)
+-- to_cancelled.go            # Strategy: -> cancelled (metadata: reason -> notes)
+-- to_completed.go            # Strategy: -> completed
+-- to_refunded.go             # Strategy: -> refunded
+-- to_failed.go               # Strategy: -> failed
```
