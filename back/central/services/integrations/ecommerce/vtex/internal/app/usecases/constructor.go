package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type IVTEXUseCase interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error

	SyncOrders(ctx context.Context, integrationID string) error
	SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error
	ProcessWebhook(ctx context.Context, payload *domain.VTEXWebhookPayload) error

	SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error)
	ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	UpdateProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error

	AssertIntegrationOwned(ctx context.Context, integrationID string, businessID uint) error
	CreateWebhooks(ctx context.Context, integrationID, baseURL string, force bool) error
	ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookItem, error)
	InspectWebhook(ctx context.Context, integrationID, baseURL string) (*domain.WebhookItem, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error

	GetWarehouses(ctx context.Context, integrationID string, businessID uint) (*domain.WarehousesInfo, error)
	SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error
	PushStock(ctx context.Context, integrationID, productID, productExternalID string, quantity int) error
}

type vtexUseCase struct {
	client         domain.IVTEXClient
	service        domain.IIntegrationService
	publisher      domain.OrderPublisher
	productRepo    domain.IProductRepository
	rabbit         rabbitmq.IQueue
	webhookBaseURL string
	logger         log.ILogger
}

func New(
	client domain.IVTEXClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	productRepo domain.IProductRepository,
	rabbit rabbitmq.IQueue,
	webhookBaseURL string,
	logger log.ILogger,
) IVTEXUseCase {
	return &vtexUseCase{
		client:         client,
		service:        service,
		publisher:      publisher,
		productRepo:    productRepo,
		rabbit:         rabbit,
		webhookBaseURL: webhookBaseURL,
		logger:         logger.WithModule("vtex"),
	}
}
