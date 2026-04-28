package pay

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/queue/consumer"
	payqueue "github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/secondary/queue"
	payredis "github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/services/integrations/core"
)

func New(
	router *gin.RouterGroup,
	database db.IDatabase,
	logger log.ILogger,
	config env.IConfig,
	rabbitMQ rabbitmq.IQueue,
	redisClient redis.IRedis,
	integrationCore core.IIntegrationCore,
) {
	ctx := context.Background()
	moduleLogger := logger.WithModule("pay")

	repo := repository.New(database, moduleLogger, integrationCore)
	requestPublisher := payqueue.New(rabbitMQ, moduleLogger)

	var ssePublisher = payredis.NewNoopSSEPublisher()
	if redisClient != nil {
		ssePublisher = payredis.NewSSEPublisher(redisClient, moduleLogger, redis.ChannelPayEvents)
		redisClient.RegisterChannel(redis.ChannelPayEvents)
	} else {
		moduleLogger.Warn(ctx).Msg("Redis no disponible - SSE de pagos deshabilitado")
	}

	// 2. CAPA DE APLICACIÓN

	useCase := app.New(repo, requestPublisher, ssePublisher, rabbitMQ, config, moduleLogger)
	walletUC := app.NewWalletUseCase(repo, useCase, config, moduleLogger)

	// 3. INFRAESTRUCTURA PRIMARIA

	handler := handlers.New(useCase, moduleLogger)
	handler.RegisterRoutes(router)

	walletHandler := handlers.NewWalletHandler(walletUC, moduleLogger)
	walletHandler.RegisterWalletRoutes(router)

	if rabbitMQ != nil {
		consumers := consumer.NewConsumers(rabbitMQ, useCase, repo, ssePublisher, moduleLogger)

		go func() {
			if err := consumers.Response.Start(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("Error al iniciar consumer de respuestas de pagos")
			}
		}()

		go consumers.Retry.Start(ctx)

		go func() {
			if err := consumers.BoldWebhook.Start(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("Error al iniciar consumer de bold webhook events")
			}
		}()

		moduleLogger.Info(ctx).Msg("Consumers de pagos iniciados: responses, retry, bold_webhook")
	} else {
		moduleLogger.Warn(ctx).Msg("RabbitMQ no disponible - consumers de pagos deshabilitados")
	}
}
