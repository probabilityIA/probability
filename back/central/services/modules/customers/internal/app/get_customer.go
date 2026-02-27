package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) GetClient(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
	client, err := uc.repo.GetByID(ctx, businessID, clientID)
	if err != nil {
		return nil, err
	}

	// Enriquecer con stats de órdenes
	orderCount, totalSpent, lastOrderAt, err := uc.repo.GetOrderStats(ctx, clientID)
	if err != nil {
		return nil, err
	}

	client.OrderCount = orderCount
	client.TotalSpent = totalSpent
	client.LastOrderAt = lastOrderAt

	return client, nil
}
