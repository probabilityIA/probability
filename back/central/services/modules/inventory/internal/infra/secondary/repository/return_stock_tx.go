package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ReturnStockTx devuelve stock por reembolso: Quantity += qty, AvailableQty = Quantity - ReservedQty.
// Usa movement type "return" (ID 5).
func (r *Repository) ReturnStockTx(ctx context.Context, params dtos.ReturnStockTxParams) error {
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. SELECT FOR UPDATE
		level, err := r.getOrCreateLevelTx(tx, params.ProductID, params.WarehouseID, nil, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx: %w", err)
		}

		previousQty := level.Quantity

		// 2. UPDATE: incrementar stock real
		level.Quantity += params.Quantity
		level.AvailableQty = level.Quantity - level.ReservedQty

		if err := r.updateLevelTx(tx, level); err != nil {
			return fmt.Errorf("updateLevelTx: %w", err)
		}

		// 3. INSERT stock_movement (in: Quantity positivo)
		refType := "order"
		movement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.WarehouseID,
			BusinessID:     params.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         fmt.Sprintf("Devolución por reembolso de orden %s", params.OrderID),
			Quantity:       params.Quantity,
			PreviousQty:    previousQty,
			NewQty:         level.Quantity,
			ReferenceType:  &refType,
			ReferenceID:    &params.OrderID,
			Notes:          fmt.Sprintf("Devuelto: %d, Stock: %d→%d", params.Quantity, previousQty, level.Quantity),
		}
		if err := r.createMovementTx(tx, movement); err != nil {
			return fmt.Errorf("createMovementTx: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Cache invalidation
	if r.cache != nil {
		go r.cache.InvalidateProduct(context.Background(), params.ProductID, params.BusinessID)
		go r.cache.InvalidateLevel(context.Background(), params.ProductID, params.WarehouseID)
	}

	return nil
}
