package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) ListSubscriptionTypes(ctx context.Context, activeOnly bool) ([]entities.SubscriptionType, error) {
	return uc.repo.ListSubscriptionTypes(ctx, activeOnly)
}
