package usecaseorder

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

// GetOrderRaw obtiene los datos crudos de una orden
func (uc *UseCaseOrder) GetOrderRaw(ctx context.Context, id string) (*dtos.OrderRawResponse, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	metadata, err := uc.repo.GetOrderRaw(ctx, id)
	if err != nil {
		// Preservar el mensaje de error específico para "not found"
		if err.Error() == "raw data not found for this order" {
			return nil, errors.New("raw data not found for this order")
		}
		return nil, fmt.Errorf("error getting raw order data: %w", err)
	}

	if metadata == nil {
		return nil, errors.New("raw data not found for this order")
	}

	return &dtos.OrderRawResponse{
		OrderID:       metadata.OrderID,
		ChannelSource: metadata.ChannelSource,
		RawData:       metadata.RawData,
	}, nil
}
