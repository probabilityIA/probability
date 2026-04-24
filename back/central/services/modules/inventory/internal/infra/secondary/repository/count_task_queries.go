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

func (r *Repository) CreateCountTask(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
	m := mappers.CountTaskEntityToModel(t)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.CountTaskModelToEntity(m), nil
}

func (r *Repository) GetCountTaskByID(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
	var m models.CycleCountTask
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrCountTaskNotFound
		}
		return nil, err
	}
	return mappers.CountTaskModelToEntity(&m), nil
}

func (r *Repository) ListCountTasks(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error) {
	var ml []models.CycleCountTask
	var total int64
	q := r.db.Conn(ctx).Model(&models.CycleCountTask{}).Where("business_id = ?", params.BusinessID)
	if params.WarehouseID != nil {
		q = q.Where("warehouse_id = ?", *params.WarehouseID)
	}
	if params.PlanID != nil {
		q = q.Where("plan_id = ?", *params.PlanID)
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
	out := make([]entities.CycleCountTask, len(ml))
	for i := range ml {
		out[i] = *mappers.CountTaskModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdateCountTask(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
	updates := map[string]any{
		"scope_type":    t.ScopeType,
		"scope_id":      t.ScopeID,
		"status":        t.Status,
		"assigned_to_id": t.AssignedToID,
		"started_at":    t.StartedAt,
		"finished_at":   t.FinishedAt,
	}
	res := r.db.Conn(ctx).Model(&models.CycleCountTask{}).
		Where("id = ? AND business_id = ?", t.ID, t.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrCountTaskNotFound
	}
	return r.GetCountTaskByID(ctx, t.BusinessID, t.ID)
}

func (r *Repository) GenerateCountLinesForTask(ctx context.Context, task *entities.CycleCountTask, strategy string) ([]entities.CycleCountLine, error) {
	type row struct {
		ProductID  string
		LocationID *uint
		LotID      *uint
		Quantity   int
	}
	var rows []row

	q := r.db.Conn(ctx).
		Table("inventory_levels il").
		Select("il.product_id, il.location_id, il.lot_id, il.quantity").
		Where("il.business_id = ? AND il.warehouse_id = ? AND il.deleted_at IS NULL", task.BusinessID, task.WarehouseID)

	if strategy == "abc" {
		q = q.Joins("INNER JOIN product_velocities pv ON pv.product_id = il.product_id AND pv.warehouse_id = il.warehouse_id AND pv.deleted_at IS NULL").
			Where("pv.rank = 'A'")
	} else if strategy == "zone" && task.ScopeID != nil {
		q = q.Joins("INNER JOIN warehouse_locations wl ON wl.id = il.location_id AND wl.deleted_at IS NULL").
			Joins("INNER JOIN warehouse_rack_levels wrl ON wrl.id = wl.level_id AND wrl.deleted_at IS NULL").
			Joins("INNER JOIN warehouse_racks wr ON wr.id = wrl.rack_id AND wr.deleted_at IS NULL").
			Joins("INNER JOIN warehouse_aisles wa ON wa.id = wr.aisle_id AND wa.deleted_at IS NULL").
			Where("wa.zone_id = ?", *task.ScopeID)
	} else if strategy == "random" {
		q = q.Order("RANDOM()").Limit(50)
	}

	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	lines := make([]entities.CycleCountLine, 0, len(rows))
	for _, row := range rows {
		line := &entities.CycleCountLine{
			TaskID:      task.ID,
			BusinessID:  task.BusinessID,
			ProductID:   row.ProductID,
			LocationID:  row.LocationID,
			LotID:       row.LotID,
			ExpectedQty: row.Quantity,
			Status:      "pending",
		}
		m := mappers.CountLineEntityToModel(line)
		if err := r.db.Conn(ctx).Create(m).Error; err != nil {
			continue
		}
		lines = append(lines, *mappers.CountLineModelToEntity(m))
	}
	return lines, nil
}
