package app

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	maxConcurrentTracks = 10
)

type ISyncUseCase interface {
	SyncBatch(ctx context.Context, req domain.SyncBatchRequest) (*domain.SyncBatchResult, error)
}

type syncUseCase struct {
	client    domain.IEnvioClickClient
	repo      domain.IWebhookLogRepository
	publisher domain.IWebhookResponsePublisher
	log       log.ILogger
}

func NewSyncUseCase(
	client domain.IEnvioClickClient,
	repo domain.IWebhookLogRepository,
	publisher domain.IWebhookResponsePublisher,
	logger log.ILogger,
) ISyncUseCase {
	return &syncUseCase{
		client:    client,
		repo:      repo,
		publisher: publisher,
		log:       logger.WithModule("transport.envioclick.sync_usecase"),
	}
}

func (uc *syncUseCase) SyncBatch(ctx context.Context, req domain.SyncBatchRequest) (*domain.SyncBatchResult, error) {
	result := &domain.SyncBatchResult{Total: len(req.Items)}
	if len(req.Items) == 0 {
		return result, nil
	}

	uc.processIndividualTracks(ctx, req, req.Items, result)

	uc.log.Info(ctx).
		Str("correlation_id", req.CorrelationID).
		Uint("business_id", req.BusinessID).
		Int("total", result.Total).
		Int("processed", result.Processed).
		Int("failed", result.Failed).
		Int("unknown", result.Unknown).
		Int("not_found", result.NotFound).
		Msg("Envioclick sync batch completed")

	return result, nil
}

func (uc *syncUseCase) processIndividualTracks(ctx context.Context, req domain.SyncBatchRequest, items []domain.SyncBatchItem, result *domain.SyncBatchResult) {
	sem := make(chan struct{}, maxConcurrentTracks)
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		sem <- struct{}{}

		go func(it domain.SyncBatchItem) {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := uc.client.Track(req.BaseURL, req.APIKey, it.TrackingNumber)
			if err != nil {
				uc.recordFailure(ctx, req, it, err.Error(), result)
				return
			}

			rawBody, _ := json.Marshal(resp.Data)
			uc.processSingleResult(ctx, req, it, resp.Data.Status, resp.Data.StatusDetail, resp.Data.RealPickupDate, resp.Data.RealDeliveryDate, rawBody, result)
		}(item)
	}

	wg.Wait()
}

func (uc *syncUseCase) processSingleResult(
	ctx context.Context,
	req domain.SyncBatchRequest,
	item domain.SyncBatchItem,
	apiStatus, apiStatusDetail string,
	realPickupDate, realDeliveryDate *string,
	rawBody []byte,
	result *domain.SyncBatchResult,
) {
	step := domain.ApiStatusToStep(apiStatus, apiStatusDetail)
	probStatus, unknown := domain.MapStatusStepToProbability(step, false)

	if unknown {
		uc.log.Warn(ctx).
			Str("tracking_number", item.TrackingNumber).
			Str("api_status", apiStatus).
			Str("api_status_detail", apiStatusDetail).
			Msg("Unknown API status during sync — falling back to in_transit")
	}

	itemCorrelationID := req.CorrelationID + "-" + uuid.New().String()[:8]
	mappedStatus := probStatus.String()
	trackingNumber := item.TrackingNumber

	logEntry := &domain.WebhookLog{
		Source:        domain.WebhookSourceEnvioclick,
		EventType:     domain.WebhookEventSync,
		URL:           req.URL,
		Method:        "INTERNAL",
		RequestBody:   rawBody,
		RemoteIP:      req.RemoteIP,
		Status:        domain.WebhookLogStatusReceived,
		ResponseCode:  200,
		ShipmentID:    &item.ShipmentID,
		BusinessID:    &req.BusinessID,
		CorrelationID: &itemCorrelationID,
		TrackingNumber: &trackingNumber,
		MappedStatus:  &mappedStatus,
		RawStatus:     &step,
	}
	if err := uc.repo.Save(ctx, logEntry); err != nil {
		uc.log.Warn(ctx).Err(err).Str("tracking_number", item.TrackingNumber).Msg("Failed to save sync log")
	}

	shipmentIDCopy := item.ShipmentID
	msg := &domain.WebhookUpdateMessage{
		ShipmentID:     &shipmentIDCopy,
		BusinessID:     req.BusinessID,
		CorrelationID:  itemCorrelationID,
		TrackingNumber: item.TrackingNumber,
		Status:         probStatus,
		RawStatus:      step,
		HasIncidence:   false,
		IsUnknown:      unknown,
		Description:    apiStatusDetail,
		ShippedAt:      realPickupDate,
		DeliveredAt:    realDeliveryDate,
	}
	if err := uc.publisher.PublishWebhookUpdate(ctx, msg); err != nil {
		errMsg := "publish failed: " + err.Error()
		_ = uc.repo.MarkProcessed(ctx, logEntry.ID, &errMsg)
		result.Failed++
		return
	}

	if err := uc.repo.MarkProcessed(ctx, logEntry.ID, nil); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to mark sync log as processed")
	}
	result.Processed++
	if unknown {
		result.Unknown++
	}
}

func (uc *syncUseCase) recordFailure(ctx context.Context, req domain.SyncBatchRequest, item domain.SyncBatchItem, errMsg string, result *domain.SyncBatchResult) {
	uc.log.Warn(ctx).
		Str("tracking_number", item.TrackingNumber).
		Str("error", errMsg).
		Msg("Sync item failed")

	trackingCopy := item.TrackingNumber
	shipmentCopy := item.ShipmentID
	itemCorrID := req.CorrelationID + "-" + uuid.New().String()[:8]
	errCopy := errMsg

	logEntry := &domain.WebhookLog{
		Source:         domain.WebhookSourceEnvioclick,
		EventType:      domain.WebhookEventSync,
		URL:            req.URL,
		Method:         "INTERNAL",
		RequestBody:    []byte("{}"),
		RemoteIP:       req.RemoteIP,
		Status:         domain.WebhookLogStatusFailed,
		ResponseCode:   500,
		ShipmentID:     &shipmentCopy,
		BusinessID:     &req.BusinessID,
		CorrelationID:  &itemCorrID,
		TrackingNumber: &trackingCopy,
		ErrorMessage:   &errCopy,
	}
	if err := uc.repo.Save(ctx, logEntry); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to save failed sync log")
	}
	result.Failed++
}

