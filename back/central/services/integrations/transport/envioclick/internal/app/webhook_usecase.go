package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IWebhookUseCase interface {
	Process(ctx context.Context, req WebhookRequest) (*WebhookResult, error)
}

type WebhookRequest struct {
	URL         string
	Method      string
	RemoteIP    string
	Headers     map[string][]string
	RawBody     []byte
	Payload     domain.WebhookPayload
}

type WebhookResult struct {
	LogID             uuid.UUID
	CorrelationID     string
	TrackingNumber    string
	ProbabilityStatus domain.ProbabilityShipmentStatus
	RawStatusStep     string
	IsUnknownStatus   bool
	IsIgnored         bool
	IgnoredReason     string
}

type webhookUseCase struct {
	repo      domain.IWebhookLogRepository
	publisher domain.IWebhookResponsePublisher
	log       log.ILogger
}

func NewWebhookUseCase(
	repo domain.IWebhookLogRepository,
	publisher domain.IWebhookResponsePublisher,
	logger log.ILogger,
) IWebhookUseCase {
	return &webhookUseCase{
		repo:      repo,
		publisher: publisher,
		log:       logger.WithModule("transport.envioclick.webhook_usecase"),
	}
}

func (uc *webhookUseCase) Process(ctx context.Context, req WebhookRequest) (*WebhookResult, error) {
	correlationID := "wh-" + uuid.New().String()
	headersBytes, _ := json.Marshal(req.Headers)

	logEntry := &domain.WebhookLog{
		Source:        domain.WebhookSourceEnvioclick,
		EventType:     domain.WebhookEventTrackingUpdate,
		URL:           req.URL,
		Method:        req.Method,
		Headers:       headersBytes,
		RequestBody:   req.RawBody,
		RemoteIP:      req.RemoteIP,
		Status:        domain.WebhookLogStatusReceived,
		ResponseCode:  200,
		CorrelationID: &correlationID,
	}

	if req.Payload.TrackingCode != "" {
		tracking := req.Payload.TrackingCode
		logEntry.TrackingNumber = &tracking
	}

	normalized := req.Payload.ToNormalizedUpdate()
	if normalized == nil {
		logEntry.Status = domain.WebhookLogStatusIgnored
		reason := "Sin eventos en el payload"
		logEntry.ErrorMessage = &reason
		if err := uc.repo.Save(ctx, logEntry); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Failed to save ignored webhook log")
		}
		go uc.trimAsync(domain.WebhookSourceEnvioclick)
		return &WebhookResult{
			LogID:         logEntry.ID,
			CorrelationID: correlationID,
			IsIgnored:     true,
			IgnoredReason: reason,
		}, nil
	}

	mappedStatus := normalized.ProbabilityStatus.String()
	logEntry.MappedStatus = &mappedStatus
	logEntry.RawStatus = &normalized.RawStatusStep

	if err := uc.repo.Save(ctx, logEntry); err != nil {
		uc.log.Error(ctx).Err(err).Str("correlation_id", correlationID).Msg("Failed to save webhook log")
		return nil, fmt.Errorf("failed to save webhook log: %w", err)
	}
	go uc.trimAsync(domain.WebhookSourceEnvioclick)

	if normalized.IsUnknownStatus {
		uc.log.Warn(ctx).
			Str("raw_status_step", normalized.RawStatusStep).
			Str("tracking_number", normalized.TrackingNumber).
			Str("correlation_id", correlationID).
			Msg("Unknown Envioclick statusStep — falling back to in_transit")
	}

	publishMsg := &domain.WebhookUpdateMessage{
		BusinessID:      0,
		CorrelationID:   correlationID,
		TrackingNumber:  normalized.TrackingNumber,
		Carrier:         normalized.Carrier,
		Status:          normalized.ProbabilityStatus,
		RawStatus:       normalized.RawStatusStep,
		RawStatusDetail: normalized.RawStatusDetail,
		HasIncidence:    normalized.HasIncidence,
		IsUnknown:       normalized.IsUnknownStatus,
		Description:     normalized.EventDescription,
		EventTimestamp:  normalized.EventTimestamp,
		ShippedAt:       normalized.ShippedAt,
		DeliveredAt:     normalized.DeliveredAt,
	}

	if err := uc.publisher.PublishWebhookUpdate(ctx, publishMsg); err != nil {
		errMsg := "Failed to publish webhook update: " + err.Error()
		if markErr := uc.repo.MarkProcessed(ctx, logEntry.ID, &errMsg); markErr != nil {
			uc.log.Error(ctx).Err(markErr).Msg("Failed to mark webhook log as failed")
		}
		return nil, fmt.Errorf("failed to publish webhook update: %w", err)
	}

	if markErr := uc.repo.MarkProcessed(ctx, logEntry.ID, nil); markErr != nil {
		uc.log.Warn(ctx).Err(markErr).Msg("Failed to mark webhook log as processed")
	}

	return &WebhookResult{
		LogID:             logEntry.ID,
		CorrelationID:     correlationID,
		TrackingNumber:    normalized.TrackingNumber,
		ProbabilityStatus: normalized.ProbabilityStatus,
		RawStatusStep:     normalized.RawStatusStep,
		IsUnknownStatus:   normalized.IsUnknownStatus,
	}, nil
}

func (uc *webhookUseCase) trimAsync(source string) {
	ctx := context.Background()
	if err := uc.repo.TrimOldBySource(ctx, source, domain.WebhookLogMaxPerSource); err != nil {
		uc.log.Warn(ctx).Err(err).Str("source", source).Msg("Failed to trim old webhook logs")
	}
}
