package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *shopifyClient) SetInventoryLevel(ctx context.Context, storeName, accessToken string, locationID, inventoryItemID int64, available int) error {
	url := buildURL(storeName, "/admin/api/2024-10/inventory_levels/set.json")

	body := map[string]interface{}{
		"location_id":       locationID,
		"inventory_item_id": inventoryItemID,
		"available":         available,
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-Shopify-Access-Token", accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("error al fijar inventario en Shopify (codigo %d): %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}
