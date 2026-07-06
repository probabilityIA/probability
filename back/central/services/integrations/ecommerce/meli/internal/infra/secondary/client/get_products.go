package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type meliItemsSearchResponse struct {
	Results []string `json:"results"`
	Paging  struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"paging"`
}

type meliItemDetail struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	Price             float64 `json:"price"`
	AvailableQuantity int     `json:"available_quantity"`
	SellerCustomField string  `json:"seller_custom_field"`
	Attributes        []struct {
		ID        string `json:"id"`
		ValueName string `json:"value_name"`
	} `json:"attributes"`
}

type meliMultigetItem struct {
	Code int            `json:"code"`
	Body meliItemDetail `json:"body"`
}

func extractSKU(item meliItemDetail) string {
	if strings.TrimSpace(item.SellerCustomField) != "" {
		return item.SellerCustomField
	}
	for _, attr := range item.Attributes {
		if attr.ID == "SELLER_SKU" {
			return attr.ValueName
		}
	}
	return ""
}

func (c *MeliClient) GetProducts(ctx context.Context, accessToken string, sellerID int64) ([]domain.MeliProduct, error) {
	itemIDs := make([]string, 0)
	offset := 0
	limit := 50

	for {
		endpoint := fmt.Sprintf("%s/users/%d/items/search?limit=%d&offset=%d", c.baseURL, sellerID, limit, offset)
		req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
		if err != nil {
			return nil, err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("meli client: items search failed: %w", err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return nil, domain.ErrTokenExpired
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("meli client: items search status %d: %s", resp.StatusCode, string(body))
		}
		var searchResp meliItemsSearchResponse
		if err := json.Unmarshal(body, &searchResp); err != nil {
			return nil, fmt.Errorf("meli client: parsing items search: %w", err)
		}
		itemIDs = append(itemIDs, searchResp.Results...)
		offset += limit
		if len(searchResp.Results) < limit || offset >= searchResp.Paging.Total {
			break
		}
		if offset >= 1000 {
			break
		}
	}

	products := make([]domain.MeliProduct, 0, len(itemIDs))
	for i := 0; i < len(itemIDs); i += 20 {
		end := i + 20
		if end > len(itemIDs) {
			end = len(itemIDs)
		}
		batch := itemIDs[i:end]
		details, err := c.multigetItems(ctx, accessToken, batch)
		if err != nil {
			return nil, err
		}
		for _, d := range details {
			products = append(products, domain.MeliProduct{
				ID:            d.ID,
				SKU:           extractSKU(d),
				Name:          d.Title,
				Price:         d.Price,
				StockQuantity: d.AvailableQuantity,
			})
		}
	}

	return products, nil
}

func (c *MeliClient) multigetItems(ctx context.Context, accessToken string, ids []string) ([]meliItemDetail, error) {
	attrs := "id,title,price,available_quantity,seller_custom_field,attributes"
	endpoint := fmt.Sprintf("%s/items?ids=%s&attributes=%s", c.baseURL, strings.Join(ids, ","), attrs)
	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("meli client: multiget failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: multiget status %d: %s", resp.StatusCode, string(body))
	}
	var items []meliMultigetItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("meli client: parsing multiget: %w", err)
	}
	out := make([]meliItemDetail, 0, len(items))
	for _, it := range items {
		if it.Code == http.StatusOK && it.Body.ID != "" {
			out = append(out, it.Body)
		}
	}
	return out, nil
}
