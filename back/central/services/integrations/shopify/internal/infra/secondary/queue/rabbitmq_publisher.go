package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	OrdersQueueName = "probability.orders.canonical"
)

type rabbitMQPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

func New(queue rabbitmq.IQueue, logger log.ILogger) domain.OrderPublisher {
	return &rabbitMQPublisher{
		queue:  queue,
		logger: logger,
	}
}

func (p *rabbitMQPublisher) Publish(ctx context.Context, order *domain.UnifiedOrder) error {
	// Map UnifiedOrder to CanonicalOrderDTO to ensure correct JSON structure for Consumer
	// Consumer expects flattened address fields or CanonicalAddressDTOs
	canonicalOrder := mapUnifiedToCanonical(order)

	// Serializar la orden canónica a JSON
	orderJSON, err := json.Marshal(canonicalOrder)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("order_number", order.OrderNumber).
			Msg("Failed to marshal order to JSON")
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	// Publicar a la cola de RabbitMQ
	if err := p.queue.Publish(ctx, OrdersQueueName, orderJSON); err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("queue", OrdersQueueName).
			Str("order_number", order.OrderNumber).
			Msg("Failed to publish order to queue")
		return fmt.Errorf("failed to publish order to queue: %w", err)
	}

	p.logger.Info(ctx).
		Str("queue", OrdersQueueName).
		Str("order_number", order.OrderNumber).
		Str("platform", order.Platform).
		Uint("integration_id", order.IntegrationID).
		Msg("Order published to queue successfully")

	return nil
}
