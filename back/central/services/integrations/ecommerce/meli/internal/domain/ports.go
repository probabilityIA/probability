package domain

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

// IMeliClient define las operaciones del cliente HTTP de MercadoLibre.
// Implementado en infra/secondary/client.
type IMeliClient interface {
	// TestConnection verifica que las credenciales sean válidas
	TestConnection(ctx context.Context, accessToken string) error
}

// IIntegrationService define las operaciones del core de integraciones
// que el módulo de MercadoLibre necesita.
type IIntegrationService interface {
	GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error)
	DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error)
	UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error
}

// OrderPublisher publica órdenes al canal canónico de RabbitMQ.
// Implementado en infra/secondary/queue.
type OrderPublisher interface {
	Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
}
