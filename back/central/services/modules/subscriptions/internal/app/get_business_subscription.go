package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) GetBusinessSubscription(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
	return uc.repo.GetLatestByBusinessID(ctx, businessID)
}
