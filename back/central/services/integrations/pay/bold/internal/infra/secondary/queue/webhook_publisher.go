package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueuePayBoldWebhookEvents = rabbitmq.QueuePayBoldWebhookEvents

type WebhookPublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

func NewWebhookPublisher(queue rabbitmq.IQueue, logger log.ILogger) ports.IWebhookPublisher {
	return &WebhookPublisher{
		queue: queue,
		log:   logger.WithModule("bold.webhook_publisher"),
	}
}

func (p *WebhookPublisher) PublishWebhookEvent(ctx context.Context, msg *ports.BoldWebhookMessage) error {
	if p.queue == nil {
		p.log.Warn(ctx).
			Str("bold_event_id", msg.BoldEventID).
			Msg("RabbitMQ not available - webhook event not published")
		return fmt.Errorf("rabbitmq unavailable")
	}

	if err := p.queue.DeclareQueue(QueuePayBoldWebhookEvents, true); err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	type enriched struct {
		*ports.BoldWebhookMessage
		PublishedAt time.Time `json:"published_at"`
	}
	payload, err := json.Marshal(enriched{
		BoldWebhookMessage: msg,
		PublishedAt:        time.Now(),
	})
	if err != nil {
		return fmt.Errorf("marshal webhook event: %w", err)
	}

	if err := p.queue.Publish(ctx, QueuePayBoldWebhookEvents, payload); err != nil {
		return fmt.Errorf("publish webhook event: %w", err)
	}

	p.log.Info(ctx).
		Str("bold_event_id", msg.BoldEventID).
		Str("type", msg.Type).
		Msg("bold webhook event published")
	return nil
}
