package app

import (
	"context"
	"fmt"
)

// ClearRechargeHistory elimina TODOS el historial de transacciones de un negocio (admin)
// Borra TODAS las transacciones sin considerar el tipo
func (uc *walletUseCase) ClearRechargeHistory(ctx context.Context, businessID uint) error {
	wallet, err := uc.GetWallet(ctx, businessID)
	if err != nil {
		return fmt.Errorf("error getting wallet for business %d: %w", businessID, err)
	}

	if wallet == nil {
		return fmt.Errorf("wallet not found for business %d", businessID)
	}

	// Obtener transacciones antes de borrar para logging
	allTxs, _ := uc.repo.GetTransactionsByWalletID(ctx, wallet.ID)
	fmt.Printf("DEBUG: Business %d, Wallet %s has %d total transactions before delete\n", businessID, wallet.ID, len(allTxs))
	for _, tx := range allTxs {
		fmt.Printf("  - TX %s: Type=%s, Amount=%.2f, Status=%s\n", tx.ID, tx.Type, tx.Amount, tx.Status)
	}

	// Borrar TODAS las transacciones de esta billetera
	fmt.Printf("DEBUG: Deleting ALL transactions for wallet %s\n", wallet.ID)
	if err := uc.repo.DeleteAllTransactionsByWalletID(ctx, wallet.ID); err != nil {
		return fmt.Errorf("error deleting all transactions: %w", err)
	}

	// Verificar después de borrar
	remainingTxs, _ := uc.repo.GetTransactionsByWalletID(ctx, wallet.ID)
	fmt.Printf("DEBUG: After delete, wallet %s has %d remaining transactions\n", wallet.ID, len(remainingTxs))
	if len(remainingTxs) > 0 {
		for _, tx := range remainingTxs {
			fmt.Printf("  - TX %s: Type=%s, Amount=%.2f, Status=%s\n", tx.ID, tx.Type, tx.Amount, tx.Status)
		}
	}

	return nil
}
