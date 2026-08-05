package usecaseshipment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func u(v uint) *uint {
	return &v
}

func b(v bool) *bool {
	return &v
}

func TestUpdateShipment_ActualizaTodosLosCamposOpcionales(t *testing.T) {
	var saved *domain.Shipment
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, Status: "pending"}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			saved = shipment
			return nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	entrega := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	nombreBodega := "Bodega Norte"
	nombreConductor := "Pedro"

	resp, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{
		TrackingURL:       s("https://track/1"),
		Carrier:           s("Servientrega"),
		CarrierCode:       s("servientrega"),
		GuideID:           s("G-1"),
		GuideURL:          s("https://s3/g1.pdf"),
		ShippedAt:         &entrega,
		DeliveredAt:       &entrega,
		ShippingAddressID: u(3),
		ShippingCost:      f64(10000),
		InsuranceCost:     f64(500),
		TotalCost:         f64(11000),
		Weight:            f64(1.5),
		Height:            f64(10),
		Width:             f64(20),
		Length:            f64(30),
		WarehouseID:       u(2),
		WarehouseName:     &nombreBodega,
		DriverID:          u(9),
		DriverName:        &nombreConductor,
		IsLastMile:        b(true),
		EstimatedDelivery: &entrega,
		DeliveryNotes:     s("dejar en porteria"),
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, saved)

	assert.Equal(t, "https://track/1", *saved.TrackingURL)
	assert.Equal(t, "Servientrega", *saved.Carrier)
	assert.Equal(t, "servientrega", *saved.CarrierCode)
	assert.Equal(t, "G-1", *saved.GuideID)
	assert.Equal(t, "https://s3/g1.pdf", *saved.GuideURL)
	assert.Equal(t, uint(3), *saved.ShippingAddressID)
	assert.InDelta(t, 10000.0, *saved.ShippingCost, 0.001)
	assert.InDelta(t, 500.0, *saved.InsuranceCost, 0.001)
	assert.InDelta(t, 11000.0, *saved.TotalCost, 0.001)
	assert.InDelta(t, 1.5, *saved.Weight, 0.001)
	assert.InDelta(t, 10.0, *saved.Height, 0.001)
	assert.InDelta(t, 20.0, *saved.Width, 0.001)
	assert.InDelta(t, 30.0, *saved.Length, 0.001)
	assert.Equal(t, uint(2), *saved.WarehouseID)
	assert.Equal(t, "Bodega Norte", saved.WarehouseName)
	assert.Equal(t, uint(9), *saved.DriverID)
	assert.Equal(t, "Pedro", saved.DriverName)
	assert.True(t, saved.IsLastMile)
	assert.Equal(t, "dejar en porteria", *saved.DeliveryNotes)
}

func TestUpdateShipment_CambioDeCodCarrierFee_RecalculaMargenes(t *testing.T) {
	var saved *domain.Shipment
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:                   id,
				Status:               "pending",
				CarrierCost:          f64(999),
				AppliedMargin:        f64(111),
				CodProbabilityMargin: f64(222),
			}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error {
			saved = shipment
			return nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	_, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{
		CodCarrierFee: f64(4000),
	})

	require.NoError(t, err)
	require.NotNil(t, saved)
	assert.InDelta(t, 4000.0, *saved.CodCarrierFee, 0.001)
	assert.Nil(t, saved.CarrierCost)
	assert.Nil(t, saved.AppliedMargin)
	assert.Nil(t, saved.CodProbabilityMargin)
}

func TestUpdateShipment_SinCambios_NoRompe(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, Status: "pending", ClientName: "Ana"}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error { return nil },
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{})

	require.NoError(t, err)
	assert.Equal(t, "Ana", resp.ClientName)
	assert.Equal(t, "pending", resp.Status)
}

func TestUpdateShipment_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("update failed")
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id}, nil
		},
		UpdateShipmentFn: func(ctx context.Context, shipment *domain.Shipment) error { return dbErr },
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestUpdateShipment_FallaLeerEnvio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) { return nil, dbErr },
	}
	uc := newUseCaseForCost(repo, nil)

	resp, err := uc.UpdateShipment(context.Background(), 7, &domain.UpdateShipmentRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestDeleteShipment_IDCero_RetornaError(t *testing.T) {
	uc := newUseCaseForCost(&mocks.RepositoryMock{}, nil)

	err := uc.DeleteShipment(context.Background(), 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "shipment ID is required")
}

func TestDeleteShipment_NoExiste_RetornaNotFound(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) { return nil, nil },
	}
	uc := newUseCaseForCost(repo, nil)

	err := uc.DeleteShipment(context.Background(), 7)

	assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
}

