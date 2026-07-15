package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (c *MeliClient) SendShipmentStatus(ctx context.Context, accessToken string, shipmentID int64, status string) error {
	payload := map[string]interface{}{"status": status}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("%s/shipments/%d", c.baseURL, shipmentID)

	resp, respBody, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequestWithBody(ctx, http.MethodPut, endpoint, accessToken, bodyBytes)
	})
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrTokenExpired
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("meli client: send shipment status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
