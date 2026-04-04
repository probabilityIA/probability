package ai_sales

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/infra/secondary/ai_adapter"
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/infra/secondary/cache"
	configprovider "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/infra/secondary/config"
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/bedrock"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el modulo ai_sales (solo queue-based, sin HTTP handlers)
func New(database db.IDatabase, logger log.ILogger, rabbitMQ rabbitmq.IQueue, redisClient redis.IRedis, bedrockClient bedrock.IBedrock) {
	// 1. Infraestructura secundaria
	aiProvider := ai_adapter.New(bedrockClient, logger)
	sessionCache := cache.New(redisClient, logger)
	productRepo := repository.New(database, logger)
	responsePublisher := queue.NewResponsePublisher(rabbitMQ, logger)
	orderPublisher := queue.NewOrderPublisher(rabbitMQ, logger)
	configProvider := configprovider.New(redisClient, logger)

	// 2. Caso de uso
	useCase := app.New(aiProvider, sessionCache, productRepo, responsePublisher, orderPublisher, configProvider, logger)

	// 3. Consumer (infraestructura primaria)
	aiConsumer := consumer.New(rabbitMQ, useCase, logger)

	// 4. Iniciar consumer en background
	go func() {
		if err := aiConsumer.Start(context.Background()); err != nil {
			logger.Error(context.Background()).
				Err(err).
				Msg("Error iniciando AI Sales consumer")
		}
	}()

	logger.Info(context.Background()).Msg("Modulo AI Sales inicializado")
}
