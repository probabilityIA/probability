package whatsapp

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app/usecasetestconnection"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IWhatsAppBundle define la interfaz del bundle de WhatsApp
type IWhatsAppBundle interface {
	// SendMessage envía un mensaje de WhatsApp con el número de orden y número de teléfono
	SendMessage(ctx context.Context, orderNumber, phoneNumber string) (string, error)
	// TestConnection prueba la conexión (implementa core.ITestIntegration)
	TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error
}

type bundle struct {
	wa          domain.IWhatsApp
	usecase     app.IUseCaseSendMessage
	testUsecase usecasetestconnection.ITestConnectionUseCase
}

// New crea una nueva instancia del bundle de WhatsApp y retorna la interfaz
func New(config env.IConfig, logger log.ILogger) IWhatsAppBundle {
	logger.WithModule("whatsapp")
	wa := client.New(config)
	usecase := app.New(wa, logger, config)
	testUsecase := usecasetestconnection.New(config, logger)

	return &bundle{
		wa:          wa,
		usecase:     usecase,
		testUsecase: testUsecase,
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

// TestConnection prueba la conexión enviando un mensaje de prueba
func (b *bundle) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	// Factory para crear clientes de WhatsApp con configuración dinámica
	clientFactory := func(cfg env.IConfig) domain.IWhatsApp {
		return client.New(cfg)
	}

	// Delegar al caso de uso pasando los mapas directamente
	return b.testUsecase.TestConnection(ctx, config, credentials, clientFactory)
}
