package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (c *shopifyClient) GetLocations(ctx context.Context, storeName, accessToken string) ([]domain.ShopifyLocation, error) {
	url := buildURL(storeName, "/admin/api/2024-10/locations.json")

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error al obtener locations de Shopify (codigo %d)", resp.StatusCode())
	}

	var parsed struct {
		Locations []struct {
			ID     int64  `json:"id"`
			Name   string `json:"name"`
			Active bool   `json:"active"`
		} `json:"locations"`
	}
	if err := json.Unmarshal(resp.Body(), &parsed); err != nil {
		return nil, fmt.Errorf("error unmarshalling locations: %w", err)
	}

	locations := make([]domain.ShopifyLocation, 0, len(parsed.Locations))
	for _, l := range parsed.Locations {
		locations = append(locations, domain.ShopifyLocation{ID: l.ID, Name: l.Name})
	}
	return locations, nil
}
