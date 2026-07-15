package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type InventoryRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func NewInventory(database db.IDatabase, logger log.ILogger) domain.IInventoryRepository {
	return &InventoryRepository{
		db:  database,
		log: logger.WithModule("shopify.inventory_repository"),
	}
}

func (r *InventoryRepository) ListMappedItems(ctx context.Context, integrationID uint) ([]domain.MappedItem, error) {
	var rows []struct {
		ProductID         string
		SKU               string
		ExternalProductID string
	}
	err := r.db.Conn(ctx).
		Table("product_business_integrations AS pbi").
		Select("pbi.product_id, p.sku, pbi.external_product_id").
		Joins("JOIN products p ON p.id = pbi.product_id").
		Where("pbi.integration_id = ? AND pbi.deleted_at IS NULL AND pbi.external_product_id <> '' AND p.deleted_at IS NULL", integrationID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	items := make([]domain.MappedItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.MappedItem{
			ProductID:      row.ProductID,
			SKU:            row.SKU,
			ExternalItemID: row.ExternalProductID,
		})
	}
	return items, nil
}

func (r *InventoryRepository) GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error) {
	result := make(map[string]int)
	if len(productIDs) == 0 {
		return result, nil
	}
	var rows []struct {
		ProductID string
		Qty       int
	}
	query := r.db.Conn(ctx).
		Table("inventory_levels").
		Select("product_id, COALESCE(SUM(available_qty), 0) AS qty").
		Where("product_id IN ? AND deleted_at IS NULL", productIDs)
	if len(warehouseIDs) > 0 {
		query = query.Where("warehouse_id IN ?", warehouseIDs)
	}
	err := query.Group("product_id").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.ProductID] = row.Qty
	}
	return result, nil
}

func (r *InventoryRepository) GetInventoryByWarehouses(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]map[uint]int, error) {
	result := make(map[string]map[uint]int)
	if len(productIDs) == 0 {
		return result, nil
	}
	var rows []struct {
		ProductID   string
		WarehouseID uint
		Qty         int
	}
	query := r.db.Conn(ctx).
		Table("inventory_levels").
		Select("product_id, warehouse_id, COALESCE(SUM(available_qty), 0) AS qty").
		Where("product_id IN ? AND deleted_at IS NULL", productIDs)
	if len(warehouseIDs) > 0 {
		query = query.Where("warehouse_id IN ?", warehouseIDs)
	}
	err := query.Group("product_id, warehouse_id").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if result[row.ProductID] == nil {
			result[row.ProductID] = make(map[uint]int)
		}
		result[row.ProductID][row.WarehouseID] = row.Qty
	}
	return result, nil
}
