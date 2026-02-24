package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/client/response"
)

// GetOrders obtiene órdenes paginadas del vendedor.
// GET https://api.mercadolibre.com/orders/search?seller={seller_id}&...
// MeLi usa offset/limit (no page), máximo 50 por request.
// Retorna las órdenes tipadas, los bytes crudos por orden, y error.
func (c *MeliClient) GetOrders(ctx context.Context, accessToken string, sellerID int64, params *domain.GetOrdersParams) (*domain.GetOrdersResult, [][]byte, error) {
	queryStr := ""
	if params != nil {
		queryStr = params.ToQueryString(sellerID)
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
