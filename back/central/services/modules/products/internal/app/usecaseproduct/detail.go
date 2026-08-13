package usecaseproduct

import (
	"context"
	"errors"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

const maxMovementsInDetail = 10

func (uc *UseCaseProduct) GetProductDetail(ctx context.Context, businessID uint, id string) (*domain.ProductDetailResponse, error) {
	if id == "" {
		return nil, errors.New("product ID is required")
	}

	product, err := uc.repo.GetProductByID(ctx, businessID, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, domain.ErrProductNotFound
	}

	detail := &domain.ProductDetailResponse{
		ProductResponse: *mapProductToResponse(product),
		Origins:         []domain.FieldOriginView{},
		Channels:        []domain.ProductChannel{},
		Stock:           []domain.ProductStockLevel{},
		Movements:       []domain.ProductMovement{},
	}

	origins, err := uc.repo.GetProductFieldOrigins(ctx, businessID, product.ID)
	if err != nil {
		return nil, err
	}
	if origins != nil {
		detail.Origins = origins
	}

	channels, err := uc.repo.GetProductChannels(ctx, businessID, product.ID)
	if err != nil {
		return nil, err
	}
	for i := range channels {
		channels[i].SKUMatches = sameSKU(channels[i].ExternalSKU, product.SKU)
	}
	if channels != nil {
		detail.Channels = channels
	}

	stock, err := uc.repo.GetProductStock(ctx, businessID, product.ID)
	if err != nil {
		return nil, err
	}
	if stock != nil {
		detail.Stock = stock
	}
	for _, level := range detail.Stock {
		detail.StockTotal += level.Quantity
		detail.ReservedTotal += level.ReservedQty
		detail.AvailableTotal += level.AvailableQty
	}

	movements, err := uc.repo.GetProductMovements(ctx, businessID, product.ID, maxMovementsInDetail)
	if err != nil {
		return nil, err
	}
	if movements != nil {
		detail.Movements = movements
	}

	if product.FamilyID != nil && *product.FamilyID > 0 {
		siblings, err := uc.repo.CountFamilySiblings(ctx, businessID, *product.FamilyID, product.ID)
		if err != nil {
			return nil, err
		}
		detail.SiblingVariants = siblings
	}

	return detail, nil
}

func sameSKU(channelSKU *string, ownSKU string) bool {
	if channelSKU == nil || strings.TrimSpace(*channelSKU) == "" {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(*channelSKU), strings.TrimSpace(ownSKU))
}
