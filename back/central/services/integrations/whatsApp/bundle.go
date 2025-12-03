package whatsapp

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/tester"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IWhatsAppBundle define la interfaz del bundle de WhatsApp
type IWhatsAppBundle interface {
	// SendMessage envía un mensaje de WhatsApp con el número de orden y número de teléfono
	SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error)
}

type bundle struct {
	wa      domain.IWhatsApp
	usecase app.IUseCaseSendMessage
}

// New crea una nueva instancia del bundle de WhatsApp y retorna la interfaz
// Si se proporciona integrationCore, registra el tester de WhatsApp
func New(config env.IConfig, integrationCore core.IIntegrationCore) IWhatsAppBundle {
	logger := log.New()
	wa := client.New(config)
	usecase := app.New(wa, logger, config)

	// Registrar tester de WhatsApp si se proporciona integrationCore
	if integrationCore != nil {
		whatsAppTester := tester.NewWhatsAppTester(logger)
		if err := integrationCore.RegisterTester(core.IntegrationTypeWhatsApp, whatsAppTester); err != nil {
			logger.Error().Err(err).Msg("Error al registrar tester de WhatsApp")
		} else {

		}
	}

	return &bundle{
		wa:      wa,
		usecase: usecase,
	}
}

// SendMessage expone el método simplificado para enviar mensajes
func (b *bundle) SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error) {
	req := domain.SendMessageRequest{
		OrderNumber: orderNumber,
		PhoneNumber: phoneNumber,
	}
	return b.usecase.SendMessage(ctx, req)
}
