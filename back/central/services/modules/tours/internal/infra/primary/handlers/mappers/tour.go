package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/infra/primary/handlers/response"
)

func ToResponse(progress entities.TourProgress) response.TourProgress {
	item := response.TourProgress{
		TourKey:   progress.TourKey,
		Version:   progress.Version,
		Status:    progress.Status,
		StepIndex: progress.StepIndex,
	}
	if progress.CompletedAt != nil {
		item.CompletedAt = progress.CompletedAt.Format(time.RFC3339)
	}
	if !progress.UpdatedAt.IsZero() {
		item.UpdatedAt = progress.UpdatedAt.Format(time.RFC3339)
	}
	return item
}

func ToResponseList(items []entities.TourProgress) []response.TourProgress {
	salida := make([]response.TourProgress, 0, len(items))
	for _, item := range items {
		salida = append(salida, ToResponse(item))
	}
	return salida
}
