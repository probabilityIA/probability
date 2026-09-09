package channel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/events/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

var pushDataFields = []string{
	"shipment_id", "tracking_number", "carrier", "order_id", "order_number",
	"customer_name", "business_name", "status_code", "status_name",
	"total_amount", "cod_total", "balance", "threshold",
}

func (p *channelPublisher) PublishToPush(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	payload := map[string]any{
		"event_type":        event.Type,
		"event_category":    event.Category,
		"business_id":       event.BusinessID,
		"integration_id":    event.IntegrationID,
		"config_id":         config.ID,
		"notification_type": "push",
	}

	for _, field := range pushDataFields {
		if val, ok := event.Data[field]; ok && val != nil && val != "" {
			payload[field] = val
		}
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error(ctx).Err(err).Str("event_id", event.ID).Msg("Error serializando payload para push")
		return fmt.Errorf("%w: push payload: %v", domainerrors.ErrSerializeFailed, err)
	}

	targetQueue := rabbitmq.QueuePushNotificationRequests
	if err := p.rabbitMQ.Publish(ctx, targetQueue, jsonBytes); err != nil {
		p.logger.Error(ctx).Err(err).Str("event_id", event.ID).Str("queue", targetQueue).Msg("Error publicando a la cola de push")
		return fmt.Errorf("%w: push queue: %v", domainerrors.ErrPublishFailed, err)
	}

	p.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Uint("config_id", config.ID).
		Str("queue", targetQueue).
		Msg("Evento encolado para push")
	return nil
}
