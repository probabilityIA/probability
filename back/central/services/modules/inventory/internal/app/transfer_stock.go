package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) TransferStock(ctx context.Context, dto dtos.TransferStockDTO) error {
	if dto.Quantity <= 0 {
		return domainerrors.ErrTransferQtyNeg
	}
	if dto.FromWarehouseID == dto.ToWarehouseID {
		return domainerrors.ErrSameWarehouse
	}

	// Verificar producto
	_, _, trackInventory, err := uc.repo.GetProductByID(ctx, dto.ProductID, dto.BusinessID)
	if err != nil {
		return domainerrors.ErrProductNotFound
	}
	if !trackInventory {
		return domainerrors.ErrProductNoTracking
	}

	// Verificar bodegas
	for _, whID := range []uint{dto.FromWarehouseID, dto.ToWarehouseID} {
		exists, err := uc.repo.WarehouseExists(ctx, whID, dto.BusinessID)
		if err != nil {
			return err
		}
		if !exists {
			return domainerrors.ErrWarehouseNotFound
		}
	}

	// Obtener tipo de movimiento "transfer"
	transferTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "transfer")
	if err != nil {
		return err
	}

	// Ejecutar transferencia dentro de transacción con SELECT FOR UPDATE
	txResult, err := uc.repo.TransferStockTx(ctx, dtos.TransferStockTxParams{
		ProductID:       dto.ProductID,
		FromWarehouseID: dto.FromWarehouseID,
		ToWarehouseID:   dto.ToWarehouseID,
		FromLocationID:  dto.FromLocationID,
		ToLocationID:    dto.ToLocationID,
		BusinessID:      dto.BusinessID,
		Quantity:        dto.Quantity,
		MovementTypeID:  transferTypeID,
		Reason:          dto.Reason,
		Notes:           dto.Notes,
		ReferenceType:   "manual",
		CreatedByID:     dto.CreatedByID,
	})
	if err != nil {
		return err
	}

	// Publicar sync (el stock total no cambió, pero las bodegas sí)
	uc.publishSync(ctx, dto.ProductID, dto.BusinessID, txResult.FromNewQty+txResult.ToNewQty, dto.FromWarehouseID, "transfer")

	return nil
}
