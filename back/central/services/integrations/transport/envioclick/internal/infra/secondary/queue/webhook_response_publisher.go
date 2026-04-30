package queue

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

const (
	OperationWebhookUpdate = "webhook_update"
	ProviderEnvioclick     = "envioclick"
)

type WebhookResponsePublisher struct {
	inner *ResponsePublisher
}

func NewWebhookResponsePublisher(inner *ResponsePublisher) domain.IWebhookResponsePublisher {
	return &WebhookResponsePublisher{inner: inner}
}

func (p *WebhookResponsePublisher) PublishWebhookUpdate(ctx context.Context, msg *domain.WebhookUpdateMessage) error {
	data := map[string]any{
		"tracking_number":    msg.TrackingNumber,
		"carrier":            msg.Carrier,
		"probability_status": msg.Status.String(),
		"raw_status":         msg.RawStatus,
		"raw_status_detail":  msg.RawStatusDetail,
		"has_incidence":      msg.HasIncidence,
		"is_unknown_status":  msg.IsUnknown,
		"event_description":  msg.Description,
		"event_timestamp":    msg.EventTimestamp,
	}
	if msg.ShippedAt != nil {
		data["shipped_at"] = *msg.ShippedAt
	}
	if msg.DeliveredAt != nil {
		data["delivered_at"] = *msg.DeliveredAt
	}

	response := &TransportResponseMessage{
		ShipmentID:    msg.ShipmentID,
		BusinessID:    msg.BusinessID,
		Provider:      ProviderEnvioclick,
		Operation:     OperationWebhookUpdate,
		Status:        "success",
		CorrelationID: msg.CorrelationID,
		Timestamp:     time.Now(),
		Data:          data,
	}

	return p.inner.PublishResponse(ctx, response)
}
