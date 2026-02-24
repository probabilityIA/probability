package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

// IExitoClient define las operaciones del cliente HTTP de Exito.
// Implementado en infra/secondary/client.
type IExitoClient interface {
	// TestConnection verifica que las credenciales sean validas
	TestConnection(ctx context.Context, apiKey, sellerID string) error
}

// IIntegrationService define las operaciones del core de integraciones
// que el modulo de Exito necesita.
type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
}

// OrderPublisher publica ordenes al canal canonico de RabbitMQ.
// Implementado en infra/secondary/queue.
type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
