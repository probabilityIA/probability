package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRoutesUseCase(repo *mocks.RepositoryMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo)
}

func uintPtr(v uint) *uint      { return &v }
func strPtr(v string) *string   { return &v }
func f64Ptr(v float64) *float64 { return &v }

func repoConRuta(status string, driverID *uint, stops []entities.RouteStop) *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return &entities.Route{
				ID: routeID, BusinessID: businessID, Status: status,
				DriverID: driverID, DriverName: "Conductor Demo",
				TotalStops: len(stops), Stops: stops,
			}, nil
		},
	}
}

func TestCreateRoute_NaceEnPlaneada(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26, Date: time.Now(),
	})

	require.NoError(t, err)
	require.NotNil(t, repo.CreatedRoute)
	assert.Equal(t, "planned", repo.CreatedRoute.Status)
}

func TestCreateRoute_NumeraLasParadasDesdeUno(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26,
		Stops: []dtos.CreateRouteStopDTO{
			{Address: "Cra 1", City: "Bogota"},
			{Address: "Cra 2", City: "Bogota"},
			{Address: "Cra 3", City: "Bogota"},
		},
	})

	require.NoError(t, err)
	require.Len(t, repo.CreatedStops, 3)
	assert.Equal(t, 1, repo.CreatedStops[0].Sequence, "la secuencia arranca en 1, no en 0")
	assert.Equal(t, 2, repo.CreatedStops[1].Sequence)
	assert.Equal(t, 3, repo.CreatedStops[2].Sequence)
	assert.Equal(t, 3, repo.CreatedRoute.TotalStops)
}

func TestCreateRoute_TodasLasParadasNacenPendientes(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26,
		Stops: []dtos.CreateRouteStopDTO{
			{Address: "Cra 1"}, {Address: "Cra 2"},
		},
	})

	require.NoError(t, err)
	for i, s := range repo.CreatedStops {
		assert.Equal(t, "pending", s.Status, "parada %d", i)
	}
}

func TestCreateRoute_SinConductorNiVehiculo_NoLosConsulta(t *testing.T) {
	consultadoConductor := false
	consultadoVehiculo := false
	repo := &mocks.RepositoryMock{
		GetDriverNameByIDFn: func(ctx context.Context, driverID uint) (string, error) {
			consultadoConductor = true
			return "", nil
		},
		GetVehiclePlateByIDFn: func(ctx context.Context, vehicleID uint) (string, error) {
			consultadoVehiculo = true
			return "", nil
		},
	}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{BusinessID: 26})

	require.NoError(t, err)
	assert.False(t, consultadoConductor)
	assert.False(t, consultadoVehiculo)
	assert.Empty(t, repo.CreatedRoute.DriverName)
	assert.Empty(t, repo.CreatedRoute.VehiclePlate)
}

func TestCreateRoute_DesnormalizaNombreDelConductorYPlaca(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetDriverNameByIDFn: func(ctx context.Context, driverID uint) (string, error) {
			return "Juan Perez", nil
		},
		GetVehiclePlateByIDFn: func(ctx context.Context, vehicleID uint) (string, error) {
			return "XYZ789", nil
		},
	}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26, DriverID: uintPtr(7), VehicleID: uintPtr(3),
	})

	require.NoError(t, err)
	assert.Equal(t, "Juan Perez", repo.CreatedRoute.DriverName)
	assert.Equal(t, "XYZ789", repo.CreatedRoute.VehiclePlate)
}

func TestCreateRoute_ConductorInexistente_NoCrea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetDriverNameByIDFn: func(ctx context.Context, driverID uint) (string, error) {
			return "", domainerrors.ErrDriverNotFound
		},
	}

	got, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26, DriverID: uintPtr(404),
	})

	assert.ErrorIs(t, err, domainerrors.ErrDriverNotFound)
	assert.Nil(t, got)
	assert.Nil(t, repo.CreatedRoute)
}

func TestCreateRoute_VehiculoInexistente_NoCrea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetVehiclePlateByIDFn: func(ctx context.Context, vehicleID uint) (string, error) {
			return "", domainerrors.ErrVehicleNotFound
		},
	}

	got, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26, VehicleID: uintPtr(404),
	})

	assert.ErrorIs(t, err, domainerrors.ErrVehicleNotFound)
	assert.Nil(t, got)
	assert.Nil(t, repo.CreatedRoute)
}

func TestCreateRoute_ConConductor_MarcaLasOrdenesComoUltimaMilla(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetDriverNameByIDFn: func(ctx context.Context, driverID uint) (string, error) {
			return "Juan Perez", nil
		},
	}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26, DriverID: uintPtr(7),
		Stops: []dtos.CreateRouteStopDTO{
			{OrderID: strPtr("ORD-1")},
			{OrderID: nil},
			{OrderID: strPtr("ORD-2")},
		},
	})

	require.NoError(t, err)
	require.Len(t, repo.OrderDriverCalls, 2, "solo las paradas con orden asociada")
	assert.Equal(t, "ORD-1", repo.OrderDriverCalls[0].OrderID)
	assert.Equal(t, "Juan Perez", repo.OrderDriverCalls[0].DriverName)
	assert.True(t, repo.OrderDriverCalls[0].IsLastMile)
	assert.Equal(t, "ORD-2", repo.OrderDriverCalls[1].OrderID)
}

