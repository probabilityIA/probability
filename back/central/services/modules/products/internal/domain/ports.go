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
	GetProductByID(ctx context.Context, id uint) (*Product, error)
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*Product, error)
	ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]Product, int64, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id uint) error

	// Validation
	ProductExists(ctx context.Context, businessID uint, sku string) (bool, error)
}

