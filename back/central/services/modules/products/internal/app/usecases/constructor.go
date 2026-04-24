package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/app/usecasefamily"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/app/usecaseproduct"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// UseCases contiene todos los casos de uso del módulo products
type UseCases struct {
	repo domain.IRepository

	// Casos de uso modulares
	ProductCRUD *usecaseproduct.UseCaseProduct
	FamilyCRUD  *usecasefamily.UseCaseFamily
}

// New crea una nueva instancia de UseCases
func New(repo domain.IRepository) *UseCases {
	return &UseCases{
		repo:        repo,
		ProductCRUD: usecaseproduct.New(repo),
		FamilyCRUD:  usecasefamily.New(repo),
	}
}

// MÉTODOS DE COMPATIBILIDAD - Delegar al CRUD

// CreateProduct delega al caso de uso CRUD
func (uc *UseCases) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.CreateProduct(ctx, req)
}

// GetProductByID delega al caso de uso CRUD
func (uc *UseCases) GetProductByID(ctx context.Context, businessID uint, id string) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.GetProductByID(ctx, businessID, id)
}

// ListProducts delega al caso de uso CRUD
func (uc *UseCases) ListProducts(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) (*domain.ProductsListResponse, error) {
	return uc.ProductCRUD.ListProducts(ctx, businessID, page, pageSize, filters)
}

// UpdateProduct delega al caso de uso CRUD
func (uc *UseCases) UpdateProduct(ctx context.Context, businessID uint, id string, req *domain.UpdateProductRequest) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.UpdateProduct(ctx, businessID, id, req)
}

// DeleteProduct delega al caso de uso CRUD
func (uc *UseCases) DeleteProduct(ctx context.Context, businessID uint, id string) error {
	return uc.ProductCRUD.DeleteProduct(ctx, businessID, id)
}

// MÉTODOS DE INTEGRACIÓN - Delegar al CRUD

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

// CreateProductFamily delega al caso de uso CRUD de familias.
func (uc *UseCases) CreateProductFamily(ctx context.Context, req *domain.CreateProductFamilyStandaloneRequest) (*domain.ProductFamilyResponse, error) {
	return uc.FamilyCRUD.CreateProductFamily(ctx, req)
}

// GetProductFamilyByID delega al caso de uso CRUD de familias.
func (uc *UseCases) GetProductFamilyByID(ctx context.Context, businessID uint, familyID uint) (*domain.ProductFamilyResponse, error) {
	return uc.FamilyCRUD.GetProductFamilyByID(ctx, businessID, familyID)
}

// ListProductFamilies delega al caso de uso CRUD de familias.
func (uc *UseCases) ListProductFamilies(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) (*domain.ProductFamiliesListResponse, error) {
	return uc.FamilyCRUD.ListProductFamilies(ctx, businessID, page, pageSize, filters)
}

// UpdateProductFamily delega al caso de uso CRUD de familias.
func (uc *UseCases) UpdateProductFamily(ctx context.Context, businessID uint, familyID uint, req *domain.UpdateProductFamilyRequest) (*domain.ProductFamilyResponse, error) {
	return uc.FamilyCRUD.UpdateProductFamily(ctx, businessID, familyID, req)
}

// DeleteProductFamily delega al caso de uso CRUD de familias.
func (uc *UseCases) DeleteProductFamily(ctx context.Context, businessID uint, familyID uint) error {
	return uc.FamilyCRUD.DeleteProductFamily(ctx, businessID, familyID)
}

// ListProductsByFamilyID delega al repositorio para listar variantes de una familia.
func (uc *UseCases) ListProductsByFamilyID(ctx context.Context, businessID uint, familyID uint) ([]domain.Product, error) {
	return uc.repo.ListProductsByFamilyID(ctx, businessID, familyID)
}
