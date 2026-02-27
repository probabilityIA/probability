package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/app/usecaseproduct"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

// UseCases contiene todos los casos de uso del módulo products
type UseCases struct {
	repo domain.IRepository

	// Casos de uso modulares
	ProductCRUD *usecaseproduct.UseCaseProduct
}

// New crea una nueva instancia de UseCases
func New(repo domain.IRepository) *UseCases {
	return &UseCases{
		repo:        repo,
		ProductCRUD: usecaseproduct.New(repo),
	}
}

// ───────────────────────────────────────────
// MÉTODOS DE COMPATIBILIDAD - Delegar al CRUD
// ───────────────────────────────────────────

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

// ───────────────────────────────────────────
// MÉTODOS DE INTEGRACIÓN - Delegar al CRUD
// ───────────────────────────────────────────

// AddProductIntegration delega al caso de uso CRUD
func (uc *UseCases) AddProductIntegration(ctx context.Context, businessID uint, productID string, req *domain.AddProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	return uc.ProductCRUD.AddProductIntegration(ctx, businessID, productID, req)
}

// RemoveProductIntegration delega al caso de uso CRUD
func (uc *UseCases) RemoveProductIntegration(ctx context.Context, businessID uint, productID string, integrationID uint) error {
	return uc.ProductCRUD.RemoveProductIntegration(ctx, businessID, productID, integrationID)
}

// GetProductIntegrations delega al caso de uso CRUD
func (uc *UseCases) GetProductIntegrations(ctx context.Context, businessID uint, productID string) ([]domain.ProductBusinessIntegration, error) {
	return uc.ProductCRUD.GetProductIntegrations(ctx, businessID, productID)
}

// GetProductsByIntegration delega al caso de uso CRUD
func (uc *UseCases) GetProductsByIntegration(ctx context.Context, integrationID uint) ([]domain.Product, error) {
	return uc.ProductCRUD.GetProductsByIntegration(ctx, integrationID)
}
