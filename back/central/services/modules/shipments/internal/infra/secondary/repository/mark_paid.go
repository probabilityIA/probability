package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) MarkOrderPaidCOD(ctx context.Context, orderID string, amount float64, paymentMethodID uint, notes string) error {
	now := time.Now().UTC()
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		payment := &models.Payment{
			OrderID:         orderID,
			PaymentMethodID: paymentMethodID,
			Amount:          amount,
			Status:          "completed",
			PaidAt:          &now,
			ProcessedAt:     &now,
		}
		if notes != "" {
			ref := notes
			payment.PaymentReference = &ref
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		return tx.Model(&models.Order{}).
			Where("id = ? AND deleted_at IS NULL", orderID).
			Updates(map[string]any{
				"is_paid": true,
				"paid_at": now,
			}).Error
	})
}
