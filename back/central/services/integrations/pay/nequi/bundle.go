package nequi

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/infra/secondary/client"
	nequiqueue "github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el módulo de integración Nequi
func New(
	config env.IConfig,
	logger log.ILogger,
	database db.IDatabase,
	rabbit rabbitmq.IQueue,
) {
	logger = logger.WithModule("nequi")
	ctx := context.Background()

	encryptionKey := config.Get("ENCRYPTION_KEY")
	if encryptionKey == "" {
		logger.Warn(ctx).Msg("ENCRYPTION_KEY not set - using default (INSECURE)")
		encryptionKey = "default-encryption-key-change-me-in-production"
	}

	// 1. Cliente HTTP Nequi
	nequiClient := client.New(logger)

	// 2. Repository para credenciales
	integrationRepo := repository.New(database, logger, encryptionKey)

	// 3. Response Publisher
	responsePublisher := nequiqueue.New(rabbit, logger)

	// 4. Use Case
	useCase := app.New(nequiClient, integrationRepo, responsePublisher, logger)

	// 5. Consumer
	if rabbit != nil {
		nequiConsumer := consumer.New(rabbit, useCase, logger)
		go func() {
			if err := nequiConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("Nequi consumer failed")
			}
		}()
		logger.Info(ctx).Msg("Nequi consumer started")
	} else {
		logger.Warn(ctx).Msg("RabbitMQ not available - Nequi consumer disabled")
	}
}
