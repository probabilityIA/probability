package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/scheduled"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type scheduledSendPublisher struct {
	rabbit rabbitmq.IQueue
	logger log.ILogger
}

func NewScheduledSendPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) scheduled.ISendPublisher {
	return &scheduledSendPublisher{
		rabbit: rabbit,
		logger: logger.WithModule("scheduled_send_publisher"),
	}
}

func (p *scheduledSendPublisher) PublishScheduledSend(ctx context.Context, message dtos.ScheduledSendMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error serializando el envio programado: %w", err)
	}

	return p.rabbit.Publish(ctx, rabbitmq.QueueScheduledNotificationSends, payload)
}
