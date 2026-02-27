package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
)

// ReleaseStockForOrder libera reservas cuando una orden es cancelada.
// ReservedQty -= qty, AvailableQty += qty.
func (uc *UseCase) ReleaseStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error) {
	// Resolver warehouse
	whID, err := uc.resolveWarehouse(ctx, warehouseID, businessID)
	if err != nil {
		return nil, err
	}

	// Obtener movement type
	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "release")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get release movement type")
		return nil, err
	}

	result := &dtos.OrderStockResult{
		OrderID:     orderID,
		BusinessID:  businessID,
		WarehouseID: whID,
		Success:     true,
	}

	for _, item := range items {
		itemResult := dtos.ItemStockResult{
			ProductID: item.ProductID,
			SKU:       item.SKU,
			Requested: item.Quantity,
		}

		// Verificar producto y tracking
		_, _, trackInventory, err := uc.repo.GetProductByID(ctx, item.ProductID, businessID)
		if err != nil {
			itemResult.ErrorMessage = "producto no encontrado"
			result.ItemResults = append(result.ItemResults, itemResult)
			continue
		}
		if !trackInventory {
			itemResult.Processed = item.Quantity
			itemResult.Sufficient = true
			result.ItemResults = append(result.ItemResults, itemResult)
			continue
		}

		// Ejecutar liberación transaccional
		err = uc.repo.ReleaseStockTx(ctx, dtos.ReleaseTxParams{
			ProductID:      item.ProductID,
			WarehouseID:    whID,
			BusinessID:     businessID,
			Quantity:       item.Quantity,
			MovementTypeID: movTypeID,
			OrderID:        orderID,
		})
		if err != nil {
			itemResult.ErrorMessage = err.Error()
			result.ItemResults = append(result.ItemResults, itemResult)
			result.Success = false
			continue
		}

		itemResult.Processed = item.Quantity
		itemResult.Sufficient = true
		result.ItemResults = append(result.ItemResults, itemResult)

		// Actualizar stock total del producto (best-effort)
		uc.updateProductTotalStock(ctx, item.ProductID, businessID)
	}

	// Publicar evento
	uc.publishEvent(ctx, "inventory.released", orderID, businessID, whID, result)

	return result, nil
}
