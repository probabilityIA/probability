package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

const meliSiteID = "MCO"

const meliCurrencyID = "COP"

const meliListingTypeID = "bronze"

func firstNonEmpty(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func (c *MeliClient) predictCategory(ctx context.Context, accessToken, siteID, title string) (string, error) {
	endpoint := fmt.Sprintf("%s/sites/%s/domain_discovery/search?limit=1&q=%s", c.baseURL, siteID, url.QueryEscape(title))
	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("meli client: category predict failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("meli client: category predict status %d: %s", resp.StatusCode, string(body))
	}
	var results []struct {
		CategoryID string `json:"category_id"`
	}
	if err := json.Unmarshal(body, &results); err != nil {
		return "", fmt.Errorf("meli client: parsing category predict: %w", err)
	}
	if len(results) == 0 || results[0].CategoryID == "" {
		return "", fmt.Errorf("meli client: no category predicted for %q", title)
	}
	return results[0].CategoryID, nil
}

func (c *MeliClient) CreateProduct(ctx context.Context, accessToken string, input domain.CreateProductInput) (string, error) {
	siteID := firstNonEmpty(input.SiteID, meliSiteID)
	currencyID := firstNonEmpty(input.CurrencyID, meliCurrencyID)
	listingTypeID := firstNonEmpty(input.ListingTypeID, meliListingTypeID)

	categoryID, err := c.predictCategory(ctx, accessToken, siteID, input.Name)
	if err != nil {
		return "", err
	}

	quantity := input.StockQuantity
	if quantity < 1 {
		quantity = 1
	}

	payload := map[string]interface{}{
		"title":              input.Name,
		"category_id":        categoryID,
		"price":              input.Price,
		"currency_id":        currencyID,
		"available_quantity": quantity,
		"buying_mode":        "buy_it_now",
		"listing_type_id":    listingTypeID,
		"condition":          "new",
		"attributes": []map[string]interface{}{
			{"id": "SELLER_SKU", "value_name": input.SKU},
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/items", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("meli client: creating item request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("meli client: create item failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", domain.ErrInvalidCredentials
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("meli client: create item status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("meli client: parsing create item response: %w", err)
	}
	if result.ID == "" {
		return "", fmt.Errorf("meli client: empty item id in create response")
	}
	return result.ID, nil
}
