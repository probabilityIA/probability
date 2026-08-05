package app

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	domainErrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newStatusUseCase(repo *mocks.RepositoryMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo, mocks.NewSilentLogger())
}

func boolPtr(v bool) *bool { return &v }

func TestCreateOrderStatusMapping_Duplicado_Rechaza(t *testing.T) {
	creado := false
	var vistoTipo uint
	var vistoOriginal string
	repo := &mocks.RepositoryMock{
		ExistsFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error) {
			vistoTipo, vistoOriginal = integrationTypeID, originalStatus
			return true, nil
		},
		CreateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error {
			creado = true
			return nil
		},
	}

	got, err := newStatusUseCase(repo).CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 2, OriginalStatus: "paid", OrderStatusID: 5})

	assert.ErrorIs(t, err, domainErrors.ErrMappingAlreadyExists)
	assert.Nil(t, got)
	assert.False(t, creado)
	assert.Equal(t, uint(2), vistoTipo)
	assert.Equal(t, "paid", vistoOriginal,
		"la unicidad es por canal + estado original, no por estado destino")
}

func TestCreateOrderStatusMapping_ErrorAlVerificarDuplicado_NoCrea(t *testing.T) {
	dbErr := stderrors.New("timeout")
	creado := false
	repo := &mocks.RepositoryMock{
		ExistsFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error) {
			return false, dbErr
		},
		CreateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error {
			creado = true
			return nil
		},
	}

	_, err := newStatusUseCase(repo).CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 2, OriginalStatus: "paid"})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, creado)
}

func TestCreateOrderStatusMapping_NaceActivoAunqueLlegueInactivo(t *testing.T) {
	var creado *entities.OrderStatusMapping
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error {
			creado = m
			m.ID = 9
			return nil
		},
	}

	_, err := newStatusUseCase(repo).CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 2, OriginalStatus: "paid", IsActive: false})

	require.NoError(t, err)
	require.NotNil(t, creado)
	assert.True(t, creado.IsActive, "el mapeo nuevo siempre nace activo, el flag de entrada se ignora")
}

func TestCreateOrderStatusMapping_RecargaConLasRelaciones(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error {
			m.ID = 9
			return nil
		},
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			return &entities.OrderStatusMapping{
				ID:              id,
				OriginalStatus:  "paid",
				IntegrationType: &entities.IntegrationTypeInfo{ID: 2, Code: "shopify", Name: "Shopify"},
				OrderStatus:     &entities.OrderStatusInfo{ID: 5, Code: "confirmed", Name: "Confirmada"},
			}, nil
		},
	}

	got, err := newStatusUseCase(repo).CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 2, OriginalStatus: "paid", OrderStatusID: 5})

	require.NoError(t, err)
	require.NotNil(t, got.IntegrationType)
	assert.Equal(t, "shopify", got.IntegrationType.Code)
	require.NotNil(t, got.OrderStatus)
	assert.Equal(t, "confirmed", got.OrderStatus.Code)
}

func TestCreateOrderStatusMapping_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error { return dbErr },
	}

	got, err := newStatusUseCase(repo).CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 2, OriginalStatus: "paid"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateOrderStatusMapping_ErrorAlRecargar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 2, OriginalStatus: "paid"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateOrderStatusMapping_MismoEstadoOriginalEnCanalesDistintos_AmbosPasan(t *testing.T) {
	var consultas []struct {
		tipo     uint
		original string
	}
	repo := &mocks.RepositoryMock{
		ExistsFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (bool, error) {
			consultas = append(consultas, struct {
				tipo     uint
				original string
			}{integrationTypeID, originalStatus})
			return false, nil
		},
	}
	uc := newStatusUseCase(repo)

	_, err := uc.CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 2, OriginalStatus: "paid", OrderStatusID: 5})
	require.NoError(t, err)

	_, err = uc.CreateOrderStatusMapping(context.Background(),
		&entities.OrderStatusMapping{IntegrationTypeID: 3, OriginalStatus: "paid", OrderStatusID: 7})
	require.NoError(t, err)

	require.Len(t, consultas, 2)
	assert.Equal(t, uint(2), consultas[0].tipo)
	assert.Equal(t, uint(3), consultas[1].tipo)
}

