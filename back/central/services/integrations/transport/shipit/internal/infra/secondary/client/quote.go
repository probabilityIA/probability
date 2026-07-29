package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
)

func (c *Client) Rates(baseURL string, creds domain.Credentials, req domain.RatesRequest) (*domain.RatesResponse, error) {
	ctx := context.Background()
	url := resolveBaseURL(baseURL) + "/v/rates"

	c.log.Info(ctx).
		Uint("origin_id", req.Parcel.OriginID).
		Uint("destiny_id", req.Parcel.DestinyID).
		Msg("Requesting Shipit rates")

	var apiResp domain.RatesResponse
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
			Msg("Shipit rates request failed")
		return nil, fmt.Errorf("%s", parseShipitError(resp.StatusCode(), resp.Body()))
	}

	c.log.Info(ctx).Int("prices", len(apiResp.Prices)).Msg("Shipit rates received")
	return &apiResp, nil
}
