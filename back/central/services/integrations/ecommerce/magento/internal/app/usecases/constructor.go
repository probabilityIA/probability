package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/magento/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IMagentoUseCase define las operaciones de negocio de Magento.
type IMagentoUseCase interface {
	// TestConnection verifica que las credenciales de una integración sean válidas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type magentoUseCase struct {
	client    domain.IMagentoClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de Magento con todas sus dependencias.
func New(
	client domain.IMagentoClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) IMagentoUseCase {
	return &magentoUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("magento"),
	}
}
