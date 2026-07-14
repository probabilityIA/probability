package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) GetSubscriptionType(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
	return uc.repo.GetSubscriptionType(ctx, id)
}
