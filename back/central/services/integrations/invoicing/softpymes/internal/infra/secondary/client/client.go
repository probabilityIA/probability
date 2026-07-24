package client

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Client struct {
	httpClient *httpclient.Client
	tokenCache *TokenCache
	log        log.ILogger
}

func New(logger log.ILogger) ports.ISoftpymesClient {
	httpConfig := httpclient.HTTPClientConfig{
		Timeout:   90 * time.Second,
		RetryWait: 3 * time.Second,
		Debug:     true,
	}

	httpClient := httpclient.New(httpConfig, logger)
	httpClient.GetRestyClient().SetRetryCount(0)
	httpClient.SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &Client{
		httpClient: httpClient,
		tokenCache: NewTokenCache(),
		log:        logger.WithModule("softpymes.client"),
	}
}

func (c *Client) resolveURL(baseURL, path string) string {
	return baseURL + path
}
