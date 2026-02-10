package client

import (
<<<<<<< HEAD
=======
	"context"
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implementa ISoftpymesClient para comunicarse con la API de Softpymes
type Client struct {
	baseURL    string
	httpClient *httpclient.Client
	tokenCache *TokenCache
	log        log.ILogger
}

// New crea un nuevo cliente de Softpymes
func New(baseURL string, logger log.ILogger) ports.ISoftpymesClient {
<<<<<<< HEAD
=======
	logger.Info(context.Background()).
		Str("base_url", baseURL).
		Msg("🔍 DEBUG: Creating Softpymes HTTP client with URL")

>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	// Configurar cliente HTTP usando el cliente compartido
	httpConfig := httpclient.HTTPClientConfig{
		BaseURL:    baseURL,
		Timeout:    30 * time.Second,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      true, // ✅ Debug habilitado para ver todas las peticiones HTTP
	}

	httpClient := httpclient.New(httpConfig, logger)

	// Establecer headers comunes
	httpClient.SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		tokenCache: NewTokenCache(),
		log:        logger.WithModule("softpymes.client"),
	}
}
