package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueuePayBancolombiaWebhookEvents = rabbitmq.QueuePayBancolombiaWebhookEvents

type WebhookPublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

func NewWebhookPublisher(queue rabbitmq.IQueue, logger log.ILogger) ports.IWebhookPublisher {
	return &WebhookPublisher{
		queue: queue,
		log:   logger.WithModule("bancolombia.webhook_publisher"),
	}
}

func (p *WebhookPublisher) PublishWebhookEvent(ctx context.Context, msg *ports.BancolombiaWebhookMessage) error {
	if p.queue == nil {
		p.log.Warn(ctx).
			Str("event_id", msg.EventID).
			Msg("RabbitMQ not available - webhook event not published")
		return fmt.Errorf("rabbitmq unavailable")
	}

	if err := p.queue.DeclareQueue(QueuePayBancolombiaWebhookEvents, true); err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	type enriched struct {
		*ports.BancolombiaWebhookMessage
		PublishedAt time.Time `json:"published_at"`
	}
	payload, err := json.Marshal(enriched{
		BancolombiaWebhookMessage: msg,
		PublishedAt:               time.Now(),
	})
	if err != nil {
		return fmt.Errorf("marshal webhook event: %w", err)
	}

	if err := p.queue.Publish(ctx, QueuePayBancolombiaWebhookEvents, payload); err != nil {
		return fmt.Errorf("publish webhook event: %w", err)
	}

	p.log.Info(ctx).
		Str("event_id", msg.EventID).
		Str("type", msg.Type).
		Msg("bancolombia webhook event published")
	return nil
}