func TestDeleteShipment_Exitoso(t *testing.T) {
	borrado := uint(0)
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id}, nil
		},
		DeleteShipmentFn: func(ctx context.Context, id uint) error {
			borrado = id
			return nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	err := uc.DeleteShipment(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, uint(7), borrado)
}

func TestDeleteShipment_FallaBorrar_PropagaError(t *testing.T) {
	dbErr := errors.New("delete failed")
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id}, nil
		},
		DeleteShipmentFn: func(ctx context.Context, id uint) error { return dbErr },
	}
	uc := newUseCaseForCost(repo, nil)

	err := uc.DeleteShipment(context.Background(), 7)

	assert.ErrorIs(t, err, dbErr)
}

func TestDeleteShipment_FallaLeer_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetShipmentByIDFn: func(ctx context.Context, id uint) (*domain.Shipment, error) { return nil, dbErr },
	}
	uc := newUseCaseForCost(repo, nil)

	err := uc.DeleteShipment(context.Background(), 7)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetShipmentsByOrderID_OrderIDVacio_RetornaError(t *testing.T) {
	uc := newUseCaseForCost(&mocks.RepositoryMock{}, nil)

	got, err := uc.GetShipmentsByOrderID(context.Background(), "")

	assert.ErrorIs(t, err, domain.ErrOrderIDRequired)
	assert.Nil(t, got)
}

func TestGetShipmentsByOrderID_Exitoso(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentsByOrderIDFn: func(ctx context.Context, orderID string) ([]domain.Shipment, error) {
			return []domain.Shipment{{ID: 1}, {ID: 2}}, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	got, err := uc.GetShipmentsByOrderID(context.Background(), "order-uuid")

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestGetShipmentsByOrderID_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetShipmentsByOrderIDFn: func(ctx context.Context, orderID string) ([]domain.Shipment, error) {
			return nil, dbErr
		},
	}
	uc := newUseCaseForCost(repo, nil)

	got, err := uc.GetShipmentsByOrderID(context.Background(), "order-uuid")

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetShipmentByTrackingNumber_Vacio_RetornaError(t *testing.T) {
	uc := newUseCaseForCost(&mocks.RepositoryMock{}, nil)

	got, err := uc.GetShipmentByTrackingNumber(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "tracking number is required")
}

func TestGetShipmentByTrackingNumber_NoExiste_RetornaNotFound(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByTrackingNumberFn: func(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
			return nil, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	got, err := uc.GetShipmentByTrackingNumber(context.Background(), "034057376067")

	assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
	assert.Nil(t, got)
}

func TestGetShipmentByTrackingNumber_Exitoso(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetShipmentByTrackingNumberFn: func(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
			return &domain.Shipment{ID: 34678, TrackingNumber: &trackingNumber}, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	got, err := uc.GetShipmentByTrackingNumber(context.Background(), "034057376067")

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint(34678), got.ID)
}

func TestGetShipmentByTrackingNumber_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetShipmentByTrackingNumberFn: func(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
			return nil, dbErr
		},
	}
	uc := newUseCaseForCost(repo, nil)

	got, err := uc.GetShipmentByTrackingNumber(context.Background(), "034057376067")

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestResolveGeozoneBestEffort_NoLlamaCuandoFaltanDatos(t *testing.T) {
	llamadas := 0
	repo := &mocks.RepositoryMock{
		GetOrderBusinessIDFn: func(ctx context.Context, orderUUID string) (uint, error) {
			llamadas++
			return 36, nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	uc.resolveGeozoneBestEffort(context.Background(), nil)
	uc.resolveGeozoneBestEffort(context.Background(), &domain.Shipment{ID: 0, OrderID: s("x")})
	uc.resolveGeozoneBestEffort(context.Background(), &domain.Shipment{ID: 1, OrderID: nil})
	uc.resolveGeozoneBestEffort(context.Background(), &domain.Shipment{ID: 1, OrderID: s("")})

	assert.Equal(t, 0, llamadas)
}

func TestResolveGeozoneBestEffort_ResuelveConDatosCompletos(t *testing.T) {
	var resueltoShipment, resueltoBusiness uint
	repo := &mocks.RepositoryMock{
		GetOrderBusinessIDFn: func(ctx context.Context, orderUUID string) (uint, error) { return 36, nil },
		ResolveShipmentGeozoneFn: func(ctx context.Context, shipmentID, businessID uint) error {
			resueltoShipment, resueltoBusiness = shipmentID, businessID
			return nil
		},
	}
	uc := newUseCaseForCost(repo, nil)

	uc.resolveGeozoneBestEffort(context.Background(), &domain.Shipment{ID: 55, OrderID: s("order-uuid")})

	assert.Equal(t, uint(55), resueltoShipment)
	assert.Equal(t, uint(36), resueltoBusiness)
}
