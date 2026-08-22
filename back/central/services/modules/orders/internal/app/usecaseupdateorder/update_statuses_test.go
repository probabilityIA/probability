package usecaseupdateorder

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateOrderStatus_SinDatos_NoCambiaNada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Status: "pending"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
	assert.Equal(t, "pending", order.Status)
}

func TestUpdateOrderStatus_MismoEstado_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Status: "confirmed"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{Status: "confirmed"})

	assert.False(t, changed)
}

func TestUpdateOrderStatus_CambioValido_SeAplica(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Status: "pending"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{Status: "confirmed"})

	assert.True(t, changed)
	assert.Equal(t, "confirmed", order.Status)
}

func TestUpdateOrderStatus_SaltoDeEstado_SeAceptaPorSerFuenteExterna(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{ID: "order-uuid", Status: "pending"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status:          "delivered",
		IntegrationType: "shopify",
	})

	assert.True(t, changed)
	assert.Equal(t, "delivered", order.Status,
		"las integraciones son fuente de verdad y pueden saltar estados")
}

func TestUpdateOrderStatus_Cancelacion_NoSeConsideraSalto(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Status: "pending"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status: string(entities.OrderStatusCancelled),
	})

	assert.True(t, changed)
	assert.Equal(t, string(entities.OrderStatusCancelled), order.Status)
}

func TestUpdateOrderStatus_CambiaOriginalStatus_RemapeaStatusID(t *testing.T) {
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			return uintPtr(9), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{Status: "confirmed", OriginalStatus: "pending"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		OriginalStatus:  "paid",
		IntegrationType: "shopify",
	})

	assert.True(t, changed)
	assert.Equal(t, "paid", order.OriginalStatus)
	require.NotNil(t, order.StatusID)
	assert.Equal(t, uint(9), *order.StatusID)
}

func TestUpdateOrderStatus_MismoOriginalStatus_NoRemapea(t *testing.T) {
	consultado := false
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			consultado = true
			return uintPtr(9), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{OriginalStatus: "paid"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		OriginalStatus:  "paid",
		IntegrationType: "shopify",
	})

	assert.False(t, changed)
	assert.False(t, consultado)
}

func TestUpdateOrderStatus_MapeoNulo_LimpiaStatusID(t *testing.T) {
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			return nil, nil
		},
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return nil, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{StatusID: uintPtr(5), OriginalStatus: "pending"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		OriginalStatus:  "desconocido",
		IntegrationType: "shopify",
	})

	assert.True(t, changed)
	assert.Nil(t, order.StatusID)
}

func TestUpdatePaymentStatus_MapeaYSincronizaIsPaid(t *testing.T) {
	paidID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return &paidID, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updatePaymentStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		PaymentStatusID: &paidID,
	})

	assert.True(t, changed)
	require.NotNil(t, order.PaymentStatusID)
	assert.Equal(t, paidID, *order.PaymentStatusID)
	assert.True(t, order.IsPaid)
}

func TestUpdatePaymentStatus_MismoEstado_NoMarcaCambio(t *testing.T) {
	otroID := uint(1)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(3), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{PaymentStatusID: &otroID}

	changed := uc.updatePaymentStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		PaymentStatusID: &otroID,
	})

	assert.False(t, changed)
	assert.False(t, order.IsPaid)
}

func TestUpdatePaymentStatus_SinMapeo_MantieneElEstadoActual(t *testing.T) {
	actual := uint(1)
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{PaymentStatusID: &actual}

	changed := uc.updatePaymentStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
	require.NotNil(t, order.PaymentStatusID)
	assert.Equal(t, actual, *order.PaymentStatusID)
}

func TestUpdatePaymentStatus_SoloSincronizacionDeIsPaid_MarcaCambio(t *testing.T) {
	paidID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return &paidID, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{PaymentStatusID: &paidID, IsPaid: false}

	changed := uc.updatePaymentStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{})

	assert.True(t, changed, "sincronizar is_paid con un estado ya pagado cuenta como cambio")
	assert.True(t, order.IsPaid)
}

