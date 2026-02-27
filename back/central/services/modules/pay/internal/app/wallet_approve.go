package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

// ApproveTransaction aprueba una transacción pendiente y acredita el saldo
func (uc *walletUseCase) ApproveTransaction(ctx context.Context, transactionID string) error {
	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return payerrs.ErrTransactionNotFound
	}

	tx, err := uc.repo.GetWalletTransactionByID(ctx, txUUID)
	if err != nil {
		return err
	}

	if tx.Status != entities.WalletTxStatusPending {
		return payerrs.ErrTransactionNotPending
	}

	tx.Status = entities.WalletTxStatusCompleted
	if err := uc.repo.UpdateWalletTransaction(ctx, tx); err != nil {
		return err
	}

	wallet, err := uc.repo.GetWalletByID(ctx, tx.WalletID)
	if err != nil {
		return err
	}

	wallet.Balance += tx.Amount
	if err := uc.repo.UpdateWallet(ctx, wallet); err != nil {
		return err
	}

	uc.log.Info(ctx).
		Str("tx_id", transactionID).
		Float64("amount", tx.Amount).
		Float64("new_balance", wallet.Balance).
		Msg("Wallet transaction approved")

	return nil
}

// RejectTransaction rechaza una transacción pendiente
func (uc *walletUseCase) RejectTransaction(ctx context.Context, transactionID string) error {
	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return payerrs.ErrTransactionNotFound
	}

	tx, err := uc.repo.GetWalletTransactionByID(ctx, txUUID)
	if err != nil {
		return err
	}

	if tx.Status != entities.WalletTxStatusPending {
		return payerrs.ErrTransactionNotPending
	}

	tx.Status = entities.WalletTxStatusFailed
	if err := uc.repo.UpdateWalletTransaction(ctx, tx); err != nil {
		return err
	}

	uc.log.Info(ctx).
		Str("tx_id", transactionID).
		Msg("Wallet transaction rejected")

	return nil
}
