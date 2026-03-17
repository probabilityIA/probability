package whatsapp_conversation_consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/whatsapp_conversation_consumer/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Start inicia el consumer de conversaciones WhatsApp
func (c *ConversationConsumer) Start(ctx context.Context) error {
	c.logger.Info(ctx).
		Str("queue", rabbitmq.QueueWhatsAppConversationEvents).
		Msg("Iniciando consumer de WhatsApp conversation events")

	if err := c.rabbitMQ.DeclareQueue(rabbitmq.QueueWhatsAppConversationEvents, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", rabbitmq.QueueWhatsAppConversationEvents, err)
	}

	return c.rabbitMQ.Consume(ctx, rabbitmq.QueueWhatsAppConversationEvents, c.handleMessage)
}

// handleMessage procesa un evento de conversación individual
func (c *ConversationConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var event request.ConversationEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("raw", string(body)).
			Msg("Error deserializando conversation event")
		return nil // No reintentar mensajes malformados
	}

	switch event.EventType {
	case "conversation.created":
		c.handleCreated(ctx, &event)
	case "conversation.updated":
		c.handleUpdated(ctx, &event)
	case "conversation.expired":
		c.handleExpired(ctx, &event)
	default:
		c.logger.Warn(ctx).
			Str("event_type", event.EventType).
			Msg("Tipo de evento de conversación no soportado, descartando")
	}

	return nil
}

func (c *ConversationConsumer) handleCreated(ctx context.Context, event *request.ConversationEvent) {
	conv := toConversationEntity(&event.Conversation)

	if err := c.persister.CreateConversation(ctx, conv); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("conversation_id", event.Conversation.ID).
			Msg("Error persistiendo conversación creada")
		return
	}

	c.logger.Info(ctx).
		Str("conversation_id", event.Conversation.ID).
		Str("phone_number", event.Conversation.PhoneNumber).
		Msg("Conversación WhatsApp persistida")
}

func (c *ConversationConsumer) handleUpdated(ctx context.Context, event *request.ConversationEvent) {
	conv := toConversationEntity(&event.Conversation)

	if err := c.persister.UpdateConversation(ctx, conv); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("conversation_id", event.Conversation.ID).
			Msg("Error actualizando conversación")
		return
	}

	c.logger.Info(ctx).
		Str("conversation_id", event.Conversation.ID).
		Str("state", event.Conversation.CurrentState).
		Msg("Conversación WhatsApp actualizada")
}

func (c *ConversationConsumer) handleExpired(ctx context.Context, event *request.ConversationEvent) {
	if err := c.persister.ExpireConversation(ctx, event.Conversation.ID); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("conversation_id", event.Conversation.ID).
			Msg("Error expirando conversación")
		return
	}

	c.logger.Info(ctx).
		Str("conversation_id", event.Conversation.ID).
		Msg("Conversación WhatsApp expirada")
}

// toConversationEntity convierte el payload de RabbitMQ a entidad de dominio
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
