package whatsapp_messagelog_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/whatsapp_messagelog_consumer/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Start inicia el consumer de message logs WhatsApp
func (c *MessageLogConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rabbitmq.QueueWhatsAppMessageLogEvents).
		Msg("Iniciando consumer de WhatsApp message log events")

	if err := c.rabbitMQ.DeclareQueue(rabbitmq.QueueWhatsAppMessageLogEvents, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", rabbitmq.QueueWhatsAppMessageLogEvents, err)
	}

	return c.rabbitMQ.Consume(ctx, rabbitmq.QueueWhatsAppMessageLogEvents, c.handleMessage)
}

// handleMessage procesa un evento de message log individual
func (c *MessageLogConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var event request.MessageLogEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("raw", string(body)).
			Msg("Error deserializando message log event")
		return nil // No reintentar mensajes malformados
	}

	switch event.EventType {
	case "messagelog.created":
		c.handleCreated(ctx, &event)
	case "messagelog.status_updated":
		c.handleStatusUpdated(ctx, &event)
	default:
		c.logger.Warn(ctx).
			Str("event_type", event.EventType).
			Msg("Tipo de evento de message log no soportado, descartando")
	}

	return nil
}

func (c *MessageLogConsumer) handleCreated(ctx context.Context, event *request.MessageLogEvent) {
	if event.MessageLog == nil {
		c.logger.Warn(ctx).Msg("Evento messagelog.created sin payload de message_log")
		return
	}

	entry := toMessageLogEntity(event.MessageLog)

	if err := c.persister.CreateMessageLog(ctx, entry); err != nil {
		// Si es FK violation, la conversación puede no haberse persistido aún (race condition entre consumers)
		// Reintentar una vez después de esperar a que el conversation consumer la persista
		if isFKViolation(err) {
			c.logger.Warn(ctx).
				Str("conversation_id", event.MessageLog.ConversationID).
				Msg("FK violation en message log, esperando a que la conversación se persista...")

			time.Sleep(500 * time.Millisecond)

			if retryErr := c.persister.CreateMessageLog(ctx, entry); retryErr != nil {
				c.logger.Error(ctx).
					Err(retryErr).
					Str("message_id", event.MessageLog.MessageID).
					Str("conversation_id", event.MessageLog.ConversationID).
					Msg("Error persistiendo message log después de retry")
				return
			}

			c.logger.Info(ctx).
				Str("message_id", event.MessageLog.MessageID).
				Str("conversation_id", event.MessageLog.ConversationID).
				Msg("Message log WhatsApp persistido después de retry")
			return
		}

		c.logger.Error(ctx).
			Err(err).
			Str("message_id", event.MessageLog.MessageID).
			Msg("Error persistiendo message log")
		return
	}

	c.logger.Info(ctx).
		Str("message_id", event.MessageLog.MessageID).
		Str("conversation_id", event.MessageLog.ConversationID).
		Msg("Message log WhatsApp persistido")
}

// isFKViolation verifica si el error es una violación de foreign key (SQLSTATE 23503)
func isFKViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23503")
}

func (c *MessageLogConsumer) handleStatusUpdated(ctx context.Context, event *request.MessageLogEvent) {
	if event.Update == nil {
		c.logger.Warn(ctx).Msg("Evento messagelog.status_updated sin payload de update")
		return
	}

	var deliveredAt, readAt *string
	if v, ok := event.Update.Timestamps["delivered_at"]; ok {
		deliveredAt = &v
	}
	if v, ok := event.Update.Timestamps["read_at"]; ok {
		readAt = &v
	}

	if err := c.persister.UpdateMessageLogStatus(ctx, event.Update.MessageID, event.Update.Status, deliveredAt, readAt); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("message_id", event.Update.MessageID).
			Str("status", event.Update.Status).
			Msg("Error actualizando status de message log")
		return
	}

	c.logger.Info(ctx).
		Str("message_id", event.Update.MessageID).
		Str("status", event.Update.Status).
		Msg("Status de message log WhatsApp actualizado")
}

// toMessageLogEntity convierte el payload de RabbitMQ a entidad de dominio
func toMessageLogEntity(p *request.MessageLogPayload) *entities.WhatsAppMessageLogEntry {
	return &entities.WhatsAppMessageLogEntry{
		ID:             p.ID,
		ConversationID: p.ConversationID,
		Direction:      p.Direction,
		MessageID:      p.MessageID,
		TemplateName:   p.TemplateName,
		Content:        p.Content,
		Status:         p.Status,
		DeliveredAt:    p.DeliveredAt,
		ReadAt:         p.ReadAt,
		CreatedAt:      p.CreatedAt,
	}
}
