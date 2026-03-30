package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

// CancelBatch cancela envíos en lote en EnvioClick
// Endpoint: POST /v2cancellation/batch/order
func (c *Client) CancelBatch(baseURL, apiKey string, req domain.CancelBatchRequest) (*domain.CancelBatchResponse, error) {
	ctx := context.Background()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Int("order_count", len(req.Orders)).
		Msg("🗑️ Cancelling EnvioClick shipments in batch")

	var apiResp domain.CancelBatchResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(req).
		SetResult(&apiResp).
		Post(strings.TrimRight(baseURL, "/") + "/v2cancellation/batch/order")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ EnvioClick cancel batch request failed - network error")
		return nil, fmt.Errorf("error de red al conectar con el servicio de transporte: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ EnvioClick cancel batch failed")
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	if apiResp.Status == "" && strings.Contains(strings.ToLower(string(resp.Body())), "success") {
		apiResp.Status = "success"
	}

	c.log.Info(ctx).
		Str("status", apiResp.Status).
		Msg("✅ EnvioClick shipments batch cancellation processed")

	return &apiResp, nil
}
