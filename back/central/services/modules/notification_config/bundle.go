package notification_config

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/scheduled"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/scheduled_rule"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/whatsapp_template"
	deliveryConsumer "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/consumer"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/scheduled_result_consumer"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/template_result_consumer"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/whatsapp_persistence_consumer"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/worker"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/cache"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa y registra el módulo de configuración de notificaciones
func New(router *gin.RouterGroup, database db.IDatabase, redisClient redisclient.IRedis, logger log.ILogger, rabbitMQ rabbitmq.IQueue) {
	logger = logger.WithModule("notification_config")

	// 1. Infraestructura secundaria (adaptadores de salida)
	repo := repository.New(database, logger)
	notificationTypeRepo := repository.NewNotificationTypeRepository(database, logger)
	notificationEventTypeRepo := repository.NewNotificationEventTypeRepository(database, logger)
	orderStatusQuerier := repository.NewOrderStatusQuerier(database, logger)
	messageAuditQuerier := repository.NewMessageAuditQuerier(database, logger)
	deliveryLogRepo := repository.NewDeliveryLogRepository(database, logger)

	// Cache Manager
	cacheManager := cache.New(redisClient, repo, orderStatusQuerier, logger)

	// AI Pause Checker (lee estado de IA pausada desde Redis)
	aiPauseChecker := cache.NewAIPauseChecker(redisClient)

	// Warmup inicial del cache
	ctx := context.Background()
	if err := cacheManager.WarmupCache(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("❌ Error en warmup de cache de notification configs - sistema continuará sin cache")
	}

	// 2. Capa de aplicación (casos de uso) - inyectar cache manager
	useCase := app.New(repo, notificationTypeRepo, notificationEventTypeRepo, cacheManager, messageAuditQuerier, aiPauseChecker, logger)

	// 3. Infraestructura primaria (adaptadores de entrada)
	configHandler := notification_config.New(useCase, logger)
	typeHandler := notification_type.New(useCase, logger)
	eventTypeHandler := notification_event_type.New(useCase, logger)
	auditHandler := message_audit.New(useCase, logger)

	templateRepo := repository.NewWhatsappTemplateRepository(database, logger)
	var templatePublisher templates.ISubmissionPublisher
	if rabbitMQ != nil {
		templatePublisher = queue.NewTemplatePublisher(rabbitMQ, logger)
	}
	templatesUseCase := templates.New(templateRepo, templatePublisher, logger)
	templateHandler := whatsapp_template.New(templatesUseCase, logger)

	scheduledRuleRepo := repository.NewScheduledRuleRepository(database, logger)
	scheduledRunRepo := repository.NewScheduledRunRepository(database, logger)
	scheduledSendRepo := repository.NewScheduledSendRepository(database, logger)
	segmentQuerier := repository.NewSegmentQuerier(database, logger)

	var scheduledPublisher scheduled.ISendPublisher
	if rabbitMQ != nil {
		scheduledPublisher = queue.NewScheduledSendPublisher(rabbitMQ, logger)
	}
	scheduledUseCase := scheduled.New(
		scheduledRuleRepo,
		scheduledRunRepo,
		scheduledSendRepo,
		segmentQuerier,
		templateRepo,
		scheduledPublisher,
		logger,
	)
	scheduledHandler := scheduled_rule.New(scheduledUseCase, logger)

	// 4. Registrar rutas HTTP
	configHandler.RegisterRoutes(router)
	typeHandler.RegisterRoutes(router)
	eventTypeHandler.RegisterRoutes(router)
	auditHandler.RegisterRoutes(router)
	templateHandler.RegisterRoutes(router)
	scheduledHandler.RegisterRoutes(router)

	// 5. Consumer de resultados de entrega (email, SMS futuro, etc.)
	if rabbitMQ != nil {
		consumer := deliveryConsumer.New(rabbitMQ, deliveryLogRepo, logger)
		go func() {
			if err := consumer.Start(context.Background()); err != nil {
				logger.Error(ctx).
					Err(err).
					Msg("Error al iniciar consumer de delivery results")
			}
		}()

		whatsappPersister := repository.NewWhatsAppPersister(database, logger)

		persistenceConsumer := whatsapp_persistence_consumer.New(rabbitMQ, whatsappPersister, logger)
		go func() {
			if err := persistenceConsumer.Start(context.Background()); err != nil {
				logger.Error(ctx).
					Err(err).
					Msg("Error al iniciar consumer de WhatsApp persistence")
			}
		}()

		templateResultConsumer := template_result_consumer.New(rabbitMQ, templatesUseCase, logger)
		go func() {
			if err := templateResultConsumer.Start(context.Background()); err != nil {
				logger.Error(ctx).
					Err(err).
					Msg("Error al iniciar consumer de resultados de plantillas")
			}
		}()

		scheduledResultConsumer := scheduled_result_consumer.New(rabbitMQ, scheduledUseCase, logger)
		go func() {
			if err := scheduledResultConsumer.Start(context.Background()); err != nil {
				logger.Error(ctx).
					Err(err).
					Msg("Error al iniciar consumer de resultados de envios programados")
			}
		}()

		go worker.NewScheduler(scheduledUseCase, logger).Start(context.Background())
	}
}
