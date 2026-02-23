package usecases

// Este archivo contiene todos los mocks necesarios para los tests del paquete usecases.
// Los mocks implementan las interfaces definidas en domain/ports.go y los paquetes compartidos.

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// ─── Mock: IIntegrationService ──────────────────────────────────────────────

type mockIntegrationService struct {
	GetIntegrationByIDFn         func(ctx context.Context, integrationID string) (*domain.Integration, error)
	GetIntegrationByExternalIDFn func(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error)
	DecryptCredentialFn          func(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfigFn    func(ctx context.Context, integrationID string, config map[string]interface{}) error
}

func (m *mockIntegrationService) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	if m.GetIntegrationByIDFn != nil {
		return m.GetIntegrationByIDFn(ctx, integrationID)
	}
	return nil, nil
}

func (m *mockIntegrationService) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.Integration, error) {
	if m.GetIntegrationByExternalIDFn != nil {
		return m.GetIntegrationByExternalIDFn(ctx, externalID, integrationType)
	}
	return nil, nil
}

func (m *mockIntegrationService) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	if m.DecryptCredentialFn != nil {
		return m.DecryptCredentialFn(ctx, integrationID, fieldName)
	}
	return "mock-token", nil
}

func (m *mockIntegrationService) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	if m.UpdateIntegrationConfigFn != nil {
		return m.UpdateIntegrationConfigFn(ctx, integrationID, config)
	}
	return nil
}

// ─── Mock: ShopifyClient ────────────────────────────────────────────────────

type mockShopifyClient struct {
	ValidateTokenFn func(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error)
	GetOrdersFn     func(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error)
	GetOrderFn      func(ctx context.Context, storeName, accessToken string, orderID string) (*domain.ShopifyOrder, error)
	CreateWebhookFn func(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error)
	ListWebhooksFn  func(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error)
	DeleteWebhookFn func(ctx context.Context, storeName, accessToken, webhookID string) error
	SetDebugFn      func(enabled bool)
}

func (m *mockShopifyClient) ValidateToken(ctx context.Context, storeName, accessToken string) (bool, map[string]interface{}, error) {
	if m.ValidateTokenFn != nil {
		return m.ValidateTokenFn(ctx, storeName, accessToken)
	}
	return true, nil, nil
}

func (m *mockShopifyClient) GetOrders(ctx context.Context, storeName, accessToken string, params *domain.GetOrdersParams) ([]domain.ShopifyOrder, string, error) {
	if m.GetOrdersFn != nil {
		return m.GetOrdersFn(ctx, storeName, accessToken, params)
	}
	return nil, "", nil
}

func (m *mockShopifyClient) GetOrder(ctx context.Context, storeName, accessToken string, orderID string) (*domain.ShopifyOrder, error) {
	if m.GetOrderFn != nil {
		return m.GetOrderFn(ctx, storeName, accessToken, orderID)
	}
	return nil, nil
}

func (m *mockShopifyClient) CreateWebhook(ctx context.Context, storeName, accessToken, webhookURL, event string) (string, error) {
	if m.CreateWebhookFn != nil {
		return m.CreateWebhookFn(ctx, storeName, accessToken, webhookURL, event)
	}
	return "mock-webhook-id", nil
}

func (m *mockShopifyClient) ListWebhooks(ctx context.Context, storeName, accessToken string) ([]domain.WebhookInfo, error) {
	if m.ListWebhooksFn != nil {
		return m.ListWebhooksFn(ctx, storeName, accessToken)
	}
	return nil, nil
}

func (m *mockShopifyClient) DeleteWebhook(ctx context.Context, storeName, accessToken, webhookID string) error {
	if m.DeleteWebhookFn != nil {
		return m.DeleteWebhookFn(ctx, storeName, accessToken, webhookID)
	}
	return nil
}

func (m *mockShopifyClient) SetDebug(enabled bool) {
	if m.SetDebugFn != nil {
		m.SetDebugFn(enabled)
	}
}

// ─── Mock: OrderPublisher ───────────────────────────────────────────────────

type mockOrderPublisher struct {
	PublishFn       func(ctx context.Context, order *domain.ProbabilityOrderDTO) error
	PublishedOrders []*domain.ProbabilityOrderDTO
}

func (m *mockOrderPublisher) Publish(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
	if m.PublishFn != nil {
		return m.PublishFn(ctx, order)
	}
	m.PublishedOrders = append(m.PublishedOrders, order)
	return nil
}

// ─── Mock: IDatabase ────────────────────────────────────────────────────────

type mockDatabase struct{}

func (m *mockDatabase) Connect(ctx context.Context) error        { return nil }
func (m *mockDatabase) Close() error                             { return nil }
func (m *mockDatabase) Conn(ctx context.Context) *gorm.DB        { return nil }
func (m *mockDatabase) WithContext(ctx context.Context) *gorm.DB { return nil }
func (m *mockDatabase) DebugConn(ctx context.Context) *gorm.DB   { return nil }

// ─── Mock: log.ILogger ──────────────────────────────────────────────────────

// mockLoggerILogger implementa log.ILogger usando zerolog.Nop() para que los
// logs de los use cases no emitan salida ni fallen durante la ejecucion de tests.
type mockLoggerILogger struct{}

func (m *mockLoggerILogger) Info(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Info()
}

func (m *mockLoggerILogger) Error(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Error()
}

func (m *mockLoggerILogger) Warn(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Warn()
}

func (m *mockLoggerILogger) Debug(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Debug()
}

func (m *mockLoggerILogger) Fatal(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Fatal()
}

func (m *mockLoggerILogger) Panic(ctx ...context.Context) *zerolog.Event {
	nop := zerolog.Nop()
	return nop.Panic()
}

func (m *mockLoggerILogger) With() zerolog.Context {
	nop := zerolog.Nop()
	return nop.With()
}

func (m *mockLoggerILogger) WithService(service string) log.ILogger     { return m }
func (m *mockLoggerILogger) WithModule(module string) log.ILogger       { return m }
func (m *mockLoggerILogger) WithBusinessID(businessID uint) log.ILogger { return m }

// ─── Helpers ─────────────────────────────────────────────────────────────────

// newTestUseCase construye una instancia de SyncOrdersUseCase con todos los mocks
// inyectados directamente, sin depender de infraestructura real (DB, logger, etc.).
func newTestUseCase(
	integrationSvc *mockIntegrationService,
	shopifyClient *mockShopifyClient,
	publisher *mockOrderPublisher,
) *SyncOrdersUseCase {
	return &SyncOrdersUseCase{
		integrationService: integrationSvc,
		shopifyClient:      shopifyClient,
		orderPublisher:     publisher,
		db:                 &mockDatabase{},
		log:                &mockLoggerILogger{},
	}
}

// newIntegration es un helper para crear una Integration de prueba con valores por defecto.
func newIntegration(id uint, storeName string) *domain.Integration {
	businessID := uint(99)
	return &domain.Integration{
		ID:         id,
		BusinessID: &businessID,
		Name:       storeName,
		Config: map[string]interface{}{
			"store_name": storeName,
		},
	}
}
