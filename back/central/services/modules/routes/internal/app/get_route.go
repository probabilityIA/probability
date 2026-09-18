package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

func (uc *UseCase) GetRoute(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
	route, err := uc.repo.GetRouteByID(ctx, businessID, routeID)
	if err != nil {
		return nil, err
	}

	if route.OriginLat == nil || route.OriginLng == nil {
		if wh, err := uc.resolveOriginWarehouse(ctx, route); err == nil && wh != nil && wh.Lat != nil && wh.Lng != nil {
			applyOriginWarehouse(route, wh)
		}
	}

	return route, nil
}
