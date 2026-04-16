package usecaseorder

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

// GetOrderHistory obtiene el historial de cambios de estado de una orden
func (uc *UseCaseOrder) GetOrderHistory(ctx context.Context, orderID string) ([]dtos.OrderHistoryResponse, error) {
	if orderID == "" {
		return nil, errors.New("order ID is required")
	}

	history, err := uc.repo.GetOrderHistory(ctx, orderID)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.OrderHistoryResponse, len(history))
	for i, h := range history {
		result[i] = dtos.OrderHistoryResponse{
			ID:             h.ID,
			CreatedAt:      h.CreatedAt,
			OrderID:        h.OrderID,
			PreviousStatus: h.PreviousStatus,
			NewStatus:      h.NewStatus,
			ChangedBy:      h.ChangedBy,
			ChangedByName:  h.ChangedByName,
			Reason:         h.Reason,
		}
	}

	return result, nil
}
