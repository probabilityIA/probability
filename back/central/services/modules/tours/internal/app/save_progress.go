package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"
	tourerrors "github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/errors"
)

func (uc *useCase) SaveProgress(ctx context.Context, input dtos.SaveProgressInput) (*entities.TourProgress, error) {
	if input.UserID == 0 {
		return nil, tourerrors.ErrUserRequired
	}
	if input.TourKey == "" {
		return nil, tourerrors.ErrTourKeyRequired
	}
	if !entities.IsValidStatus(input.Status) {
		return nil, tourerrors.ErrInvalidStatus
	}
	if input.Version <= 0 {
		return nil, tourerrors.ErrInvalidVersion
	}
	if input.StepIndex < 0 {
		input.StepIndex = 0
	}

	progress := entities.TourProgress{
		UserID:     input.UserID,
		BusinessID: input.BusinessID,
		TourKey:    input.TourKey,
		Version:    input.Version,
		Status:     input.Status,
		StepIndex:  input.StepIndex,
	}
	if input.Status == entities.StatusCompleted || input.Status == entities.StatusSkipped {
		ahora := time.Now()
		progress.CompletedAt = &ahora
	}

	guardado, err := uc.repo.Upsert(ctx, progress)
	if err != nil {
		return nil, fmt.Errorf("error al guardar el progreso del tour: %w", err)
	}

	return guardado, nil
}
