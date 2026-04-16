package app

import "context"

func (uc *UseCase) DeleteQuantityDiscount(ctx context.Context, businessID, discountID uint) error {
	return uc.repo.DeleteQuantityDiscount(ctx, businessID, discountID)
}
