package whatsapp_persistence_consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/whatsapp_persistence_consumer/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (c *PersistenceConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rabbitmq.QueueWhatsAppPersistenceEvents).
		Msg("Iniciando consumer unificado de WhatsApp persistence events")

	if err := c.rabbitMQ.DeclareQueue(rabbitmq.QueueWhatsAppPersistenceEvents, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", rabbitmq.QueueWhatsAppPersistenceEvents, err)
	}

	return c.rabbitMQ.Consume(ctx, rabbitmq.QueueWhatsAppPersistenceEvents, c.handleMessage)
}

func (c *PersistenceConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var event request.PersistenceEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("raw", string(body)).
			Msg("Error deserializando persistence event")
		return nil
	}

	switch event.EventType {
	case "conversation.created":
		c.handleConversationCreated(ctx, &event)
	case "conversation.updated":
		c.handleConversationUpdated(ctx, &event)
	case "conversation.upsert":
		c.handleConversationUpsert(ctx, &event)
	case "conversation.expired":
		c.handleConversationExpired(ctx, &event)
	case "messagelog.created":
		c.handleMessageLogCreated(ctx, &event)
	case "messagelog.status_updated":
		c.handleMessageLogStatusUpdated(ctx, &event)
	default:
		c.logger.Warn(ctx).
			Str("event_type", event.EventType).
			Msg("Tipo de evento de persistence no soportado, descartando")
	}

	return nil
}

func (c *PersistenceConsumer) handleConversationCreated(ctx context.Context, event *request.PersistenceEvent) {
	if event.Conversation == nil {
		c.logger.Warn(ctx).Msg("Evento conversation.created sin payload")
		return
	}
	conv := toConversationEntity(event.Conversation)
	if err := c.persister.CreateConversation(ctx, conv); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("conversation_id", event.Conversation.ID).
			Msg("Error persistiendo conversacion creada")
		return
	}
	c.logger.Info(ctx).
		Str("conversation_id", event.Conversation.ID).
		Str("phone_number", event.Conversation.PhoneNumber).
		Msg("Conversacion WhatsApp persistida")
}

func (c *PersistenceConsumer) handleConversationUpdated(ctx context.Context, event *request.PersistenceEvent) {
	if event.Conversation == nil {
		c.logger.Warn(ctx).Msg("Evento conversation.updated sin payload")
		return
	}
	conv := toConversationEntity(event.Conversation)
	if err := c.persister.UpdateConversation(ctx, conv); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("conversation_id", event.Conversation.ID).
			Msg("Error actualizando conversacion")
		return
	}
	c.logger.Info(ctx).
		Str("conversation_id", event.Conversation.ID).
		Str("state", event.Conversation.CurrentState).
		Msg("Conversacion WhatsApp actualizada")
}

func (c *PersistenceConsumer) handleConversationUpsert(ctx context.Context, event *request.PersistenceEvent) {
	if event.Conversation == nil {
		c.logger.Warn(ctx).Msg("Evento conversation.upsert sin payload")
		return
	}
	conv := toConversationEntity(event.Conversation)
	if err := c.persister.UpdateConversation(ctx, conv); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("conversation_id", event.Conversation.ID).
			Msg("Error en upsert de conversacion")
		return
	}
	c.logger.Info(ctx).
		Str("conversation_id", event.Conversation.ID).
		Str("phone_number", event.Conversation.PhoneNumber).
		Str("state", event.Conversation.CurrentState).
		Msg("Conversacion WhatsApp upserted")
}

func (c *PersistenceConsumer) handleConversationExpired(ctx context.Context, event *request.PersistenceEvent) {
	if event.Conversation == nil {
		c.logger.Warn(ctx).Msg("Evento conversation.expired sin payload")
		return
	}
	if err := c.persister.ExpireConversation(ctx, event.Conversation.ID); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("conversation_id", event.Conversation.ID).
			Msg("Error expirando conversacion")
		return
	}
	c.logger.Info(ctx).
		Str("conversation_id", event.Conversation.ID).
		Msg("Conversacion WhatsApp expirada")
}

func (c *PersistenceConsumer) handleMessageLogCreated(ctx context.Context, event *request.PersistenceEvent) {
	if event.MessageLog == nil {
		c.logger.Warn(ctx).Msg("Evento messagelog.created sin payload")
		return
	}
	entry := toMessageLogEntity(event.MessageLog)
	if err := c.persister.CreateMessageLog(ctx, entry); err != nil {
		c.logger.Error(ctx).Err(err).
			Str("message_id", event.MessageLog.MessageID).
			Str("conversation_id", event.MessageLog.ConversationID).
			Msg("Error persistiendo message log")
		return
	}
	c.logger.Info(ctx).
		Str("message_id", event.MessageLog.MessageID).
		Str("conversation_id", event.MessageLog.ConversationID).
		Msg("Message log WhatsApp persistido")
}

func (c *PersistenceConsumer) handleMessageLogStatusUpdated(ctx context.Context, event *request.PersistenceEvent) {
	if event.Update == nil {
		c.logger.Warn(ctx).Msg("Evento messagelog.status_updated sin payload")
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
		c.logger.Error(ctx).Err(err).
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

func toConversationEntity(p *request.ConversationPayload) *entities.WhatsAppConversation {
	return &entities.WhatsAppConversation{
		ID:             p.ID,
		PhoneNumber:    p.PhoneNumber,
		OrderNumber:    p.OrderNumber,
		BusinessID:     p.BusinessID,
		CurrentState:   p.CurrentState,
		LastMessageID:  p.LastMessageID,
		LastTemplateID: p.LastTemplateID,
		Metadata:       p.Metadata,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		ExpiresAt:      p.ExpiresAt,
	}
}

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
