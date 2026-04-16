package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type persistencePublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewPersistencePublisher crea un publisher que guarda sesiones AI en BD via las queues de WhatsApp
func NewPersistencePublisher(rabbit rabbitmq.IQueue, logger log.ILogger) domain.IAIPersistencePublisher {
	return &persistencePublisher{
		rabbit: rabbit,
		log:    logger.WithModule("ai-sales-persistence"),
	}
}

// PublishConversationUpsert publica un evento de conversación AI para persistencia en BD.
// Reutiliza la queue whatsapp.conversation.events que ya tiene consumer y modelo de BD.
func (p *persistencePublisher) PublishConversationUpsert(ctx context.Context, session *domain.AISession) error {
	event := map[string]any{
		"event_type": "conversation.upsert",
		"conversation": map[string]any{
			"id":            session.ID,
			"phone_number":  session.PhoneNumber,
			"order_number":  "",
			"business_id":   session.BusinessID,
			"current_state": "AI_ACTIVE",
			"metadata": map[string]any{
				"source":        "ai_sales",
				"message_count": len(session.Messages),
			},
			"created_at": session.CreatedAt,
			"updated_at": session.UpdatedAt,
			"expires_at": session.ExpiresAt,
		},
		"timestamp": time.Now().Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializing conversation upsert: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppConversationEvents, data); err != nil {
		p.log.Error(ctx).Err(err).
			Str("session_id", session.ID).
			Str("phone", session.PhoneNumber).
			Msg("Error publicando conversation upsert")
		return err
	}

	p.log.Info(ctx).
		Str("session_id", session.ID).
		Str("phone", session.PhoneNumber).
		Msg("Conversation AI publicada para persistencia")

	return nil
}

// PublishMessageLog publica un mensaje AI para persistencia en whatsapp_message_logs.
func (p *persistencePublisher) PublishMessageLog(ctx context.Context, conversationID, phoneNumber, direction, content string) error {
	event := map[string]any{
		"event_type": "messagelog.created",
		"message_log": map[string]any{
			"id":              fmt.Sprintf("ai-%s-%d", phoneNumber, time.Now().UnixNano()),
			"conversation_id": conversationID,
			"direction":       direction,
			"message_id":      "",
			"template_name":   "",
			"content":         content,
			"status":          "delivered",
			"created_at":      time.Now(),
		},
		"timestamp": time.Now().Unix(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializing message log: %w", err)
	}

	if err := p.rabbit.Publish(ctx, rabbitmq.QueueWhatsAppMessageLogEvents, data); err != nil {
		p.log.Error(ctx).Err(err).
			Str("conversation_id", conversationID).
			Str("direction", direction).
			Msg("Error publicando message log AI")
		return err
	}

	return nil
}
