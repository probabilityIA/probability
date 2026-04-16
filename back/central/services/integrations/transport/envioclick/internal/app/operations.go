package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

func (uc *useCase) Quote(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	uc.log.Info(ctx).Msg("Quoting shipment")
	return uc.client.Quote(baseURL, apiKey, req)
}

func (uc *useCase) Generate(ctx context.Context, baseURL, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	uc.log.Info(ctx).Msg("Generating guide")
	return uc.client.Generate(baseURL, apiKey, req)
}

func (uc *useCase) Track(ctx context.Context, baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
	uc.log.Info(ctx).Str("tracking_number", trackingNumber).Msg("Tracking shipment")
	return uc.client.Track(baseURL, apiKey, trackingNumber)
}

func (uc *useCase) Cancel(ctx context.Context, baseURL, apiKey string, trackingNumber string, idOrder int64) (*domain.CancelResponse, error) {
	uc.log.Info(ctx).
		Str("tracking_number", trackingNumber).
		Int64("id_order", idOrder).
		Msg("Verifying and canceling shipment")

	trackResp, err := uc.client.Track(baseURL, apiKey, trackingNumber)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("tracking_number", trackingNumber).Msg("Verification failed: tracking failed")
		return nil, fmt.Errorf("no se pudo verificar el estado del envio: %w", err)
	}

	status := trackResp.Data.Status
	statusLower := strings.ToLower(status)
	if !strings.Contains(statusLower, "pendiente") {
		uc.log.Warn(ctx).
			Str("status", status).
			Str("tracking_number", trackingNumber).
			Msg("Cancellation aborted: invalid status")
		return nil, fmt.Errorf("el envio no puede ser cancelado porque se encuentra en estado: %s", status)
	}

	if idOrder != 0 {
		uc.log.Info(ctx).Int64("id_order", idOrder).Msg("Proceeding to cancel via Batch API v2")
		batchResp, err := uc.client.CancelBatch(baseURL, apiKey, domain.CancelBatchRequest{
			IDOrders: []int64{idOrder},
		})
		if err != nil {
			return nil, err
		}

		msg := "Cancelado exitosamente"
		resStatus := "success"
		if len(batchResp.Data.NotValidOrders) > 0 {
			resStatus = "error"
			msg = "EnvioClick reporto la orden como no valida para cancelacion"
		}

		return &domain.CancelResponse{
			Status:  resStatus,
			Message: msg,
		}, nil
	}

	uc.log.Warn(ctx).Str("tracking_number", trackingNumber).Msg("No idOrder provided, falling back to singular DELETE")
	return uc.client.Cancel(baseURL, apiKey, trackingNumber)
}

func (uc *useCase) CancelBatch(ctx context.Context, baseURL, apiKey string, req domain.CancelBatchRequest) (*domain.CancelBatchResponse, error) {
	uc.log.Info(ctx).Int("order_count", len(req.IDOrders)).Msg("Canceling shipments in batch")
	return uc.client.CancelBatch(baseURL, apiKey, req)
}

func (uc *useCase) TrackByOrdersBatch(ctx context.Context, baseURL, apiKey string, orders []int64) (*domain.TrackingResponse, error) {
	uc.log.Info(ctx).Int("order_count", len(orders)).Msg("Tracking shipments in batch")
	return uc.client.TrackByOrdersBatch(baseURL, apiKey, orders)
}

