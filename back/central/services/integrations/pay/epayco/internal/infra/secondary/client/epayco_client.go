package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/pay/epayco/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TODO: ePayco API docs: https://docs.epayco.co/
// Base URL: https://api.epayco.co
// Auth: Bearer <public_key> + p_key en parámetros
// Endpoint crear pago: POST /payment/process

// EPaycoClient implementa ports.IEPaycoClient
type EPaycoClient struct {
	log log.ILogger
}

// New crea una nueva instancia del cliente ePayco
func New(logger log.ILogger) ports.IEPaycoClient {
	return &EPaycoClient{
		log: logger.WithModule("epayco.client"),
	}
}

// CreateCheckout crea un checkout en ePayco
// TODO: implementar llamada real a la API de ePayco
func (c *EPaycoClient) CreateCheckout(ctx context.Context, config *ports.EPaycoConfig, amount float64, currency, reference, description string) (string, string, error) {
	c.log.Warn(ctx).Msg("ePayco client not yet implemented")
	return "", "", fmt.Errorf("epayco integration not yet implemented - pending API credentials")
}
