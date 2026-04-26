# Submodulo Envios Contra Entrega (COD)

## Contexto

Hoy las ordenes con `CodTotal > 0` (pago contra entrega) se mezclan con el resto en `/orders` y `/shipments`. No hay vista dedicada para que el negocio vea cuanto le adeudan los couriers ni para conciliar cobros. Necesitamos una pestana "Contra entrega" en el navbar de Ordenes que liste los envios COD, muestre orden + guia + transportadora + mapa origen/destino + historial de tracking, y permita marcar el cobro reusando el flujo de pagos existente (`Order.IsPaid` + registro en `payments`).

## Decisiones tomadas

- **Identificacion COD**: `Order.CodTotal > 0` (sin migracion ni seed nuevo).
- **Estado de cobro**: reusar `Order.IsPaid` + `Order.PaidAt`. No hay campos COD nuevos en shipments.
- **Marcar recaudado**: crear `Payment` (status=completed, amount=CodTotal) + setear `Order.IsPaid=true`, `Order.PaidAt=now()`.
- **Validacion**: solo permitir recaudo si `shipment.status='delivered'`.
- **Alcance**: extender modulo shipments existente + nueva pestana frontend.

## Backend - `back/central/services/modules/shipments/`

### Domain
- `internal/domain/entities.go`: agregar a `ShipmentResponse` campos derivados de Order: `cod_total *float64`, `is_paid bool`, `paid_at *time.Time`, `payment_method_code string`. (No tocar entidad Shipment).
- `internal/domain/ports.go`: agregar a la interfaz del repo:
  - `ListCODShipments(ctx, businessID, filters PaginationParams + status?, isPaid?) (PaginatedResponse, error)`
  - `MarkOrderPaid(ctx, orderID string, amount float64, notes string, userID uint) error` (transaccion: insert payment + update order).

### Repository (`internal/infra/secondary/repository/`)
- `cod_queries.go` (archivo nuevo): query con JOIN a `orders` filtrando `orders.cod_total > 0`, preload de payment_method para devolver el `code`. Reusar paginacion de `dtos.PaginationParams`.
- `mark_paid.go` (archivo nuevo): transaccion GORM:
  1. `INSERT INTO payments (order_id, amount, status='completed', payment_method_id, paid_at, notes, created_by)`
  2. `UPDATE orders SET is_paid=true, paid_at=NOW() WHERE id=?`
  - Usa modelos de `migration/shared/models/payment.go` y `order.go`. NUNCA `.Table()`.

### Use cases (`internal/app/usecases/`)
- `list_cod_shipments.go`: invoca repo, mapea a response.
- `collect_cod.go`: valida que `shipment.status == "delivered"` y `order.is_paid == false` y `order.cod_total > 0`; llama `MarkOrderPaid`. Errores en `internal/domain/errors.go`: `ErrShipmentNotDelivered`, `ErrOrderAlreadyPaid`, `ErrOrderNotCOD`.

### Handlers (`internal/infra/primary/handlers/`)
- `list_cod_handler.go`: `GET /api/v1/shipments/cod?page&page_size&status&is_paid`. Soporta super admin con `?business_id=` (ver `resolveBusinessID` patron en customers).
- `collect_cod_handler.go`: `POST /api/v1/shipments/:id/collect-cod` body `{ notes?: string }`.
- Registrar en `routes.go`.

### Bundle
- `bundle.go`: instanciar nuevos usecases y handlers; pasar `OrderRepo`/`PaymentRepo` no - replicar consultas localmente segun convencion (`backend-conventions.md` regla 1).

## Frontend - `front/central/src/`

### Navbar
- `shared/ui/orders-subnavbar.tsx`: agregar tab "Contra entrega" en seccion OPERACIONES -> ruta `/shipments/cod`. Reusar el patron de permisos existente.

### Modulo (`services/modules/shipments/`)
- `domain/types.ts`: agregar `CODShipment` (extiende Shipment con `codTotal`, `isPaid`, `paidAt`, `paymentMethodCode`). `CollectCODRequest { notes?: string }`.
- `domain/ports.ts`: `listCODShipments(params)`, `collectCOD(shipmentId, payload)`.
- `infra/repository/api-repository.ts`: implementar las 2 llamadas (`/shipments/cod`, `/shipments/:id/collect-cod`).
- `infra/actions/index.ts`: server action `collectCODAction` con `revalidatePath('/shipments/cod')`.
- `app/use-cases.ts`: orquestar listado y cobro.

