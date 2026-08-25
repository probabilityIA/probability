package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

func (c *Client) Track(baseURL, apiKey string, trackingNumber string, meta *domain.SyncMeta) (*domain.TrackingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Str("tracking_number", trackingNumber).
		Msg("🔍 Tracking EnvioClick shipment")

	payload := map[string]string{"trackingCode": trackingNumber}

	var apiResp domain.TrackingResponse
	url := strings.TrimRight(baseURL, "/") + "/track"
	started := time.Now()
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(payload).
		SetResult(&apiResp).
		SetDebug(true).
		Post(url)
	captureMeta(meta, "POST", url, payload, started, resp, err)

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ EnvioClick track request failed")
		return nil, fmt.Errorf("error de red al conectar con el servicio de transporte: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ EnvioClick tracking failed")
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	c.log.Info(ctx).
		Str("status", apiResp.Data.Status).
		Str("detail", apiResp.Data.StatusDetail).
		Msg("✅ EnvioClick tracking received")

	return &apiResp, nil
}

func (c *Client) TrackByOrdersBatch(baseURL, apiKey string, orders []int64, meta *domain.SyncMeta) (*domain.TrackingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Int("order_count", len(orders)).
		Msg("🔍 Tracking EnvioClick shipments in batch")

	payload := map[string]interface{}{"orders": orders}

	var apiResp domain.TrackingResponse
	url := strings.TrimRight(baseURL, "/") + "/track-by-orders"
	started := time.Now()
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(payload).
		SetResult(&apiResp).
		Post(url)
	captureMeta(meta, "POST", url, payload, started, resp, err)

	if err != nil {
		return nil, fmt.Errorf("error de red: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	return &apiResp, nil
}
