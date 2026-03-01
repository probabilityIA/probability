package eventpublisher

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type syncEventPublisher struct {
	queue rabbitmq.IQueue
}

// New crea un publisher de eventos de sincronización que publica al exchange de eventos via RabbitMQ
func New(queue rabbitmq.IQueue) domain.ISyncEventPublisher {
	return &syncEventPublisher{queue: queue}
}

func (p *syncEventPublisher) PublishSyncEvent(ctx context.Context, integrationID uint, businessID *uint, eventType string, data map[string]interface{}) {
	var bID uint
	if businessID != nil {
		bID = *businessID
	}

	rabbitmq.PublishEvent(ctx, p.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
		Type:          eventType,
		Category:      "integration",
		BusinessID:    bID,
		IntegrationID: integrationID,
		Data:          data,
	})
}
