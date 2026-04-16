package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) CreateImages(ctx context.Context, images []entities.AnnouncementImage) error {
	if len(images) == 0 {
		return nil
	}

	var imgModels []models.AnnouncementImage
	for _, img := range images {
		imgModels = append(imgModels, models.AnnouncementImage{
			AnnouncementID: img.AnnouncementID,
			ImageURL:       img.ImageURL,
			SortOrder:      img.SortOrder,
		})
	}

	return r.db.Conn(ctx).Create(&imgModels).Error
}

func (r *Repository) DeleteImagesByAnnouncementID(ctx context.Context, announcementID uint) error {
	return r.db.Conn(ctx).
		Where("announcement_id = ?", announcementID).
		Delete(&models.AnnouncementImage{}).Error
}

func (r *Repository) GetImageByID(ctx context.Context, id uint) (*entities.AnnouncementImage, error) {
	var row models.AnnouncementImage
	err := r.db.Conn(ctx).First(&row, id).Error
	if err != nil {
		return nil, err
	}
	img := imageToEntity(&row)
	return &img, nil
}

func (r *Repository) DeleteImageByID(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&models.AnnouncementImage{}, id).Error
}

func (r *Repository) GetImagesByAnnouncementID(ctx context.Context, announcementID uint) ([]entities.AnnouncementImage, error) {
	var rows []models.AnnouncementImage
	err := r.db.Conn(ctx).
		Where("announcement_id = ?", announcementID).
		Order("sort_order ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.AnnouncementImage, len(rows))
	for i, row := range rows {
		result[i] = imageToEntity(&row)
	}
	return result, nil
}
