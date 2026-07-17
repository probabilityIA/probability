package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/client/response"
)

func (c *JumpsellerClient) GetStoreInfo(ctx context.Context, cred domain.Credential) (*domain.StoreInfo, error) {
	raw, err := c.do(ctx, cred, http.MethodGet, "/store/info.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var envelope response.StoreInfoEnvelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("jumpseller client: parsing store info: %w", err)
	}

	if envelope.Store.Code == "" {
		return nil, domain.ErrInvalidCredentials
	}

	info := envelope.ToDomain()
	return &info, nil
}

func (c *JumpsellerClient) TestConnection(ctx context.Context, cred domain.Credential) error {
	_, err := c.GetStoreInfo(ctx, cred)
	return err
}
