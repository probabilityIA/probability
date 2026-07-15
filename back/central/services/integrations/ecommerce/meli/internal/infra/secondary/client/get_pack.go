package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type packResponse struct {
	ID     int64 `json:"id"`
	Orders []struct {
		ID int64 `json:"id"`
	} `json:"orders"`
}

func (c *MeliClient) GetPack(ctx context.Context, accessToken string, packID int64) (*domain.MeliPack, error) {
	endpoint := fmt.Sprintf("%s/packs/%d", c.baseURL, packID)

	resp, body, err := c.do(ctx, func() (*http.Request, error) {
		return c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrTokenExpired
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrOrderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: pack status %d: %s", resp.StatusCode, string(body))
	}

	var parsed packResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("meli client: parsing pack: %w", err)
	}

	pack := &domain.MeliPack{ID: parsed.ID}
	for _, o := range parsed.Orders {
		if o.ID > 0 {
			pack.OrderIDs = append(pack.OrderIDs, o.ID)
		}
	}
	return pack, nil
}
