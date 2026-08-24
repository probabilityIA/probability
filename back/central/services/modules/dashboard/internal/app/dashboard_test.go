package app

import (
	"context"
	stderrors "errors"
	"sync"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDashboardUseCase(repo *mocks.RepositoryMock) domain.IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo, nil, mocks.NewSilentLogger())
}

func uintPtr(v uint) *uint { return &v }

func statsDe(t *testing.T, repo *mocks.RepositoryMock, businessID *uint) *domain.DashboardStats {
	t.Helper()
	got, err := newDashboardUseCase(repo).
		GetDashboardStats(context.Background(), businessID, nil, nil, nil, nil, false)
	require.NoError(t, err)
	require.NotNil(t, got)
	return got
}

func TestGetDashboardStats_ConsultasCriticas_AbortanElDashboard(t *testing.T) {
	dbErr := stderrors.New("timeout")

	casos := []struct {
		nombre string
		aplica func(*mocks.RepositoryMock)
	}{
		{"total de ordenes", func(m *mocks.RepositoryMock) {
			m.GetTotalOrdersFn = func(ctx context.Context, b, i *uint, s, e *time.Time) (int64, error) {
				return 0, dbErr
			}
		}},
		{"ordenes por tipo de integracion", func(m *mocks.RepositoryMock) {
			m.GetOrdersByIntegrationTypeFn = func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.OrderCountByIntegrationType, error) {
				return nil, dbErr
			}
		}},
		{"top clientes", func(m *mocks.RepositoryMock) {
			m.GetTopCustomersFn = func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopCustomer, error) {
				return nil, dbErr
			}
		}},
		{"ordenes por ubicacion", func(m *mocks.RepositoryMock) {
			m.GetOrdersByLocationFn = func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.OrderCountByLocation, error) {
				return nil, dbErr
			}
		}},
		{"top transportadores", func(m *mocks.RepositoryMock) {
			m.GetTopDriversFn = func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopDriver, error) {
				return nil, dbErr
			}
		}},
		{"transportadores por ubicacion", func(m *mocks.RepositoryMock) {
			m.GetDriversByLocationFn = func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.DriverByLocation, error) {
				return nil, dbErr
			}
		}},
		{"top productos", func(m *mocks.RepositoryMock) {
			m.GetTopProductsFn = func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopProduct, error) {
				return nil, dbErr
			}
		}},
		{"productos por categoria", func(m *mocks.RepositoryMock) {
			m.GetProductsByCategoryFn = func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ProductByCategory, error) {
				return nil, dbErr
			}
		}},
		{"productos por marca", func(m *mocks.RepositoryMock) {
			m.GetProductsByBrandFn = func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ProductByBrand, error) {
				return nil, dbErr
			}
		}},
		{"envios por estado", func(m *mocks.RepositoryMock) {
			m.GetShipmentsByStatusFilteredFn = func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ShipmentsByStatus, error) {
				return nil, dbErr
			}
		}},
		{"envios por transportista", func(m *mocks.RepositoryMock) {
			m.GetShipmentsByCarrierFn = func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ShipmentsByCarrier, error) {
				return nil, dbErr
			}
		}},
		{"envios por almacen", func(m *mocks.RepositoryMock) {
			m.GetShipmentsByWarehouseFn = func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.ShipmentsByWarehouse, error) {
				return nil, dbErr
			}
		}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}
			tc.aplica(repo)

			got, err := newDashboardUseCase(repo).
				GetDashboardStats(context.Background(), nil, nil, nil, nil, nil, false)

			assert.ErrorIs(t, err, dbErr)
			assert.Nil(t, got)
		})
	}
}

