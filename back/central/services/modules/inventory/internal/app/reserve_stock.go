package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
)

// ReserveStockForOrder reserva stock para una orden nueva.
// Para cada item: ReservedQty += qty, AvailableQty -= qty.
// Si no hay stock suficiente, reserva parcial + publica evento "inventory.insufficient".
func (uc *useCase) ReserveStockForOrder(ctx context.Context, orderID string, businessID uint, warehouseID *uint, items []dtos.OrderInventoryItem) (*dtos.OrderStockResult, error) {
	// Resolver warehouse
	whID, err := uc.resolveWarehouse(ctx, warehouseID, businessID)
	if err != nil {
		return nil, err
	}

	// Obtener movement type
	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "reserve")
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Failed to get reserve movement type")
		return nil, err
	}

	result := &dtos.OrderStockResult{
		OrderID:     orderID,
		BusinessID:  businessID,
		WarehouseID: whID,
		Success:     true,
	}

	allSufficient := true

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
			itemResult.Sufficient = false
			result.ItemResults = append(result.ItemResults, itemResult)
			continue
		}
		if !trackInventory {
			// Producto sin tracking → skip (no afecta inventario)
			itemResult.Processed = item.Quantity
			itemResult.Sufficient = true
			result.ItemResults = append(result.ItemResults, itemResult)
			continue
		}

		// Ejecutar reserva transaccional
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

		// Actualizar stock total del producto (best-effort)
		uc.updateProductTotalStock(ctx, item.ProductID, businessID)
	}

	// Publicar evento de inventario (fire-and-forget)
	eventType := "inventory.reserved"
	if !allSufficient {
		eventType = "inventory.insufficient"
	}
	uc.publishEvent(ctx, eventType, orderID, businessID, whID, result)

	return result, nil
}

// resolveWarehouse determina la bodega a usar: explícita o default
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

// publishEvent publica un evento de inventario a Redis SSE (fire-and-forget)
func (uc *useCase) publishEvent(ctx context.Context, eventType string, orderID string, businessID uint, warehouseID uint, data interface{}) {
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
			Data: map[string]interface{}{
				"result": data,
			},
		})
	}()
}
