package orders

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecasecreateorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseupdateorder"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/app/usecaseupdatestatus"
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

type Bundle struct {
	CreateUC                ports.IOrderCreateUseCase
	SendGuideNotificationUC ports.ISendGuideNotificationUseCase
	RequestConfirmationUC   ports.IRequestConfirmationUseCase
}

func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue) *Bundle {
	repo := repository.New(database, environment)

	rabbitPublisher := initRabbitPublisher(rabbitMQ, logger)
	integrationEventPub := eventpublisher.New(rabbitMQ)

	orderCRUD := usecaseorder.New(repo, rabbitPublisher, logger)

	updateUC := usecaseupdateorder.New(repo, logger, rabbitPublisher, integrationEventPub)
	createUC := usecasecreateorder.New(repo, logger, rabbitPublisher, integrationEventPub, updateUC)

	statusUC := usecaseupdatestatus.New(repo, logger, rabbitPublisher)
	requestConfirmationUC := initRequestConfirmationUseCase(repo, rabbitPublisher, logger)
	sendGuideNotificationUC := initSendGuideNotificationUseCase(repo, rabbitPublisher, logger)

	h := handlers.New(orderCRUD, createUC, requestConfirmationUC, sendGuideNotificationUC, statusUC, logger)
	h.RegisterRoutes(router)

	startRabbitMQConsumer(rabbitMQ, logger, createUC, repo, integrationEventPub)
	startWhatsAppConsumer(rabbitMQ, logger, repo, rabbitPublisher)
	startInventoryFeedbackConsumer(rabbitMQ, logger, repo, rabbitPublisher)

	return &Bundle{
		CreateUC:                createUC,
		SendGuideNotificationUC: sendGuideNotificationUC,
		RequestConfirmationUC:   requestConfirmationUC,
	}
}

func initRabbitPublisher(rabbitMQ rabbitmq.IQueue, logger log.ILogger) ports.IOrderRabbitPublisher {
	if rabbitMQ == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ not available, rabbit publisher disabled")
		return nil
	}

	setupOrdersExchange(rabbitMQ, logger)

	publisher := rabbitqueue.NewOrderRabbitPublisher(rabbitMQ, logger)

	return publisher
}

func setupOrdersExchange(rabbitMQ rabbitmq.IQueue, logger log.ILogger) {
	ctx := context.Background()
	if err := rabbitMQ.DeclareExchange(rabbitmq.ExchangeOrderEvents, "fanout", true); err != nil {
		logger.Error(ctx).Err(err).Msg("Error al declarar exchange de ordenes")
		return
	}

	queues := []string{
		rabbitmq.QueueOrdersToInvoicing,
		rabbitmq.QueueOrdersToScore,
		rabbitmq.QueueOrdersToInventory,
		rabbitmq.QueueOrdersToEvents,
		rabbitmq.QueueOrdersToCustomers,
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
			Msg("Cola bindeada al fanout de ordenes")
	}
}

func initRequestConfirmationUseCase(repo ports.IRepository, rabbitPublisher ports.IOrderRabbitPublisher, logger log.ILogger) ports.IRequestConfirmationUseCase {
	if rabbitPublisher == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ publisher not available, request confirmation use case disabled")
		return nil
	}

	return usecaseorder.NewRequestConfirmationUseCase(repo, rabbitPublisher, logger)
}

func initSendGuideNotificationUseCase(repo ports.IRepository, rabbitPublisher ports.IOrderRabbitPublisher, logger log.ILogger) ports.ISendGuideNotificationUseCase {
	if rabbitPublisher == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ publisher not available, send guide notification use case disabled")
		return nil
	}
	return usecaseorder.NewSendGuideNotificationUseCase(repo, rabbitPublisher, logger)
}

func startRabbitMQConsumer(rabbitMQ rabbitmq.IQueue, logger log.ILogger, createUC ports.IOrderCreateUseCase, repo ports.IRepository, eventPub ports.IIntegrationEventPublisher) {
	if rabbitMQ == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ no disponible, consumer de ordenes deshabilitado")
		return
	}

	consumer := queue.New(rabbitMQ, logger, createUC, repo, eventPub)

	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			logger.Error().
				Err(err).
				Msg("Consumer de ordenes detenido con error")
		}
	}()
}

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

func startInventoryFeedbackConsumer(rabbitMQ rabbitmq.IQueue, logger log.ILogger, repo ports.IRepository, rabbitPublisher ports.IOrderRabbitPublisher) {
	if rabbitMQ == nil {
		return
	}

	consumer := queue.NewInventoryConsumer(rabbitMQ, repo, rabbitPublisher, logger)
	consumer.Start(context.Background())
}
