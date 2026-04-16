package eventpublisher

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type integrationEventPublisher struct {
	queue rabbitmq.IQueue
}

// New crea un publisher de eventos de integración que publica al exchange de eventos via RabbitMQ
func New(queue rabbitmq.IQueue) ports.IIntegrationEventPublisher {
	return &integrationEventPublisher{queue: queue}
}

func (p *integrationEventPublisher) PublishSyncOrderCreated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	var bID uint
	if businessID != nil {
		bID = *businessID
	}

	rabbitmq.PublishEvent(ctx, p.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
		Type:          "integration.sync.order.created",
		Category:      "integration",
		BusinessID:    bID,
		IntegrationID: integrationID,
		Data:          data,
	})
}

func (p *integrationEventPublisher) PublishSyncOrderUpdated(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	var bID uint
	if businessID != nil {
		bID = *businessID
	}

	rabbitmq.PublishEvent(ctx, p.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
		Type:          "integration.sync.order.updated",
		Category:      "integration",
		BusinessID:    bID,
		IntegrationID: integrationID,
		Data:          data,
	})
}

func (p *integrationEventPublisher) PublishSyncOrderRejected(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
	var bID uint
	if businessID != nil {
		bID = *businessID
	}

	rabbitmq.PublishEvent(ctx, p.queue, rabbitmq.EventEnvelope{ //nolint:errcheck
		Type:          "integration.sync.order.rejected",
		Category:      "integration",
		BusinessID:    bID,
		IntegrationID: integrationID,
		Data:          data,
	})
}
