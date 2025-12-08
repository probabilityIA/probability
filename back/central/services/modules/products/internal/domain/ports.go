package domain

import (
	"context"
)

// ───────────────────────────────────────────
//
//	REPOSITORY INTERFACE
//
// ───────────────────────────────────────────

// IRepository define todos los métodos de repositorio del módulo products
type IRepository interface {
	// CRUD Operations
	CreateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*Product, error)
	ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]Product, int64, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id string) error

	// Validation
	ProductExists(ctx context.Context, businessID uint, sku string) (bool, error)

	// Product-Integration Management
	AddProductIntegration(ctx context.Context, productID string, integrationID uint, externalProductID string) (*ProductBusinessIntegration, error)
	RemoveProductIntegration(ctx context.Context, productID string, integrationID uint) error
	GetProductIntegrations(ctx context.Context, productID string) ([]ProductBusinessIntegration, error)
	GetProductsByIntegration(ctx context.Context, integrationID uint) ([]Product, error)
	ProductIntegrationExists(ctx context.Context, productID string, integrationID uint) (bool, error)
}
