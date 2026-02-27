package orders

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseordermapping"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/queue"
	rabbitqueue "github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue"
	redisevents "github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

const defaultRedisChannel = redisclient.ChannelOrdersEvents

// New inicializa el módulo de orders y retorna el caso de uso de mapping para integraciones
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redisclient.IRedis) ports.IOrderMappingUseCase {
	// 1. Inicializar Repository
	// Nota: El repositorio de orders incluye métodos de consulta a tablas de estados
	// (order_statuses, payment_statuses, fulfillment_statuses) replicados localmente.
	// NO se comparten repositorios entre módulos - solo consultas SQL directas.
	repo := repository.New(database, environment)

	// 2. Inicializar Publishers
	eventPublisher := initRedisPublisher(redisClient, logger)
	rabbitPublisher := initRabbitPublisher(rabbitMQ, logger)

	// 3. Inicializar Use Cases
	scoreUseCase := usecaseorderscore.New(repo)
	orderCRUD := usecaseorder.New(repo, eventPublisher, rabbitPublisher, logger, scoreUseCase)
	orderMapping := usecaseordermapping.New(repo, logger, eventPublisher, rabbitPublisher)
	requestConfirmationUC := initRequestConfirmationUseCase(repo, rabbitPublisher, logger)

	// 4. Inicializar Handlers y Registrar Rutas
	h := handlers.New(orderCRUD, orderMapping, requestConfirmationUC)
	h.RegisterRoutes(router)

	// 5. Inicializar Consumers (background goroutines)
	startRedisEventConsumer(redisClient, logger, scoreUseCase)
	startRabbitMQConsumer(rabbitMQ, logger, orderMapping, repo)
	startWhatsAppConsumer(rabbitMQ, logger, orderCRUD, repo, eventPublisher)

	return orderMapping
}

// initRedisPublisher inicializa el publicador de eventos de Redis
func initRedisPublisher(redisClient redisclient.IRedis, logger log.ILogger) ports.IOrderEventPublisher {
	if redisClient == nil {
		logger.Warn(context.Background()).Msg("Redis client not available, event publisher disabled")
		return nil
	}

	publisher := redisevents.NewOrderEventPublisher(redisClient, logger, defaultRedisChannel)

	// Registrar canal para mostrar en startup logs
	redisClient.RegisterChannel(defaultRedisChannel)

	return publisher
}

// initRabbitPublisher inicializa el publicador de RabbitMQ
func initRabbitPublisher(rabbitMQ rabbitmq.IQueue, logger log.ILogger) ports.IOrderRabbitPublisher {
	if rabbitMQ == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ not available, rabbit publisher disabled")
		return nil
	}

	// Configurar exchange y bindings para distribuir eventos a múltiples consumers
	setupOrdersExchange(rabbitMQ, logger)

	publisher := rabbitqueue.NewOrderRabbitPublisher(rabbitMQ, logger)

	return publisher
}

// setupOrdersExchange configura el exchange de órdenes y sus bindings
func setupOrdersExchange(rabbitMQ rabbitmq.IQueue, logger log.ILogger) {
	ctx := context.Background()
	exchangeName := "orders.events"

	// 1. Declarar exchange tipo fanout (envía a TODAS las colas bindeadas)
	if err := rabbitMQ.DeclareExchange(exchangeName, "fanout", true); err != nil {
		logger.Error(ctx).Err(err).Msg("Error al declarar exchange de órdenes")
		return
	}

	// 2. Declarar y bindear las 3 colas destino
	queues := []string{
		"orders.events.invoicing",
		"orders.events.whatsapp",
		"orders.events.score",
		"orders.events.inventory",
	}

	for _, queueName := range queues {
		if err := rabbitMQ.DeclareQueue(queueName, true); err != nil {
			logger.Error(ctx).Err(err).Str("queue", queueName).Msg("Error al declarar cola")
			continue
		}

		if err := rabbitMQ.BindQueue(queueName, exchangeName, ""); err != nil {
			logger.Error(ctx).Err(err).Str("queue", queueName).Msg("Error al bindear cola")
			continue
		}
	}

	// Exchange configurado - las colas se muestran en LogStartupInfo()
}

// initRequestConfirmationUseCase inicializa el caso de uso de confirmación por WhatsApp
func initRequestConfirmationUseCase(repo ports.IRepository, rabbitPublisher ports.IOrderRabbitPublisher, logger log.ILogger) usecaseorder.IRequestConfirmationUseCase {
	if rabbitPublisher == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ publisher not available, request confirmation use case disabled")
		return nil
	}

	useCase := usecaseorder.NewRequestConfirmationUseCase(repo, rabbitPublisher, logger)

	return useCase
}

// startRedisEventConsumer inicia el consumer de eventos de Redis para cálculo de score
func startRedisEventConsumer(redisClient redisclient.IRedis, logger log.ILogger, scoreUseCase ports.IOrderScoreUseCase) {
	if redisClient == nil {
		logger.Warn(context.Background()).Msg("Redis client not available, event consumer disabled")
		return
	}

	consumer := redisevents.NewOrderEventConsumer(redisClient, logger, defaultRedisChannel, scoreUseCase)

	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			logger.Error().
				Err(err).
				Msg("Order event consumer stopped with error")
		}
	}()
}

// startRabbitMQConsumer inicia el consumer de RabbitMQ para órdenes
func startRabbitMQConsumer(rabbitMQ rabbitmq.IQueue, logger log.ILogger, orderMapping ports.IOrderMappingUseCase, repo ports.IRepository) {
	if rabbitMQ == nil {
		logger.Warn(context.Background()).Msg("❌ RabbitMQ no disponible, consumer de órdenes deshabilitado")
		return
	}

	consumer := queue.New(rabbitMQ, logger, orderMapping, repo)

	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			logger.Error().
				Err(err).
				Msg("❌ Consumer de órdenes detenido con error")
		}
	}()
}

// startWhatsAppConsumer inicia el consumer de WhatsApp para confirmaciones
func startWhatsAppConsumer(rabbitMQ rabbitmq.IQueue, logger log.ILogger, orderCRUD ports.IOrderUseCase, repo ports.IRepository, eventPublisher ports.IOrderEventPublisher) {
	if rabbitMQ == nil || eventPublisher == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ or event publisher not available, WhatsApp consumer disabled")
		return
	}

	consumer := queue.NewWhatsAppConsumer(rabbitMQ, orderCRUD, repo, eventPublisher, logger)

	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			logger.Error().
				Err(err).
				Msg("WhatsApp consumer stopped with error")
		}
	}()
}

