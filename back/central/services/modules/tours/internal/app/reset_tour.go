package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	tourerrors "github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/errors"
)

func (uc *useCase) ResetTour(ctx context.Context, input dtos.ResetInput) error {
	if input.UserID == 0 {
		return tourerrors.ErrUserRequired
	}
	if input.TourKey == "" {
		return tourerrors.ErrTourKeyRequired
	}
	if err := uc.repo.Delete(ctx, input.UserID, input.BusinessID, input.TourKey); err != nil {
		return fmt.Errorf("error al reiniciar el tour: %w", err)
	}
	return nil
}

func (uc *useCase) ResetAll(ctx context.Context, input dtos.ListProgressInput) error {
	if input.UserID == 0 {
		return tourerrors.ErrUserRequired
	}
	if err := uc.repo.DeleteAll(ctx, input.UserID, input.BusinessID); err != nil {
		return fmt.Errorf("error al reiniciar los tours: %w", err)
	}
	return nil
}
