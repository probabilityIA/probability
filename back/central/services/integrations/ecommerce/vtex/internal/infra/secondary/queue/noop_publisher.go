package queue

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type noOpPublisher struct {
	logger log.ILogger
}

// NewNoOpPublisher crea un publisher que descarta las órdenes (cuando RabbitMQ no está disponible).
func NewNoOpPublisher(logger log.ILogger) domain.OrderPublisher {
	return &noOpPublisher{logger: logger}
}

func (p *noOpPublisher) Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error {
	p.logger.Warn(ctx).
		Str("order_number", order.OrderNumber).
		Uint("integration_id", order.IntegrationID).
		Msg("RabbitMQ not available, VTEX order not published to queue (no-op)")
	return nil
}
