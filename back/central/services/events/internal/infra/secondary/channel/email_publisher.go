package channel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// PublishToEmail publica un evento a la queue de notificaciones por email
func (p *channelPublisher) PublishToEmail(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	// Extraer customer_email de los datos del evento
	customerEmail := ""
	if email, ok := event.Data["customer_email"]; ok {
		if emailStr, ok := email.(string); ok {
			customerEmail = emailStr
		}
	}

	if customerEmail == "" {
		p.logger.Warn(ctx).
			Str("event_id", event.ID).
			Uint("config_id", config.ID).
			Msg("No se encontró customer_email en el evento, saltando publicación a Email")
		return nil
	}

	payload := map[string]interface{}{
		"event_type":     event.Type,
		"business_id":    event.BusinessID,
		"integration_id": event.IntegrationID,
		"config_id":      config.ID,
		"customer_email": customerEmail,
		"event_data":     event.Data,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Msg("Error serializando payload para Email queue")
		return fmt.Errorf("error serializando payload Email: %w", err)
	}

	if err := p.rabbitMQ.Publish(ctx, rabbitmq.QueueMessagingEmailRequests, jsonBytes); err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Str("queue", rabbitmq.QueueMessagingEmailRequests).
			Msg("Error publicando a Email queue")
		return fmt.Errorf("error publicando a Email queue: %w", err)
	}

	p.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Str("to", customerEmail).
		Uint("config_id", config.ID).
		Str("queue", rabbitmq.QueueMessagingEmailRequests).
		Msg("Evento encolado para Email")

	return nil
}
