package usecaseorder

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseOrder) linkItemsAndComputePackage(ctx context.Context, order *entities.ProbabilityOrder, items []entities.ProbabilityOrderItem) {
	if order.BusinessID == nil || *order.BusinessID == 0 {
		return
	}
	businessID := *order.BusinessID

	var totalWeight, maxLength, maxWidth, maxHeight float64
	totalQuantity := 0

	for i := range items {
		item := &items[i]
		totalQuantity += item.Quantity

		var product *entities.Product
		if item.ProductID != nil && *item.ProductID != "" {
			p, err := uc.repo.GetProductByID(ctx, businessID, *item.ProductID)
			if err != nil {
				uc.logger.Warn(ctx).Err(err).Str("order_id", order.ID).Str("product_id", *item.ProductID).Msg("failed to load product by id for order item")
			}
			product = p
		}
		if product == nil && item.ProductSKU != "" {
			p, err := uc.repo.GetProductBySKU(ctx, businessID, item.ProductSKU)
			if err != nil {
				uc.logger.Warn(ctx).Err(err).Str("order_id", order.ID).Str("sku", item.ProductSKU).Msg("failed to resolve product by sku for order item")
			}
			if p != nil {
				product = p
				item.ProductID = &p.ID
			}
		}

		if product != nil {
			if product.Weight != nil {
				totalWeight += *product.Weight * float64(item.Quantity)
			}
			if product.Length != nil && *product.Length > maxLength {
				maxLength = *product.Length
			}
			if product.Width != nil && *product.Width > maxWidth {
				maxWidth = *product.Width
			}
			if product.Height != nil && *product.Height > maxHeight {
				maxHeight = *product.Height
			}
		}
	}

	uc.applyShippingPackageDimensions(ctx, order, totalQuantity, totalWeight, maxLength, maxWidth, maxHeight)
}

func (uc *UseCaseOrder) applyShippingPackageDimensions(ctx context.Context, order *entities.ProbabilityOrder, totalQuantity int, totalWeight, maxLength, maxWidth, maxHeight float64) {
	var config *entities.ShippingPackageConfig
	if order.BusinessID != nil && *order.BusinessID > 0 {
		cfg, err := uc.repo.GetShippingPackageConfig(ctx, *order.BusinessID, order.WarehouseID)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).Str("order_id", order.ID).Msg("failed to load shipping package config")
		} else {
			config = cfg
		}
	}

	if config != nil && config.Strategy == entities.ShippingPackageStrategyStandardBox {
		box := config.SelectBox(totalQuantity, maxLength, maxWidth, maxHeight)
		if box != nil {
			if box.Weight != nil {
				order.Weight = box.Weight
			}
			if box.Length != nil {
				order.Length = box.Length
			}
			if box.Width != nil {
				order.Width = box.Width
			}
			if box.Height != nil {
				order.Height = box.Height
			}
			return
		}
	}

	if totalWeight > 0 {
		order.Weight = &totalWeight
		order.Length = &maxLength
		order.Width = &maxWidth
		order.Height = &maxHeight
	}
}
