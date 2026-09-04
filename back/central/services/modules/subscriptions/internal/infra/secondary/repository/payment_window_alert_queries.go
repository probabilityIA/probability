package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

const platformIntegrationTypeID = 6

func (r *Repository) GetWhatsAppContact(ctx context.Context, businessID uint) (*entities.WhatsAppContact, error) {
	var integrationCount int64
	if err := r.db.Conn(ctx).
		Table("integrations").
		Where("business_id = ? AND integration_type_id = ? AND is_active = true AND deleted_at IS NULL", businessID, platformIntegrationTypeID).
		Count(&integrationCount).Error; err != nil {
		return nil, err
	}
	if integrationCount == 0 {
		return nil, nil
	}

	var row struct {
		BusinessName string
		Phone        string
	}
	err := r.db.Conn(ctx).
		Table("business b").
		Select(`b.name AS business_name, COALESCE(wh.phone, '') AS phone`).
		Joins(`LEFT JOIN warehouses wh ON wh.business_id = b.id AND wh.is_default = true AND wh.deleted_at IS NULL`).
		Where("b.id = ? AND b.deleted_at IS NULL", businessID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &entities.WhatsAppContact{
		Phone:        row.Phone,
		BusinessName: row.BusinessName,
	}, nil
}

func (r *Repository) HasAuditLogSince(ctx context.Context, businessID uint, action string, since time.Time) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Table("subscription_audit_logs").
		Where("business_id = ? AND action = ? AND created_at >= ?", businessID, action, since).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
