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

// GetUserMe obtiene los datos del usuario autenticado.
// GET https://api.mercadolibre.com/users/me
// Se usa para extraer el seller_id al conectar la integración.
func (c *MeliClient) GetUserMe(ctx context.Context, accessToken string) (*domain.MeliSeller, error) {
	endpoint := fmt.Sprintf("%s/users/me", c.baseURL)

	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("meli client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.ErrInvalidCredentials
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("meli client: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var userResp response.MeliUserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("meli client: parsing user response: %w", err)
	}

	seller := userResp.ToDomain()
	return &seller, nil
}
