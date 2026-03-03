package consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/consumer/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Start inicia el consumer de la cola de resultados de entrega
func (c *DeliveryResultConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rabbitmq.QueueNotificationDeliveryResults).
		Msg("Iniciando consumer de delivery results")

	return c.rabbitMQ.Consume(ctx, rabbitmq.QueueNotificationDeliveryResults, c.handleMessage)
}

// handleMessage procesa un mensaje individual de resultado de entrega
func (c *DeliveryResultConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var msg request.DeliveryResult
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("raw", string(body)).
			Msg("Error deserializando delivery result")
		return nil // No reintentar mensajes malformados
	}

	switch msg.Channel {
	case "email":
		c.handleEmailResult(ctx, &msg)
	default:
		c.logger.Warn(ctx).
			Str("channel", msg.Channel).
			Msg("Canal de delivery result no soportado, descartando")
	}

	return nil
}

// handleEmailResult persiste un resultado de entrega de email
func (c *DeliveryResultConsumer) handleEmailResult(ctx context.Context, msg *request.DeliveryResult) {
	entry := &entities.EmailDeliveryLog{
		BusinessID:    msg.BusinessID,
		IntegrationID: msg.IntegrationID,
		ConfigID:      msg.ConfigID,
		To:            msg.To,
		Subject:       msg.Subject,
		EventType:     msg.EventType,
		Status:        msg.Status,
		ErrorMessage:  msg.ErrorMessage,
		SentAt:        msg.SentAt,
	}

	if err := c.deliveryLogRepo.CreateEmailLog(ctx, entry); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("to", msg.To).
			Str("status", msg.Status).
			Msg("Error persistiendo email delivery log")
		return
	}

	c.logger.Info(ctx).
		Str("to", msg.To).
		Str("status", msg.Status).
		Str("event_type", msg.EventType).
		Msg("Email delivery log persistido")
}
