package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	// OrdersCanonicalQueueName es el nombre de la cola donde se publican órdenes canónicas
	OrdersCanonicalQueueName = "probability.orders.canonical"
)

// OrderPublisher publica órdenes canónicas a RabbitMQ
type OrderPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

// New crea una nueva instancia del publicador
func New(queue rabbitmq.IQueue, logger log.ILogger) domain.IOrderPublisher {
	return &OrderPublisher{
		queue:  queue,
		logger: logger,
	}
}

// PublishCanonicalOrder publica una orden canónica a la cola de RabbitMQ
func (p *OrderPublisher) PublishCanonicalOrder(ctx context.Context, order *domain.CanonicalOrderDTO) error {
	// Serializar la orden a JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		p.logger.Error().
			Err(err).
			Msg("Failed to marshal canonical order to JSON")
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	// Publicar a la cola de RabbitMQ
	if err := p.queue.Publish(ctx, OrdersCanonicalQueueName, orderJSON); err != nil {
		p.logger.Error().
			Err(err).
			Str("queue", OrdersCanonicalQueueName).
			Msg("Failed to publish canonical order to queue")
		return fmt.Errorf("failed to publish order to queue: %w", err)
	}

	p.logger.Info().
		Str("queue", OrdersCanonicalQueueName).
		Msg("Canonical order published to queue successfully")

	return nil
}
