package usecaseupdatestatus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder/mapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
)

func (uc *UseCaseUpdateStatus) ChangeStatus(ctx context.Context, orderID string, req *dtos.ChangeStatusRequest) (*dtos.OrderResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}
	if req.Status == "" {
		return nil, domainerrors.ErrInvalidStatus
	}

	targetStatus := entities.OrderStatus(req.Status)
	if !targetStatus.IsValid() {
		return nil, fmt.Errorf("%w: %s", domainerrors.ErrInvalidStatus, req.Status)
	}

	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	currentStatus := entities.OrderStatus(order.Status)
	if currentStatus.IsTerminal() {
		return nil, fmt.Errorf("%w: current status is %s", domainerrors.ErrOrderInTerminalState, order.Status)
	}

	if !currentStatus.CanTransitionTo(targetStatus) {
		return nil, fmt.Errorf("%w: cannot transition from %s to %s", domainerrors.ErrInvalidStatusTransition, order.Status, req.Status)
	}

	if currentStatus == entities.OrderStatusInventoryIssue && targetStatus == entities.OrderStatusPicking {
		if err := uc.validateStockForPicking(ctx, order); err != nil {
			return nil, err
		}
	}

	previousStatus := order.Status

	uc.executeStrategy(order, req)

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

	order.StatusSource = entities.StatusSourceUser
	order.StatusChangedBy = req.UserName
	now := time.Now()
	order.StatusChangedAt = &now

	if err := uc.repo.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	uc.saveOrderHistory(ctx, order, previousStatus, req)

	uc.publishStatusChangeEvents(ctx, order, previousStatus)

	uc.logger.Info(ctx).
		Str("order_id", orderID).
		Str("previous_status", previousStatus).
		Str("new_status", req.Status).
		Msg("Estado de orden actualizado")

	return mapper.ToOrderResponse(order), nil
}

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
	case entities.OrderStatusCancelRequested:
		uc.toCancelRequested(order, req)
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

	changedByName := req.UserName
	if changedByName == "" && req.UserID != nil && *req.UserID > 0 {
		changedByName = uc.repo.GetUserDisplayName(ctx, *req.UserID)
	}

	history := &entities.OrderHistory{
		OrderID:        order.ID,
		PreviousStatus: previousStatus,
		NewStatus:      req.Status,
		ChangedBy:      req.UserID,
		ChangedByName:  changedByName,
		Source:         entities.StatusSourceUser,
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
