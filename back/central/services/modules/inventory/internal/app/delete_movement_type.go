package app

import "context"

func (uc *useCase) DeleteMovementType(ctx context.Context, id uint) error {
	return uc.repo.DeleteMovementType(ctx, id)
}
