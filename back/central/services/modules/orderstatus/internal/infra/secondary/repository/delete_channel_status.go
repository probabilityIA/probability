package repository

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/secondary/repository/models"
	"gorm.io/gorm"
)

func (r *repository) DeleteChannelStatus(ctx context.Context, id uint) error {
	// Obtener el canal status para conocer su integration_type_id y code
	var model models.IntegrationChannelStatus
	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domainerrors.ErrChannelStatusNotFound
		}
		return err
	}

	// Verificar que no haya mapeos que usen este estado del canal
	var count int64
	r.db.Conn(ctx).Model(&models.OrderStatusMapping{}).
		Where("integration_type_id = ? AND original_status = ?", model.IntegrationTypeID, model.Code).
		Count(&count)
	if count > 0 {
		return domainerrors.ErrChannelStatusHasMappings
	}

	return r.db.Conn(ctx).Delete(&model).Error
}
