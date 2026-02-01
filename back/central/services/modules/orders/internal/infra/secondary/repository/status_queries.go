package repository

import (
	"context"

	"gorm.io/gorm"
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
		Table("order_statuses os").
		Select("os.id").
		Joins("INNER JOIN order_status_mappings osm ON os.id = osm.order_status_id").
		Where("osm.integration_type_id = ?", integrationTypeID).
		Where("osm.original_status = ?", originalStatus).
		Where("osm.deleted_at IS NULL").
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
	var result struct {
		ID uint
	}

	err := r.db.Conn(ctx).
		Table("payment_statuses").
		Select("id").
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Limit(1).
		First(&result).Error

	if err != nil {
		// Si no existe, retornar nil sin error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &result.ID, nil
}

// GetFulfillmentStatusIDByCode obtiene el ID de un estado de fulfillment por su código
func (r *Repository) GetFulfillmentStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	var result struct {
		ID uint
	}

	err := r.db.Conn(ctx).
		Table("fulfillment_statuses").
		Select("id").
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Limit(1).
		First(&result).Error

	if err != nil {
		// Si no existe, retornar nil sin error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &result.ID, nil
}
