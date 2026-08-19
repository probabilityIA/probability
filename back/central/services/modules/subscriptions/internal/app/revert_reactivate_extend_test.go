package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevertPayment_Exitoso_RegistraAuditoria(t *testing.T) {
	var vistoBusinessID, vistoActor uint
	var vistaAccion string
	repo := &mocks.RepositoryMock{
		RevertSubscriptionAndRecalculateFn: func(ctx context.Context, subscriptionID uint) (*entities.BusinessSubscription, error) {
			return &entities.BusinessSubscription{ID: subscriptionID, BusinessID: 26, Amount: 999, SubscriptionTypeName: "Basico"}, nil
		},
		CreateAuditLogFn: func(ctx context.Context, log *entities.SubscriptionAuditLog) error {
			vistoBusinessID = log.BusinessID
			vistoActor = *log.ActorUserID
			vistaAccion = log.Action
			return nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).RevertPayment(context.Background(), 55, 7)

	require.NoError(t, err)
	assert.Equal(t, uint(26), got.BusinessID)
	assert.Equal(t, uint(26), vistoBusinessID)
	assert.Equal(t, uint(7), vistoActor)
	assert.Equal(t, entities.AuditActionPaymentReverted, vistaAccion)
}

func TestRevertPayment_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("ya estaba revertido")
	repo := &mocks.RepositoryMock{
		RevertSubscriptionAndRecalculateFn: func(ctx context.Context, subscriptionID uint) (*entities.BusinessSubscription, error) {
			return nil, dbErr
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).RevertPayment(context.Background(), 55, 7)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestReactivateSubscription_Exitoso_RestauraElEstadoActivo(t *testing.T) {
	fin := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	var vistoStatus string
	var vistoFin *time.Time
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return &entities.BusinessSubscription{ID: 9, BusinessID: businessID, EndDate: fin}, nil
		},
		UpdateBusinessSubscriptionStatusFn: func(ctx context.Context, businessID uint, status string, endDate *time.Time) error {
			vistoStatus, vistoFin = status, endDate
			return nil
		},
	}

	err := newSubsUseCase(repo, nil, nil).ReactivateSubscription(context.Background(), 26, 1)

	require.NoError(t, err)
	assert.Equal(t, entities.BusinessStatusActive, vistoStatus)
	require.NotNil(t, vistoFin)
	assert.Equal(t, fin, *vistoFin)
}

func TestReactivateSubscription_SinSuscripcionPrevia_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
	}

	err := newSubsUseCase(repo, nil, nil).ReactivateSubscription(context.Background(), 26, 1)

	assert.ErrorIs(t, err, errs.ErrNothingToReactivate)
}

func TestExtendCourtesy_DiasInvalidos_Rechaza(t *testing.T) {
	for _, dias := range []int{0, -5} {
		_, err := newSubsUseCase(nil, nil, nil).ExtendCourtesy(context.Background(),
			dtos.ExtendCourtesyDTO{BusinessID: 26, Days: dias, Reason: "cortesia"}, 1)

		assert.ErrorIs(t, err, errs.ErrInvalidDays)
	}
}

func TestExtendCourtesy_ExtiendeDesdeElVencimientoVigente(t *testing.T) {
	finActual := time.Now().AddDate(0, 0, 10)
	inicio := time.Now().AddDate(0, -1, 0)
	var vistoFin time.Time
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return &entities.BusinessSubscription{ID: 9, StartDate: inicio, EndDate: finActual}, nil
		},
		UpdateSubscriptionDatesFn: func(ctx context.Context, id uint, s, e time.Time) error {
			vistoFin = e
			return nil
		},
	}

	got, err := newSubsUseCase(repo, nil, nil).ExtendCourtesy(context.Background(),
		dtos.ExtendCourtesyDTO{BusinessID: 26, Days: 7, Reason: "cliente pidio prorroga"}, 1)

	require.NoError(t, err)
	assert.Equal(t, finActual.AddDate(0, 0, 7).Day(), vistoFin.Day())
	assert.Equal(t, finActual.AddDate(0, 0, 7).Day(), got.EndDate.Day())
}

func TestExtendCourtesy_SinSuscripcion_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetLatestByBusinessIDFn: func(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
			return nil, nil
		},
	}

	_, err := newSubsUseCase(repo, nil, nil).ExtendCourtesy(context.Background(),
		dtos.ExtendCourtesyDTO{BusinessID: 26, Days: 7, Reason: "x"}, 1)

	assert.ErrorIs(t, err, errs.ErrSubscriptionNotFound)
}

func TestRegisterPayment_MetodoWallet_DescuentaSaldo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, Code: "pro", Price: 1000}, nil
		},
	}
	wallet := &mocks.WalletDebiterMock{}
	metodo := entities.PaymentMethodWallet

	_, err := newSubsUseCase(repo, wallet, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 2, PaymentMethod: &metodo}, 7)

	require.NoError(t, err)
	require.Len(t, wallet.DebitCalls, 1, "con metodo Wallet si debe descontar saldo real")
	assert.InDelta(t, 2000.0, wallet.DebitCalls[0].Amount, 0.001)
	assert.Equal(t, uint(26), wallet.DebitCalls[0].BusinessID)
	assert.Equal(t, uint(7), wallet.DebitCalls[0].UserID)
}

func TestRegisterPayment_ErrorAlDescontarWallet_SePropaga(t *testing.T) {
	dbErr := stderrors.New("saldo insuficiente")
	repo := &mocks.RepositoryMock{
		GetSubscriptionTypeFn: func(ctx context.Context, id uint) (*entities.SubscriptionType, error) {
			return &entities.SubscriptionType{ID: id, Code: "pro", Price: 1000}, nil
		},
	}
	wallet := &mocks.WalletDebiterMock{
		DebitFn: func(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error {
			return dbErr
		},
	}
	metodo := entities.PaymentMethodWallet

	_, err := newSubsUseCase(repo, wallet, nil).RegisterPayment(context.Background(),
		dtos.RegisterPaymentDTO{BusinessID: 26, SubscriptionTypeID: 3, Months: 1, PaymentMethod: &metodo}, 7)

	assert.ErrorIs(t, err, dbErr)
}
