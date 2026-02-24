package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/amazon/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IAmazonUseCase define las operaciones de negocio de Amazon.
type IAmazonUseCase interface {
	// TestConnection verifica que las credenciales de una integracion sean validas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type amazonUseCase struct {
	client    domain.IAmazonClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de Amazon con todas sus dependencias.
func New(
	client domain.IAmazonClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) IAmazonUseCase {
	return &amazonUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("amazon"),
	}
}
