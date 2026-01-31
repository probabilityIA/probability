package softpymes

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implementa IInvoicingProviderClient para Softpymes
type Client struct {
	baseURL    string
	httpClient *httpclient.Client
	tokenCache *TokenCache
	log        log.ILogger
}

// New crea un nuevo cliente de Softpymes
func New(baseURL string, logger log.ILogger) ports.IInvoicingProviderClient {
	// Configurar cliente HTTP usando el cliente compartido
	httpConfig := httpclient.HTTPClientConfig{
		BaseURL:    baseURL,
		Timeout:    30 * time.Second,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      false,
	}

	client := httpclient.New(httpConfig, logger)

	// Establecer headers comunes
	client.SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &Client{
		baseURL:    baseURL,
		httpClient: client,
		tokenCache: NewTokenCache(),
		log:        logger.WithModule("softpymes.client"),
	}
}
