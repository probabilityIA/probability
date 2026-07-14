package app

import "context"

func (uc *UseCase) DeleteSubscriptionType(ctx context.Context, id uint) error {
	return uc.repo.DeleteSubscriptionType(ctx, id)
}
