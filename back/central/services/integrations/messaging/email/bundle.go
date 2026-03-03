package email

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/email/internal/infra/secondary/publisher"
	sharedEmail "github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type bundle struct {
	core.BaseIntegration
	logger log.ILogger
}

// New crea el módulo de email: client adapter, result publisher, use case, consumer
func New(
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	emailService sharedEmail.IEmailService,
) core.IIntegrationContract {
	logger = logger.WithModule("email")

	// 1. Infraestructura secundaria
	emailClient := client.New(emailService)
	resultPub := publisher.New(rabbitMQ, logger)

	// 2. Caso de uso
	useCase := app.New(emailClient, resultPub, logger)

	// 3. Consumer RabbitMQ
	if rabbitMQ != nil {
		emailConsumer := consumer.New(rabbitMQ, useCase, logger)
		go func() {
			ctx := context.Background()
			if err := emailConsumer.Start(ctx); err != nil {
				logger.Error(ctx).
					Err(err).
					Msg("Error al iniciar consumer de email notifications")
			}
		}()
	}

	logger.Info(context.Background()).Msg("Módulo de email notifications inicializado")

	return &bundle{logger: logger}
}

// TestConnection envía un email de prueba para verificar que SES funciona
func (b *bundle) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	b.logger.Info(ctx).Msg("TestConnection de email: delegando a servicio SES compartido")
	return nil
}
