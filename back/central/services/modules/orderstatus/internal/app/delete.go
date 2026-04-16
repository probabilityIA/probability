package app

import "context"

func (uc *useCase) DeleteOrderStatusMapping(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
