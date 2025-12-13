package orders

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseordermapping"
	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo de orders y retorna el caso de uso de mapping para integraciones
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redisclient.IRedis) domain.IOrderMappingUseCase {
	// 1. Init Repositories
	repo := repository.New(database)

	// 2. Init Event Publisher (si Redis está disponible)
	var eventPublisher domain.IOrderEventPublisher
	if redisClient != nil {
		redisChannel := environment.Get("REDIS_ORDER_EVENTS_CHANNEL")
		if redisChannel == "" {
			redisChannel = "probability:orders:events" // Valor por defecto
		}
		eventPublisher = redis.NewOrderEventPublisher(redisClient, logger, redisChannel)
		logger.Info(context.Background()).
			Str("channel", redisChannel).
			Msg("Order event publisher initialized")
	}

	// 3. Init Use Cases
	orderCRUD := usecaseorder.New(repo, eventPublisher)
	orderMapping := usecaseordermapping.New(repo, logger, eventPublisher)

	// 4. Init Handlers
	h := handlers.New(orderCRUD, orderMapping)

	// 5. Register Routes
	h.RegisterRoutes(router)

	// 6. Init RabbitMQ Consumer (si RabbitMQ está disponible)
	if rabbitMQ != nil {
		orderConsumer := queue.New(rabbitMQ, logger, orderMapping)
		go func() {
			if err := orderConsumer.Start(context.Background()); err != nil {
				logger.Error().
					Err(err).
					Msg("Order consumer stopped with error")
			}
		}()
	}

	return orderMapping
}
