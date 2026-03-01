package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/client/response"
)

// buildOrdersQueryString construye los query params para la API de MeLi.
// El sellerID es obligatorio para buscar órdenes.
// Esta lógica vive en infra porque usa net/url, que es un detalle de transporte HTTP.
func buildOrdersQueryString(sellerID int64, p *domain.GetOrdersParams) string {
	params := url.Values{}
	params.Set("seller", fmt.Sprintf("%d", sellerID))

	if p.Status != "" {
		params.Set("order.status", p.Status)
	}
	if p.DateFrom != nil {
		params.Set("order.date_created.from", p.DateFrom.Format(time.RFC3339))
	}
	if p.DateTo != nil {
		params.Set("order.date_created.to", p.DateTo.Format(time.RFC3339))
	}
	if p.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", p.Offset))
	}
	if p.Limit > 0 {
		limit := p.Limit
		if limit > 50 {
			limit = 50
		}
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	if p.Sort != "" {
		params.Set("sort", p.Sort)
	}

	return params.Encode()
}

// GetOrders obtiene órdenes paginadas del vendedor.
// GET https://api.mercadolibre.com/orders/search?seller={seller_id}&...
// MeLi usa offset/limit (no page), máximo 50 por request.
// Retorna las órdenes tipadas, los bytes crudos por orden, y error.
func (c *MeliClient) GetOrders(ctx context.Context, accessToken string, sellerID int64, params *domain.GetOrdersParams) (*domain.GetOrdersResult, [][]byte, error) {
	queryStr := ""
	if params != nil {
		queryStr = buildOrdersQueryString(sellerID, params)
	} else {
		queryStr = fmt.Sprintf("seller=%d", sellerID)
	}

	endpoint := fmt.Sprintf("%s/orders/search?%s", c.baseURL, queryStr)

	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("meli client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, nil, domain.ErrTokenExpired
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, nil, domain.ErrRateLimited
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("meli client: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("meli client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp response.MeliOrdersSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, nil, fmt.Errorf("meli client: parsing response: %w", err)
	}

	// Convertir cada orden a dominio y preservar bytes crudos
	orders := make([]domain.MeliOrder, 0, len(searchResp.Results))
	rawBytes := make([][]byte, 0, len(searchResp.Results))

	for _, orderResp := range searchResp.Results {
		orders = append(orders, orderResp.ToDomain())
		// Serializar de vuelta cada orden individual para RawData
		rawOrder, err := json.Marshal(orderResp)
		if err != nil {
			continue
		}
		rawBytes = append(rawBytes, rawOrder)
	}

	return &domain.GetOrdersResult{
		Orders: orders,
		Total:  searchResp.Paging.Total,
		Offset: searchResp.Paging.Offset,
		Limit:  searchResp.Paging.Limit,
	}, rawBytes, nil
}
