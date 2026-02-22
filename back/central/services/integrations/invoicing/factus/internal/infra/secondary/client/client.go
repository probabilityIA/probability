package client

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implementa IFactusClient para comunicarse con la API de Factus
type Client struct {
	baseURL    string
	httpClient *httpclient.Client
	tokenCache *TokenCache
	log        log.ILogger
}

// New crea un nuevo cliente de Factus
func New(baseURL string, logger log.ILogger) ports.IFactusClient {
	logger.Info(context.Background()).
		Str("base_url", baseURL).
		Msg("🔍 Creating Factus HTTP client")

	httpConfig := httpclient.HTTPClientConfig{
		BaseURL:    baseURL,
		Timeout:    30 * time.Second,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      true,
	}

	httpClient := httpclient.New(httpConfig, logger)

	// Factus acepta JSON para el body de factura, pero form-data para auth
	// No establecemos Content-Type global ya que varía por endpoint
	httpClient.SetHeader("Accept", "application/json")

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		tokenCache: NewTokenCache(),
		log:        logger.WithModule("factus.client"),
	}
}

// endpointURL construye la URL completa para un endpoint.
// Si baseOverride está vacío, retorna solo el path (resty antepone la baseURL del cliente).
// Si baseOverride tiene valor, retorna la URL completa con ese base.
func (c *Client) endpointURL(baseOverride, path string) string {
	if baseOverride != "" {
		return strings.TrimRight(baseOverride, "/") + path
	}
	return path
}
