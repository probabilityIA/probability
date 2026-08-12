package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) domain.IRepository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	dbProduct := mappers.ToDBProduct(product)
	if err := r.db.Conn(ctx).Create(dbProduct).Error; err != nil {
		return err
	}
	product.ID = dbProduct.ID
	return nil
}

func (r *Repository) GetProductByID(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
	var product models.Product
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Family").
		Where("id = ? AND business_id = ?", id, businessID).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return mappers.ToDomainProduct(&product), nil
}

func (r *Repository) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
	var product models.Product
	err := r.db.Conn(ctx).
		Preload("Business").
		Preload("Family").
		Where("business_id = ? AND sku = ?", businessID, sku).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return mappers.ToDomainProduct(&product), nil
}

func (r *Repository) ListProducts(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Conn(ctx).Model(&models.Product{}).Where("products.business_id = ?", businessID)

	if integrationID, ok := filters["integration_id"].(uint); ok && integrationID > 0 {
		query = query.
			Joins("INNER JOIN business ON products.business_id = business.id").
			Joins("INNER JOIN integrations ON business.id = integrations.business_id").
			Where("integrations.id = ?", integrationID).
			Where("integrations.is_active = ?", true)
	}

	if integrationType, ok := filters["integration_type"].(string); ok && integrationType != "" {
		query = query.
			Joins("INNER JOIN business ON products.business_id = business.id").
			Joins("INNER JOIN integrations ON business.id = integrations.business_id").
			Joins("INNER JOIN integration_types ON integrations.integration_type_id = integration_types.id").
			Where("integration_types.code = ?", integrationType).
			Where("integrations.is_active = ?", true)
	}

	if sku, ok := filters["sku"].(string); ok && sku != "" {
		query = query.Where("products.sku ILIKE ?", "%"+sku+"%")
	}

	if skus, ok := filters["skus"].([]string); ok && len(skus) > 0 {
		query = query.Where("products.sku IN ?", skus)
	}

	if familyID, ok := filters["family_id"].(uint); ok && familyID > 0 {
		query = query.Where("products.family_id = ?", familyID)
	}

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("products.name ILIKE ?", "%"+name+"%")
	}

	if barcode, ok := filters["barcode"].(string); ok && barcode != "" {
		query = query.Where("products.barcode = ?", barcode)
	}

	if externalID, ok := filters["external_id"].(string); ok && externalID != "" {
		query = query.Where("products.external_id = ?", externalID)
	}

	if externalIDs, ok := filters["external_ids"].([]string); ok && len(externalIDs) > 0 {
		query = query.Where("products.external_id IN ?", externalIDs)
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("products.created_at >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("products.created_at <= ?", endDate)
	}

	if createdAfter, ok := filters["created_after"].(string); ok && createdAfter != "" {
		query = query.Where("products.created_at >= ?", createdAfter)
	}

	if createdBefore, ok := filters["created_before"].(string); ok && createdBefore != "" {
		query = query.Where("products.created_at <= ?", createdBefore)
	}

	if updatedAfter, ok := filters["updated_after"].(string); ok && updatedAfter != "" {
		query = query.Where("products.updated_at >= ?", updatedAfter)
	}

	if updatedBefore, ok := filters["updated_before"].(string); ok && updatedBefore != "" {
		query = query.Where("products.updated_at <= ?", updatedBefore)
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("products.status = ?", status)
	}

	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("products.category ILIKE ?", "%"+category+"%")
	}

	if brand, ok := filters["brand"].(string); ok && brand != "" {
		query = query.Where("products.brand ILIKE ?", "%"+brand+"%")
	}

	if hasFamily, ok := filters["has_family"].(bool); ok {
		if hasFamily {
			query = query.Where("products.family_id IS NOT NULL")
		} else {
			query = query.Where("products.family_id IS NULL")
		}
	}

	if priceMin, ok := filters["price_min"].(float64); ok {
		query = query.Where("products.price >= ?", priceMin)
	}

	if priceMax, ok := filters["price_max"].(float64); ok {
		query = query.Where("products.price <= ?", priceMax)
	}

	if stockMin, ok := filters["stock_min"].(int); ok {
		query = query.Where("products.stock_quantity >= ?", stockMin)
	}

	if stockMax, ok := filters["stock_max"].(int); ok {
		query = query.Where("products.stock_quantity <= ?", stockMax)
	}

	if weightMin, ok := filters["weight_min"].(float64); ok {
		query = query.Where("products.weight >= ?", weightMin)
	}

	if weightMax, ok := filters["weight_max"].(float64); ok {
		query = query.Where("products.weight <= ?", weightMax)
	}

	hasJoins := filters["integration_id"] != nil || filters["integration_type"] != nil
	if hasJoins {
		query = query.Distinct("products.id")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "products.created_at"
	if sort, ok := filters["sort_by"].(string); ok && sort != "" {
		sortFieldMap := map[string]string{
			"id":          "products.id",
			"sku":         "products.sku",
			"name":        "products.name",
			"created_at":  "products.created_at",
			"updated_at":  "products.updated_at",
			"business_id": "products.business_id",
		}
		if mappedField, exists := sortFieldMap[sort]; exists {
			sortBy = mappedField
		}
	}

	sortOrder := "desc"
	if order, ok := filters["sort_order"].(string); ok && order != "" {
		sortOrder = order
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	query = query.Preload("Business").Preload("Family")

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	domainProducts := make([]domain.Product, len(products))
	for i, product := range products {
		domainProducts[i] = *mappers.ToDomainProduct(&product)
	}

	return domainProducts, total, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	dbProduct := mappers.ToDBProduct(product)
	return r.db.Conn(ctx).Save(dbProduct).Error
}

func (r *Repository) DeleteProduct(ctx context.Context, businessID uint, id string) error {
	return r.db.Conn(ctx).Where("id = ? AND business_id = ?", id, businessID).Delete(&models.Product{}).Error
}

func (r *Repository) ListProductsByFamilyID(ctx context.Context, businessID uint, familyID uint) ([]domain.Product, error) {
	var products []models.Product
	err := r.db.Conn(ctx).
		Where("business_id = ? AND family_id = ?", businessID, familyID).
		Preload("Business").
		Preload("Family").
		Order("created_at asc").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	domainProducts := make([]domain.Product, len(products))
	for i := range products {
		domainProducts[i] = *mappers.ToDomainProduct(&products[i])
	}

	return domainProducts, nil
}

func (r *Repository) ProductExists(ctx context.Context, businessID uint, sku string) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Model(&models.Product{}).
		Where("business_id = ? AND sku = ?", businessID, sku).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) VariantExistsInFamily(ctx context.Context, businessID uint, familyID uint, variantSignature string, excludeProductID *string) (bool, error) {
	if variantSignature == "" {
		return false, nil
	}

	query := r.db.Conn(ctx).
		Model(&models.Product{}).
		Where("business_id = ? AND family_id = ? AND variant_signature = ?", businessID, familyID, variantSignature)

	if excludeProductID != nil && *excludeProductID != "" {
		query = query.Where("id <> ?", *excludeProductID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) CreateProductFamily(ctx context.Context, family *domain.ProductFamily) error {
	dbFamily := mappers.ToDBProductFamily(family)
	if err := r.db.Conn(ctx).Create(dbFamily).Error; err != nil {
		return err
	}

	family.ID = dbFamily.ID
	family.CreatedAt = dbFamily.CreatedAt
	family.UpdatedAt = dbFamily.UpdatedAt
	return nil
}

func (r *Repository) GetProductFamilyByID(ctx context.Context, businessID uint, familyID uint) (*domain.ProductFamily, error) {
	var family models.ProductFamily
	err := r.db.Conn(ctx).
		Model(&models.ProductFamily{}).
		Select("product_families.*, COUNT(products.id) AS variant_count").
		Joins("LEFT JOIN products ON products.family_id = product_families.id AND products.deleted_at IS NULL").
		Where("product_families.id = ? AND product_families.business_id = ?", familyID, businessID).
		Group("product_families.id").
		First(&family).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductFamilyNotFound
		}
		return nil, err
	}

	return mappers.ToDomainProductFamily(&family), nil
}

