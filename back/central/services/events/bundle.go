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
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa el módulo unificado de eventos.
// NO recibe database — este módulo lee configs de Redis cache.
func New(
	router *gin.RouterGroup,
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	redisClient redisclient.IRedis,
	emailService email.IEmailService,
) {
	// 1. Declarar infraestructura RabbitMQ (exchange + queue + binding)
	rmqInfra.SetupInfrastructure(rabbitMQ, logger)

	// 2. SSE EventManager
	eventManager := sse.New(logger)

	// 2. Notification config cache reader (lee de Redis, no de BD)
	configCache := cache.New(redisClient, logger)

	// 3. Channel publishers (WhatsApp → RabbitMQ queue)
	channelPub := channel.New(rabbitMQ, logger)

	// 4. Event Dispatcher (capa de aplicación)
	dispatcher := app.New(eventManager, configCache, channelPub, emailService, logger)

	// 5. RabbitMQ consumer → background
	eventConsumer := consumer.New(rabbitMQ, dispatcher, logger)
	go func() {
		ctx := context.Background()
		if err := eventConsumer.Start(ctx); err != nil {
			logger.Error(ctx).
				Err(err).
				Msg("Error al iniciar consumer de eventos unificado")
		}
	}()

	// 6. SSE HTTP handler + routes
	sseHandler := handlers.New(eventManager, logger)
	sseHandler.RegisterRoutes(router)

	logger.Info(context.Background()).
		Msg("Módulo de eventos unificado inicializado")
}
