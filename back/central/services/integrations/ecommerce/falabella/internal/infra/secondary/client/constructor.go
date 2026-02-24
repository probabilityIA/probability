package client

import (
	"context"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella/internal/domain"
)

// FalabellaClient implementa domain.IFalabellaClient usando la Falabella Seller Center API.
type FalabellaClient struct {
	httpClient *http.Client
}

// New crea un nuevo cliente HTTP para Falabella Seller Center.
func New() domain.IFalabellaClient {
	return &FalabellaClient{
		httpClient: &http.Client{},
	}
}

// TestConnection verifica las credenciales llamando a la Falabella Seller Center API.
// TODO: implementar llamada real a Falabella Seller Center API
func (c *FalabellaClient) TestConnection(ctx context.Context, apiKey, userID string) error {
	// TODO: Falabella Seller Center API
	// La API de Falabella usa firma HMAC con apiKey + userID + timestamp + action
	// Endpoint de prueba: GetOrder o similar para verificar credenciales
	if apiKey == "" {
		return domain.ErrMissingAPIKey
	}
	if userID == "" {
		return domain.ErrMissingUserID
	}

	return nil
}
