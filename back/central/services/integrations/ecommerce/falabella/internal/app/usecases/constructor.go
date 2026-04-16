package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IFalabellaUseCase define las operaciones de negocio de Falabella.
type IFalabellaUseCase interface {
	// TestConnection verifica que las credenciales de una integración sean válidas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type falabellaUseCase struct {
	client    domain.IFalabellaClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de Falabella con todas sus dependencias.
func New(
	client domain.IFalabellaClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) IFalabellaUseCase {
	return &falabellaUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("falabella"),
	}
}
