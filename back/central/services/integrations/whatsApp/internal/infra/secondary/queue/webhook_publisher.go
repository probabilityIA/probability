package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// WebhookPublisher implementa la interfaz IEventPublisher
type WebhookPublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewWebhookPublisher crea una nueva instancia del publicador
func NewWebhookPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) domain.IEventPublisher {
	return &WebhookPublisher{
		rabbit: rabbit,
		log:    logger.WithModule("whatsapp-publisher"),
	}
}

// PublishOrderConfirmed publica un evento cuando un pedido es confirmado
func (p *WebhookPublisher) PublishOrderConfirmed(ctx context.Context, orderNumber, phoneNumber string, businessID uint) error {
	p.log.Info(ctx).
		Str("order_number", orderNumber).
		Str("phone_number", phoneNumber).
		Uint("business_id", businessID).
		Msg("[WhatsApp Publisher] - publicando evento order.confirmed")

	event := map[string]interface{}{
		"event_type":   "order.confirmed",
		"order_number": orderNumber,
		"phone_number": phoneNumber,
		"business_id":  businessID,
		"source":       "whatsapp",
		"timestamp":    time.Now().Unix(),
	}

	return p.publish(ctx, "orders.whatsapp.confirmed", event)
}

// PublishOrderCancelled publica un evento cuando un pedido es cancelado
func (p *WebhookPublisher) PublishOrderCancelled(ctx context.Context, orderNumber, reason, phoneNumber string, businessID uint) error {
	p.log.Info(ctx).
		Str("order_number", orderNumber).
		Str("reason", reason).
		Str("phone_number", phoneNumber).
		Uint("business_id", businessID).
		Msg("[WhatsApp Publisher] - publicando evento order.cancelled")

	event := map[string]interface{}{
		"event_type":          "order.cancelled",
		"order_number":        orderNumber,
		"cancellation_reason": reason,
		"phone_number":        phoneNumber,
		"business_id":         businessID,
		"source":              "whatsapp",
		"timestamp":           time.Now().Unix(),
	}

	return p.publish(ctx, "orders.whatsapp.cancelled", event)
}

// PublishNoveltyRequested publica un evento cuando se solicita una novedad
func (p *WebhookPublisher) PublishNoveltyRequested(ctx context.Context, orderNumber, noveltyType, phoneNumber string, businessID uint) error {
	p.log.Info(ctx).
		Str("order_number", orderNumber).
		Str("novelty_type", noveltyType).
		Str("phone_number", phoneNumber).
		Uint("business_id", businessID).
		Msg("[WhatsApp Publisher] - publicando evento order.novelty_requested")

	event := map[string]interface{}{
		"event_type":   "order.novelty_requested",
		"order_number": orderNumber,
		"novelty_type": noveltyType, // "change_address", "change_products", "change_payment"
		"phone_number": phoneNumber,
		"business_id":  businessID,
		"source":       "whatsapp",
		"timestamp":    time.Now().Unix(),
	}

	return p.publish(ctx, "orders.whatsapp.novelty", event)
}

// PublishHandoffRequested publica un evento cuando se solicita atención humana
func (p *WebhookPublisher) PublishHandoffRequested(ctx context.Context, orderNumber, phoneNumber string, businessID uint, conversationID string) error {
	p.log.Info(ctx).
		Str("order_number", orderNumber).
		Str("phone_number", phoneNumber).
		Uint("business_id", businessID).
		Str("conversation_id", conversationID).
		Msg("[WhatsApp Publisher] - publicando evento customer.handoff_requested")

	event := map[string]interface{}{
		"event_type":      "customer.handoff_requested",
		"order_number":    orderNumber,
		"phone_number":    phoneNumber,
		"business_id":     businessID,
		"conversation_id": conversationID,
		"source":          "whatsapp",
		"timestamp":       time.Now().Unix(),
		"status":          "pending_human_agent",
	}

	return p.publish(ctx, "customer.whatsapp.handoff", event)
}

// publish es el método interno que serializa y publica el evento
func (p *WebhookPublisher) publish(ctx context.Context, queueName string, event map[string]interface{}) error {
	// Serializar evento a JSON
	payload, err := json.Marshal(event)
	if err != nil {
		p.log.Error(ctx).Err(err).
			Str("queue", queueName).
			Msg("[WhatsApp Publisher] - error serializando evento")
		return fmt.Errorf("error al serializar evento: %w", err)
	}

	// Publicar en RabbitMQ
	if err := p.rabbit.Publish(ctx, queueName, payload); err != nil {
		p.log.Error(ctx).Err(err).
			Str("queue", queueName).
			Msg("[WhatsApp Publisher] - error publicando evento")
		return fmt.Errorf("error al publicar evento en RabbitMQ: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", queueName).
		Str("event_type", event["event_type"].(string)).
		Msg("[WhatsApp Publisher] - evento publicado exitosamente")

	return nil
}
