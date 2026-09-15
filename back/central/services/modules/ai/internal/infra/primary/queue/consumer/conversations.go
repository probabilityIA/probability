package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type conversationEvent struct {
	Type      string          `json:"type"`
	Message   *messagePayload `json:"message"`
	UserID    uint            `json:"user_id"`
	MessageID string          `json:"message_id"`
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

func (c *Consumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueueAIAssistantConversations
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		return fmt.Errorf("declarar cola %s: %w", queueName, err)
	}
	if err := c.queue.Consume(ctx, queueName, c.handle); err != nil {
		return fmt.Errorf("consumir cola %s: %w", queueName, err)
	}
	c.log.Info(ctx).Str("queue", queueName).Msg("[ai.assistant] consumidor de conversaciones iniciado")
	return nil
}

func (c *Consumer) handle(body []byte) error {
	ctx := context.Background()

	var event conversationEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.log.Warn(ctx).Err(err).Msg("[ai.assistant] evento de conversacion malformado (ACK)")
		return nil
	}

	var err error
	switch event.Type {
	case "message":
		if event.Message == nil || event.Message.ID == "" || event.Message.UserID == 0 {
			c.log.Warn(ctx).Msg("[ai.assistant] mensaje sin datos minimos (ACK)")
			return nil
		}
		err = c.uc.PersistMessage(ctx, toRecord(*event.Message))
	case "feedback":
		err = c.uc.PersistFeedback(ctx, event.UserID, event.MessageID, event.Value, event.At)
	case "click":
		err = c.uc.PersistClick(ctx, event.UserID, event.MessageID, event.At)
	default:
		c.log.Warn(ctx).Str("type", event.Type).Msg("[ai.assistant] tipo de evento desconocido (ACK)")
		return nil
	}

	if err == nil {
		return nil
	}
	if isPermanent(err) {
		c.log.Warn(ctx).Err(err).Str("type", event.Type).Msg("[ai.assistant] evento descartado por error permanente (ACK)")
		return nil
	}
	if errors.Is(err, domainerrors.ErrMessageNotFound) {
		c.log.Warn(ctx).Str("message_id", event.MessageID).Msg("[ai.assistant] el mensaje aun no esta guardado; se reintenta")
		return err
	}
	c.log.Error(ctx).Err(err).Str("type", event.Type).Msg("[ai.assistant] error guardando evento; se reintenta")
	return err
}

func isPermanent(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	switch pgErr.Code {
	case "23503", "23502", "22P02", "22001":
		return true
	}
	return false
}

func toRecord(p messagePayload) entities.MessageRecord {
	return entities.MessageRecord{
		ID:               p.ID,
		ConversationID:   p.ConversationID,
		BusinessID:       p.BusinessID,
		UserID:           p.UserID,
		Pathname:         p.Pathname,
		Question:         p.Question,
		Answer:           p.Answer,
		DestinationKey:   p.DestinationKey,
		DestinationRoute: p.DestinationRoute,
		ErrorCode:        p.ErrorCode,
		Model:            p.Model,
		InputTokens:      p.InputTokens,
		OutputTokens:     p.OutputTokens,
		LatencyMs:        p.LatencyMs,
		CreatedAt:        p.CreatedAt,
	}
}
