package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

func (uc *UseCase) GetClientPricingRule(ctx context.Context, businessID, ruleID uint) (*entities.ClientPricingRule, error) {
	return uc.repo.GetClientPricingRule(ctx, businessID, ruleID)
}
