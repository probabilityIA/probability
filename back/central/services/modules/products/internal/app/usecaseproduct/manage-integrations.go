package usecaseproduct

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (uc *UseCaseProduct) AddProductIntegration(ctx context.Context, businessID uint, productID string, req *domain.AddProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	product, err := uc.repo.GetProductByID(ctx, businessID, productID)
	if err != nil {
		return nil, err
	}
	return uc.repo.AddProductIntegration(ctx, product.ID, req.IntegrationID, req.ExternalProductID, req.ExternalVariantID, req.ExternalSKU, req.ExternalBarcode)
}

func (uc *UseCaseProduct) UpdateProductIntegration(ctx context.Context, businessID uint, productID string, integrationID uint, req *domain.UpdateProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	if _, err := uc.repo.GetProductByID(ctx, businessID, productID); err != nil {
		return nil, err
	}
	return uc.repo.UpdateProductIntegration(ctx, productID, integrationID, req)
}

func (uc *UseCaseProduct) RemoveProductIntegration(ctx context.Context, businessID uint, productID string, integrationID uint) error {
	if _, err := uc.repo.GetProductByID(ctx, businessID, productID); err != nil {
		return err
	}
	return uc.repo.RemoveProductIntegration(ctx, productID, integrationID)
}

func (uc *UseCaseProduct) GetProductIntegrations(ctx context.Context, businessID uint, productID string) ([]domain.ProductBusinessIntegration, error) {
	if _, err := uc.repo.GetProductByID(ctx, businessID, productID); err != nil {
		return nil, err
	}
	return uc.repo.GetProductIntegrations(ctx, productID)
}

func (uc *UseCaseProduct) GetProductsByIntegration(ctx context.Context, integrationID uint) ([]domain.Product, error) {
	return uc.repo.GetProductsByIntegration(ctx, integrationID)
}

func (uc *UseCaseProduct) LookupProductByExternalRef(ctx context.Context, businessID uint, integrationID uint, externalVariantID, externalSKU, externalProductID, externalBarcode *string) (*domain.Product, error) {
	return uc.repo.LookupProductByExternalRef(ctx, businessID, integrationID, externalVariantID, externalSKU, externalProductID, externalBarcode)
}
