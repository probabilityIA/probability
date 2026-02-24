package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

// Generate crea un envío (genera guía) en EnvioClick
// Endpoint: POST /shipment
func (c *Client) Generate(baseURL, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	ctx := context.Background()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Str("reference", req.MyShipmentReference).
		Int64("rate_id", req.IDRate).
		Msg("🚀 Generating EnvioClick shipment")

	var apiResp domain.GenerateResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(req).
		SetResult(&apiResp).
		Post(strings.TrimRight(baseURL, "/") + "/shipment")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ EnvioClick generate request failed - network error")
		return nil, fmt.Errorf("error de red al conectar con EnvioClick: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ EnvioClick generate failed")
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	c.log.Info(ctx).
		Str("tracking_number", apiResp.Data.TrackingNumber).
		Str("label_url", apiResp.Data.LabelURL).
		Msg("✅ EnvioClick shipment generated")

	return &apiResp, nil
}
