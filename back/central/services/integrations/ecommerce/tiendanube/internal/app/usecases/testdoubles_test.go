package usecases

import (
	"context"
	"encoding/json"
	"strconv"
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
	tokenVacio  bool
}

func (f *fakeService) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	return f.integration, nil
}

func (f *fakeService) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	if f.tokenVacio {
		return "", nil
	}
	return "token-" + fieldName, nil
}

func (f *fakeService) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	return nil
}

type fakeClient struct {
	domain.ITiendanubeClient
	products      []domain.TiendanubeProduct
	channelStock  []domain.ChannelStock
	stockErr      error
	created       []domain.CreateProductInput
	actualizados  []domain.UpdateProductInput
	variantes     []domain.UpdateVariantInput
	stockEscrito  map[string]int
	objetivo      *domain.StockTarget
	listaErr      error
	fallaCrearSKU string
	tienda        *domain.StoreInfo
	tiendaErr     error
}

func (f *fakeClient) GetStoreInfo(ctx context.Context, cred domain.Credential) (*domain.StoreInfo, error) {
	if f.tiendaErr != nil {
		return nil, f.tiendaErr
	}
	if f.tienda != nil {
		return f.tienda, nil
	}
	return &domain.StoreInfo{Name: "Tienda de prueba"}, nil
}

func (f *fakeClient) GetProducts(ctx context.Context, cred domain.Credential) ([]domain.TiendanubeProduct, error) {
	if f.listaErr != nil {
		return nil, f.listaErr
	}
	return f.products, nil
}

func (f *fakeClient) GetProductsStock(ctx context.Context, cred domain.Credential, externalIDs []string) ([]domain.ChannelStock, error) {
	if f.stockErr != nil {
		return nil, f.stockErr
	}
	return f.channelStock, nil
}

func (f *fakeClient) UpdateProduct(ctx context.Context, cred domain.Credential, productID int64, input domain.UpdateProductInput) error {
	f.actualizados = append(f.actualizados, input)
	return nil
}

func (f *fakeClient) UpdateVariant(ctx context.Context, cred domain.Credential, productID, variantID int64, input domain.UpdateVariantInput) error {
	f.variantes = append(f.variantes, input)
	return nil
}

func (f *fakeClient) SetVariantStock(ctx context.Context, cred domain.Credential, productID, variantID int64, stock int) error {
	if f.stockEscrito == nil {
		f.stockEscrito = make(map[string]int)
	}
	f.stockEscrito[strconv.FormatInt(productID, 10)+":"+strconv.FormatInt(variantID, 10)] = stock
	return nil
}

func (f *fakeClient) ResolveStockTarget(ctx context.Context, cred domain.Credential, sku string) (*domain.StockTarget, error) {
	if f.objetivo != nil {
		return f.objetivo, nil
	}
	return &domain.StockTarget{Found: false}, nil
}

func (f *fakeClient) CreateProduct(ctx context.Context, cred domain.Credential, input domain.CreateProductInput) (int64, int64, error) {
	if f.fallaCrearSKU != "" && f.fallaCrearSKU == input.SKU {
		return 0, 0, errTiendanubeTest("el canal rechazo el producto")
	}
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

type colaFalsa struct {
	publicados map[string][][]byte
	fallaAl    string
}

func nuevaCola() *colaFalsa {
	return &colaFalsa{publicados: make(map[string][][]byte)}
}

func (c *colaFalsa) Publish(ctx context.Context, queueName string, message []byte) error {
	if c.fallaAl != "" && c.fallaAl == queueName {
		return errTiendanubeTest("cola caida")
	}
	c.publicados[queueName] = append(c.publicados[queueName], message)
	return nil
}

func (c *colaFalsa) PublishToExchange(ctx context.Context, exchangeName, routingKey string, message []byte) error {
	c.publicados[exchangeName] = append(c.publicados[exchangeName], message)
	return nil
}

func (c *colaFalsa) Consume(ctx context.Context, queueName string, handler func([]byte) error) error {
	return nil
}

func (c *colaFalsa) ConsumeConcurrent(ctx context.Context, queueName string, handler func([]byte) error, workers int) error {
	return nil
}

func (c *colaFalsa) Close() error { return nil }

func (c *colaFalsa) DeclareQueue(queueName string, durable bool) error { return nil }

func (c *colaFalsa) DeclareExchange(exchangeName, exchangeType string, durable bool) error {
	return nil
}

func (c *colaFalsa) BindQueue(queueName, exchangeName, routingKey string) error { return nil }

func (c *colaFalsa) Ping() error { return nil }

func (c *colaFalsa) mensajes(cola string) []providerUpsertMsg {
	out := make([]providerUpsertMsg, 0)
	for _, raw := range c.publicados[cola] {
		var m providerUpsertMsg
		if err := json.Unmarshal(raw, &m); err == nil {
			out = append(out, m)
		}
	}
	return out
}

func newTestUseCaseConCola(client domain.ITiendanubeClient, repo domain.IProductRepository, cola *colaFalsa) *tiendanubeUseCase {
	uc := newTestUseCase(client, repo)
	uc.rabbit = cola
	return uc
}
