package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/dtos"
	tourerrors "github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/errors"
)

func (uc *useCase) ListProgress(ctx context.Context, input dtos.ListProgressInput) (*dtos.ListProgressResult, error) {
	if input.UserID == 0 {
		return nil, tourerrors.ErrUserRequired
	}

	items, err := uc.repo.ListByUser(ctx, input.UserID, input.BusinessID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar el progreso de tours: %w", err)
	}

	return &dtos.ListProgressResult{Items: items}, nil
}
