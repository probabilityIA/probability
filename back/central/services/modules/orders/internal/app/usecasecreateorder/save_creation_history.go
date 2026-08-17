package usecasecreateorder

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseCreateOrder) saveCreationHistory(ctx context.Context, order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) {
	if order == nil || order.ID == "" {
		return
	}

	source := entities.StatusSourceSalesChannel
	if dto.IsManualOrder {
		source = entities.StatusSourceUser
	}

	reason := creationReason(dto)

	history := &entities.OrderHistory{
		OrderID:        order.ID,
		PreviousStatus: "",
		NewStatus:      order.Status,
		ChangedBy:      order.UserID,
		ChangedByName:  creationAuthor(order, dto),
		Source:         source,
		Reason:         &reason,
	}

	if err := uc.repo.CreateOrderHistory(ctx, history); err != nil {
		uc.logger.Warn(ctx).Err(err).
			Str("order_id", order.ID).
			Str("status", order.Status).
			Msg("No se pudo registrar la creacion en el historial de la orden")
	}
}

func creationAuthor(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) string {
	if name := strings.TrimSpace(order.UserName); name != "" {
		return name
	}
	if dto.IsManualOrder {
		return "Creacion manual"
	}
	if platform := strings.TrimSpace(order.Platform); platform != "" {
		return platform
	}
	return "Sistema"
}

func creationReason(dto *dtos.ProbabilityOrderDTO) string {
	if dto.IsManualOrder {
		return "Orden creada manualmente"
	}
	if platform := strings.TrimSpace(dto.Platform); platform != "" {
		return "Orden importada desde " + platform
	}
	return "Orden creada"
}
