package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) DisableSubscription(ctx context.Context, businessID uint, actorUserID uint) error {
	now := time.Now()
	if err := uc.repo.UpdateBusinessSubscriptionStatus(ctx, businessID, entities.BusinessStatusCancelled, &now); err != nil {
		return err
	}
	uc.recordAudit(ctx, businessID, actorUserID, entities.AuditActionSubscriptionSuspended, "suspendio la suscripcion")
	return nil
}
