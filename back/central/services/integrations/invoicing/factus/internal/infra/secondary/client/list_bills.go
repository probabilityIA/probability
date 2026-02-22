package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/mappers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/response"
)

// ListBills consulta la lista paginada de facturas emitidas en Factus
// Endpoint: GET /v1/bills
func (c *Client) ListBills(ctx context.Context, credentials dtos.Credentials, params dtos.ListBillsParams) (*dtos.ListBillsResult, error) {
	token, err := c.authenticate(
		ctx,
		credentials.BaseURL,
		credentials.ClientID,
		credentials.ClientSecret,
		credentials.Username,
		credentials.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("factus list_bills: authentication failed: %w", err)
	}

	queryParams := map[string]string{}
	if params.Page > 1 {
		queryParams["page"] = fmt.Sprintf("%d", params.Page)
	}
	if params.PerPage > 0 {
		queryParams["filter[per_page]"] = fmt.Sprintf("%d", params.PerPage)
	}
	if params.Number != "" {
		queryParams["filter[number]"] = params.Number
	}
	if params.Prefix != "" {
		queryParams["filter[prefix]"] = params.Prefix
	}
	if params.Identification != "" {
		queryParams["filter[identification]"] = params.Identification
	}
	if params.Names != "" {
		queryParams["filter[names]"] = params.Names
	}
	if params.ReferenceCode != "" {
		queryParams["filter[reference_code]"] = params.ReferenceCode
	}
	if params.Status != nil {
		queryParams["filter[status]"] = fmt.Sprintf("%d", *params.Status)
	}
	if params.StartDate != "" {
		queryParams["filter[created_at][start_date]"] = params.StartDate
	}
	if params.EndDate != "" {
		queryParams["filter[created_at][end_date]"] = params.EndDate
	}

	logEvent := c.log.Info(ctx).
		Int("page", params.Page).
		Int("per_page", params.PerPage).
		Str("number", params.Number).
		Str("prefix", params.Prefix).
		Str("identification", params.Identification).
		Str("names", params.Names).
		Str("reference_code", params.ReferenceCode).
		Str("start_date", params.StartDate).
		Str("end_date", params.EndDate)
	if params.Status != nil {
		logEvent = logEvent.Int("status", *params.Status)
	}
	logEvent.Msg("📋 Listing Factus bills")

	var apiResp response.Bills

	req := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(&apiResp)

	if len(queryParams) > 0 {
		req = req.SetQueryParams(queryParams)
	}

	resp, err := req.Get(c.endpointURL(credentials.BaseURL, "/v1/bills"))
	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Factus list_bills request failed - network error")
		return nil, fmt.Errorf("factus list_bills request failed: %w", err)
	}

	if resp.IsError() {
		if resp.StatusCode() == 401 {
			c.tokenCache.Clear()
			return nil, fmt.Errorf("factus list_bills: authentication token expired (401)")
		}
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ Factus list_bills failed")
		return nil, fmt.Errorf("factus list_bills failed (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	result := mappers.BillsToListResult(&apiResp)

	c.log.Info(ctx).
		Int("total", result.Pagination.Total).
		Int("page", result.Pagination.CurrentPage).
		Int("last_page", result.Pagination.LastPage).
		Int("returned", len(result.Bills)).
		Msg("✅ Factus bills retrieved successfully")

	return result, nil
}
