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

	// Crear el modelo de producto
	product := &domain.Product{
		BusinessID: req.BusinessID,
		SKU:        req.SKU,
		Name:       req.Name,
		ExternalID: req.ExternalID,
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
func (uc *UseCaseProduct) GetProductByID(ctx context.Context, id uint) (*domain.ProductResponse, error) {
	if id == 0 {
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
func (uc *UseCaseProduct) UpdateProduct(ctx context.Context, id uint, req *domain.UpdateProductRequest) (*domain.ProductResponse, error) {
	if id == 0 {
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
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.ExternalID != nil {
		product.ExternalID = *req.ExternalID
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

// DeleteProduct elimina (soft delete) un producto
func (uc *UseCaseProduct) DeleteProduct(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("product ID is required")
	}

	// Verificar que el producto existe
	product, err := uc.repo.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting product: %w", err)
	}

	if product == nil {
		return domain.ErrProductNotFound
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
		ID:         product.ID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
		DeletedAt:  product.DeletedAt,
		BusinessID: product.BusinessID,
		SKU:        product.SKU,
		Name:       product.Name,
		ExternalID: product.ExternalID,
	}
}

