package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) seedCustomerHistory(ctx context.Context) error {
	db := r.db.Conn(ctx)

	var count int64
	if err := db.Table("customer_summaries").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check customer_summaries: %w", err)
	}
	if count > 0 {
		return nil
	}

	var hasOrders int64
	if err := db.Table("orders").Where("deleted_at IS NULL AND customer_id IS NOT NULL").Count(&hasOrders).Error; err != nil {
		return fmt.Errorf("failed to check orders: %w", err)
	}
	if hasOrders == 0 {
		return nil
	}

	if err := r.seedCustomerSummaries(ctx, db); err != nil {
		return fmt.Errorf("failed to seed customer summaries: %w", err)
	}

	if err := r.seedCustomerAddresses(ctx, db); err != nil {
		return fmt.Errorf("failed to seed customer addresses: %w", err)
	}

	if err := r.seedCustomerOrderItems(ctx, db); err != nil {
		return fmt.Errorf("failed to seed customer order items: %w", err)
	}

	if err := r.seedCustomerProductHistories(ctx, db); err != nil {
		return fmt.Errorf("failed to seed customer product histories: %w", err)
	}

	return nil
}

type summaryRow struct {
	CustomerID       uint
	BusinessID       uint
	TotalOrders      int
	DeliveredOrders  int
	CancelledOrders  int
	InProgressOrders int
	TotalSpent       float64
	AvgTicket        float64
	TotalPaidOrders  int
	AvgDeliveryScore float64
	FirstOrderAt     *time.Time
	LastOrderAt      *time.Time
	Platform         string
}

func (r *Repository) seedCustomerSummaries(_ context.Context, db *gorm.DB) error {
	var rows []summaryRow
	err := db.Raw(`
		SELECT
			o.customer_id,
			o.business_id,
			COUNT(*) AS total_orders,
			COUNT(*) FILTER (WHERE os.category = 'completed') AS delivered_orders,
			COUNT(*) FILTER (WHERE os.category = 'cancelled') AS cancelled_orders,
			COUNT(*) FILTER (WHERE os.category = 'active') AS in_progress_orders,
			COALESCE(SUM(o.total_amount), 0) AS total_spent,
			CASE WHEN COUNT(*) > 0 THEN COALESCE(SUM(o.total_amount), 0) / COUNT(*) ELSE 0 END AS avg_ticket,
			COUNT(*) FILTER (WHERE o.is_paid = true) AS total_paid_orders,
			COALESCE(AVG(o.delivery_probability), 0) AS avg_delivery_score,
			MIN(o.created_at) AS first_order_at,
			MAX(o.created_at) AS last_order_at,
			MODE() WITHIN GROUP (ORDER BY o.platform) AS platform
		FROM orders o
		LEFT JOIN order_statuses os ON o.status_id = os.id
		WHERE o.deleted_at IS NULL
			AND o.customer_id IS NOT NULL
			AND o.business_id IS NOT NULL
		GROUP BY o.customer_id, o.business_id
	`).Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("failed to query order summaries: %w", err)
	}

	const batchSize = 500
	now := time.Now()
	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[i:end]

		values := make([]map[string]interface{}, len(batch))
		for j, row := range batch {
			values[j] = map[string]interface{}{
				"customer_id":       row.CustomerID,
				"business_id":      row.BusinessID,
				"total_orders":     row.TotalOrders,
				"delivered_orders":  row.DeliveredOrders,
				"cancelled_orders":  row.CancelledOrders,
				"in_progress_orders": row.InProgressOrders,
				"total_spent":       row.TotalSpent,
				"avg_ticket":        row.AvgTicket,
				"total_paid_orders": row.TotalPaidOrders,
				"avg_delivery_score": row.AvgDeliveryScore,
				"first_order_at":    row.FirstOrderAt,
				"last_order_at":     row.LastOrderAt,
				"preferred_platform": row.Platform,
				"last_updated_at":   now,
				"created_at":        now,
				"updated_at":        now,
			}
		}

		if err := db.Table("customer_summaries").Create(values).Error; err != nil {
			return fmt.Errorf("failed to insert customer_summaries batch %d: %w", i/batchSize, err)
		}
	}

	return nil
}

type addressRow struct {
	CustomerID uint
	BusinessID uint
	Street     string
	City       string
	State      string
	Country    string
	PostalCode string
	TimesUsed  int
	LastUsedAt time.Time
}

func (r *Repository) seedCustomerAddresses(_ context.Context, db *gorm.DB) error {
	var rows []addressRow
	err := db.Raw(`
		SELECT
			o.customer_id,
			o.business_id,
			o.shipping_street AS street,
			o.shipping_city AS city,
			o.shipping_state AS state,
			o.shipping_country AS country,
			o.shipping_postal_code AS postal_code,
			COUNT(*) AS times_used,
			MAX(o.created_at) AS last_used_at
		FROM orders o
		WHERE o.deleted_at IS NULL
			AND o.customer_id IS NOT NULL
			AND o.business_id IS NOT NULL
			AND (o.shipping_street != '' OR o.shipping_city != '')
		GROUP BY o.customer_id, o.business_id, o.shipping_street, o.shipping_city, o.shipping_state, o.shipping_country, o.shipping_postal_code
	`).Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("failed to query addresses: %w", err)
	}

	const batchSize = 500
	now := time.Now()
	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[i:end]

		values := make([]map[string]interface{}, len(batch))
		for j, row := range batch {
			values[j] = map[string]interface{}{
				"customer_id": row.CustomerID,
				"business_id": row.BusinessID,
				"street":      row.Street,
				"city":        row.City,
				"state":       row.State,
				"country":     row.Country,
				"postal_code": row.PostalCode,
				"times_used":  row.TimesUsed,
				"last_used_at": row.LastUsedAt,
				"created_at":  now,
				"updated_at":  now,
			}
		}

		if err := db.Table("customer_addresses").Create(values).Error; err != nil {
			return fmt.Errorf("failed to insert customer_addresses batch %d: %w", i/batchSize, err)
		}
	}

	return nil
}

