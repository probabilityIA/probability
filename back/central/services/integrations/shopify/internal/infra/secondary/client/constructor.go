package client

import (
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/httpclient"
)

type shopifyClient struct {
	httpClient *http.Client
}

func New() domain.ShopifyClient {
	return &shopifyClient{
		httpClient: httpclient.NewHTTPClient(httpclient.HTTPClientConfig{
			Timeout:         30 * time.Second,
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		}),
	}
}
