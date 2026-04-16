package consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/infra/primary/queue/consumer/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Start inicia el consumer de la cola de emails
func (c *emailConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rabbitmq.QueueMessagingEmailRequests).
		Msg("Iniciando consumer de email notifications")

	return c.rabbitMQ.Consume(ctx, rabbitmq.QueueMessagingEmailRequests, c.handleMessage)
}

// handleMessage procesa un mensaje individual de la cola
func (c *emailConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var event request.EmailNotificationEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("raw", string(body)).
			Msg("Error deserializando mensaje de email")
		return nil // No reintentar mensajes malformados
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

	// El use case maneja el envío + logging internamente
	_ = c.useCase.SendNotificationEmail(ctx, dto)

	return nil
}
