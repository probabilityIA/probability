package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
)

// BankAccountResponse representa una cuenta bancaria de Softpymes.
type BankAccountResponse struct {
	AccountNumber string `json:"accountNumber"`
	Name          string `json:"name"`
	NameType      string `json:"nameType"`
}

// ListBankAccounts lista cuentas bancarias de Softpymes.
// Endpoint: GET /app/integration/bank_accounts?branchCode=XXX
func (c *Client) ListBankAccounts(ctx context.Context, apiKey, apiSecret, referer, baseURL, branchCode string) ([]ports.BankAccount, error) {
	c.log.Info(ctx).
		Str("branch_code", branchCode).
		Msg("🏦 Listing bank accounts from Softpymes")

	// Autenticar usando la URL efectiva
	token, err := c.authenticate(ctx, apiKey, apiSecret, referer, baseURL)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	var accounts []BankAccountResponse

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Referer", referer).
		SetHeader("Content-Type", "application/json").
		SetQueryParam("branchCode", branchCode).
		SetResult(&accounts).
		Get(c.resolveURL(baseURL, "/app/integration/bank_accounts"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to list bank accounts")
		return nil, fmt.Errorf("list bank accounts request failed: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Msg("Received list bank accounts response")

	if resp.IsError() {
		var errorBody map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &errorBody); err == nil {
			if msg, ok := errorBody["message"].(string); ok {
				c.log.Error(ctx).
					Int("status", resp.StatusCode()).
					Str("error", msg).
					Msg("List bank accounts failed")
				return nil, fmt.Errorf("list bank accounts failed (status %d): %s", resp.StatusCode(), msg)
			}
		}
		return nil, fmt.Errorf("list bank accounts failed (status %d): %s", resp.StatusCode(), resp.Status())
	}

	// Mapear a tipos del dominio
	result := make([]ports.BankAccount, 0, len(accounts))
	for _, acc := range accounts {
		result = append(result, ports.BankAccount{
			AccountNumber: acc.AccountNumber,
			Name:          acc.Name,
			NameType:      acc.NameType,
		})
	}

	c.log.Info(ctx).
		Int("accounts_count", len(result)).
		Msg("Bank accounts retrieved successfully")

	return result, nil
}
