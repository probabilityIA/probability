package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
)

func (uc *useCase) ReleaseStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error) {
	whID, err := uc.resolveWarehouse(ctx, warehouseID, businessID)
	if err != nil {
		return nil, err
	}

	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "release")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get release movement type")
		return nil, err
	}

	result := &response.OrderStockResult{
		OrderID:     orderID,
		BusinessID:  businessID,
		WarehouseID: whID,
		Success:     true,
	}

	for _, item := range items {
		itemResult := response.ItemStockResult{
			ProductID: item.ProductID,
			SKU:       item.SKU,
			Requested: item.Quantity,
		}

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

		uc.updateProductTotalStock(ctx, item.ProductID, businessID)
	}

	uc.publishEvent(ctx, "inventory.released", orderID, businessID, whID, result)

	return result, nil
}
