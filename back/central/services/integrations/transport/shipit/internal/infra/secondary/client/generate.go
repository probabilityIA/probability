package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
)

func (c *Client) CreateShipment(baseURL string, creds domain.Credentials, req domain.ShipmentRequest) (*domain.ShipmentResponse, error) {
	ctx := context.Background()
	url := resolveBaseURL(baseURL) + "/v/shipments"

	c.log.Info(ctx).
		Str("reference", req.Shipment.Reference).
		Bool("sandbox", req.Shipment.Sandbox).
		Msg("Creating Shipit shipment")

	var apiResp domain.ShipmentResponse
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeaders(authHeaders(creds)).
		SetBody(req).
		SetResult(&apiResp).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("error de red al conectar con Shipit: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Shipit shipment creation failed")
		return nil, fmt.Errorf("%s", parseShipitError(resp.StatusCode(), resp.Body()))
	}

	c.log.Info(ctx).
		Int64("shipment_id", apiResp.ID).
		Str("status", apiResp.Status).
		Msg("Shipit shipment created")

	return &apiResp, nil
}
