package usecaseorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// GetOrderRaw obtiene los datos crudos de una orden
func (uc *UseCaseOrder) GetOrderRaw(ctx context.Context, id string) (*domain.OrderRawResponse, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	metadata, err := uc.repo.GetOrderRaw(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting raw order data: %w", err)
	}

	if metadata == nil {
		return nil, errors.New("raw data not found for this order")
	}

	return &domain.OrderRawResponse{
		OrderID:       metadata.OrderID,
		ChannelSource: metadata.ChannelSource,
		RawData:       metadata.RawData,
	}, nil
}
