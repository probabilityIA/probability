package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) RevertPayment(ctx context.Context, subscriptionID uint, actorUserID uint) (*entities.BusinessSubscription, error) {
	reverted, err := uc.repo.RevertSubscriptionAndRecalculate(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	uc.recordAudit(ctx, reverted.BusinessID, actorUserID, entities.AuditActionPaymentReverted,
		fmt.Sprintf("revirtio un pago de %.0f (%s)", reverted.Amount, reverted.SubscriptionTypeName))

	return reverted, nil
}