func TestCreateRoute_SinConductor_NoTocaLasOrdenes(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26,
		Stops:      []dtos.CreateRouteStopDTO{{OrderID: strPtr("ORD-1")}},
	})

	require.NoError(t, err)
	assert.Empty(t, repo.OrderDriverCalls,
		"sin conductor asignado no hay a quien atribuirle la orden")
}

func TestCreateRoute_ErrorAlPersistir_NoTocaLasOrdenes(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreateRouteFn: func(ctx context.Context, r *entities.Route, s []entities.RouteStop) (*entities.Route, error) {
			return nil, dbErr
		},
	}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26, DriverID: uintPtr(7),
		Stops: []dtos.CreateRouteStopDTO{{OrderID: strPtr("ORD-1")}},
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.OrderDriverCalls)
}

func TestCreateRoute_CopiaLosDatosDeCadaParada(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newRoutesUseCase(repo).CreateRoute(context.Background(), dtos.CreateRouteDTO{
		BusinessID: 26,
		Stops: []dtos.CreateRouteStopDTO{{
			OrderID: strPtr("ORD-1"), Address: "Cra 1 #2-3", City: "Bogota",
			Lat: f64Ptr(4.65), Lng: f64Ptr(-74.05),
			CustomerName: "Ana", CustomerPhone: "3001234567",
			DeliveryNotes: strPtr("timbre 2"),
		}},
	})

	require.NoError(t, err)
	require.Len(t, repo.CreatedStops, 1)
	s := repo.CreatedStops[0]
	assert.Equal(t, "Cra 1 #2-3", s.Address)
	assert.Equal(t, "Bogota", s.City)
	assert.InDelta(t, 4.65, *s.Lat, 1e-9)
	assert.Equal(t, "Ana", s.CustomerName)
	assert.Equal(t, "3001234567", s.CustomerPhone)
	assert.Equal(t, "timbre 2", *s.DeliveryNotes)
}

func TestStartRoute_SoloDesdePlaneada(t *testing.T) {
	for _, estado := range []string{"in_progress", "completed", "cancelled"} {
		t.Run(estado, func(t *testing.T) {
			repo := repoConRuta(estado, uintPtr(7), nil)

			err := newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1)

			assert.ErrorIs(t, err, domainerrors.ErrRouteNotPlanned)
			assert.Empty(t, repo.RouteStatusCalls)
			assert.Zero(t, repo.ActualStartCalls)
		})
	}
}

func TestStartRoute_PasaAEnProgresoYMarcaLaHora(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)

	err := newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Equal(t, []string{"in_progress"}, repo.RouteStatusCalls)
	assert.Equal(t, 1, repo.ActualStartCalls)
}

func TestStartRoute_PoneAlConductorEnRuta(t *testing.T) {
	repo := repoConRuta("planned", uintPtr(7), nil)

	err := newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	require.Len(t, repo.DriverStatusCalls, 1)
	assert.Equal(t, uint(7), repo.DriverStatusCalls[0].DriverID)
	assert.Equal(t, "on_route", repo.DriverStatusCalls[0].Status)
}

func TestStartRoute_SinConductor_NoCambiaEstadoDeNadie(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)

	err := newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Empty(t, repo.DriverStatusCalls)
}

func TestStartRoute_FalloAlMarcarElConductor_NoRompeElArranque(t *testing.T) {
	repo := repoConRuta("planned", uintPtr(7), nil)
	repo.UpdateDriverStatusFn = func(ctx context.Context, driverID uint, status string) error {
		return stderrors.New("drivers caido")
	}

	err := newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1)

	require.NoError(t, err, "el estado del conductor es secundario, la ruta ya arranco")
}

func TestStartRoute_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al leer la ruta", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
				return nil, dbErr
			},
		}
		assert.ErrorIs(t, newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1), dbErr)
	})

	t.Run("al cambiar el estado", func(t *testing.T) {
		repo := repoConRuta("planned", nil, nil)
		repo.UpdateRouteStatusFn = func(ctx context.Context, routeID uint, status string) error { return dbErr }
		assert.ErrorIs(t, newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1), dbErr)
	})

	t.Run("al marcar la hora de inicio", func(t *testing.T) {
		repo := repoConRuta("planned", nil, nil)
		repo.SetRouteActualStartFn = func(ctx context.Context, routeID uint) error { return dbErr }
		assert.ErrorIs(t, newRoutesUseCase(repo).StartRoute(context.Background(), 26, 1), dbErr)
	})
}

