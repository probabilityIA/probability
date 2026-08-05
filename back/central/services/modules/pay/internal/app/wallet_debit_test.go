package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newWalletUseCaseForTest(repo ports.IRepository) ports.IWalletUseCase {
	return NewWalletUseCase(repo, nil, nil, nil, mocks.NewSilentLogger())
}

func existingWallet(businessID uint, balance float64) *entities.Wallet {
	return &entities.Wallet{
		ID:         uuid.New(),
		BusinessID: businessID,
		Balance:    balance,
	}
}

func TestManualDebit_MontoCero_RetornaErrorSinTocarRepo(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	err := uc.ManualDebit(context.Background(), &dtos.ManualDebitDTO{
		BusinessID: 1,
		Amount:     0,
		Reference:  "test",
	})

	assert.ErrorIs(t, err, payerrs.ErrInvalidAmount)
	repo.AssertNotCalled(t, "GetWalletByBusinessID", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "CreateWalletTransaction", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestManualDebit_MontoNegativo_RetornaErrorSinTocarRepo(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	err := uc.ManualDebit(context.Background(), &dtos.ManualDebitDTO{
		BusinessID: 1,
		Amount:     -500,
		Reference:  "test",
	})

	assert.ErrorIs(t, err, payerrs.ErrInvalidAmount)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestManualDebit_Exitoso_CreaTransaccionUsageYRestaSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(7, 50000)

	var createdTx *entities.WalletTransaction
	repo.On("GetWalletByBusinessID", ctx, uint(7)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.MatchedBy(func(tx *entities.WalletTransaction) bool {
		createdTx = tx
		return true
	})).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)
	userID := uint(42)

	err := uc.ManualDebit(ctx, &dtos.ManualDebitDTO{
		BusinessID: 7,
		Amount:     16353,
		Reference:  "ajuste manual",
		Concept:    entities.WalletTxConceptGuide,
		UserID:     &userID,
	})

	require.NoError(t, err)
	require.NotNil(t, createdTx)

	assert.Equal(t, wallet.ID, createdTx.WalletID)
	assert.InDelta(t, 16353.0, createdTx.Amount, 0.001)
	assert.Equal(t, entities.WalletTxTypeUsage, createdTx.Type)
	assert.Equal(t, entities.WalletTxStatusCompleted, createdTx.Status)
	assert.Equal(t, entities.WalletTxConceptGuide, createdTx.Concept)
	assert.Equal(t, uint(7), createdTx.BusinessID)
	assert.Equal(t, &userID, createdTx.UserID)
	assert.True(t, strings.HasPrefix(createdTx.Reference, "MAN_DEB_"))
	assert.True(t, strings.HasSuffix(createdTx.Reference, ": ajuste manual"))

	assert.InDelta(t, 33647.0, wallet.Balance, 0.001)
	repo.AssertExpectations(t)
}

func TestManualDebit_ConceptoInvalido_SeNormalizaAOther(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(3, 1000)

	var createdTx *entities.WalletTransaction
	repo.On("GetWalletByBusinessID", ctx, uint(3)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.MatchedBy(func(tx *entities.WalletTransaction) bool {
		createdTx = tx
		return true
	})).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ManualDebit(ctx, &dtos.ManualDebitDTO{
		BusinessID: 3,
		Amount:     100,
		Concept:    "CONCEPTO_QUE_NO_EXISTE",
	})

	require.NoError(t, err)
	require.NotNil(t, createdTx)
	assert.Equal(t, entities.WalletTxConceptOther, createdTx.Concept)
}

func TestManualDebit_FallaCrearTransaccion_NoModificaSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(9, 20000)
	dbErr := errors.New("db down")

	repo.On("GetWalletByBusinessID", ctx, uint(9)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ManualDebit(ctx, &dtos.ManualDebitDTO{BusinessID: 9, Amount: 5000})

	assert.ErrorIs(t, err, dbErr)
	assert.InDelta(t, 20000.0, wallet.Balance, 0.001)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestManualDebit_FallaObtenerWallet_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dbErr := errors.New("connection refused")

	repo.On("GetWalletByBusinessID", ctx, uint(11)).Return(nil, dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ManualDebit(ctx, &dtos.ManualDebitDTO{BusinessID: 11, Amount: 100})

	assert.ErrorIs(t, err, dbErr)
	repo.AssertNotCalled(t, "CreateWalletTransaction", mock.Anything, mock.Anything)
}

func TestManualDebit_SaldoInsuficiente_HoyPermiteBalanceNegativo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(5, 1000)

	repo.On("GetWalletByBusinessID", ctx, uint(5)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ManualDebit(ctx, &dtos.ManualDebitDTO{BusinessID: 5, Amount: 50000})

	require.NoError(t, err)
	assert.InDelta(t, -49000.0, wallet.Balance, 0.001)
	assert.Less(t, wallet.Balance, 0.0)
}

func TestDebitForGuide_Exitoso_UsaConceptoGuideYReferenciaConTracking(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(36, 100000)

	var createdTx *entities.WalletTransaction
	repo.On("GetWalletByBusinessID", ctx, uint(36)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.MatchedBy(func(tx *entities.WalletTransaction) bool {
		createdTx = tx
		return true
	})).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.DebitForGuide(ctx, &dtos.DebitForGuideDTO{
		BusinessID:     36,
		Amount:         16353,
		TrackingNumber: "034057376067",
	})

	require.NoError(t, err)
	require.NotNil(t, createdTx)
	assert.Equal(t, entities.WalletTxConceptGuide, createdTx.Concept)
	assert.Equal(t, entities.WalletTxTypeUsage, createdTx.Type)
	assert.Contains(t, createdTx.Reference, "Guide generation: 034057376067")
	assert.InDelta(t, 83647.0, wallet.Balance, 0.001)
}

func TestDebitForGuide_MontoInvalido_RetornaErrorSinTocarRepo(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	err := uc.DebitForGuide(context.Background(), &dtos.DebitForGuideDTO{
		BusinessID:     36,
		Amount:         0,
		TrackingNumber: "034057376067",
	})

	assert.ErrorIs(t, err, payerrs.ErrInvalidAmount)
	repo.AssertNotCalled(t, "GetWalletByBusinessID", mock.Anything, mock.Anything)
}

func TestDebitForGuide_DosLlamadasMismoTracking_DebitanDosVeces(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(36, 100000)

	repo.On("GetWalletByBusinessID", ctx, uint(36)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)
	dto := &dtos.DebitForGuideDTO{BusinessID: 36, Amount: 16353, TrackingNumber: "034057376067"}

	require.NoError(t, uc.DebitForGuide(ctx, dto))
	require.NoError(t, uc.DebitForGuide(ctx, dto))

	assert.InDelta(t, 67294.0, wallet.Balance, 0.001)
	repo.AssertNumberOfCalls(t, "CreateWalletTransaction", 2)
}
