package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

const meliAPIBaseURL = "https://api.mercadolibre.com"

// MeliClient implementa domain.IMeliClient usando la API REST de MercadoLibre.
type MeliClient struct {
	httpClient *http.Client
	baseURL    string
}

// New crea un nuevo cliente HTTP para MercadoLibre.
func New() domain.IMeliClient {
	return &MeliClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: meliAPIBaseURL,
	}
}

// newAuthorizedRequest crea un request HTTP con el Bearer token de MeLi.
func (c *MeliClient) newAuthorizedRequest(ctx context.Context, method, url, accessToken string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("meli client: creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	return req, nil
}

// TestConnection verifica que el access_token sea válido llamando a GET /users/me.
func (c *MeliClient) TestConnection(ctx context.Context, accessToken string) error {
	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, fmt.Sprintf("%s/users/me", c.baseURL), accessToken)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("meli client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("meli client: unexpected status %d", resp.StatusCode)
	}

	return nil
}