func TestCompleteRoute_SoloDesdeEnProgreso(t *testing.T) {
	for _, estado := range []string{"planned", "completed", "cancelled"} {
		t.Run(estado, func(t *testing.T) {
			repo := repoConRuta(estado, nil, nil)

			err := newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1)

			assert.ErrorIs(t, err, domainerrors.ErrRouteNotInProgress)
			assert.Empty(t, repo.RouteStatusCalls)
			assert.Empty(t, repo.PendingSetTo)
		})
	}
}

func TestCompleteRoute_MarcaLasParadasPendientesComoOmitidas(t *testing.T) {
	repo := repoConRuta("in_progress", nil, nil)

	err := newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Equal(t, []string{"skipped"}, repo.PendingSetTo,
		"al cerrar la ruta las paradas que quedaron pendientes se marcan omitidas, no entregadas")
}

func TestCompleteRoute_CierraLaRutaYRefrescaContadores(t *testing.T) {
	repo := repoConRuta("in_progress", nil, nil)

	err := newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Equal(t, []string{"completed"}, repo.RouteStatusCalls)
	assert.Equal(t, 1, repo.ActualEndCalls)
	assert.Equal(t, 1, repo.CountersRefreshed)
}

func TestCompleteRoute_DevuelveElConductorAActivo(t *testing.T) {
	repo := repoConRuta("in_progress", uintPtr(7), nil)

	err := newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	require.Len(t, repo.DriverStatusCalls, 1)
	assert.Equal(t, "active", repo.DriverStatusCalls[0].Status,
		"si no vuelve a activo el conductor queda bloqueado para la siguiente ruta")
}

func TestCompleteRoute_FalloAlOmitirPendientes_NoAbortaElCierre(t *testing.T) {
	repo := repoConRuta("in_progress", nil, nil)
	repo.SetPendingStopsStatusFn = func(ctx context.Context, routeID uint, status string) error {
		return stderrors.New("db caida")
	}

	err := newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Equal(t, []string{"completed"}, repo.RouteStatusCalls)
}

func TestCompleteRoute_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("al cambiar el estado", func(t *testing.T) {
		repo := repoConRuta("in_progress", nil, nil)
		repo.UpdateRouteStatusFn = func(ctx context.Context, routeID uint, status string) error { return dbErr }
		assert.ErrorIs(t, newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1), dbErr)
	})

	t.Run("al marcar la hora de fin", func(t *testing.T) {
		repo := repoConRuta("in_progress", nil, nil)
		repo.SetRouteActualEndFn = func(ctx context.Context, routeID uint) error { return dbErr }
		assert.ErrorIs(t, newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1), dbErr)
	})

	t.Run("al refrescar contadores", func(t *testing.T) {
		repo := repoConRuta("in_progress", nil, nil)
		repo.UpdateRouteCountersFn = func(ctx context.Context, routeID uint) error { return dbErr }
		assert.ErrorIs(t, newRoutesUseCase(repo).CompleteRoute(context.Background(), 26, 1), dbErr)
	})
}

func TestUpdateStopStatus_SoloConLaRutaEnProgreso(t *testing.T) {
	for _, estado := range []string{"planned", "completed"} {
		t.Run(estado, func(t *testing.T) {
			repo := repoConRuta(estado, nil, nil)

			err := newRoutesUseCase(repo).UpdateStopStatus(context.Background(),
				dtos.UpdateStopStatusDTO{ID: 5, RouteID: 1, BusinessID: 26, Status: "delivered"})

			assert.ErrorIs(t, err, domainerrors.ErrRouteNotInProgress,
				"no se puede entregar una parada de una ruta que no arranco")
			assert.Empty(t, repo.StopStatusCalls)
		})
	}
}

func TestUpdateStopStatus_ParadaInexistente_NoActualiza(t *testing.T) {
	repo := repoConRuta("in_progress", nil, nil)
	repo.GetStopByIDFn = func(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error) {
		return nil, domainerrors.ErrStopNotFound
	}

	err := newRoutesUseCase(repo).UpdateStopStatus(context.Background(),
		dtos.UpdateStopStatusDTO{ID: 404, RouteID: 1, BusinessID: 26, Status: "delivered"})

	assert.ErrorIs(t, err, domainerrors.ErrStopNotFound)
	assert.Empty(t, repo.StopStatusCalls)
}

func TestUpdateStopStatus_PasaEvidenciaYMotivoDeFallo(t *testing.T) {
	repo := repoConRuta("in_progress", nil, nil)
	motivo := "cliente ausente"

	err := newRoutesUseCase(repo).UpdateStopStatus(context.Background(), dtos.UpdateStopStatusDTO{
		ID: 5, RouteID: 1, BusinessID: 26, Status: "failed",
		FailureReason: &motivo, SignatureURL: "s3://firma.png", PhotoURL: "s3://foto.jpg",
	})

	require.NoError(t, err)
	require.Len(t, repo.StopStatusCalls, 1)
	c := repo.StopStatusCalls[0]
	assert.Equal(t, uint(5), c.StopID)
	assert.Equal(t, "failed", c.Status)
	require.NotNil(t, c.FailureReason)
	assert.Equal(t, motivo, *c.FailureReason)
	assert.Equal(t, "s3://firma.png", c.SignatureURL)
	assert.Equal(t, "s3://foto.jpg", c.PhotoURL)
}

