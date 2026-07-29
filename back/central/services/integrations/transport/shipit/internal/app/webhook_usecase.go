package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IWebhookUseCase interface {
	Process(ctx context.Context, req WebhookRequest) (*WebhookResult, error)
}

type WebhookRequest struct {
	URL      string
	Method   string
	RemoteIP string
	Headers  map[string][]string
	RawBody  []byte
	Payload  domain.WebhookPayload
}

type WebhookResult struct {
	CorrelationID     string
	TrackingNumber    string
	ProbabilityStatus domain.ProbabilityShipmentStatus
	RawStatus         string
	IsUnknownStatus   bool
	IsIgnored         bool
	IgnoredReason     string
}

type webhookUseCase struct {
	publisher domain.IWebhookResponsePublisher
	log       log.ILogger
}

func NewWebhookUseCase(publisher domain.IWebhookResponsePublisher, logger log.ILogger) IWebhookUseCase {
	return &webhookUseCase{
		publisher: publisher,
		log:       logger.WithModule("transport.shipit.webhook_usecase"),
	}
}

func (uc *webhookUseCase) Process(ctx context.Context, req WebhookRequest) (*WebhookResult, error) {
	correlationID := "wh-" + uuid.New().String()

	normalized := req.Payload.ToNormalizedUpdate()
	if normalized == nil {
		return &WebhookResult{
			CorrelationID: correlationID,
			IsIgnored:     true,
			IgnoredReason: "Payload sin tracking_number o sin estado",
		}, nil
	}

	if normalized.IsUnknownStatus {
		uc.log.Warn(ctx).
			Str("raw_status", normalized.RawStatus).
			Str("tracking_number", normalized.TrackingNumber).
			Str("correlation_id", correlationID).
			Msg("Estado Shipit desconocido, se usa in_transit como fallback")
	}

	carrier := normalized.Carrier
	if carrier == "" {
		carrier = "shipit"
	}

	msg := &domain.WebhookUpdateMessage{
		BusinessID:      0,
		CorrelationID:   correlationID,
		TrackingNumber:  normalized.TrackingNumber,
		Carrier:         carrier,
		Status:          normalized.ProbabilityStatus,
		RawStatus:       normalized.RawStatus,
		RawStatusDetail: normalized.RawStatusDetail,
		IsUnknown:       normalized.IsUnknownStatus,
		Description:     normalized.RawStatusDetail,
		EventTimestamp:  normalized.EventTimestamp,
		ShippedAt:       normalized.ShippedAt,
		DeliveredAt:     normalized.DeliveredAt,
	}

	if err := uc.publisher.PublishWebhookUpdate(ctx, msg); err != nil {
		return nil, err
	}

	return &WebhookResult{
		CorrelationID:     correlationID,
		TrackingNumber:    normalized.TrackingNumber,
		ProbabilityStatus: normalized.ProbabilityStatus,
		RawStatus:         normalized.RawStatus,
		IsUnknownStatus:   normalized.IsUnknownStatus,
	}, nil
}