func TestUpdateOrderStatusMapping_NoPermiteCambiarElCanalNiElEstadoActivo(t *testing.T) {
	original := &entities.OrderStatusMapping{
		ID: 9, IntegrationTypeID: 2, OriginalStatus: "paid", OrderStatusID: 5,
		IsActive: true, Description: "vieja",
	}
	var guardado *entities.OrderStatusMapping
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			return original, nil
		},
		UpdateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error {
			guardado = m
			return nil
		},
	}

	_, err := newStatusUseCase(repo).UpdateOrderStatusMapping(context.Background(), 9,
		&entities.OrderStatusMapping{
			IntegrationTypeID: 99, OriginalStatus: "fulfilled", OrderStatusID: 7,
			IsActive: false, Description: "nueva",
		})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, uint(2), guardado.IntegrationTypeID,
		"el canal de origen no se cambia por update: seria mover el mapeo a otra integracion")
	assert.True(t, guardado.IsActive, "el estado activo se cambia solo por el toggle")
	assert.Equal(t, "fulfilled", guardado.OriginalStatus)
	assert.Equal(t, uint(7), guardado.OrderStatusID)
	assert.Equal(t, "nueva", guardado.Description)
}

func TestUpdateOrderStatusMapping_MappingInexistente_NoActualiza(t *testing.T) {
	dbErr := stderrors.New("record not found")
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			return nil, dbErr
		},
		UpdateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error {
			actualizado = true
			return nil
		},
	}

	_, err := newStatusUseCase(repo).UpdateOrderStatusMapping(context.Background(), 404,
		&entities.OrderStatusMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, actualizado)
}

func TestUpdateOrderStatusMapping_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("deadlock")
	repo := &mocks.RepositoryMock{
		UpdateFn: func(ctx context.Context, m *entities.OrderStatusMapping) error { return dbErr },
	}

	got, err := newStatusUseCase(repo).UpdateOrderStatusMapping(context.Background(), 9,
		&entities.OrderStatusMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateOrderStatusMapping_RecargaDespuesDeGuardar(t *testing.T) {
	llamadas := 0
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			llamadas++
			if llamadas == 1 {
				return &entities.OrderStatusMapping{ID: id, IntegrationTypeID: 2}, nil
			}
			return &entities.OrderStatusMapping{
				ID: id, OrderStatus: &entities.OrderStatusInfo{ID: 7, Code: "shipped"},
			}, nil
		},
	}

	got, err := newStatusUseCase(repo).UpdateOrderStatusMapping(context.Background(), 9,
		&entities.OrderStatusMapping{OrderStatusID: 7})

	require.NoError(t, err)
	assert.Equal(t, 2, llamadas, "se lee antes de guardar y se recarga despues")
	require.NotNil(t, got.OrderStatus)
	assert.Equal(t, "shipped", got.OrderStatus.Code)
}

func TestUpdateOrderStatusMapping_ErrorAlRecargar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	llamadas := 0
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			llamadas++
			if llamadas == 1 {
				return &entities.OrderStatusMapping{ID: id}, nil
			}
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).UpdateOrderStatusMapping(context.Background(), 9,
		&entities.OrderStatusMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetOrderStatusMapping_DelegaYPropaga(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			visto = id
			return &entities.OrderStatusMapping{ID: id, OriginalStatus: "paid"}, nil
		},
	}

	got, err := newStatusUseCase(repo).GetOrderStatusMapping(context.Background(), 9)

	require.NoError(t, err)
	assert.Equal(t, uint(9), visto)
	assert.Equal(t, "paid", got.OriginalStatus)

	dbErr := stderrors.New("no encontrado")
	repoErr := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			return nil, dbErr
		},
	}
	_, err = newStatusUseCase(repoErr).GetOrderStatusMapping(context.Background(), 9)
	assert.ErrorIs(t, err, dbErr)
}

func TestListOrderStatusMappings_PasaFiltrosYTotal(t *testing.T) {
	esperados := map[string]interface{}{"integration_type_id": uint(2), "is_active": true}
	var recibidos map[string]interface{}
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error) {
			recibidos = filters
			return []entities.OrderStatusMapping{{ID: 1}, {ID: 2}}, 2, nil
		},
	}

	got, total, err := newStatusUseCase(repo).ListOrderStatusMappings(context.Background(), esperados)

	require.NoError(t, err)
	assert.Equal(t, esperados, recibidos)
	assert.Len(t, got, 2)
	assert.Equal(t, int64(2), total)
}