func TestUpdateStopStatus_RefrescaLosContadores(t *testing.T) {
	repo := repoConRuta("in_progress", nil, nil)

	err := newRoutesUseCase(repo).UpdateStopStatus(context.Background(),
		dtos.UpdateStopStatusDTO{ID: 5, RouteID: 1, BusinessID: 26, Status: "delivered"})

	require.NoError(t, err)
	assert.Equal(t, 1, repo.CountersRefreshed)
}

func TestUpdateStopStatus_ErrorAlGuardar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := repoConRuta("in_progress", nil, nil)
	repo.UpdateStopStatusFn = func(ctx context.Context, stopID uint, status string, fr *string, s, p string) error {
		return dbErr
	}

	err := newRoutesUseCase(repo).UpdateStopStatus(context.Background(),
		dtos.UpdateStopStatusDTO{ID: 5, RouteID: 1, BusinessID: 26, Status: "delivered"})

	assert.ErrorIs(t, err, dbErr)
	assert.Zero(t, repo.CountersRefreshed)
}

func TestAddStop_LaSecuenciaContinuaDesdeElTotal(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return &entities.Route{ID: routeID, TotalStops: 4}, nil
		},
	}

	_, err := newRoutesUseCase(repo).AddStop(context.Background(),
		dtos.AddStopDTO{RouteID: 1, BusinessID: 26, Address: "Cra 5"})

	require.NoError(t, err)
	require.Len(t, repo.AddedStops, 1)
	assert.Equal(t, 5, repo.AddedStops[0].Sequence, "la nueva parada va al final")
	assert.Equal(t, "pending", repo.AddedStops[0].Status)
}

func TestAddStop_ConOrdenYConductor_MarcaLaOrden(t *testing.T) {
	repo := repoConRuta("planned", uintPtr(7), nil)

	_, err := newRoutesUseCase(repo).AddStop(context.Background(),
		dtos.AddStopDTO{RouteID: 1, BusinessID: 26, OrderID: strPtr("ORD-9")})

	require.NoError(t, err)
	require.Len(t, repo.OrderDriverCalls, 1)
	assert.Equal(t, "ORD-9", repo.OrderDriverCalls[0].OrderID)
	assert.Equal(t, "Conductor Demo", repo.OrderDriverCalls[0].DriverName)
	assert.True(t, repo.OrderDriverCalls[0].IsLastMile)
}

func TestAddStop_SinOrdenOSinConductor_NoMarcaNada(t *testing.T) {
	t.Run("sin orden", func(t *testing.T) {
		repo := repoConRuta("planned", uintPtr(7), nil)
		_, err := newRoutesUseCase(repo).AddStop(context.Background(),
			dtos.AddStopDTO{RouteID: 1, BusinessID: 26})
		require.NoError(t, err)
		assert.Empty(t, repo.OrderDriverCalls)
	})

	t.Run("sin conductor", func(t *testing.T) {
		repo := repoConRuta("planned", nil, nil)
		_, err := newRoutesUseCase(repo).AddStop(context.Background(),
			dtos.AddStopDTO{RouteID: 1, BusinessID: 26, OrderID: strPtr("ORD-9")})
		require.NoError(t, err)
		assert.Empty(t, repo.OrderDriverCalls)
	})
}

func TestAddStop_RefrescaContadores(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)

	_, err := newRoutesUseCase(repo).AddStop(context.Background(),
		dtos.AddStopDTO{RouteID: 1, BusinessID: 26})

	require.NoError(t, err)
	assert.Equal(t, 1, repo.CountersRefreshed)
}

func TestAddStop_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("ruta inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
				return nil, domainerrors.ErrRouteNotFound
			},
		}
		_, err := newRoutesUseCase(repo).AddStop(context.Background(), dtos.AddStopDTO{RouteID: 404})
		assert.ErrorIs(t, err, domainerrors.ErrRouteNotFound)
		assert.Empty(t, repo.AddedStops)
	})

	t.Run("al agregar", func(t *testing.T) {
		repo := repoConRuta("planned", nil, nil)
		repo.AddStopFn = func(ctx context.Context, s *entities.RouteStop) (*entities.RouteStop, error) {
			return nil, dbErr
		}
		_, err := newRoutesUseCase(repo).AddStop(context.Background(), dtos.AddStopDTO{RouteID: 1})
		assert.ErrorIs(t, err, dbErr)
		assert.Zero(t, repo.CountersRefreshed)
	})
}

