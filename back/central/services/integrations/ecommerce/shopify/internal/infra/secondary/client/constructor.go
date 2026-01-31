package client

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

type shopifyClient struct {
	client *resty.Client
}

func New() domain.ShopifyClient {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(0)

	return &shopifyClient{
		client: client,
	}
}

// SetDebug habilita o deshabilita el modo debug de Resty (muestra request/response completo)
func (c *shopifyClient) SetDebug(enabled bool) {
	c.client.SetDebug(enabled)
}
