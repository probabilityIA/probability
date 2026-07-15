package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type IMeliUseCase interface {
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error

	ProcessNotification(ctx context.Context, notification *domain.MeliNotification) error

	SyncOrders(ctx context.Context, integrationID string) error

	SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error

	EnsureValidToken(ctx context.Context, integrationID string) (string, error)

	ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error)
	ApplyProductsToMeli(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	UpdateItemStock(ctx context.Context, integrationID, itemID string, quantity int) error

	PushOrderStatus(ctx context.Context, integrationID string, shipmentID int64, status string) error

	RetryBilling(ctx context.Context, integrationID string, orderID int64) (bool, error)
}

type meliUseCase struct {
	client        domain.IMeliClient
	service       domain.IIntegrationService
	publisher     domain.OrderPublisher
	productRepo   domain.IProductRepository
	inventoryRepo domain.IInventoryRepository
	rabbit        rabbitmq.IQueue
	logger        log.ILogger
}

func New(
	client domain.IMeliClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	productRepo domain.IProductRepository,
	inventoryRepo domain.IInventoryRepository,
	rabbit rabbitmq.IQueue,
	logger log.ILogger,
) IMeliUseCase {
	return &meliUseCase{
		client:        client,
		service:       service,
		publisher:     publisher,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		rabbit:        rabbit,
		logger:        logger.WithModule("meli"),
	}
}
