package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// TransferStockTx ejecuta la transferencia de stock entre bodegas dentro de una transacción con SELECT FOR UPDATE.
// Bloquea ambos niveles en orden ascendente de WarehouseID para evitar deadlocks.
// Garantiza atomicidad: ambos niveles y ambos movimientos se modifican en una sola transacción.
func (r *Repository) TransferStockTx(ctx context.Context, params dtos.TransferStockTxParams) (*dtos.TransferStockTxResult, error) {
	var result dtos.TransferStockTxResult

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		// Ordenar por WarehouseID para evitar deadlocks:
		// Si dos requests concurrentes transfieren A→B y B→A, ambos bloquean
		// en el mismo orden (menor ID primero), evitando circular wait.
		firstWH, secondWH := params.FromWarehouseID, params.ToWarehouseID
		firstLoc, secondLoc := params.FromLocationID, params.ToLocationID
		swapped := false
		if firstWH > secondWH {
			firstWH, secondWH = secondWH, firstWH
			firstLoc, secondLoc = secondLoc, firstLoc
			swapped = true
		}

		// 1. SELECT FOR UPDATE del primer nivel (menor WarehouseID)
		firstLevel, err := r.getOrCreateLevelTx(tx, params.ProductID, firstWH, firstLoc, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx first: %w", err)
		}

		// 2. SELECT FOR UPDATE del segundo nivel (mayor WarehouseID)
		secondLevel, err := r.getOrCreateLevelTx(tx, params.ProductID, secondWH, secondLoc, params.BusinessID)
		if err != nil {
			return fmt.Errorf("getOrCreateLevelTx second: %w", err)
		}

		// Determinar cuál es origen y cuál destino
		var fromLevel, toLevel *models.InventoryLevel
		if swapped {
			fromLevel = secondLevel // FromWarehouseID era el mayor
			toLevel = firstLevel
		} else {
			fromLevel = firstLevel
			toLevel = secondLevel
		}

		// 3. Validar stock suficiente en origen
		if fromLevel.Quantity < params.Quantity {
			return domainerrors.ErrInsufficientStock
		}

		fromPrev := fromLevel.Quantity
		toPrev := toLevel.Quantity

		// 4. UPDATE ambos niveles
		fromLevel.Quantity -= params.Quantity
		fromLevel.AvailableQty = fromLevel.Quantity - fromLevel.ReservedQty
		toLevel.Quantity += params.Quantity
		toLevel.AvailableQty = toLevel.Quantity - toLevel.ReservedQty

		if err := r.updateLevelTx(tx, fromLevel); err != nil {
			return fmt.Errorf("updateLevelTx from: %w", err)
		}
		if err := r.updateLevelTx(tx, toLevel); err != nil {
			return fmt.Errorf("updateLevelTx to: %w", err)
		}

		// 5. INSERT 2 stock_movements
		refType := params.ReferenceType

		// Movimiento de salida (origen)
		outMovement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.FromWarehouseID,
			LocationID:     params.FromLocationID,
			BusinessID:     params.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         params.Reason,
			Quantity:       -params.Quantity,
			PreviousQty:    fromPrev,
			NewQty:         fromLevel.Quantity,
			ReferenceType:  &refType,
			Notes:          params.Notes,
			CreatedByID:    params.CreatedByID,
		}
		if err := r.createMovementTx(tx, outMovement); err != nil {
			return fmt.Errorf("createMovementTx out: %w", err)
		}

		// Movimiento de entrada (destino)
		inMovement := &models.StockMovement{
			ProductID:      params.ProductID,
			WarehouseID:    params.ToWarehouseID,
			LocationID:     params.ToLocationID,
			BusinessID:     params.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         params.Reason,
			Quantity:       params.Quantity,
			PreviousQty:    toPrev,
			NewQty:         toLevel.Quantity,
			ReferenceType:  &refType,
			Notes:          params.Notes,
			CreatedByID:    params.CreatedByID,
		}
		if err := r.createMovementTx(tx, inMovement); err != nil {
			return fmt.Errorf("createMovementTx in: %w", err)
		}

		result = dtos.TransferStockTxResult{
			FromNewQty: fromLevel.Quantity,
			ToNewQty:   toLevel.Quantity,
		}

		return nil // commit automático
	})

	if err != nil {
		return nil, err
	}

	// Cache invalidation después del commit (fire-and-forget)
	if r.cache != nil {
		go r.cache.InvalidateProduct(context.Background(), params.ProductID, params.BusinessID)
		go r.cache.InvalidateLevel(context.Background(), params.ProductID, params.FromWarehouseID)
		go r.cache.InvalidateLevel(context.Background(), params.ProductID, params.ToWarehouseID)
	}

	return &result, nil
}
