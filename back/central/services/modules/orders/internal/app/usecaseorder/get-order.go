package usecaseorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder/mapper"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
)

// GetOrderByID obtiene una orden por su ID
func (uc *UseCaseOrder) GetOrderByID(ctx context.Context, id string) (*domain.OrderResponse, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	order, err := uc.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	return mapper.ToOrderResponse(order), nil
}
