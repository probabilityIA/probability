package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ============================================
// BUSINESS ID RESOLUTION QUERIES
// (Replicated locally — module isolation rule)
// Table consulted: orders, shipments
// ============================================

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
// Joins shipments → orders to get the business_id.
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
// Joins shipments → orders to get the business_id.
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
