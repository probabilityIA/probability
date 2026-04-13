package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ReplaceTargets(ctx context.Context, announcementID uint, targets []entities.AnnouncementTarget) error {
	if err := r.db.Conn(ctx).Where("announcement_id = ?", announcementID).Delete(&models.AnnouncementTarget{}).Error; err != nil {
		return err
	}

	if len(targets) == 0 {
		return nil
	}

	var targetModels []models.AnnouncementTarget
	for _, t := range targets {
		targetModels = append(targetModels, models.AnnouncementTarget{
			AnnouncementID: announcementID,
			BusinessID:     t.BusinessID,
		})
	}

	return r.db.Conn(ctx).Create(&targetModels).Error
}
