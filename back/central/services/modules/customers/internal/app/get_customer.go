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

	summary, err := uc.repo.GetCustomerSummary(ctx, businessID, clientID)
	if err != nil {
		return client, nil
	}

	if summary != nil {
		client.OrderCount = int64(summary.TotalOrders)
		client.TotalSpent = summary.TotalSpent
		client.LastOrderAt = summary.LastOrderAt
	}

	return client, nil
}
