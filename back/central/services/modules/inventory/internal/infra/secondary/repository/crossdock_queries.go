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

func (r *Repository) CreateCrossDockLink(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
	m := mappers.CrossDockLinkEntityToModel(l)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return mappers.CrossDockLinkModelToEntity(m), nil
}

func (r *Repository) GetCrossDockLinkByID(ctx context.Context, businessID, id uint) (*entities.CrossDockLink, error) {
	var m models.CrossDockLink
	if err := r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrCrossDockNotFound
		}
		return nil, err
	}
	return mappers.CrossDockLinkModelToEntity(&m), nil
}

func (r *Repository) ListCrossDockLinks(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error) {
	var ml []models.CrossDockLink
	var total int64
	q := r.db.Conn(ctx).Model(&models.CrossDockLink{}).Where("business_id = ?", params.BusinessID)
	if params.Status != "" {
		q = q.Where("status = ?", params.Status)
	}
	if params.OutboundOrderID != "" {
		q = q.Where("outbound_order_id = ?", params.OutboundOrderID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(params.Offset()).Limit(params.PageSize).Order("id DESC").Find(&ml).Error; err != nil {
		return nil, 0, err
	}
	out := make([]entities.CrossDockLink, len(ml))
	for i := range ml {
		out[i] = *mappers.CrossDockLinkModelToEntity(&ml[i])
	}
	return out, total, nil
}

func (r *Repository) UpdateCrossDockLink(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
	updates := map[string]any{
		"status":      l.Status,
		"executed_at": l.ExecutedAt,
	}
	res := r.db.Conn(ctx).Model(&models.CrossDockLink{}).
		Where("id = ? AND business_id = ?", l.ID, l.BusinessID).
		Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domainerrors.ErrCrossDockNotFound
	}
	return r.GetCrossDockLinkByID(ctx, l.BusinessID, l.ID)
}
