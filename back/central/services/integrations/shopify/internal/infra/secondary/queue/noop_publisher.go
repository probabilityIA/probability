package queue

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// noOpPublisher es un publisher que no hace nada (cuando RabbitMQ no está disponible)
type noOpPublisher struct {
	logger log.ILogger
}

// NewNoOpPublisher crea un publisher que no publica nada
func NewNoOpPublisher(logger log.ILogger) domain.OrderPublisher {
	return &noOpPublisher{
		logger: logger,
	}
}

// Publish no hace nada, solo registra un warning
func (p *noOpPublisher) Publish(ctx context.Context, order *domain.ProbabilityOrderDTO) error {
	p.logger.Warn(ctx).
		Str("order_number", order.OrderNumber).
		Str("platform", order.Platform).
		Uint("integration_id", order.IntegrationID).
		Msg("RabbitMQ not available, order not published to queue (no-op publisher)")
	return nil
}


