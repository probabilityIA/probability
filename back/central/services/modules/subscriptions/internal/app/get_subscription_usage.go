package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) GetSubscriptionUsage(ctx context.Context, businessID uint) (*entities.SubscriptionUsage, error) {
	return uc.repo.GetSubscriptionUsage(ctx, businessID)
}
