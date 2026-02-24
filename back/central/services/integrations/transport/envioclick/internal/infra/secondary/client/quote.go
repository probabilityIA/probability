package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

// Quote obtiene cotizaciones de envío desde EnvioClick
// Endpoint: POST /quotation
func (c *Client) Quote(baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	ctx := context.Background()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Str("origin_dane", req.Origin.DaneCode).
		Str("dest_dane", req.Destination.DaneCode).
		Int("packages", len(req.Packages)).
		Msg("📦 Requesting EnvioClick quote")

	var apiResp domain.QuoteResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(req).
		SetResult(&apiResp).
		SetDebug(true).
		Post(strings.TrimRight(baseURL, "/") + "/quotation")

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ EnvioClick quote request failed - network error")
		return nil, fmt.Errorf("error de red al conectar con EnvioClick: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ EnvioClick quote failed")
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	c.log.Info(ctx).
		Int("rates_count", len(apiResp.Data.Rates)).
		Msg("✅ EnvioClick quote received")

	return &apiResp, nil
}
