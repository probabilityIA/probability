package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (uc *UseCase) UpdateStopStatus(ctx context.Context, dto dtos.UpdateStopStatusDTO) error {
	route, err := uc.repo.GetRouteByID(ctx, dto.BusinessID, dto.RouteID)
	if err != nil {
		return err
	}

	if route.Status != "in_progress" {
		return domainerrors.ErrRouteNotInProgress
	}

	_, err = uc.repo.GetStopByID(ctx, dto.RouteID, dto.ID)
	if err != nil {
		return err
	}

	if err := uc.repo.UpdateStopStatus(ctx, dto.ID, dto.Status, dto.FailureReason, dto.SignatureURL, dto.PhotoURL); err != nil {
		return err
	}

	_ = uc.repo.UpdateRouteCounters(ctx, dto.RouteID)

	return nil
}
