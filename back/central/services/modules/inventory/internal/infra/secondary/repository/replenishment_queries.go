package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateReplenishmentTask(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
	m := mappers.ReplenishmentTaskEntityToModel(t)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.ReplenishmentTaskModelToEntity(m), nil
}

func (r *Repository) GetReplenishmentTaskByID(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
	var m models.ReplenishmentTask
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrReplenishmentNotFound
		}
		return nil, err
	}
	return mappers.ReplenishmentTaskModelToEntity(&m), nil
}

func (r *Repository) ListReplenishmentTasks(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error) {
	var ml []models.ReplenishmentTask
	var total int64
	q := r.db.Conn(ctx).Model(&models.ReplenishmentTask{}).Where("business_id = ?", params.BusinessID)
	if params.WarehouseID != nil {
		q = q.Where("warehouse_id = ?", *params.WarehouseID)
	}
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if params.AssignedTo != nil {
		q = q.Where("assigned_to_id = ?", *params.AssignedTo)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.ReplenishmentTask, len(ml))
	for i := range ml {
		out[i] = *mappers.ReplenishmentTaskModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdateReplenishmentTask(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
	updates := map[string]any{
		"from_location_id": t.FromLocationID,
		"to_location_id":   t.ToLocationID,
		"quantity":         t.Quantity,
		"status":           t.Status,
		"assigned_to_id":   t.AssignedToID,
		"assigned_at":      t.AssignedAt,
		"completed_at":     t.CompletedAt,
		"notes":            t.Notes,
	}
	res := r.db.Conn(ctx).Model(&models.ReplenishmentTask{}).
		Where("id = ? AND business_id = ?", t.ID, t.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrReplenishmentNotFound
	}
	return r.GetReplenishmentTaskByID(ctx, t.BusinessID, t.ID)
}

func (r *Repository) DetectReplenishmentCandidates(ctx context.Context, businessID uint) ([]entities.ReplenishmentTask, error) {
	type row struct {
		ProductID   string
		WarehouseID uint
		LocationID  *uint
		Deficit     int
	}
	var results []row

	err := r.db.Conn(ctx).
		Table("inventory_levels il").
		Select("il.product_id, il.warehouse_id, il.location_id, (COALESCE(il.reorder_point, il.min_stock, 0) - il.available_qty) AS deficit").
		Where("il.business_id = ? AND il.deleted_at IS NULL", businessID).
		Where("COALESCE(il.reorder_point, il.min_stock, 0) > 0").
		Where("il.available_qty < COALESCE(il.reorder_point, il.min_stock)").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	tasks := make([]entities.ReplenishmentTask, 0, len(results))
	for _, row := range results {
		tasks = append(tasks, entities.ReplenishmentTask{
			BusinessID:   businessID,
			ProductID:    row.ProductID,
			WarehouseID:  row.WarehouseID,
			ToLocationID: row.LocationID,
			Quantity:     row.Deficit,
			Status:       "pending",
			TriggeredBy:  "auto",
		})
	}
	return tasks, nil
}
