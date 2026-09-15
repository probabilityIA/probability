package businessdata

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

const (
	maxMatches = 3
	maxItems   = 10

	withoutGuideSQL = "(orders.tracking_number IS NULL OR orders.tracking_number = '') AND NOT EXISTS (" +
		"SELECT 1 FROM shipments s WHERE s.order_id = orders.id AND s.deleted_at IS NULL " +
		"AND COALESCE(s.tracking_number, '') <> '' AND s.status NOT IN ('cancelled', 'failed'))"

	withNoveltySQL = "((orders.novelty IS NOT NULL AND orders.novelty <> '') OR EXISTS (" +
		"SELECT 1 FROM shipments s WHERE s.order_id = orders.id AND s.deleted_at IS NULL AND s.status = 'on_hold'))"
)

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func selectIDName(db *gorm.DB) *gorm.DB {
	return db.Select("id", "name")
}

func latestShipments(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}

func (r *Reader) businessOrders(ctx context.Context, businessID uint) *gorm.DB {
	return r.db.Conn(ctx).Model(&models.Order{}).
		Where("orders.business_id = ? AND orders.deleted_at IS NULL", businessID)
}

func (r *Reader) FindOrders(ctx context.Context, businessID uint, number string) ([]entities.OrderInfo, error) {
	clean := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(number), "#"))
	if clean == "" {
		return nil, nil
	}

	var rows []models.Order
	err := r.businessOrders(ctx, businessID).
		Where("(LOWER(orders.order_number) = LOWER(?) OR LOWER(orders.order_number) = LOWER(?) OR LOWER(orders.internal_number) = LOWER(?) OR orders.external_id = ?)",
			clean, "#"+clean, clean, clean).
		Preload("OrderStatus", selectIDName).
		Preload("PaymentStatus", selectIDName).
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_id", "product_name", "variant_label", "quantity").Order("id")
		}).
		Preload("Shipments", latestShipments).
		Order("orders.created_at DESC").
		Limit(maxMatches).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrderInfo, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOrderInfo(row))
	}
	return out, nil
}

func (r *Reader) ListOrders(ctx context.Context, businessID uint, query dtos.OrderQuery) ([]entities.OrderSummary, int64, error) {
	build := func() *gorm.DB {
		db := r.businessOrders(ctx, businessID)
		if query.From != nil {
			db = db.Where("orders.created_at >= ?", *query.From)
		}
		if query.To != nil {
			db = db.Where("orders.created_at < ?", *query.To)
		}
		if status := strings.TrimSpace(query.Status); status != "" {
			pattern := "%" + likeEscaper.Replace(status) + "%"
			db = db.Joins("LEFT JOIN order_statuses ON order_statuses.id = orders.status_id").
				Where("(order_statuses.name ILIKE ? OR orders.status ILIKE ?)", pattern, pattern)
		}
		if query.CodOnly {
			db = db.Where("orders.is_cod = true")
		}
		if query.WithoutGuide {
			db = db.Where(withoutGuideSQL)
		}
		if query.WithNovelty {
			db = db.Where(withNoveltySQL)
		}
		return db
	}

	var total int64
	if err := build().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.Order
	err := build().
		Select("orders.*").
		Preload("OrderStatus", selectIDName).
		Preload("Shipments", latestShipments).
		Order("orders.created_at DESC").
		Limit(query.Limit).
		Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	out := make([]entities.OrderSummary, 0, len(rows))
	for _, row := range rows {
		summary := entities.OrderSummary{
			Number:       orderNumber(row),
			CreatedAt:    orderDate(row),
			CustomerName: row.CustomerName,
			Total:        row.TotalAmount,
			Currency:     row.Currency,
			StatusName:   statusName(row),
			IsCod:        row.IsCod,
			IsTest:       row.IsTest,
		}
		if len(row.Shipments) > 0 {
			summary.TrackingNumber = deref(row.Shipments[0].TrackingNumber)
			summary.ShipmentStatus = row.Shipments[0].Status
		}
		if summary.TrackingNumber == "" {
			summary.TrackingNumber = deref(row.TrackingNumber)
		}
		out = append(out, summary)
	}
	return out, total, nil
}

func toOrderInfo(row models.Order) entities.OrderInfo {
	info := entities.OrderInfo{
		Number:        orderNumber(row),
		Platform:      row.Platform,
		CreatedAt:     orderDate(row),
		StatusName:    statusName(row),
		PaymentStatus: row.PaymentStatus.Name,
		IsPaid:        row.IsPaid,
		Total:         row.TotalAmount,
		Currency:      row.Currency,
		IsCod:         row.IsCod,
		CodTotal:      row.CodTotal,
		CustomerName:  row.CustomerName,
		City:          row.ShippingCity,
		IsConfirmed:   row.IsConfirmed != nil && *row.IsConfirmed,
		Novelty:       deref(row.Novelty),
		HasInvoice:    deref(row.InvoiceID) != "",
		IsTest:        row.IsTest,
	}
	for i, item := range row.OrderItems {
		if i >= maxItems {
			break
		}
		info.Items = append(info.Items, entities.OrderItemInfo{Name: item.ProductName, Variant: item.VariantLabel, Quantity: item.Quantity})
	}
	if len(row.Shipments) > 0 {
		shipment := toShipmentInfo(row.Shipments[0], info.Number)
		info.Shipment = &shipment
	}
	return info
}

func orderNumber(row models.Order) string {
	if row.OrderNumber != "" {
		return row.OrderNumber
	}
	return row.InternalNumber
}

func orderDate(row models.Order) time.Time {
	if !row.OccurredAt.IsZero() {
		return row.OccurredAt
	}
	return row.CreatedAt
}

func statusName(row models.Order) string {
	if row.OrderStatus.Name != "" {
		return row.OrderStatus.Name
	}
	return row.Status
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
