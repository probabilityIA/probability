package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/pay/wompi/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TODO: Wompi API docs: https://docs.wompi.co/
// Base URL sandbox:    https://sandbox.wompi.co/v1
// Base URL producción: https://production.wompi.co/v1
// Auth: Bearer <private_key>
// Endpoint crear transacción: POST /transactions

// WompiClient implementa ports.IWompiClient
type WompiClient struct {
	log log.ILogger
}

// New crea una nueva instancia del cliente Wompi
func New(logger log.ILogger) ports.IWompiClient {
	return &WompiClient{
		log: logger.WithModule("wompi.client"),
	}
}

// CreateTransaction crea una transacción en Wompi
// TODO: implementar llamada real a la API de Wompi
func (c *WompiClient) CreateTransaction(ctx context.Context, config *ports.WompiConfig, amount float64, currency, reference, description string) (string, string, error) {
	c.log.Warn(ctx).Msg("Wompi client not yet implemented")
	return "", "", fmt.Errorf("wompi integration not yet implemented - pending API credentials")
}
