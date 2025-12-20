package orders

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseordermapping"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/queue"
	redisevents "github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/repository"
	orderstatusrepo "github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo de orders y retorna el caso de uso de mapping para integraciones
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redisclient.IRedis) domain.IOrderMappingUseCase {
	// 1. Init Repositories
	repo := repository.New(database, environment)

	// 2. Init Event Publisher (si Redis está disponible)
	var eventPublisher domain.IOrderEventPublisher
	var eventConsumer redisevents.IOrderEventConsumer
	if redisClient != nil {
		redisChannel := environment.Get("REDIS_ORDER_EVENTS_CHANNEL")
		if redisChannel == "" {
			redisChannel = "probability:orders:events" // Valor por defecto
		}
		eventPublisher = redisevents.NewOrderEventPublisher(redisClient, logger, redisChannel)
		logger.Info(context.Background()).
			Str("channel", redisChannel).
			Msg("Order event publisher initialized")
	}

	// 3. Init Use Cases
	orderCRUD := usecaseorder.New(repo, eventPublisher)

	// 3.0. Init OrderStatus Repository (para mapeo de estados)
	orderStatusRepo := orderstatusrepo.New(database, logger)
	orderMapping := usecaseordermapping.New(repo, logger, eventPublisher, orderStatusRepo)

	// 3.1. Init Score Use Case
	scoreUseCase := usecaseorderscore.New(repo)

	// 3.2. Init Event Consumer (si Redis está disponible)
	if redisClient != nil && eventPublisher != nil {
		redisChannel := environment.Get("REDIS_ORDER_EVENTS_CHANNEL")
		if redisChannel == "" {
			redisChannel = "probability:orders:events" // Valor por defecto
		}
		eventConsumer = redisevents.NewOrderEventConsumer(redisClient, logger, redisChannel, scoreUseCase)
		logger.Info(context.Background()).
			Str("channel", redisChannel).
			Msg("Order event consumer initialized")
	}

	// 4. Init Handlers
	h := handlers.New(orderCRUD, orderMapping)

	// 5. Register Routes
	h.RegisterRoutes(router)

	// 6. Init RabbitMQ Consumer (si RabbitMQ está disponible)
	if rabbitMQ != nil {
		orderConsumer := queue.New(rabbitMQ, logger, orderMapping, repo)
		go func() {
			if err := orderConsumer.Start(context.Background()); err != nil {
				logger.Error().
					Err(err).
					Msg("Order consumer stopped with error")
			}
		}()
	}

	// 7. Init Redis Event Consumer (si Redis está disponible)
	if eventConsumer != nil {
		go func() {
			fmt.Printf("[Bundle] Iniciando consumer de eventos de Redis para cálculo de score...\n")
			if err := eventConsumer.Start(context.Background()); err != nil {
				logger.Error().
					Err(err).
					Msg("Order event consumer stopped with error")
			}
		}()
	}

	return orderMapping
}
