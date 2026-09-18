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

func ptrF(v float64) *float64 { return &v }
func ptrS(v string) *string   { return &v }

func rutaSinOrigen(stops []entities.RouteStop) (*mocks.RepositoryMock, **entities.Route) {
	var guardada *entities.Route
	repo := &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return &entities.Route{ID: routeID, BusinessID: businessID, Status: "planned", Stops: stops}, nil
		},
		UpdateRouteFn: func(ctx context.Context, route *entities.Route) (*entities.Route, error) {
			guardada = route
			return route, nil
		},
	}
	return repo, &guardada
}

func optimizadorQueRegistraOrigen(origen *dtos.GeoPoint) *mocks.OptimizerMock {
	return &mocks.OptimizerMock{
		OptimizeFn: func(ctx context.Context, origin dtos.GeoPoint, stops []dtos.GeoPoint) (dtos.OptimizedRoute, error) {
			*origen = origin
			return dtos.OptimizedRoute{Order: []int{0, 1}, DistanceKm: 10, DurationMin: 20}, nil
		},
	}
}

func paradasConOrden() []entities.RouteStop {
	a := paradaCon(1, 4.71, -74.07)
	a.OrderID = ptrS("o-1")
	b := paradaCon(2, 4.60, -74.08)
	b.OrderID = ptrS("o-2")
	return []entities.RouteStop{a, b}
}

func TestOptimizeRouteSinOrigenUsaBodegaComunDeLasOrdenes(t *testing.T) {
	repo, guardada := rutaSinOrigen(paradasConOrden())
	repo.GetOrdersWarehouseIDsFn = func(ctx context.Context, businessID uint, orderIDs []string) ([]uint, error) {
		assert.ElementsMatch(t, []string{"o-1", "o-2"}, orderIDs)
		return []uint{22}, nil
	}
	repo.GetWarehouseOriginFn = func(ctx context.Context, businessID, warehouseID uint) (*entities.OriginWarehouse, error) {
		require.Equal(t, uint(22), warehouseID)
		return &entities.OriginWarehouse{ID: 22, Address: "Cl. 47a #51-18", City: "COPACABANA", Lat: ptrF(6.34), Lng: ptrF(-75.50)}, nil
	}
	repo.GetDefaultWarehouseOriginFn = func(ctx context.Context, businessID uint) (*entities.OriginWarehouse, error) {
		t.Fatal("no debia consultar la bodega por defecto")
		return nil, nil
	}

	var origen dtos.GeoPoint
	uc := newRoutesUseCaseCon(repo, optimizadorQueRegistraOrigen(&origen))
	_, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 4, BusinessID: 46})

	require.NoError(t, err)
	assert.Equal(t, dtos.GeoPoint{Lat: 6.34, Lng: -75.50}, origen)
	require.NotNil(t, (*guardada).OriginWarehouseID)
	assert.Equal(t, uint(22), *(*guardada).OriginWarehouseID)
	assert.Equal(t, "Cl. 47a #51-18, COPACABANA", (*guardada).OriginAddress)
}

func TestOptimizeRouteOrdenesDeVariasBodegasUsaLaPorDefecto(t *testing.T) {
	repo, _ := rutaSinOrigen(paradasConOrden())
	repo.GetOrdersWarehouseIDsFn = func(ctx context.Context, businessID uint, orderIDs []string) ([]uint, error) {
		return []uint{3, 7}, nil
	}
	repo.GetDefaultWarehouseOriginFn = func(ctx context.Context, businessID uint) (*entities.OriginWarehouse, error) {
		return &entities.OriginWarehouse{ID: 3, Lat: ptrF(4.65), Lng: ptrF(-74.10)}, nil
	}

	var origen dtos.GeoPoint
	uc := newRoutesUseCaseCon(repo, optimizadorQueRegistraOrigen(&origen))
	_, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 4, BusinessID: 26})

	require.NoError(t, err)
	assert.Equal(t, dtos.GeoPoint{Lat: 4.65, Lng: -74.10}, origen)
}

func TestOptimizeRouteBodegaSinUbicacionDevuelveErrorClaro(t *testing.T) {
	repo, _ := rutaSinOrigen(paradasConOrden())
	repo.GetDefaultWarehouseOriginFn = func(ctx context.Context, businessID uint) (*entities.OriginWarehouse, error) {
		return &entities.OriginWarehouse{ID: 3, Name: "principal"}, nil
	}

	uc := newRoutesUseCaseCon(repo, nil)
	_, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 4, BusinessID: 26})

	assert.ErrorIs(t, err, domainerrors.ErrOriginWithoutLocation)
}

func TestOptimizeRouteNegocioSinBodegas(t *testing.T) {
	repo, _ := rutaSinOrigen(paradasConOrden())

	uc := newRoutesUseCaseCon(repo, nil)
	_, err := uc.OptimizeRoute(context.Background(), dtos.OptimizeRouteDTO{RouteID: 4, BusinessID: 26})

	assert.ErrorIs(t, err, domainerrors.ErrOriginMissing)
}

func TestCreateRouteSinOrigenTomaLaBodegaPorDefecto(t *testing.T) {
	var creada *entities.Route
	repo := &mocks.RepositoryMock{
		GetDefaultWarehouseOriginFn: func(ctx context.Context, businessID uint) (*entities.OriginWarehouse, error) {
			return &entities.OriginWarehouse{ID: 22, Address: "Cl. 47a #51-18", City: "COPACABANA", Lat: ptrF(6.34), Lng: ptrF(-75.50)}, nil
		},
		CreateRouteFn: func(ctx context.Context, route *entities.Route, stops []entities.RouteStop) (*entities.Route, error) {
			creada = route
			return route, nil
		},
	}

	uc := newRoutesUseCaseCon(repo, nil)
	_, err := uc.CreateRoute(context.Background(), dtos.CreateRouteDTO{BusinessID: 46})

	require.NoError(t, err)
	require.NotNil(t, creada.OriginWarehouseID)
	assert.Equal(t, uint(22), *creada.OriginWarehouseID)
	require.NotNil(t, creada.OriginLat)
	assert.Nil(t, creada.Stops)
}

func TestGetRouteSinOrigenMuestraLaBodegaSinGuardar(t *testing.T) {
	repo, guardada := rutaSinOrigen(paradasConOrden())
	repo.GetDefaultWarehouseOriginFn = func(ctx context.Context, businessID uint) (*entities.OriginWarehouse, error) {
		return &entities.OriginWarehouse{ID: 22, Address: "Cl. 47a #51-18", Lat: ptrF(6.34), Lng: ptrF(-75.50)}, nil
	}

	uc := newRoutesUseCaseCon(repo, nil)
	route, err := uc.GetRoute(context.Background(), 46, 4)

	require.NoError(t, err)
	require.NotNil(t, route.OriginLat)
	assert.InDelta(t, 6.34, *route.OriginLat, 0.0001)
	assert.Nil(t, *guardada)
}
