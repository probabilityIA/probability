package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
)

func (c *Client) Track(baseURL string, creds domain.Credentials, trackingNumber string) (*domain.TrackingResponse, error) {
	ctx := context.Background()
	endpoint := resolveTrackingBaseURL(baseURL) + "/v/trackings/number/" + url.PathEscape(trackingNumber)

	c.log.Info(ctx).Str("tracking_number", trackingNumber).Msg("Tracking Shipit shipment")

	var apiResp domain.TrackingResponse
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeaders(authHeaders(creds)).
		SetResult(&apiResp).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("error de red al conectar con Shipit: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Shipit tracking failed")
		return nil, fmt.Errorf("%s", parseShipitError(resp.StatusCode(), resp.Body()))
	}

	c.log.Info(ctx).
		Str("courier", apiResp.Data.Courier).
		Int("statuses", len(apiResp.Data.Statuses)).
		Msg("Shipit tracking received")

	return &apiResp, nil
}
