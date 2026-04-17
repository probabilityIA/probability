package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ============================================
// BUSINESS ID RESOLUTION QUERIES
// (Replicated locally — module isolation rule)
// Table consulted: orders, shipments
// ============================================
// NOTE: UpdateOrderGuideLink is also here (replicated write — module isolation rule)

// GetOrderBusinessID retrieves the business_id for an order by its UUID.
// Used by super admin handlers to resolve which business owns the order.
func (r *Repository) GetOrderBusinessID(ctx context.Context, orderUUID string) (uint, error) {
	var result struct {
		BusinessID uint `gorm:"column:business_id"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select("business_id").
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("orden %s no encontrada", orderUUID)
		}
		return 0, err
	}
	if result.BusinessID == 0 {
		return 0, fmt.Errorf("orden %s no encontrada o sin business asociado", orderUUID)
	}

	return result.BusinessID, nil
}

// GetShipmentBusinessIDByTracking resolves the business_id for a shipment via its tracking number.
// Joins shipments -> orders to get the business_id.
func (r *Repository) GetShipmentBusinessIDByTracking(ctx context.Context, trackingNumber string) (uint, error) {
	var result struct {
		BusinessID uint `gorm:"column:business_id"`
	}

	err := r.db.Conn(ctx).
		Table("shipments s").
		Select("o.business_id").
		Joins("INNER JOIN orders o ON s.order_id = o.id").
		Where("s.tracking_number = ?", trackingNumber).
		Where("s.deleted_at IS NULL").
		Where("o.deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("envío con tracking %s no encontrado", trackingNumber)
		}
		return 0, err
	}
	if result.BusinessID == 0 {
		return 0, fmt.Errorf("envío con tracking %s no encontrado o sin business asociado", trackingNumber)
	}

	return result.BusinessID, nil
}

// GetShipmentBusinessIDByID resolves the business_id for a shipment via its DB ID.
// Joins shipments -> orders to get the business_id.
func (r *Repository) GetShipmentBusinessIDByID(ctx context.Context, shipmentID uint) (uint, error) {
	var result struct {
		BusinessID uint `gorm:"column:business_id"`
	}

	err := r.db.Conn(ctx).
		Table("shipments s").
		Select("o.business_id").
		Joins("INNER JOIN orders o ON s.order_id = o.id").
		Where("s.id = ?", shipmentID).
		Where("s.deleted_at IS NULL").
		Where("o.deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("envío con ID %d no encontrado", shipmentID)
		}
		return 0, err
	}
	if result.BusinessID == 0 {
		return 0, fmt.Errorf("envío con ID %d no encontrado o sin business asociado", shipmentID)
	}

	return result.BusinessID, nil
}

// GetOrderIntegrationID retrieves the integration_id for an order by its UUID.
// Replicated query — module isolation rule.
func (r *Repository) GetOrderIntegrationID(ctx context.Context, orderUUID string) (uint, error) {
	var result struct {
		IntegrationID uint `gorm:"column:integration_id"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select("integration_id").
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("orden %s no encontrada", orderUUID)
		}
		return 0, err
	}

	return result.IntegrationID, nil
}

// UpdateOrderGuideLink updates guide_link, tracking_number, and carrier on the orders table
// after a guide is generated. Replicated write — orders table is owned by the orders
// module but we update it directly to avoid inter-module repository sharing.
func (r *Repository) UpdateOrderGuideLink(ctx context.Context, orderID string, guideLink string, trackingNumber string, carrier string) error {
	updates := map[string]interface{}{}
	if guideLink != "" {
		updates["guide_link"] = guideLink
	}
	if trackingNumber != "" {
		updates["tracking_number"] = trackingNumber
	}
	if carrier != "" {
		updates["carrier"] = carrier
	}
	if len(updates) == 0 {
		return nil
	}

	return r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Where("deleted_at IS NULL").
		Updates(updates).Error
}

func (r *Repository) UpdateOrderStatusByOrderID(ctx context.Context, orderID string, status string) error {
	if orderID == "" || status == "" {
		return nil
	}

	return r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Where("deleted_at IS NULL").
		Update("status", status).Error
}

func (r *Repository) ClearOrderGuideData(ctx context.Context, orderID string) error {
	if orderID == "" {
		return nil
	}
	return r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ? AND deleted_at IS NULL", orderID).
		Updates(map[string]any{
			"tracking_number": nil,
			"tracking_link":   nil,
			"guide_link":      nil,
			"guide_id":        nil,
			"carrier":         nil,
			"status":          "pending",
		}).Error
}

func (r *Repository) EnsureAllBusinessesActive(ctx context.Context) error {
	return r.db.Conn(ctx).
		Table("businesses").
		Where("deleted_at IS NULL").
		Updates(map[string]interface{}{
			"status":     "paid",
			"expiration": "2030-01-01 00:00:00",
		}).Error
}
