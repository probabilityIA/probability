package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) ChangeStateTx(ctx context.Context, params dtos.ChangeInventoryStateTxParams) (*entities.StockMovement, error) {
	var result *entities.StockMovement

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		var fromState models.InventoryState
		if err := tx.Where("code = ?", params.FromStateCode).First(&fromState).Error; err != nil {
			return domainerrors.ErrStateNotFound
		}
		var toState models.InventoryState
		if err := tx.Where("code = ?", params.ToStateCode).First(&toState).Error; err != nil {
			return domainerrors.ErrStateNotFound
		}

		var level models.InventoryLevel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND business_id = ?", params.LevelID, params.BusinessID).
			First(&level).Error; err != nil {
			return fmt.Errorf("level not found: %w", err)
		}

		if level.Quantity < params.Quantity {
			return domainerrors.ErrInsufficientStock
		}

		level.Quantity -= params.Quantity
		level.AvailableQty = level.Quantity - level.ReservedQty
		if err := tx.Model(&level).Updates(map[string]any{
			"quantity":      level.Quantity,
			"available_qty": level.AvailableQty,
		}).Error; err != nil {
			return err
		}

		var destLevel models.InventoryLevel
		selQuery := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ? AND business_id = ? AND state_id = ?",
				level.ProductID, level.WarehouseID, level.BusinessID, toState.ID)
		if level.LocationID != nil {
			selQuery = selQuery.Where("location_id = ?", *level.LocationID)
		} else {
			selQuery = selQuery.Where("location_id IS NULL")
		}
		if level.LotID != nil {
			selQuery = selQuery.Where("lot_id = ?", *level.LotID)
		} else {
			selQuery = selQuery.Where("lot_id IS NULL")
		}

		if err := selQuery.First(&destLevel).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			destLevel = models.InventoryLevel{
				ProductID:    level.ProductID,
				WarehouseID:  level.WarehouseID,
				LocationID:   level.LocationID,
				LotID:        level.LotID,
				StateID:      &toState.ID,
				BusinessID:   level.BusinessID,
				Quantity:     params.Quantity,
				ReservedQty:  0,
				AvailableQty: params.Quantity,
			}
			if err := tx.Create(&destLevel).Error; err != nil {
				return err
			}
		} else {
			destLevel.Quantity += params.Quantity
			destLevel.AvailableQty = destLevel.Quantity - destLevel.ReservedQty
			if err := tx.Model(&destLevel).Updates(map[string]any{
				"quantity":      destLevel.Quantity,
				"available_qty": destLevel.AvailableQty,
			}).Error; err != nil {
				return err
			}
		}

		var movType models.StockMovementType
		_ = tx.Where("code = ?", "state_change").First(&movType).Error
		var movTypeID uint
		if movType.ID > 0 {
			movTypeID = movType.ID
		} else {
			var fallback models.StockMovementType
			_ = tx.Where("code = ?", "adjustment").First(&fallback).Error
			movTypeID = fallback.ID
		}

		refType := "state_change"
		movement := &models.StockMovement{
			ProductID:      level.ProductID,
			WarehouseID:    level.WarehouseID,
			LocationID:     level.LocationID,
			LotID:          level.LotID,
			FromStateID:    &fromState.ID,
			ToStateID:      &toState.ID,
			BusinessID:     level.BusinessID,
			MovementTypeID: movTypeID,
			Reason:         params.Reason,
			Quantity:       params.Quantity,
			PreviousQty:    level.Quantity + params.Quantity,
			NewQty:         level.Quantity,
			ReferenceType:  &refType,
			CreatedByID:    params.CreatedByID,
		}
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		result = mappers.MovementModelToEntity(movement)
		return nil
	})

	return result, err
}
