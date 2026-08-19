package app

import (
	"context"

	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) ReactivateSubscription(ctx context.Context, businessID uint, actorUserID uint) error {
	latest, err := uc.repo.GetLatestByBusinessID(ctx, businessID)
	if err != nil {
		return err
	}
	if latest == nil {
		return errs.ErrNothingToReactivate
	}

	if err := uc.repo.UpdateBusinessSubscriptionStatus(ctx, businessID, entities.BusinessStatusActive, &latest.EndDate); err != nil {
		return err
	}

	uc.recordAudit(ctx, businessID, actorUserID, entities.AuditActionSubscriptionReactivated, "reactivo la suscripcion")
	return nil
}
