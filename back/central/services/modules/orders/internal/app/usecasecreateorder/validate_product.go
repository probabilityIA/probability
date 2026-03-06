package usecasecreateorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// GetOrCreateProduct verifica si el producto existe, si no, lo crea
func (uc *UseCaseCreateOrder) GetOrCreateProduct(ctx context.Context, businessID uint, itemDTO dtos.ProbabilityOrderItemDTO) (*entities.Product, error) {
	if itemDTO.ProductSKU == "" {
		return nil, fmt.Errorf("product SKU is required")
	}

	// 1. Buscar producto existente
	product, err := uc.repo.GetProductBySKU(ctx, businessID, itemDTO.ProductSKU)
	if err != nil {
		return nil, fmt.Errorf("error searching product: %w", err)
	}

	if product != nil {
		// Update price if it changed
		if itemDTO.UnitPrice > 0 && itemDTO.UnitPrice != product.Price {
			if err := uc.repo.UpdateProductPrice(ctx, product.ID, itemDTO.UnitPrice); err != nil {
				uc.logger.Warn(ctx).Err(err).Str("product_id", product.ID).Msg("failed to update product price")
			} else {
				product.Price = itemDTO.UnitPrice
			}
		}
		return product, nil
	}

	// 2. Crear nuevo producto si no existe
	// Nota: Usamos el ProductID externo como ExternalID si está disponible
	var externalID string
	if itemDTO.ProductID != nil {
		externalID = *itemDTO.ProductID
	}

	newProduct := &entities.Product{
		BusinessID: businessID,
		SKU:        itemDTO.ProductSKU,
		Name:       itemDTO.ProductName,
		ExternalID: externalID,
		Price:      itemDTO.UnitPrice,
	}

	if err := uc.repo.CreateProduct(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return newProduct, nil
}
