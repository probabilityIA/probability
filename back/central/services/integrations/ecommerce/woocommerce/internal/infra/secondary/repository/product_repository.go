package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/productmatch"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type ProductRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func New(database db.IDatabase, logger log.ILogger) domain.IProductRepository {
	return &ProductRepository{
		db:  database,
		log: logger.WithModule("woocommerce.product_repository"),
	}
}

func (r *ProductRepository) ListProductsByBusiness(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
	var rows []struct {
		ID                string
		SKU               string
		Barcode           string
		ExternalID        string
		Name              string
		Description       string
		Price             float64
		StockQuantity     int
		TrackInventory    bool
		ImageURL          string
		FamilyID          *uint
		FamilyName        string
		FamilyDescription string
		FamilyImageURL    string
		SubfamilyID       *uint
		SubfamilyName     string
		VariantLabel      string
		VariantAttributes datatypes.JSON
	}

	err := r.db.Conn(ctx).
		Table("products p").
		Select(`p.id, p.sku, COALESCE(p.barcode, '') AS barcode, p.external_id, p.name, p.description,
			p.price, p.stock_quantity, p.track_inventory, p.image_url,
			COALESCE(mf.id, f.id) AS family_id,
			COALESCE(mf.name, f.name, '') AS family_name,
			COALESCE(mf.description, f.description, '') AS family_description,
			COALESCE(NULLIF(mf.image_url, ''), NULLIF(f.image_url, ''), '') AS family_image_url,
			CASE WHEN mf.id IS NOT NULL THEN f.id END AS subfamily_id,
			CASE WHEN mf.id IS NOT NULL THEN COALESCE(f.name, '') END AS subfamily_name,
			COALESCE(p.variant_label, '') AS variant_label, p.variant_attributes`).
		Joins("LEFT JOIN product_families f ON f.id = p.family_id AND f.deleted_at IS NULL").
		Joins("LEFT JOIN product_families mf ON mf.id = f.parent_family_id AND mf.deleted_at IS NULL").
		Where("p.business_id = ? AND p.deleted_at IS NULL AND p.is_active = ?", businessID, true).
		Order("p.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	products := make([]domain.ProductForSync, 0, len(rows))
	for _, row := range rows {
		familyID := ""
		if row.FamilyID != nil {
			familyID = strconv.FormatUint(uint64(*row.FamilyID), 10)
		}
		subfamilyID := ""
		if row.SubfamilyID != nil {
			subfamilyID = strconv.FormatUint(uint64(*row.SubfamilyID), 10)
		}
		var attrs map[string]string
		if len(row.VariantAttributes) > 0 {
			_ = json.Unmarshal(row.VariantAttributes, &attrs)
		}
		products = append(products, domain.ProductForSync{
			ID:                row.ID,
			SKU:               row.SKU,
			Barcode:           row.Barcode,
			ExternalID:        row.ExternalID,
			Name:              row.Name,
			Description:       row.Description,
			Price:             row.Price,
			StockQuantity:     row.StockQuantity,
			TrackInventory:    row.TrackInventory,
			ImageURL:          row.ImageURL,
			FamilyID:          familyID,
			FamilyName:        row.FamilyName,
			FamilyDescription: row.FamilyDescription,
			FamilyImageURL:    row.FamilyImageURL,
			SubfamilyID:       subfamilyID,
			SubfamilyName:     row.SubfamilyName,
			VariantLabel:      row.VariantLabel,
			VariantAttributes: attrs,
		})
	}
	return products, nil
}

func (r *ProductRepository) ListMappedItems(ctx context.Context, integrationID uint) ([]domain.MappedItem, error) {
	var rows []struct {
		ProductID         string
		SKU               string
		Name              string
		ImageURL          string
		Barcode           string
		ExternalProductID string
		ExternalVariantID string
		ExternalSKU       string
		ExternalBarcode   string
	}
	err := r.db.Conn(ctx).
		Table("product_business_integrations AS pbi").
		Select(`pbi.product_id, p.sku, COALESCE(p.name, '') AS name, COALESCE(p.image_url, '') AS image_url, COALESCE(p.barcode, '') AS barcode, pbi.external_product_id,
			COALESCE(pbi.external_variant_id, '') AS external_variant_id,
			COALESCE(pbi.external_sku, '') AS external_sku,
			COALESCE(pbi.external_barcode, '') AS external_barcode`).
		Joins("JOIN products p ON p.id = pbi.product_id").
		Where("pbi.integration_id = ? AND pbi.deleted_at IS NULL AND pbi.external_product_id <> '' AND p.deleted_at IS NULL", integrationID).
		Order("p.sku").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	items := make([]domain.MappedItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.MappedItem{
			ProductID:         row.ProductID,
			SKU:               row.SKU,
			Name:              row.Name,
			ImageURL:          row.ImageURL,
			Barcode:           row.Barcode,
			ExternalItemID:    row.ExternalProductID,
			ExternalVariantID: row.ExternalVariantID,
			ExternalSKU:       row.ExternalSKU,
			ExternalBarcode:   row.ExternalBarcode,
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

func (r *ProductRepository) GetProductIDByExternalRef(ctx context.Context, integrationID uint, externalProductID, externalVariantID string) (string, error) {
	q := r.db.Conn(ctx).
		Table("product_business_integrations").
		Select("product_id").
		Where("integration_id = ? AND external_product_id = ? AND deleted_at IS NULL", integrationID, externalProductID)
	if externalVariantID == "" {
		q = q.Where("external_variant_id IS NULL")
	} else {
		q = q.Where("external_variant_id = ?", externalVariantID)
	}
	var productID string
	if err := q.Limit(1).Scan(&productID).Error; err != nil {
		return "", err
	}
	return productID, nil
}

func (r *ProductRepository) GetExternalRefs(ctx context.Context, productID string, integrationID uint) (string, string, error) {
	var row struct {
		ExternalProductID string
		ExternalVariantID *string
	}
	err := r.db.Conn(ctx).
		Table("product_business_integrations").
		Select("external_product_id, external_variant_id").
		Where("product_id = ? AND integration_id = ? AND deleted_at IS NULL", productID, integrationID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return "", "", err
	}
	variant := ""
	if row.ExternalVariantID != nil {
		variant = *row.ExternalVariantID
	}
	return row.ExternalProductID, variant, nil
}

func (r *ProductRepository) GetProductSKUByID(ctx context.Context, productID string, businessID uint) (string, error) {
	var sku string
	err := r.db.Conn(ctx).
		Table("products").
		Select("sku").
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", productID, businessID).
		Limit(1).
		Scan(&sku).Error
	if err != nil {
		return "", err
	}
	return sku, nil
}

func optionalRef(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (r *ProductRepository) UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
	var existing models.ProductBusinessIntegration
	err := r.db.Conn(ctx).
		Where("product_id = ? AND integration_id = ?", productID, integrationID).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		record := models.ProductBusinessIntegration{
			ProductID:         productID,
			BusinessID:        businessID,
			IntegrationID:     integrationID,
			ExternalProductID: refs.ProductID,
			ExternalVariantID: optionalRef(refs.VariantID),
			ExternalSKU:       optionalRef(refs.SKU),
			ExternalBarcode:   optionalRef(refs.Barcode),
		}
		return r.db.Conn(ctx).Create(&record).Error
	}
	if err != nil {
		return err
	}

	existing.ExternalProductID = refs.ProductID
	if v := optionalRef(refs.VariantID); v != nil {
		existing.ExternalVariantID = v
	}
	if v := optionalRef(refs.SKU); v != nil {
		existing.ExternalSKU = v
	}
	if v := optionalRef(refs.Barcode); v != nil {
		existing.ExternalBarcode = v
	}
	return r.db.Conn(ctx).Save(&existing).Error
}
