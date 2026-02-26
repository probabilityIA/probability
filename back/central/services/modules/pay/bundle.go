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
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo de pagos
func New(
	router *gin.RouterGroup,
	database db.IDatabase,
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	redisClient redis.IRedis,
) {
	ctx := context.Background()
	moduleLogger := logger.WithModule("pay")

	// ═══════════════════════════════════════════════════════════════
	// 1. INFRAESTRUCTURA SECUNDARIA
	// ═══════════════════════════════════════════════════════════════

	repo := repository.New(database, moduleLogger)
	requestPublisher := payqueue.New(rabbitMQ, moduleLogger)

	var ssePublisher = payredis.NewNoopSSEPublisher()
	if redisClient != nil {
		ssePublisher = payredis.NewSSEPublisher(redisClient, moduleLogger, redis.ChannelPayEvents)
		redisClient.RegisterChannel(redis.ChannelPayEvents)
	} else {
		moduleLogger.Warn(ctx).Msg("Redis no disponible - SSE de pagos deshabilitado")
	}

	// ═══════════════════════════════════════════════════════════════
	// 2. CAPA DE APLICACIÓN
	// ═══════════════════════════════════════════════════════════════

	useCase := app.New(repo, requestPublisher, ssePublisher, moduleLogger)

	// ═══════════════════════════════════════════════════════════════
	// 3. INFRAESTRUCTURA PRIMARIA
	// ═══════════════════════════════════════════════════════════════

	handler := handlers.New(useCase, moduleLogger)
	handler.RegisterRoutes(router)

	if rabbitMQ != nil {
		consumers := consumer.NewConsumers(rabbitMQ, useCase, repo, ssePublisher, moduleLogger)

		go func() {
			if err := consumers.Response.Start(ctx); err != nil {
				moduleLogger.Error(ctx).Err(err).Msg("Error al iniciar consumer de respuestas de pagos")
			}
		}()

		go consumers.Retry.Start(ctx)

		moduleLogger.Info(ctx).Msg("Consumers de pagos iniciados: responses, retry")
	} else {
		moduleLogger.Warn(ctx).Msg("RabbitMQ no disponible - consumers de pagos deshabilitados")
	}
}