func TestGetDashboardStats_ConsultasSecundarias_DegradanSinTumbarElDashboard(t *testing.T) {
	dbErr := stderrors.New("timeout")

	t.Run("ordenes de hoy cae a cero", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetTotalOrdersFn: func(ctx context.Context, b, i *uint, s, e *time.Time) (int64, error) {
				return 500, nil
			},
			GetOrdersTodayFn: func(ctx context.Context, b, i *uint, s, e *time.Time) (int64, error) {
				return 0, dbErr
			},
		}

		got := statsDe(t, repo, nil)

		assert.EqualValues(t, 500, got.TotalOrders, "el resto del dashboard sigue vivo")
		assert.Zero(t, got.OrdersToday)
	})

	t.Run("envios por transportista de hoy cae a lista vacia", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetShipmentsByCarrierTodayFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ShipmentsByCarrier, error) {
				return nil, dbErr
			},
		}

		got := statsDe(t, repo, nil)

		require.NotNil(t, got.ShipmentsByCarrierToday)
		assert.Empty(t, got.ShipmentsByCarrierToday)
	})

	t.Run("envios por dia de la semana cae a lista vacia", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetShipmentsByDayOfWeekFn: func(ctx context.Context, b, i *uint, w *time.Time) ([]domain.ShipmentsByDayOfWeek, error) {
				return nil, dbErr
			},
		}

		got := statsDe(t, repo, nil)

		require.NotNil(t, got.ShipmentsByDayOfWeek)
		assert.Empty(t, got.ShipmentsByDayOfWeek)
	})

	t.Run("ordenes por departamento cae a lista vacia", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetOrdersByDepartmentFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.OrdersByDepartment, error) {
				return nil, dbErr
			},
		}

		got := statsDe(t, repo, nil)

		require.NotNil(t, got.OrdersByDepartment)
		assert.Empty(t, got.OrdersByDepartment)
	})

	t.Run("ordenes por mes cae a lista vacia", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetOrdersByMonthFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.OrdersByMonth, error) {
				return nil, dbErr
			},
		}

		got := statsDe(t, repo, nil)

		require.NotNil(t, got.OrdersByMonth)
		assert.Empty(t, got.OrdersByMonth)
	})

	t.Run("ordenes por semana cae a lista vacia", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetOrdersByWeekFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.OrdersByWeek, error) {
				return nil, dbErr
			},
		}

		got := statsDe(t, repo, nil)

		require.NotNil(t, got.OrdersByWeek)
		assert.Empty(t, got.OrdersByWeek)
	})

	t.Run("ordenes por business solo se loguea", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetOrdersByBusinessFn: func(ctx context.Context, l int, s, e *time.Time) ([]domain.OrdersByBusiness, error) {
				return nil, dbErr
			},
		}

		got := statsDe(t, repo, nil)

		assert.Empty(t, got.OrdersByBusiness)
	})
}

func TestGetDashboardStats_SuperAdminSinFiltro_ConsultaOrdenesPorBusiness(t *testing.T) {
	var vistoLimite int
	repo := &mocks.RepositoryMock{
		GetOrdersByBusinessFn: func(ctx context.Context, limit int, s, e *time.Time) ([]domain.OrdersByBusiness, error) {
			vistoLimite = limit
			return []domain.OrdersByBusiness{{BusinessID: 26, OrderCount: 100}}, nil
		},
	}

	got := statsDe(t, repo, nil)

	assert.True(t, repo.WasCalled("GetOrdersByBusiness"))
	assert.Equal(t, 10, vistoLimite)
	require.Len(t, got.OrdersByBusiness, 1)
	assert.EqualValues(t, 26, got.OrdersByBusiness[0].BusinessID)
}

func TestGetDashboardStats_ConFiltroDeBusiness_NoConsultaNiMuestraOtrosNegocios(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetOrdersByBusinessFn: func(ctx context.Context, limit int, s, e *time.Time) ([]domain.OrdersByBusiness, error) {
			return []domain.OrdersByBusiness{{BusinessID: 99, OrderCount: 999}}, nil
		},
	}

	got := statsDe(t, repo, uintPtr(26))

	assert.False(t, repo.WasCalled("GetOrdersByBusiness"),
		"un negocio filtrado no debe ver el ranking de todos los negocios")
	require.NotNil(t, got.OrdersByBusiness)
	assert.Empty(t, got.OrdersByBusiness)
}

func TestGetDashboardStats_PropagaElBusinessIDATodasLasConsultas(t *testing.T) {
	negocio := uintPtr(26)
	var muVistos sync.Mutex
	var vistos []*uint
	anotaVisto := func(b *uint) {
		muVistos.Lock()
		defer muVistos.Unlock()
		vistos = append(vistos, b)
	}
	repo := &mocks.RepositoryMock{
		GetTotalOrdersFn: func(ctx context.Context, b, i *uint, s, e *time.Time) (int64, error) {
			anotaVisto(b)
			return 0, nil
		},
		GetTopCustomersFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopCustomer, error) {
			anotaVisto(b)
			return nil, nil
		},
		GetShipmentsByCarrierFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ShipmentsByCarrier, error) {
			anotaVisto(b)
			return nil, nil
		},
	}

	statsDe(t, repo, negocio)

	require.Len(t, vistos, 3)
	for i, v := range vistos {
		require.NotNil(t, v, "consulta %d recibio nil en vez del negocio", i)
		assert.Equal(t, uint(26), *v,
			"si una consulta pierde el business_id se filtrarian datos de otros negocios")
	}
}

