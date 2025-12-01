package client

import (
	"net/http"

	httpclient "github.com/secamc93/probability/back/central/shared/client"
	"github.com/secamc93/probability/back/central/shared/env"
)

type WhatsappClient struct {
	httpClient *http.Client
	baseURL    string
}

func New(config env.IConfig) *WhatsappClient {
	return &WhatsappClient{
		httpClient: httpclient.NewHTTPClient(httpclient.HTTPClientConfig{}),
		baseURL:    config.Get("WHATSAPP_BASE_URL"),
	}
}
