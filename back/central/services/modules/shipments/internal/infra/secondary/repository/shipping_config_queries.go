package repository

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/shippingpkg"
)

func (r *Repository) GetBusinessPackageConfig(ctx context.Context, businessID uint, warehouseID *uint) (*domain.PackageConfig, error) {
	var rows []struct {
		WarehouseID     *uint
		PackageStrategy string
		Boxes           []byte
	}

	err := r.db.Conn(ctx).
		Table("shipping_configs").
		Select("warehouse_id, package_strategy, boxes").
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	pick := -1
	for i := range rows {
		if warehouseID != nil && rows[i].WarehouseID != nil && *rows[i].WarehouseID == *warehouseID {
			pick = i
			break
		}
		if rows[i].WarehouseID == nil && pick < 0 {
			pick = i
		}
	}
	if pick < 0 {
		return nil, nil
	}

	cfg := &domain.PackageConfig{Strategy: rows[pick].PackageStrategy}
	if len(rows[pick].Boxes) > 0 {
		var boxes []shippingpkg.Box
		if err := json.Unmarshal(rows[pick].Boxes, &boxes); err == nil {
			cfg.Boxes = boxes
		}
	}
	return cfg, nil
}

func (r *Repository) GetOrderPackageItems(ctx context.Context, orderID string) ([]shippingpkg.PackageItem, uint, error) {
	var order struct {
		WarehouseID *uint
	}
	if err := r.db.Conn(ctx).
		Table("orders").
		Select("warehouse_id").
		Where("id = ?", orderID).
		Limit(1).
		Scan(&order).Error; err != nil {
		return nil, 0, err
	}

	var rows []struct {
		SKU      string
		Quantity int
		Weight   *float64
		Length   *float64
		Width    *float64
		Height   *float64
	}

	err := r.db.Conn(ctx).
		Table("order_items oi").
		Select("oi.product_sku as sku, oi.quantity, p.weight, p.length, p.width, p.height").
		Joins("LEFT JOIN products p ON p.id = oi.product_id AND p.deleted_at IS NULL").
		Where("oi.order_id = ? AND oi.deleted_at IS NULL AND oi.product_sku <> ''", orderID).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]shippingpkg.PackageItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, shippingpkg.PackageItem{
			SKU:      row.SKU,
			Quantity: row.Quantity,
			Weight:   row.Weight,
			Length:   row.Length,
			Width:    row.Width,
			Height:   row.Height,
		})
	}

	var warehouseID uint
	if order.WarehouseID != nil {
		warehouseID = *order.WarehouseID
	}
	return items, warehouseID, nil
}
