package queue

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type SSEPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

func NewSSEPublisher(queue rabbitmq.IQueue, logger log.ILogger) domain.IShipmentSSEPublisher {
	return &SSEPublisher{
		queue:  queue,
		logger: logger.WithModule("shipments.sse_publisher"),
	}
}

func (p *SSEPublisher) PublishQuoteReceived(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	p.publish(ctx, "shipment.quote_received", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"quotes":         data,
	})
}

func (p *SSEPublisher) PublishQuoteFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.quote_failed", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

func (p *SSEPublisher) PublishGuideGenerated(ctx context.Context, businessID uint, shipmentID uint, correlationID string, trackingNumber string, labelURL string, carrier string, notification *domain.GuideNotificationData) {
	data := map[string]interface{}{
		"shipment_id":     shipmentID,
		"correlation_id":  correlationID,
		"tracking_number": trackingNumber,
		"label_url":       labelURL,
		"carrier":         carrier,
	}

	if notification != nil {
		data["customer_name"] = notification.CustomerName
		data["customer_phone"] = notification.CustomerPhone
		data["order_number"] = notification.OrderNumber
		data["business_name"] = notification.BusinessName
		data["integration_id"] = notification.IntegrationID
		if notification.CodTotal != nil {
			data["cod_total"] = *notification.CodTotal
		}
		if notification.CodCarrierFee != nil {
			data["cod_carrier_fee"] = *notification.CodCarrierFee
		}
		if notification.TrackingURL != "" {
			data["tracking_url"] = notification.TrackingURL
		}
	}

	p.publish(ctx, "shipment.guide_generated", businessID, data)
}

func (p *SSEPublisher) PublishGuideFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.guide_failed", businessID, map[string]interface{}{
		"shipment_id":    shipmentID,
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

func (p *SSEPublisher) PublishTrackingUpdated(ctx context.Context, businessID uint, correlationID string, data map[string]interface{}) {
	p.publish(ctx, "shipment.tracking_updated", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"tracking":       data,
	})
}

func (p *SSEPublisher) PublishTrackingFailed(ctx context.Context, businessID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.tracking_failed", businessID, map[string]interface{}{
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

func (p *SSEPublisher) PublishShipmentCancelled(ctx context.Context, businessID uint, shipmentID uint) {
	p.publish(ctx, "shipment.cancelled", businessID, map[string]interface{}{
		"shipment_id": shipmentID,
	})
}

func (p *SSEPublisher) PublishCancelFailed(ctx context.Context, businessID uint, shipmentID uint, correlationID string, errorMsg string) {
	p.publish(ctx, "shipment.cancel_failed", businessID, map[string]interface{}{
		"shipment_id":    shipmentID,
		"correlation_id": correlationID,
		"error_message":  errorMsg,
	})
}

func (p *SSEPublisher) PublishSyncBatchCompleted(ctx context.Context, businessID uint, correlationID string, summary map[string]interface{}) {
	data := map[string]interface{}{"correlation_id": correlationID}
	for k, v := range summary {
		data[k] = v
	}
	p.publish(ctx, "shipment.sync_batch_completed", businessID, data)
}

func (p *SSEPublisher) publish(ctx context.Context, eventType string, businessID uint, data map[string]interface{}) {

	var integrationID uint
	if rawID, ok := data["integration_id"]; ok {
		switch v := rawID.(type) {
		case uint:
			integrationID = v
		case float64:
			integrationID = uint(v)
		case int:
			integrationID = uint(v)
		}
	}

	go func() {
		if err := rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
			Type:          eventType,
			Category:      "shipment",
			BusinessID:    businessID,
			IntegrationID: integrationID,
			Data:          data,
		}); err != nil {
			p.logger.Error(ctx).
				Err(err).
				Str("event_type", eventType).
				Uint("business_id", businessID).
				Msg("Error publishing shipment SSE event to RabbitMQ")
		}
	}()
}

type noopSSEPublisher struct{}

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
func (n *noopSSEPublisher) PublishSyncBatchCompleted(_ context.Context, _ uint, _ string, _ map[string]interface{}) {
}

var _ domain.IShipmentSSEPublisher = (*SSEPublisher)(nil)
var _ domain.IShipmentSSEPublisher = (*noopSSEPublisher)(nil)
