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

func (r *Repository) CreateCountPlan(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
	m := mappers.CountPlanEntityToModel(p)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.CountPlanModelToEntity(m), nil
}

func (r *Repository) GetCountPlanByID(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
	var m models.CycleCountPlan
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrCountPlanNotFound
		}
		return nil, err
	}
	return mappers.CountPlanModelToEntity(&m), nil
}

func (r *Repository) ListCountPlans(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error) {
	var ml []models.CycleCountPlan
	var total int64
	q := r.db.Conn(ctx).Model(&models.CycleCountPlan{}).Where("business_id = ?", params.BusinessID)
	if params.WarehouseID != nil {
		q = q.Where("warehouse_id = ?", *params.WarehouseID)
	}
	if params.ActiveOnly {
		q = q.Where("is_active = true")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.CycleCountPlan, len(ml))
	for i := range ml {
		out[i] = *mappers.CountPlanModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdateCountPlan(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
	updates := map[string]any{
		"warehouse_id":   p.WarehouseID,
		"name":           p.Name,
		"strategy":       p.Strategy,
		"frequency_days": p.FrequencyDays,
		"next_run_at":    p.NextRunAt,
		"is_active":      p.IsActive,
	}
	res := r.db.Conn(ctx).Model(&models.CycleCountPlan{}).
		Where("id = ? AND business_id = ?", p.ID, p.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrCountPlanNotFound
	}
	return r.GetCountPlanByID(ctx, p.BusinessID, p.ID)
}

func (r *Repository) DeleteCountPlan(ctx context.Context, businessID, id uint) error {
	res := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).Delete(&models.CycleCountPlan{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domainerrors.ErrCountPlanNotFound
	}
	return nil
}
