package redis

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// shipmentSSEEvent es la estructura del evento publicado a Redis
type shipmentSSEEvent struct {
	ID         string                 `json:"id"`
	EventType  string                 `json:"event_type"`
	BusinessID uint                   `json:"business_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}

// SSEPublisher publica eventos de envíos a Redis Pub/Sub para SSE
type SSEPublisher struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
}

// NewSSEPublisher crea un nuevo publicador SSE de envíos
func NewSSEPublisher(redisClient redisclient.IRedis, logger log.ILogger, channel string) domain.IShipmentSSEPublisher {
	return &SSEPublisher{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
	}
}

// PublishQuoteReceived publica evento de cotización recibida
func (p *SSEPublisher) PublishQuoteReceived(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.quote_received",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"correlation_id": correlationID,
			"quotes":         data,
		},
	}
	p.publish(ctx, event)
}

// PublishQuoteFailed publica evento de cotización fallida
func (p *SSEPublisher) PublishQuoteFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.quote_failed",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"correlation_id": correlationID,
			"error_message":  errorMsg,
		},
	}
	p.publish(ctx, event)
}

// PublishGuideGenerated publica evento de guía generada exitosamente
func (p *SSEPublisher) PublishGuideGenerated(ctx context.Context, businessID uint, shipmentID uint, trackingNumber string, labelURL string) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.guide_generated",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"shipment_id":     shipmentID,
			"tracking_number": trackingNumber,
			"label_url":       labelURL,
		},
	}
	p.publish(ctx, event)
}

// PublishGuideFailed publica evento de generación de guía fallida
func (p *SSEPublisher) PublishGuideFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.guide_failed",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"shipment_id":    shipmentID,
			"correlation_id": correlationID,
			"error_message":  errorMsg,
		},
	}
	p.publish(ctx, event)
}

// PublishTrackingUpdated publica evento de tracking actualizado
func (p *SSEPublisher) PublishTrackingUpdated(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.tracking_updated",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"correlation_id": correlationID,
			"tracking":       data,
		},
	}
	p.publish(ctx, event)
}

// PublishTrackingFailed publica evento de tracking fallido
func (p *SSEPublisher) PublishTrackingFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.tracking_failed",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"correlation_id": correlationID,
			"error_message":  errorMsg,
		},
	}
	p.publish(ctx, event)
}

// PublishShipmentCancelled publica evento de envío cancelado
func (p *SSEPublisher) PublishShipmentCancelled(ctx context.Context, businessID uint, shipmentID uint) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.cancelled",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"shipment_id": shipmentID,
		},
	}
	p.publish(ctx, event)
}

// PublishCancelFailed publica evento de cancelación fallida
func (p *SSEPublisher) PublishCancelFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	event := shipmentSSEEvent{
		ID:         generateEventID(),
		EventType:  "shipment.cancel_failed",
		BusinessID: businessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"shipment_id":    shipmentID,
			"correlation_id": correlationID,
			"error_message":  errorMsg,
		},
	}
	p.publish(ctx, event)
}

// publish serializa y publica el evento a Redis de forma no-bloqueante
func (p *SSEPublisher) publish(ctx context.Context, event shipmentSSEEvent) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_type", event.EventType).
			Msg("Error serializing shipment SSE event")
		return
	}

	// Publicar de forma no-bloqueante
	go func() {
		publishCtx := context.Background()
		if pubErr := p.redisClient.Client(publishCtx).Publish(publishCtx, p.channel, eventJSON).Err(); pubErr != nil {
			p.logger.Error(publishCtx).
				Err(pubErr).
				Str("event_type", event.EventType).
				Str("channel", p.channel).
				Msg("Error publishing shipment SSE event to Redis")
			return
		}

		p.logger.Info(publishCtx).
			Str("event_id", event.ID).
			Str("event_type", event.EventType).
			Uint("business_id", event.BusinessID).
			Str("channel", p.channel).
			Msg("Shipment SSE event published to Redis")
	}()
}

// ptrToString convierte un puntero de string a string vacío si es nil
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// generateEventID genera un ID único para el evento
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString genera una cadena aleatoria
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// noopSSEPublisher es una implementación no-op para cuando Redis no está disponible
type noopSSEPublisher struct{}

// NewNoopSSEPublisher crea un publisher que no hace nada (para cuando Redis no está disponible)
func NewNoopSSEPublisher() domain.IShipmentSSEPublisher {
	return &noopSSEPublisher{}
}

func (n *noopSSEPublisher) PublishQuoteReceived(_ context.Context, _ uint, _ string, _ map[string]interface{}) {
}
func (n *noopSSEPublisher) PublishQuoteFailed(_ context.Context, _ uint, _ string, _ string) {}
func (n *noopSSEPublisher) PublishGuideGenerated(_ context.Context, _ uint, _ uint, _ string, _ string) {
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
