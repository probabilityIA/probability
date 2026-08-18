package usecaseshipment

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validCreateRequest() *domain.CreateShipmentRequest {
	return &domain.CreateShipmentRequest{
		OrderID:            s("order-uuid"),
		ClientName:         "Juan Perez",
		DestinationAddress: "Calle 1 # 2-3",
	}
}

func TestCreateShipment_SinOrderIDNiDatosDeCliente_RetornaError(t *testing.T) {
	uc := newUseCaseForCost(&mocks.RepositoryMock{}, nil)

	resp, err := uc.CreateShipment(context.Background(), &domain.CreateShipmentRequest{})

	assert.ErrorIs(t, err, domain.ErrOrderIDRequired)
	assert.Nil(t, resp)
}

func TestCreateShipment_SinOrderIDPeroConClienteYDireccion_EsValido(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CreateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			shipment.ID = 100
			return nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.CreateShipment(context.Background(), &domain.CreateShipmentRequest{
		ClientName:         "Juan Perez",
		DestinationAddress: "Calle 1 # 2-3",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint(100), resp.ID)
}

func TestCreateShipment_StatusPorDefectoEsPending(t *testing.T) {
	var saved *domain.Shipment
	repo := &mocks.RepositoryMock{
		CreateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			shipment.ID = 1
			saved = shipment
			return nil
		},
		GetOrderBusinessIDFn: func(ctx context.Context, orderUUID string) (uint, error) { return 0, nil },
	}
	uc := newUseCaseForCost(repo, nil)

	_, err := uc.CreateShipment(context.Background(), validCreateRequest())

	require.NoError(t, err)
	require.NotNil(t, saved)
	assert.Equal(t, "pending", saved.Status)
}

func TestCreateShipment_StatusExplicitoSeRespeta(t *testing.T) {
	var saved *domain.Shipment
	repo := &mocks.RepositoryMock{
		CreateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			shipment.ID = 1
			saved = shipment
			return nil
		},
		GetOrderBusinessIDFn: func(ctx context.Context, orderUUID string) (uint, error) { return 0, nil },
	}
	uc := newUseCaseForCost(repo, nil)

	req := validCreateRequest()
	req.Status = "in_transit"

	_, err := uc.CreateShipment(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "in_transit", saved.Status)
}

func TestCreateShipment_TrackingDuplicadoEnLaMismaOrden_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ShipmentExistsFn: func(ctx context.Context, orderID, trackingNumber string) (bool, error) {
			return true, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	req := validCreateRequest()
	req.TrackingNumber = s("034057376067")

	resp, err := uc.CreateShipment(context.Background(), req)

	assert.ErrorIs(t, err, domain.ErrShipmentAlreadyExists)
	assert.Nil(t, resp)
}

func TestCreateShipment_FallaVerificarDuplicado_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ShipmentExistsFn: func(ctx context.Context, orderID, trackingNumber string) (bool, error) {
			return false, dbErr
		},
	}
	uc := newUseCaseForCost(repo, nil)

	req := validCreateRequest()
	req.TrackingNumber = s("034057376067")

	resp, err := uc.CreateShipment(context.Background(), req)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestCreateShipment_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	repo := &mocks.RepositoryMock{
		CreateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error { return dbErr },
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.CreateShipment(context.Background(), validCreateRequest())

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestCreateShipment_AplicaCostoDeCarrierAntesDeGuardar(t *testing.T) {
	var saved *domain.Shipment
	repo := &mocks.RepositoryMock{
		CreateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			shipment.ID = 5
			saved = shipment
			return nil
		},
		GetOrderBusinessIDFn: func(ctx context.Context, orderUUID string) (uint, error) { return 36, nil },
	}
	reader := &marginReaderStub{margin: domain.ShippingMargin{MarginAmount: 2000}}
	uc := newUseCaseForCost(repo, reader)

	req := validCreateRequest()
	req.TotalCost = f64(20000)
	req.CarrierCode = s("servientrega")

	_, err := uc.CreateShipment(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, saved.CarrierCost)
	assert.InDelta(t, 18000.0, *saved.CarrierCost, 0.001)
}

func TestGetShipmentByID_IDCero_RetornaError(t *testing.T) {
	uc := newUseCaseForCost(&mocks.RepositoryMock{}, nil)

	resp, err := uc.GetShipmentByID(context.Background(), 0)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "shipment ID is required")
}

func TestGetShipmentByID_NoExiste_RetornaNotFound(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) { return nil, nil },
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.GetShipmentByID(context.Background(), 42)

	assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
	assert.Nil(t, resp)
}

func TestGetShipmentByID_Existe_RetornaResponse(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, Status: "delivered", ClientName: "Ana"}, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.GetShipmentByID(context.Background(), 42)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint(42), resp.ID)
	assert.Equal(t, "delivered", resp.Status)
}