func TestListOrderStatusMappings_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, filters map[string]interface{}) ([]entities.OrderStatusMapping, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, total, err := newStatusUseCase(repo).ListOrderStatusMappings(context.Background(), nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.Zero(t, total)
}

func TestDeleteOrderStatusMapping_DelegaYPropaga(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}

	require.NoError(t, newStatusUseCase(repo).DeleteOrderStatusMapping(context.Background(), 9))
	assert.Equal(t, uint(9), borradoID)

	dbErr := stderrors.New("fk violation")
	repoErr := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, id uint) error { return dbErr },
	}
	assert.ErrorIs(t, newStatusUseCase(repoErr).DeleteOrderStatusMapping(context.Background(), 9), dbErr)
}

func TestToggleOrderStatusMappingActive_DelegaYDevuelveElEstadoNuevo(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		ToggleActiveFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			visto = id
			return &entities.OrderStatusMapping{ID: id, IsActive: false}, nil
		},
	}

	got, err := newStatusUseCase(repo).ToggleOrderStatusMappingActive(context.Background(), 9)

	require.NoError(t, err)
	assert.Equal(t, uint(9), visto)
	assert.False(t, got.IsActive)
}

func TestToggleOrderStatusMappingActive_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("no encontrado")
	repo := &mocks.RepositoryMock{
		ToggleActiveFn: func(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).ToggleOrderStatusMappingActive(context.Background(), 9)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListOrderStatuses_PasaElFiltroIsActive(t *testing.T) {
	for _, filtro := range []*bool{nil, boolPtr(true), boolPtr(false)} {
		var visto *bool
		repo := &mocks.RepositoryMock{
			ListOrderStatusesFn: func(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error) {
				visto = isActive
				return []entities.OrderStatusInfo{{ID: 1, Code: "pending"}}, nil
			},
		}

		got, err := newStatusUseCase(repo).ListOrderStatuses(context.Background(), filtro)

		require.NoError(t, err)
		assert.Equal(t, filtro, visto)
		assert.Len(t, got, 1)
	}
}

func TestListOrderStatuses_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListOrderStatusesFn: func(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error) {
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).ListOrderStatuses(context.Background(), nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListFulfillmentStatuses_PasaElFiltroYPropaga(t *testing.T) {
	var visto *bool
	repo := &mocks.RepositoryMock{
		ListFulfillmentStatusesFn: func(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error) {
			visto = isActive
			return []entities.FulfillmentStatusInfo{{ID: 1, Code: "unfulfilled"}}, nil
		},
	}

	got, err := newStatusUseCase(repo).ListFulfillmentStatuses(context.Background(), boolPtr(true))

	require.NoError(t, err)
	require.NotNil(t, visto)
	assert.True(t, *visto)
	assert.Equal(t, "unfulfilled", got[0].Code)

	dbErr := stderrors.New("timeout")
	repoErr := &mocks.RepositoryMock{
		ListFulfillmentStatusesFn: func(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error) {
			return nil, dbErr
		},
	}
	_, err = newStatusUseCase(repoErr).ListFulfillmentStatuses(context.Background(), nil)
	assert.ErrorIs(t, err, dbErr)
}

func TestCreateOrderStatus_PasaLaEntidadTalCual(t *testing.T) {
	entrada := &entities.OrderStatusInfo{
		Code: "confirmed", Name: "Confirmada", Category: "open",
		Color: "#0F0", Priority: 3, IsActive: true,
	}
	var recibido *entities.OrderStatusInfo
	repo := &mocks.RepositoryMock{
		CreateOrderStatusFn: func(ctx context.Context, s *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
			recibido = s
			s.ID = 9
			return s, nil
		},
	}

	got, err := newStatusUseCase(repo).CreateOrderStatus(context.Background(), entrada)

	require.NoError(t, err)
	assert.Same(t, entrada, recibido)
	assert.Equal(t, uint(9), got.ID)
	assert.Equal(t, 3, got.Priority)
}

func TestCreateOrderStatus_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("duplicate code")
	repo := &mocks.RepositoryMock{
		CreateOrderStatusFn: func(ctx context.Context, s *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).CreateOrderStatus(context.Background(), &entities.OrderStatusInfo{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetOrderStatus_DelegaYPropaga(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		GetOrderStatusByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusInfo, error) {
			visto = id
			return &entities.OrderStatusInfo{ID: id, Code: "confirmed"}, nil
		},
	}

	got, err := newStatusUseCase(repo).GetOrderStatus(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), visto)
	assert.Equal(t, "confirmed", got.Code)

	repoErr := &mocks.RepositoryMock{
		GetOrderStatusByIDFn: func(ctx context.Context, id uint) (*entities.OrderStatusInfo, error) {
			return nil, domainErrors.ErrOrderStatusNotFound
		},
	}
	_, err = newStatusUseCase(repoErr).GetOrderStatus(context.Background(), 404)
	assert.ErrorIs(t, err, domainErrors.ErrOrderStatusNotFound)
}

func TestUpdateOrderStatus_PasaIDYEntidad(t *testing.T) {
	var vistoID uint
	var recibido *entities.OrderStatusInfo
	entrada := &entities.OrderStatusInfo{Name: "Confirmada", Priority: 5}
	repo := &mocks.RepositoryMock{
		UpdateOrderStatusFn: func(ctx context.Context, id uint, s *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
			vistoID, recibido = id, s
			return s, nil
		},
	}

	got, err := newStatusUseCase(repo).UpdateOrderStatus(context.Background(), 5, entrada)

	require.NoError(t, err)
	assert.Equal(t, uint(5), vistoID)
	assert.Same(t, entrada, recibido)
	assert.Equal(t, 5, got.Priority)
}

func TestUpdateOrderStatus_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		UpdateOrderStatusFn: func(ctx context.Context, id uint, s *entities.OrderStatusInfo) (*entities.OrderStatusInfo, error) {
			return nil, domainErrors.ErrOrderStatusNotFound
		},
	}

	got, err := newStatusUseCase(repo).UpdateOrderStatus(context.Background(), 404, &entities.OrderStatusInfo{})

	assert.ErrorIs(t, err, domainErrors.ErrOrderStatusNotFound)
	assert.Nil(t, got)
}

