package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) GetBusinessSubscriptionMeta(ctx context.Context, businessID uint) (*entities.BusinessSubscriptionMeta, error) {
	return uc.repo.GetBusinessSubscriptionMeta(ctx, businessID)
}
