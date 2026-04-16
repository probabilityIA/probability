package client

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implementa IHelisaClient para comunicarse con la API de Helisa
// La URL base y credenciales se obtienen de cada integración, no del cliente
type Client struct {
	httpClient *httpclient.Client
	log        log.ILogger
}

// New crea un nuevo cliente de Helisa
// La URL base se obtiene de las credenciales almacenadas en la base de datos (req.Credentials.BaseURL)
func New(logger log.ILogger) ports.IHelisaClient {
	logger.Info(context.Background()).Msg("🔍 Creating Helisa HTTP client")

	httpConfig := httpclient.HTTPClientConfig{
		Timeout:    30 * time.Second,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      true,
	}

	httpClient := httpclient.New(httpConfig, logger)
	httpClient.SetHeader("Accept", "application/json")
	httpClient.SetHeader("Content-Type", "application/json")

	return &Client{
		httpClient: httpClient,
		log:        logger.WithModule("helisa.client"),
	}
}

// TestAuthentication verifica que las credenciales sean válidas
// TODO: Implementar autenticación real con la API de Helisa
func (c *Client) TestAuthentication(ctx context.Context, username, password, companyID, baseURL string) error {
	c.log.Warn(ctx).
		Str("username", username).
		Str("company_id", companyID).
		Bool("has_base_url", baseURL != "").
		Msg("⚠️ Helisa TestAuthentication not yet implemented")
	return fmt.Errorf("helisa: TestAuthentication not yet implemented")
}

// CreateInvoice crea una factura en Helisa
// TODO: Implementar creación de factura con la API de Helisa
func (c *Client) CreateInvoice(ctx context.Context, req *dtos.CreateInvoiceRequest) (*dtos.CreateInvoiceResult, error) {
	c.log.Warn(ctx).
		Str("order_id", req.OrderID).
		Msg("⚠️ Helisa CreateInvoice not yet implemented")
	return nil, fmt.Errorf("helisa: CreateInvoice not yet implemented")
}
