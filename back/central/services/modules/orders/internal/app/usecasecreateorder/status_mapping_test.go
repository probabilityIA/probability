package usecasecreateorder

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uintPtr(v uint) *uint {
	return &v
}

func TestMapOrderStatusID_MapeoPorIntegracionConStatus(t *testing.T) {
	var pedidoTipo uint
	var pedidoStatus string
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			pedidoTipo, pedidoStatus = integrationTypeID, originalStatus
			return uintPtr(7), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "open",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(7), *got)
	assert.Equal(t, uint(1), pedidoTipo)
	assert.Equal(t, "open", pedidoStatus)
}

func TestMapOrderStatusID_CaeEnOriginalStatusSiStatusNoMapea(t *testing.T) {
	consultas := []string{}
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			consultas = append(consultas, originalStatus)
			if originalStatus == "paid" {
				return uintPtr(9), nil
			}
			return nil, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "desconocido",
		OriginalStatus:  "paid",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(9), *got)
	assert.Equal(t, []string{"desconocido", "paid"}, consultas)
}

func TestMapOrderStatusID_CaeEnCodigoDirectoSinMapeoDeIntegracion(t *testing.T) {
	var codigoBuscado string
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			return nil, nil
		},
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigoBuscado = code
			return uintPtr(4), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "confirmed",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(4), *got)
	assert.Equal(t, "confirmed", codigoBuscado)
}

func TestMapOrderStatusID_OrdenManualSinIntegracion_UsaCodigoDirecto(t *testing.T) {
	consultoMapeo := false
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			consultoMapeo = true
			return nil, nil
		},
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(2), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{Status: "pending"})

	require.NotNil(t, got)
	assert.Equal(t, uint(2), *got)
	assert.False(t, consultoMapeo, "sin integration_type no se consulta el mapeo por integracion")
}

func TestMapOrderStatusID_TipoDeIntegracionDesconocido_SaltaAlCodigoDirecto(t *testing.T) {
	consultoMapeo := false
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			consultoMapeo = true
			return nil, nil
		},
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(2), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "canal-que-no-existe",
		Status:          "pending",
	})

	require.NotNil(t, got)
	assert.False(t, consultoMapeo)
}

func TestMapOrderStatusID_ErrorEnOriginalStatus_SeRegistraYContinua(t *testing.T) {
	repo := &mockRepository{
		GetOrderStatusIDByIntegrationTypeAndOriginalStatusFn: func(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
			return nil, errors.New("db down")
		},
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(5), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		Status:          "open",
		OriginalStatus:  "paid",
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(5), *got)
}

func TestMapOrderStatusID_SinNadaQueMapear_RetornaNil(t *testing.T) {
	uc := newTestCreateUseCase(&mockRepository{}, nil, nil, nil)

	got := uc.mapOrderStatusID(context.Background(), &dtos.ProbabilityOrderDTO{})

	assert.Nil(t, got)
}

func TestMapPaymentStatusID_DTOConIDExplicito_SeUsaDirecto(t *testing.T) {
	consulto := false
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			consulto = true
			return uintPtr(99), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		PaymentStatusID: uintPtr(3),
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(3), *got)
	assert.False(t, consulto)
}

func TestMapPaymentStatusID_IDCeroSeIgnora(t *testing.T) {
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return nil, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		PaymentStatusID: uintPtr(0),
	})

	assert.Nil(t, got)
}

func TestMapPaymentStatusID_PagoCompletado_MapeaAPaid(t *testing.T) {
	var codigo string
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigo = code
			return uintPtr(3), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{Status: "completed"}},
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(3), *got)
	assert.Equal(t, "paid", codigo)
}

func TestMapPaymentStatusID_PagoPendiente_NoMapeaAPaid(t *testing.T) {
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(3), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{Status: "pending"}},
	})

	assert.Nil(t, got)
}

func TestMapPaymentStatusID_PagoCompletadoPeroFallaConsulta_RetornaNil(t *testing.T) {
	llamadas := 0
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			llamadas++
			return nil, errors.New("db down")
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		Payments: []dtos.ProbabilityPaymentDTO{{Status: "completed"}, {Status: "completed"}},
	})

	assert.Nil(t, got)
	assert.Equal(t, 1, llamadas, "tras el primer pago completado corta el recorrido")
}

func TestMapPaymentStatusID_Shopify_DerivaDeFinancialStatus(t *testing.T) {
	var codigo string
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigo = code
			return uintPtr(8), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		PaymentDetails:  []byte(`{"financial_status":"refunded"}`),
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(8), *got)
	assert.Equal(t, "refunded", codigo)
}

