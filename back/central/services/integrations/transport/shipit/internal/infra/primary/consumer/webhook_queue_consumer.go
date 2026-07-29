package consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type WebhookQueueMessage struct {
	URL      string              `json:"url"`
	Method   string              `json:"method"`
	RemoteIP string              `json:"remote_ip"`
	Headers  map[string][]string `json:"headers"`
	RawBody  json.RawMessage     `json:"raw_body"`
}

type WebhookQueueConsumer struct {
	queue rabbitmq.IQueue
	uc    app.IWebhookUseCase
	log   log.ILogger
}

func NewWebhookQueueConsumer(queue rabbitmq.IQueue, uc app.IWebhookUseCase, logger log.ILogger) *WebhookQueueConsumer {
	return &WebhookQueueConsumer{
		queue: queue,
		uc:    uc,
		log:   logger.WithModule("transport.shipit.webhook_consumer"),
	}
}

func (c *WebhookQueueConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueWebhooksShipitReceived, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al declarar la cola de webhooks Shipit")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueWebhooksShipitReceived, func(body []byte) error {
			var msg WebhookQueueMessage
			if err := json.Unmarshal(body, &msg); err != nil {
				c.log.Error(ctx).Err(err).Msg("Mensaje de webhook Shipit invalido, descartado")
				return nil
			}

			var payload domain.WebhookPayload
			if err := json.Unmarshal(msg.RawBody, &payload); err != nil {
				c.log.Warn(ctx).Err(err).Msg("Payload de webhook Shipit invalido, descartado")
				return nil
			}

			bg := context.Background()
			result, err := c.uc.Process(bg, app.WebhookRequest{
				URL:      msg.URL,
				Method:   msg.Method,
				RemoteIP: msg.RemoteIP,
				Headers:  msg.Headers,
				RawBody:  msg.RawBody,
				Payload:  payload,
			})
			if err != nil {
				c.log.Error(bg).Err(err).Str("tracking_number", payload.TrackingNumber).Msg("Webhook processing failed")
				return nil
			}

			if result.IsIgnored {
				c.log.Info(bg).Str("tracking_number", payload.TrackingNumber).Str("reason", result.IgnoredReason).Msg("Webhook ignorado")
				return nil
			}

			c.log.Info(bg).
				Str("tracking_number", result.TrackingNumber).
				Str("probability_status", result.ProbabilityStatus.String()).
				Str("correlation_id", result.CorrelationID).
				Msg("Webhook processed")
			return nil
		})
		if err != nil {
			c.log.Error(ctx).Err(err).Msg("Error al consumir la cola de webhooks Shipit")
		}
	}()

	c.log.Info(ctx).Msg("Consumer de webhooks Shipit iniciado")
}