func TestUpdateStop_ActualizaLosDatosDeEntrega(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)
	repo.GetStopByIDFn = func(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error) {
		return &entities.RouteStop{ID: stopID, RouteID: routeID, Sequence: 3, Status: "pending"}, nil
	}

	_, err := newRoutesUseCase(repo).UpdateStop(context.Background(), dtos.UpdateStopDTO{
		ID: 5, RouteID: 1, BusinessID: 26,
		Address: "Cra 9", City: "Cali", Lat: f64Ptr(3.4), Lng: f64Ptr(-76.5),
		CustomerName: "Pedro", CustomerPhone: "3009876543", DeliveryNotes: strPtr("apto 501"),
	})

	require.NoError(t, err)
	require.Len(t, repo.UpdatedStops, 1)
	s := repo.UpdatedStops[0]
	assert.Equal(t, "Cra 9", s.Address)
	assert.Equal(t, "Cali", s.City)
	assert.Equal(t, "Pedro", s.CustomerName)
	assert.Equal(t, "apto 501", *s.DeliveryNotes)
	assert.Equal(t, 3, s.Sequence, "editar los datos no reordena la parada")
	assert.Equal(t, "pending", s.Status, "ni cambia su estado de entrega")
}

func TestUpdateStop_ErroresDeLectura_NoActualizan(t *testing.T) {
	t.Run("ruta inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
				return nil, domainerrors.ErrRouteNotFound
			},
		}
		_, err := newRoutesUseCase(repo).UpdateStop(context.Background(), dtos.UpdateStopDTO{RouteID: 404})
		assert.ErrorIs(t, err, domainerrors.ErrRouteNotFound)
		assert.Empty(t, repo.UpdatedStops)
	})

	t.Run("parada inexistente", func(t *testing.T) {
		repo := repoConRuta("planned", nil, nil)
		repo.GetStopByIDFn = func(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error) {
			return nil, domainerrors.ErrStopNotFound
		}
		_, err := newRoutesUseCase(repo).UpdateStop(context.Background(), dtos.UpdateStopDTO{ID: 404, RouteID: 1})
		assert.ErrorIs(t, err, domainerrors.ErrStopNotFound)
		assert.Empty(t, repo.UpdatedStops)
	})
}

func TestDeleteStop_LiberaLaOrdenAntesDeBorrar(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)
	repo.GetStopByIDFn = func(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error) {
		return &entities.RouteStop{ID: stopID, OrderID: strPtr("ORD-3")}, nil
	}

	err := newRoutesUseCase(repo).DeleteStop(context.Background(), 26, 1, 5)

	require.NoError(t, err)
	assert.Equal(t, []string{"ORD-3"}, repo.ClearedOrders,
		"la orden debe quedar libre para reasignarse a otra ruta")
	assert.Equal(t, 1, repo.CountersRefreshed)
}

func TestDeleteStop_ParadaSinOrden_NoLiberaNada(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)

	err := newRoutesUseCase(repo).DeleteStop(context.Background(), 26, 1, 5)

	require.NoError(t, err)
	assert.Empty(t, repo.ClearedOrders)
}

func TestDeleteStop_ErrorAlBorrar_NoRefrescaContadores(t *testing.T) {
	dbErr := stderrors.New("db caida")
	repo := repoConRuta("planned", nil, nil)
	repo.DeleteStopFn = func(ctx context.Context, routeID, stopID uint) error { return dbErr }

	err := newRoutesUseCase(repo).DeleteStop(context.Background(), 26, 1, 5)

	assert.ErrorIs(t, err, dbErr)
	assert.Zero(t, repo.CountersRefreshed)
}

func TestDeleteStop_ErroresDeLectura_NoBorran(t *testing.T) {
	t.Run("ruta inexistente", func(t *testing.T) {
		borrado := false
		repo := &mocks.RepositoryMock{
			GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
				return nil, domainerrors.ErrRouteNotFound
			},
			DeleteStopFn: func(ctx context.Context, routeID, stopID uint) error {
				borrado = true
				return nil
			},
		}
		assert.ErrorIs(t, newRoutesUseCase(repo).DeleteStop(context.Background(), 26, 404, 5),
			domainerrors.ErrRouteNotFound)
		assert.False(t, borrado)
	})

	t.Run("parada inexistente", func(t *testing.T) {
		borrado := false
		repo := repoConRuta("planned", nil, nil)
		repo.GetStopByIDFn = func(ctx context.Context, routeID, stopID uint) (*entities.RouteStop, error) {
			return nil, domainerrors.ErrStopNotFound
		}
		repo.DeleteStopFn = func(ctx context.Context, routeID, stopID uint) error {
			borrado = true
			return nil
		}
		assert.ErrorIs(t, newRoutesUseCase(repo).DeleteStop(context.Background(), 26, 1, 404),
			domainerrors.ErrStopNotFound)
		assert.False(t, borrado)
	})
}