func TestGetDashboardStats_PropagaElIntegrationIDYElRangoDeFechas(t *testing.T) {
	integracion := uintPtr(7)
	desde := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	hasta := desde.AddDate(0, 1, 0)

	var vistaIntegracion *uint
	var vistoDesde, vistoHasta *time.Time
	repo := &mocks.RepositoryMock{
		GetTotalOrdersFn: func(ctx context.Context, b, i *uint, s, e *time.Time) (int64, error) {
			vistaIntegracion, vistoDesde, vistoHasta = i, s, e
			return 0, nil
		},
	}

	_, err := newDashboardUseCase(repo).GetDashboardStats(
		context.Background(), nil, integracion, nil, &desde, &hasta, false)

	require.NoError(t, err)
	require.NotNil(t, vistaIntegracion)
	assert.Equal(t, uint(7), *vistaIntegracion)
	require.NotNil(t, vistoDesde)
	assert.Equal(t, desde, *vistoDesde)
	require.NotNil(t, vistoHasta)
	assert.Equal(t, hasta, *vistoHasta)
}

func TestGetDashboardStats_ElInicioDeSemanaSoloVaALaConsultaDeDiaDeSemana(t *testing.T) {
	inicioSemana := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	var visto *time.Time
	repo := &mocks.RepositoryMock{
		GetShipmentsByDayOfWeekFn: func(ctx context.Context, b, i *uint, w *time.Time) ([]domain.ShipmentsByDayOfWeek, error) {
			visto = w
			return nil, nil
		},
	}

	_, err := newDashboardUseCase(repo).GetDashboardStats(
		context.Background(), nil, nil, &inicioSemana, nil, nil, false)

	require.NoError(t, err)
	require.NotNil(t, visto)
	assert.Equal(t, inicioSemana, *visto)
}

func TestGetDashboardStats_UsaLaVersionFiltradaDeEnviosPorEstado(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	statsDe(t, repo, nil)

	assert.True(t, repo.WasCalled("GetShipmentsByStatusFiltered"))
	assert.False(t, repo.WasCalled("GetShipmentsByStatus"),
		"el dashboard usa la version filtrada; la sin filtrar quedo sin uso en el port")
}

func TestGetDashboardStats_LosTopsPidenDiez(t *testing.T) {
	var muLimites sync.Mutex
	limites := map[string]int{}
	anota := func(nombre string, l int) {
		muLimites.Lock()
		defer muLimites.Unlock()
		limites[nombre] = l
	}
	repo := &mocks.RepositoryMock{
		GetTopCustomersFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopCustomer, error) {
			anota("customers", l)
			return nil, nil
		},
		GetOrdersByLocationFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.OrderCountByLocation, error) {
			anota("location", l)
			return nil, nil
		},
		GetTopDriversFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopDriver, error) {
			anota("drivers", l)
			return nil, nil
		},
		GetDriversByLocationFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.DriverByLocation, error) {
			anota("driversByLocation", l)
			return nil, nil
		},
		GetTopProductsFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopProduct, error) {
			anota("products", l)
			return nil, nil
		},
		GetShipmentsByWarehouseFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.ShipmentsByWarehouse, error) {
			anota("warehouse", l)
			return nil, nil
		},
	}

	statsDe(t, repo, nil)

	require.Len(t, limites, 6)
	for nombre, l := range limites {
		assert.Equal(t, 10, l, "el top de %s deberia pedir 10", nombre)
	}
}

