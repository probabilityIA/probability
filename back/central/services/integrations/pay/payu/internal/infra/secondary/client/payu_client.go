package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/pay/payu/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TODO: PayU API docs: https://developers.payulatam.com/
// Base URL sandbox:    https://sandbox.api.payulatam.com/payments-api/4.0/service.cgi
// Base URL producción: https://api.payulatam.com/payments-api/4.0/service.cgi
// Auth: Authorization header con apiLogin:apiKey en base64
// Endpoint: POST (body con command: "SUBMIT_TRANSACTION")

// PayUClient implementa ports.IPayUClient
type PayUClient struct {
	log log.ILogger
}

// New crea una nueva instancia del cliente PayU
func New(logger log.ILogger) ports.IPayUClient {
	return &PayUClient{
		log: logger.WithModule("payu.client"),
	}
}

// CreateTransaction crea una transacción en PayU
// TODO: implementar llamada real a la API de PayU
func (c *PayUClient) CreateTransaction(ctx context.Context, config *ports.PayUConfig, amount float64, currency, reference, description string) (string, string, error) {
	c.log.Warn(ctx).Msg("PayU client not yet implemented")
	return "", "", fmt.Errorf("payu integration not yet implemented - pending API credentials")
}
