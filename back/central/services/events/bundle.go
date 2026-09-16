package events

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/events/internal/app"
	"github.com/secamc93/probability/back/central/services/events/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/events/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/events/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/events/internal/infra/secondary/channel"
	rmqInfra "github.com/secamc93/probability/back/central/services/events/internal/infra/secondary/rabbitmq"
	"github.com/secamc93/probability/back/central/services/events/internal/infra/secondary/sse"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

func New(
	router *gin.RouterGroup,
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	redisClient redisclient.IRedis,
) {
	if rabbitMQ == nil {
		logger.Warn(context.Background()).
			Msg("RabbitMQ no disponible - eventos por cola deshabilitados, SSE activo")
		degradedManager := sse.New(logger)
		handlers.New(degradedManager, logger).RegisterRoutes(router)
		return
	}

	if err := rmqInfra.SetupInfrastructure(rabbitMQ, logger); err != nil {
		logger.Error(context.Background()).
			Err(err).
			Msg("Error critico al configurar infraestructura RabbitMQ de eventos")
		return
	}

	eventManager := sse.New(logger)

	configCache := cache.New(redisClient, logger)

	channelPub := channel.New(rabbitMQ, logger)

	alertPub := channel.NewAssistantAlertPublisher(rabbitMQ, logger)

	dispatcher := app.New(eventManager, configCache, channelPub, alertPub, logger)

	eventConsumer := consumer.New(rabbitMQ, dispatcher, logger)
	go func() {
		ctx := context.Background()
		if err := eventConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Msg("Error al iniciar consumer de eventos unificado")
		}
	}()

	orderEventConsumer := consumer.NewOrderEventConsumer(rabbitMQ, dispatcher, logger)
	go func() {
		ctx := context.Background()
		if err := orderEventConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Msg("Error al iniciar consumer de eventos de órdenes")
		}
	}()

	sseHandler := handlers.New(eventManager, logger)
	sseHandler.RegisterRoutes(router)

	logger.Info(context.Background()).
		Msg("Módulo de eventos unificado inicializado")
}
