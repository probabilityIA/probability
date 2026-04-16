package app

import "context"

func (uc *UseCase) DeleteClientPricingRule(ctx context.Context, businessID, ruleID uint) error {
	return uc.repo.DeleteClientPricingRule(ctx, businessID, ruleID)
}
