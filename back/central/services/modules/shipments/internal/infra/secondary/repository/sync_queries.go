package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"gorm.io/gorm"
)

func (r *Repository) ListShipmentsForSync(ctx context.Context, filter domain.SyncShipmentsFilter) ([]domain.SyncShipmentRow, error) {
	q := r.db.Conn(ctx).
		Table("shipments s").
		Select(`s.id AS shipment_id,
			s.tracking_number AS tracking_number,
			COALESCE((s.metadata->>'envioclick_id_order')::bigint, 0) AS envioclick_id_order`).
		Joins("JOIN orders o ON o.id::text = s.order_id").
		Where("s.deleted_at IS NULL").
		Where("o.deleted_at IS NULL").
		Where("s.tracking_number IS NOT NULL AND s.tracking_number != ''")

	if filter.BusinessID > 0 {
		q = q.Where("o.business_id = ?", filter.BusinessID)
	}

	if filter.Provider == domain.SyncProviderEnvioclick {
		q = q.Where("LOWER(s.carrier) LIKE ? OR s.metadata ? 'envioclick_id_order'", "%envioclick%")
	}

	if len(filter.Statuses) > 0 {
		q = q.Where("s.status IN ?", filter.Statuses)
	}

	if filter.DateFrom != nil {
		q = q.Where("s.created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		q = q.Where("s.created_at <= ?", *filter.DateTo)
	}

	type row struct {
		ShipmentID        uint   `gorm:"column:shipment_id"`
		TrackingNumber    string `gorm:"column:tracking_number"`
		EnvioclickIDOrder int64  `gorm:"column:envioclick_id_order"`
	}
	var rows []row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list shipments for sync: %w", err)
	}

	out := make([]domain.SyncShipmentRow, 0, len(rows))
	for _, r := range rows {
		item := domain.SyncShipmentRow{
			ShipmentID:     r.ShipmentID,
			TrackingNumber: r.TrackingNumber,
		}
		if r.EnvioclickIDOrder > 0 {
			id := r.EnvioclickIDOrder
			item.EnvioclickIDOrder = &id
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *Repository) GetBusinessActiveIntegration(ctx context.Context, businessID uint, providerCode string) (uint, string, error) {
	var result struct {
		ID      uint   `gorm:"column:id"`
		BaseURL string `gorm:"column:base_url"`
	}
	err := r.db.Conn(ctx).
		Table("integrations i").
		Select("i.id, it.base_url").
		Joins("JOIN integration_types it ON it.id = i.integration_type_id").
		Where("i.business_id = ?", businessID).
		Where("it.code = ?", providerCode).
		Where("i.is_active = TRUE").
		Where("i.deleted_at IS NULL").
		Where("it.deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", fmt.Errorf("no active %s integration for business %d", providerCode, businessID)
		}
		return 0, "", err
	}
	if result.ID == 0 {
		return 0, "", fmt.Errorf("no active %s integration for business %d", providerCode, businessID)
	}
	return result.ID, result.BaseURL, nil
}
