package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// ClearRechargeHistory elimina el historial de recargas de un negocio (admin)
func (uc *walletUseCase) ClearRechargeHistory(ctx context.Context, businessID uint) error {
	wallet, err := uc.GetWallet(ctx, businessID)
	if err != nil {
		return err
	}

	return uc.repo.DeleteTransactionsByWalletIDAndType(ctx, wallet.ID, entities.WalletTxTypeRecharge)
}
