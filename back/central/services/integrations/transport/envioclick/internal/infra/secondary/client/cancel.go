package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

// Cancel cancela un envío en EnvioClick
// Endpoint: DELETE /shipment/{idShipment}
func (c *Client) Cancel(baseURL, apiKey string, idShipment string) (*domain.CancelResponse, error) {
	ctx := context.Background()

	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c.log.Info(ctx).
		Str("shipment_id", idShipment).
		Msg("🗑️ Cancelling EnvioClick shipment")

	var apiResp domain.CancelResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", apiKey).
		SetResult(&apiResp).
		Delete(strings.TrimRight(baseURL, "/") + fmt.Sprintf("/shipment/%s", idShipment))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ EnvioClick cancel request failed - network error")
		return nil, fmt.Errorf("error de red al conectar con EnvioClick: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ EnvioClick cancel failed")
		return nil, fmt.Errorf("%s", parseEnvioClickError(resp.Body()))
	}

	// EnvioClick a veces retorna texto plano con "success" en vez de JSON estructurado
	if apiResp.Status == "" && strings.Contains(strings.ToLower(string(resp.Body())), "success") {
		apiResp = domain.CancelResponse{Status: "success", Message: "Cancelación exitosa"}
	}

	c.log.Info(ctx).
		Str("status", apiResp.Status).
		Msg("✅ EnvioClick shipment cancelled")

	return &apiResp, nil
}
