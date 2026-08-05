package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

type RepositoryMock struct {
	CreateProductFn               func(ctx context.Context, product *domain.Product) error
	GetProductByIDFn              func(ctx context.Context, businessID uint, id string) (*domain.Product, error)
	GetProductBySKUFn             func(ctx context.Context, businessID uint, sku string) (*domain.Product, error)
	ListProductsFn                func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error)
	ListSKUsFn                    func(ctx context.Context, businessID uint, prefix string) ([]string, error)
	GetNextSKUFn                  func(ctx context.Context, businessID uint, prefix string) (string, error)
	GetNextSKUBatchFn             func(ctx context.Context, businessID uint, prefix string, count int) ([]string, error)
	UpdateProductFn               func(ctx context.Context, product *domain.Product) error
	DeleteProductFn               func(ctx context.Context, businessID uint, id string) error
	ListProductsByFamilyIDFn      func(ctx context.Context, businessID uint, familyID uint) ([]domain.Product, error)
	ProductExistsFn               func(ctx context.Context, businessID uint, sku string) (bool, error)
	VariantExistsInFamilyFn       func(ctx context.Context, businessID uint, familyID uint, variantSignature string, excludeProductID *string) (bool, error)
	CreateProductFamilyFn         func(ctx context.Context, family *domain.ProductFamily) error
	GetProductFamilyByIDFn        func(ctx context.Context, businessID uint, familyID uint) (*domain.ProductFamily, error)
	ListProductFamiliesFn         func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.ProductFamily, int64, error)
	UpdateProductFamilyFn         func(ctx context.Context, family *domain.ProductFamily) error
	DeleteProductFamilyFn         func(ctx context.Context, businessID uint, familyID uint) error
	HasFamilyActiveVariantsFn     func(ctx context.Context, businessID uint, familyID uint) (bool, error)
	AddProductIntegrationFn       func(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error)
	UpdateProductIntegrationFn    func(ctx context.Context, productID string, integrationID uint, req *domain.UpdateProductIntegrationRequest) (*domain.ProductBusinessIntegration, error)
	RemoveProductIntegrationFn    func(ctx context.Context, productID string, integrationID uint) error
	GetProductIntegrationsFn      func(ctx context.Context, productID string) ([]domain.ProductBusinessIntegration, error)
	GetProductsByIntegrationFn    func(ctx context.Context, integrationID uint) ([]domain.Product, error)
	ProductIntegrationExistsFn    func(ctx context.Context, productID string, integrationID uint) (bool, error)
	LookupProductByExternalRefFn  func(ctx context.Context, businessID uint, integrationID uint, externalVariantID, externalSKU, externalProductID, externalBarcode *string) (*domain.Product, error)
}

var _ domain.IRepository = (*RepositoryMock)(nil)

func (m *RepositoryMock) CreateProduct(ctx context.Context, product *domain.Product) error {
	if m.CreateProductFn != nil {
		return m.CreateProductFn(ctx, product)
	}
	return nil
}

func (m *RepositoryMock) GetProductByID(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
	if m.GetProductByIDFn != nil {
		return m.GetProductByIDFn(ctx, businessID, id)
	}
	return nil, nil
}

func (m *RepositoryMock) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
	if m.GetProductBySKUFn != nil {
		return m.GetProductBySKUFn(ctx, businessID, sku)
	}
	return nil, nil
}

