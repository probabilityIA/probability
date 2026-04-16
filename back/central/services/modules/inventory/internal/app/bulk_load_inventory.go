package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) BulkLoadInventory(ctx context.Context, dto dtos.BulkLoadDTO) (*dtos.BulkLoadResult, error) {
	if len(dto.Items) == 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}

	// Validar warehouse
	exists, err := uc.repo.WarehouseExists(ctx, dto.WarehouseID, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWarehouseNotFound
	}

	// Resolver movement type "inbound"
	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "inbound")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve inbound movement type: %w", err)
	}

	result := &dtos.BulkLoadResult{
		TotalItems: len(dto.Items),
		Items:      make([]dtos.BulkLoadItemResult, 0, len(dto.Items)),
	}

	for _, item := range dto.Items {
		itemResult := uc.processBulkLoadItem(ctx, dto, item, movTypeID)
		if itemResult.Success {
			result.SuccessCount++
		} else {
			result.FailureCount++
		}
		result.Items = append(result.Items, itemResult)
	}

	return result, nil
}

func (uc *useCase) processBulkLoadItem(ctx context.Context, dto dtos.BulkLoadDTO, item dtos.BulkLoadItem, movTypeID uint) dtos.BulkLoadItemResult {
	itemResult := dtos.BulkLoadItemResult{
		SKU: item.SKU,
	}

	if item.Quantity <= 0 {
		itemResult.Error = "quantity must be positive"
		return itemResult
	}

	// Resolver producto por SKU
	productID, _, trackInventory, err := uc.repo.GetProductBySKU(ctx, item.SKU, dto.BusinessID)
	if err != nil {
		itemResult.Error = fmt.Sprintf("product not found for SKU %s", item.SKU)
		return itemResult
	}
	itemResult.ProductID = productID

	// Auto-habilitar track_inventory si no está activo
	if !trackInventory {
		if err := uc.repo.EnableProductTrackInventory(ctx, productID); err != nil {
			itemResult.Error = fmt.Sprintf("failed to enable inventory tracking: %v", err)
			return itemResult
		}
	}

	// Ejecutar ajuste usando AdjustStockTx existente
	reason := dto.Reason
	if reason == "" {
		reason = "bulk_load"
	}

	txResult, err := uc.repo.AdjustStockTx(ctx, dtos.AdjustStockTxParams{
		ProductID:      productID,
		WarehouseID:    dto.WarehouseID,
		LocationID:     nil,
		BusinessID:     dto.BusinessID,
		Quantity:       item.Quantity,
		MovementTypeID: movTypeID,
		Reason:         reason,
		Notes:          fmt.Sprintf("Bulk load - SKU: %s", item.SKU),
		ReferenceType:  "bulk_load",
		CreatedByID:    dto.CreatedByID,
	})
	if err != nil {
		itemResult.Error = fmt.Sprintf("failed to adjust stock: %v", err)
		return itemResult
	}

	itemResult.Success = true
	if txResult.Movement != nil {
		itemResult.PreviousQty = txResult.Movement.PreviousQty
		itemResult.NewQty = txResult.Movement.NewQty
	}

	// Actualizar min/max/reorder en el inventory level si se proporcionaron
	if txResult.Level != nil && (item.MinStock != nil || item.MaxStock != nil || item.ReorderPoint != nil) {
		level := txResult.Level
		if item.MinStock != nil {
			level.MinStock = item.MinStock
		}
		if item.MaxStock != nil {
			level.MaxStock = item.MaxStock
		}
		if item.ReorderPoint != nil {
			level.ReorderPoint = item.ReorderPoint
		}
		_ = uc.repo.UpdateLevel(ctx, level)
	}

	// Actualizar stock total del producto (best-effort)
	uc.updateProductTotalStock(ctx, productID, dto.BusinessID)

	// Publicar sync a canales de venta (fire-and-forget)
	uc.publishSync(ctx, productID, dto.BusinessID, txResult.NewQuantity, dto.WarehouseID, "bulk_load")

	return itemResult
}
