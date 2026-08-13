package domain

import (
	"context"
)

type IRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, businessID uint, id string) (*Product, error)
	GetProductBySKU(ctx context.Context, businessID uint, sku string) (*Product, error)
	ListProducts(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]Product, int64, error)
	ListSKUs(ctx context.Context, businessID uint, prefix string) ([]string, error)
	GetNextSKU(ctx context.Context, businessID uint, prefix string) (string, error)
	GetNextSKUBatch(ctx context.Context, businessID uint, prefix string, count int) ([]string, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, businessID uint, id string) error
	ListProductsByFamilyID(ctx context.Context, businessID uint, familyID uint) ([]Product, error)

	IDataApplyRepository

	RecordFieldChanges(ctx context.Context, changes []FieldChange) error
	GetProductFieldOrigins(ctx context.Context, businessID uint, productID string) ([]FieldOriginView, error)

	GetProductChannels(ctx context.Context, businessID uint, productID string) ([]ProductChannel, error)
	GetProductStock(ctx context.Context, businessID uint, productID string) ([]ProductStockLevel, error)
	GetProductMovements(ctx context.Context, businessID uint, productID string, limit int) ([]ProductMovement, error)
	CountFamilySiblings(ctx context.Context, businessID uint, familyID uint, productID string) (int64, error)

	ProductExists(ctx context.Context, businessID uint, sku string) (bool, error)
	VariantExistsInFamily(ctx context.Context, businessID uint, familyID uint, variantSignature string, excludeProductID *string) (bool, error)

	CreateProductFamily(ctx context.Context, family *ProductFamily) error
	GetProductFamilyByID(ctx context.Context, businessID uint, familyID uint) (*ProductFamily, error)
	GetProductFamilyByName(ctx context.Context, businessID uint, name string) (*ProductFamily, error)
	ListProductFamilies(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]ProductFamily, int64, error)
	UpdateProductFamily(ctx context.Context, family *ProductFamily) error
	DeleteProductFamily(ctx context.Context, businessID uint, familyID uint) error
	HasFamilyActiveVariants(ctx context.Context, businessID uint, familyID uint) (bool, error)

	AddProductIntegration(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*ProductBusinessIntegration, error)
	UpdateProductIntegration(ctx context.Context, productID string, integrationID uint, req *UpdateProductIntegrationRequest) (*ProductBusinessIntegration, error)
	RemoveProductIntegration(ctx context.Context, productID string, integrationID uint) error
	GetProductIntegrations(ctx context.Context, productID string) ([]ProductBusinessIntegration, error)
	GetProductsByIntegration(ctx context.Context, integrationID uint) ([]Product, error)
	ProductIntegrationExists(ctx context.Context, productID string, integrationID uint) (bool, error)
	LookupProductByExternalRef(ctx context.Context, businessID uint, integrationID uint, externalVariantID, externalSKU, externalProductID, externalBarcode *string) (*Product, error)
}
