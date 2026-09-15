package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type conversationEvent struct {
	Type      string          `json:"type"`
	Message   *messagePayload `json:"message,omitempty"`
	UserID    uint            `json:"user_id,omitempty"`
	MessageID string          `json:"message_id,omitempty"`
	Value     int             `json:"value"`
	At        time.Time       `json:"at"`
}

type messagePayload struct {
	ID               string    `json:"id"`
	ConversationID   string    `json:"conversation_id"`
	BusinessID       *uint     `json:"business_id"`
	UserID           uint      `json:"user_id"`
	Pathname         string    `json:"pathname"`
	Question         string    `json:"question"`
	Answer           string    `json:"answer"`
	DestinationKey   string    `json:"destination_key"`
	DestinationRoute string    `json:"destination_route"`
	ErrorCode        string    `json:"error_code"`
	Model            string    `json:"model"`
	InputTokens      int       `json:"input_tokens"`
	OutputTokens     int       `json:"output_tokens"`
	LatencyMs        int       `json:"latency_ms"`
	CreatedAt        time.Time `json:"created_at"`
}

func (r *Recorder) RecordMessage(ctx context.Context, record entities.MessageRecord) error {
	return r.publish(ctx, conversationEvent{
		Type: "message",
		Message: &messagePayload{
			ID:               record.ID,
			ConversationID:   record.ConversationID,
			BusinessID:       record.BusinessID,
			UserID:           record.UserID,
			Pathname:         record.Pathname,
			Question:         record.Question,
			Answer:           record.Answer,
			DestinationKey:   record.DestinationKey,
			DestinationRoute: record.DestinationRoute,
			ErrorCode:        record.ErrorCode,
			Model:            record.Model,
			InputTokens:      record.InputTokens,
			OutputTokens:     record.OutputTokens,
			LatencyMs:        record.LatencyMs,
			CreatedAt:        record.CreatedAt,
		},
		At: record.CreatedAt,
	})
}

func (r *Recorder) RecordFeedback(ctx context.Context, userID uint, messageID string, value int, at time.Time) error {
	return r.publish(ctx, conversationEvent{Type: "feedback", UserID: userID, MessageID: messageID, Value: value, At: at})
}

func (r *Recorder) RecordClick(ctx context.Context, userID uint, messageID string, at time.Time) error {
	return r.publish(ctx, conversationEvent{Type: "click", UserID: userID, MessageID: messageID, At: at})
}

func (r *Recorder) publish(ctx context.Context, event conversationEvent) error {
	r.declared.Do(func() {
		_ = r.queue.DeclareQueue(rabbitmq.QueueAIAssistantConversations, true)
	})

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("serializar evento del asistente: %w", err)
	}
	return r.queue.Publish(ctx, rabbitmq.QueueAIAssistantConversations, body)
}
