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

func (r *Repository) CreateCountLine(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
	m := mappers.CountLineEntityToModel(line)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.CountLineModelToEntity(m), nil
}

func (r *Repository) GetCountLineByID(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
	var m models.CycleCountLine
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrCountLineNotFound
		}
		return nil, err
	}
	return mappers.CountLineModelToEntity(&m), nil
}

func (r *Repository) ListCountLines(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error) {
	var ml []models.CycleCountLine
	var total int64
	q := r.db.Conn(ctx).Model(&models.CycleCountLine{}).Where("business_id = ? AND task_id = ?", params.BusinessID, params.TaskID)
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id ASC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.CycleCountLine, len(ml))
	for i := range ml {
		out[i] = *mappers.CountLineModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdateCountLine(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
	updates := map[string]any{
		"counted_qty": line.CountedQty,
		"variance":    line.Variance,
		"status":      line.Status,
	}
	res := r.db.Conn(ctx).Model(&models.CycleCountLine{}).
		Where("id = ? AND business_id = ?", line.ID, line.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrCountLineNotFound
	}
	return r.GetCountLineByID(ctx, line.BusinessID, line.ID)
}
