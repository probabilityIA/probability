package app

import "context"

func (uc *UseCase) ForceRedisplay(ctx context.Context, id uint) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existing.ForceRedisplay = true
	if _, err := uc.repo.Update(ctx, existing); err != nil {
		return err
	}

	return uc.repo.DeleteViewsByAnnouncementID(ctx, id)
}
