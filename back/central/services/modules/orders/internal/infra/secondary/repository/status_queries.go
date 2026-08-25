package repository

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/migration/shared/models"
)

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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if statusID == nil {
		return nil, nil
	}

	return statusID, nil
}

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

func (r *Repository) resolveOrderStatusByCode(ctx context.Context, orders []models.Order) {
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

	var statuses []models.OrderStatus
	if err := r.db.Conn(ctx).
		Where("code IN ?", codes).
		Where("deleted_at IS NULL").
		Find(&statuses).Error; err != nil {
		return
	}

	statusMap := make(map[string]models.OrderStatus, len(statuses))
	for _, s := range statuses {
		statusMap[s.Code] = s
	}

	for i := range orders {
		if orders[i].OrderStatus.ID == 0 && orders[i].Status != "" {
			if status, ok := statusMap[orders[i].Status]; ok {
				orders[i].OrderStatus = status
			}
		}
	}
}

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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &paymentStatus.ID, nil
}

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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &fulfillmentStatus.ID, nil
}

func (r *Repository) IsChannelStatusInboundEnabled(ctx context.Context, integrationID uint) (bool, error) {
	if integrationID == 0 {
		return true, nil
	}

	var result struct {
		Config []byte
	}
	err := r.db.Conn(ctx).
		Table("integrations").
		Select("config").
		Where("id = ?", integrationID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return true, err
	}
	if len(result.Config) == 0 {
		return true, nil
	}

	var config map[string]interface{}
	if err := json.Unmarshal(result.Config, &config); err != nil {
		return true, nil
	}

	valor, existe := config["status_inbound_enabled"]
	if !existe {
		return true, nil
	}
	habilitado, ok := valor.(bool)
	if !ok {
		return true, nil
	}
	return habilitado, nil
}
