package orders

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecasecreateorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorderscore"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseupdateorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/eventpublisher"
	rabbitqueue "github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el módulo de orders y retorna el caso de uso de create para integraciones
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue) ports.IOrderCreateUseCase {
	// 1. Inicializar Repository
	// Nota: El repositorio de orders incluye métodos de consulta a tablas de estados
	// (order_statuses, payment_statuses, fulfillment_statuses) replicados localmente.
	// NO se comparten repositorios entre módulos - solo consultas SQL directas.
	repo := repository.New(database, environment)

	// 2. Inicializar Publishers
	rabbitPublisher := initRabbitPublisher(rabbitMQ, logger)
	integrationEventPub := eventpublisher.New(rabbitMQ)

	// 3. Inicializar Use Cases
	scoreUseCase := usecaseorderscore.New(repo)
	orderCRUD := usecaseorder.New(repo, rabbitPublisher, logger, scoreUseCase)

	// Update use case se crea primero (no depende de create)
	updateUC := usecaseupdateorder.New(repo, logger, rabbitPublisher, integrationEventPub)
	// Create use case recibe update como dependencia
	createUC := usecasecreateorder.New(repo, logger, rabbitPublisher, integrationEventPub, updateUC)

	requestConfirmationUC := initRequestConfirmationUseCase(repo, rabbitPublisher, logger)

	// 4. Inicializar Handlers y Registrar Rutas
	h := handlers.New(orderCRUD, createUC, requestConfirmationUC, logger)
	h.RegisterRoutes(router)

	// 5. Inicializar Consumers (background goroutines)
	startRabbitMQConsumer(rabbitMQ, logger, createUC, repo)
	startWhatsAppConsumer(rabbitMQ, logger, repo, rabbitPublisher)

	return createUC
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
	// 1. Declarar exchange tipo fanout (envía a TODAS las colas bindeadas)
	if err := rabbitMQ.DeclareExchange(rabbitmq.ExchangeOrderEvents, "fanout", true); err != nil {
		logger.Error(ctx).Err(err).Msg("Error al declarar exchange de órdenes")
		return
	}

	// 2. Declarar y bindear las 5 colas destino
	queues := []string{
		rabbitmq.QueueOrdersToInvoicing,
		rabbitmq.QueueOrdersToWhatsApp,
		rabbitmq.QueueOrdersToScore,
		rabbitmq.QueueOrdersToInventory,
		rabbitmq.QueueOrdersToEvents,
	}

	for _, queueName := range queues {
		if err := rabbitMQ.DeclareQueue(queueName, true); err != nil {
			logger.Error(ctx).Err(err).Str("queue", queueName).Msg("Error al declarar cola")
			continue
		}

		if err := rabbitMQ.BindQueue(queueName, rabbitmq.ExchangeOrderEvents, ""); err != nil {
			logger.Error(ctx).Err(err).Str("queue", queueName).Msg("Error al bindear cola")
			continue
		}

		logger.Info(ctx).
			Str("queue", queueName).
			Str("exchange", rabbitmq.ExchangeOrderEvents).
			Msg("✅ Cola bindeada al fanout de órdenes")
	}

	logger.Info(ctx).
		Int("queues", len(queues)).
		Str("exchange", rabbitmq.ExchangeOrderEvents).
		Msg("Exchange de órdenes configurado")
}

// initRequestConfirmationUseCase inicializa el caso de uso de confirmación por WhatsApp
func initRequestConfirmationUseCase(repo ports.IRepository, rabbitPublisher ports.IOrderRabbitPublisher, logger log.ILogger) ports.IRequestConfirmationUseCase {
	if rabbitPublisher == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ publisher not available, request confirmation use case disabled")
		return nil
	}

	useCase := usecaseorder.NewRequestConfirmationUseCase(repo, rabbitPublisher, logger)

	return useCase
}

// startRabbitMQConsumer inicia el consumer de RabbitMQ para órdenes
func startRabbitMQConsumer(rabbitMQ rabbitmq.IQueue, logger log.ILogger, createUC ports.IOrderCreateUseCase, repo ports.IRepository) {
	if rabbitMQ == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ no disponible, consumer de órdenes deshabilitado")
		return
	}

	consumer := queue.New(rabbitMQ, logger, createUC, repo)

	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			logger.Error().
				Err(err).
				Msg("Consumer de órdenes detenido con error")
		}
	}()
}

// startWhatsAppConsumer inicia el consumer de WhatsApp para confirmaciones
func startWhatsAppConsumer(rabbitMQ rabbitmq.IQueue, logger log.ILogger, repo ports.IRepository, rabbitPublisher ports.IOrderRabbitPublisher) {
	if rabbitMQ == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ not available, WhatsApp consumer disabled")
		return
	}

	consumer := queue.NewWhatsAppConsumer(rabbitMQ, repo, rabbitPublisher, logger)

	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			logger.Error().
				Err(err).
				Msg("WhatsApp consumer stopped with error")
		}
	}()
}
