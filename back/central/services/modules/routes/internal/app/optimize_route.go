package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
)

func (uc *UseCase) OptimizeRoute(ctx context.Context, dto dtos.OptimizeRouteDTO) (*dtos.OptimizeRouteResult, error) {
	if uc.optimizer == nil || !uc.optimizer.IsConfigured() {
		return nil, domainerrors.ErrOptimizerNotConfigured
	}

	route, err := uc.repo.GetRouteByID(ctx, dto.BusinessID, dto.RouteID)
	if err != nil {
		return nil, err
	}

	if route.Status != "planned" {
		return nil, domainerrors.ErrRouteNotEditable
	}

	if route.OriginLat == nil || route.OriginLng == nil {
		return nil, domainerrors.ErrOriginMissing
	}

	conCoords, sinCoords := splitByCoords(route.Stops)
	if len(conCoords) < 2 {
		return nil, domainerrors.ErrNotEnoughStops
	}

	puntos := make([]dtos.GeoPoint, 0, len(conCoords))
	for _, s := range conCoords {
		puntos = append(puntos, dtos.GeoPoint{Lat: *s.Lat, Lng: *s.Lng})
	}

	optimizada, err := uc.optimizer.Optimize(ctx,
		dtos.GeoPoint{Lat: *route.OriginLat, Lng: *route.OriginLng},
		puntos,
	)
	if err != nil {
		return nil, err
	}

	ordenados := make([]uint, 0, len(route.Stops))
	for _, idx := range optimizada.Order {
		if idx < 0 || idx >= len(conCoords) {
			continue
		}
		ordenados = append(ordenados, conCoords[idx].ID)
	}
	for _, s := range sinCoords {
		ordenados = append(ordenados, s.ID)
	}

	if len(ordenados) != len(route.Stops) {
		return nil, domainerrors.ErrStopIDsMismatch
	}

	if err := uc.repo.ReorderStops(ctx, route.ID, ordenados); err != nil {
		return nil, err
	}

	route.TotalDistanceKm = &optimizada.DistanceKm
	route.TotalDurationMin = &optimizada.DurationMin
	if _, err := uc.repo.UpdateRoute(ctx, route); err != nil {
		return nil, err
	}

	return &dtos.OptimizeRouteResult{
		StopIDs:        ordenados,
		DistanceKm:     optimizada.DistanceKm,
		DurationMin:    optimizada.DurationMin,
		StopsOptimized: len(conCoords),
		StopsSinCoords: len(sinCoords),
	}, nil
}

func splitByCoords(stops []entities.RouteStop) (con, sin []entities.RouteStop) {
	for _, s := range stops {
		if s.Lat != nil && s.Lng != nil {
			con = append(con, s)
			continue
		}
		sin = append(sin, s)
	}
	return con, sin
}
