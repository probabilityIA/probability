package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IExitoUseCase define las operaciones de negocio de Exito.
type IExitoUseCase interface {
	// TestConnection verifica que las credenciales de una integracion sean validas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type exitoUseCase struct {
	client    domain.IExitoClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de Exito con todas sus dependencias.
func New(
	client domain.IExitoClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) IExitoUseCase {
	return &exitoUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("exito"),
	}
}
