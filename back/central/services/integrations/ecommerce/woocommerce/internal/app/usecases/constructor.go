package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IWooCommerceUseCase define las operaciones de negocio de WooCommerce.
type IWooCommerceUseCase interface {
	// TestConnection verifica que las credenciales de una integración sean válidas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type wooCommerceUseCase struct {
	client    domain.IWooCommerceClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de WooCommerce con todas sus dependencias.
func New(
	client domain.IWooCommerceClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) IWooCommerceUseCase {
	return &wooCommerceUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("woocommerce"),
	}
}
