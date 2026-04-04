package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

// SearchProducts busca productos por query ILIKE en name, description, short_description
// Tabla consultada: products (gestionada por modulo products)
// Replicado localmente para evitar compartir repositorios entre modulos
func (r *repository) SearchProducts(ctx context.Context, businessID uint, query string, limit int) ([]domain.ProductSearchResult, error) {
	if limit <= 0 || limit > 20 {
		limit = 5
	}

	var products []models.Product
	searchPattern := fmt.Sprintf("%%%s%%", query)

	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Where(
			"name ILIKE ? OR description ILIKE ? OR short_description ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		).
		Limit(limit).
		Find(&products).Error

	if err != nil {
		return nil, fmt.Errorf("error searching products: %w", err)
	}

	results := make([]domain.ProductSearchResult, 0, len(products))
	for _, p := range products {
		results = append(results, mapProductToDomain(&p))
	}

	return results, nil
}

// GetProductBySKU obtiene un producto por su SKU dentro de un negocio
func (r *repository) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*domain.ProductSearchResult, error) {
	var product models.Product

	err := r.db.Conn(ctx).
		Where("business_id = ? AND sku = ? AND is_active = ? AND deleted_at IS NULL", businessID, sku, true).
		First(&product).Error

	if err != nil {
		return nil, &domain.ErrProductNotFound{SKU: sku}
	}

	result := mapProductToDomain(&product)
	return &result, nil
}

func mapProductToDomain(p *models.Product) domain.ProductSearchResult {
	return domain.ProductSearchResult{
		ID:               p.ID,
		SKU:              p.SKU,
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Price:            p.Price,
		Currency:         p.Currency,
		StockQuantity:    p.StockQuantity,
		Category:         p.Category,
		Brand:            p.Brand,
		ImageURL:         p.ImageURL,
		IsActive:         p.IsActive,
	}
}
