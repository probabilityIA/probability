package client

import (
	"context"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito/internal/domain"
)

// ExitoClient implementa domain.IExitoClient usando la API del marketplace Exito.
type ExitoClient struct {
	httpClient *http.Client
}

// New crea un nuevo cliente HTTP para Exito.
func New() domain.IExitoClient {
	return &ExitoClient{
		httpClient: &http.Client{},
	}
}

// TestConnection verifica las credenciales contra la API de Exito.
//
// TODO: implementar llamada real a la API del marketplace Exito
// para validar api_key y seller_id.
func (c *ExitoClient) TestConnection(ctx context.Context, apiKey, sellerID string) error {
	// Stub: la implementacion real debera llamar al endpoint de
	// autenticacion/validacion del marketplace Exito.
	_ = apiKey
	_ = sellerID
	return nil
}
