package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type ITiendanubeUseCase interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error

	ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error)
	ReconcileProductsAsync(ctx context.Context, integrationID string, businessID, integIDUint uint, correlationID string)

	SyncProducts(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	ApplyProductsToTiendanube(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error
	UpdateProductsToTiendanube(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error
	ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error
	UpdateProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string, skus ...string) error
	AssociateProducts(ctx context.Context, integrationID string, businessID uint, correlationID string, skus []string) error

	SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error
	CompareInventory(ctx context.Context, integrationID string, businessID uint, page, pageSize int, skus ...string) (*inventorycompare.Page, error)
	LoadInventoryCompare(ctx context.Context, integrationID string, businessID uint, opts inventorycompare.LoadOptions) (*inventorycompare.Page, error)

	SyncOrders(ctx context.Context, integrationID string, filters domain.OrderFilters) (int, error)
	ProcessOrderEvent(ctx context.Context, integrationID, event, orderID string) error

	CreateWebhooks(ctx context.Context, integrationID, baseURL string) (*domain.CreateWebhooksResult, error)
	ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookItem, error)
	DeleteWebhook(ctx context.Context, integrationID, webhookID string) error
}

type tiendanubeUseCase struct {
	client      domain.ITiendanubeClient
	service     domain.IIntegrationService
	publisher   domain.OrderPublisher
	productRepo domain.IProductRepository
	rabbit      rabbitmq.IQueue
	logger      log.ILogger
}

func New(
	client domain.ITiendanubeClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	productRepo domain.IProductRepository,
	rabbit rabbitmq.IQueue,
	logger log.ILogger,
) ITiendanubeUseCase {
	return &tiendanubeUseCase{
		client:      client,
		service:     service,
		publisher:   publisher,
		productRepo: productRepo,
		rabbit:      rabbit,
		logger:      logger.WithModule("tiendanube"),
	}
}
