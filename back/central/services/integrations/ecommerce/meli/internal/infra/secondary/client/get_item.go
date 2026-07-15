package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type itemResponse struct {
	ID                string   `json:"id"`
	Status            string   `json:"status"`
	SubStatus         []string `json:"sub_status"`
	AvailableQuantity int      `json:"available_quantity"`
	UserProductID     string   `json:"user_product_id"`
	CatalogListing    bool     `json:"catalog_listing"`
	Shipping          struct {
		LogisticType string   `json:"logistic_type"`
		Tags         []string `json:"tags"`
	} `json:"shipping"`
	Attributes []itemAttributeResp `json:"attributes"`
	Variations []struct {
		ID                int64               `json:"id"`
		UserProductID     string              `json:"user_product_id"`
		AvailableQuantity int                 `json:"available_quantity"`
		SellerCustomField string              `json:"seller_custom_field"`
		Attributes        []itemAttributeResp `json:"attributes"`
	} `json:"variations"`
}

type itemAttributeResp struct {
	ID        string `json:"id"`
	ValueName string `json:"value_name"`
}

func sellerSKUFromAttrs(sellerCustom string, attrs []itemAttributeResp) string {
	if sellerCustom != "" {
		return sellerCustom
	}
	for _, a := range attrs {
		if a.ID == "SELLER_SKU" {
			return a.ValueName
		}
	}
	return ""
}

func (c *MeliClient) GetItem(ctx context.Context, accessToken, itemID string) (*domain.MeliItemDetail, error) {
	endpoint := fmt.Sprintf("%s/items/%s?include_attributes=all", c.baseURL, itemID)

	resp, body, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrTokenExpired
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrItemNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: item status %d: %s", resp.StatusCode, string(body))
	}

	var parsed itemResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("meli client: parsing item: %w", err)
	}

	detail := &domain.MeliItemDetail{
		ID:                parsed.ID,
		Status:            parsed.Status,
		SubStatus:         parsed.SubStatus,
		AvailableQuantity: parsed.AvailableQuantity,
		UserProductID:     parsed.UserProductID,
		LogisticType:      parsed.Shipping.LogisticType,
		Tags:              parsed.Shipping.Tags,
		CatalogListing:    parsed.CatalogListing,
		SellerSKU:         sellerSKUFromAttrs("", parsed.Attributes),
	}

	for _, v := range parsed.Variations {
		detail.Variations = append(detail.Variations, domain.MeliItemVariation{
			ID:                v.ID,
			UserProductID:     v.UserProductID,
			AvailableQuantity: v.AvailableQuantity,
			SellerSKU:         sellerSKUFromAttrs(v.SellerCustomField, v.Attributes),
		})
	}

	return detail, nil
}
