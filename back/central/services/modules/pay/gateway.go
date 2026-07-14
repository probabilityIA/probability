package pay

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
)

func (b *Bundle) GetBalance(ctx context.Context, businessID uint) (float64, error) {
	wallet, err := b.WalletUseCase.GetWallet(ctx, businessID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (b *Bundle) Debit(ctx context.Context, businessID uint, amount float64, reference, concept string, userID uint) error {
	return b.WalletUseCase.ManualDebit(ctx, &dtos.ManualDebitDTO{
		BusinessID: businessID,
		Amount:     amount,
		Reference:  reference,
		Concept:    concept,
		UserID:     &userID,
	})
}
