package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (uc *UseCase) ReorderStops(ctx context.Context, dto dtos.ReorderStopsDTO) error {
	route, err := uc.repo.GetRouteByID(ctx, dto.BusinessID, dto.RouteID)
	if err != nil {
		return err
	}

	if len(dto.StopIDs) != len(route.Stops) {
		return domainerrors.ErrStopIDsMismatch
	}

	return uc.repo.ReorderStops(ctx, dto.RouteID, dto.StopIDs)
}
