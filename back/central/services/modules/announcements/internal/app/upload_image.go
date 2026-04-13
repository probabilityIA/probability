package app

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func (uc *UseCase) UploadImage(ctx context.Context, announcementID uint, file *multipart.FileHeader, sortOrder int) (*entities.AnnouncementImage, error) {
	_, err := uc.repo.GetByID(ctx, announcementID)
	if err != nil {
		return nil, fmt.Errorf("announcement not found: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	folder := fmt.Sprintf("announcements/%d", announcementID)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	contentType := file.Header.Get("Content-Type")

	path, err := uc.storage.UploadFile(ctx, folder, filename, buf, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	img := entities.AnnouncementImage{
		AnnouncementID: announcementID,
		ImageURL:       path,
		SortOrder:      sortOrder,
	}

	if err := uc.repo.CreateImages(ctx, []entities.AnnouncementImage{img}); err != nil {
		return nil, fmt.Errorf("failed to save image record: %w", err)
	}

	images, err := uc.repo.GetImagesByAnnouncementID(ctx, announcementID)
	if err != nil {
		return nil, err
	}

	for _, i := range images {
		if i.ImageURL == path {
			return &i, nil
		}
	}

	return &img, nil
}
