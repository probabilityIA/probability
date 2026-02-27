package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// GetAllWallets obtiene todas las billeteras (admin)
func (uc *walletUseCase) GetAllWallets(ctx context.Context) ([]*entities.Wallet, error) {
	return uc.repo.GetAllWallets(ctx)
}

// GetPendingTransactions obtiene transacciones de recarga pendientes (admin)
func (uc *walletUseCase) GetPendingTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	return uc.repo.GetPendingRechargeTransactions(ctx)
}

// GetProcessedTransactions obtiene transacciones procesadas (admin)
func (uc *walletUseCase) GetProcessedTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	return uc.repo.GetProcessedTransactions(ctx)
}

// GetTransactionsByBusinessID obtiene transacciones de un negocio
func (uc *walletUseCase) GetTransactionsByBusinessID(ctx context.Context, businessID uint) ([]*entities.WalletTransaction, error) {
	wallet, err := uc.GetWallet(ctx, businessID)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetTransactionsByWalletID(ctx, wallet.ID)
}
