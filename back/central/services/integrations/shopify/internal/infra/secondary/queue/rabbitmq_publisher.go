package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/queue/mapper"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type rabbitMQPublisher struct {
	queue     rabbitmq.IQueue
	logger    log.ILogger
	queueName string
}

func New(queue rabbitmq.IQueue, logger log.ILogger, config env.IConfig) domain.OrderPublisher {
	queueName := config.Get("RABBITMQ_ORDERS_CANONICAL_QUEUE")
	if queueName == "" {
		// Fallback al valor por defecto si no está configurado
		queueName = "probability.orders.canonical"
		logger.Warn(context.Background()).
			Str("queue_name", queueName).
			Msg("RABBITMQ_ORDERS_CANONICAL_QUEUE not set, using default queue name")
	}

	return &rabbitMQPublisher{
		queue:     queue,
		logger:    logger,
		queueName: queueName,
	}
}

func (p *rabbitMQPublisher) Publish(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
	// Validar que la cola de RabbitMQ esté disponible
	if p.queue == nil {
		p.logger.Warn(ctx).
			Str("order_number", order.OrderNumber).
			Str("platform", order.Platform).
			Uint("integration_id", order.IntegrationID).
			Msg("RabbitMQ queue not available, skipping order publication")
		return fmt.Errorf("rabbitmq queue not available")
	}

	// Mapear desde dominio (sin etiquetas) a estructura de serialización (con etiquetas JSON)
	orderForSerialization := mapper.MapDomainToSerializable(order)

	// Serializar la orden canónica a JSON
	orderJSON, err := json.Marshal(orderForSerialization)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("order_number", order.OrderNumber).
			Msg("Failed to marshal order to JSON")
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	// Publicar a la cola de RabbitMQ
	if err := p.queue.Publish(ctx, p.queueName, orderJSON); err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("queue", p.queueName).
			Str("order_number", order.OrderNumber).
			Msg("Failed to publish order to queue")
		return fmt.Errorf("failed to publish order to queue: %w", err)
	}

	p.logger.Info(ctx).
		Str("queue", p.queueName).
		Str("order_number", order.OrderNumber).
		Str("platform", order.Platform).
		Uint("integration_id", order.IntegrationID).
		Msg("Order published to queue successfully")

	return nil
}
