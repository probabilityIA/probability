package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
)

func (r *repository) ListOrderStatuses(ctx context.Context, isActive *bool) ([]entities.OrderStatusInfo, error) {
	var modelsList []models.OrderStatus
	query := r.db.Conn(ctx).Model(&models.OrderStatus{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("code ASC").Find(&modelsList).Error
	if err != nil {
		return nil, err
	}

	// Convertir a dominio
	result := make([]entities.OrderStatusInfo, len(modelsList))
	for i, status := range modelsList {
		result[i] = status.ToDomain()
	}

	return result, nil
}
