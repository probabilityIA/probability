package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

type fakeProductRepo struct {
	domain.IProductRepository
	GetProductIDByExternalRefFn       func(ctx context.Context, integrationID uint, externalProductID, externalVariantID string) (string, error)
	GetProductSKUByIDFn               func(ctx context.Context, productID string, businessID uint) (string, error)
	GetExternalRefsFn                 func(ctx context.Context, productID string, integrationID uint) (string, string, error)
	ListProductsByBusinessFn          func(ctx context.Context, businessID uint) ([]domain.ProductForSync, error)
	UpsertProductIntegrationMappingFn func(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error
}

func (f *fakeProductRepo) ListProductsByBusiness(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
	if f.ListProductsByBusinessFn != nil {
		return f.ListProductsByBusinessFn(ctx, businessID)
	}
	return nil, nil
}

func (f *fakeProductRepo) ListMappedItems(ctx context.Context, integrationID uint) ([]domain.MappedItem, error) {
	return nil, nil
}

func (f *fakeProductRepo) UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
	if f.UpsertProductIntegrationMappingFn != nil {
		return f.UpsertProductIntegrationMappingFn(ctx, productID, businessID, integrationID, refs)
	}
	return nil
}

func (f *fakeProductRepo) GetProductIDByExternalRef(ctx context.Context, integrationID uint, externalProductID, externalVariantID string) (string, error) {
	if f.GetProductIDByExternalRefFn != nil {
		return f.GetProductIDByExternalRefFn(ctx, integrationID, externalProductID, externalVariantID)
	}
	return "", nil
}

func (f *fakeProductRepo) GetProductSKUByID(ctx context.Context, productID string, businessID uint) (string, error) {
	if f.GetProductSKUByIDFn != nil {
		return f.GetProductSKUByIDFn(ctx, productID, businessID)
	}
	return "", nil
}

func (f *fakeProductRepo) GetExternalRefs(ctx context.Context, productID string, integrationID uint) (string, string, error) {
	if f.GetExternalRefsFn != nil {
		return f.GetExternalRefsFn(ctx, productID, integrationID)
	}
	return "", "", nil
}

func businessIDPtr(v uint) *uint { return &v }

func TestUpdateInventory_NotFound_SinMapeoEnProbability_DevuelveErrorOriginal(t *testing.T) {
	businessID := businessIDPtr(46)
	svc := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(ctx context.Context, integrationID string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         221,
				BusinessID: businessID,
				Config: map[string]interface{}{
					"store_url":              "https://vigaropadeportiva.com",
					"inventory_sync_enabled": true,
				},
			}, nil
		},
	}
	client := &mocks.WooClientMock{
		UpdateProductStockFn: func(ctx context.Context, storeURL, consumerKey, consumerSecret, productExternalID string, quantity int) error {
			return domain.ErrProductNotFoundInStore
		},
	}
	repo := &fakeProductRepo{
		GetProductIDByExternalRefFn: func(ctx context.Context, integrationID uint, externalProductID, externalVariantID string) (string, error) {
			return "", nil
		},
	}

	uc := New(client, svc, nil, repo, nil, log.New())
	err := uc.UpdateInventory(context.Background(), "221", "10859:10863", 11)

	if err == nil {
		t.Fatal("esperaba error cuando no se puede autocorregir el mapeo, obtuve nil")
	}
	if !errors.Is(err, domain.ErrProductNotFoundInStore) {
		t.Errorf("error = %v, esperaba que envolviera ErrProductNotFoundInStore", err)
	}
}

func TestUpdateInventory_NotFound_ReasociaYReintentaConElIDNuevo(t *testing.T) {
	businessID := businessIDPtr(46)
	svc := &mocks.IntegrationServiceMock{
		GetIntegrationByIDFn: func(ctx context.Context, integrationID string) (*domain.Integration, error) {
			return &domain.Integration{
				ID:         221,
				BusinessID: businessID,
				Config: map[string]interface{}{
					"store_url":              "https://vigaropadeportiva.com",
					"inventory_sync_enabled": true,
				},
			}, nil
		},
	}

	calls := 0
	client := &mocks.WooClientMock{
		UpdateProductStockFn: func(ctx context.Context, storeURL, consumerKey, consumerSecret, productExternalID string, quantity int) error {
			calls++
			if productExternalID == "10859:10863" {
				return domain.ErrProductNotFoundInStore
			}
			if productExternalID == "10859:14710" {
				return nil
			}
			t.Fatalf("ID inesperado en UpdateProductStock: %s", productExternalID)
			return nil
		},
		GetProductsFn: func(ctx context.Context, storeURL, consumerKey, consumerSecret string) ([]domain.WooProduct, error) {
			return []domain.WooProduct{
				{ID: "14710", ParentID: "10859", SKU: "SD313-4XL"},
			}, nil
		},
	}

	repo := &fakeProductRepo{
		ListProductsByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
			return []domain.ProductForSync{
				{ID: "PRD_test", SKU: "SD313-4XL", TrackInventory: true},
			}, nil
		},
		GetProductIDByExternalRefFn: func(ctx context.Context, integrationID uint, externalProductID, externalVariantID string) (string, error) {
			if externalProductID == "10859" && externalVariantID == "10863" {
				return "PRD_test", nil
			}
			return "", nil
		},
		GetProductSKUByIDFn: func(ctx context.Context, productID string, businessID uint) (string, error) {
			if productID == "PRD_test" {
				return "SD313-4XL", nil
			}
			return "", nil
		},
		GetExternalRefsFn: func(ctx context.Context, productID string, integrationID uint) (string, string, error) {
			if productID == "PRD_test" {
				return "10859", "14710", nil
			}
			return "", "", nil
		},
	}

	uc := &wooCommerceUseCase{
		client:      client,
		service:     svc,
		productRepo: repo,
		logger:      log.New(),
	}

	err := uc.UpdateInventory(context.Background(), "221", "10859:10863", 11)
	if err != nil {
		t.Fatalf("esperaba nil tras autocorregir el mapeo y reintentar, obtuve: %v", err)
	}
	if calls != 2 {
		t.Errorf("UpdateProductStock se llamo %d veces, esperaba 2 (original + reintento)", calls)
	}
}
