package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IMeliUseCase define las operaciones de negocio de MercadoLibre.
type IMeliUseCase interface {
	// TestConnection verifica que las credenciales de una integración sean válidas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error

	// ProcessNotification procesa una notificación IPN de MercadoLibre.
	// Recibe la notificación parseada, obtiene la orden completa de la API y la publica.
	ProcessNotification(ctx context.Context, notification *domain.MeliNotification) error

	// SyncOrders sincroniza órdenes de los últimos 30 días.
	SyncOrders(ctx context.Context, integrationID string) error

	// SyncOrdersWithParams sincroniza órdenes con parámetros personalizados.
	SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error

	// EnsureValidToken verifica si el token actual está vigente, y si no, lo renueva.
	// Retorna un access_token listo para usar.
	EnsureValidToken(ctx context.Context, integrationID string) (string, error)

	// ReconcileProducts cruza los productos de ambos lados por SKU.
	ReconcileProducts(ctx context.Context, integrationID string, businessID uint) (*domain.ReconcileResult, error)
	// ApplyProductsToMeli crea en MercadoLibre los productos que solo existen en Probability.
	ApplyProductsToMeli(ctx context.Context, integrationID string, businessID uint, correlationID string) error
	// ApplyProductsToProbability crea en Probability los productos que solo existen en MercadoLibre.
	ApplyProductsToProbability(ctx context.Context, integrationID string, businessID uint, correlationID string) error

	// SyncInventory empuja el stock de Probability a las publicaciones de MercadoLibre.
	SyncInventory(ctx context.Context, integrationID string, businessID uint, correlationID string) error
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

// New crea el use case de MercadoLibre con todas sus dependencias.
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
