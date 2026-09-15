package openrouter

import (
	"net/http"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type Client struct {
	logger        log.ILogger
	httpClient    *http.Client
	transportData map[string]interface{}
}

func New(logger log.ILogger) ports.IRecommendationProvider {
	client := &Client{
		logger:     logger,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
	client.loadTransportData()
	return client
}
