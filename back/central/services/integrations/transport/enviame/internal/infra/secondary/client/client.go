package client

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/domain"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implements IEnviameClient for the Enviame API
type Client struct {
	httpClient *httpclient.Client
	log        log.ILogger
}

// New creates a new Enviame HTTP client
func New(logger log.ILogger) domain.IEnviameClient {
	logger.Info(context.Background()).Msg("🔍 Creating Enviame HTTP client")

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
		log:        logger.WithModule("enviame.client"),
	}
}

// TestAuthentication verifies that the credentials are valid
// TODO: Implement real authentication with Enviame API
func (c *Client) TestAuthentication(ctx context.Context, apiKey string) error {
	c.log.Warn(ctx).Msg("⚠️ Enviame TestAuthentication not yet implemented")
	return fmt.Errorf("enviame: TestAuthentication not yet implemented")
}

// Quote gets shipping rates
// TODO: Implement real quote with Enviame API
func (c *Client) Quote(apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	return nil, fmt.Errorf("enviame: Quote not yet implemented")
}

// Generate creates a shipment
// TODO: Implement real shipment generation with Enviame API
func (c *Client) Generate(apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	return nil, fmt.Errorf("enviame: Generate not yet implemented")
}

// Track gets tracking info
// TODO: Implement real tracking with Enviame API
func (c *Client) Track(apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
	return nil, fmt.Errorf("enviame: Track not yet implemented")
}

// Cancel cancels a shipment
// TODO: Implement real cancellation with Enviame API
func (c *Client) Cancel(apiKey string, idShipment string) (*domain.CancelResponse, error) {
	return nil, fmt.Errorf("enviame: Cancel not yet implemented")
}
