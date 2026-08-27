package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

var (
	cancelVerifyAttempts = 3
	cancelVerifyDelay    = 2 * time.Second
)

func appendMeta(metas *[]domain.SyncMeta, m domain.SyncMeta) {
	if metas != nil {
		*metas = append(*metas, m)
	}
}

func (uc *useCase) Quote(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, metas *[]domain.SyncMeta) (*domain.QuoteResponse, error) {
	uc.log.Info(ctx).Msg("Quoting shipment")
	var meta domain.SyncMeta
	resp, err := uc.client.Quote(baseURL, apiKey, req, &meta)
	appendMeta(metas, meta)
	return resp, err
}

func (uc *useCase) Generate(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest, metas *[]domain.SyncMeta) (*domain.GenerateResponse, error) {
	uc.log.Info(ctx).Msg("Generating guide")
	var meta domain.SyncMeta
	resp, err := uc.client.Generate(baseURL, apiKey, req, &meta)
	appendMeta(metas, meta)
	return resp, err
}

func (uc *useCase) Track(ctx context.Context, baseURL, apiKey string, trackingNumber string, metas *[]domain.SyncMeta) (*domain.TrackingResponse, error) {
	uc.log.Info(ctx).Str("tracking_number", trackingNumber).Msg("Tracking shipment")
	var meta domain.SyncMeta
	resp, err := uc.client.Track(baseURL, apiKey, trackingNumber, &meta)
	appendMeta(metas, meta)
	return resp, err
}

func (uc *useCase) Cancel(ctx context.Context, baseURL, apiKey string, trackingNumber string, idOrder int64, metas *[]domain.SyncMeta) (*domain.CancelResponse, error) {
	uc.log.Info(ctx).
		Str("tracking_number", trackingNumber).
		Int64("id_order", idOrder).
		Msg("Verifying and canceling shipment")

	var trackMeta domain.SyncMeta
	trackResp, err := uc.client.Track(baseURL, apiKey, trackingNumber, &trackMeta)
	appendMeta(metas, trackMeta)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("tracking_number", trackingNumber).Msg("Verification failed: tracking failed")
		return nil, fmt.Errorf("no se pudo verificar el estado del envio: %w", err)
	}

	status := trackResp.Data.Status
	statusLower := strings.ToLower(status)
	if strings.Contains(statusLower, "cancelad") {
		uc.log.Info(ctx).
			Str("tracking_number", trackingNumber).
			Msg("Shipment already cancelled in carrier — returning success to sync DB")
		return &domain.CancelResponse{
			Status:  "success",
			Message: "El envio ya estaba cancelado en el carrier",
		}, nil
	}
	if !strings.Contains(statusLower, "pendiente") {
		uc.log.Warn(ctx).
			Str("status", status).
			Str("tracking_number", trackingNumber).
			Msg("Cancellation aborted: invalid status")
		return nil, fmt.Errorf("el envio no puede ser cancelado porque se encuentra en estado: %s", status)
	}

	if idOrder != 0 {
		uc.log.Info(ctx).Int64("id_order", idOrder).Msg("Proceeding to cancel via Batch API v2")
		var batchMeta domain.SyncMeta
		batchResp, err := uc.client.CancelBatch(baseURL, apiKey, domain.CancelBatchRequest{
			IDOrders: []int64{idOrder},
		}, &batchMeta)
		appendMeta(metas, batchMeta)
		if err != nil {
			return nil, err
		}

		if len(batchResp.Data.NotValidOrders) > 0 {
			uc.log.Warn(ctx).
				Int64("id_order", idOrder).
				Msg("Cancellation rejected: order reported as not valid")
			return nil, fmt.Errorf("el envio no puede ser cancelado: la transportadora lo reporto como no valido")
		}

		return uc.verifyCancelled(ctx, baseURL, apiKey, trackingNumber, idOrder, metas)
	}

	uc.log.Warn(ctx).Str("tracking_number", trackingNumber).Msg("No idOrder provided, falling back to singular DELETE")
	var cancelMeta domain.SyncMeta
	resp, err := uc.client.Cancel(baseURL, apiKey, trackingNumber, &cancelMeta)
	appendMeta(metas, cancelMeta)
	return resp, err
}

func (uc *useCase) verifyCancelled(ctx context.Context, baseURL, apiKey, trackingNumber string, idOrder int64, metas *[]domain.SyncMeta) (*domain.CancelResponse, error) {
	var lastStatus string

	for attempt := 1; attempt <= cancelVerifyAttempts; attempt++ {
		if attempt > 1 {
			time.Sleep(cancelVerifyDelay)
		}

		var verifyMeta domain.SyncMeta
		verifyResp, err := uc.client.Track(baseURL, apiKey, trackingNumber, &verifyMeta)
		appendMeta(metas, verifyMeta)
		if err != nil {
			uc.log.Warn(ctx).Err(err).
				Str("tracking_number", trackingNumber).
				Int("attempt", attempt).
				Msg("Cancellation verification failed: could not track shipment")
			continue
		}

		lastStatus = verifyResp.Data.Status
		if strings.Contains(strings.ToLower(lastStatus), "cancelad") {
			uc.log.Info(ctx).
				Str("tracking_number", trackingNumber).
				Int("attempt", attempt).
				Msg("Cancellation confirmed by carrier")
			return &domain.CancelResponse{
				Status:  "success",
				Message: "Cancelado exitosamente",
			}, nil
		}

		uc.log.Warn(ctx).
			Str("tracking_number", trackingNumber).
			Str("carrier_status", lastStatus).
			Int("attempt", attempt).
			Msg("Cancellation not reflected by carrier yet")
	}

	uc.log.Error(ctx).
		Str("tracking_number", trackingNumber).
		Int64("id_order", idOrder).
		Str("carrier_status", lastStatus).
		Msg("Cancellation accepted by EnvioClick but shipment is still active - NOT marking as cancelled")

	if lastStatus == "" {
		return nil, fmt.Errorf("no se pudo confirmar la cancelacion con la transportadora: comunicate con soporte para cancelar la guia")
	}
	return nil, fmt.Errorf("la transportadora acepto la solicitud pero el envio sigue activo (estado: %s): comunicate con soporte para cancelar la guia", lastStatus)
}

func (uc *useCase) CancelBatch(ctx context.Context, baseURL, apiKey string, req domain.CancelBatchRequest, metas *[]domain.SyncMeta) (*domain.CancelBatchResponse, error) {
	uc.log.Info(ctx).Int("order_count", len(req.IDOrders)).Msg("Canceling shipments in batch")
	var meta domain.SyncMeta
	resp, err := uc.client.CancelBatch(baseURL, apiKey, req, &meta)
	appendMeta(metas, meta)
	return resp, err
}

func (uc *useCase) TrackByOrdersBatch(ctx context.Context, baseURL, apiKey string, orders []int64, metas *[]domain.SyncMeta) (*domain.TrackingResponse, error) {
	uc.log.Info(ctx).Int("order_count", len(orders)).Msg("Tracking shipments in batch")
	var meta domain.SyncMeta
	resp, err := uc.client.TrackByOrdersBatch(baseURL, apiKey, orders, &meta)
	appendMeta(metas, meta)
	return resp, err
}
