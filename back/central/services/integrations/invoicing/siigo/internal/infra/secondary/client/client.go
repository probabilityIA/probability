package client

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Client struct {
	httpClient *httpclient.Client
	tokenStore *TokenStore
	log        log.ILogger
}

func New(logger log.ILogger) ports.ISiigoClient {
	logger.Info(context.Background()).Msg("Creating Siigo HTTP client")

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
		tokenStore: NewTokenStore(),
		log:        logger.WithModule("siigo.client"),
	}
}

func (c *Client) endpointURL(baseOverride, path string) string {
	if baseOverride != "" {
		return strings.TrimRight(baseOverride, "/") + path
	}
	return path
}
