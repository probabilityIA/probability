package usecasecreateorder

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncIsPaidFromPaymentStatus_EstadoPagado_MarcaPagadaYFecha(t *testing.T) {
	paidID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			assert.Equal(t, "paid", code)
			return &paidID, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	order := &entities.ProbabilityOrder{}

	antes := time.Now()
	uc.syncIsPaidFromPaymentStatus(context.Background(), order, &paidID)

	assert.True(t, order.IsPaid)
	require.NotNil(t, order.PaidAt)
	assert.WithinDuration(t, antes, *order.PaidAt, 5*time.Second)
}

func TestSyncIsPaidFromPaymentStatus_YaTeniaFecha_NoLaSobreescribe(t *testing.T) {
	paidID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return &paidID, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)

	original := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	order := &entities.ProbabilityOrder{PaidAt: &original}

	uc.syncIsPaidFromPaymentStatus(context.Background(), order, &paidID)

	assert.True(t, order.IsPaid)
	require.NotNil(t, order.PaidAt)
	assert.True(t, order.PaidAt.Equal(original), "la fecha de pago original debe conservarse")
}

func TestSyncIsPaidFromPaymentStatus_SinEstado_NoHaceNada(t *testing.T) {
	consultado := false
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			consultado = true
			return nil, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	order := &entities.ProbabilityOrder{}

	uc.syncIsPaidFromPaymentStatus(context.Background(), order, nil)

	assert.False(t, order.IsPaid)
	assert.Nil(t, order.PaidAt)
	assert.False(t, consultado)
}

func TestSyncIsPaidFromPaymentStatus_EstadoDistintoDePagado_NoMarca(t *testing.T) {
	paidID := uint(3)
	pendingID := uint(1)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return &paidID, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	order := &entities.ProbabilityOrder{}

	uc.syncIsPaidFromPaymentStatus(context.Background(), order, &pendingID)

	assert.False(t, order.IsPaid)
	assert.Nil(t, order.PaidAt)
}

func TestSyncIsPaidFromPaymentStatus_FallaConsulta_NoMarca(t *testing.T) {
	algunID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return nil, errors.New("db down")
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	order := &entities.ProbabilityOrder{}

	uc.syncIsPaidFromPaymentStatus(context.Background(), order, &algunID)

	assert.False(t, order.IsPaid, "un fallo de base de datos no debe marcar una orden como pagada")
	assert.Nil(t, order.PaidAt)
}

func TestSyncIsPaidFromPaymentStatus_CatalogoSinPaid_NoMarca(t *testing.T) {
	algunID := uint(3)
	repo := &mockRepository{
		GetPaymentStatusIDByCodeFn: func(ctx context.Context, code string) (*uint, error) {
			return nil, nil
		},
	}
	uc := newTestCreateUseCase(repo, nil, nil, nil)
	order := &entities.ProbabilityOrder{}

	uc.syncIsPaidFromPaymentStatus(context.Background(), order, &algunID)

	assert.False(t, order.IsPaid)
}
