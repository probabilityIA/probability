package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ReplaceLinks(ctx context.Context, announcementID uint, links []entities.AnnouncementLink) error {
	if err := r.db.Conn(ctx).Where("announcement_id = ?", announcementID).Delete(&models.AnnouncementLink{}).Error; err != nil {
		return err
	}

	if len(links) == 0 {
		return nil
	}

	var linkModels []models.AnnouncementLink
	for _, l := range links {
		linkModels = append(linkModels, models.AnnouncementLink{
			AnnouncementID: announcementID,
			Label:          l.Label,
			URL:            l.URL,
			SortOrder:      l.SortOrder,
		})
	}

	return r.db.Conn(ctx).Create(&linkModels).Error
}
