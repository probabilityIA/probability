package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
)

func (uc *walletUseCase) UpdateWalletKPISelection(ctx context.Context, req *dtos.UpdateWalletKPISelectionRequest) (*dtos.WalletKPISelectionResponse, error) {
	selection := &entities.WalletKPISelection{
		ID:                  1,
		SelectedBusinessIDs: req.SelectedBusinessIDs,
	}

	if err := uc.repo.UpdateWalletKPISelection(ctx, selection); err != nil {
		return nil, err
	}

	return &dtos.WalletKPISelectionResponse{
		ID:                  selection.ID,
		SelectedBusinessIDs: selection.SelectedBusinessIDs,
	}, nil
}
