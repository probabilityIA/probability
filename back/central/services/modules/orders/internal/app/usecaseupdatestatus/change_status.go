package usecaseupdatestatus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder/mapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
)

// ChangeStatus es el orquestador principal para cambios de estado de órdenes.
// Valida la transición, delega la lógica específica al strategy correspondiente,
// persiste los cambios, registra el historial y publica eventos.
func (uc *UseCaseUpdateStatus) ChangeStatus(ctx context.Context, orderID string, req *dtos.ChangeStatusRequest) (*dtos.OrderResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}
	if req.Status == "" {
		return nil, domainerrors.ErrInvalidStatus
	}

	// 1. Validar que el estado destino es válido
	targetStatus := entities.OrderStatus(req.Status)
	if !targetStatus.IsValid() {
		return nil, fmt.Errorf("%w: %s", domainerrors.ErrInvalidStatus, req.Status)
	}

	// 2. Obtener la orden
	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	// 3. Validar que el estado actual no es terminal
	currentStatus := entities.OrderStatus(order.Status)
	if currentStatus.IsTerminal() {
		return nil, fmt.Errorf("%w: current status is %s", domainerrors.ErrOrderInTerminalState, order.Status)
	}

	// 4. Validar la transición
	if !currentStatus.CanTransitionTo(targetStatus) {
		return nil, fmt.Errorf("%w: cannot transition from %s to %s", domainerrors.ErrInvalidStatusTransition, order.Status, req.Status)
	}

	// 5. Guardar estado anterior
	previousStatus := order.Status

	// 6. Delegar al strategy correspondiente
	uc.executeStrategy(order, req)

	// 7. Resolver StatusID desde el código
	statusID, err := uc.repo.GetOrderStatusIDByCode(ctx, req.Status)
	if err != nil {
		uc.logger.Warn(ctx).
			Err(err).
			Str("order_id", orderID).
			Str("status_code", req.Status).
			Msg("No se pudo resolver status_id para el código de estado")
	}
	if statusID != nil {
		order.StatusID = statusID
	}

	// 8. Persistir cambios
	if err := uc.repo.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	// 9. Registrar historial de cambio de estado
	uc.saveOrderHistory(ctx, order, previousStatus, req)

	// 10. Publicar eventos
	uc.publishStatusChangeEvents(ctx, order, previousStatus)

	uc.logger.Info(ctx).
		Str("order_id", orderID).
		Str("previous_status", previousStatus).
		Str("new_status", req.Status).
		Msg("Estado de orden actualizado")

	return mapper.ToOrderResponse(order), nil
}

// executeStrategy delega al método específico según el estado destino
func (uc *UseCaseUpdateStatus) executeStrategy(order *entities.ProbabilityOrder, req *dtos.ChangeStatusRequest) {
	switch entities.OrderStatus(req.Status) {
	case entities.OrderStatusPicking:
		uc.toPicking(order, req)
	case entities.OrderStatusPacking:
		uc.toPacking(order, req)
	case entities.OrderStatusReadyToShip:
		uc.toReadyToShip(order, req)
	case entities.OrderStatusAssignedToDriver:
		uc.toAssignedToDriver(order, req)
	case entities.OrderStatusPickedUp:
		uc.toPickedUp(order, req)
	case entities.OrderStatusInTransit:
		uc.toInTransit(order, req)
	case entities.OrderStatusOutForDelivery:
		uc.toOutForDelivery(order, req)
	case entities.OrderStatusDelivered:
		uc.toDelivered(order, req)
	case entities.OrderStatusDeliveryNovelty:
		uc.toDeliveryNovelty(order, req)
	case entities.OrderStatusDeliveryFailed:
		uc.toDeliveryFailed(order, req)
	case entities.OrderStatusRejected:
		uc.toRejected(order, req)
	case entities.OrderStatusReturnInTransit:
		uc.toReturnInTransit(order, req)
	case entities.OrderStatusReturned:
		uc.toReturned(order, req)
	case entities.OrderStatusInventoryIssue:
		uc.toInventoryIssue(order, req)
	case entities.OrderStatusCancelled:
		uc.toCancelled(order, req)
	case entities.OrderStatusOnHold:
		uc.toOnHold(order, req)
	case entities.OrderStatusCompleted:
		uc.toCompleted(order, req)
	case entities.OrderStatusRefunded:
		uc.toRefunded(order, req)
	case entities.OrderStatusFailed:
		uc.toFailed(order, req)
	default:
		order.Status = req.Status
	}
}

// saveOrderHistory registra el cambio de estado en la tabla de historial
func (uc *UseCaseUpdateStatus) saveOrderHistory(ctx context.Context, order *entities.ProbabilityOrder, previousStatus string, req *dtos.ChangeStatusRequest) {
	var reason *string
	if req.Metadata != nil {
		if r, ok := req.Metadata["reason"].(string); ok {
			reason = &r
		}
	}

	var metadataBytes []byte
	if req.Metadata != nil {
		metadataBytes, _ = json.Marshal(req.Metadata)
	}

	history := &entities.OrderHistory{
		OrderID:        order.ID,
		PreviousStatus: previousStatus,
		NewStatus:      req.Status,
		ChangedBy:      req.UserID,
		ChangedByName:  req.UserName,
		Reason:         reason,
		Metadata:       metadataBytes,
	}

	if err := uc.repo.CreateOrderHistory(ctx, history); err != nil {
		uc.logger.Error(ctx).
			Err(err).
			Str("order_id", order.ID).
			Str("previous_status", previousStatus).
			Str("new_status", req.Status).
			Msg("Error al guardar historial de cambio de estado")
	}
}
