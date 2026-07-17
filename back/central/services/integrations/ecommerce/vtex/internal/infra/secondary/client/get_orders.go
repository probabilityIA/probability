package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/client/response"
)

func (c *VTEXClient) GetOrders(ctx context.Context, cred domain.Credential, page, perPage int, filters map[string]string) (*domain.VTEXOrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}

	endpoint := fmt.Sprintf("%s/api/oms/pvt/orders?page=%d&per_page=%d", baseURL(cred), page, perPage)

	for key, value := range filters {
		endpoint += fmt.Sprintf("&%s=%s", key, value)
	}

	body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		return nil, err
	}

	var listResp response.VTEXOrderListAPIResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, fmt.Errorf("vtex client: parsing response: %w", err)
	}

	result := listResp.ToDomain()
	return &result, nil
}
