package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/response"
)

// GetCustomerByIdentification busca un cliente en Siigo por número de identificación
// Endpoint: GET /v1/customers?identification=xxx
func (c *Client) GetCustomerByIdentification(ctx context.Context, credentials dtos.Credentials, identification string) (*dtos.CustomerResult, error) {
	c.log.Info(ctx).
		Str("identification", identification).
		Msg("🔍 Getting Siigo customer by identification")

	// Autenticar
	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	var listResp response.ListCustomersResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetQueryParam("identification", identification).
		SetResult(&listResp).
		Get(c.endpointURL(credentials.BaseURL, "/v1/customers"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Siigo get customer request failed - network error")
		return nil, fmt.Errorf("error de red al buscar cliente en Siigo: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Int("results_count", len(listResp.Results)).
		Msg("📥 Siigo get customer response received")

	if resp.IsError() {
		c.log.Warn(ctx).
			Int("status", resp.StatusCode()).
			Str("identification", identification).
			Msg("Siigo customer not found or error")
		return nil, nil // No encontrado no es un error crítico
	}

	if len(listResp.Results) == 0 {
		c.log.Info(ctx).
			Str("identification", identification).
			Msg("Customer not found in Siigo")
		return nil, nil
	}

	customer := mappers.CustomerToDTO(&listResp.Results[0])
	c.log.Info(ctx).
		Str("customer_id", customer.ID).
		Str("identification", identification).
		Msg("✅ Siigo customer found")

	return customer, nil
}
