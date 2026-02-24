package shipments

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/primary/handlers"
	queueconsumer "github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/queue"
	shipmentsredis "github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo de shipments
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redis.IRedis) {
	// 1. Init Repositories
	repo := repository.New(database)

	// 2. Init Use Cases
	uc := usecases.New(repo)

	// 3. Transport Request Publisher (async via queue)
	transportPub := queue.NewTransportRequestPublisher(rabbitMQ, logger)

	// 4. Init SSE Publisher (Redis Pub/Sub for real-time notifications)
	ssePublisher := shipmentsredis.NewNoopSSEPublisher()
	if redisClient != nil {
		redisClient.RegisterChannel(redis.ChannelShipmentsEvents)
		ssePublisher = shipmentsredis.NewSSEPublisher(redisClient, logger, redis.ChannelShipmentsEvents)
	}

	// 5. Transport Response Consumer
	if rabbitMQ != nil {
		responseConsumer := queueconsumer.NewResponseConsumer(rabbitMQ, repo, logger, ssePublisher)
		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting transport response consumer in background...")
			if err := responseConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Transport response consumer failed")
			}
		}()
	}

	// 6. Init Handlers (repo satisfies ICarrierResolver via GetActiveShippingCarrier)
	h := handlers.New(uc, transportPub, repo)

	// 7. Register Routes
	h.RegisterRoutes(router)
}
