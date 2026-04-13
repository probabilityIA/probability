package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) GetStats(ctx context.Context, announcementID uint) (*entities.AnnouncementStats, error) {
	var stats entities.AnnouncementStats

	conn := r.db.Conn(ctx).Model(&models.AnnouncementView{}).
		Where("announcement_id = ?", announcementID)

	if err := conn.Count(&stats.TotalViews).Error; err != nil {
		return nil, err
	}

	if err := r.db.Conn(ctx).Model(&models.AnnouncementView{}).
		Where("announcement_id = ?", announcementID).
		Distinct("user_id").
		Count(&stats.UniqueUsers).Error; err != nil {
		return nil, err
	}

	if err := r.db.Conn(ctx).Model(&models.AnnouncementView{}).
		Where("announcement_id = ? AND action = ?", announcementID, "clicked_link").
		Count(&stats.TotalClicks).Error; err != nil {
		return nil, err
	}

	if err := r.db.Conn(ctx).Model(&models.AnnouncementView{}).
		Where("announcement_id = ? AND action = ?", announcementID, "accepted").
		Count(&stats.TotalAcceptances).Error; err != nil {
		return nil, err
	}

	if err := r.db.Conn(ctx).Model(&models.AnnouncementView{}).
		Where("announcement_id = ? AND action = ?", announcementID, "closed").
		Count(&stats.TotalClosed).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
