package channel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// whatsAppConfirmationQueue usa la constante centralizada.
const whatsAppConfirmationQueue = rabbitmq.QueueOrdersConfirmationRequested

type channelPublisher struct {
	rabbitMQ rabbitmq.IQueue
	logger   log.ILogger
}

// New crea un nuevo ChannelPublisher para publicar a canales destino (WhatsApp, etc.)
func New(rabbitMQ rabbitmq.IQueue, logger log.ILogger) ports.IChannelPublisher {
	return &channelPublisher{
		rabbitMQ: rabbitMQ,
		logger:   logger,
	}
}

// PublishToWhatsApp publica un evento a la queue de confirmación de WhatsApp
func (p *channelPublisher) PublishToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	// Construir payload compatible con el consumer de WhatsApp OrderConfirmation
	payload := map[string]interface{}{
		"event_type":       "order.confirmation_requested",
		"business_id":      event.BusinessID,
		"integration_id":   event.IntegrationID,
		"config_id":        config.ID,
		"notification_type": "whatsapp",
	}

	// Copiar datos relevantes del evento
	if orderID, ok := event.Data["order_id"]; ok {
		payload["order_id"] = orderID
	}
	if event.Metadata != nil {
		if orderID, ok := event.Metadata["order_id"]; ok {
			payload["order_id"] = orderID
		}
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Msg("Error serializando payload para WhatsApp queue")
		return fmt.Errorf("error serializando payload WhatsApp: %w", err)
	}

	if err := p.rabbitMQ.Publish(ctx, whatsAppConfirmationQueue, jsonBytes); err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_id", event.ID).
			Str("queue", whatsAppConfirmationQueue).
			Msg("Error publicando a WhatsApp queue")
		return fmt.Errorf("error publicando a WhatsApp queue: %w", err)
	}

	p.logger.Info(ctx).
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Uint("config_id", config.ID).
		Str("queue", whatsAppConfirmationQueue).
		Msg("Evento encolado para WhatsApp")

	return nil
}
