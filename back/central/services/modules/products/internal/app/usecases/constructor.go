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
func (uc *UseCases) GetProductByID(ctx context.Context, id string) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.GetProductByID(ctx, id)
}

// ListProducts delega al caso de uso CRUD
func (uc *UseCases) ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.ProductsListResponse, error) {
	return uc.ProductCRUD.ListProducts(ctx, page, pageSize, filters)
}

// UpdateProduct delega al caso de uso CRUD
func (uc *UseCases) UpdateProduct(ctx context.Context, id string, req *domain.UpdateProductRequest) (*domain.ProductResponse, error) {
	return uc.ProductCRUD.UpdateProduct(ctx, id, req)
}

// DeleteProduct delega al caso de uso CRUD
func (uc *UseCases) DeleteProduct(ctx context.Context, id string) error {
	return uc.ProductCRUD.DeleteProduct(ctx, id)
}

// ───────────────────────────────────────────
// MÉTODOS DE INTEGRACIÓN - Delegar al CRUD
// ───────────────────────────────────────────

// AddProductIntegration delega al caso de uso CRUD
func (uc *UseCases) AddProductIntegration(ctx context.Context, productID string, req *domain.AddProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	return uc.ProductCRUD.AddProductIntegration(ctx, productID, req)
}

// RemoveProductIntegration delega al caso de uso CRUD
func (uc *UseCases) RemoveProductIntegration(ctx context.Context, productID string, integrationID uint) error {
	return uc.ProductCRUD.RemoveProductIntegration(ctx, productID, integrationID)
}

// GetProductIntegrations delega al caso de uso CRUD
func (uc *UseCases) GetProductIntegrations(ctx context.Context, productID string) ([]domain.ProductBusinessIntegration, error) {
	return uc.ProductCRUD.GetProductIntegrations(ctx, productID)
}

// GetProductsByIntegration delega al caso de uso CRUD
func (uc *UseCases) GetProductsByIntegration(ctx context.Context, integrationID uint) ([]domain.Product, error) {
	return uc.ProductCRUD.GetProductsByIntegration(ctx, integrationID)
}
