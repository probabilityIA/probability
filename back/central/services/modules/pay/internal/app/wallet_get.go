package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

// GetWallet obtiene la billetera de un negocio, creándola si no existe
func (uc *walletUseCase) GetWallet(ctx context.Context, businessID uint) (*entities.Wallet, error) {
	wallet, err := uc.repo.GetWalletByBusinessID(ctx, businessID)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		wallet = &entities.Wallet{
			BusinessID: businessID,
			Balance:    0,
		}
		if err = uc.repo.CreateWallet(ctx, wallet); err != nil {
			return nil, err
		}
	}

	return wallet, nil
}
