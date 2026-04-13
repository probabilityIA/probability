package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) RegisterView(ctx context.Context, view *entities.AnnouncementView) error {
	m := models.AnnouncementView{
		AnnouncementID: view.AnnouncementID,
		UserID:         view.UserID,
		BusinessID:     view.BusinessID,
		Action:         string(view.Action),
		LinkID:         view.LinkID,
		ViewedAt:       view.ViewedAt,
	}
	return r.db.Conn(ctx).Create(&m).Error
}

func (r *Repository) GetUserViews(ctx context.Context, userID, announcementID uint) ([]entities.AnnouncementView, error) {
	var rows []models.AnnouncementView
	err := r.db.Conn(ctx).
		Where("user_id = ? AND announcement_id = ?", userID, announcementID).
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.AnnouncementView, len(rows))
	for i, row := range rows {
		result[i] = viewToEntity(&row)
	}
	return result, nil
}

func (r *Repository) DeleteViewsByAnnouncementID(ctx context.Context, announcementID uint) error {
	return r.db.Conn(ctx).
		Where("announcement_id = ?", announcementID).
		Delete(&models.AnnouncementView{}).Error
}
