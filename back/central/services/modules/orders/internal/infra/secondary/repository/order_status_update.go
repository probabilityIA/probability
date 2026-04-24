package repository

import (
	"context"

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