### UI
- `app/shipments/cod/page.tsx` (nuevo, Server Component): fetch inicial paginado (`cache: 'no-store'`).
- `services/modules/shipments/ui/components/CODShipmentList.tsx` (Client Component): split panel 1/3 + 2/3 igual que `ShipmentList.tsx` actual:
  - Izquierda: cards con badge `Por cobrar`/`Recaudado` (segun `isPaid`), monto `codTotal`, cliente, ciudad, tracking, fecha.
  - Filtros: estado envio, pagado/no pagado.
  - Derecha (`CODShipmentDetail.tsx`): orden + guia + transportadora + monto + `MiniAddressMap` origen y destino (reusar `services/modules/shipments/ui/components/MiniAddressMap.tsx`) + historial de tracking (reusar logica de `TrackingPanel.tsx`) + boton "Marcar como recaudado" (deshabilitado si `status != delivered` o `isPaid`).
  - Modal `CollectCODModal.tsx`: confirma monto a cobrar (readonly), input notas opcionales, dispara `collectCODAction`.

## Archivos criticos

- Backend (modificar):
  - `back/central/services/modules/shipments/internal/domain/entities.go`
  - `back/central/services/modules/shipments/internal/domain/ports.go`
  - `back/central/services/modules/shipments/internal/domain/errors.go`
  - `back/central/services/modules/shipments/internal/infra/primary/handlers/routes.go`
  - `back/central/services/modules/shipments/bundle.go`
- Backend (crear):
  - `internal/infra/secondary/repository/cod_queries.go`
  - `internal/infra/secondary/repository/mark_paid.go`
  - `internal/app/usecases/list_cod_shipments.go`
  - `internal/app/usecases/collect_cod.go`
  - `internal/infra/primary/handlers/list_cod_handler.go`
  - `internal/infra/primary/handlers/collect_cod_handler.go`
- Frontend (modificar):
  - `front/central/src/shared/ui/orders-subnavbar.tsx`
  - `front/central/src/services/modules/shipments/domain/{types,ports}.ts`
  - `front/central/src/services/modules/shipments/infra/repository/api-repository.ts`
  - `front/central/src/services/modules/shipments/infra/actions/index.ts`
  - `front/central/src/services/modules/shipments/app/use-cases.ts`
- Frontend (crear):
  - `front/central/src/app/shipments/cod/page.tsx`
  - `services/modules/shipments/ui/components/CODShipmentList.tsx`
  - `services/modules/shipments/ui/components/CODShipmentDetail.tsx`
  - `services/modules/shipments/ui/components/CollectCODModal.tsx`

## Reuso

- `MiniAddressMap` (Leaflet) ya existente.
- `TrackingPanel` existente para timeline.
- `dtos.PaginationParams` / `PaginatedResponse` Go.
- Modelos `migration/shared/models/{order,payment,shipment}.go` (sin nuevas migraciones).
- Patron super admin `resolveBusinessID` (ver customers module).

## Verificacion

1. `cd back/central && go build ./... && go test ./services/modules/shipments/...`
2. Backend: `./scripts/dev-services.sh restart backend`
3. Test E2E backend (con `mcp__fetch__fetch`):
   - Login como business: `POST /auth/login` con credenciales de `test-credentials.md`.
   - `GET /shipments/cod?page=1&page_size=10` -> verifica que solo aparecen ordenes con `cod_total > 0`.
   - Crear orden COD via API + generar guia + simular delivered (UPDATE shipment.status='delivered' via API existente).
   - `POST /shipments/:id/collect-cod` -> 200; verificar con MCP postgres: `SELECT is_paid, paid_at FROM orders WHERE id=?` y `SELECT * FROM payments WHERE order_id=?`.
   - Reintentar collect -> debe fallar con `ErrOrderAlreadyPaid`.
   - Probar collect en envio no entregado -> `ErrShipmentNotDelivered`.
4. Frontend: `./scripts/dev-services.sh restart frontend`, navegar a `/shipments/cod`, verificar tab visible, listado, detalle con mapa, marcar recaudado y validar refresco del badge a "Recaudado".
