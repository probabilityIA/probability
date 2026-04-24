package usecasefamily

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (uc *UseCaseFamily) CreateProductFamily(ctx context.Context, req *domain.CreateProductFamilyStandaloneRequest) (*domain.ProductFamilyResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	family := &domain.ProductFamily{
		BusinessID:  req.BusinessID,
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		Slug:        req.Slug,
		Category:    req.Category,
		Brand:       req.Brand,
		ImageURL:    req.ImageURL,
		Status:      req.Status,
		IsActive:    isActive,
		VariantAxes: req.VariantAxes,
		Metadata:    req.Metadata,
	}

	if err := uc.repo.CreateProductFamily(ctx, family); err != nil {
		return nil, fmt.Errorf("error creating product family: %w", err)
	}

	return mapProductFamilyToResponse(family), nil
}

func (uc *UseCaseFamily) GetProductFamilyByID(ctx context.Context, businessID uint, familyID uint) (*domain.ProductFamilyResponse, error) {
	family, err := uc.repo.GetProductFamilyByID(ctx, businessID, familyID)
	if err != nil {
		return nil, fmt.Errorf("error getting product family: %w", err)
	}

	response := mapProductFamilyToResponse(family)
	variants, err := uc.repo.ListProductsByFamilyID(ctx, businessID, familyID)
	if err != nil {
		return nil, fmt.Errorf("error getting family variants: %w", err)
	}
	response.Variants = mapProductsToResponses(variants)

	return response, nil
}

func (uc *UseCaseFamily) ListProductFamilies(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) (*domain.ProductFamiliesListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	families, total, err := uc.repo.ListProductFamilies(ctx, businessID, page, pageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing product families: %w", err)
	}

	responses := make([]domain.ProductFamilyResponse, len(families))
	for i := range families {
		responses[i] = *mapProductFamilyToResponse(&families[i])
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	return &domain.ProductFamiliesListResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (uc *UseCaseFamily) UpdateProductFamily(ctx context.Context, businessID uint, familyID uint, req *domain.UpdateProductFamilyRequest) (*domain.ProductFamilyResponse, error) {
	if familyID == 0 {
		return nil, errors.New("product family ID is required")
	}

	family, err := uc.repo.GetProductFamilyByID(ctx, businessID, familyID)
	if err != nil {
		return nil, fmt.Errorf("error getting product family: %w", err)
	}

	if req.Name != nil {
		family.Name = *req.Name
	}
	if req.Title != nil {
		family.Title = *req.Title
	}
	if req.Description != nil {
		family.Description = *req.Description
	}
	if req.Slug != nil {
		family.Slug = *req.Slug
	}
	if req.Category != nil {
		family.Category = *req.Category
	}
	if req.Brand != nil {
		family.Brand = *req.Brand
	}
	if req.ImageURL != nil {
		family.ImageURL = *req.ImageURL
	}
	if req.Status != nil {
		family.Status = *req.Status
	}
	if req.IsActive != nil {
		family.IsActive = *req.IsActive
	}
	if req.VariantAxes != nil {
		family.VariantAxes = req.VariantAxes
	}
	if req.Metadata != nil {
		family.Metadata = req.Metadata
	}

	if err := uc.repo.UpdateProductFamily(ctx, family); err != nil {
		return nil, fmt.Errorf("error updating product family: %w", err)
	}

	return mapProductFamilyToResponse(family), nil
}

func (uc *UseCaseFamily) DeleteProductFamily(ctx context.Context, businessID uint, familyID uint) error {
	if familyID == 0 {
		return errors.New("product family ID is required")
	}

	if _, err := uc.repo.GetProductFamilyByID(ctx, businessID, familyID); err != nil {
		return err
	}

	hasVariants, err := uc.repo.HasFamilyActiveVariants(ctx, businessID, familyID)
	if err != nil {
		return fmt.Errorf("error checking family variants: %w", err)
	}
	if hasVariants {
		return domain.ErrFamilyHasActiveVariants
	}

	if err := uc.repo.DeleteProductFamily(ctx, businessID, familyID); err != nil {
		return fmt.Errorf("error deleting product family: %w", err)
	}

	return nil
}

func mapProductFamilyToResponse(family *domain.ProductFamily) *domain.ProductFamilyResponse {
	if family == nil {
		return nil
	}

	return &domain.ProductFamilyResponse{
		ID:           family.ID,
		BusinessID:   family.BusinessID,
		Name:         family.Name,
		Title:        family.Title,
		Description:  family.Description,
		Slug:         family.Slug,
		Category:     family.Category,
		Brand:        family.Brand,
		ImageURL:     family.ImageURL,
		Status:       family.Status,
		IsActive:     family.IsActive,
		VariantAxes:  family.VariantAxes,
		Metadata:     family.Metadata,
		VariantCount: family.VariantCount,
		CreatedAt:    family.CreatedAt,
		UpdatedAt:    family.UpdatedAt,
	}
}

func mapProductsToResponses(products []domain.Product) []domain.ProductResponse {
	if len(products) == 0 {
		return nil
	}

	responses := make([]domain.ProductResponse, len(products))
	for i := range products {
		product := products[i]
		responses[i] = domain.ProductResponse{
			ID:                product.ID,
			CreatedAt:         product.CreatedAt,
			UpdatedAt:         product.UpdatedAt,
			DeletedAt:         product.DeletedAt,
			BusinessID:        product.BusinessID,
			SKU:               product.SKU,
			ExternalID:        product.ExternalID,
			Barcode:           product.Barcode,
			FamilyID:          product.FamilyID,
			Name:              product.Name,
			Title:             product.Title,
			Description:       product.Description,
			ShortDescription:  product.ShortDescription,
			Slug:              product.Slug,
			VariantLabel:      product.VariantLabel,
			VariantAttributes: product.VariantAttributes,
			Price:             product.Price,
			CompareAtPrice:    product.CompareAtPrice,
			CostPrice:         product.CostPrice,
			Currency:          product.Currency,
			StockQuantity:     product.StockQuantity,
			TrackInventory:    product.TrackInventory,
			AllowBackorder:    product.AllowBackorder,
			LowStockThreshold: product.LowStockThreshold,
			ImageURL:          product.ImageURL,
			Images:            product.Images,
			VideoURL:          product.VideoURL,
			Weight:            product.Weight,
			WeightUnit:        product.WeightUnit,
			Length:            product.Length,
			Width:             product.Width,
			Height:            product.Height,
			DimensionUnit:     product.DimensionUnit,
			Category:          product.Category,
			Tags:              product.Tags,
			Brand:             product.Brand,
			Status:            product.Status,
			IsActive:          product.IsActive,
			IsFeatured:        product.IsFeatured,
			Metadata:          product.Metadata,
		}
	}

	return responses
}
