package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ReserveStockTx reserva stock para una orden dentro de una transacción con SELECT FOR UPDATE.
// ReservedQty += qty, AvailableQty = Quantity - ReservedQty.
// Si no hay stock suficiente, reserva parcial (lo que haya disponible).
func (r *Repository) ReserveStockTx(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
	var result dtos.ReserveStockTxResult

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. SELECT FOR UPDATE del inventory_level (o crear si no existe)
		level, err := r.getOrCreateLevelTx(tx, params.ProductID, params.WarehouseID, nil, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx: %w", err)
		}

		result.PreviousAvailable = level.AvailableQty

		// 2. Calcular cuánto se puede reservar
		toReserve := params.Quantity
		if level.AvailableQty < toReserve {
			toReserve = level.AvailableQty
			if toReserve < 0 {
				toReserve = 0
			}
			result.Sufficient = false
		} else {
			result.Sufficient = true
		}
		result.Reserved = toReserve

		if toReserve == 0 {
			// Nada que reservar, pero no es error — insuficiente
			result.NewAvailable = level.AvailableQty
			result.NewReserved = level.ReservedQty
			return nil
		}

		// 3. UPDATE inventory_level
		level.ReservedQty += toReserve
		level.AvailableQty = level.Quantity - level.ReservedQty
		if err := r.updateLevelTx(tx, level); err != nil {
			return fmt.Errorf("updateLevelTx: %w", err)
		}

		result.NewAvailable = level.AvailableQty
		result.NewReserved = level.ReservedQty

		// 4. INSERT stock_movement (neutral: no cambia Quantity real)
		refType := "order"
		movement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.WarehouseID,
			BusinessID:     params.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         fmt.Sprintf("Reserva por orden %s", params.OrderID),
			Quantity:       0, // neutral — no afecta Quantity total
			PreviousQty:    level.Quantity,
			NewQty:         level.Quantity,
			ReferenceType:  &refType,
			ReferenceID:    &params.OrderID,
			Notes:          fmt.Sprintf("Reservado: %d, Disponible: %d→%d", toReserve, result.PreviousAvailable, level.AvailableQty),
		}
		if err := r.createMovementTx(tx, movement); err != nil {
			return fmt.Errorf("createMovementTx: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Cache invalidation después del commit
	if r.cache != nil {
		go r.cache.InvalidateProduct(context.Background(), params.ProductID, params.BusinessID)
		go r.cache.InvalidateLevel(context.Background(), params.ProductID, params.WarehouseID)
	}

	return &result, nil
}
