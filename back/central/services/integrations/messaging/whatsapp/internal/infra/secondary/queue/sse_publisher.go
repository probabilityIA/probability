package queue

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	whatsappMessageReceived      = "whatsapp.message_received"
	whatsappConversationStarted  = "whatsapp.conversation_started"
	whatsappMessageStatusUpdated = "whatsapp.message_status_updated"
)

// SSEPublisher publica eventos WhatsApp al exchange de eventos para notificar al frontend vía SSE.
type SSEPublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewSSEPublisher crea un publicador de eventos SSE para WhatsApp.
func NewSSEPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) ports.ISSEEventPublisher {
	return &SSEPublisher{rabbit: rabbit, log: logger}
}

// PublishMessageReceived publica evento cuando llega un mensaje inbound del cliente.
func (p *SSEPublisher) PublishMessageReceived(ctx context.Context, businessID uint, conversationID, phoneNumber, messageID, content string) error {
	envelope := rabbitmq.EventEnvelope{
		ID:         uuid.New().String(),
		Type:       whatsappMessageReceived,
		Category:   "whatsapp",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"conversation_id": conversationID,
			"phone_number":    phoneNumber,
			"message_id":      messageID,
			"content":         content,
			"direction":       "inbound",
		},
	}

	if err := rabbitmq.PublishEvent(ctx, p.rabbit, envelope); err != nil {
		p.log.Error(ctx).Err(err).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp SSE] - error publicando whatsapp.message_received")
		return err
	}

	p.log.Info(ctx).
		Uint("business_id", businessID).
		Str("conversation_id", conversationID).
		Msg("[WhatsApp SSE] - whatsapp.message_received publicado")
	return nil
}

// PublishConversationStarted publica evento cuando se inicia una nueva conversación.
func (p *SSEPublisher) PublishConversationStarted(ctx context.Context, businessID uint, conversationID, phoneNumber string) error {
	envelope := rabbitmq.EventEnvelope{
		ID:         uuid.New().String(),
		Type:       whatsappConversationStarted,
		Category:   "whatsapp",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"conversation_id": conversationID,
			"phone_number":    phoneNumber,
		},
	}

	if err := rabbitmq.PublishEvent(ctx, p.rabbit, envelope); err != nil {
		p.log.Error(ctx).Err(err).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp SSE] - error publicando whatsapp.conversation_started")
		return err
	}

	p.log.Info(ctx).
		Uint("business_id", businessID).
		Str("conversation_id", conversationID).
		Msg("[WhatsApp SSE] - whatsapp.conversation_started publicado")
	return nil
}

// PublishMessageStatusUpdated publica evento cuando cambia el estado de un mensaje outbound.
func (p *SSEPublisher) PublishMessageStatusUpdated(ctx context.Context, businessID uint, messageID, status string) error {
	envelope := rabbitmq.EventEnvelope{
		ID:         uuid.New().String(),
		Type:       whatsappMessageStatusUpdated,
		Category:   "whatsapp",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"message_id": messageID,
			"status":     status,
		},
	}

	if err := rabbitmq.PublishEvent(ctx, p.rabbit, envelope); err != nil {
		p.log.Error(ctx).Err(err).
			Str("message_id", messageID).
			Msg("[WhatsApp SSE] - error publicando whatsapp.message_status_updated")
		return err
	}

	p.log.Info(ctx).
		Uint("business_id", businessID).
		Str("message_id", messageID).
		Str("status", status).
		Msg("[WhatsApp SSE] - whatsapp.message_status_updated publicado")
	return nil
}

// noopSSEPublisher no hace nada (cuando RabbitMQ no está disponible).
type noopSSEPublisher struct{}

func NewNoopSSEPublisher() ports.ISSEEventPublisher { return &noopSSEPublisher{} }

func (n *noopSSEPublisher) PublishMessageReceived(_ context.Context, _ uint, _, _, _, _ string) error {
	return nil
}
func (n *noopSSEPublisher) PublishConversationStarted(_ context.Context, _ uint, _, _ string) error {
	return nil
}
func (n *noopSSEPublisher) PublishMessageStatusUpdated(_ context.Context, _ uint, _, _ string) error {
	return nil
}

// Compile-time interface checks
var _ ports.ISSEEventPublisher = (*SSEPublisher)(nil)
var _ ports.ISSEEventPublisher = (*noopSSEPublisher)(nil)
