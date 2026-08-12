package usecases

import (
	"context"
	"sync"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

func testLogger() log.ILogger { return log.New() }

// futureToken devuelve una config con el token vigente, para que EnsureValidToken
// tome la rama que desencripta en vez de la que refresca contra MercadoLibre.
func futureToken() map[string]interface{} {
	return map[string]interface{}{
		"token_expires_at":       time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		"inventory_sync_enabled": true,
	}
}

type stockCall struct {
	ItemID    string
	VariantID string
	Quantity  int
}

type mockClient struct {
	domain.IMeliClient

	item    *domain.MeliItemDetail
	getErr  error
	putErrs []error

	mu        sync.Mutex
	stockPuts []stockCall
	getCalls  int
}

func (m *mockClient) WithBaseURL(string) domain.IMeliClient { return m }

func (m *mockClient) GetItem(context.Context, string, string) (*domain.MeliItemDetail, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getCalls++
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.item, nil
}

func (m *mockClient) UpdateStock(_ context.Context, _, itemID, variantID string, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stockPuts = append(m.stockPuts, stockCall{ItemID: itemID, VariantID: variantID, Quantity: quantity})
	if len(m.putErrs) > 0 {
		err := m.putErrs[0]
		m.putErrs = m.putErrs[1:]
		return err
	}
	return nil
}

func (m *mockClient) puts() []stockCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]stockCall, len(m.stockPuts))
	copy(out, m.stockPuts)
	return out
}

type mockService struct {
	integration *domain.Integration
	getErr      error
	token       string
	decryptErr  error

	decryptCalls int
}

func (m *mockService) GetIntegrationByID(context.Context, string) (*domain.Integration, error) {
	return m.integration, m.getErr
}

func (m *mockService) DecryptCredential(_ context.Context, _, field string) (string, error) {
	m.decryptCalls++
	if m.decryptErr != nil {
		return "", m.decryptErr
	}
	return m.token, nil
}

func (m *mockService) UpdateIntegrationConfig(context.Context, string, map[string]interface{}) error {
	return nil
}

func (m *mockService) UpdateIntegrationCredentials(context.Context, string, map[string]interface{}) error {
	return nil
}

func (m *mockService) GetIntegrationByStoreID(context.Context, string) (*domain.Integration, error) {
	return m.integration, nil
}

type markCall struct {
	IntegrationID uint
	ProductID     string
	VariantID     string
	Quantity      int
}

type mockInventoryRepo struct {
	mapped    []domain.MappedItem
	mappedErr error

	stock    map[string]int
	stockErr error

	pushState    *domain.PushState
	pushStateErr error

	markErr     error
	logisticErr error

	mu           sync.Mutex
	marks        []markCall
	logisticSets []string
	stateQueries []string
}

func (m *mockInventoryRepo) ListMappedItems(context.Context, uint) ([]domain.MappedItem, error) {
	return m.mapped, m.mappedErr
}

func (m *mockInventoryRepo) GetStockForProducts(context.Context, []string, []uint) (map[string]int, error) {
	if m.stockErr != nil {
		return nil, m.stockErr
	}
	if m.stock == nil {
		return map[string]int{}, nil
	}
	return m.stock, nil
}

func (m *mockInventoryRepo) GetInventoryByWarehouses(context.Context, []string, []uint) (map[string]map[uint]int, error) {
	return map[string]map[uint]int{}, nil
}

func (m *mockInventoryRepo) UpdateLogisticType(_ context.Context, _ uint, externalItemID, logisticType string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logisticSets = append(m.logisticSets, externalItemID+"="+logisticType)
	return m.logisticErr
}

func (m *mockInventoryRepo) MarkPushedQty(_ context.Context, integrationID uint, productID, variantID string, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.marks = append(m.marks, markCall{IntegrationID: integrationID, ProductID: productID, VariantID: variantID, Quantity: quantity})
	return m.markErr
}

func (m *mockInventoryRepo) GetPushState(_ context.Context, _ uint, productID, variantID string) (*domain.PushState, error) {
	m.mu.Lock()
	m.stateQueries = append(m.stateQueries, productID+"|"+variantID)
	m.mu.Unlock()
	return m.pushState, m.pushStateErr
}

func (m *mockInventoryRepo) markedQty() []markCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]markCall, len(m.marks))
	copy(out, m.marks)
	return out
}

func intPtr(v int) *int { return &v }

func newUseCase(cli domain.IMeliClient, svc domain.IIntegrationService, repo domain.IInventoryRepository) *meliUseCase {
	return &meliUseCase{
		client:        cli,
		service:       svc,
		inventoryRepo: repo,
		logger:        testLogger(),
	}
}
