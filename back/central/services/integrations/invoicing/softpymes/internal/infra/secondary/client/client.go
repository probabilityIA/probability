package client

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implementa ISoftpymesClient para comunicarse con la API de Softpymes.
// La URL base SIEMPRE viene del tipo de integración (integration_types.base_url /
// base_url_test) y se pasa por llamada. No hay URL quemada en el constructor.
type Client struct {
	httpClient *httpclient.Client
	tokenCache *TokenCache
	log        log.ILogger
}

// New crea un nuevo cliente de Softpymes sin URL base fija.
// La URL efectiva se pasa en cada operación (viene de integration_types.base_url).
func New(logger log.ILogger) ports.ISoftpymesClient {
	httpConfig := httpclient.HTTPClientConfig{
		Timeout:    30 * time.Second,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      true,
	}

	httpClient := httpclient.New(httpConfig, logger)
	httpClient.SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &Client{
		httpClient: httpClient,
		tokenCache: NewTokenCache(),
		log:        logger.WithModule("softpymes.client"),
	}
}

// resolveURL construye la URL absoluta para una llamada HTTP.
// baseURL debe ser no-vacío (viene de integration_types.base_url o base_url_test).
func (c *Client) resolveURL(baseURL, path string) string {
	return baseURL + path
}
