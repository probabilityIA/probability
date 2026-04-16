package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (uc *useCase) GetProductInventory(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
	// Verificar que el producto existe
	_, _, _, err := uc.repo.GetProductByID(ctx, params.ProductID, params.BusinessID)
	if err != nil {
		return nil, domainerrors.ErrProductNotFound
	}

	return uc.repo.GetProductInventory(ctx, params)
}
