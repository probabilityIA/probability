package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type resultPublisher struct {
	rabbitMQ rabbitmq.IQueue
	logger   log.ILogger
}

// New crea un publisher que envía resultados de entrega a la cola de notification_config
func New(rabbitMQ rabbitmq.IQueue, logger log.ILogger) ports.IResultPublisher {
	return &resultPublisher{
		rabbitMQ: rabbitMQ,
		logger:   logger,
	}
}

// PublishResult serializa el DeliveryResult y lo publica a notification.delivery.results
func (p *resultPublisher) PublishResult(ctx context.Context, result *entities.DeliveryResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("error serializando delivery result: %w", err)
	}

	if err := p.rabbitMQ.Publish(ctx, rabbitmq.QueueNotificationDeliveryResults, data); err != nil {
		return fmt.Errorf("error publicando delivery result: %w", err)
	}

	p.logger.Info(ctx).
		Str("channel", result.Channel).
		Str("status", result.Status).
		Str("to", result.To).
		Msg("Resultado de entrega publicado")

	return nil
}
