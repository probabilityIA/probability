package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

// Generate crea un envío (genera guía) en EnvioClick
// Endpoint: POST /shipment
func (c *Client) Generate(baseURL, apiKey string, req domain.QuoteRequest, meta *domain.SyncMeta) (*domain.GenerateResponse, error) {
	ctx := context.Background()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	if req.CODValue == 0 {
		req.CODPaymentMethod = ""
	}

	c.log.Info(ctx).
		Str("reference", req.MyShipmentReference).
		Int64("rate_id", req.IDRate).
		Msg("🚀 Generating EnvioClick shipment")

	url := strings.TrimRight(baseURL, "/") + "/shipment"
	started := time.Now()
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetBody(req).
		SetDebug(true).
		Post(url)
	captureMeta(meta, "POST", url, req, started, resp, err)

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ EnvioClick generate request failed - network error")
		return nil, fmt.Errorf("error de red al conectar con el servicio de transporte: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ EnvioClick generate failed")
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	// Parsear manualmente con interface{} en tracker e idOrder para tolerar variaciones
	var raw struct {
		Status string `json:"status"`
		Data   struct {
			Tracker          interface{} `json:"tracker"`
			IDOrder          interface{} `json:"idOrder"`
			URL              string      `json:"url"`
			Carrier          string      `json:"carrier"`
			MyGuideReference string      `json:"myGuideReference"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &raw); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(resp.Body())).Msg("❌ Error parsing EnvioClick response")
		return nil, fmt.Errorf("error parseando respuesta del servicio de transporte: %w", err)
	}

	// Convertir tracker a string
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

	// Convertir idOrder a int64
	var idOrder int64
	switch v := raw.Data.IDOrder.(type) {
	case float64:
		idOrder = int64(v)
	case int:
		idOrder = int64(v)
	case string:
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			idOrder = n
		}
	}

	apiResp := &domain.GenerateResponse{
		Status: raw.Status,
		Data: domain.GenerateData{
			TrackingNumber:   trackingNumber,
			LabelURL:         raw.Data.URL,
			MyGuideReference: raw.Data.MyGuideReference,
			IDOrder:          idOrder,
			Carrier:          raw.Data.Carrier,
		},
	}

	c.log.Info(ctx).
		Str("tracking_number", apiResp.Data.TrackingNumber).
		Str("label_url", apiResp.Data.LabelURL).
		Msg("✅ EnvioClick shipment generated")

	return apiResp, nil
}
