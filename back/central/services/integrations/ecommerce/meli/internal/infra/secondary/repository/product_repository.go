package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
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
		log: logger.WithModule("meli.product_repository"),
	}
}

func NewInventory(database db.IDatabase, logger log.ILogger) domain.IInventoryRepository {
	return &ProductRepository{
		db:  database,
		log: logger.WithModule("meli.inventory_repository"),
	}
}

func (r *ProductRepository) ListMappedItems(ctx context.Context, integrationID uint) ([]domain.MappedItem, error) {
	var rows []struct {
		ProductID         string
		SKU               string
		Name              string
		Barcode           string
		ExternalProductID string
		ExternalVariantID string
		ExternalSKU       string
		ExternalBarcode   string
		LogisticType      string
		LastPushedQty     *int
	}
	err := r.db.Conn(ctx).
		Table("product_business_integrations AS pbi").
		Select(`pbi.product_id, p.sku, COALESCE(p.name, '') AS name, COALESCE(p.barcode, '') AS barcode, pbi.external_product_id,
			COALESCE(pbi.external_variant_id, '') AS external_variant_id,
			COALESCE(pbi.external_sku, '') AS external_sku,
			COALESCE(pbi.external_barcode, '') AS external_barcode,
			COALESCE(pbi.external_logistic_type, '') AS logistic_type,
			pbi.last_pushed_qty`).
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
			Barcode:           row.Barcode,
			ExternalItemID:    row.ExternalProductID,
			ExternalVariantID: row.ExternalVariantID,
			ExternalSKU:       row.ExternalSKU,
			ExternalBarcode:   row.ExternalBarcode,
			LogisticType:      row.LogisticType,
			LastPushedQty:     row.LastPushedQty,
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

func (r *ProductRepository) GetInventoryByWarehouses(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]map[uint]int, error) {
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
		Brand             string
		Category          string
		MeliCategoryID    string
		VariantAttributes []byte
	}

	err := r.db.Conn(ctx).
		Table("products AS p").
		Select(`p.id, p.sku, COALESCE(p.barcode, '') AS barcode, p.external_id, p.name, p.description, p.price, p.stock_quantity, p.track_inventory, p.image_url,
			COALESCE(NULLIF(p.brand, ''), NULLIF(f.brand, ''), '') AS brand,
			COALESCE(NULLIF(p.category, ''), NULLIF(f.category, ''), '') AS category,
			COALESCE(p.channel_categories->>'meli', '') AS meli_category_id,
			p.variant_attributes AS variant_attributes`).
		Joins("LEFT JOIN product_families f ON f.id = p.family_id AND f.deleted_at IS NULL").
		Where("p.business_id = ? AND p.deleted_at IS NULL AND p.is_active = ?", businessID, true).
		Order("p.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	products := make([]domain.ProductForSync, 0, len(rows))
	for _, row := range rows {
		products = append(products, domain.ProductForSync{
			ID:             row.ID,
			SKU:            row.SKU,
			Barcode:        row.Barcode,
			ExternalID:     row.ExternalID,
			Name:           row.Name,
			Description:    row.Description,
			Price:          row.Price,
			StockQuantity:  row.StockQuantity,
			TrackInventory: row.TrackInventory,
			ImageURL:       row.ImageURL,
			Brand:          row.Brand,
			Category:       row.Category,
			MeliCategoryID: row.MeliCategoryID,
			VariantAttrs:   decodeVariantAttrs(row.VariantAttributes),
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

func optionalRef(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (r *ProductRepository) UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
	var existing models.ProductBusinessIntegration
	query := r.db.Conn(ctx).
		Where("product_id = ? AND integration_id = ?", productID, integrationID)
	if refs.VariantID != "" {
		query = query.Where("external_variant_id = ?", refs.VariantID)
	} else {
		query = query.Where("external_variant_id IS NULL OR external_variant_id = ''")
	}
	err := query.First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		record := models.ProductBusinessIntegration{
			ProductID:            productID,
			BusinessID:           businessID,
			IntegrationID:        integrationID,
			ExternalProductID:    refs.ProductID,
			ExternalVariantID:    optionalRef(refs.VariantID),
			ExternalSKU:          optionalRef(refs.SKU),
			ExternalBarcode:      optionalRef(refs.Barcode),
			ExternalLogisticType: optionalRef(refs.LogisticType),
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
	if v := optionalRef(refs.LogisticType); v != nil {
		existing.ExternalLogisticType = v
	}
	return r.db.Conn(ctx).Save(&existing).Error
}

func decodeVariantAttrs(raw []byte) map[string]string {
	if len(raw) == 0 {
		return nil
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil
	}
	attrs := make(map[string]string, len(parsed))
	for key, value := range parsed {
		if text, ok := value.(string); ok && strings.TrimSpace(text) != "" {
			attrs[strings.ToUpper(strings.TrimSpace(key))] = text
		}
	}
	return attrs
}

func (r *ProductRepository) UpdateLogisticType(ctx context.Context, integrationID uint, externalItemID, logisticType string) error {
	if externalItemID == "" || logisticType == "" {
		return nil
	}
	return r.db.Conn(ctx).
		Table("product_business_integrations").
		Where("integration_id = ? AND external_product_id = ? AND deleted_at IS NULL", integrationID, externalItemID).
		Where("external_logistic_type IS DISTINCT FROM ?", logisticType).
		Update("external_logistic_type", logisticType).Error
}

func variantScope(q *gorm.DB, variantID string) *gorm.DB {
	if variantID != "" {
		return q.Where("external_variant_id = ?", variantID)
	}
	return q.Where("external_variant_id IS NULL OR external_variant_id = ''")
}

func (r *ProductRepository) MarkPushedQty(ctx context.Context, integrationID uint, productID, variantID string, quantity int) error {
	q := r.db.Conn(ctx).
		Table("product_business_integrations").
		Where("integration_id = ? AND product_id = ? AND deleted_at IS NULL", integrationID, productID)
	return variantScope(q, variantID).Update("last_pushed_qty", quantity).Error
}

func (r *ProductRepository) GetPushState(ctx context.Context, integrationID uint, productID, variantID string) (*domain.PushState, error) {
	var row struct {
		LogisticType  string
		LastPushedQty *int
	}
	q := r.db.Conn(ctx).
		Table("product_business_integrations").
		Select("COALESCE(external_logistic_type, '') AS logistic_type, last_pushed_qty").
		Where("integration_id = ? AND product_id = ? AND deleted_at IS NULL", integrationID, productID)
	err := variantScope(q, variantID).Limit(1).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return &domain.PushState{LogisticType: row.LogisticType, LastPushedQty: row.LastPushedQty}, nil
}
