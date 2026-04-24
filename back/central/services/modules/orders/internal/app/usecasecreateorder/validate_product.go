package usecasecreateorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseCreateOrder) GetOrCreateProduct(ctx context.Context, businessID uint, integrationID uint, itemDTO dtos.ProbabilityOrderItemDTO) (*entities.Product, error) {
	if itemDTO.ProductSKU == "" && itemDTO.VariantID == nil && itemDTO.ProductID == nil {
		return nil, fmt.Errorf("product identifier is required")
	}

	product, err := uc.repo.ResolveProductForOrderItem(ctx, businessID, integrationID, itemDTO)
	if err != nil {
		return nil, fmt.Errorf("error searching product: %w", err)
	}

	if product != nil {
		if itemDTO.UnitPrice > 0 && itemDTO.UnitPrice != product.Price {
			if err := uc.repo.UpdateProductPrice(ctx, product.ID, itemDTO.UnitPrice); err != nil {
				uc.logger.Warn(ctx).Err(err).Str("product_id", product.ID).Msg("failed to update product price")
			} else {
				product.Price = itemDTO.UnitPrice
			}
		}
		return product, nil
	}

	var externalID string
	if itemDTO.ProductID != nil {
		externalID = *itemDTO.ProductID
	} else if itemDTO.VariantID != nil {
		externalID = *itemDTO.VariantID
	}

	sku := itemDTO.ProductSKU
	if sku == "" {
		if itemDTO.VariantID != nil && *itemDTO.VariantID != "" {
			sku = fmt.Sprintf("VAR-%s", *itemDTO.VariantID)
		} else if itemDTO.ProductID != nil && *itemDTO.ProductID != "" {
			sku = fmt.Sprintf("PROD-%s", *itemDTO.ProductID)
		}
	}

	if sku == "" {
		uc.logger.Warn(ctx).
			Interface("item", itemDTO.ProductName).
			Msg("unmapped variant: no resolvable identifier, skipping integration mapping")
	}

	newProduct := &entities.Product{
		BusinessID: businessID,
		SKU:        sku,
		Name:       itemDTO.ProductName,
		ExternalID: externalID,
		Price:      itemDTO.UnitPrice,
	}

	if err := uc.repo.CreateProduct(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	if integrationID > 0 && sku != "" {
		if err := uc.repo.UpsertProductIntegrationMapping(ctx, newProduct.ID, businessID, integrationID, itemDTO); err != nil {
			uc.logger.Warn(ctx).Err(err).Str("product_id", newProduct.ID).Uint("integration_id", integrationID).Msg("failed to persist product integration mapping")
		}
	}

	return newProduct, nil
}
