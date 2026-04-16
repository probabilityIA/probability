package queue

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type storefrontPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

// NewStorefrontPublisher creates a publisher that sends orders to probability.orders.canonical
func NewStorefrontPublisher(queue rabbitmq.IQueue, logger log.ILogger) ports.IStorefrontPublisher {
	return &storefrontPublisher{
		queue:  queue,
		logger: logger,
	}
}

func (p *storefrontPublisher) PublishOrder(ctx context.Context, order []byte) error {
	if p.queue == nil {
		p.logger.Warn(ctx).Msg("Cola RabbitMQ no disponible, omitiendo publicacion de orden storefront")
		return fmt.Errorf("cola rabbitmq no disponible")
	}

	if err := p.queue.Publish(ctx, rabbitmq.QueueOrdersCanonical, order); err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("queue", rabbitmq.QueueOrdersCanonical).
			Msg("Error publicando orden storefront a cola")
		return fmt.Errorf("error publicando orden: %w", err)
	}

	p.logger.Info(ctx).
		Str("queue", rabbitmq.QueueOrdersCanonical).
		Msg("Orden storefront publicada exitosamente")

	return nil
}
