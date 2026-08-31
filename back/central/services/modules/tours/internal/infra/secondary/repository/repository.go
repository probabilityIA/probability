package repository

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) ListByUser(ctx context.Context, userID, businessID uint) ([]entities.TourProgress, error) {
	var rows []models.UserTourProgress
	err := r.db.Conn(ctx).
		Where("user_id = ? AND business_id = ?", userID, businessID).
		Order("tour_key ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return mappers.ToDomainList(rows), nil
}

func (r *Repository) Upsert(ctx context.Context, progress entities.TourProgress) (*entities.TourProgress, error) {
	row := mappers.ToModel(progress)
	err := r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "business_id"}, {Name: "tour_key"}},
			DoUpdates: clause.AssignmentColumns([]string{"version", "status", "step_index", "completed_at", "updated_at"}),
		}).
		Create(&row).Error
	if err != nil {
		return nil, err
	}

	guardado := mappers.ToDomain(row)
	return &guardado, nil
}

func (r *Repository) Delete(ctx context.Context, userID, businessID uint, tourKey string) error {
	return r.db.Conn(ctx).
		Where("user_id = ? AND business_id = ? AND tour_key = ?", userID, businessID, tourKey).
		Delete(&models.UserTourProgress{}).Error
}

func (r *Repository) DeleteAll(ctx context.Context, userID, businessID uint) error {
	return r.db.Conn(ctx).
		Where("user_id = ? AND business_id = ?", userID, businessID).
		Delete(&models.UserTourProgress{}).Error
}

func (r *Repository) UpsertMany(ctx context.Context, items []entities.TourProgress) error {
	if len(items) == 0 {
		return nil
	}

	filas := make([]models.UserTourProgress, 0, len(items))
	for _, item := range items {
		filas = append(filas, mappers.ToModel(item))
	}

	return r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "business_id"}, {Name: "tour_key"}},
			DoUpdates: clause.AssignmentColumns([]string{"version", "status", "step_index", "completed_at", "updated_at"}),
		}).
		Create(&filas).Error
}
