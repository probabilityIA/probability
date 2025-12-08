package usecaseproduct

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// AddProductIntegration asocia un producto con una integración
func (uc *UseCaseProduct) AddProductIntegration(ctx context.Context, productID string, req *domain.AddProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	// Validar que el producto existe
	product, err := uc.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Agregar la integración (el repositorio valida que pertenezca al mismo negocio)
	integration, err := uc.repo.AddProductIntegration(
		ctx,
		product.ID,
		req.IntegrationID,
		req.ExternalProductID,
	)
	if err != nil {
		return nil, err
	}

	return integration, nil
}

// RemoveProductIntegration remueve la asociación entre un producto y una integración
func (uc *UseCaseProduct) RemoveProductIntegration(ctx context.Context, productID string, integrationID uint) error {
	// Validar que el producto existe
	_, err := uc.repo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	// Remover la integración
	return uc.repo.RemoveProductIntegration(ctx, productID, integrationID)
}

// GetProductIntegrations obtiene todas las integraciones asociadas a un producto
func (uc *UseCaseProduct) GetProductIntegrations(ctx context.Context, productID string) ([]domain.ProductBusinessIntegration, error) {
	// Validar que el producto existe
	_, err := uc.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Obtener las integraciones
	return uc.repo.GetProductIntegrations(ctx, productID)
}

// GetProductsByIntegration obtiene todos los productos asociados a una integración
func (uc *UseCaseProduct) GetProductsByIntegration(ctx context.Context, integrationID uint) ([]domain.Product, error) {
	return uc.repo.GetProductsByIntegration(ctx, integrationID)
}
