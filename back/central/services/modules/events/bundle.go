package events

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/primary"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/events"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo de eventos adaptado a Gin y con soporte para eventos de órdenes
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, redisClient redisclient.IRedis) {
	// 1. Obtener canal Redis desde variable de entorno
	redisChannel := environment.Get("REDIS_ORDER_EVENTS_CHANNEL")

	// 2. Init Event Manager (para SSE y eventos en tiempo real)
	eventManager := events.New(logger)

	// 3. Init Repositories
	notificationConfigRepo := repository.New(database)

	// 4. Init Redis Subscriber (consumidor de eventos de órdenes)
	orderEventSubscriber := redis.New(redisClient, logger, redisChannel)

	// 5. Init Order Event Consumer
	orderEventConsumer := app.New(
		orderEventSubscriber,
		eventManager,
		notificationConfigRepo,
		logger,
	)

	// 6. Iniciar consumidor Redis en background
	go func() {
		ctx := context.Background()
		if err := orderEventConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Str("channel", redisChannel).
				Msg("Error al iniciar consumidor de eventos de órdenes")
		}
	}()

	// 7. Init Invoice Event Subscriber y Consumer (facturación en tiempo real)
	invoiceRedisChannel := environment.Get("REDIS_INVOICE_EVENTS_CHANNEL")
	if invoiceRedisChannel == "" {
		invoiceRedisChannel = "probability:invoicing:events"
	}

	invoiceEventSubscriber := redis.NewInvoiceEventSubscriber(redisClient, logger, invoiceRedisChannel)

	invoiceEventConsumer := app.NewInvoiceEventConsumer(
		invoiceEventSubscriber,
		eventManager,
		logger,
	)

	// 8. Iniciar consumidor Redis de facturación en background
	go func() {
		ctx := context.Background()
		if err := invoiceEventConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Str("channel", invoiceRedisChannel).
				Msg("Error al iniciar consumidor de eventos de facturación")
		}
	}()

	// 9. Init SSE Handler (adaptado a Gin)
	sseHandler := handlers.New(eventManager, logger)

	// 10. Init Routes (adaptado a Gin)
	routes := primary.New(sseHandler)

	// 9. Register Routes
	routes.RegisterRoutes(router)
}
