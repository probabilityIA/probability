package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (r *repository) ToggleActive(ctx context.Context, id uint) (*entities.OrderStatusMapping, error) {
	mapping, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mapping.IsActive = !mapping.IsActive
	if err := r.Update(ctx, mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}
