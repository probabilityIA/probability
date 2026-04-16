package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListCategories(ctx context.Context) ([]entities.AnnouncementCategory, error) {
	var rows []models.AnnouncementCategory
	if err := r.db.Conn(ctx).Order("name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]entities.AnnouncementCategory, len(rows))
	for i, row := range rows {
		result[i] = categoryToEntity(&row)
	}
	return result, nil
}
