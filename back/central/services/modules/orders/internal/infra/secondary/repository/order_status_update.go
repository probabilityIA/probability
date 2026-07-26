package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID string, status string, statusID *uint) error {
	return r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]any{
			"status":    status,
			"status_id": statusID,
		}).Error
}

func (r *Repository) UpdateOrderStatusWithSource(ctx context.Context, orderID string, status string, statusID *uint, source string, changedBy string) error {
	return r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]any{
			"status":            status,
			"status_id":         statusID,
			"status_source":     source,
			"status_changed_by": changedBy,
			"status_changed_at": time.Now(),
		}).Error
}
