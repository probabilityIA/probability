package queue

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// SSEPublisher publica eventos de envíos al dispatcher central via RabbitMQ (ExchangeEvents)
type SSEPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

// NewSSEPublisher crea un nuevo publicador SSE de envíos via RabbitMQ
func NewSSEPublisher(queue rabbitmq.IQueue, logger log.ILogger) domain.IShipmentSSEPublisher {
	return &SSEPublisher{
		queue:  queue,
		logger: logger.WithModule("shipments.sse_publisher"),
	}
}

// PublishQuoteReceived publica evento de cotización recibida
func (p *SSEPublisher) PublishQuoteReceived(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	p.publish(ctx, "shipment.quote_received", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"quotes":         data,
	})
}

// PublishQuoteFailed publica evento de cotización fallida
func (p *SSEPublisher) PublishQuoteFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.quote_failed", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

// PublishGuideGenerated publica evento de guía generada exitosamente
func (p *SSEPublisher) PublishGuideGenerated(ctx context.Context, businessID uint, shipmentID uint, correlationID string, trackingNumber string, labelURL string, carrier string, notification *domain.GuideNotificationData) {
	data := map[string]interface{}{
		"shipment_id":     shipmentID,
		"correlation_id":  correlationID,
		"tracking_number": trackingNumber,
		"label_url":       labelURL,
		"carrier":         carrier,
	}

	// Enrich with notification data for WhatsApp routing
	if notification != nil {
		data["customer_name"] = notification.CustomerName
		data["customer_phone"] = notification.CustomerPhone
		data["order_number"] = notification.OrderNumber
		data["business_name"] = notification.BusinessName
		data["integration_id"] = notification.IntegrationID
	}

	p.publish(ctx, "shipment.guide_generated", businessID, data)
}

// PublishGuideFailed publica evento de generación de guía fallida
func (p *SSEPublisher) PublishGuideFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.guide_failed", businessID, map[string]interface{}{
		"shipment_id":    shipmentID,
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

// PublishTrackingUpdated publica evento de tracking actualizado
func (p *SSEPublisher) PublishTrackingUpdated(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	p.publish(ctx, "shipment.tracking_updated", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"tracking":       data,
	})
}

// PublishTrackingFailed publica evento de tracking fallido
func (p *SSEPublisher) PublishTrackingFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.tracking_failed", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

// PublishShipmentCancelled publica evento de envío cancelado
func (p *SSEPublisher) PublishShipmentCancelled(ctx context.Context, businessID uint, shipmentID uint) {
	p.publish(ctx, "shipment.cancelled", businessID, map[string]interface{}{
		"shipment_id": shipmentID,
	})
}

// PublishCancelFailed publica evento de cancelación fallida
func (p *SSEPublisher) PublishCancelFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.cancel_failed", businessID, map[string]interface{}{
		"shipment_id":    shipmentID,
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

// publish publica un evento al dispatcher central de forma no-bloqueante
func (p *SSEPublisher) publish(ctx context.Context, eventType string, businessID uint, data map[string]interface{}) {
	go func() {
		if err := rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
			Type:       eventType,
			Category:   "shipment",
			BusinessID: businessID,
			Data:       data,
		}); err != nil {
			p.logger.Error(ctx).
				Err(err).
				Str("event_type", eventType).
				Uint("business_id", businessID).
				Msg("Error publishing shipment SSE event to RabbitMQ")
		}
	}()
}

// noopSSEPublisher es una implementación no-op para cuando RabbitMQ no está disponible
type noopSSEPublisher struct{}

// NewNoopSSEPublisher crea un publisher que no hace nada
func NewNoopSSEPublisher() domain.IShipmentSSEPublisher {
	return &noopSSEPublisher{}
}

func (n *noopSSEPublisher) PublishQuoteReceived(_ context.Context, _ uint, _ string, _ map[string]interface{}) {
}
func (n *noopSSEPublisher) PublishQuoteFailed(_ context.Context, _ uint, _ string, _ string) {}
func (n *noopSSEPublisher) PublishGuideGenerated(_ context.Context, _ uint, _ uint, _ string, _ string, _ string, _ string, _ *domain.GuideNotificationData) {
}
func (n *noopSSEPublisher) PublishGuideFailed(_ context.Context, _ uint, _ uint, _ string, _ string) {
}
func (n *noopSSEPublisher) PublishTrackingUpdated(_ context.Context, _ uint, _ string, _ map[string]interface{}) {
}
func (n *noopSSEPublisher) PublishTrackingFailed(_ context.Context, _ uint, _ string, _ string) {}
func (n *noopSSEPublisher) PublishShipmentCancelled(_ context.Context, _ uint, _ uint)          {}
func (n *noopSSEPublisher) PublishCancelFailed(_ context.Context, _ uint, _ uint, _ string, _ string) {
}

// Compile-time interface checks
var _ domain.IShipmentSSEPublisher = (*SSEPublisher)(nil)
var _ domain.IShipmentSSEPublisher = (*noopSSEPublisher)(nil)
