package client

import (
	"context"
	"encoding/json"
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

	// No usar SetResult para evitar error de unmarshal cuando EnvioClick
	// devuelve el campo "tracker" como número en lugar de string
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(req).
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

	// Parsear manualmente con interface{} en tracker para tolerar string o número
	var raw struct {
		Status string `json:"status"`
		Data   struct {
			Tracker          interface{} `json:"tracker"`
			URL              string      `json:"url"`
			MyGuideReference string      `json:"myGuideReference"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &raw); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(resp.Body())).Msg("❌ Error parsing EnvioClick response")
		return nil, fmt.Errorf("error parseando respuesta de EnvioClick: %w", err)
	}

	// Convertir tracker a string sin importar si viene como string o número
	var trackingNumber string
	switch v := raw.Data.Tracker.(type) {
	case string:
		trackingNumber = v
	case float64:
		trackingNumber = fmt.Sprintf("%.0f", v)
	default:
		if v != nil {
			trackingNumber = fmt.Sprintf("%v", v)
		}
	}

	apiResp := &domain.GenerateResponse{
		Status: raw.Status,
		Data: domain.GenerateData{
			TrackingNumber:   trackingNumber,
			LabelURL:         raw.Data.URL,
			MyGuideReference: raw.Data.MyGuideReference,
		},
	}

	c.log.Info(ctx).
		Str("tracking_number", apiResp.Data.TrackingNumber).
		Str("label_url", apiResp.Data.LabelURL).
		Msg("✅ EnvioClick shipment generated")

	return apiResp, nil
}
