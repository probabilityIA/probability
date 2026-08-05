package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRechargeWallet_MontoInvalido_RetornaError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	for _, amount := range []float64{0, -1, -99999} {
		tx, err := uc.RechargeWallet(context.Background(), &dtos.RechargeWalletDTO{
			BusinessID: 1,
			Amount:     amount,
		})

		require.Error(t, err)
		assert.Nil(t, tx)
		assert.Contains(t, err.Error(), "amount must be greater than 0")
	}

	repo.AssertNotCalled(t, "CreateWalletTransaction", mock.Anything, mock.Anything)
}

func TestRechargeWallet_Exitoso_CreaPendienteSinAcreditarSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(15, 5000)

	repo.On("GetWalletByBusinessID", ctx, uint(15)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	tx, err := uc.RechargeWallet(ctx, &dtos.RechargeWalletDTO{
		BusinessID: 15,
		Amount:     300000,
		Reference:  "REF-BOLD-123",
		Concept:    entities.WalletTxConceptRecharge,
	})

	require.NoError(t, err)
	require.NotNil(t, tx)
	assert.Equal(t, entities.WalletTxTypeRecharge, tx.Type)
	assert.Equal(t, entities.WalletTxStatusPending, tx.Status)
	assert.Equal(t, "REF-BOLD-123", tx.Reference)
	assert.Equal(t, wallet.ID, tx.WalletID)
	assert.Equal(t, uint(15), tx.BusinessID)
	assert.InDelta(t, 300000.0, tx.Amount, 0.001)

	assert.InDelta(t, 5000.0, wallet.Balance, 0.001)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestRechargeWallet_SinReferencia_GeneraReferenciaManual(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(15, 0)

	repo.On("GetWalletByBusinessID", ctx, uint(15)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	tx, err := uc.RechargeWallet(ctx, &dtos.RechargeWalletDTO{BusinessID: 15, Amount: 1000})

	require.NoError(t, err)
	require.NotNil(t, tx)
	assert.True(t, strings.HasPrefix(tx.Reference, "MANUAL_"))

	_, parseErr := uuid.Parse(strings.TrimPrefix(tx.Reference, "MANUAL_"))
	assert.NoError(t, parseErr)
}

func TestRechargeWallet_ConceptoInvalido_UsaRecharge(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(15, 0)

	repo.On("GetWalletByBusinessID", ctx, uint(15)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	tx, err := uc.RechargeWallet(ctx, &dtos.RechargeWalletDTO{
		BusinessID: 15,
		Amount:     1000,
		Concept:    "INVENTADO",
	})

	require.NoError(t, err)
	assert.Equal(t, entities.WalletTxConceptRecharge, tx.Concept)
}

func TestRechargeWallet_FallaCrearTransaccion_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(15, 0)
	dbErr := errors.New("insert failed")

	repo.On("GetWalletByBusinessID", ctx, uint(15)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.Anything).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	tx, err := uc.RechargeWallet(ctx, &dtos.RechargeWalletDTO{BusinessID: 15, Amount: 1000})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, tx)
}

func TestRechargeWallet_FallaObtenerWallet_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dbErr := errors.New("db down")

	repo.On("GetWalletByBusinessID", ctx, uint(15)).Return(nil, dbErr)

	uc := newWalletUseCaseForTest(repo)

	tx, err := uc.RechargeWallet(ctx, &dtos.RechargeWalletDTO{BusinessID: 15, Amount: 1000})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, tx)
}

func TestRechargeYAprobacion_FlujoCompleto_AcreditaUnaSolaVez(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(15, 1000)

	var created *entities.WalletTransaction
	repo.On("GetWalletByBusinessID", ctx, uint(15)).Return(wallet, nil)
	repo.On("CreateWalletTransaction", ctx, mock.MatchedBy(func(tx *entities.WalletTransaction) bool {
		created = tx
		return true
	})).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	tx, err := uc.RechargeWallet(ctx, &dtos.RechargeWalletDTO{BusinessID: 15, Amount: 50000})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.InDelta(t, 1000.0, wallet.Balance, 0.001)

	tx.ID = uuid.New()
	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)
	repo.On("UpdateWalletTransaction", ctx, tx).Return(nil)
	repo.On("GetWalletByID", ctx, wallet.ID).Return(wallet, nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	require.NoError(t, uc.ApproveTransaction(ctx, tx.ID.String()))
	assert.InDelta(t, 51000.0, wallet.Balance, 0.001)

	err = uc.ApproveTransaction(ctx, tx.ID.String())
	require.Error(t, err)
	assert.InDelta(t, 51000.0, wallet.Balance, 0.001)
}

