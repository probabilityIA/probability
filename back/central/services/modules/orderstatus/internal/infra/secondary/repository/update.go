package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) Update(ctx context.Context, mapping *entities.OrderStatusMapping) error {
	model := models.FromDomain(*mapping)
	return r.db.Conn(ctx).Save(model).Error
}
