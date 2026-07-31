package client

import (
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type TikTokClient struct {
	httpClient *http.Client
	logger     log.ILogger
}

func New(logger log.ILogger) domain.ITikTokClient {
	return &TikTokClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger.WithModule("tiktok.client"),
	}
}