func TestReorderStops_CantidadDistinta_Rechaza(t *testing.T) {
	casos := []struct {
		nombre   string
		enRuta   int
		enviados []uint
	}{
		{"faltan ids", 3, []uint{1, 2}},
		{"sobran ids", 2, []uint{1, 2, 3}},
		{"lista vacia con paradas", 2, nil},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			stops := make([]entities.RouteStop, tc.enRuta)
			repo := repoConRuta("planned", nil, stops)

			err := newRoutesUseCase(repo).ReorderStops(context.Background(),
				dtos.ReorderStopsDTO{RouteID: 1, BusinessID: 26, StopIDs: tc.enviados})

			assert.ErrorIs(t, err, domainerrors.ErrStopIDsMismatch,
				"un reorden parcial dejaria secuencias duplicadas o huecos")
			assert.Nil(t, repo.ReorderedIDs)
		})
	}
}

func TestReorderStops_MismaCantidad_Aplica(t *testing.T) {
	stops := make([]entities.RouteStop, 3)
	repo := repoConRuta("planned", nil, stops)

	err := newRoutesUseCase(repo).ReorderStops(context.Background(),
		dtos.ReorderStopsDTO{RouteID: 1, BusinessID: 26, StopIDs: []uint{3, 1, 2}})

	require.NoError(t, err)
	assert.Equal(t, []uint{3, 1, 2}, repo.ReorderedIDs)
}

func TestReorderStops_RutaSinParadasYListaVacia_Aplica(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)

	err := newRoutesUseCase(repo).ReorderStops(context.Background(),
		dtos.ReorderStopsDTO{RouteID: 1, BusinessID: 26, StopIDs: nil})

	require.NoError(t, err)
}

func TestReorderStops_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("ruta inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
				return nil, domainerrors.ErrRouteNotFound
			},
		}
		assert.ErrorIs(t, newRoutesUseCase(repo).ReorderStops(context.Background(),
			dtos.ReorderStopsDTO{RouteID: 404}), domainerrors.ErrRouteNotFound)
	})

	t.Run("al reordenar", func(t *testing.T) {
		repo := repoConRuta("planned", nil, nil)
		repo.ReorderStopsFn = func(ctx context.Context, routeID uint, ids []uint) error { return dbErr }
		assert.ErrorIs(t, newRoutesUseCase(repo).ReorderStops(context.Background(),
			dtos.ReorderStopsDTO{RouteID: 1}), dbErr)
	})
}

func TestUpdateRoute_QuitarConductorLimpiaElNombreDesnormalizado(t *testing.T) {
	var guardada *entities.Route
	repo := &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return &entities.Route{
				ID: routeID, DriverID: uintPtr(7), DriverName: "Juan Perez",
				VehicleID: uintPtr(3), VehiclePlate: "XYZ789",
			}, nil
		},
		UpdateRouteFn: func(ctx context.Context, r *entities.Route) (*entities.Route, error) {
			guardada = r
			return r, nil
		},
	}

	_, err := newRoutesUseCase(repo).UpdateRoute(context.Background(),
		dtos.UpdateRouteDTO{ID: 1, BusinessID: 26, DriverID: nil, VehicleID: nil})

	require.NoError(t, err)
	require.NotNil(t, guardada)
	assert.Nil(t, guardada.DriverID)
	assert.Empty(t, guardada.DriverName,
		"si se quita el conductor el nombre desnormalizado no puede quedar pegado")
	assert.Nil(t, guardada.VehicleID)
	assert.Empty(t, guardada.VehiclePlate)
}

func TestUpdateRoute_CambiarConductorRefrescaElNombre(t *testing.T) {
	var guardada *entities.Route
	repo := &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return &entities.Route{ID: routeID, DriverID: uintPtr(7), DriverName: "Viejo"}, nil
		},
		GetDriverNameByIDFn: func(ctx context.Context, driverID uint) (string, error) {
			return "Nuevo Conductor", nil
		},
		GetVehiclePlateByIDFn: func(ctx context.Context, vehicleID uint) (string, error) {
			return "NEW123", nil
		},
		UpdateRouteFn: func(ctx context.Context, r *entities.Route) (*entities.Route, error) {
			guardada = r
			return r, nil
		},
	}

	_, err := newRoutesUseCase(repo).UpdateRoute(context.Background(),
		dtos.UpdateRouteDTO{ID: 1, BusinessID: 26, DriverID: uintPtr(9), VehicleID: uintPtr(4)})

	require.NoError(t, err)
	assert.Equal(t, "Nuevo Conductor", guardada.DriverName)
	assert.Equal(t, "NEW123", guardada.VehiclePlate)
	assert.Equal(t, uint(9), *guardada.DriverID)
}

