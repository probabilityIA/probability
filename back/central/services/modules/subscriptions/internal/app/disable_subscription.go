package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) DisableSubscription(ctx context.Context, businessID uint) error {
	now := time.Now()
	return uc.repo.UpdateBusinessSubscriptionStatus(ctx, businessID, entities.BusinessStatusCancelled, &now)
}
