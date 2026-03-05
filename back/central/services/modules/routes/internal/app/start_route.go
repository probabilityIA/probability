package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (uc *UseCase) StartRoute(ctx context.Context, businessID, routeID uint) error {
	route, err := uc.repo.GetRouteByID(ctx, businessID, routeID)
	if err != nil {
		return err
	}

	if route.Status != "planned" {
		return domainerrors.ErrRouteNotPlanned
	}

	if err := uc.repo.UpdateRouteStatus(ctx, routeID, "in_progress"); err != nil {
		return err
	}

	if err := uc.repo.SetRouteActualStart(ctx, routeID); err != nil {
		return err
	}

	// Set driver to on_route
	if route.DriverID != nil {
		_ = uc.repo.UpdateDriverStatus(ctx, *route.DriverID, "on_route")
	}

	return nil
}
