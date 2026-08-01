package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/client/response"
)

const shopInfoFields = "shop_name,shop_rating,shop_review_count,item_sold_count,shop_id,shop_performance_value"

func (c *TikTokClient) QueryShopInfo(ctx context.Context, cred domain.Credential, shopName string, limit int) ([]domain.ShopInfo, error) {
	if limit < 1 || limit > 10 {
		limit = 10
	}

	endpoint := strings.TrimSuffix(cred.BaseURL, "/") + "/v2/research/tts/shop/?fields=" + shopInfoFields

	payload := map[string]interface{}{
		"shop_name": shopName,
		"fields":    shopInfoFields,
		"limit":     limit,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("tiktok: marshaling shop query: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("tiktok: building shop query request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cred.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tiktok: shop query request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("tiktok: reading shop query response: %w", err)
	}

	var shopResp response.ShopQueryResponse
	if err := json.Unmarshal(respBody, &shopResp); err != nil {
		return nil, fmt.Errorf("tiktok: parsing shop query response (status %d): %w", resp.StatusCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		if shopResp.Error.Code != "" && shopResp.Error.Code != "ok" {
			return nil, fmt.Errorf("tiktok: shop query error %s: %s", shopResp.Error.Code, shopResp.Error.Message)
		}
		return nil, fmt.Errorf("tiktok: shop query returned status %d", resp.StatusCode)
	}

	return shopResp.ToDomain(), nil
}
