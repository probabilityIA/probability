package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
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
}

type meliUseCase struct {
	client    domain.IMeliClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de MercadoLibre con todas sus dependencias.
func New(
	client domain.IMeliClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) IMeliUseCase {
	return &meliUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("meli"),
	}
}
