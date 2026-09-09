package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/mocks"
)

func rutaConParadas(status string, stops []entities.RouteStop) *mocks.RepositoryMock {
	lat, lng := 4.6888, -74.0717
	return &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return &entities.Route{
				ID:         routeID,
				BusinessID: businessID,
				Status:     status,
				OriginLat:  &lat,
				OriginLng:  &lng,
				Stops:      stops,
			}, nil
		},
		UpdateRouteFn: func(ctx context.Context, route *entities.Route) (*entities.Route, error) {
			return route, nil
		},
	}
}

func paradaCon(id uint, lat, lng float64) entities.RouteStop {
	return entities.RouteStop{ID: id, Lat: &lat, Lng: &lng}
}

func TestOptimizeRouteReordenaYGuardaTotales(t *testing.T) {
	repo := rutaConParadas("planned", []entities.RouteStop{
		paradaCon(10, 4.7110, -74.0721),
		paradaCon(20, 4.6097, -74.0817),
		paradaCon(30, 4.6700, -74.0500),
	})

	var guardadas []uint
	repo.ReorderStopsFn = func(ctx context.Context, routeID uint, stopIDs []uint) error {
		guardadas = stopIDs
		return nil
	}

	var rutaGuardada *entities.Route
	repo.UpdateRouteFn = func(ctx context.Context, route *entities.Route) (*entities.Route, error) {
		rutaGuardada = route
		return route, nil
	}

	optimizer := &mocks.OptimizerMock{
		OptimizeFn: func(ctx context.Context, origin dtos.GeoPoint, stops []dtos.GeoPoint) (dtos.OptimizedRoute, error) {
			return dtos.OptimizedRoute{Order: []int{2, 0, 1}, DistanceKm: 29.061, DurationMin: 67}, nil
		},
	}

	uc := newRoutesUseCaseCon(repo, optimizer)
	result, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 1, BusinessID: 26})

	require.NoError(t, err)
	assert.Equal(t, []uint{30, 10, 20}, guardadas)
	assert.Equal(t, []uint{30, 10, 20}, result.StopIDs)
	assert.Equal(t, 3, result.StopsOptimized)
	assert.Equal(t, 0, result.StopsSinCoords)
	require.NotNil(t, rutaGuardada.TotalDistanceKm)
	assert.InDelta(t, 29.061, *rutaGuardada.TotalDistanceKm, 0.001)
	require.NotNil(t, rutaGuardada.TotalDurationMin)
	assert.Equal(t, 67, *rutaGuardada.TotalDurationMin)
}

func TestOptimizeRouteDejaAlFinalLasParadasSinCoordenadas(t *testing.T) {
	repo := rutaConParadas("planned", []entities.RouteStop{
		paradaCon(10, 4.7110, -74.0721),
		{ID: 99},
		paradaCon(20, 4.6097, -74.0817),
	})

	var guardadas []uint
	repo.ReorderStopsFn = func(ctx context.Context, routeID uint, stopIDs []uint) error {
		guardadas = stopIDs
		return nil
	}

	optimizer := &mocks.OptimizerMock{
		OptimizeFn: func(ctx context.Context, origin dtos.GeoPoint, stops []dtos.GeoPoint) (dtos.OptimizedRoute, error) {
			assert.Len(t, stops, 2)
			return dtos.OptimizedRoute{Order: []int{1, 0}, DistanceKm: 12, DurationMin: 30}, nil
		},
	}

	uc := newRoutesUseCaseCon(repo, optimizer)
	result, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 1, BusinessID: 26})

	require.NoError(t, err)
	assert.Equal(t, []uint{20, 10, 99}, guardadas)
	assert.Equal(t, 2, result.StopsOptimized)
	assert.Equal(t, 1, result.StopsSinCoords)
}

func TestOptimizeRouteRechazaRutaYaIniciada(t *testing.T) {
	repo := rutaConParadas("in_progress", []entities.RouteStop{
		paradaCon(10, 4.7110, -74.0721),
		paradaCon(20, 4.6097, -74.0817),
	})

	llamado := false
	repo.ReorderStopsFn = func(ctx context.Context, routeID uint, stopIDs []uint) error {
		llamado = true
		return nil
	}

	uc := newRoutesUseCaseCon(repo, nil)
	_, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 1, BusinessID: 26})

	assert.ErrorIs(t, err, domainerrors.ErrRouteNotEditable)
	assert.False(t, llamado, "no debe tocar el orden de una ruta en curso")
}

func TestOptimizeRouteExigeDosParadasConCoordenadas(t *testing.T) {
	repo := rutaConParadas("planned", []entities.RouteStop{
		paradaCon(10, 4.7110, -74.0721),
		{ID: 99},
	})

	uc := newRoutesUseCaseCon(repo, nil)
	_, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 1, BusinessID: 26})

	assert.ErrorIs(t, err, domainerrors.ErrNotEnoughStops)
}

func TestOptimizeRouteSinOptimizadorConfigurado(t *testing.T) {
	repo := rutaConParadas("planned", []entities.RouteStop{
		paradaCon(10, 4.7110, -74.0721),
		paradaCon(20, 4.6097, -74.0817),
	})

	optimizer := &mocks.OptimizerMock{IsConfiguredFn: func() bool { return false }}

	uc := newRoutesUseCaseCon(repo, optimizer)
	_, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 1, BusinessID: 26})

	assert.ErrorIs(t, err, domainerrors.ErrOptimizerNotConfigured)
}
