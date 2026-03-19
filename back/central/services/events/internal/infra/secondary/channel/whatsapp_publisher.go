package channel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/events/internal/domain/errors"
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

// PublishToWhatsApp publica un evento a la queue de WhatsApp apropiada según la categoría
func (p *channelPublisher) PublishToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	// Route by category
	switch event.Category {
	case "shipment":
		return p.publishShipmentToWhatsApp(ctx, event, config)
	default:
		return p.publishOrderToWhatsApp(ctx, event, config)
	}
}

// publishOrderToWhatsApp publica un evento de orden a la queue de confirmación de WhatsApp
func (p *channelPublisher) publishOrderToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	payload := map[string]interface{}{
		"event_type":        "order.confirmation_requested",
		"business_id":       event.BusinessID,
		"integration_id":    event.IntegrationID,
		"config_id":         config.ID,
		"notification_type": "whatsapp",
	}

	dataFields := []string{
		"order_id", "order_number", "internal_number", "external_id",
		"customer_name", "customer_phone", "customer_email",
		"total_amount", "currency", "platform",
		"items_summary", "shipping_address", "business_name",
	}
	for _, field := range dataFields {
		if val, ok := event.Data[field]; ok && val != nil && val != "" {
			payload[field] = val
		}
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error(ctx).Err(err).Str("event_id", event.ID).Msg("Error serializando payload para WhatsApp queue")
		return fmt.Errorf("%w: WhatsApp payload: %v", domainerrors.ErrSerializeFailed, err)
	}

	if err := p.rabbitMQ.Publish(ctx, whatsAppConfirmationQueue, jsonBytes); err != nil {
		p.logger.Error(ctx).Err(err).Str("event_id", event.ID).Str("queue", whatsAppConfirmationQueue).Msg("Error publicando a WhatsApp queue")
		return fmt.Errorf("%w: WhatsApp queue: %v", domainerrors.ErrPublishFailed, err)
	}

	p.logger.Info(ctx).Str("event_id", event.ID).Str("event_type", event.Type).Uint("config_id", config.ID).Str("queue", whatsAppConfirmationQueue).Msg("Evento encolado para WhatsApp")
	return nil
}

// publishShipmentToWhatsApp publica un evento de envío a la queue de guía generada para WhatsApp
func (p *channelPublisher) publishShipmentToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	payload := map[string]interface{}{
		"event_type":        event.Type,
		"business_id":       event.BusinessID,
		"integration_id":    event.IntegrationID,
		"config_id":         config.ID,
		"notification_type": "whatsapp",
	}

	dataFields := []string{
		"shipment_id", "tracking_number", "label_url", "carrier",
		"customer_name", "customer_phone", "order_number", "business_name",
		"correlation_id",
	}
	for _, field := range dataFields {
		if val, ok := event.Data[field]; ok && val != nil && val != "" {
			payload[field] = val
		}
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error(ctx).Err(err).Str("event_id", event.ID).Msg("Error serializando payload para WhatsApp shipment queue")
		return fmt.Errorf("%w: WhatsApp shipment payload: %v", domainerrors.ErrSerializeFailed, err)
	}

	targetQueue := rabbitmq.QueueShipmentsWhatsAppGuideNotification
	if err := p.rabbitMQ.Publish(ctx, targetQueue, jsonBytes); err != nil {
		p.logger.Error(ctx).Err(err).Str("event_id", event.ID).Str("queue", targetQueue).Msg("Error publicando a WhatsApp shipment queue")
		return fmt.Errorf("%w: WhatsApp shipment queue: %v", domainerrors.ErrPublishFailed, err)
	}

	p.logger.Info(ctx).Str("event_id", event.ID).Str("event_type", event.Type).Uint("config_id", config.ID).Str("queue", targetQueue).Msg("Evento de envío encolado para WhatsApp")
	return nil
}