func (m *RepositoryMock) ListProducts(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error) {
	if m.ListProductsFn != nil {
		return m.ListProductsFn(ctx, businessID, page, pageSize, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) ListSKUs(ctx context.Context, businessID uint, prefix string) ([]string, error) {
	if m.ListSKUsFn != nil {
		return m.ListSKUsFn(ctx, businessID, prefix)
	}
	return nil, nil
}

func (m *RepositoryMock) GetNextSKU(ctx context.Context, businessID uint, prefix string) (string, error) {
	if m.GetNextSKUFn != nil {
		return m.GetNextSKUFn(ctx, businessID, prefix)
	}
	return "", nil
}

func (m *RepositoryMock) GetNextSKUBatch(ctx context.Context, businessID uint, prefix string, count int) ([]string, error) {
	if m.GetNextSKUBatchFn != nil {
		return m.GetNextSKUBatchFn(ctx, businessID, prefix, count)
	}
	return nil, nil
}

func (m *RepositoryMock) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if m.UpdateProductFn != nil {
		return m.UpdateProductFn(ctx, product)
	}
	return nil
}

func (m *RepositoryMock) DeleteProduct(ctx context.Context, businessID uint, id string) error {
	if m.DeleteProductFn != nil {
		return m.DeleteProductFn(ctx, businessID, id)
	}
	return nil
}

func (m *RepositoryMock) ListProductsByFamilyID(ctx context.Context, businessID uint, familyID uint) ([]domain.Product, error) {
	if m.ListProductsByFamilyIDFn != nil {
		return m.ListProductsByFamilyIDFn(ctx, businessID, familyID)
	}
	return nil, nil
}

func (m *RepositoryMock) ProductExists(ctx context.Context, businessID uint, sku string) (bool, error) {
	if m.ProductExistsFn != nil {
		return m.ProductExistsFn(ctx, businessID, sku)
	}
	return false, nil
}

func (m *RepositoryMock) VariantExistsInFamily(ctx context.Context, businessID uint, familyID uint, variantSignature string, excludeProductID *string) (bool, error) {
	if m.VariantExistsInFamilyFn != nil {
		return m.VariantExistsInFamilyFn(ctx, businessID, familyID, variantSignature, excludeProductID)
	}
	return false, nil
}

func (m *RepositoryMock) CreateProductFamily(ctx context.Context, family *domain.ProductFamily) error {
	if m.CreateProductFamilyFn != nil {
		return m.CreateProductFamilyFn(ctx, family)
	}
	return nil
}

func (m *RepositoryMock) GetProductFamilyByID(ctx context.Context, businessID uint, familyID uint) (*domain.ProductFamily, error) {
	if m.GetProductFamilyByIDFn != nil {
		return m.GetProductFamilyByIDFn(ctx, businessID, familyID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListProductFamilies(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.ProductFamily, int64, error) {
	if m.ListProductFamiliesFn != nil {
		return m.ListProductFamiliesFn(ctx, businessID, page, pageSize, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) UpdateProductFamily(ctx context.Context, family *domain.ProductFamily) error {
	if m.UpdateProductFamilyFn != nil {
		return m.UpdateProductFamilyFn(ctx, family)
	}
	return nil
}

func (m *RepositoryMock) DeleteProductFamily(ctx context.Context, businessID uint, familyID uint) error {
	if m.DeleteProductFamilyFn != nil {
		return m.DeleteProductFamilyFn(ctx, businessID, familyID)
	}
	return nil
}

func (m *RepositoryMock) HasFamilyActiveVariants(ctx context.Context, businessID uint, familyID uint) (bool, error) {
	if m.HasFamilyActiveVariantsFn != nil {
		return m.HasFamilyActiveVariantsFn(ctx, businessID, familyID)
	}
	return false, nil
}

func (m *RepositoryMock) AddProductIntegration(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error) {
	if m.AddProductIntegrationFn != nil {
		return m.AddProductIntegrationFn(ctx, productID, integrationID, externalProductID, externalVariantID, externalSKU, externalBarcode)
	}
	return nil, nil
}

func (m *RepositoryMock) UpdateProductIntegration(ctx context.Context, productID string, integrationID uint, req *domain.UpdateProductIntegrationRequest) (*domain.ProductBusinessIntegration, error) {
	if m.UpdateProductIntegrationFn != nil {
		return m.UpdateProductIntegrationFn(ctx, productID, integrationID, req)
	}
	return nil, nil
}

func (m *RepositoryMock) RemoveProductIntegration(ctx context.Context, productID string, integrationID uint) error {
	if m.RemoveProductIntegrationFn != nil {
		return m.RemoveProductIntegrationFn(ctx, productID, integrationID)
	}
	return nil
}

func (m *RepositoryMock) GetProductIntegrations(ctx context.Context, productID string) ([]domain.ProductBusinessIntegration, error) {
	if m.GetProductIntegrationsFn != nil {
		return m.GetProductIntegrationsFn(ctx, productID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetProductsByIntegration(ctx context.Context, integrationID uint) ([]domain.Product, error) {
	if m.GetProductsByIntegrationFn != nil {
		return m.GetProductsByIntegrationFn(ctx, integrationID)
	}
	return nil, nil
}

func (m *RepositoryMock) ProductIntegrationExists(ctx context.Context, productID string, integrationID uint) (bool, error) {
	if m.ProductIntegrationExistsFn != nil {
		return m.ProductIntegrationExistsFn(ctx, productID, integrationID)
	}
	return false, nil
}

func (m *RepositoryMock) LookupProductByExternalRef(ctx context.Context, businessID uint, integrationID uint, externalVariantID, externalSKU, externalProductID, externalBarcode *string) (*domain.Product, error) {
	if m.LookupProductByExternalRefFn != nil {
		return m.LookupProductByExternalRefFn(ctx, businessID, integrationID, externalVariantID, externalSKU, externalProductID, externalBarcode)
	}
	return nil, nil
}
