package app

import "context"

func (uc *useCase) DeleteChannelStatus(ctx context.Context, id uint) error {
	return uc.repo.DeleteChannelStatus(ctx, id)
}
