package publisher

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type logPublisher struct {
	logger log.ILogger
}

func New(logger log.ILogger) domain.OrderPublisher {
	return &logPublisher{
		logger: logger,
	}
}

func (p *logPublisher) Publish(ctx context.Context, order *domain.UnifiedOrder) error {
	// In a real implementation, this would publish to RabbitMQ.
	// For now, we just log the order.

	orderJSON, _ := json.Marshal(order)
	p.logger.Info(ctx).
		Str("component", "shopify_publisher").
		Str("order_number", order.OrderNumber).
		RawJSON("order_payload", orderJSON).
		Msg("Publishing unified order to queue (simulated)")

	return nil
}
