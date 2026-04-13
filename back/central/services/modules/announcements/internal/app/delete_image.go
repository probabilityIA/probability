package app

import (
	"context"
	"fmt"
)

func (uc *UseCase) DeleteImage(ctx context.Context, announcementID, imageID uint) error {
	img, err := uc.repo.GetImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	if img.AnnouncementID != announcementID {
		return fmt.Errorf("image does not belong to this announcement")
	}

	if img.ImageURL != "" {
		_ = uc.storage.DeleteFile(ctx, img.ImageURL)
	}

	return uc.repo.DeleteImageByID(ctx, imageID)
}
