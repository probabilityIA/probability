package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ReleaseStockTx libera reserva por cancelación: ReservedQty -= qty.
// Clamp: si ReservedQty < qty, solo libera lo que había reservado.
func (r *Repository) ReleaseStockTx(ctx context.Context, params dtos.ReleaseTxParams) error {
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. SELECT FOR UPDATE
		level, err := r.getOrCreateLevelTx(tx, params.ProductID, params.WarehouseID, nil, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx: %w", err)
		}

		// 2. Clamp a lo que realmente está reservado
		toRelease := params.Quantity
		if toRelease > level.ReservedQty {
			toRelease = level.ReservedQty
		}
		if toRelease <= 0 {
			return nil // nada que liberar
		}

		// 3. UPDATE: solo ReservedQty baja, Quantity no cambia
		level.ReservedQty -= toRelease
		level.AvailableQty = level.Quantity - level.ReservedQty

		if err := r.updateLevelTx(tx, level); err != nil {
			return fmt.Errorf("updateLevelTx: %w", err)
		}

		// 4. INSERT stock_movement (neutral: no cambia Quantity)
		refType := "order"
		movement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.WarehouseID,
			BusinessID:     params.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         fmt.Sprintf("Liberación de reserva por cancelación de orden %s", params.OrderID),
			Quantity:       0,
			PreviousQty:    level.Quantity,
			NewQty:         level.Quantity,
			ReferenceType:  &refType,
			ReferenceID:    &params.OrderID,
			Notes:          fmt.Sprintf("Liberado: %d, Reserva: %d→%d", toRelease, toRelease+level.ReservedQty, level.ReservedQty),
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
