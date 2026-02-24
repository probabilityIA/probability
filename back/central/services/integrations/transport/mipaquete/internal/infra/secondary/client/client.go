package client

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/transport/mipaquete/internal/domain"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

// Client implements IMiPaqueteClient for the MiPaquete API
type Client struct {
	httpClient *httpclient.Client
	log        log.ILogger
}

// New creates a new MiPaquete HTTP client
func New(logger log.ILogger) domain.IMiPaqueteClient {
	logger.Info(context.Background()).Msg("🔍 Creating MiPaquete HTTP client")

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
		log:        logger.WithModule("mipaquete.client"),
	}
}

// TestAuthentication verifies that the credentials are valid
// TODO: Implement real authentication with MiPaquete API
func (c *Client) TestAuthentication(ctx context.Context, apiKey string) error {
	c.log.Warn(ctx).Msg("⚠️ MiPaquete TestAuthentication not yet implemented")
	return fmt.Errorf("mipaquete: TestAuthentication not yet implemented")
}

// Quote gets shipping rates
// TODO: Implement real quote with MiPaquete API
func (c *Client) Quote(apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
	return nil, fmt.Errorf("mipaquete: Quote not yet implemented")
}

// Generate creates a shipment
// TODO: Implement real shipment generation with MiPaquete API
func (c *Client) Generate(apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
	return nil, fmt.Errorf("mipaquete: Generate not yet implemented")
}

// Track gets tracking info
// TODO: Implement real tracking with MiPaquete API
func (c *Client) Track(apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
	return nil, fmt.Errorf("mipaquete: Track not yet implemented")
}

// Cancel cancels a shipment
// TODO: Implement real cancellation with MiPaquete API
func (c *Client) Cancel(apiKey string, idShipment string) (*domain.CancelResponse, error) {
	return nil, fmt.Errorf("mipaquete: Cancel not yet implemented")
}
