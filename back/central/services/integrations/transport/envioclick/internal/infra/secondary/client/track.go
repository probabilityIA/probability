package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

// Track obtiene el estado de tracking de un envío en EnvioClick usando trackingCode
// Endpoint: POST /track
func (c *Client) Track(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
	ctx := context.Background()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Str("tracking_number", trackingNumber).
		Msg("🔍 Tracking EnvioClick shipment")

	payload := map[string]string{"trackingCode": trackingNumber}

	var apiResp domain.TrackingResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(payload).
		SetResult(&apiResp).
		SetDebug(true).
		Post(strings.TrimRight(baseURL, "/") + "/track")

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

// TrackByOrdersBatch obtiene el estado de múltiples envíos usando sus idOrder
// Endpoint: POST /track-by-orders
func (c *Client) TrackByOrdersBatch(baseURL, apiKey string, orders []int64) (*domain.TrackingResponse, error) {
	ctx := context.Background()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Int("order_count", len(orders)).
		Msg("🔍 Tracking EnvioClick shipments in batch")

	payload := map[string]interface{}{"orders": orders}

	var apiResp domain.TrackingResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(payload).
		SetResult(&apiResp).
		Post(strings.TrimRight(baseURL, "/") + "/track-by-orders")

	if err != nil {
		return nil, fmt.Errorf("error de red: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	return &apiResp, nil
}
