package app

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func pendingRecharge(walletID uuid.UUID, amount float64) *entities.WalletTransaction {
	return &entities.WalletTransaction{
		ID:       uuid.New(),
		WalletID: walletID,
		Amount:   amount,
		Type:     entities.WalletTxTypeRecharge,
		Status:   entities.WalletTxStatusPending,
		Concept:  entities.WalletTxConceptRecharge,
	}
}

func TestApproveTransaction_IDNoEsUUID_RetornaNotFound(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	err := uc.ApproveTransaction(context.Background(), "no-es-un-uuid")

	assert.ErrorIs(t, err, payerrs.ErrTransactionNotFound)
	repo.AssertNotCalled(t, "GetWalletTransactionByID", mock.Anything, mock.Anything)
}

func TestApproveTransaction_Exitoso_AcreditaSaldoYCompleta(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(20, 10000)
	tx := pendingRecharge(wallet.ID, 250000)

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)
	repo.On("UpdateWalletTransaction", ctx, tx).Return(nil)
	repo.On("GetWalletByID", ctx, wallet.ID).Return(wallet, nil)
	repo.On("UpdateWallet", ctx, wallet).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ApproveTransaction(ctx, tx.ID.String())

	require.NoError(t, err)
	assert.Equal(t, entities.WalletTxStatusCompleted, tx.Status)
	assert.InDelta(t, 260000.0, wallet.Balance, 0.001)
	repo.AssertExpectations(t)
}

func TestApproveTransaction_YaCompletada_NoAcreditaDosVeces(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(20, 10000)
	tx := pendingRecharge(wallet.ID, 250000)
	tx.Status = entities.WalletTxStatusCompleted

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ApproveTransaction(ctx, tx.ID.String())

	assert.ErrorIs(t, err, payerrs.ErrTransactionNotPending)
	assert.InDelta(t, 10000.0, wallet.Balance, 0.001)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "UpdateWalletTransaction", mock.Anything, mock.Anything)
}

func TestApproveTransaction_Rechazada_NoSePuedeAprobar(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(20, 10000)
	tx := pendingRecharge(wallet.ID, 250000)
	tx.Status = entities.WalletTxStatusFailed

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ApproveTransaction(ctx, tx.ID.String())

	assert.ErrorIs(t, err, payerrs.ErrTransactionNotPending)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestApproveTransaction_FallaBuscarTransaccion_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	id := uuid.New()
	dbErr := errors.New("wallet transaction not found")

	repo.On("GetWalletTransactionByID", ctx, id).Return(nil, dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ApproveTransaction(ctx, id.String())

	assert.ErrorIs(t, err, dbErr)
}

func TestApproveTransaction_FallaActualizarTransaccion_NoAcreditaSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(20, 10000)
	tx := pendingRecharge(wallet.ID, 250000)
	dbErr := errors.New("update failed")

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)
	repo.On("UpdateWalletTransaction", ctx, tx).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ApproveTransaction(ctx, tx.ID.String())

	assert.ErrorIs(t, err, dbErr)
	assert.InDelta(t, 10000.0, wallet.Balance, 0.001)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestApproveTransaction_FallaObtenerWallet_DejaTransaccionCompletadaSinAcreditar(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	walletID := uuid.New()
	tx := pendingRecharge(walletID, 250000)
	dbErr := errors.New("wallet not found")

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)
	repo.On("UpdateWalletTransaction", ctx, tx).Return(nil)
	repo.On("GetWalletByID", ctx, walletID).Return(nil, dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.ApproveTransaction(ctx, tx.ID.String())

	assert.ErrorIs(t, err, dbErr)
	assert.Equal(t, entities.WalletTxStatusCompleted, tx.Status)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestRejectTransaction_IDNoEsUUID_RetornaNotFound(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	err := uc.RejectTransaction(context.Background(), "abc")

	assert.ErrorIs(t, err, payerrs.ErrTransactionNotFound)
}

func TestRejectTransaction_Exitoso_MarcaFailedSinTocarSaldo(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	wallet := existingWallet(20, 10000)
	tx := pendingRecharge(wallet.ID, 250000)

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)
	repo.On("UpdateWalletTransaction", ctx, tx).Return(nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.RejectTransaction(ctx, tx.ID.String())

	require.NoError(t, err)
	assert.Equal(t, entities.WalletTxStatusFailed, tx.Status)
	assert.InDelta(t, 10000.0, wallet.Balance, 0.001)
	repo.AssertNotCalled(t, "UpdateWallet", mock.Anything, mock.Anything)
}

func TestRejectTransaction_YaCompletada_NoSePuedeRechazar(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	tx := pendingRecharge(uuid.New(), 250000)
	tx.Status = entities.WalletTxStatusCompleted

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)

	uc := newWalletUseCaseForTest(repo)

	err := uc.RejectTransaction(ctx, tx.ID.String())

	assert.ErrorIs(t, err, payerrs.ErrTransactionNotPending)
	assert.Equal(t, entities.WalletTxStatusCompleted, tx.Status)
	repo.AssertNotCalled(t, "UpdateWalletTransaction", mock.Anything, mock.Anything)
}

func TestRejectTransaction_FallaActualizar_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	tx := pendingRecharge(uuid.New(), 250000)
	dbErr := errors.New("update failed")

	repo.On("GetWalletTransactionByID", ctx, tx.ID).Return(tx, nil)
	repo.On("UpdateWalletTransaction", ctx, tx).Return(dbErr)

	uc := newWalletUseCaseForTest(repo)

	err := uc.RejectTransaction(ctx, tx.ID.String())

	assert.ErrorIs(t, err, dbErr)
}
