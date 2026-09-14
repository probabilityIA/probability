package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/infra/primary/queue/consumer/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (c *emailConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rabbitmq.QueueMessagingEmailRequests).
		Msg("Iniciando consumer de email notifications")

	if err := c.rabbitMQ.DeclareQueue(rabbitmq.QueueMessagingEmailRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", rabbitmq.QueueMessagingEmailRequests, err)
	}

	return c.rabbitMQ.Consume(ctx, rabbitmq.QueueMessagingEmailRequests, c.handleMessage)
}

func (c *emailConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var event request.EmailNotificationEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("raw", string(body)).
			Msg("Error deserializando mensaje de email")
		return nil
	}

	if event.CustomerEmail == "" {
		c.logger.Warn(ctx).
			Str("event_type", event.EventType).
			Uint("business_id", event.BusinessID).
			Msg("Mensaje sin customer_email, descartando")
		return nil
	}

	dto := dtos.SendEmailDTO{
		EventType:     event.EventType,
		BusinessID:    event.BusinessID,
		IntegrationID: event.IntegrationID,
		ConfigID:      event.ConfigID,
		CustomerEmail: event.CustomerEmail,
		EventData:     event.EventData,
	}

	_ = c.useCase.SendNotificationEmail(ctx, dto)

	return nil
}
