package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type paymentResponse struct {
	OrderID int64 `json:"order_id"`
	Order   struct {
		ID int64 `json:"id"`
	} `json:"order"`
}

func (c *MeliClient) GetPaymentOrderID(ctx context.Context, accessToken string, paymentID int64) (int64, error) {
	endpoint := fmt.Sprintf("%s/payments/%d", c.baseURL, paymentID)

	resp, body, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	})
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return 0, domain.ErrTokenExpired
	}
	if resp.StatusCode == http.StatusNotFound {
		return 0, domain.ErrOrderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("meli client: payment status %d: %s", resp.StatusCode, string(body))
	}

	var parsed paymentResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, fmt.Errorf("meli client: parsing payment: %w", err)
	}

	if parsed.OrderID > 0 {
		return parsed.OrderID, nil
	}
	if parsed.Order.ID > 0 {
		return parsed.Order.ID, nil
	}
	return 0, domain.ErrOrderNotFound
}
