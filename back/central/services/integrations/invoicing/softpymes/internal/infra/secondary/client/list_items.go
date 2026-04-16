package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
)

// ItemResponse representa un ítem del catálogo de Softpymes.
// Softpymes retorna itemPrice y unitCost como strings ("18000.00").
type ItemResponse struct {
	ItemCode      string `json:"itemCode"`
	ItemName      string `json:"itemName"`
	ItemPrice     string `json:"itemPrice"`
	UnitCost      string `json:"unitCost"`
	Description   string `json:"description"`
	MinimumStock  string `json:"minimumStock"`
	OrderQuantity string `json:"orderQuantity"`
}

// ListItems lista ítems del catálogo de Softpymes.
// Endpoint: GET /app/integration/items?page=X&pageSize=Y
func (c *Client) ListItems(ctx context.Context, apiKey, apiSecret, referer, baseURL string, page, pageSize int) ([]ports.ListedItem, error) {
	c.log.Info(ctx).
		Int("page", page).
		Int("page_size", pageSize).
		Msg("📦 Listing items from Softpymes")

	// Autenticar usando la URL efectiva
	token, err := c.authenticate(ctx, apiKey, apiSecret, referer, baseURL)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	var items []ItemResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetHeader("Content-Type", "application/json").
		SetQueryParam("page", strconv.Itoa(page)).
		SetQueryParam("pageSize", strconv.Itoa(pageSize)).
		SetResult(&items).
		Get(c.resolveURL(baseURL, "/app/integration/items"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Failed to list items")
		return nil, fmt.Errorf("list items request failed: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Msg("📥 Received list items response")

	if resp.IsError() {
		var errorBody map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorBody); err == nil {
			if msg, ok := errorBody["message"].(string); ok {
				c.log.Error(ctx).
					Int("status", resp.StatusCode()).
					Str("error", msg).
					Msg("❌ List items failed")
				return nil, fmt.Errorf("list items failed (status %d): %s", resp.StatusCode(), msg)
			}
		}
		return nil, fmt.Errorf("list items failed (status %d): %s", resp.StatusCode(), resp.Status())
	}

	// Mapear a tipos del dominio (parsear strings a float64)
	result := make([]ports.ListedItem, 0, len(items))
	for _, item := range items {
		price, _ := strconv.ParseFloat(item.ItemPrice, 64)
		cost, _ := strconv.ParseFloat(item.UnitCost, 64)
		result = append(result, ports.ListedItem{
			ItemCode:      item.ItemCode,
			ItemName:      item.ItemName,
			ItemPrice:     price,
			UnitCost:      cost,
			Description:   item.Description,
			MinimumStock:  item.MinimumStock,
			OrderQuantity: item.OrderQuantity,
		})
	}

	c.log.Info(ctx).
		Int("items_count", len(result)).
		Msg("✅ Items retrieved successfully")

	return result, nil
}
