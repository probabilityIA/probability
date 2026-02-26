package client

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/pay/melipago/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TODO: MercadoPago API docs: https://www.mercadopago.com.co/developers/es/reference
// Base URL sandbox:    https://api.mercadopago.com (con access_token de prueba)
// Base URL producción: https://api.mercadopago.com
// Auth: Bearer <access_token>
// Endpoint crear preferencia: POST /checkout/preferences

// MeliPagoClient implementa ports.IMeliPagoClient
type MeliPagoClient struct {
	log log.ILogger
}

// New crea una nueva instancia del cliente MercadoPago
func New(logger log.ILogger) ports.IMeliPagoClient {
	return &MeliPagoClient{
		log: logger.WithModule("melipago.client"),
	}
}

// CreatePreference crea una preferencia de pago en MercadoPago
// TODO: implementar llamada real a la API de MercadoPago
func (c *MeliPagoClient) CreatePreference(ctx context.Context, config *ports.MeliPagoConfig, amount float64, currency, reference, description string) (string, string, error) {
	c.log.Warn(ctx).Msg("MercadoPago client not yet implemented")
	return "", "", fmt.Errorf("melipago integration not yet implemented - pending API credentials")
}
