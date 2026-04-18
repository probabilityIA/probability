package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
)

func (uc *useCase) ReserveStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*response.OrderStockResult, error) {
	whID, err := uc.resolveWarehouse(ctx, warehouseID, businessID)
	if err != nil {
		return nil, err
	}

	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "reserve")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get reserve movement type")
		return nil, err
	}

	result := &response.OrderStockResult{
		OrderID:     orderID,
		BusinessID:  businessID,
		WarehouseID: whID,
		Success:     true,
	}

	allSufficient := true

	for _, item := range items {
		itemResult := response.ItemStockResult{
			ProductID: item.ProductID,
			SKU:       item.SKU,
			Requested: item.Quantity,
		}

		_, _, trackInventory, err := uc.repo.GetProductByID(ctx, item.ProductID, businessID)
		if err != nil {
			itemResult.ErrorMessage = "producto no encontrado"
			itemResult.Sufficient = false
			result.ItemResults = append(result.ItemResults, itemResult)
			continue
		}
		if !trackInventory {
			itemResult.Processed = item.Quantity
			itemResult.Sufficient = true
			result.ItemResults = append(result.ItemResults, itemResult)
			continue
		}

		txResult, err := uc.repo.ReserveStockTx(ctx, dtos.ReserveStockTxParams{
			ProductID:      item.ProductID,
			WarehouseID:    whID,
			BusinessID:     businessID,
			Quantity:       item.Quantity,
			MovementTypeID: movTypeID,
			OrderID:        orderID,
		})
		if err != nil {
			itemResult.ErrorMessage = err.Error()
			itemResult.Sufficient = false
			result.ItemResults = append(result.ItemResults, itemResult)
			result.Success = false
			continue
		}

		itemResult.Processed = txResult.Reserved
		itemResult.Sufficient = txResult.Sufficient
		if !txResult.Sufficient {
			allSufficient = false
		}
		result.ItemResults = append(result.ItemResults, itemResult)

		uc.updateProductTotalStock(ctx, item.ProductID, businessID)
	}

	eventType := "inventory.reserved"
	if !allSufficient {
		eventType = "inventory.insufficient"
	}
	uc.publishEvent(ctx, eventType, orderID, businessID, whID, result)

	return result, nil
}

func (uc *useCase) resolveWarehouse(ctx context.Context, warehouseID *uint, businessID uint) (uint, error) {
	if warehouseID != nil && *warehouseID > 0 {
		return *warehouseID, nil
	}

	whID, err := uc.repo.GetDefaultWarehouseID(ctx, businessID)
	if err != nil {
		uc.log.Warn(ctx).Uint("business_id", businessID).Msg("No default warehouse found")
		return 0, domainerrors.ErrNoDefaultWarehouse
	}
	return whID, nil
}

func (uc *useCase) publishEvent(ctx context.Context, eventType string, orderID string, businessID uint, warehouseID uint, data any) {
	if uc.eventPublisher == nil {
		return
	}

	go func() {
		_ = uc.eventPublisher.PublishInventoryEvent(context.Background(), ports.InventoryEvent{
			EventType:   eventType,
			OrderID:     orderID,
			BusinessID:  businessID,
			WarehouseID: warehouseID,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Data: map[string]any{
				"result": data,
			},
		})
	}()
}
