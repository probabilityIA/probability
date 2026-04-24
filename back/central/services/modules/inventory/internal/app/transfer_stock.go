package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) TransferStock(ctx context.Context, dto request.TransferStockDTO) error {
	if dto.Quantity <= 0 {
		return domainerrors.ErrTransferQtyNeg
	}
	if dto.FromWarehouseID == dto.ToWarehouseID {
		return domainerrors.ErrSameWarehouse
	}

	_, _, trackInventory, err := uc.repo.GetProductByID(ctx, dto.ProductID, dto.BusinessID)
	if err != nil {
		return domainerrors.ErrProductNotFound
	}

	if !trackInventory {
		return domainerrors.ErrProductNoTracking
	}

	for _, whID := range []uint{dto.FromWarehouseID, dto.ToWarehouseID} {
		exists, err := uc.repo.WarehouseExists(ctx, whID, dto.BusinessID)
		if err != nil {
			return err
		}
		if !exists {
			return domainerrors.ErrWarehouseNotFound
		}
	}

	transferTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "transfer")
	if err != nil {
		return err
	}

	txResult, err := uc.repo.TransferStockTx(ctx, mappers.TransferStockDTOToTxParams(dto, transferTypeID, "manual"))
	if err != nil {
		return err
	}

	uc.publishSync(ctx, dto.ProductID, dto.BusinessID, txResult.FromNewQty+txResult.ToNewQty, dto.FromWarehouseID, "transfer")

	return nil
}
