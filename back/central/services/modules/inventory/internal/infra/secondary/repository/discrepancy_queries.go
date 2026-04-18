package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateDiscrepancy(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
	m := mappers.DiscrepancyEntityToModel(d)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.DiscrepancyModelToEntity(m), nil
}

func (r *Repository) GetDiscrepancyByID(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
	var m models.InventoryDiscrepancy
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrDiscrepancyNotFound
		}
		return nil, err
	}
	return mappers.DiscrepancyModelToEntity(&m), nil
}

func (r *Repository) ListDiscrepancies(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error) {
	var ml []models.InventoryDiscrepancy
	var total int64
	q := r.db.Conn(ctx).Model(&models.InventoryDiscrepancy{}).Where("business_id = ?", params.BusinessID)
	if params.TaskID != nil {
		q = q.Where("task_id = ?", *params.TaskID)
	}
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.InventoryDiscrepancy, len(ml))
	for i := range ml {
		out[i] = *mappers.DiscrepancyModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdateDiscrepancy(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
	updates := map[string]any{
		"status":                 d.Status,
		"resolution_movement_id": d.ResolutionMovementID,
		"reviewed_by_id":         d.ReviewedByID,
		"reviewed_at":            d.ReviewedAt,
		"notes":                  d.Notes,
	}
	res := r.db.Conn(ctx).Model(&models.InventoryDiscrepancy{}).
		Where("id = ? AND business_id = ?", d.ID, d.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrDiscrepancyNotFound
	}
	return r.GetDiscrepancyByID(ctx, d.BusinessID, d.ID)
}

func (r *Repository) ApproveDiscrepancyTx(ctx context.Context, params dtos.ApproveDiscrepancyTxParams) (*entities.InventoryDiscrepancy, error) {
	var result *entities.InventoryDiscrepancy

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		var disc models.InventoryDiscrepancy
		if err := tx.Where("id = ? AND business_id = ?", params.DiscrepancyID, params.BusinessID).First(&disc).Error; err != nil {
			return domainerrors.ErrDiscrepancyNotFound
		}
		if disc.Status == "approved" || disc.Status == "rejected" {
			return domainerrors.ErrDiscrepancyResolved
		}

		var line models.CycleCountLine
		if err := tx.First(&line, disc.LineID).Error; err != nil {
			return err
		}
		if line.CountedQty == nil {
			return domainerrors.ErrCountLineNotFound
		}

		delta := *line.CountedQty - line.ExpectedQty

		var level models.InventoryLevel
		lq := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("business_id = ? AND product_id = ?", disc.BusinessID, line.ProductID)
		if line.LocationID != nil {
			lq = lq.Where("location_id = ?", *line.LocationID)
		} else {
			lq = lq.Where("location_id IS NULL")
		}
		if line.LotID != nil {
			lq = lq.Where("lot_id = ?", *line.LotID)
		} else {
			lq = lq.Where("lot_id IS NULL")
		}
		if err := lq.First(&level).Error; err != nil {
			return err
		}

		previousQty := level.Quantity
		newQty := previousQty + delta
		if newQty < 0 {
			newQty = 0
		}
		level.Quantity = newQty
		level.AvailableQty = newQty - level.ReservedQty
		if err := tx.Model(&level).Updates(map[string]any{
			"quantity":      level.Quantity,
			"available_qty": level.AvailableQty,
		}).Error; err != nil {
			return err
		}

		refType := "count_adjustment"
		refID := ""
		movement := &models.StockMovement{
			ProductID:      line.ProductID,
			WarehouseID:    level.WarehouseID,
			LocationID:     line.LocationID,
			LotID:          line.LotID,
			BusinessID:     disc.BusinessID,
			MovementTypeID: params.MovementTypeID,
			Reason:         "count_adjustment",
			Quantity:       delta,
			PreviousQty:    previousQty,
			NewQty:         newQty,
			ReferenceType:  &refType,
			ReferenceID:    &refID,
			Notes:          params.Notes,
			CreatedByID:    &params.ReviewerID,
		}
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		now := time.Now()
		disc.Status = "approved"
		disc.ResolutionMovementID = &movement.ID
		disc.ReviewedByID = &params.ReviewerID
		disc.ReviewedAt = &now
		if params.Notes != "" {
			disc.Notes = params.Notes
		}
		if err := tx.Save(&disc).Error; err != nil {
			return err
		}

		line.Status = "resolved"
		if err := tx.Save(&line).Error; err != nil {
			return err
		}

		result = mappers.DiscrepancyModelToEntity(&disc)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
