package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) GetActiveAnnouncements(ctx context.Context, params dtos.ActiveAnnouncementsParams) ([]entities.Announcement, error) {
	now := time.Now()

	query := r.db.Conn(ctx).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Links", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Targets").
		Where("status = ?", "active").
		Where("(starts_at IS NULL OR starts_at <= ?)", now).
		Where("(ends_at IS NULL OR ends_at >= ?)", now)

	if params.BusinessID > 0 {
		query = query.Where(
			r.db.Conn(ctx).Where("is_global = ?", true).
				Or("id IN (?)",
					r.db.Conn(ctx).Model(&models.AnnouncementTarget{}).
						Select("announcement_id").
						Where("business_id = ?", params.BusinessID),
				),
		)
	}

	query = query.Order("force_redisplay DESC, created_at DESC")

	var rows []models.Announcement
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]entities.Announcement, len(rows))
	for i, row := range rows {
		result[i] = *announcementToEntity(&row)
	}
	return result, nil
}
