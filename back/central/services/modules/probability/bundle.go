package probability

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/probability/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/modules/probability/internal/infra/secondary/publisher"
	"github.com/secamc93/probability/back/central/services/modules/probability/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el modulo de probability (calculo de score de entrega).
// Es un modulo consumer-only: no tiene endpoints HTTP, solo consume eventos de RabbitMQ.
func New(database db.IDatabase, logger log.ILogger, rabbitMQ rabbitmq.IQueue) {
	if rabbitMQ == nil {
		logger.Warn().Msg("RabbitMQ no disponible, modulo de probability no se inicializara")
		return
	}

	ctx := context.Background()

	// 1. Repository
	repo := repository.New(database)

	// 2. Event Publisher
	eventPublisher := publisher.New(rabbitMQ, logger)

	// 3. Use Case
	useCase := app.New(repo, eventPublisher, logger)

	// 4. Consumer (start in background)
	scoreConsumer := consumer.New(rabbitMQ, logger, useCase)
	go func() {
		if err := scoreConsumer.Start(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("Error iniciando consumer de probability score")
		}
	}()

	logger.Info(ctx).Msg("Modulo de probability (score) inicializado - consumiendo orders.events.score")
}

