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
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// Canales Redis — usar siempre las constantes de shared/redis para evitar desincronización
const (
	channelOrders      = redisclient.ChannelOrdersEvents
	channelInvoicing   = redisclient.ChannelInvoicingEvents
	channelIntegration = redisclient.ChannelIntegrationsSyncOrders
)

// New inicializa el módulo de eventos adaptado a Gin y con soporte para eventos de órdenes
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, redisClient redisclient.IRedis) {
	// 1. Init Event Manager (para SSE y eventos en tiempo real)
	eventManager := events.New(logger)

	// 2. Init Repositories
	notificationConfigRepo := repository.New(database)

	// 3. Init Redis Subscriber (consumidor de eventos de órdenes)
	orderEventSubscriber := redis.New(redisClient, logger, channelOrders)

	// 4. Init Order Event Consumer
	orderEventConsumer := app.New(
		orderEventSubscriber,
		eventManager,
		notificationConfigRepo,
		logger,
	)

	// 5. Iniciar consumidor Redis en background
	go func() {
		ctx := context.Background()
		if err := orderEventConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Str("channel", channelOrders).
				Msg("Error al iniciar consumidor de eventos de órdenes")
		}
	}()

	// 6. Init Invoice Event Subscriber y Consumer (facturación en tiempo real)
	invoiceEventSubscriber := redis.NewInvoiceEventSubscriber(redisClient, logger, channelInvoicing)

	invoiceEventConsumer := app.NewInvoiceEventConsumer(
		invoiceEventSubscriber,
		eventManager,
		logger,
	)

	// 7. Iniciar consumidor Redis de facturación en background
	go func() {
		ctx := context.Background()
		if err := invoiceEventConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Str("channel", channelInvoicing).
				Msg("Error al iniciar consumidor de eventos de facturación")
		}
	}()

	// 8. Init Integration Event Subscriber y Consumer (sincronización de integraciones)
	integrationSubscriber := redis.NewIntegrationEventSubscriber(redisClient, logger, channelIntegration)

	integrationConsumer := app.NewIntegrationEventConsumer(
		integrationSubscriber,
		eventManager,
		notificationConfigRepo,
		logger,
	)

	// 9. Iniciar consumidor Redis de integration events en background
	go func() {
		ctx := context.Background()
		if err := integrationConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Str("channel", channelIntegration).
				Msg("Error al iniciar consumidor de integration events")
		}
	}()

	// 10. Init SSE Handler (adaptado a Gin)
	sseHandler := handlers.New(eventManager, logger)

	// 11. Register Routes
	primary.New(sseHandler).RegisterRoutes(router)
}
