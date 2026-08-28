package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
	tourerrors "github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/errors"
)

func (uc *useCase) SkipAll(ctx context.Context, input dtos.SkipAllInput) error {
	if input.UserID == 0 {
		return tourerrors.ErrUserRequired
	}
	if len(input.Tours) == 0 {
		return tourerrors.ErrTourKeyRequired
	}

	ahora := time.Now()
	registros := make([]entities.TourProgress, 0, len(input.Tours))

	for _, tour := range input.Tours {
		if tour.TourKey == "" {
			return tourerrors.ErrTourKeyRequired
		}
		if tour.Version <= 0 {
			return tourerrors.ErrInvalidVersion
		}
		registros = append(registros, entities.TourProgress{
			UserID:      input.UserID,
			BusinessID:  input.BusinessID,
			TourKey:     tour.TourKey,
			Version:     tour.Version,
			Status:      entities.StatusSkipped,
			StepIndex:   0,
			CompletedAt: &ahora,
		})
	}

	if err := uc.repo.UpsertMany(ctx, registros); err != nil {
		return fmt.Errorf("error al omitir los tutoriales: %w", err)
	}

	uc.log.Info(ctx).Uint("user_id", input.UserID).Int("tours", len(registros)).Msg("Tutoriales omitidos por el usuario")
	return nil
}