type orderItemRow struct {
	CustomerID   uint
	BusinessID   uint
	OrderID      string
	OrderNumber  string
	ProductID    *string
	ProductName  string
	ProductSKU   string
	ProductImage *string
	Quantity     int
	UnitPrice    float64
	TotalPrice   float64
	OrderStatus  string
	OrderedAt    time.Time
}

func (r *Repository) seedCustomerOrderItems(_ context.Context, db *gorm.DB) error {
	var rows []orderItemRow
	err := db.Raw(`
		SELECT
			o.customer_id,
			o.business_id,
			o.id AS order_id,
			o.order_number,
			oi.product_id,
			COALESCE(p.name, '') AS product_name,
			COALESCE(p.sku, '') AS product_sku,
			p.image_url AS product_image,
			oi.quantity,
			oi.unit_price,
			oi.total_price,
			COALESCE(os.code, '') AS order_status,
			o.created_at AS ordered_at
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id AND o.deleted_at IS NULL
		LEFT JOIN products p ON p.id = oi.product_id AND p.deleted_at IS NULL
		LEFT JOIN order_statuses os ON os.id = o.status_id
		WHERE oi.deleted_at IS NULL
			AND o.customer_id IS NOT NULL
			AND o.business_id IS NOT NULL
	`).Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("failed to query order items: %w", err)
	}

	const batchSize = 500
	now := time.Now()
	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[i:end]

		values := make([]map[string]interface{}, len(batch))
		for j, row := range batch {
			values[j] = map[string]interface{}{
				"customer_id":   row.CustomerID,
				"business_id":  row.BusinessID,
				"order_id":     row.OrderID,
				"order_number": row.OrderNumber,
				"product_id":   row.ProductID,
				"product_name": row.ProductName,
				"product_sku":  row.ProductSKU,
				"product_image": row.ProductImage,
				"quantity":     row.Quantity,
				"unit_price":   row.UnitPrice,
				"total_price":  row.TotalPrice,
				"order_status": row.OrderStatus,
				"ordered_at":   row.OrderedAt,
				"created_at":   now,
				"updated_at":   now,
			}
		}

		if err := db.Table("customer_order_items").Create(values).Error; err != nil {
			return fmt.Errorf("failed to insert customer_order_items batch %d: %w", i/batchSize, err)
		}
	}

	return nil
}

type productHistoryRow struct {
	CustomerID    uint
	BusinessID    uint
	ProductID     string
	ProductName   string
	ProductSKU    string
	ProductImage  *string
	TimesOrdered  int
	TotalQuantity int
	TotalSpent    float64
	FirstOrderedAt time.Time
	LastOrderedAt  time.Time
}

func (r *Repository) seedCustomerProductHistories(_ context.Context, db *gorm.DB) error {
	var rows []productHistoryRow
	err := db.Raw(`
		SELECT
			o.customer_id,
			o.business_id,
			oi.product_id,
			COALESCE(p.name, '') AS product_name,
			COALESCE(p.sku, '') AS product_sku,
			p.image_url AS product_image,
			COUNT(DISTINCT o.id) AS times_ordered,
			COALESCE(SUM(oi.quantity), 0) AS total_quantity,
			COALESCE(SUM(oi.total_price), 0) AS total_spent,
			MIN(o.created_at) AS first_ordered_at,
			MAX(o.created_at) AS last_ordered_at
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id AND o.deleted_at IS NULL
		LEFT JOIN products p ON p.id = oi.product_id AND p.deleted_at IS NULL
		WHERE oi.deleted_at IS NULL
			AND o.customer_id IS NOT NULL
			AND o.business_id IS NOT NULL
			AND oi.product_id IS NOT NULL
		GROUP BY o.customer_id, o.business_id, oi.product_id, p.name, p.sku, p.image_url
	`).Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("failed to query product histories: %w", err)
	}

	const batchSize = 500
	now := time.Now()
	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[i:end]

		values := make([]map[string]interface{}, len(batch))
		for j, row := range batch {
			values[j] = map[string]interface{}{
				"customer_id":     row.CustomerID,
				"business_id":    row.BusinessID,
				"product_id":     row.ProductID,
				"product_name":   row.ProductName,
				"product_sku":    row.ProductSKU,
				"product_image":  row.ProductImage,
				"times_ordered":  row.TimesOrdered,
				"total_quantity": row.TotalQuantity,
				"total_spent":   row.TotalSpent,
				"first_ordered_at": row.FirstOrderedAt,
				"last_ordered_at":  row.LastOrderedAt,
				"created_at":      now,
				"updated_at":      now,
			}
		}

		if err := db.Table("customer_product_histories").Create(values).Error; err != nil {
			return fmt.Errorf("failed to insert customer_product_histories batch %d: %w", i/batchSize, err)
		}
	}

	return nil
}
