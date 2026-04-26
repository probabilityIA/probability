package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) ListCODShipments(ctx context.Context, filter domain.CODFilter) ([]domain.Shipment, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := r.db.Conn(ctx).
		Model(&models.Shipment{}).
		Joins("INNER JOIN orders ON orders.id = shipments.order_id").
		Where("orders.deleted_at IS NULL").
		Where("orders.cod_total IS NOT NULL AND orders.cod_total > 0")

	if filter.BusinessID > 0 {
		query = query.Where("orders.business_id = ?", filter.BusinessID)
	}
	if filter.Status != "" {
		query = query.Where("shipments.status = ?", filter.Status)
	}
	if filter.IsPaid != nil {
		query = query.Where("orders.is_paid = ?", *filter.IsPaid)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var dbShipments []models.Shipment
	err := query.
		Order("shipments.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Preload("Order").
		Preload("Order.PaymentMethod").
		Preload("ShippingAddress").
		Find(&dbShipments).Error
	if err != nil {
		return nil, 0, err
	}

	out := make([]domain.Shipment, len(dbShipments))
	for i := range dbShipments {
		out[i] = *mappers.ToDomainShipment(&dbShipments[i])
		if dbShipments[i].Order != nil && dbShipments[i].Order.PaymentMethodID > 0 {
			out[i].PaymentMethodCode = dbShipments[i].Order.PaymentMethod.Code
		}
	}

	return out, total, nil
}

func (r *Repository) GetOrderPublicTrackingByNumber(ctx context.Context, orderNumber string) (*domain.OrderPublicTracking, error) {
	if orderNumber == "" {
		return nil, fmt.Errorf("order_number requerido")
	}
	var result struct {
		ID                 string
		OrderNumber        string     `gorm:"column:order_number"`
		BusinessID         *uint      `gorm:"column:business_id"`
		BusinessName       string     `gorm:"column:business_name"`
		Status             string
		IsPaid             bool       `gorm:"column:is_paid"`
		TotalAmount        float64    `gorm:"column:total_amount"`
		CodTotal           *float64   `gorm:"column:cod_total"`
		Currency           string
		CustomerName       string     `gorm:"column:customer_name"`
		CustomerPhone      string     `gorm:"column:customer_phone"`
		ShippingStreet     string     `gorm:"column:shipping_street"`
		ShippingCity       string     `gorm:"column:shipping_city"`
		ShippingState      string     `gorm:"column:shipping_state"`
		ShippingPostalCode string     `gorm:"column:shipping_postal_code"`
		CreatedAt          time.Time  `gorm:"column:created_at"`
		OccurredAt         *time.Time `gorm:"column:occurred_at"`
	}
	err := r.db.Conn(ctx).
		Table("orders o").
		Select(`o.id, o.order_number, o.business_id, COALESCE(b.name,'') AS business_name,
			o.status, o.is_paid, o.total_amount, o.cod_total, o.currency,
			o.customer_name, o.customer_phone,
			o.shipping_street, o.shipping_city, o.shipping_state, o.shipping_postal_code,
			o.created_at, o.occurred_at`).
		Joins("LEFT JOIN business b ON b.id = o.business_id").
		Joins(`LEFT JOIN LATERAL (
			SELECT 1 AS has_real
			FROM shipments s
			WHERE s.order_id = o.id
			  AND s.deleted_at IS NULL
			  AND s.tracking_number IS NOT NULL
			  AND s.tracking_number <> ''
			LIMIT 1
		) shp ON true`).
		Where("o.order_number = ? AND o.deleted_at IS NULL", orderNumber).
		Order("shp.has_real DESC NULLS LAST, o.created_at DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if result.ID == "" {
		return nil, nil
	}
	out := &domain.OrderPublicTracking{
		ID:                 result.ID,
		OrderNumber:        result.OrderNumber,
		BusinessName:       result.BusinessName,
		Status:             result.Status,
		IsPaid:             result.IsPaid,
		TotalAmount:        result.TotalAmount,
		CodTotal:           result.CodTotal,
		Currency:           result.Currency,
		CustomerName:       result.CustomerName,
		CustomerPhone:      result.CustomerPhone,
		ShippingStreet:     result.ShippingStreet,
		ShippingCity:       result.ShippingCity,
		ShippingState:      result.ShippingState,
		ShippingPostalCode: result.ShippingPostalCode,
		CreatedAt:          result.CreatedAt,
		OccurredAt:         result.OccurredAt,
	}
	if result.BusinessID != nil {
		out.BusinessID = *result.BusinessID
	}
	return out, nil
}

func (r *Repository) GetOrderCODInfo(ctx context.Context, orderID string) (*domain.OrderCODInfo, error) {
	var result struct {
		ID                string
		BusinessID        *uint
		CodTotal          *float64
		TotalAmount       float64
		Currency          string
		IsPaid            bool
		PaidAt            *time.Time
		PaymentMethodID   uint
		PaymentMethodCode string
	}

	err := r.db.Conn(ctx).
		Table("orders o").
		Select("o.id, o.business_id, o.cod_total, o.total_amount, o.currency, o.is_paid, o.paid_at, o.payment_method_id, pm.code AS payment_method_code").
		Joins("LEFT JOIN payment_methods pm ON pm.id = o.payment_method_id").
		Where("o.id = ? AND o.deleted_at IS NULL", orderID).
		Limit(1).
		Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("orden %s no encontrada", orderID)
		}
		return nil, err
	}
	if result.ID == "" {
		return nil, fmt.Errorf("orden %s no encontrada", orderID)
	}

	info := &domain.OrderCODInfo{
		OrderID:           result.ID,
		CodTotal:          result.CodTotal,
		TotalAmount:       result.TotalAmount,
		Currency:          result.Currency,
		IsPaid:            result.IsPaid,
		PaidAt:            result.PaidAt,
		PaymentMethodID:   result.PaymentMethodID,
		PaymentMethodCode: result.PaymentMethodCode,
	}
	if result.BusinessID != nil {
		info.BusinessID = *result.BusinessID
	}
	return info, nil
}
