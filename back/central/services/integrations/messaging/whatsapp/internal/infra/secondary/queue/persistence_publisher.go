package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// persistencePublisher publica eventos para persistencia asíncrona en DB
type persistencePublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewPersistencePublisher crea un publisher que envía eventos de WhatsApp a notification_config para persistir
func NewPersistencePublisher(rabbit rabbitmq.IQueue, logger log.ILogger) ports.IPersistencePublisher {
	return &persistencePublisher{
		rabbit: rabbit,
		log:    logger.WithModule("whatsapp-persistence-publisher"),
	}
}

// ============================================
// Conversation events
// ============================================

// conversationEvent es el payload para eventos de conversación
type conversationEvent struct {
	EventType    string                 `json:"event_type"`
	Conversation conversationPayload    `json:"conversation"`
	Timestamp    int64                  `json:"timestamp"`
}

type conversationPayload struct {
	ID             string                 `json:"id"`
	PhoneNumber    string                 `json:"phone_number"`
	OrderNumber    string                 `json:"order_number"`
	BusinessID     uint                   `json:"business_id"`
	CurrentState   string                 `json:"current_state"`
	LastMessageID  string                 `json:"last_message_id"`
	LastTemplateID string                 `json:"last_template_id"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ExpiresAt      time.Time              `json:"expires_at"`
}

func toConversationPayload(c *entities.Conversation) conversationPayload {
	return conversationPayload{
		ID:             c.ID,
		PhoneNumber:    c.PhoneNumber,
		OrderNumber:    c.OrderNumber,
		BusinessID:     c.BusinessID,
		CurrentState:   string(c.CurrentState),
		LastMessageID:  c.LastMessageID,
		LastTemplateID: c.LastTemplateID,
		Metadata:       c.Metadata,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		ExpiresAt:      c.ExpiresAt,
	}
}

// PublishConversationCreated publica evento de conversación creada
func (p *persistencePublisher) PublishConversationCreated(ctx context.Context, conversation *entities.Conversation) error {
	event := conversationEvent{
		EventType:    "conversation.created",
		Conversation: toConversationPayload(conversation),
		Timestamp:    time.Now().Unix(),
	}
	return p.publishConversation(ctx, event)
}

// PublishConversationUpdated publica evento de conversación actualizada
func (p *persistencePublisher) PublishConversationUpdated(ctx context.Context, conversation *entities.Conversation) error {
	event := conversationEvent{
		EventType:    "conversation.updated",
		Conversation: toConversationPayload(conversation),
		Timestamp:    time.Now().Unix(),
	}
	return p.publishConversation(ctx, event)
}

// PublishConversationExpired publica evento de conversación expirada
func (p *persistencePublisher) PublishConversationExpired(ctx context.Context, conversationID string) error {
	event := conversationEvent{
		EventType: "conversation.expired",
		Conversation: conversationPayload{
			ID: conversationID,
		},
		Timestamp: time.Now().Unix(),
	}
	return p.publishConversation(ctx, event)
}

func (p *persistencePublisher) publishConversation(ctx context.Context, event conversationEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando conversation event: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppConversationEvents, data); err != nil {
		p.log.Error(ctx).Err(err).
			Str("event_type", event.EventType).
			Str("conversation_id", event.Conversation.ID).
			Msg("[PersistencePublisher] Error publicando conversation event")
		return fmt.Errorf("error publicando conversation event: %w", err)
	}

	p.log.Info(ctx).
		Str("event_type", event.EventType).
		Str("conversation_id", event.Conversation.ID).
		Msg("[PersistencePublisher] Conversation event publicado")

	return nil
}

// ============================================
// MessageLog events
// ============================================

// messageLogEvent es el payload para eventos de message log
type messageLogEvent struct {
	EventType  string            `json:"event_type"`
	MessageLog messageLogPayload `json:"message_log"`
	Timestamp  int64             `json:"timestamp"`
}

type messageLogPayload struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	Direction      string     `json:"direction"`
	MessageID      string     `json:"message_id"`
	TemplateName   string     `json:"template_name"`
	Content        string     `json:"content"`
	Status         string     `json:"status"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func toMessageLogPayload(m *entities.MessageLog) messageLogPayload {
	return messageLogPayload{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		Direction:      string(m.Direction),
		MessageID:      m.MessageID,
		TemplateName:   m.TemplateName,
		Content:        m.Content,
		Status:         string(m.Status),
		DeliveredAt:    m.DeliveredAt,
		ReadAt:         m.ReadAt,
		CreatedAt:      m.CreatedAt,
	}
}

// PublishMessageLogCreated publica evento de message log creado
func (p *persistencePublisher) PublishMessageLogCreated(ctx context.Context, messageLog *entities.MessageLog) error {
	event := messageLogEvent{
		EventType:  "messagelog.created",
		MessageLog: toMessageLogPayload(messageLog),
		Timestamp:  time.Now().Unix(),
	}
	return p.publishMessageLog(ctx, event)
}

// messageStatusUpdatePayload es el payload para actualización de estado
type messageStatusUpdatePayload struct {
	MessageID  string            `json:"message_id"`
	Status     string            `json:"status"`
	Timestamps map[string]string `json:"timestamps"`
}

// PublishMessageStatusUpdated publica evento de actualización de estado de mensaje
func (p *persistencePublisher) PublishMessageStatusUpdated(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error {
	// Convertir timestamps a strings para serialización
	tsStrings := make(map[string]string, len(timestamps))
	for k, v := range timestamps {
		tsStrings[k] = v.Format(time.RFC3339)
	}

	event := struct {
		EventType string                     `json:"event_type"`
		Update    messageStatusUpdatePayload `json:"update"`
		Timestamp int64                      `json:"timestamp"`
	}{
		EventType: "messagelog.status_updated",
		Update: messageStatusUpdatePayload{
			MessageID:  messageID,
			Status:     string(status),
			Timestamps: tsStrings,
		},
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando status update event: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppMessageLogEvents, data); err != nil {
		p.log.Error(ctx).Err(err).
			Str("message_id", messageID).
			Str("status", string(status)).
			Msg("[PersistencePublisher] Error publicando message status update")
		return fmt.Errorf("error publicando message status update: %w", err)
	}

	p.log.Info(ctx).
		Str("message_id", messageID).
		Str("status", string(status)).
		Msg("[PersistencePublisher] Message status update publicado")

	return nil
}

func (p *persistencePublisher) publishMessageLog(ctx context.Context, event messageLogEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando messagelog event: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppMessageLogEvents, data); err != nil {
		p.log.Error(ctx).Err(err).
			Str("event_type", event.EventType).
			Str("message_id", event.MessageLog.MessageID).
			Msg("[PersistencePublisher] Error publicando messagelog event")
		return fmt.Errorf("error publicando messagelog event: %w", err)
	}

	p.log.Info(ctx).
		Str("event_type", event.EventType).
		Str("message_id", event.MessageLog.MessageID).
		Msg("[PersistencePublisher] MessageLog event publicado")

	return nil
}
