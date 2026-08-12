package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (c *MeliClient) UpdateStock(ctx context.Context, accessToken, itemID, variantID string, quantity int) error {
	payload := map[string]interface{}{
		"available_quantity": quantity,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/items/%s", c.baseURL, itemID)
	if variantID != "" {
		endpoint = fmt.Sprintf("%s/items/%s/variations/%s", c.baseURL, itemID, variantID)
	}
	resp, respBody, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequestWithBody(ctx, http.MethodPut, endpoint, accessToken, bodyBytes)
	})
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrTokenExpired
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("meli client: update stock status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