func (r *Repository) ListProductFamilies(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.ProductFamily, int64, error) {
	var families []models.ProductFamily
	var total int64

	base := r.db.Conn(ctx).Model(&models.ProductFamily{}).Where("product_families.business_id = ?", businessID)

	if name, ok := filters["name"].(string); ok && name != "" {
		base = base.Where("product_families.name ILIKE ?", "%"+name+"%")
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		base = base.Where("product_families.category ILIKE ?", "%"+category+"%")
	}
	if brand, ok := filters["brand"].(string); ok && brand != "" {
		base = base.Where("product_families.brand ILIKE ?", "%"+brand+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		base = base.Where("product_families.status = ?", status)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "product_families.created_at"
	if sort, ok := filters["sort_by"].(string); ok && sort != "" {
		sortFieldMap := map[string]string{
			"id":         "product_families.id",
			"name":       "product_families.name",
			"created_at": "product_families.created_at",
			"updated_at": "product_families.updated_at",
		}
		if mappedField, exists := sortFieldMap[sort]; exists {
			sortBy = mappedField
		}
	}

	sortOrder := "desc"
	if order, ok := filters["sort_order"].(string); ok && order != "" {
		sortOrder = order
	}

	offset := (page - 1) * pageSize
	query := r.db.Conn(ctx).
		Model(&models.ProductFamily{}).
		Select("product_families.*, COUNT(products.id) AS variant_count").
		Joins("LEFT JOIN products ON products.family_id = product_families.id AND products.deleted_at IS NULL").
		Where("product_families.business_id = ?", businessID)

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("product_families.name ILIKE ?", "%"+name+"%")
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("product_families.category ILIKE ?", "%"+category+"%")
	}
	if brand, ok := filters["brand"].(string); ok && brand != "" {
		query = query.Where("product_families.brand ILIKE ?", "%"+brand+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("product_families.status = ?", status)
	}

	query = query.Group("product_families.id").Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(pageSize)

	if err := query.Find(&families).Error; err != nil {
		return nil, 0, err
	}

	result := make([]domain.ProductFamily, len(families))
	for i := range families {
		result[i] = *mappers.ToDomainProductFamily(&families[i])
	}

	return result, total, nil
}

func (r *Repository) UpdateProductFamily(ctx context.Context, family *domain.ProductFamily) error {
	dbFamily := mappers.ToDBProductFamily(family)
	return r.db.Conn(ctx).Save(dbFamily).Error
}

func (r *Repository) DeleteProductFamily(ctx context.Context, businessID uint, familyID uint) error {
	return r.db.Conn(ctx).
		Where("id = ? AND business_id = ?", familyID, businessID).
		Delete(&models.ProductFamily{}).Error
}

func (r *Repository) AddProductIntegration(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error) {
	var product models.Product
	if err := r.db.Conn(ctx).Where("id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	var integration models.Integration
	if err := r.db.Conn(ctx).Where("id = ?", integrationID).First(&integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("integration not found")
		}
		return nil, err
	}

	if integration.BusinessID == nil || *integration.BusinessID != product.BusinessID {
		return nil, fmt.Errorf("integration does not belong to the same business as the product")
	}

	var existingCount int64
	err := r.db.Conn(ctx).
		Model(&models.ProductBusinessIntegration{}).
		Where("product_id = ? AND integration_id = ?", productID, integrationID).
		Count(&existingCount).Error
	if err != nil {
		return nil, err
	}
	if existingCount > 0 {
		return nil, fmt.Errorf("product is already associated with this integration")
	}

	dbPI := &models.ProductBusinessIntegration{
		ProductID:         productID,
		BusinessID:        product.BusinessID,
		IntegrationID:     integrationID,
		ExternalProductID: externalProductID,
		ExternalVariantID: externalVariantID,
		ExternalSKU:       externalSKU,
		ExternalBarcode:   externalBarcode,
	}

	if err := r.db.Conn(ctx).Create(dbPI).Error; err != nil {
		return nil, err
	}

	return mappers.ToDomainProductIntegration(dbPI), nil
}

func (r *Repository) RemoveProductIntegration(ctx context.Context, productID string, integrationID uint) error {
	result := r.db.Conn(ctx).
		Where("product_id = ? AND integration_id = ?", productID, integrationID).
		Delete(&models.ProductBusinessIntegration{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product integration association not found")
	}

	return nil
}

func (r *Repository) GetProductIntegrations(ctx context.Context, productID string) ([]domain.ProductBusinessIntegration, error) {
	var integrations []models.ProductBusinessIntegration
	err := r.db.Conn(ctx).
		Preload("Integration").
		Preload("Integration.IntegrationType").
		Where("product_id = ?", productID).
		Find(&integrations).Error

	if err != nil {
		return nil, err
	}

	return mappers.ToDomainProductIntegrations(integrations), nil
}

func (r *Repository) GetProductsByIntegration(ctx context.Context, integrationID uint) ([]domain.Product, error) {
	var products []models.Product
	err := r.db.Conn(ctx).
		Joins("INNER JOIN product_business_integrations ON products.id = product_business_integrations.product_id").
		Where("product_business_integrations.integration_id = ?", integrationID).
		Preload("Business").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	domainProducts := make([]domain.Product, len(products))
	for i, product := range products {
		domainProducts[i] = *mappers.ToDomainProduct(&product)
	}

	return domainProducts, nil
}

func (r *Repository) ProductIntegrationExists(ctx context.Context, productID string, integrationID uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Model(&models.ProductBusinessIntegration{}).
		Where("product_id = ? AND integration_id = ?", productID, integrationID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) HasFamilyActiveVariants(ctx context.Context, businessID uint, familyID uint) (bool, error) {
	var count int64
	err := r.db.Conn(ctx).
		Model(&models.Product{}).
		Where("business_id = ? AND family_id = ? AND deleted_at IS NULL", businessID, familyID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) UpdateProductIntegration(ctx context.Context, productID string, integrationID uint, req *domain.UpdateProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	var existing models.ProductBusinessIntegration
	err := r.db.Conn(ctx).
		Where("product_id = ? AND integration_id = ?", productID, integrationID).
		First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductIntegrationNotFound
		}
		return nil, err
	}

	if req.ExternalProductID != nil {
		existing.ExternalProductID = *req.ExternalProductID
	}
	if req.ExternalVariantID != nil {
		existing.ExternalVariantID = req.ExternalVariantID
	}
	if req.ExternalSKU != nil {
		existing.ExternalSKU = req.ExternalSKU
	}
	if req.ExternalBarcode != nil {
		existing.ExternalBarcode = req.ExternalBarcode
	}

	if err := r.db.Conn(ctx).Save(&existing).Error; err != nil {
		return nil, err
	}

	return mappers.ToDomainProductIntegration(&existing), nil
}

func (r *Repository) LookupProductByExternalRef(ctx context.Context, businessID uint, integrationID uint, externalVariantID, externalSKU, externalProductID, externalBarcode *string) (*domain.Product, error) {
	type candidate struct {
		field string
		value *string
	}

	candidates := []candidate{
		{"external_variant_id", externalVariantID},
		{"external_sku", externalSKU},
		{"external_barcode", externalBarcode},
		{"external_product_id", externalProductID},
	}

	for _, c := range candidates {
		if c.value == nil || *c.value == "" {
			continue
		}
		var products []models.Product
		err := r.db.Conn(ctx).
			Model(&models.Product{}).
			Joins("INNER JOIN product_business_integrations ON product_business_integrations.product_id = products.id").
			Where("products.business_id = ? AND product_business_integrations.integration_id = ? AND product_business_integrations."+c.field+" = ?", businessID, integrationID, *c.value).
			Limit(2).
			Find(&products).Error
		if err != nil {
			return nil, err
		}
		if len(products) == 1 {
			return mappers.ToDomainProduct(&products[0]), nil
		}
	}

	return nil, nil
}

func (r *Repository) ListSKUs(ctx context.Context, businessID uint, prefix string) ([]string, error) {
	var skus []string
	query := r.db.Conn(ctx).
		Model(&models.Product{}).
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Distinct("sku").
		Order("sku ASC")

	if prefix != "" {
		query = query.Where("LOWER(sku) LIKE LOWER(?)", prefix+"%")
	}

	err := query.Pluck("sku", &skus).Error
	if err != nil {
		return nil, err
	}

	return skus, nil
}

func (r *Repository) GetNextSKU(ctx context.Context, businessID uint, prefix string) (string, error) {
	var skus []string
	query := r.db.Conn(ctx).
		Model(&models.Product{}).
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Distinct("sku").
		Order("sku ASC")

	if prefix != "" {
		query = query.Where("LOWER(sku) LIKE LOWER(?)", prefix+"%")
	}

	err := query.Pluck("sku", &skus).Error
	if err != nil {
		return "", err
	}

	if len(skus) == 0 {
		return fmt.Sprintf("%s-%03d", prefix, 1), nil
	}

	nextNum := 1
	for _, sku := range skus {
		var num int
		if _, err := fmt.Sscanf(sku, prefix+"-%d", &num); err == nil {
			if num == nextNum {
				nextNum++
			} else if num > nextNum {
				break
			}
		}
	}

	return fmt.Sprintf("%s-%03d", prefix, nextNum), nil
}

func (r *Repository) GetNextSKUBatch(ctx context.Context, businessID uint, prefix string, count int) ([]string, error) {
	var skus []string
	query := r.db.Conn(ctx).
		Model(&models.Product{}).
		Where("business_id = ? AND deleted_at IS NULL", businessID).
		Distinct("sku").
		Order("sku ASC")

	if prefix != "" {
		query = query.Where("LOWER(sku) LIKE LOWER(?)", prefix+"%")
	}

	err := query.Pluck("sku", &skus).Error
	if err != nil {
		return nil, err
	}

	if count <= 0 {
		count = 1
	}

	result := make([]string, 0, count)
	nextNum := 1

	for len(result) < count {
		found := false
		for _, sku := range skus {
			var num int
			if _, err := fmt.Sscanf(sku, prefix+"-%d", &num); err == nil {
				if num == nextNum {
					found = true
					break
				}
			}
		}

		if !found {
			result = append(result, fmt.Sprintf("%s-%03d", prefix, nextNum))
		}
		nextNum++
	}

	return result, nil
}
