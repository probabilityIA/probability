package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/app/usecasefamily"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/app/usecaseproduct"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

type UseCases struct {
	repo    domain.IRepository
	eventos *publicadorEventos

	ProductCRUD *usecaseproduct.UseCaseProduct
	FamilyCRUD  *usecasefamily.UseCaseFamily
}

func New(repo domain.IRepository) *UseCases {
	return &UseCases{
		repo:        repo,
		ProductCRUD: usecaseproduct.New(repo),
		FamilyCRUD:  usecasefamily.New(repo),
	}
}

func (uc *UseCases) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.CreateProduct(ctx, req)
}

func (uc *UseCases) GetProductByID(ctx context.Context, businessID uint, id string) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.GetProductByID(ctx, businessID, id)
}

func (uc *UseCases) ListProducts(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) (*domain.ProductsListResponse, error) {
	return uc.ProductCRUD.ListProducts(ctx, businessID, page, pageSize, filters)
}

func (uc *UseCases) UpdateProduct(ctx context.Context, businessID uint, id string, req *domain.UpdateProductRequest) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.UpdateProduct(ctx, businessID, id, req)
}

func (uc *UseCases) DeleteProduct(ctx context.Context, businessID uint, id string) error {
	return uc.ProductCRUD.DeleteProduct(ctx, businessID, id)
}

func (uc *UseCases) AddProductIntegration(ctx context.Context, businessID uint, productID string, req *domain.AddProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	return uc.ProductCRUD.AddProductIntegration(ctx, businessID, productID, req)
}

func (uc *UseCases) UpdateProductIntegration(ctx context.Context, businessID uint, productID string, integrationID uint, req *domain.UpdateProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	return uc.ProductCRUD.UpdateProductIntegration(ctx, businessID, productID, integrationID, req)
}

func (uc *UseCases) RemoveProductIntegration(ctx context.Context, businessID uint, productID string, integrationID uint) error {
	return uc.ProductCRUD.RemoveProductIntegration(ctx, businessID, productID, integrationID)
}

func (uc *UseCases) GetProductIntegrations(ctx context.Context, businessID uint, productID string) ([]domain.ProductBusinessIntegration, error) {
	return uc.ProductCRUD.GetProductIntegrations(ctx, businessID, productID)
}

func (uc *UseCases) GetProductsByIntegration(ctx context.Context, integrationID uint) ([]domain.Product, error) {
	return uc.ProductCRUD.GetProductsByIntegration(ctx, integrationID)
}

func (uc *UseCases) LookupProductByExternalRef(ctx context.Context, businessID uint, integrationID uint, externalVariantID, externalSKU, externalProductID, externalBarcode *string) (*domain.Product, error) {
	return uc.ProductCRUD.LookupProductByExternalRef(ctx, businessID, integrationID, externalVariantID, externalSKU, externalProductID, externalBarcode)
}

func (uc *UseCases) CreateProductFamily(ctx context.Context, req *domain.CreateProductFamilyStandaloneRequest) (*domain.ProductFamilyResponse, error) {
	return uc.FamilyCRUD.CreateProductFamily(ctx, req)
}

func (uc *UseCases) GetProductFamilyByID(ctx context.Context, businessID uint, familyID uint) (*domain.ProductFamilyResponse, error) {
	return uc.FamilyCRUD.GetProductFamilyByID(ctx, businessID, familyID)
}

func (uc *UseCases) ListProductFamilies(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) (*domain.ProductFamiliesListResponse, error) {
	return uc.FamilyCRUD.ListProductFamilies(ctx, businessID, page, pageSize, filters)
}

func (uc *UseCases) UpdateProductFamily(ctx context.Context, businessID uint, familyID uint, req *domain.UpdateProductFamilyRequest) (*domain.ProductFamilyResponse, error) {
	return uc.FamilyCRUD.UpdateProductFamily(ctx, businessID, familyID, req)
}

func (uc *UseCases) DeleteProductFamily(ctx context.Context, businessID uint, familyID uint) error {
	return uc.FamilyCRUD.DeleteProductFamily(ctx, businessID, familyID)
}

func (uc *UseCases) ListProductsByFamilyID(ctx context.Context, businessID uint, familyID uint) ([]domain.Product, error) {
	return uc.repo.ListProductsByFamilyID(ctx, businessID, familyID)
}

func (uc *UseCases) ListSKUs(ctx context.Context, businessID uint, prefix string) ([]string, error) {
	return uc.ProductCRUD.ListSKUs(ctx, businessID, prefix)
}

func (uc *UseCases) GetNextSKU(ctx context.Context, businessID uint, prefix string) (string, error) {
	return uc.ProductCRUD.GetNextSKU(ctx, businessID, prefix)
}

func (uc *UseCases) GetNextSKUBatch(ctx context.Context, businessID uint, prefix string, count int) ([]string, error) {
	return uc.ProductCRUD.GetNextSKUBatch(ctx, businessID, prefix, count)
}

func (uc *UseCases) ExportDimensions(ctx context.Context, businessID uint) ([]domain.Product, error) {
	return uc.ProductCRUD.ExportDimensions(ctx, businessID)
}

func (uc *UseCases) BulkUpdateDimensions(ctx context.Context, businessID uint, rows []domain.DimensionRow) (*domain.BulkDimensionsResult, error) {
	return uc.ProductCRUD.BulkUpdateDimensions(ctx, businessID, rows)
}

func (uc *UseCases) MatchFamilyImportRows(ctx context.Context, businessID uint, rows []domain.DimensionRow) (*domain.FamilyImportPreview, error) {
	return uc.ProductCRUD.MatchFamilyImportRows(ctx, businessID, rows)
}
