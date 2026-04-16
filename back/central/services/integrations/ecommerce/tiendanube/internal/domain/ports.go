package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

// ITiendanubeClient define las operaciones del cliente HTTP de Tiendanube.
// Implementado en infra/secondary/client.
type ITiendanubeClient interface {
	// TestConnection verifica que las credenciales sean validas
	TestConnection(ctx context.Context, storeURL, accessToken string) error
}

// IIntegrationService define las operaciones del core de integraciones
// que el modulo de Tiendanube necesita.
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
