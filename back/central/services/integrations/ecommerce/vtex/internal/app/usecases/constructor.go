package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IVTEXUseCase define las operaciones de negocio de VTEX.
type IVTEXUseCase interface {
	// TestConnection verifica que las credenciales de una integración sean válidas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error

	// SyncOrders sincroniza órdenes de los últimos 30 días.
	SyncOrders(ctx context.Context, integrationID string) error

	// SyncOrdersWithParams sincroniza órdenes con parámetros personalizados.
	SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error

	// ProcessWebhook procesa un webhook de VTEX (cambio de estado de orden).
	ProcessWebhook(ctx context.Context, payload *domain.VTEXWebhookPayload) error
}

type vtexUseCase struct {
	client    domain.IVTEXClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de VTEX con todas sus dependencias.
func New(
	client domain.IVTEXClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) IVTEXUseCase {
	return &vtexUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("vtex"),
	}
}
