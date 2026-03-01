package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella/internal/infra/secondary/queue/mapper"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type rabbitMQPublisher struct {
	queue     rabbitmq.IQueue
	logger    log.ILogger
	queueName string
}

// New crea el publisher de órdenes hacia RabbitMQ.
func New(queue rabbitmq.IQueue, logger log.ILogger, config env.IConfig) domain.OrderPublisher {
	queueName := config.Get("RABBITMQ_ORDERS_CANONICAL_QUEUE")
	if queueName == "" {
		queueName = rabbitmq.QueueOrdersCanonical
		logger.Warn(context.Background()).
			Str("queue_name", queueName).
			Msg("RABBITMQ_ORDERS_CANONICAL_QUEUE not set, using default")
	}
	return &rabbitMQPublisher{
		queue:     queue,
		logger:    logger,
		queueName: queueName,
	}
}

// Publish serializa la orden canónica y la publica en RabbitMQ.
func (p *rabbitMQPublisher) Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error {
	if p.queue == nil {
		return fmt.Errorf("rabbitmq queue not available")
	}

	serializable := mapper.MapDomainToSerializable(order)

	body, err := json.Marshal(serializable)
	if err != nil {
		p.logger.Error(ctx).Err(err).Str("order_number", order.OrderNumber).
			Msg("Failed to marshal order to JSON")
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	if err := p.queue.Publish(ctx, p.queueName, body); err != nil {
		p.logger.Error(ctx).Err(err).
			Str("queue", p.queueName).
			Str("order_number", order.OrderNumber).
			Msg("Failed to publish order to queue")
		return fmt.Errorf("failed to publish order to queue: %w", err)
	}

	p.logger.Info(ctx).
		Str("queue", p.queueName).
		Str("order_number", order.OrderNumber).
		Uint("integration_id", order.IntegrationID).
		Msg("Order published to queue successfully")

	return nil
}