func TestMapPaymentStatusID_Shopify_JSONInvalido_RetornaNil(t *testing.T) {
	uc := newTestCreateUseCase(&mockRepository{}, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
		PaymentDetails:  []byte(`{no es json`),
	})

	assert.Nil(t, got)
}

func TestMapPaymentStatusID_Shopify_SinFinancialStatus_RetornaNil(t *testing.T) {
	uc := newTestCreateUseCase(&mockRepository{}, nil, nil, nil)

	cases := [][]byte{
		[]byte(`{}`),
		[]byte(`{"financial_status":""}`),
		[]byte(`{"financial_status":null}`),
		[]byte(`{"financial_status":123}`),
	}

	for _, payload := range cases {
		got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
			IntegrationType: "shopify",
			PaymentDetails:  payload,
		})
		assert.Nil(t, got)
	}
}

func TestMapPaymentStatusID_CanalNoShopify_NoLeePaymentDetails(t *testing.T) {
	uc := newTestCreateUseCase(&mockRepository{}, nil, nil, nil)

	got := uc.mapPaymentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "woocommerce",
		PaymentDetails:  []byte(`{"financial_status":"paid"}`),
	})

	assert.Nil(t, got)
}

func TestMapFulfillmentStatusID_DTOConIDExplicito_SeUsaDirecto(t *testing.T) {
	uc := newTestCreateUseCase(&mockRepository{}, nil, nil, nil)

	got := uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		FulfillmentStatusID: uintPtr(2),
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(2), *got)
}

func TestMapFulfillmentStatusID_Shopify_DerivaDelDetalle(t *testing.T) {
	var codigo string
	repo := &mockRepository{
		GetFulfillmentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigo = code
			return uintPtr(6), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType:    "shopify",
		FulfillmentDetails: []byte(`{"fulfillment_status":"fulfilled"}`),
	})

	require.NotNil(t, got)
	assert.Equal(t, uint(6), *got)
	assert.Equal(t, "fulfilled", codigo)
}

func TestMapFulfillmentStatusID_Shopify_JSONInvalido_CaeEnUnfulfilled(t *testing.T) {
	var codigo string
	repo := &mockRepository{
		GetFulfillmentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigo = code
			return uintPtr(1), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType:    "shopify",
		FulfillmentDetails: []byte(`{roto`),
	})

	require.NotNil(t, got)
	assert.Equal(t, "unfulfilled", codigo)
}

func TestMapFulfillmentStatusID_Shopify_StatusNulo_CaeEnUnfulfilled(t *testing.T) {
	var codigo string
	repo := &mockRepository{
		GetFulfillmentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			codigo = code
			return uintPtr(1), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType:    "shopify",
		FulfillmentDetails: []byte(`{"fulfillment_status":null}`),
	})

	require.NotNil(t, got)
	assert.Equal(t, "unfulfilled", codigo)
}

func TestMapFulfillmentStatusID_SinDetalles_RetornaNil(t *testing.T) {
	uc := newTestCreateUseCase(&mockRepository{}, nil, nil, nil)

	got := uc.mapFulfillmentStatusID(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType: "shopify",
	})

	assert.Nil(t, got)
}

func TestGetFulfillmentStatusIDByCode_ErrorORetornoNil_DevuelveNil(t *testing.T) {
	cases := []struct {
		name string
		fn   func(ctx context.Context, code string) (*uint, error)
	}{
		{name: "error", fn: func(ctx context.Context, code string) (*uint, error) { return nil, errors.New("db down") }},
		{name: "no encontrado", fn: func(ctx context.Context, code string) (*uint, error) { return nil, nil }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc := newTestCreateUseCase(&mockRepository{GetFulfillmentStatusIDByCodeFn: tc.fn}, nil, nil, nil)

			assert.Nil(t, uc.getFulfillmentStatusIDByCode(context.Background(), "fulfilled"))
		})
	}
}

func TestMapOrderStatuses_CombinaLosTresMapeos(t *testing.T) {
	repo := &mockRepository{
		GetOrderStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(1), nil
		},
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(2), nil
		},
		GetFulfillmentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return uintPtr(3), nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	got := uc.mapOrderStatuses(context.Background(), &dtos.ProbabilityOrderDTO{
		IntegrationType:    "shopify",
		Status:             "confirmed",
		Payments:           []dtos.ProbabilityPaymentDTO{{Status: "completed"}},
		FulfillmentDetails: []byte(`{"fulfillment_status":"shipped"}`),
	})

	require.NotNil(t, got.OrderStatusID)
	require.NotNil(t, got.PaymentStatusID)
	require.NotNil(t, got.FulfillmentStatusID)
	assert.Equal(t, uint(1), *got.OrderStatusID)
	assert.Equal(t, uint(2), *got.PaymentStatusID)
	assert.Equal(t, uint(3), *got.FulfillmentStatusID)
}
