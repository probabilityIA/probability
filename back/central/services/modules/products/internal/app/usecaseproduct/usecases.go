package usecaseproduct

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// ───────────────────────────────────────────
//
//	CREATE PRODUCT
//
// ───────────────────────────────────────────

// CreateProduct crea un nuevo producto
func (uc *UseCaseProduct) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.ProductResponse, error) {
	// Validar que no exista un producto con el mismo SKU para el mismo negocio
	exists, err := uc.repo.ProductExists(ctx, req.BusinessID, req.SKU)
	if err != nil {
		return nil, fmt.Errorf("error checking if product exists: %w", err)
	}
	if exists {
		return nil, domain.ErrProductAlreadyExists
	}

	// Crear el modelo de producto con todos los campos
	product := &domain.Product{
		// Identificadores
		BusinessID: req.BusinessID,
		SKU:        req.SKU,
		ExternalID: req.ExternalID,

		// Información Básica
		Name:             req.Name,
		Title:            req.Title,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Slug:             req.Slug,

		// Pricing
		Price:          req.Price,
		CompareAtPrice: req.CompareAtPrice,
		CostPrice:      req.CostPrice,
		Currency:       req.Currency,

		// Inventory
		StockQuantity:     req.StockQuantity,
		TrackInventory:    req.TrackInventory,
		AllowBackorder:    req.AllowBackorder,
		LowStockThreshold: req.LowStockThreshold,

		// Media
		ImageURL: req.ImageURL,
		Images:   req.Images,
		VideoURL: req.VideoURL,

		// Dimensiones y Peso
		Weight:        req.Weight,
		WeightUnit:    req.WeightUnit,
		Length:        req.Length,
		Width:         req.Width,
		Height:        req.Height,
		DimensionUnit: req.DimensionUnit,

		// Categorización
		Category: req.Category,
		Tags:     req.Tags,
		Brand:    req.Brand,

		// Estado
		Status:     req.Status,
		IsActive:   req.IsActive,
		IsFeatured: req.IsFeatured,

		// Metadata
		Metadata: req.Metadata,
	}

	// Guardar en la base de datos
	if err := uc.repo.CreateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	// Retornar la respuesta
	return mapProductToResponse(product), nil
}

// ───────────────────────────────────────────
//
//	GET PRODUCT BY ID
//
// ───────────────────────────────────────────

// GetProductByID obtiene un producto por su ID
func (uc *UseCaseProduct) GetProductByID(ctx context.Context, id string) (*domain.ProductResponse, error) {
	if id == "" {
		return nil, errors.New("product ID is required")
	}

	product, err := uc.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	if product == nil {
		return nil, domain.ErrProductNotFound
	}

	return mapProductToResponse(product), nil
}

// ───────────────────────────────────────────
//
//	LIST PRODUCTS
//
// ───────────────────────────────────────────

