package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
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
		log: logger.WithModule("meli.product_repository"),
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
	}

	err := r.db.Conn(ctx).
		Table("products").
		Select("id, sku, name, description, price, stock_quantity, track_inventory, image_url").
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
		})
	}
	return products, nil
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
