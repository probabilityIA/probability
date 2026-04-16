package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// DeleteAllOrders hard-deletes all orders and related records for a business.
// Deletion order: invoice children -> invoices -> order children -> orders
func (r *Repository) DeleteAllOrders(ctx context.Context, businessID uint) (int64, error) {
	db := r.db.Conn(ctx)

	// Sub-query for invoice IDs related to this business's orders
	invoiceSubQuery := db.Model(&models.Invoice{}).
		Select("id").
		Where("order_id IN (?)",
			db.Model(&models.Order{}).Select("id").Where("business_id = ?", businessID),
		)

	// 1. InvoiceSyncLog
	if err := db.Unscoped().
		Where("invoice_id IN (?)", invoiceSubQuery).
		Delete(&models.InvoiceSyncLog{}).Error; err != nil {
		return 0, fmt.Errorf("delete invoice_sync_logs: %w", err)
	}

	// 2. CreditNote
	if err := db.Unscoped().
		Where("invoice_id IN (?)", invoiceSubQuery).
		Delete(&models.CreditNote{}).Error; err != nil {
		return 0, fmt.Errorf("delete credit_notes: %w", err)
	}

	// 3. InvoiceItem
	if err := db.Unscoped().
		Where("invoice_id IN (?)", invoiceSubQuery).
		Delete(&models.InvoiceItem{}).Error; err != nil {
		return 0, fmt.Errorf("delete invoice_items: %w", err)
	}

	// 4. Invoice
	if err := db.Unscoped().
		Where("order_id IN (?)",
			db.Model(&models.Order{}).Select("id").Where("business_id = ?", businessID),
		).
		Delete(&models.Invoice{}).Error; err != nil {
		return 0, fmt.Errorf("delete invoices: %w", err)
	}

	// Sub-query for order IDs of this business
	orderSubQuery := db.Model(&models.Order{}).
		Select("id").
		Where("business_id = ?", businessID)

	// 5. Shipment
	if err := db.Unscoped().
		Where("order_id IN (?)", orderSubQuery).
		Delete(&models.Shipment{}).Error; err != nil {
		return 0, fmt.Errorf("delete shipments: %w", err)
	}

	// 6. Payment
	if err := db.Unscoped().
		Where("order_id IN (?)", orderSubQuery).
		Delete(&models.Payment{}).Error; err != nil {
		return 0, fmt.Errorf("delete payments: %w", err)
	}

	// 7. OrderChannelMetadata
	if err := db.Unscoped().
		Where("order_id IN (?)", orderSubQuery).
		Delete(&models.OrderChannelMetadata{}).Error; err != nil {
		return 0, fmt.Errorf("delete order_channel_metadata: %w", err)
	}

	// 8. OrderItem
	if err := db.Unscoped().
		Where("order_id IN (?)", orderSubQuery).
		Delete(&models.OrderItem{}).Error; err != nil {
		return 0, fmt.Errorf("delete order_items: %w", err)
	}

	// 9. Order (returns count)
	result := db.Unscoped().
		Where("business_id = ?", businessID).
		Delete(&models.Order{})
	if result.Error != nil {
		return 0, fmt.Errorf("delete orders: %w", result.Error)
	}

	return result.RowsAffected, nil
}
