package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// ============================================
// MÉTODOS DE CONSULTA A TABLAS DE ESTADOS
// (Replicados localmente - no compartir repos entre módulos)
// ============================================

// GetOrderStatusIDByIntegrationTypeAndOriginalStatus obtiene el ID de un estado de orden
// basado en el tipo de integración y el estado original de la plataforma
func (r *Repository) GetOrderStatusIDByIntegrationTypeAndOriginalStatus(ctx context.Context, integrationTypeID uint, originalStatus string) (*uint, error) {
	var statusID *uint

	err := r.db.Conn(ctx).
		Model(&models.OrderStatus{}).
		Select("order_statuses.id").
		Joins("INNER JOIN order_status_mappings ON order_statuses.id = order_status_mappings.order_status_id").
		Where("order_status_mappings.integration_type_id = ?", integrationTypeID).
		Where("order_status_mappings.original_status = ?", originalStatus).
		Where("order_status_mappings.deleted_at IS NULL").
		Limit(1).
		Scan(&statusID).Error

	if err != nil {
		// Si no hay mapeo, retornar nil sin error (permitir fallback)
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// Si no se encontró, retornar nil
	if statusID == nil {
		return nil, nil
	}

	return statusID, nil
}

// GetPaymentStatusIDByCode obtiene el ID de un estado de pago por su código
func (r *Repository) GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	var paymentStatus models.PaymentStatus

	err := r.db.Conn(ctx).
		Model(&models.PaymentStatus{}).
		Select("id").
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Limit(1).
		First(&paymentStatus).Error

	if err != nil {
		// Si no existe, retornar nil sin error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &paymentStatus.ID, nil
}

// GetFulfillmentStatusIDByCode obtiene el ID de un estado de fulfillment por su código
func (r *Repository) GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	var fulfillmentStatus models.FulfillmentStatus

	err := r.db.Conn(ctx).
		Model(&models.FulfillmentStatus{}).
		Select("id").
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Limit(1).
		First(&fulfillmentStatus).Error

	if err != nil {
		// Si no existe, retornar nil sin error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &fulfillmentStatus.ID, nil
}
