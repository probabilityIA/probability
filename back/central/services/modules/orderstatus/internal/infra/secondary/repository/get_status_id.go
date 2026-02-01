package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
	var model models.OrderStatusMapping

	err := r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Where("integration_type_id = ? AND original_status = ? AND is_active = ?", integrationTypeID, originalStatus, true).
		Order("priority DESC"). // Si hay múltiples, tomar el de mayor prioridad
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No se encontró mapeo, retornar nil sin error
			return nil, nil
		}
		return nil, err
	}

	return &model.OrderStatusID, nil
}
