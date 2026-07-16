package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) AdjustStockTx(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
	var result dtos.AdjustStockTxResult

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		level, err := r.getOrCreateLevelKeyTx(tx, params.ProductID, params.WarehouseID, params.LocationID, params.LotID, params.StateID, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx: %w", err)
		}

		newQty := level.Quantity + params.Quantity
		previousQty := level.Quantity

		level.Quantity = newQty
		level.AvailableQty = newQty - level.ReservedQty
		if err := r.updateLevelTx(tx, level); err != nil {
			return fmt.Errorf("updateLevelTx: %w", err)
		}

		movement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.WarehouseID,
			LocationID:     params.LocationID,
			LotID:          params.LotID,
			ToStateID:      params.StateID,
			UomID:          params.UomID,
			BusinessID:     params.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         params.Reason,
			Quantity:       params.Quantity,
			PreviousQty:    previousQty,
			NewQty:         newQty,
			ReferenceType:  &params.ReferenceType,
			Notes:          params.Notes,
			CreatedByID:    params.CreatedByID,
		}
		if err := r.createMovementTx(tx, movement); err != nil {
			return fmt.Errorf("createMovementTx: %w", err)
		}

		if err := tx.Exec(
			"UPDATE products SET stock_quantity = COALESCE((SELECT SUM(quantity) FROM inventory_levels WHERE product_id = ? AND deleted_at IS NULL), 0) WHERE id = ?",
			params.ProductID, params.ProductID,
		).Error; err != nil {
			return fmt.Errorf("update product stock_quantity: %w", err)
		}

		levelEntity := mappers.LevelModelToEntity(level)
		result = dtos.AdjustStockTxResult{
			Movement:    mappers.MovementModelToEntity(movement),
			NewQuantity: newQty,
			Level:       levelEntity,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if r.cache != nil {
		go r.cache.InvalidateProduct(context.Background(), params.ProductID, params.BusinessID)
		go r.cache.InvalidateLevel(context.Background(), params.ProductID, params.WarehouseID)
	}

	return &result, nil
}