// ListProducts obtiene una lista paginada de productos con filtros
func (uc *UseCaseProduct) ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.ProductsListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Obtener productos del repositorio
	products, total, err := uc.repo.ListProducts(ctx, page, pageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	// Mapear a respuestas
	productResponses := make([]domain.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = *mapProductToResponse(&product)
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &domain.ProductsListResponse{
		Data:       productResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ───────────────────────────────────────────
//
//	UPDATE PRODUCT
//
// ───────────────────────────────────────────

// UpdateProduct actualiza un producto existente
func (uc *UseCaseProduct) UpdateProduct(ctx context.Context, id string, req *domain.UpdateProductRequest) (*domain.ProductResponse, error) {
	if id == "" {
		return nil, errors.New("product ID is required")
	}

	// Obtener el producto existente
	product, err := uc.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	if product == nil {
		return nil, domain.ErrProductNotFound
	}

	// Actualizar solo los campos proporcionados
	// Identificadores
	if req.SKU != nil {
		// Si se cambia el SKU, verificar que no exista otro producto con ese SKU
		if *req.SKU != product.SKU {
			exists, err := uc.repo.ProductExists(ctx, product.BusinessID, *req.SKU)
			if err != nil {
				return nil, fmt.Errorf("error checking if product exists: %w", err)
			}
			if exists {
				return nil, domain.ErrProductAlreadyExists
			}
		}
		product.SKU = *req.SKU
	}
	if req.ExternalID != nil {
		product.ExternalID = *req.ExternalID
	}

	// Información Básica
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Title != nil {
		product.Title = *req.Title
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.ShortDescription != nil {
		product.ShortDescription = *req.ShortDescription
	}
	if req.Slug != nil {
		product.Slug = *req.Slug
	}

	// Pricing
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CompareAtPrice != nil {
		product.CompareAtPrice = req.CompareAtPrice
	}
	if req.CostPrice != nil {
		product.CostPrice = req.CostPrice
	}
	if req.Currency != nil {
		product.Currency = *req.Currency
	}

	// Inventory
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
	}
	if req.TrackInventory != nil {
		product.TrackInventory = *req.TrackInventory
	}
	if req.AllowBackorder != nil {
		product.AllowBackorder = *req.AllowBackorder
	}
	if req.LowStockThreshold != nil {
		product.LowStockThreshold = req.LowStockThreshold
	}

	// Media
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.VideoURL != nil {
		product.VideoURL = req.VideoURL
	}

	// Dimensiones y Peso
	if req.Weight != nil {
		product.Weight = req.Weight
	}
	if req.WeightUnit != nil {
		product.WeightUnit = *req.WeightUnit
	}
	if req.Length != nil {
		product.Length = req.Length
	}
	if req.Width != nil {
		product.Width = req.Width
	}
	if req.Height != nil {
		product.Height = req.Height
	}
	if req.DimensionUnit != nil {
		product.DimensionUnit = *req.DimensionUnit
	}

	// Categorización
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.Tags != nil {
		product.Tags = req.Tags
	}
	if req.Brand != nil {
		product.Brand = *req.Brand
	}

	// Estado
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}

	// Metadata
	if req.Metadata != nil {
		product.Metadata = req.Metadata
	}

	// Guardar cambios
	if err := uc.repo.UpdateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return mapProductToResponse(product), nil
}

// ───────────────────────────────────────────
//
//	DELETE PRODUCT
//
// ───────────────────────────────────────────

// DeleteProduct elimina un producto por su ID
func (uc *UseCaseProduct) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product ID is required")
	}

	// Validar que el producto existe
	_, err := uc.repo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	// Eliminar el producto
	if err := uc.repo.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}

// ───────────────────────────────────────────
//
//	HELPER FUNCTIONS
//
// ───────────────────────────────────────────

// mapProductToResponse convierte un modelo Product a ProductResponse
func mapProductToResponse(product *domain.Product) *domain.ProductResponse {
	return &domain.ProductResponse{
		// Timestamps
		ID:        product.ID,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		DeletedAt: product.DeletedAt,

		// Identificadores
		BusinessID: product.BusinessID,
		SKU:        product.SKU,
		ExternalID: product.ExternalID,

		// Información Básica
		Name:             product.Name,
		Title:            product.Title,
		Description:      product.Description,
		ShortDescription: product.ShortDescription,
		Slug:             product.Slug,

		// Pricing
		Price:          product.Price,
		CompareAtPrice: product.CompareAtPrice,
		CostPrice:      product.CostPrice,
		Currency:       product.Currency,

		// Inventory
		StockQuantity:     product.StockQuantity,
		TrackInventory:    product.TrackInventory,
		AllowBackorder:    product.AllowBackorder,
		LowStockThreshold: product.LowStockThreshold,

		// Media
		ImageURL: product.ImageURL,
		Images:   product.Images,
		VideoURL: product.VideoURL,

		// Dimensiones y Peso
		Weight:        product.Weight,
		WeightUnit:    product.WeightUnit,
		Length:        product.Length,
		Width:         product.Width,
		Height:        product.Height,
		DimensionUnit: product.DimensionUnit,

		// Categorización
		Category: product.Category,
		Tags:     product.Tags,
		Brand:    product.Brand,

		// Estado
		Status:     product.Status,
		IsActive:   product.IsActive,
		IsFeatured: product.IsFeatured,

		// Metadata
		Metadata: product.Metadata,
	}
}
