package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/pay/stripe/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TODO: Stripe API docs: https://stripe.com/docs/api
// Base URL: https://api.stripe.com/v1
// Auth: Bearer <secret_key>  (Basic auth: secret_key como usuario)
// Endpoint crear PaymentIntent: POST /payment_intents

// StripeClient implementa ports.IStripeClient
type StripeClient struct {
	log log.ILogger
}

// New crea una nueva instancia del cliente Stripe
func New(logger log.ILogger) ports.IStripeClient {
	return &StripeClient{
		log: logger.WithModule("stripe.client"),
	}
}

// CreatePaymentIntent crea un PaymentIntent en Stripe
// TODO: implementar llamada real a la API de Stripe
func (c *StripeClient) CreatePaymentIntent(ctx context.Context, config *ports.StripeConfig, amount float64, currency, reference, description string) (string, string, error) {
	c.log.Warn(ctx).Msg("Stripe client not yet implemented")
	return "", "", fmt.Errorf("stripe integration not yet implemented - pending API credentials")
}
