package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type ProductRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func New(database db.IDatabase, logger log.ILogger) domain.IProductRepository {
	return &ProductRepository{
		db:  database,
		log: logger.WithModule("vtex.product_repository"),
	}
}

func (r *ProductRepository) ListProductsByBusiness(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
	var rows []struct {
		ID             string
		SKU            string
		Name           string
		Description    string
		Price          float64
		StockQuantity  int
		TrackInventory bool
		ImageURL       string
		Weight         *float64
		WeightUnit     string
		Length         *float64
		Width          *float64
		Height         *float64
		DimensionUnit  string
	}

	err := r.db.Conn(ctx).
		Table("products").
		Select("id, sku, name, description, price, stock_quantity, track_inventory, image_url, weight, weight_unit, length, width, height, dimension_unit").
		Where("business_id = ? AND deleted_at IS NULL AND is_active = ?", businessID, true).
		Order("created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	products := make([]domain.ProductForSync, 0, len(rows))
	for _, row := range rows {
		products = append(products, domain.ProductForSync{
			ID:             row.ID,
			SKU:            row.SKU,
			Name:           row.Name,
			Description:    row.Description,
			Price:          row.Price,
			StockQuantity:  row.StockQuantity,
			TrackInventory: row.TrackInventory,
			ImageURL:       row.ImageURL,
			Weight:         row.Weight,
			WeightUnit:     row.WeightUnit,
			Length:         row.Length,
			Width:          row.Width,
			Height:         row.Height,
			DimensionUnit:  row.DimensionUnit,
		})
	}
	return products, nil
}

func (r *ProductRepository) ListMappedItems(ctx context.Context, integrationID uint) ([]domain.MappedItem, error) {
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

func (r *ProductRepository) GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error) {
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

func (r *ProductRepository) GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error) {
	var result struct {
		ExternalProductID string
	}
	err := r.db.Conn(ctx).
		Table("product_business_integrations").
		Select("external_product_id").
		Where("product_id = ? AND integration_id = ? AND deleted_at IS NULL", productID, integrationID).
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return "", false, err
	}
	if result.ExternalProductID == "" {
		return "", false, nil
	}
	return result.ExternalProductID, true, nil
}

func (r *ProductRepository) UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, externalProductID string) error {
	var existing models.ProductBusinessIntegration
	err := r.db.Conn(ctx).
		Where("product_id = ? AND integration_id = ?", productID, integrationID).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		record := models.ProductBusinessIntegration{
			ProductID:         productID,
			BusinessID:        businessID,
			IntegrationID:     integrationID,
			ExternalProductID: externalProductID,
		}
		return r.db.Conn(ctx).Create(&record).Error
	}
	if err != nil {
		return err
	}

	existing.ExternalProductID = externalProductID
	return r.db.Conn(ctx).Save(&existing).Error
}
