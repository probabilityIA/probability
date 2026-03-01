package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
)

// ReturnStockForOrder devuelve stock cuando una orden es reembolsada.
// Quantity += qty, AvailableQty = Quantity - ReservedQty.
func (uc *useCase) ReturnStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error) {
	// Resolver warehouse
	whID, err := uc.resolveWarehouse(ctx, warehouseID, businessID)
	if err != nil {
		return nil, err
	}

	// Obtener movement type "return" (ID 5)
	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "return")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get return movement type")
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

		// Ejecutar devolución transaccional
		err = uc.repo.ReturnStockTx(ctx, dtos.ReturnStockTxParams{
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

		// Publicar sync a canales de venta
		uc.publishSync(ctx, item.ProductID, businessID, 0, whID, "order_return")
	}

	// Publicar evento
	uc.publishEvent(ctx, "inventory.returned", orderID, businessID, whID, result)

	return result, nil
}
