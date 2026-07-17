package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/client/response"
)

func (c *VTEXClient) GetOrderByID(ctx context.Context, cred domain.Credential, orderID string) (*domain.VTEXOrder, []byte, error) {
	endpoint := fmt.Sprintf("%s/api/oms/pvt/orders/%s", baseURL(cred), orderID)

	body, err := c.do(ctx, http.MethodGet, endpoint, cred, nil)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, nil, domain.ErrOrderNotFound
		}
		return nil, nil, err
	}

	var orderResp response.VTEXOrderDetailResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, nil, fmt.Errorf("vtex client: parsing response: %w", err)
	}

	order := orderResp.ToDomain()
	return &order, body, nil
}
