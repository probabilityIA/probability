package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) recordAudit(ctx context.Context, businessID, actorUserID uint, action, description string) {
	label, err := uc.repo.ResolveUserLabel(ctx, actorUserID)
	if err != nil || label == "" {
		label = "Sistema"
	}
	if err := uc.repo.CreateAuditLog(ctx, &entities.SubscriptionAuditLog{
		BusinessID:  businessID,
		ActorUserID: &actorUserID,
		ActorLabel:  label,
		Action:      action,
		Description: description,
	}); err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("business_id", businessID).
			Str("action", action).
			Msg("failed to record subscription audit log")
	}
}
