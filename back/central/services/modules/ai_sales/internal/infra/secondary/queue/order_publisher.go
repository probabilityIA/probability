package queue

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (p *orderPublisher) PublishOrder(ctx context.Context, orderPayload []byte) error {
	if err := p.rabbit.Publish(ctx, rabbitmq.QueueOrdersCanonical, orderPayload); err != nil {
		p.log.Error(ctx).
			Err(err).
			Msg("Error publicando orden AI a probability.orders.canonical")
		return fmt.Errorf("error publishing AI order: %w", err)
	}

	p.log.Info(ctx).Msg("Orden AI publicada a probability.orders.canonical")
	return nil
}
