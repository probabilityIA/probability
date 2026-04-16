package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ITiendanubeUseCase define las operaciones de negocio de Tiendanube.
type ITiendanubeUseCase interface {
	// TestConnection verifica que las credenciales de una integracion sean validas.
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type tiendanubeUseCase struct {
	client    domain.ITiendanubeClient
	service   domain.IIntegrationService
	publisher domain.OrderPublisher
	logger    log.ILogger
}

// New crea el use case de Tiendanube con todas sus dependencias.
func New(
	client domain.ITiendanubeClient,
	service domain.IIntegrationService,
	publisher domain.OrderPublisher,
	logger log.ILogger,
) ITiendanubeUseCase {
	return &tiendanubeUseCase{
		client:    client,
		service:   service,
		publisher: publisher,
		logger:    logger.WithModule("tiendanube"),
	}
}
