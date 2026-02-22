package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

// listInvoicesResponse respuesta de Siigo para listar facturas
type listInvoicesResponse struct {
	Pagination struct {
		Page       int `json:"page"`
		PageSize   int `json:"page_size"`
		TotalPages int `json:"total_pages"`
		TotalItems int `json:"total_items"`
	} `json:"_pagination"`
	Results []struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Date     string  `json:"date"`
		Customer struct {
			Identification string `json:"identification"`
			Name           string `json:"name"`
		} `json:"customer"`
		Total  float64 `json:"total"`
		Status string  `json:"status"`
	} `json:"results"`
}

// ListInvoices consulta la lista paginada de facturas emitidas en Siigo
// Endpoint: GET /v1/invoices
func (c *Client) ListInvoices(ctx context.Context, credentials dtos.Credentials, params dtos.ListInvoicesParams) (*dtos.ListInvoicesResult, error) {
	c.log.Info(ctx).
		Int("page", params.Page).
		Int("page_size", params.PageSize).
		Msg("📋 Listing Siigo invoices")

	// Autenticar
	token, err := c.authenticate(ctx, credentials.Username, credentials.AccessKey, credentials.AccountID, credentials.PartnerID, credentials.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Siigo: %w", err)
	}

	// Valores por defecto de paginación
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req := c.httpClient.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Partner-Id", credentials.PartnerID).
		SetQueryParam("page", strconv.Itoa(page)).
		SetQueryParam("page_size", strconv.Itoa(pageSize))

	if params.DateFrom != "" {
		req = req.SetQueryParam("date_start", params.DateFrom)
	}
	if params.DateTo != "" {
		req = req.SetQueryParam("date_end", params.DateTo)
	}

	var listResp listInvoicesResponse

	resp, err := req.
		SetResult(&listResp).
		Get(c.endpointURL(credentials.BaseURL, "/v1/invoices"))

	if err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Siigo list invoices request failed - network error")
		return nil, fmt.Errorf("error de red al listar facturas en Siigo: %w", err)
	}

	c.log.Info(ctx).
		Int("status_code", resp.StatusCode()).
		Int("results_count", len(listResp.Results)).
		Msg("📥 Siigo list invoices response received")

	if resp.IsError() {
		c.log.Error(ctx).
			Int("status", resp.StatusCode()).
			Str("body", string(resp.Body())).
			Msg("❌ Siigo list invoices failed")
		return nil, fmt.Errorf("error al listar facturas en Siigo (código %d)", resp.StatusCode())
	}

	// Mapear resultados
	items := make([]dtos.InvoiceSummary, 0, len(listResp.Results))
	for _, r := range listResp.Results {
		items = append(items, dtos.InvoiceSummary{
			ID:           r.ID,
			Number:       r.Name,
			Date:         r.Date,
			CustomerName: r.Customer.Name,
			CustomerID:   r.Customer.Identification,
			Total:        r.Total,
			Status:       r.Status,
		})
	}

	return &dtos.ListInvoicesResult{
		Items:      items,
		Total:      listResp.Pagination.TotalItems,
		Page:       listResp.Pagination.Page,
		PageSize:   listResp.Pagination.PageSize,
		TotalPages: listResp.Pagination.TotalPages,
	}, nil
}
