package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) CancelInvoiceAndReleaseOrder(ctx context.Context, invoiceID uint) (bool, string, error) {
	var orderID string
	var released bool

	err := r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		var invoice models.Invoice
		if err := tx.First(&invoice, invoiceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to load invoice %d: %w", invoiceID, err)
		}

		orderID = invoice.OrderID

		if invoice.Status == constants.InvoiceStatusCancelled {
			return nil
		}

		now := time.Now()
		if err := tx.Model(&models.Invoice{}).
			Where("id = ?", invoiceID).
			Updates(map[string]interface{}{
				"status":       constants.InvoiceStatusCancelled,
				"cancelled_at": now,
			}).Error; err != nil {
			return fmt.Errorf("failed to cancel invoice %d: %w", invoiceID, err)
		}

		if orderID != "" {
			if err := tx.Model(&models.Order{}).
				Where("id = ?", orderID).
				Updates(map[string]interface{}{
					"invoice_id":  gorm.Expr("NULL"),
					"invoice_url": gorm.Expr("NULL"),
				}).Error; err != nil {
				return fmt.Errorf("failed to release order %s: %w", orderID, err)
			}
		}

		released = true
		return nil
	})

	if err != nil {
		r.log.Error(ctx).Err(err).Uint("invoice_id", invoiceID).Msg("Failed to cancel invoice and release order")
		return false, orderID, err
	}

	if released {
		r.log.Info(ctx).Uint("invoice_id", invoiceID).Str("order_id", orderID).Msg("Invoice cancelled and order released for re-invoicing")
	}

	return released, orderID, nil
}
