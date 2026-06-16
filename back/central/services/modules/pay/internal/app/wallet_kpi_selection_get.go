package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
)

func (uc *walletUseCase) GetWalletKPISelection(ctx context.Context) (*dtos.WalletKPISelectionResponse, error) {
	selection, err := uc.repo.GetWalletKPISelection(ctx)
	if err != nil {
		return nil, err
	}

	return &dtos.WalletKPISelectionResponse{
		ID:                  selection.ID,
		SelectedBusinessIDs: selection.SelectedBusinessIDs,
	}, nil
}
