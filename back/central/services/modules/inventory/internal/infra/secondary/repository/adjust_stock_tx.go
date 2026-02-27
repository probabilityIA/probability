package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// AdjustStockTx ejecuta el ajuste de stock completo dentro de una transacción con SELECT FOR UPDATE.
// Garantiza atomicidad: el nivel y el movimiento se crean/actualizan en una sola transacción.
// El cache se invalida después del commit (fire-and-forget).
func (r *Repository) AdjustStockTx(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
	var result dtos.AdjustStockTxResult

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. SELECT FOR UPDATE del inventory_level (o crear si no existe)
		level, err := r.getOrCreateLevelTx(tx, params.ProductID, params.WarehouseID, params.LocationID, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx: %w", err)
		}

		// 2. Validar que no quede negativo
		newQty := level.Quantity + params.Quantity
		if newQty < 0 {
			return domainerrors.ErrInsufficientStock
		}

		previousQty := level.Quantity

		// 3. UPDATE inventory_level dentro del tx
		level.Quantity = newQty
		level.AvailableQty = newQty - level.ReservedQty
		if err := r.updateLevelTx(tx, level); err != nil {
			return fmt.Errorf("updateLevelTx: %w", err)
		}

		// 4. INSERT stock_movement dentro del tx
		movement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.WarehouseID,
			LocationID:     params.LocationID,
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

		// Construir resultado
		levelEntity := inventoryLevelModelToEntity(level)
		result = dtos.AdjustStockTxResult{
			Movement:    stockMovementModelToEntity(movement),
			NewQuantity: newQty,
			Level:       levelEntity,
		}

		return nil // commit automático
	})

	if err != nil {
		return nil, err
	}

	// Cache invalidation después del commit (fire-and-forget)
	if r.cache != nil {
		go r.cache.InvalidateProduct(context.Background(), params.ProductID, params.BusinessID)
		go r.cache.InvalidateLevel(context.Background(), params.ProductID, params.WarehouseID)
	}

	return &result, nil
}