func TestUpdateRoute_ConductorOVehiculoInexistente_NoActualiza(t *testing.T) {
	t.Run("conductor", func(t *testing.T) {
		actualizada := false
		repo := &mocks.RepositoryMock{
			GetDriverNameByIDFn: func(ctx context.Context, driverID uint) (string, error) {
				return "", domainerrors.ErrDriverNotFound
			},
			UpdateRouteFn: func(ctx context.Context, r *entities.Route) (*entities.Route, error) {
				actualizada = true
				return r, nil
			},
		}
		_, err := newRoutesUseCase(repo).UpdateRoute(context.Background(),
			dtos.UpdateRouteDTO{ID: 1, DriverID: uintPtr(404)})
		assert.ErrorIs(t, err, domainerrors.ErrDriverNotFound)
		assert.False(t, actualizada)
	})

	t.Run("vehiculo", func(t *testing.T) {
		actualizada := false
		repo := &mocks.RepositoryMock{
			GetVehiclePlateByIDFn: func(ctx context.Context, vehicleID uint) (string, error) {
				return "", domainerrors.ErrVehicleNotFound
			},
			UpdateRouteFn: func(ctx context.Context, r *entities.Route) (*entities.Route, error) {
				actualizada = true
				return r, nil
			},
		}
		_, err := newRoutesUseCase(repo).UpdateRoute(context.Background(),
			dtos.UpdateRouteDTO{ID: 1, VehicleID: uintPtr(404)})
		assert.ErrorIs(t, err, domainerrors.ErrVehicleNotFound)
		assert.False(t, actualizada)
	})
}

func TestUpdateRoute_ConservaElEstadoYLosContadores(t *testing.T) {
	var guardada *entities.Route
	repo := &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return &entities.Route{
				ID: routeID, Status: "in_progress",
				TotalStops: 5, CompletedStops: 3, FailedStops: 1,
			}, nil
		},
		UpdateRouteFn: func(ctx context.Context, r *entities.Route) (*entities.Route, error) {
			guardada = r
			return r, nil
		},
	}

	_, err := newRoutesUseCase(repo).UpdateRoute(context.Background(),
		dtos.UpdateRouteDTO{ID: 1, BusinessID: 26, OriginAddress: "Bodega Norte"})

	require.NoError(t, err)
	assert.Equal(t, "in_progress", guardada.Status, "editar datos no cambia el estado de la ruta")
	assert.Equal(t, 5, guardada.TotalStops)
	assert.Equal(t, 3, guardada.CompletedStops)
	assert.Equal(t, "Bodega Norte", guardada.OriginAddress)
}

func TestUpdateRoute_RutaInexistente_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			return nil, domainerrors.ErrRouteNotFound
		},
	}

	_, err := newRoutesUseCase(repo).UpdateRoute(context.Background(), dtos.UpdateRouteDTO{ID: 404})

	assert.ErrorIs(t, err, domainerrors.ErrRouteNotFound)
}

func TestDeleteRoute_LiberaTodasLasOrdenesYElConductor(t *testing.T) {
	stops := []entities.RouteStop{
		{ID: 1, OrderID: strPtr("ORD-1")},
		{ID: 2, OrderID: nil},
		{ID: 3, OrderID: strPtr("ORD-2")},
	}
	repo := repoConRuta("planned", uintPtr(7), stops)

	err := newRoutesUseCase(repo).DeleteRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Equal(t, []string{"ORD-1", "ORD-2"}, repo.ClearedOrders,
		"borrar la ruta debe liberar todas sus ordenes")
	require.Len(t, repo.DriverStatusCalls, 1)
	assert.Equal(t, "active", repo.DriverStatusCalls[0].Status)
}

func TestDeleteRoute_SinConductor_NoTocaEstados(t *testing.T) {
	repo := repoConRuta("planned", nil, nil)

	err := newRoutesUseCase(repo).DeleteRoute(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Empty(t, repo.DriverStatusCalls)
}

func TestDeleteRoute_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("fk violation")

	t.Run("ruta inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
				return nil, domainerrors.ErrRouteNotFound
			},
		}
		assert.ErrorIs(t, newRoutesUseCase(repo).DeleteRoute(context.Background(), 26, 404),
			domainerrors.ErrRouteNotFound)
	})

	t.Run("al borrar", func(t *testing.T) {
		repo := repoConRuta("planned", nil, nil)
		repo.DeleteRouteFn = func(ctx context.Context, businessID, routeID uint) error { return dbErr }
		assert.ErrorIs(t, newRoutesUseCase(repo).DeleteRoute(context.Background(), 26, 1), dbErr)
	})
}

func TestListRoutes_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 20, 1, 20},
		{-5, 20, 1, 20},
		{1, 0, 1, 20},
		{1, -1, 1, 20},
		{1, 101, 1, 20},
		{1, 100, 1, 100},
		{3, 50, 3, 50},
	}

	for _, tc := range casos {
		var visto dtos.ListRoutesParams
		repo := &mocks.RepositoryMock{
			ListRoutesFn: func(ctx context.Context, p dtos.ListRoutesParams) ([]entities.Route, int64, error) {
				visto = p
				return nil, 0, nil
			},
		}

		_, _, err := newRoutesUseCase(repo).ListRoutes(context.Background(),
			dtos.ListRoutesParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
	}
}

