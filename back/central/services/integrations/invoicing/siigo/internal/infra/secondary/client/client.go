package client

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implementa ISiigoClient para comunicarse con la API de Siigo
// La URL base se obtiene de las credenciales de cada integración, no del cliente
type Client struct {
	httpClient *httpclient.Client
	tokenCache *TokenCache
	log        log.ILogger
}

// New crea un nuevo cliente de Siigo
// La URL base se obtiene de las credenciales almacenadas en la base de datos (req.Credentials.BaseURL)
func New(logger log.ILogger) ports.ISiigoClient {
	logger.Info(context.Background()).Msg("🔍 Creating Siigo HTTP client")

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
		tokenCache: NewTokenCache(),
		log:        logger.WithModule("siigo.client"),
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
