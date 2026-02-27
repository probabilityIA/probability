package app

import "context"

func (uc *useCase) DeleteOrderStatus(ctx context.Context, id uint) error {
	return uc.repo.DeleteOrderStatus(ctx, id)
}
