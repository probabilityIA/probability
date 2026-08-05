package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAdminAdjustBalance_MontoCero_RetornaError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	err := uc.AdminAdjustBalance(context.Background(), &dtos.AdminAdjustBalanceDTO{
		BusinessID: 1,
		Amount:     0,
		Reference:  "ajuste",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "amount cannot be zero")
	repo.AssertNotCalled(t, "GetWalletByBusinessID", mock.Anything, mock.Anything)
}

func TestAdminAdjustBalance_MontoPositivo_CreaRechargeYSumaSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(12, 5000)

	var createdTx *entities.WalletTransaction
	repo.On("GetWalletByBusinessID", ctx, uint(12)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.MatchedBy(func(tx *entities.WalletTransaction) bool {
		createdTx = tx
		return true
	})).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.AdminAdjustBalance(ctx, &dtos.AdminAdjustBalanceDTO{
		BusinessID: 12,
		Amount:     20000,
		Reference:  "reembolso guia duplicada",
		Concept:    entities.WalletTxConceptRefund,
	})

	require.NoError(t, err)
	require.NotNil(t, createdTx)
	assert.Equal(t, entities.WalletTxTypeRecharge, createdTx.Type)
	assert.Equal(t, entities.WalletTxStatusCompleted, createdTx.Status)
	assert.Equal(t, entities.WalletTxConceptRefund, createdTx.Concept)
	assert.InDelta(t, 20000.0, createdTx.Amount, 0.001)
	assert.NotEqual(t, "", createdTx.ID.String())
	assert.True(t, strings.HasPrefix(createdTx.Reference, "ADMIN_ADJ_"))
	assert.Contains(t, createdTx.Reference, "reembolso guia duplicada")
	assert.InDelta(t, 25000.0, wallet.Balance, 0.001)
}

func TestAdminAdjustBalance_MontoNegativo_CreaUsageConMontoPositivoYRestaSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(12, 5000)

	var createdTx *entities.WalletTransaction
	repo.On("GetWalletByBusinessID", ctx, uint(12)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.MatchedBy(func(tx *entities.WalletTransaction) bool {
		createdTx = tx
		return true
	})).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.AdminAdjustBalance(ctx, &dtos.AdminAdjustBalanceDTO{
		BusinessID: 12,
		Amount:     -3000,
		Reference:  "correccion",
		Concept:    entities.WalletTxConceptAdjustment,
	})

	require.NoError(t, err)
	require.NotNil(t, createdTx)
	assert.Equal(t, entities.WalletTxTypeUsage, createdTx.Type)
	assert.InDelta(t, 3000.0, createdTx.Amount, 0.001)
	assert.InDelta(t, 2000.0, wallet.Balance, 0.001)
}

func TestAdminAdjustBalance_ConceptoInvalido_SeNormalizaAOther(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(12, 5000)

	var createdTx *entities.WalletTransaction
	repo.On("GetWalletByBusinessID", ctx, uint(12)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.MatchedBy(func(tx *entities.WalletTransaction) bool {
		createdTx = tx
		return true
	})).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.AdminAdjustBalance(ctx, &dtos.AdminAdjustBalanceDTO{
		BusinessID: 12,
		Amount:     100,
		Concept:    "NO_EXISTE",
	})

	require.NoError(t, err)
	require.NotNil(t, createdTx)
	assert.Equal(t, entities.WalletTxConceptOther, createdTx.Concept)
}

func TestAdminAdjustBalance_FallaCrearTransaccion_NoModificaSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(12, 5000)
	dbErr := errors.New("db down")

	repo.On("GetWalletByBusinessID", ctx, uint(12)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.AdminAdjustBalance(ctx, &dtos.AdminAdjustBalanceDTO{BusinessID: 12, Amount: 1000})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.InDelta(t, 5000.0, wallet.Balance, 0.001)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestAdminAdjustBalance_FallaActualizarWallet_DejaTransaccionHuerfana(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(12, 5000)
	dbErr := errors.New("update failed")

	repo.On("GetWalletByBusinessID", ctx, uint(12)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(nil)
	repo.On("UpdateWallet", ctx, wallet).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.AdminAdjustBalance(ctx, &dtos.AdminAdjustBalanceDTO{BusinessID: 12, Amount: -1000})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	repo.AssertNumberOfCalls(t, "CreateWalletTransaction", 1)
}

func TestGetWallet_NoExiste_CreaConSaldoCero(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)

	var created *entities.Wallet
	repo.On("GetWalletByBusinessID", ctx, uint(99)).Return(nil, nil)
	repo.On("CreateWallet", ctx, mock.MatchedBy(func(w *entities.Wallet) bool {
		created = w
		return true
	})).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	wallet, err := uc.GetWallet(ctx, 99)

	require.NoError(t, err)
	require.NotNil(t, wallet)
	require.NotNil(t, created)
	assert.Equal(t, uint(99), wallet.BusinessID)
	assert.InDelta(t, 0.0, wallet.Balance, 0.001)
}

func TestGetWallet_Existe_NoCreaOtra(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(4, 777)

	repo.On("GetWalletByBusinessID", ctx, uint(4)).Return(wallet, nil)

	uc := newWalletUseCaseForTest(repo)

	got, err := uc.GetWallet(ctx, 4)

	require.NoError(t, err)
	assert.Equal(t, wallet.ID, got.ID)
	assert.InDelta(t, 777.0, got.Balance, 0.001)
	repo.AssertNotCalled(t, "CreateWallet", mock.Anything, mock.Anything)
}

func TestGetWallet_FallaCrear_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dbErr := errors.New("insert failed")

	repo.On("GetWalletByBusinessID", ctx, uint(99)).Return(nil, nil)
	repo.On("CreateWallet", ctx, mock.Anything).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	wallet, err := uc.GetWallet(ctx, 99)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, wallet)
}
