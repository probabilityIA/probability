package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/client/response"
)

func (c *TiendanubeClient) GetStoreInfo(ctx context.Context, cred domain.Credential) (*domain.StoreInfo, error) {
	raw, _, err := c.do(ctx, cred, http.MethodGet, "/store", nil, nil)
	if err != nil {
		return nil, err
	}

	var store response.Store
	if err := json.Unmarshal(raw, &store); err != nil {
		return nil, fmt.Errorf("tiendanube client: parsing store: %w", err)
	}

	return &domain.StoreInfo{
		ID:       store.ID,
		Name:     store.Name.First(),
		URL:      store.URL,
		Country:  store.Country,
		Currency: store.Currency,
		Language: store.MainLanguage,
	}, nil
}

func (c *TiendanubeClient) TestConnection(ctx context.Context, cred domain.Credential) error {
	_, err := c.GetStoreInfo(ctx, cred)
	return err
}
