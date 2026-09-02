package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) Create(ctx context.Context, sprint *entities.Sprint) (*entities.Sprint, error) {
	m := mappers.EntityToModel(sprint)

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if m.Status == statusActive {
			if err := closeOtherActiveSprints(tx, 0); err != nil {
				return err
			}
		}
		return tx.Create(m).Error
	})
	if err != nil {
		return nil, err
	}

	sprint.ID = m.ID
	sprint.CreatedAt = m.CreatedAt
	sprint.UpdatedAt = m.UpdatedAt
	return sprint, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*entities.Sprint, error) {
	var m models.Sprint
	err := r.db.Conn(ctx).Preload("CreatedBy").Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dom.ErrSprintNotFound
		}
		return nil, err
	}

	out := mappers.ModelToEntity(&m)
	counts, err := r.countTicketsBySprint(ctx, []uint{id})
	if err != nil {
		return nil, err
	}
	if c, ok := counts[id]; ok {
		out.TicketCount = c.TicketCount
		out.DoneCount = c.DoneCount
	}
	return out, nil
}

func (r *Repository) List(ctx context.Context, params dtos.ListSprintsParams) ([]entities.Sprint, int64, error) {
	q := r.db.Conn(ctx).Model(&models.Sprint{}).Preload("CreatedBy")
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var ms []models.Sprint
	offset := (params.Page - 1) * params.PageSize
	if err := q.Order("start_date DESC, id DESC").Limit(params.PageSize).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	out := make([]entities.Sprint, 0, len(ms))
	ids := make([]uint, 0, len(ms))
	for i := range ms {
		out = append(out, *mappers.ModelToEntity(&ms[i]))
		ids = append(ids, ms[i].ID)
	}

	counts, err := r.countTicketsBySprint(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	for i := range out {
		if c, ok := counts[out[i].ID]; ok {
			out[i].TicketCount = c.TicketCount
			out[i].DoneCount = c.DoneCount
		}
	}
	return out, total, nil
}

func (r *Repository) Update(ctx context.Context, id uint, updates map[string]any) (*entities.Sprint, error) {
	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		var m models.Sprint
		if err := tx.Where("id = ?", id).First(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dom.ErrSprintNotFound
			}
			return err
		}
		if status, ok := updates["status"].(string); ok && status == statusActive {
			if err := closeOtherActiveSprints(tx, id); err != nil {
				return err
			}
		}
		if len(updates) == 0 {
			return nil
		}
		return tx.Model(&models.Sprint{}).Where("id = ?", id).Updates(updates).Error
	})
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		var m models.Sprint
		if err := tx.Where("id = ?", id).First(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dom.ErrSprintNotFound
			}
			return err
		}
		if err := detachTicketsFromSprint(tx, id); err != nil {
			return err
		}
		return tx.Delete(&models.Sprint{}, id).Error
	})
}

func closeOtherActiveSprints(tx *gorm.DB, exceptID uint) error {
	return tx.Model(&models.Sprint{}).
		Where("status = ? AND id <> ?", statusActive, exceptID).
		Update("status", statusClosed).Error
}