func TestDeleteOrderStatus_ConMapeosAsociados_SePropagaElRechazo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		DeleteOrderStatusFn: func(ctx context.Context, id uint) error {
			return domainErrors.ErrOrderStatusHasMappings
		},
	}

	err := newStatusUseCase(repo).DeleteOrderStatus(context.Background(), 5)

	assert.ErrorIs(t, err, domainErrors.ErrOrderStatusHasMappings,
		"borrar un estado con mapeos dejaria ordenes apuntando a la nada")
}

func TestDeleteOrderStatus_Exitoso(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		DeleteOrderStatusFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}

	require.NoError(t, newStatusUseCase(repo).DeleteOrderStatus(context.Background(), 5))
	assert.Equal(t, uint(5), borradoID)
}

func TestListEcommerceIntegrationTypes_FiltraPorNegocio(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		ListEcommerceIntegrationTypesFn: func(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error) {
			visto = businessID
			return []entities.IntegrationTypeInfo{{ID: 2, Code: "shopify", Name: "Shopify"}}, nil
		},
	}

	got, err := newStatusUseCase(repo).ListEcommerceIntegrationTypes(context.Background(), 26)

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto,
		"los canales se listan por negocio: un negocio no debe ver los de otro")
	assert.Equal(t, "shopify", got[0].Code)
}