func TestListRoutesParams_Offset(t *testing.T) {
	assert.Equal(t, 0, dtos.ListRoutesParams{Page: 1, PageSize: 20}.Offset())
	assert.Equal(t, 20, dtos.ListRoutesParams{Page: 2, PageSize: 20}.Offset())
	assert.Equal(t, 100, dtos.ListRoutesParams{Page: 3, PageSize: 50}.Offset())
}

func TestListRoutes_PasaLosFiltros(t *testing.T) {
	desde := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	hasta := desde.AddDate(0, 1, 0)
	var visto dtos.ListRoutesParams
	repo := &mocks.RepositoryMock{
		ListRoutesFn: func(ctx context.Context, p dtos.ListRoutesParams) ([]entities.Route, int64, error) {
			visto = p
			return []entities.Route{{ID: 1}}, 1, nil
		},
	}

	got, total, err := newRoutesUseCase(repo).ListRoutes(context.Background(), dtos.ListRoutesParams{
		BusinessID: 26, Status: "planned", DriverID: uintPtr(7),
		DateFrom: &desde, DateTo: &hasta, Search: "norte", Page: 1, PageSize: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto.BusinessID)
	assert.Equal(t, "planned", visto.Status)
	assert.Equal(t, "norte", visto.Search)
	require.NotNil(t, visto.DriverID)
	assert.Equal(t, uint(7), *visto.DriverID)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestListRoutes_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListRoutesFn: func(ctx context.Context, p dtos.ListRoutesParams) ([]entities.Route, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newRoutesUseCase(repo).ListRoutes(context.Background(), dtos.ListRoutesParams{})

	assert.ErrorIs(t, err, dbErr)
}

func TestGetRoute_DelegaAlRepositorio(t *testing.T) {
	var vistoNegocio, vistoRuta uint
	repo := &mocks.RepositoryMock{
		GetRouteByIDFn: func(ctx context.Context, businessID, routeID uint) (*entities.Route, error) {
			vistoNegocio, vistoRuta = businessID, routeID
			return &entities.Route{ID: routeID}, nil
		},
	}

	got, err := newRoutesUseCase(repo).GetRoute(context.Background(), 26, 5)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"el negocio va en la consulta: una ruta de otro negocio no debe ser visible")
	assert.Equal(t, uint(5), vistoRuta)
	assert.Equal(t, uint(5), got.ID)
}

func TestListasDeOpciones_DeleganPorNegocio(t *testing.T) {
	var vistoDrivers, vistoVehicles, vistoOrders uint
	repo := &mocks.RepositoryMock{
		ListDriversForBusinessFn: func(ctx context.Context, businessID uint) ([]dtos.DriverOption, error) {
			vistoDrivers = businessID
			return []dtos.DriverOption{{ID: 1, FirstName: "Juan"}}, nil
		},
		ListVehiclesForBusinessFn: func(ctx context.Context, businessID uint) ([]dtos.VehicleOption, error) {
			vistoVehicles = businessID
			return []dtos.VehicleOption{{ID: 1, LicensePlate: "ABC123"}}, nil
		},
		ListAssignableOrdersFn: func(ctx context.Context, businessID uint) ([]dtos.AssignableOrder, error) {
			vistoOrders = businessID
			return []dtos.AssignableOrder{{ID: "ORD-1"}}, nil
		},
	}
	uc := newRoutesUseCase(repo)

	drivers, err := uc.ListDriversForBusiness(context.Background(), 26)
	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoDrivers)
	assert.Len(t, drivers, 1)

	vehicles, err := uc.ListVehiclesForBusiness(context.Background(), 26)
	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoVehicles)
	assert.Len(t, vehicles, 1)

	orders, err := uc.ListAssignableOrders(context.Background(), 26)
	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoOrders)
	assert.Len(t, orders, 1)
}

func TestListasDeOpciones_ErroresSePropagan(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListDriversForBusinessFn: func(ctx context.Context, businessID uint) ([]dtos.DriverOption, error) {
			return nil, dbErr
		},
		ListVehiclesForBusinessFn: func(ctx context.Context, businessID uint) ([]dtos.VehicleOption, error) {
			return nil, dbErr
		},
		ListAssignableOrdersFn: func(ctx context.Context, businessID uint) ([]dtos.AssignableOrder, error) {
			return nil, dbErr
		},
	}
	uc := newRoutesUseCase(repo)

	_, err := uc.ListDriversForBusiness(context.Background(), 26)
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.ListVehiclesForBusiness(context.Background(), 26)
	assert.ErrorIs(t, err, dbErr)

	_, err = uc.ListAssignableOrders(context.Background(), 26)
	assert.ErrorIs(t, err, dbErr)
}
