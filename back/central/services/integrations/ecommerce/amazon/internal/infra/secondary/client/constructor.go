package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/amazon/internal/domain"
)

// AmazonClient implementa domain.IAmazonClient usando la Amazon SP-API.
type AmazonClient struct {
	httpClient *http.Client
}

// New crea un nuevo cliente HTTP para Amazon SP-API.
func New() domain.IAmazonClient {
	return &AmazonClient{
		httpClient: &http.Client{},
	}
}

// TestConnection verifica las credenciales de Amazon SP-API.
//
// TODO: Implementar el flujo OAuth completo de Amazon SP-API:
// 1. Usar refresh_token + client_id + client_secret para obtener un access_token
//    POST https://api.amazon.com/auth/o2/token
//    grant_type=refresh_token&refresh_token=X&client_id=Y&client_secret=Z
// 2. Usar el access_token para llamar a GET /sellers/v1/marketplaceParticipations
//    con el header x-amz-access-token
// 3. Verificar que el seller_id aparece en las participaciones
//
// Referencia: https://developer-docs.amazon.com/sp-api/docs/connecting-to-the-selling-partner-api
func (c *AmazonClient) TestConnection(ctx context.Context, sellerID, refreshToken, clientID, clientSecret string) error {
	// Amazon SP-API authentication is complex (OAuth + IAM role assumption).
	// This is a stub that will be implemented when the full SP-API integration is built.
	return fmt.Errorf("amazon: TestConnection not yet implemented - SP-API OAuth flow pending")
}
