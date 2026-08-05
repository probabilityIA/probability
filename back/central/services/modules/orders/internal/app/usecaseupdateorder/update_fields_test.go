package usecaseupdateorder

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sPtr(v string) *string {
	return &v
}

func TestUpdatePaymentFields_SinPagos_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}

func TestUpdatePaymentFields_CambiaMetodoDePago(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{PaymentMethodID: 1}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{PaymentMethodID: 5}},
	})

	assert.True(t, changed)
	assert.Equal(t, uint(5), order.PaymentMethodID)
}

func TestUpdatePaymentFields_MismoMetodo_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{PaymentMethodID: 5}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{PaymentMethodID: 5}},
	})

	assert.False(t, changed)
}

func TestUpdatePaymentFields_MetodoCero_SeIgnora(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{PaymentMethodID: 7}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{PaymentMethodID: 0}},
	})

	assert.False(t, changed)
	assert.Equal(t, uint(7), order.PaymentMethodID)
}

func TestUpdatePaymentFields_PagoCompletado_MarcaPagada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{IsPaid: false}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{Status: "completed"}},
	})

	assert.True(t, changed)
	assert.True(t, order.IsPaid)
}

func TestUpdatePaymentFields_YaPagada_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{IsPaid: true}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{Status: "completed"}},
	})

	assert.False(t, changed)
}

func TestUpdatePaymentFields_NuncaDespaga(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{IsPaid: true}

	uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{Status: "pending"}},
	})

	assert.True(t, order.IsPaid, "un pago pendiente no revierte una orden ya pagada")
}

func TestUpdatePaymentFields_ActualizaFechaDePago(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	nueva := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	order := &entities.ProbabilityOrder{}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{PaidAt: &nueva}},
	})

	assert.True(t, changed)
	require.NotNil(t, order.PaidAt)
	assert.True(t, order.PaidAt.Equal(nueva))
}

func TestUpdatePaymentFields_MismaFecha_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	fecha := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	order := &entities.ProbabilityOrder{PaidAt: &fecha}

	changed := uc.updatePaymentFields(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{PaidAt: &fecha}},
	})

	assert.False(t, changed)
}

func TestUpdateFulfillmentFields_SinEnvios_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	changed := uc.updateFulfillmentFields(&entities.ProbabilityOrder{}, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}

func TestUpdateFulfillmentFields_ActualizaBodegaYConductor(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateFulfillmentFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{
			WarehouseID:   uintPtr(2),
			WarehouseName: "Bodega Norte",
			DriverID:      uintPtr(9),
			DriverName:    "Pedro",
			IsLastMile:    true,
		}},
	})

	assert.True(t, changed)
	require.NotNil(t, order.WarehouseID)
	assert.Equal(t, uint(2), *order.WarehouseID)
	assert.Equal(t, "Bodega Norte", order.WarehouseName)
	require.NotNil(t, order.DriverID)
	assert.Equal(t, uint(9), *order.DriverID)
	assert.Equal(t, "Pedro", order.DriverName)
	assert.True(t, order.IsLastMile)
}

func TestUpdateFulfillmentFields_MismosValores_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		WarehouseID:   uintPtr(2),
		WarehouseName: "Bodega Norte",
		DriverID:      uintPtr(9),
		DriverName:    "Pedro",
		IsLastMile:    true,
	}

	changed := uc.updateFulfillmentFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{
			WarehouseID:   uintPtr(2),
			WarehouseName: "Bodega Norte",
			DriverID:      uintPtr(9),
			DriverName:    "Pedro",
			IsLastMile:    true,
		}},
	})

	assert.False(t, changed)
}

func TestUpdateFulfillmentFields_NombresVacios_NoBorranLosExistentes(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		WarehouseName: "Bodega Norte",
		DriverName:    "Pedro",
	}

	uc.updateFulfillmentFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{}},
	})

	assert.Equal(t, "Bodega Norte", order.WarehouseName)
	assert.Equal(t, "Pedro", order.DriverName)
}

func TestUpdateTrackingFields_SinEnvios_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	changed := uc.updateTrackingFields(&entities.ProbabilityOrder{}, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
}

func TestUpdateTrackingFields_ActualizaGuiaYTracking(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	entregado := time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC)
	despachado := time.Date(2026, 8, 2, 9, 0, 0, 0, time.UTC)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateTrackingFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{
			TrackingNumber: sPtr("034057376067"),
			TrackingURL:    sPtr("https://track/1"),
			GuideID:        sPtr("G-1"),
			GuideURL:       sPtr("https://s3/g1.pdf"),
			DeliveredAt:    &entregado,
			ShippedAt:      &despachado,
		}},
	})

	assert.True(t, changed)
	assert.Equal(t, "034057376067", *order.TrackingNumber)
	assert.Equal(t, "https://track/1", *order.TrackingLink)
	assert.Equal(t, "G-1", *order.GuideID)
	assert.Equal(t, "https://s3/g1.pdf", *order.GuideLink)
	assert.True(t, order.DeliveredAt.Equal(entregado))
	assert.True(t, order.DeliveryDate.Equal(despachado))
}

func TestUpdateTrackingFields_MismosValores_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	entregado := time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC)
	order := &entities.ProbabilityOrder{
		TrackingNumber: sPtr("034057376067"),
		TrackingLink:   sPtr("https://track/1"),
		GuideID:        sPtr("G-1"),
		GuideLink:      sPtr("https://s3/g1.pdf"),
		DeliveredAt:    &entregado,
	}

	changed := uc.updateTrackingFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{
			TrackingNumber: sPtr("034057376067"),
			TrackingURL:    sPtr("https://track/1"),
			GuideID:        sPtr("G-1"),
			GuideURL:       sPtr("https://s3/g1.pdf"),
			DeliveredAt:    &entregado,
		}},
	})

	assert.False(t, changed)
}

func TestUpdateTrackingFields_CamposNulos_NoBorranLosExistentes(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		TrackingNumber: sPtr("034057376067"),
		GuideLink:      sPtr("https://s3/g1.pdf"),
	}

	changed := uc.updateTrackingFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{}},
	})

	assert.False(t, changed)
	assert.Equal(t, "034057376067", *order.TrackingNumber)
	assert.Equal(t, "https://s3/g1.pdf", *order.GuideLink)
}

func TestUpdateTrackingFields_SoloCambiaLoQueLlega(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{
		TrackingNumber: sPtr("viejo"),
		GuideID:        sPtr("G-viejo"),
	}

	changed := uc.updateTrackingFields(order, &dtos.ProbabilityOrderDTO{
		Shipments: []dtos.ProbabilityShipmentDTO{{TrackingNumber: sPtr("nuevo")}},
	})

	assert.True(t, changed)
	assert.Equal(t, "nuevo", *order.TrackingNumber)
	assert.Equal(t, "G-viejo", *order.GuideID)
}