func TestGetTransactionsByBusinessID_ConsultaPorWalletDelNegocio(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(21, 0)
	expected := []*entities.WalletTransaction{pendingRecharge(wallet.ID, 100)}

	repo.On("GetWalletByBusinessID", ctx, uint(21)).Return(wallet, nil)
	repo.On("GetTransactionsByWalletID", ctx, wallet.ID).Return(expected, nil)

	uc := newWalletUseCaseForTest(repo)

	got, err := uc.GetTransactionsByBusinessID(ctx, 21)

	require.NoError(t, err)
	assert.Len(t, got, 1)
	repo.AssertExpectations(t)
}

func TestGetTransactionsByBusinessID_FallaWallet_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dbErr := errors.New("db down")

	repo.On("GetWalletByBusinessID", ctx, uint(21)).Return(nil, dbErr)

	uc := newWalletUseCaseForTest(repo)

	got, err := uc.GetTransactionsByBusinessID(ctx, 21)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	repo.AssertNotCalled(t, "GetTransactionsByWalletID", mock.Anything, mock.Anything)
}

func TestGetAllWallets_DelegaEnRepositorio(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	expected := []*entities.Wallet{existingWallet(1, 10), existingWallet(2, 20)}

	repo.On("GetAllWallets", ctx).Return(expected, nil)

	uc := newWalletUseCaseForTest(repo)

	got, err := uc.GetAllWallets(ctx)

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestGetPendingTransactions_DelegaEnRepositorio(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	expected := []*entities.WalletTransaction{pendingRecharge(uuid.New(), 100)}

	repo.On("GetPendingRechargeTransactions", ctx).Return(expected, nil)

	uc := newWalletUseCaseForTest(repo)

	got, err := uc.GetPendingTransactions(ctx)

	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestGetProcessedTransactions_DelegaEnRepositorio(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)

	repo.On("GetProcessedTransactions", ctx).Return([]*entities.WalletTransaction{}, nil)

	uc := newWalletUseCaseForTest(repo)

	got, err := uc.GetProcessedTransactions(ctx)

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestClearRechargeHistory_BorraTodasLasTransaccionesDelWallet(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(30, 5000)

	repo.On("GetWalletByBusinessID", ctx, uint(30)).Return(wallet, nil)
	repo.On("GetTransactionsByWalletID", ctx, wallet.ID).Return([]*entities.WalletTransaction{}, nil)
	repo.On("DeleteAllTransactionsByWalletID", ctx, wallet.ID).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ClearRechargeHistory(ctx, 30)

	require.NoError(t, err)
	repo.AssertNumberOfCalls(t, "DeleteAllTransactionsByWalletID", 1)
	assert.InDelta(t, 5000.0, wallet.Balance, 0.001)
}

func TestClearRechargeHistory_QuedanTransacciones_NoFallaPeroLasReporta(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(31, 5000)
	sobrante := []*entities.WalletTransaction{pendingRecharge(wallet.ID, 100)}

	repo.On("GetWalletByBusinessID", ctx, uint(31)).Return(wallet, nil)
	repo.On("GetTransactionsByWalletID", ctx, wallet.ID).Return(sobrante, nil)
	repo.On("DeleteAllTransactionsByWalletID", ctx, wallet.ID).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ClearRechargeHistory(ctx, 31)

	require.NoError(t, err)
	repo.AssertNumberOfCalls(t, "GetTransactionsByWalletID", 2)
}

func TestClearRechargeHistory_FallaBorrado_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(30, 5000)
	dbErr := errors.New("delete failed")

	repo.On("GetWalletByBusinessID", ctx, uint(30)).Return(wallet, nil)
	repo.On("GetTransactionsByWalletID", ctx, wallet.ID).Return([]*entities.WalletTransaction{}, nil)
	repo.On("DeleteAllTransactionsByWalletID", ctx, wallet.ID).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ClearRechargeHistory(ctx, 30)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateWalletKPISelection_GuardaSeleccionYRetornaResponse(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)

	var saved *entities.WalletKPISelection
	repo.On("UpdateWalletKPISelection", ctx, mock.MatchedBy(func(s *entities.WalletKPISelection) bool {
		saved = s
		return true
	})).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.UpdateWalletKPISelection(ctx, &dtos.UpdateWalletKPISelectionRequest{
		SelectedBusinessIDs: []uint{1, 2, 3},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, saved)
	assert.Equal(t, []uint{1, 2, 3}, resp.SelectedBusinessIDs)
	assert.Equal(t, uint(1), saved.ID)
}

func TestUpdateWalletKPISelection_FallaGuardar_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dbErr := errors.New("update failed")

	repo.On("UpdateWalletKPISelection", ctx, mock.Anything).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.UpdateWalletKPISelection(ctx, &dtos.UpdateWalletKPISelectionRequest{
		SelectedBusinessIDs: []uint{1},
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}
