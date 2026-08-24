package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/mobile/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) ports.IRepository {
	return &Repository{db: database}
}

func (r *Repository) GetOrderFull(ctx context.Context, businessID uint, orderID string) (*entities.OrderFull, error) {
	conn := r.db.Conn(ctx)

	var order models.Order
	err := conn.
		Where("id = ? AND deleted_at IS NULL", orderID).
		First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainerrors.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	if order.BusinessID == nil || *order.BusinessID != businessID {
		return nil, domainerrors.ErrOrderNotInScope
	}

	result := &entities.OrderFull{Order: toOrderSummary(order)}

	var items []models.OrderItem
	if err := conn.
		Where("order_id = ? AND deleted_at IS NULL", orderID).
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	result.Items = toOrderLines(items)

	var shipment models.Shipment
	err = conn.
		Where("order_id = ? AND deleted_at IS NULL", orderID).
		Order("id DESC").
		First(&shipment).Error
	if err == nil {
		result.Shipment = toShipmentSummary(shipment)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var invoice models.Invoice
	err = conn.
		Where("order_id = ? AND business_id = ? AND deleted_at IS NULL", orderID, businessID).
		Order("id DESC").
		First(&invoice).Error
	if err == nil {
		result.Invoice = toInvoiceSummary(invoice)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return result, nil
}
