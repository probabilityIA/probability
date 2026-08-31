package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func ToDomain(row models.UserTourProgress) entities.TourProgress {
	return entities.TourProgress{
		ID:          row.ID,
		UserID:      row.UserID,
		BusinessID:  row.BusinessID,
		TourKey:     row.TourKey,
		Version:     row.Version,
		Status:      row.Status,
		StepIndex:   row.StepIndex,
		CompletedAt: row.CompletedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func ToDomainList(rows []models.UserTourProgress) []entities.TourProgress {
	items := make([]entities.TourProgress, 0, len(rows))
	for _, row := range rows {
		items = append(items, ToDomain(row))
	}
	return items
}

func ToModel(progress entities.TourProgress) models.UserTourProgress {
	return models.UserTourProgress{
		UserID:      progress.UserID,
		BusinessID:  progress.BusinessID,
		TourKey:     progress.TourKey,
		Version:     progress.Version,
		Status:      progress.Status,
		StepIndex:   progress.StepIndex,
		CompletedAt: progress.CompletedAt,
	}
}
