package publisher

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type ScoreEventPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

func New(queue rabbitmq.IQueue, logger log.ILogger) ports.IScoreEventPublisher {
	return &ScoreEventPublisher{queue: queue, logger: logger}
}

func (p *ScoreEventPublisher) PublishScoreCalculated(ctx context.Context, orderID, orderNumber string, businessID, integrationID uint) error {
	return rabbitmq.PublishEvent(ctx, p.queue, rabbitmq.EventEnvelope{
		Type:          "order.score_calculated",
		Category:      "order",
		BusinessID:    businessID,
		IntegrationID: integrationID,
		Data: map[string]interface{}{
			"order_id":     orderID,
			"order_number": orderNumber,
		},
	})
}
