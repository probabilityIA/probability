package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// ClearRechargeHistory elimina TODOS el historial de transacciones de un negocio (admin)
// Borra tanto recharges como debits/usage y todas las transacciones asociadas
func (uc *walletUseCase) ClearRechargeHistory(ctx context.Context, businessID uint) error {
	wallet, err := uc.GetWallet(ctx, businessID)
	if err != nil {
		return err
	}

	// Borrar transacciones RECHARGE
	if err := uc.repo.DeleteTransactionsByWalletIDAndType(ctx, wallet.ID, entities.WalletTxTypeRecharge); err != nil {
		return err
	}

	// Borrar transacciones USAGE (debits, consumos, etc.)
	if err := uc.repo.DeleteTransactionsByWalletIDAndType(ctx, wallet.ID, entities.WalletTxTypeUsage); err != nil {
		return err
	}

	return nil
}
