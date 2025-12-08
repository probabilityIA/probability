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

// Repository implementa el repositorio de productos
type Repository struct {
	db db.IDatabase
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase) domain.IRepository {
	return &Repository{
		db: database,
	}
}

// CreateProduct crea un nuevo producto en la base de datos
func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	dbProduct := mappers.ToDBProduct(product)
	if err := r.db.Conn(ctx).Create(dbProduct).Error; err != nil {
		return err
	}
	// Actualizar el ID del modelo de dominio con el ID generado
	product.ID = dbProduct.ID
	return nil
}

// GetProductByID obtiene un producto por su ID
func (r *Repository) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	var product models.Product
	err := r.db.Conn(ctx).
		Preload("Business").
		Where("id = ?", id).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return mappers.ToDomainProduct(&product), nil
}

// GetProductBySKU obtiene un producto por su SKU y BusinessID
func (r *Repository) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
	var product models.Product
	err := r.db.Conn(ctx).
		Preload("Business").
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

// ListProducts obtiene una lista paginada de productos con filtros optimizados
func (r *Repository) ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error) {
	var products []models.Product
	var total int64

	// Construir query base
	query := r.db.Conn(ctx).Model(&models.Product{})

	// Filtro por business_id
	if businessID, ok := filters["business_id"].(uint); ok && businessID > 0 {
		query = query.Where("products.business_id = ?", businessID)
	}

	// Filtro por integration_id (JOIN con Business -> Integrations)
	if integrationID, ok := filters["integration_id"].(uint); ok && integrationID > 0 {
		query = query.
			Joins("INNER JOIN business ON products.business_id = business.id").
			Joins("INNER JOIN integrations ON business.id = integrations.business_id").
			Where("integrations.id = ?", integrationID).
			Where("integrations.is_active = ?", true)
	}

	// Filtro por integration_type (JOIN con Business -> Integrations -> IntegrationType)
	if integrationType, ok := filters["integration_type"].(string); ok && integrationType != "" {
		query = query.
			Joins("INNER JOIN business ON products.business_id = business.id").
			Joins("INNER JOIN integrations ON business.id = integrations.business_id").
			Joins("INNER JOIN integration_types ON integrations.integration_type_id = integration_types.id").
			Where("integration_types.code = ?", integrationType).
			Where("integrations.is_active = ?", true)
	}

	// Filtro por SKU (búsqueda parcial, case-insensitive)
	if sku, ok := filters["sku"].(string); ok && sku != "" {
		query = query.Where("products.sku ILIKE ?", "%"+sku+"%")
	}

	// Filtro por múltiples SKUs (búsqueda exacta con IN)
	if skus, ok := filters["skus"].([]string); ok && len(skus) > 0 {
		query = query.Where("products.sku IN ?", skus)
	}

	// Filtro por nombre (búsqueda parcial, case-insensitive)
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("products.name ILIKE ?", "%"+name+"%")
	}

	// Filtro por external_id (búsqueda exacta)
	if externalID, ok := filters["external_id"].(string); ok && externalID != "" {
		query = query.Where("products.external_id = ?", externalID)
	}

	// Filtro por múltiples external_ids (búsqueda exacta con IN)
	if externalIDs, ok := filters["external_ids"].([]string); ok && len(externalIDs) > 0 {
		query = query.Where("products.external_id IN ?", externalIDs)
	}

	// Filtros de fecha (compatibilidad con formato anterior)
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("products.created_at >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("products.created_at <= ?", endDate)
	}

	// Filtros de fecha mejorados
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

	// Usar DISTINCT si hay JOINs para evitar duplicados
	hasJoins := filters["integration_id"] != nil || filters["integration_type"] != nil
	if hasJoins {
		query = query.Distinct("products.id")
	}

	// Contar total (antes de aplicar paginación y ordenamiento)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar ordenamiento
	sortBy := "products.created_at"
	if sort, ok := filters["sort_by"].(string); ok && sort != "" {
		// Mapear campos de ordenamiento
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

	// Aplicar paginación
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Precargar relaciones
	query = query.Preload("Business")

	// Ejecutar query
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	// Convertir a dominio
	domainProducts := make([]domain.Product, len(products))
	for i, product := range products {
		domainProducts[i] = *mappers.ToDomainProduct(&product)
	}

	return domainProducts, total, nil
}

// UpdateProduct actualiza un producto existente
func (r *Repository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	dbProduct := mappers.ToDBProduct(product)
	return r.db.Conn(ctx).Save(dbProduct).Error
}

// DeleteProduct elimina (soft delete) un producto
func (r *Repository) DeleteProduct(ctx context.Context, id string) error {
	return r.db.Conn(ctx).Where("id = ?", id).Delete(&models.Product{}).Error
}

// ProductExists verifica si existe un producto con el SKU para un negocio
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

// ───────────────────────────────────────────
//
//	PRODUCT-INTEGRATION MANAGEMENT
//
// ───────────────────────────────────────────

// AddProductIntegration asocia un producto con una integración
func (r *Repository) AddProductIntegration(ctx context.Context, productID string, integrationID uint, externalProductID string) (*domain.ProductBusinessIntegration, error) {
	// Verificar que el producto existe y obtener su BusinessID
	var product models.Product
	if err := r.db.Conn(ctx).Where("id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	// Verificar que la integración existe
	var integration models.Integration
	if err := r.db.Conn(ctx).Where("id = ?", integrationID).First(&integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("integration not found")
		}
		return nil, err
	}

	// Validar que la integración pertenece al mismo negocio que el producto
	if integration.BusinessID == nil || *integration.BusinessID != product.BusinessID {
		return nil, fmt.Errorf("integration does not belong to the same business as the product")
	}

	// Verificar si ya existe la asociación
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

	// Crear la asociación
	dbPI := &models.ProductBusinessIntegration{
		ProductID:         productID,
		BusinessID:        product.BusinessID,
		IntegrationID:     integrationID,
		ExternalProductID: externalProductID,
	}

	if err := r.db.Conn(ctx).Create(dbPI).Error; err != nil {
		return nil, err
	}

	return mappers.ToDomainProductIntegration(dbPI), nil
}

// RemoveProductIntegration remueve la asociación entre un producto y una integración
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

// GetProductIntegrations obtiene todas las integraciones asociadas a un producto
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

// GetProductsByIntegration obtiene todos los productos asociados a una integración
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

	// Convertir a dominio
	domainProducts := make([]domain.Product, len(products))
	for i, product := range products {
		domainProducts[i] = *mappers.ToDomainProduct(&product)
	}

	return domainProducts, nil
}

// ProductIntegrationExists verifica si existe una asociación producto-integración
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