func TestGetShipmentByID_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) { return nil, dbErr },
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.GetShipmentByID(context.Background(), 42)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestListShipments_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "pagina cero", page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "pagina negativa", page: -5, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "page size cero", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{name: "page size negativo", page: 1, pageSize: -1, wantPage: 1, wantPageSize: 10},
		{name: "page size mayor al maximo", page: 1, pageSize: 500, wantPage: 1, wantPageSize: 10},
		{name: "page size en el limite", page: 2, pageSize: 100, wantPage: 2, wantPageSize: 100},
		{name: "valores validos", page: 3, pageSize: 25, wantPage: 3, wantPageSize: 25},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPage, gotPageSize int
			repo := &mocks.RepositoryMock{
				ListShipmentsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error) {
					gotPage, gotPageSize = page, pageSize
					return []domain.Shipment{}, 0, nil
				},
			}
			uc := newUseCaseForCost(repo, nil)

			resp, err := uc.ListShipments(context.Background(), tc.page, tc.pageSize, nil)

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, gotPage)
			assert.Equal(t, tc.wantPageSize, gotPageSize)
			assert.Equal(t, tc.wantPage, resp.Page)
			assert.Equal(t, tc.wantPageSize, resp.PageSize)
		})
	}
}

func TestListShipments_CalculaTotalDePaginas(t *testing.T) {
	cases := []struct {
		total     int64
		pageSize  int
		wantPages int
	}{
		{total: 0, pageSize: 10, wantPages: 0},
		{total: 1, pageSize: 10, wantPages: 1},
		{total: 10, pageSize: 10, wantPages: 1},
		{total: 11, pageSize: 10, wantPages: 2},
		{total: 95, pageSize: 10, wantPages: 10},
		{total: 100, pageSize: 100, wantPages: 1},
	}

	for _, tc := range cases {
		repo := &mocks.RepositoryMock{
			ListShipmentsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error) {
				return []domain.Shipment{}, tc.total, nil
			},
		}
		uc := newUseCaseForCost(repo, nil)

		resp, err := uc.ListShipments(context.Background(), 1, tc.pageSize, nil)

		require.NoError(t, err)
		assert.Equal(t, tc.wantPages, resp.TotalPages)
		assert.Equal(t, tc.total, resp.Total)
	}
}

func TestListShipments_MapeaCadaEnvio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListShipmentsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error) {
			return []domain.Shipment{
				{ID: 1, Status: "pending"},
				{ID: 2, Status: "delivered"},
			}, 2, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.ListShipments(context.Background(), 1, 10, nil)

	require.NoError(t, err)
	require.Len(t, resp.Data, 2)
	assert.Equal(t, uint(1), resp.Data[0].ID)
	assert.Equal(t, uint(2), resp.Data[1].ID)
	assert.Equal(t, "delivered", resp.Data[1].Status)
}

func TestListShipments_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("query failed")
	repo := &mocks.RepositoryMock{
		ListShipmentsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.ListShipments(context.Background(), 1, 10, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestUpdateShipment_IDCero_RetornaError(t *testing.T) {
	uc := newUseCaseForCost(&mocks.RepositoryMock{}, nil)

	resp, err := uc.UpdateShipment(context.Background(), 0, &domain.UpdateShipmentRequest{})

	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateShipment_NoExiste_RetornaNotFound(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) { return nil, nil },
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{})

	assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
	assert.Nil(t, resp)
}

func TestUpdateShipment_TrackingDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, OrderID: s("order-uuid"), TrackingNumber: s("viejo")}, nil
		},
		ShipmentExistsFn: func(ctx context.Context, orderID, trackingNumber string) (bool, error) {
			return true, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{
		TrackingNumber: s("034057376067"),
	})

	assert.ErrorIs(t, err, domain.ErrShipmentAlreadyExists)
	assert.Nil(t, resp)
}

func TestUpdateShipment_CambioDeStatus_SincronizaLaOrden(t *testing.T) {
	var syncedOrder, syncedStatus string
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, OrderID: s("order-uuid"), Status: "pending"}, nil
		},
		UpdateOrderStatusByOrderIDFn: func(ctx context.Context, orderID, status string) error {
			syncedOrder, syncedStatus = orderID, status
			return nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error { return nil },
	}
	uc := newUseCaseForCost(repo, nil)

	status := "delivered"
	_, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{Status: &status})

	require.NoError(t, err)
	assert.Equal(t, "order-uuid", syncedOrder)
	assert.Equal(t, "delivered", syncedStatus)
}

func TestUpdateShipment_StatusIgual_NoSincronizaLaOrden(t *testing.T) {
	llamadas := 0
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, OrderID: s("order-uuid"), Status: "pending"}, nil
		},
		UpdateOrderStatusByOrderIDFn: func(ctx context.Context, orderID, status string) error {
			llamadas++
			return nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error { return nil },
	}
	uc := newUseCaseForCost(repo, nil)

	status := "pending"
	_, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{Status: &status})

	require.NoError(t, err)
	assert.Equal(t, 0, llamadas)
}

func TestChunkRows(t *testing.T) {
	rows := make([]domain.SyncShipmentRow, 7)

	cases := []struct {
		size      int
		wantParts int
		wantLast  int
	}{
		{size: 3, wantParts: 3, wantLast: 1},
		{size: 7, wantParts: 1, wantLast: 7},
		{size: 10, wantParts: 1, wantLast: 7},
		{size: 1, wantParts: 7, wantLast: 1},
	}

	for _, tc := range cases {
		got := chunkRows(rows, tc.size)
		require.Len(t, got, tc.wantParts)
		assert.Len(t, got[len(got)-1], tc.wantLast)

		total := 0
		for _, part := range got {
			total += len(part)
		}
		assert.Equal(t, len(rows), total)
	}
}

func TestChunkRows_SinFilas_RetornaVacio(t *testing.T) {
	assert.Empty(t, chunkRows([]domain.SyncShipmentRow{}, 5))
	assert.Empty(t, chunkRows(nil, 5))
}
