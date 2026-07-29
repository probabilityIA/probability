package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
)

func (c *Client) GetCommunes(baseURL string, creds domain.Credentials) ([]domain.Commune, error) {
	ctx := context.Background()
	url := resolveBaseURL(baseURL) + "/v/communes"

	var apiResp []domain.Commune
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeaders(authHeaders(creds)).
		SetResult(&apiResp).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error de red al conectar con Shipit: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Shipit communes request failed")
		return nil, fmt.Errorf("%s", parseShipitError(resp.StatusCode(), resp.Body()))
	}

	return apiResp, nil
}
