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

const whatsAppConfirmationQueue = rabbitmq.QueueOrdersConfirmationRequested

type channelPublisher struct {
	rabbitMQ rabbitmq.IQueue
	logger   log.ILogger
}

func New(rabbitMQ rabbitmq.IQueue, logger log.ILogger) ports.IChannelPublisher {
	return &channelPublisher{
		rabbitMQ: rabbitMQ,
		logger:   logger,
	}
}

func (p *channelPublisher) PublishToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	switch event.Category {
	case "shipment":
		return p.publishShipmentToWhatsApp(ctx, event, config)
	default:
		return p.publishOrderToWhatsApp(ctx, event, config)
	}
}

const paymentMethodIDCOD uint = 6

func isCODEvent(event entities.Event) bool {
	val, ok := event.Data["payment_method_id"]
	if !ok || val == nil {
		return false
	}
	switch v := val.(type) {
	case uint:
		return v == paymentMethodIDCOD
	case int:
		return uint(v) == paymentMethodIDCOD
	case int64:
		return uint(v) == paymentMethodIDCOD
	case float64:
		return uint(v) == paymentMethodIDCOD
	default:
		return false
	}
}

func eventCodeToTemplateName(eventCode string, isCOD bool) string {
	switch eventCode {
	case "order.shipped":
		if isCOD {
			return "pedido_en_reparto_cod"
		}
		return "pedido_en_reparto"
	case "order.delivered":
		if isCOD {
			return "pedido_entregado_cod"
		}
		return "pedido_entregado"
	default:
		if isCOD {
			return "confirmacion_pedido_contraentrega"
		}
		return "confirmacion_pedido"
	}
}

func (p *channelPublisher) publishOrderToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	templateName := eventCodeToTemplateName(config.EventCode, isCODEvent(event))

	payload := map[string]any{
		"event_type":        "order.confirmation_requested",
		"business_id":       event.BusinessID,
		"integration_id":    event.IntegrationID,
		"config_id":         config.ID,
		"notification_type": "whatsapp",
		"template_name":     templateName,
	}

	dataFields := []string{
		"order_id", "order_number", "internal_number", "external_id",
		"customer_name", "customer_phone", "customer_email",
		"total_amount", "currency", "platform",
		"items_summary", "shipping_address", "shipping_street", "shipping_city", "shipping_state",
		"business_name", "payment_method_id", "payment_method_name", "tracking_number", "carrier",
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

	p.logger.Info(ctx).Str("event_id", event.ID).Str("event_type", event.Type).Uint("config_id", config.ID).Str("template_name", templateName).Str("queue", whatsAppConfirmationQueue).Msg("Evento encolado para WhatsApp")
	return nil
}

func (p *channelPublisher) publishShipmentToWhatsApp(ctx context.Context, event entities.Event, config entities.CachedNotificationConfig) error {
	payload := map[string]any{
		"event_type":        event.Type,
		"business_id":       event.BusinessID,
		"integration_id":    event.IntegrationID,
		"config_id":         config.ID,
		"notification_type": "whatsapp",
	}

	dataFields := []string{
		"shipment_id", "tracking_number", "label_url", "carrier",
		"customer_name", "customer_phone", "order_number", "business_name",
		"correlation_id", "total_amount", "cod_total", "tracking_url",
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

	p.logger.Info(ctx).Str("event_id", event.ID).Str("event_type", event.Type).Uint("config_id", config.ID).Str("queue", targetQueue).Msg("Evento de envio encolado para WhatsApp")
	return nil
}
