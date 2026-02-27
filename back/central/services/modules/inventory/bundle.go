package inventory

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/handlers"
	orderqueue "github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/primary/queue"
	syncqueue "github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/queue"
	inventorycache "github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo de inventory
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redis.IRedis) {
	// 1. Init Cache (resiliente: si redis es nil, cache no-op)
	var cache repository.IInventoryCache
	if redisClient != nil {
		cache = inventorycache.NewInventoryCache(redisClient, environment, logger)
	}

	// 2. Init Repository
	repo := repository.New(database, cache)

	// 3. Init Sync Publisher
	publisher := syncqueue.New(rabbitMQ, logger)

	// 4. Init Event Publisher (Redis SSE)
	eventPublisher := inventorycache.NewEventPublisher(redisClient, logger)

	// 5. Init Use Cases
	uc := app.New(repo, publisher, eventPublisher, logger)

	// 6. Init Handlers
	h := handlers.New(uc)

	// 7. Register Routes
	h.RegisterRoutes(router)

	// 8. Start Order Event Consumer (RabbitMQ)
	consumer := orderqueue.NewOrderConsumer(rabbitMQ, uc, logger)
	consumer.Start(context.Background())
}
