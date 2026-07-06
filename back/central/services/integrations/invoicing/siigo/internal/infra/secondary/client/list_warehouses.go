package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

type listWarehousesResponse struct {
	Results []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"results"`
}

func (c *Client) ListWarehouses(ctx context.Context, credentials dtos.Credentials) ([]dtos.WarehouseItem, error) {
	c.log.Info(ctx).Msg("Listing Siigo warehouses")

	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	var listResp listWarehousesResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetResult(&listResp).
		Get(c.endpointURL(credentials.BaseURL, "/v1/warehouses"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Siigo list warehouses request failed - network error")
		return nil, fmt.Errorf("error de red al listar bodegas en Siigo: %w", err)
	}

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("Siigo list warehouses failed")
		return nil, fmt.Errorf("error al listar bodegas en Siigo (codigo %d)", resp.StatusCode())
	}

	items := make([]dtos.WarehouseItem, 0, len(listResp.Results))
	for _, r := range listResp.Results {
		items = append(items, dtos.WarehouseItem{ID: r.ID, Name: r.Name})
	}

	c.log.Info(ctx).Int("count", len(items)).Msg("Siigo warehouses listed")

	return items, nil
}