func TestListEcommerceIntegrationTypes_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListEcommerceIntegrationTypesFn: func(ctx context.Context, businessID uint) ([]entities.IntegrationTypeInfo, error) {
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).ListEcommerceIntegrationTypes(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListChannelStatuses_PasaCanalYFiltro(t *testing.T) {
	var vistoTipo uint
	var vistoActivo *bool
	repo := &mocks.RepositoryMock{
		ListChannelStatusesFn: func(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error) {
			vistoTipo, vistoActivo = integrationTypeID, isActive
			return []entities.ChannelStatusInfo{{ID: 1, Code: "paid"}}, nil
		},
	}

	got, err := newStatusUseCase(repo).ListChannelStatuses(context.Background(), 2, boolPtr(true))

	require.NoError(t, err)
	assert.Equal(t, uint(2), vistoTipo)
	require.NotNil(t, vistoActivo)
	assert.True(t, *vistoActivo)
	assert.Equal(t, "paid", got[0].Code)
}

func TestListChannelStatuses_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListChannelStatusesFn: func(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error) {
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).ListChannelStatuses(context.Background(), 2, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateChannelStatus_PasaLaEntidadTalCual(t *testing.T) {
	entrada := &entities.ChannelStatusInfo{
		IntegrationTypeID: 2, Code: "paid", Name: "Pagado", DisplayOrder: 3, IsActive: true,
	}
	var recibido *entities.ChannelStatusInfo
	repo := &mocks.RepositoryMock{
		CreateChannelStatusFn: func(ctx context.Context, s *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
			recibido = s
			s.ID = 9
			return s, nil
		},
	}

	got, err := newStatusUseCase(repo).CreateChannelStatus(context.Background(), entrada)

	require.NoError(t, err)
	assert.Same(t, entrada, recibido)
	assert.Equal(t, uint(9), got.ID)
	assert.Equal(t, uint(2), got.IntegrationTypeID)
}

func TestCreateChannelStatus_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("duplicate code")
	repo := &mocks.RepositoryMock{
		CreateChannelStatusFn: func(ctx context.Context, s *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
			return nil, dbErr
		},
	}

	got, err := newStatusUseCase(repo).CreateChannelStatus(context.Background(), &entities.ChannelStatusInfo{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetChannelStatus_DelegaYPropaga(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		GetChannelStatusByIDFn: func(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error) {
			visto = id
			return &entities.ChannelStatusInfo{
				ID: id, Code: "paid",
				IntegrationType: &entities.IntegrationTypeInfo{ID: 2, Code: "shopify"},
			}, nil
		},
	}

	got, err := newStatusUseCase(repo).GetChannelStatus(context.Background(), 9)

	require.NoError(t, err)
	assert.Equal(t, uint(9), visto)
	require.NotNil(t, got.IntegrationType)
	assert.Equal(t, "shopify", got.IntegrationType.Code)

	repoErr := &mocks.RepositoryMock{
		GetChannelStatusByIDFn: func(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error) {
			return nil, domainErrors.ErrChannelStatusNotFound
		},
	}
	_, err = newStatusUseCase(repoErr).GetChannelStatus(context.Background(), 404)
	assert.ErrorIs(t, err, domainErrors.ErrChannelStatusNotFound)
}

func TestUpdateChannelStatus_PasaIDYEntidad(t *testing.T) {
	var vistoID uint
	entrada := &entities.ChannelStatusInfo{Name: "Pagado", DisplayOrder: 7}
	repo := &mocks.RepositoryMock{
		UpdateChannelStatusFn: func(ctx context.Context, id uint, s *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
			vistoID = id
			return s, nil
		},
	}

	got, err := newStatusUseCase(repo).UpdateChannelStatus(context.Background(), 9, entrada)

	require.NoError(t, err)
	assert.Equal(t, uint(9), vistoID)
	assert.Equal(t, 7, got.DisplayOrder)
}

func TestUpdateChannelStatus_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		UpdateChannelStatusFn: func(ctx context.Context, id uint, s *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
			return nil, domainErrors.ErrChannelStatusNotFound
		},
	}

	got, err := newStatusUseCase(repo).UpdateChannelStatus(context.Background(), 404, &entities.ChannelStatusInfo{})

	assert.ErrorIs(t, err, domainErrors.ErrChannelStatusNotFound)
	assert.Nil(t, got)
}

func TestDeleteChannelStatus_ConMapeosAsociados_SePropagaElRechazo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		DeleteChannelStatusFn: func(ctx context.Context, id uint) error {
			return domainErrors.ErrChannelStatusHasMappings
		},
	}

	err := newStatusUseCase(repo).DeleteChannelStatus(context.Background(), 9)

	assert.ErrorIs(t, err, domainErrors.ErrChannelStatusHasMappings)
}

func TestDeleteChannelStatus_Exitoso(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		DeleteChannelStatusFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}

	require.NoError(t, newStatusUseCase(repo).DeleteChannelStatus(context.Background(), 9))
	assert.Equal(t, uint(9), borradoID)
}
