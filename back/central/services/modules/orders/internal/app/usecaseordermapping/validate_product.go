package usecaseordermapping

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
)

// GetOrCreateProduct verifica si el producto existe, si no, lo crea
func (uc *UseCaseOrderMapping) GetOrCreateProduct(ctx context.Context, businessID uint, itemDTO domain.CanonicalOrderItemDTO) (*domain.Product, error) {
	if itemDTO.ProductSKU == "" {
		return nil, fmt.Errorf("product SKU is required")
	}

	// 1. Buscar producto existente
	product, err := uc.repo.GetProductBySKU(ctx, businessID, itemDTO.ProductSKU)
	if err != nil {
		return nil, fmt.Errorf("error searching product: %w", err)
	}

	if product != nil {
		return product, nil
	}

	// 2. Crear nuevo producto si no existe
	// Nota: Usamos el ProductID externo como ExternalID si está disponible
	var externalID string
	if itemDTO.ProductID != nil {
		externalID = *itemDTO.ProductID
	}

	newProduct := &domain.Product{
		BusinessID: businessID,
		SKU:        itemDTO.ProductSKU,
		Name:       itemDTO.ProductName,
		ExternalID: externalID,
	}

	if err := uc.repo.CreateProduct(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return newProduct, nil
}
