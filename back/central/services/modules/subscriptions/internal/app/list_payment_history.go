package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) ListPaymentHistory(ctx context.Context, businessID uint) ([]entities.BusinessSubscription, error) {
	return uc.repo.ListByBusinessID(ctx, businessID)
}

func (uc *UseCase) ListAuditLogs(ctx context.Context, businessID uint, limit int) ([]entities.SubscriptionAuditLog, error) {
	return uc.repo.ListAuditLogsByBusiness(ctx, businessID, limit)
}
