package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ConfirmSaleTx confirma la venta (shipped/completed): Quantity -= qty, ReservedQty -= qty.
// Clamp: si ReservedQty < qty, solo decrementa lo que había reservado.
func (r *Repository) ConfirmSaleTx(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. SELECT FOR UPDATE
		level, err := r.getOrCreateLevelTx(tx, params.ProductID, params.WarehouseID, nil, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx: %w", err)
		}

		previousQty := level.Quantity

		// 2. Clamp reservedQty a lo que realmente está reservado
		toRelease := params.Quantity
		if toRelease > level.ReservedQty {
			toRelease = level.ReservedQty
		}

		// 3. UPDATE: reducir stock real y reserva
		level.Quantity -= params.Quantity
		level.ReservedQty -= toRelease
		if level.ReservedQty < 0 {
			level.ReservedQty = 0
		}
		level.AvailableQty = level.Quantity - level.ReservedQty

		if err := r.updateLevelTx(tx, level); err != nil {
			return fmt.Errorf("updateLevelTx: %w", err)
		}

		// 4. INSERT stock_movement (out: Quantity negativo)
		refType := "order"
		movement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.WarehouseID,
			BusinessID:     params.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         fmt.Sprintf("Venta confirmada por orden %s", params.OrderID),
			Quantity:       -params.Quantity,
			PreviousQty:    previousQty,
			NewQty:         level.Quantity,
			ReferenceType:  &refType,
			ReferenceID:    &params.OrderID,
			Notes:          fmt.Sprintf("Confirmado: %d unidades, Reserva liberada: %d", params.Quantity, toRelease),
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
