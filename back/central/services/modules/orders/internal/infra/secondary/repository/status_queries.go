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

// GetOrderStatusIDByCode obtiene el ID de un estado de orden por su código directo
// Tabla consultada: order_statuses (búsqueda directa por code, sin pasar por order_status_mappings)
func (r *Repository) GetOrderStatusIDByCode(ctx context.Context, code string) (*uint, error) {
	var orderStatus models.OrderStatus

	err := r.db.Conn(ctx).
		Model(&models.OrderStatus{}).
		Select("id").
		Where("code = ?", code).
		Where("deleted_at IS NULL").
		Limit(1).
		First(&orderStatus).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &orderStatus.ID, nil
}

// resolveOrderStatusByCodeSingle resuelve el OrderStatus para una orden individual
// que tiene status (code) pero no status_id.
func (r *Repository) resolveOrderStatusByCodeSingle(ctx context.Context, order *models.Order) {
	if order.OrderStatus.ID > 0 || order.Status == "" {
		return
	}

	var status models.OrderStatus
	if err := r.db.Conn(ctx).
		Where("code = ?", order.Status).
		Where("deleted_at IS NULL").
		First(&status).Error; err != nil {
		return
	}
	order.OrderStatus = status
}

// resolveOrderStatusByCode resuelve el OrderStatus para órdenes que tienen status (code) pero no status_id.
// Busca en batch los códigos únicos y asigna el OrderStatus correspondiente.
// Esto permite que órdenes existentes sin status_id muestren el nombre en español.
func (r *Repository) resolveOrderStatusByCode(ctx context.Context, orders []models.Order) {
	// Recolectar códigos de estado únicos de órdenes sin OrderStatus cargado
	codesSet := make(map[string]struct{})
	for _, o := range orders {
		if o.OrderStatus.ID == 0 && o.Status != "" {
			codesSet[o.Status] = struct{}{}
		}
	}

	if len(codesSet) == 0 {
		return
	}

	codes := make([]string, 0, len(codesSet))
	for code := range codesSet {
		codes = append(codes, code)
	}

	// Buscar todos los OrderStatus por código en una sola query
	var statuses []models.OrderStatus
	if err := r.db.Conn(ctx).
		Where("code IN ?", codes).
		Where("deleted_at IS NULL").
		Find(&statuses).Error; err != nil {
		return // No romper el flujo si falla
	}

	// Crear mapa code -> OrderStatus
	statusMap := make(map[string]models.OrderStatus, len(statuses))
	for _, s := range statuses {
		statusMap[s.Code] = s
	}

	// Asignar OrderStatus a órdenes que no lo tienen
	for i := range orders {
		if orders[i].OrderStatus.ID == 0 && orders[i].Status != "" {
			if status, ok := statusMap[orders[i].Status]; ok {
				orders[i].OrderStatus = status
			}
		}
	}
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
