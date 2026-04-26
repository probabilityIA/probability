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