func TestUpdateFulfillmentStatus_MapeaNuevoEstado(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{}

	changed := uc.updateFulfillmentStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		FulfillmentStatusID: uintPtr(2),
	})

	assert.True(t, changed)
	require.NotNil(t, order.FulfillmentStatusID)
	assert.Equal(t, uint(2), *order.FulfillmentStatusID)
}

func TestUpdateFulfillmentStatus_MismoEstado_NoMarcaCambio(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{FulfillmentStatusID: uintPtr(2)}

	changed := uc.updateFulfillmentStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		FulfillmentStatusID: uintPtr(2),
	})

	assert.False(t, changed)
}

func TestUpdateFulfillmentStatus_SinMapeo_MantieneElActual(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{FulfillmentStatusID: uintPtr(2)}

	changed := uc.updateFulfillmentStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{})

	assert.False(t, changed)
	require.NotNil(t, order.FulfillmentStatusID)
	assert.Equal(t, uint(2), *order.FulfillmentStatusID)
}

func TestUpdateOrderStatuses_CombinaLosTres(t *testing.T) {
	paidID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return &paidID, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{Status: "pending"}

	changed := uc.updateOrderStatuses(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status:              "confirmed",
		PaymentStatusID:     &paidID,
		FulfillmentStatusID: uintPtr(2),
	})

	assert.True(t, changed)
	assert.Equal(t, "confirmed", order.Status)
	assert.Equal(t, paidID, *order.PaymentStatusID)
	assert.Equal(t, uint(2), *order.FulfillmentStatusID)
	assert.True(t, order.IsPaid)
}

func TestUpdateOrderStatuses_SinCambios_RetornaFalse(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{Status: "confirmed"}

	changed := uc.updateOrderStatuses(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status: "confirmed",
	})

	assert.False(t, changed)
}

func TestUpdateOrderStatus_OrdenEnBodega_IgnoraElEstadoDelCanal(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{ID: "order-uuid", Status: "picking"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status:          "paid",
		IntegrationType: "tiendanube",
	})

	assert.False(t, changed)
	assert.Equal(t, "picking", order.Status,
		"el canal no puede devolver a la orden a un estado previo a la operacion")
}

func TestUpdateOrderStatus_OrdenEnTransito_IgnoraElEstadoDelCanal(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{ID: "order-uuid", Status: "in_transit"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status:          "processing",
		IntegrationType: "shopify",
	})

	assert.False(t, changed)
	assert.Equal(t, "in_transit", order.Status)
}

func TestUpdateOrderStatus_CancelacionDelCanal_PisaAunEnBodega(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{ID: "order-uuid", Status: "picking"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status:          string(entities.OrderStatusCancelled),
		IntegrationType: "tiendanube",
	})

	assert.True(t, changed)
	assert.Equal(t, string(entities.OrderStatusCancelled), order.Status)
}

func TestUpdateOrderStatus_ReembolsoDelCanal_PisaEntregada(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)
	order := &entities.ProbabilityOrder{ID: "order-uuid", Status: "delivered"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status:          string(entities.OrderStatusRefunded),
		IntegrationType: "shopify",
	})

	assert.True(t, changed)
	assert.Equal(t, string(entities.OrderStatusRefunded), order.Status)
}

func TestUpdateOrderStatus_OrdenEnBodega_RegistraOriginalStatusSinTocarStatusID(t *testing.T) {
	mapeado := uint(2)
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			return &mapeado, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	actual := uint(12)
	order := &entities.ProbabilityOrder{ID: "order-uuid", Status: "picking", StatusID: &actual, OriginalStatus: "open"}

	changed := uc.updateOrderStatus(context.Background(), order, &dtos.ProbabilityOrderDTO{
		Status:          "paid",
		OriginalStatus:  "closed",
		IntegrationType: "tiendanube",
	})

	assert.True(t, changed, "el original_status del canal si se guarda")
	assert.Equal(t, "closed", order.OriginalStatus)
	assert.Equal(t, "picking", order.Status)
	require.NotNil(t, order.StatusID)
	assert.Equal(t, uint(12), *order.StatusID, "el status_id operativo no se remapea")
}
