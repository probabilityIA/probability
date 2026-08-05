package usecaseupdateorder

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uintPtr(v uint) *uint {
	return &v
}

func TestUpdate_MapOrderStatusID_MapeoPorIntegracion(t *testing.T) {
	var pedidoTipo uint
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			pedidoTipo = integrationTypeID
			return uintPtr(7), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "open",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(7), *got)
	assert.Equal(t, uint(1), pedidoTipo)
}

func TestUpdate_MapOrderStatusID_CaeEnOriginalStatus(t *testing.T) {
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			if originalStatus == "paid" {
				return uintPtr(9), nil
			}
			return nil, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "desconocido",
		OriginalStatus:  "paid",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(9), *got)
}

func TestUpdate_MapOrderStatusID_FallbackPorCodigo(t *testing.T) {
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			return nil, nil
		},
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(4), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "confirmed",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(4), *got)
}

func TestUpdate_MapOrderStatusID_ErrorEnMapeo_ContinuaAlFallback(t *testing.T) {
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			return nil, errors.New("db down")
		},
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(5), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "open",
		OriginalStatus:  "paid",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(5), *got)
}

func TestUpdate_MapOrderStatusID_SinDatos_RetornaNil(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	assert.Nil(t, uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{}))
}

func TestUpdate_MapPaymentStatusID_IDExplicito(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		PaymentStatusID: uintPtr(3),
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(3), *got)
}

func TestUpdate_MapPaymentStatusID_Shopify_DerivaDeFinancialStatus(t *testing.T) {
	var codigo string
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigo = code
			return uintPtr(8), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		PaymentDetails:  []byte(`{"financial_status":"partially_paid"}`),
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(8), *got)
	assert.Equal(t, "partially_paid", codigo)
}

func TestUpdate_MapPaymentStatusID_Shopify_PayloadInvalido_RetornaNil(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	cases := [][]byte{
		[]byte(`{roto`),
		[]byte(`{}`),
		[]byte(`{"financial_status":""}`),
		[]byte(`{"financial_status":42}`),
	}

	for _, payload := range cases {
		got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
			IntegrationType: "shopify",
			PaymentDetails:  payload,
		})
		assert.Nil(t, got)
	}
}

func TestUpdate_MapPaymentStatusID_CanalNoShopify_RetornaNil(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "woocommerce",
		PaymentDetails:  []byte(`{"financial_status":"paid"}`),
	})

	assert.Nil(t, got)
}

func TestUpdate_MapFulfillmentStatusID_IDExplicito(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	got := uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		FulfillmentStatusID: uintPtr(2),
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(2), *got)
}

func TestUpdate_MapFulfillmentStatusID_Shopify_DerivaDelDetalle(t *testing.T) {
	var codigo string
	repo := &mockRepository{
		GetFulfillmentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigo = code
			return uintPtr(6), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	got := uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType:    "shopify",
		FulfillmentDetails: []byte(`{"fulfillment_status":"shipped"}`),
	})

	require.NotNil(t, got)
	assert.Equal(t, "shipped", codigo)
}

func TestUpdate_MapFulfillmentStatusID_Shopify_InvalidoCaeEnUnfulfilled(t *testing.T) {
	codigos := []string{}
	repo := &mockRepository{
		GetFulfillmentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigos = append(codigos, code)
			return uintPtr(1), nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	payloads := [][]byte{
		[]byte(`{roto`),
		[]byte(`{"fulfillment_status":null}`),
		[]byte(`{"fulfillment_status":""}`),
	}

	for _, payload := range payloads {
		uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
			IntegrationType:    "shopify",
			FulfillmentDetails: payload,
		})
	}

	assert.Equal(t, []string{"unfulfilled", "unfulfilled", "unfulfilled"}, codigos)
}

func TestUpdate_MapFulfillmentStatusID_SinDetalles_RetornaNil(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	assert.Nil(t, uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
	}))
}

func TestUpdate_GetFulfillmentStatusIDByCode_ErrorORetornoNil(t *testing.T) {
	cases := []struct {
		name string
		fn   func(ctx context.Context, code string) (*uint, error)
	}{
		{name: "error", fn: func(ctx context.Context, code string) (*uint, error) { return nil, errors.New("db") }},
		{name: "no encontrado", fn: func(ctx context.Context, code string) (*uint, error) { return nil, nil }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc := newTestUpdateUseCase(&mockRepository{GetFulfillmentStatusIDByCodeFn: tc.fn}, nil, nil)

			assert.Nil(t, uc.getFulfillmentStatusIDByCode(context.Background(), "fulfilled"))
		})
	}
}

func TestUpdate_SyncIsPaidFromPaymentStatus_MarcaPagada(t *testing.T) {
	paidID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return &paidID, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)
	order := &entities.ProbabilityOrder{}

	uc.syncIsPaidFromPaymentStatus(context.Background(), order, &paidID)

	assert.True(t, order.IsPaid)
	require.NotNil(t, order.PaidAt)
}

func TestUpdate_SyncIsPaidFromPaymentStatus_ConservaFechaExistente(t *testing.T) {
	paidID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return &paidID, nil
		},
	}
	uc := newTestUpdateUseCase(repo, nil, nil)

	original := time.Date(2026, 3, 1, 8, 0, 0, 0, time.UTC)
	order := &entities.ProbabilityOrder{PaidAt: &original}

	uc.syncIsPaidFromPaymentStatus(context.Background(), order, &paidID)

	assert.True(t, order.PaidAt.Equal(original))
}

func TestUpdate_SyncIsPaidFromPaymentStatus_NoMarcaCuando(t *testing.T) {
	paidID := uint(3)
	otroID := uint(1)

	cases := []struct {
		name     string
		statusID *uint
		fn       func(ctx context.Context, code string) (*uint, error)
	}{
		{
			name:     "sin estado de pago",
			statusID: nil,
			fn:       func(ctx context.Context, code string) (*uint, error) { return &paidID, nil },
		},
		{
			name:     "estado distinto de pagado",
			statusID: &otroID,
			fn:       func(ctx context.Context, code string) (*uint, error) { return &paidID, nil },
		},
		{
			name:     "falla la consulta",
			statusID: &paidID,
			fn:       func(ctx context.Context, code string) (*uint, error) { return nil, errors.New("db down") },
		},
		{
			name:     "catalogo sin paid",
			statusID: &paidID,
			fn:       func(ctx context.Context, code string) (*uint, error) { return nil, nil },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc := newTestUpdateUseCase(&mockRepository{GetPaymentStatusIDByCodeFn: tc.fn}, nil, nil)
			order := &entities.ProbabilityOrder{}

			uc.syncIsPaidFromPaymentStatus(context.Background(), order, tc.statusID)

			assert.False(t, order.IsPaid)
			assert.Nil(t, order.PaidAt)
		})
	}
}
