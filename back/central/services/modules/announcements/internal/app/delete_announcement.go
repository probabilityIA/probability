package app

import (
	"context"
	"fmt"
)

func (uc *UseCase) DeleteAnnouncement(ctx context.Context, id uint) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	for _, img := range existing.Images {
		if delErr := uc.storage.DeleteFile(ctx, img.ImageURL); delErr != nil {
			uc.log.Warn(ctx).Err(delErr).Str("url", img.ImageURL).Msg("failed to delete image from S3")
		}
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx).Err(err).Uint("announcement_id", id).Msg("failed to delete announcement")
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	return nil
}
