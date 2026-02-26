package redis

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// paySSEEvent estructura del evento publicado a Redis
type paySSEEvent struct {
	ID         string                 `json:"id"`
	EventType  string                 `json:"event_type"`
	BusinessID uint                   `json:"business_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}

// SSEPublisher publica eventos de pago a Redis Pub/Sub
type SSEPublisher struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
}

// NewSSEPublisher crea un nuevo publisher SSE de pagos
func NewSSEPublisher(redisClient redisclient.IRedis, logger log.ILogger, channel string) ports.ISSEPublisher {
	return &SSEPublisher{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
	}
}

func (p *SSEPublisher) PublishPaymentCompleted(ctx context.Context, tx *entities.PaymentTransaction) error {
	return p.publish(ctx, paySSEEvent{
		ID:         generateEventID(),
		EventType:  "pay.completed",
		BusinessID: tx.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"transaction_id": tx.ID,
			"reference":      tx.Reference,
			"amount":         tx.Amount,
			"currency":       tx.Currency,
			"gateway":        tx.GatewayCode,
			"external_id":    ptrStr(tx.ExternalID),
			"status":         tx.Status,
		},
	})
}

func (p *SSEPublisher) PublishPaymentFailed(ctx context.Context, tx *entities.PaymentTransaction, errMsg string) error {
	return p.publish(ctx, paySSEEvent{
		ID:         generateEventID(),
		EventType:  "pay.failed",
		BusinessID: tx.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"transaction_id": tx.ID,
			"reference":      tx.Reference,
			"amount":         tx.Amount,
			"currency":       tx.Currency,
			"gateway":        tx.GatewayCode,
			"status":         tx.Status,
			"error":          errMsg,
		},
	})
}

func (p *SSEPublisher) PublishPaymentProcessing(ctx context.Context, tx *entities.PaymentTransaction) error {
	return p.publish(ctx, paySSEEvent{
		ID:         generateEventID(),
		EventType:  "pay.processing",
		BusinessID: tx.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"transaction_id": tx.ID,
			"reference":      tx.Reference,
			"amount":         tx.Amount,
			"currency":       tx.Currency,
			"gateway":        tx.GatewayCode,
			"status":         tx.Status,
		},
	})
}

func (p *SSEPublisher) publish(ctx context.Context, event paySSEEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		p.logger.Error(ctx).Err(err).Str("event_type", event.EventType).Msg("Error serializing pay SSE event")
		return err
	}

	go func() {
		publishCtx := context.Background()
		if pubErr := p.redisClient.Client(publishCtx).Publish(publishCtx, p.channel, data).Err(); pubErr != nil {
			p.logger.Error(publishCtx).Err(pubErr).Str("event_type", event.EventType).Str("channel", p.channel).Msg("Error publishing pay SSE event")
		}
	}()

	return nil
}

func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// noopSSEPublisher es una implementación no-op
type noopSSEPublisher struct{}

// NewNoopSSEPublisher crea un publisher SSE que no hace nada
func NewNoopSSEPublisher() ports.ISSEPublisher {
	return &noopSSEPublisher{}
}

func (n *noopSSEPublisher) PublishPaymentCompleted(_ context.Context, _ *entities.PaymentTransaction) error {
	return nil
}
func (n *noopSSEPublisher) PublishPaymentFailed(_ context.Context, _ *entities.PaymentTransaction, _ string) error {
	return nil
}
func (n *noopSSEPublisher) PublishPaymentProcessing(_ context.Context, _ *entities.PaymentTransaction) error {
	return nil
}
