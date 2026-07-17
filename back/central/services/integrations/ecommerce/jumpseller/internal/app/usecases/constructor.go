package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type IJumpsellerUseCase interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error

	SyncOrders(ctx context.Context, integrationID string) error

	SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error

	ProcessWebhookOrder(ctx context.Context, event string, storeCode string, integrationID string, rawBody []byte) error

	ResolveHooksToken(ctx context.Context, integrationID string) (string, error)

	CreateWebhooks(ctx context.Context, integrationID, baseURL string) error

	ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookItem, error)

	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error

	UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error

	SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error)

	GetLocations(ctx context.Context, integrationID string, businessID uint) (*domain.LocationsInfo, error)

	ApplyProductsToJumpseller(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	UpdateProductsToJumpseller(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	UpdateProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error

	UpdateOrderStatus(ctx context.Context, integrationID string, externalOrderID string, probabilityStatus string, tracking domain.UpdateOrderFields) error
}

type jumpsellerUseCase struct {
	client      domain.IJumpsellerClient
	service     domain.IIntegrationService
	publisher   domain.OrderPublisher
	productRepo domain.IProductRepository
	rabbit      rabbitmq.IQueue
	logger      log.ILogger
}

func New(
	client domain.IJumpsellerClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	productRepo domain.IProductRepository,
	rabbit rabbitmq.IQueue,
	logger log.ILogger,
) IJumpsellerUseCase {
	return &jumpsellerUseCase{
		client:      client,
		service:     service,
		publisher:   publisher,
		productRepo: productRepo,
		rabbit:      rabbit,
		logger:      logger.WithModule("jumpseller"),
	}
}
