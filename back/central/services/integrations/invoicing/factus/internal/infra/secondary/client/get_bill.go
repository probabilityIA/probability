package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/response"
)

// GetBillByNumber obtiene el detalle completo de una factura por su número
// Endpoint: GET /v1/bills/show/:number
func (c *Client) GetBillByNumber(ctx context.Context, credentials dtos.Credentials, number string) (*dtos.BillDetail, error) {
	token, err := c.authenticate(
		ctx,
		credentials.BaseURL,
		credentials.ClientID,
		credentials.ClientSecret,
		credentials.Username,
		credentials.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("factus get_bill: authentication failed: %w", err)
	}

	c.log.Info(ctx).
		Str("number", number).
		Msg("🔍 Getting Factus bill by number")

	var apiResp response.GetBillDetail

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(&apiResp).
		Get(c.endpointURL(credentials.BaseURL, fmt.Sprintf("/v1/bills/show/%s", number)))
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Factus get_bill request failed - network error")
		return nil, fmt.Errorf("factus get_bill request failed: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return nil, fmt.Errorf("factus get_bill: authentication token expired (401)")
		}
		if resp.StatusCode() == 404 {
			return nil, fmt.Errorf("factus get_bill: bill not found (number=%s)", number)
		}
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ Factus get_bill failed")
		return nil, fmt.Errorf("factus get_bill failed (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	result := mappers.GetBillToDetail(&apiResp)

	c.log.Info(ctx).
		Str("number", result.Number).
		Str("cufe", result.CUFE).
		Str("total", result.Total).
		Msg("✅ Factus bill retrieved successfully")

	return result, nil
}