func TestGetDashboardStats_ArmaElAgregadoCompleto(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetTotalOrdersFn: func(ctx context.Context, b, i *uint, s, e *time.Time) (int64, error) {
			return 500, nil
		},
		GetOrdersTodayFn: func(ctx context.Context, b, i *uint, s, e *time.Time) (int64, error) {
			return 12, nil
		},
		GetOrdersByIntegrationTypeFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.OrderCountByIntegrationType, error) {
			return []domain.OrderCountByIntegrationType{{IntegrationType: "shopify", Count: 300}}, nil
		},
		GetTopCustomersFn: func(ctx context.Context, b, i *uint, l int, s, e *time.Time) ([]domain.TopCustomer, error) {
			return []domain.TopCustomer{{CustomerName: "Ana"}}, nil
		},
		GetShipmentsByStatusFilteredFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ShipmentsByStatus, error) {
			return []domain.ShipmentsByStatus{{Status: "delivered", Count: 200}}, nil
		},
		GetShipmentsByCarrierFn: func(ctx context.Context, b, i *uint, s, e *time.Time) ([]domain.ShipmentsByCarrier, error) {
			return []domain.ShipmentsByCarrier{{Carrier: "Servientrega", Count: 150}}, nil
		},
	}

	got := statsDe(t, repo, nil)

	assert.EqualValues(t, 500, got.TotalOrders)
	assert.EqualValues(t, 12, got.OrdersToday)
	require.Len(t, got.OrdersByIntegrationType, 1)
	assert.Equal(t, "shopify", got.OrdersByIntegrationType[0].IntegrationType)
	require.Len(t, got.TopCustomers, 1)
	assert.Equal(t, "Ana", got.TopCustomers[0].CustomerName)
	require.Len(t, got.ShipmentsByStatus, 1)
	assert.Equal(t, "delivered", got.ShipmentsByStatus[0].Status)
	require.Len(t, got.ShipmentsByCarrier, 1)
	assert.Equal(t, "Servientrega", got.ShipmentsByCarrier[0].Carrier)
}

func TestGetDashboardStats_ConsultaTodasLasSeccionesEnUnaSolaLlamada(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	statsDe(t, repo, nil)

	esperadas := []string{
		"GetTotalOrders", "GetOrdersToday", "GetOrdersByIntegrationType",
		"GetTopCustomers", "GetOrdersByLocation", "GetTopDrivers",
		"GetDriversByLocation", "GetTopProducts", "GetProductsByCategory",
		"GetProductsByBrand", "GetShipmentsByStatusFiltered", "GetShipmentsByCarrier",
		"GetShipmentsByCarrierToday", "GetShipmentsByWarehouse", "GetShipmentsByDayOfWeek",
		"GetOrdersByDepartment", "GetOrdersByMonth", "GetOrdersByWeek", "GetOrdersByBusiness",
	}
	for _, m := range esperadas {
		assert.True(t, repo.WasCalled(m), "falto consultar %s", m)
	}
	assert.Len(t, repo.Called(), len(esperadas),
		"cada seccion se consulta exactamente una vez")
}

func TestGetTopSellingDays_PasaLosFiltrosYElLimite(t *testing.T) {
	var vistoNegocio, vistaIntegracion *uint
	var vistoLimite int
	repo := &mocks.RepositoryMock{
		GetTopSellingDaysFn: func(ctx context.Context, b, i *uint, limit int) ([]domain.TopSellingDay, error) {
			vistoNegocio, vistaIntegracion, vistoLimite = b, i, limit
			return []domain.TopSellingDay{{Date: "2026-06-15", Total: 42}}, nil
		},
	}

	got, err := newDashboardUseCase(repo).
		GetTopSellingDays(context.Background(), uintPtr(26), uintPtr(7), 5)

	require.NoError(t, err)
	require.NotNil(t, vistoNegocio)
	assert.Equal(t, uint(26), *vistoNegocio)
	require.NotNil(t, vistaIntegracion)
	assert.Equal(t, uint(7), *vistaIntegracion)
	assert.Equal(t, 5, vistoLimite)
	require.Len(t, got, 1)
	assert.Equal(t, "2026-06-15", got[0].Date)
}

func TestGetTopSellingDays_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetTopSellingDaysFn: func(ctx context.Context, b, i *uint, limit int) ([]domain.TopSellingDay, error) {
			return nil, dbErr
		},
	}

	got, err := newDashboardUseCase(repo).
		GetTopSellingDays(context.Background(), nil, nil, 5)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetTopSellingDays_SuperAdminSinNegocio_ConsultaTodo(t *testing.T) {
	var vistoNegocio *uint
	consultado := false
	repo := &mocks.RepositoryMock{
		GetTopSellingDaysFn: func(ctx context.Context, b, i *uint, limit int) ([]domain.TopSellingDay, error) {
			consultado = true
			vistoNegocio = b
			return nil, nil
		},
	}

	_, err := newDashboardUseCase(repo).GetTopSellingDays(context.Background(), nil, nil, 5)

	require.NoError(t, err)
	assert.True(t, consultado)
	assert.Nil(t, vistoNegocio, "sin negocio el super admin ve el consolidado")
}
