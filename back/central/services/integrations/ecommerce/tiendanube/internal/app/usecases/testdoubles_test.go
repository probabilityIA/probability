package usecases

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

const (
	testBusinessID    = uint(7)
	testIntegrationID = "77"
)

type fakeService struct {
	integration *domain.Integration
}

func (f *fakeService) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	return f.integration, nil
}

func (f *fakeService) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return "token-" + fieldName, nil
}

func (f *fakeService) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	return nil
}

type fakeClient struct {
	domain.ITiendanubeClient
	products     []domain.TiendanubeProduct
	channelStock []domain.ChannelStock
	stockErr     error
	stockWrites  []string
	created      []domain.CreateProductInput
}

func (f *fakeClient) GetProducts(ctx context.Context, cred domain.Credential) ([]domain.TiendanubeProduct, error) {
	return f.products, nil
}

func (f *fakeClient) GetProductsStock(ctx context.Context, cred domain.Credential, externalIDs []string) ([]domain.ChannelStock, error) {
	if f.stockErr != nil {
		return nil, f.stockErr
	}
	return f.channelStock, nil
}

func (f *fakeClient) CreateProduct(ctx context.Context, cred domain.Credential, input domain.CreateProductInput) (int64, int64, error) {
	f.created = append(f.created, input)
	id := int64(1000 + len(f.created))
	return id, id + 5000, nil
}

type fakeRepo struct {
	probability []domain.ProductForSync
	mapped      []domain.MappedItem
	stock       map[string]int
	mappings    []productmatch.ExternalRefs
	savedRows   []inventorycompare.Row
}

func (f *fakeRepo) ListProductsByBusiness(ctx context.Context, businessID uint) ([]domain.ProductForSync, error) {
	return f.probability, nil
}

func (f *fakeRepo) GetExternalProductID(ctx context.Context, productID string, integrationID uint) (string, bool, error) {
	return "", false, nil
}

func (f *fakeRepo) UpsertProductIntegrationMapping(ctx context.Context, productID string, businessID, integrationID uint, refs productmatch.ExternalRefs) error {
	f.mappings = append(f.mappings, refs)
	return nil
}

func (f *fakeRepo) ListMappedItems(ctx context.Context, integrationID uint) ([]domain.MappedItem, error) {
	return f.mapped, nil
}

func (f *fakeRepo) GetStockForProducts(ctx context.Context, productIDs []string, warehouseIDs []uint) (map[string]int, error) {
	if f.stock == nil {
		return map[string]int{}, nil
	}
	return f.stock, nil
}

func (f *fakeRepo) SaveCompareSnapshot(ctx context.Context, businessID, integrationID uint, rows []inventorycompare.Row, checkedAt time.Time) error {
	f.savedRows = rows
	return nil
}

func (f *fakeRepo) LoadCompareSnapshot(ctx context.Context, businessID, integrationID uint, opts inventorycompare.LoadOptions) (*inventorycompare.Page, error) {
	return &inventorycompare.Page{Rows: f.savedRows, FromCache: true}, nil
}

func newTestUseCase(client domain.ITiendanubeClient, repo domain.IProductRepository) *tiendanubeUseCase {
	businessID := testBusinessID
	return &tiendanubeUseCase{
		client: client,
		service: &fakeService{
			integration: &domain.Integration{
				ID:         77,
				BusinessID: &businessID,
				Name:       "Tiendanube Test",
				StoreID:    "1234567",
				Config:     map[string]interface{}{"store_id": "1234567"},
				BaseURL:    "https://api.tiendanube.com/v1",
			},
		},
		productRepo: repo,
		logger:      log.New(),
	}
}

func productoCanal(id, variantID int64, sku, nombre string, precio float64, stock int, manejaStock bool) domain.TiendanubeProduct {
	return domain.TiendanubeProduct{
		ID:   id,
		Name: nombre,
		Variants: []domain.TiendanubeVariant{{
			ID:              variantID,
			ProductID:       id,
			SKU:             sku,
			Price:           precio,
			Stock:           stock,
			StockManagement: manejaStock,
		}},
	}
}

var errNoDebioLlamarse = errTiendanubeTest("no se debio llamar al canal")

type errTiendanubeTest string

func (e errTiendanubeTest) Error() string { return string(e) }
