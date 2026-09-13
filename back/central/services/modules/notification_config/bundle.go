package notification_config

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/campaigns"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/scheduled"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/campaign"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/scheduled_rule"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/whatsapp_template"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/button_reply_consumer"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/queue/campaign_result_consumer"
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
	"github.com/secamc93/probability/back/central/shared/storage"
)

func New(router *gin.RouterGroup, database db.IDatabase, redisClient redisclient.IRedis, logger log.ILogger, rabbitMQ rabbitmq.IQueue, s3 storage.IS3Service) {
	logger = logger.WithModule("notification_config")

	repo := repository.New(database, logger)
	notificationTypeRepo := repository.NewNotificationTypeRepository(database, logger)
	notificationEventTypeRepo := repository.NewNotificationEventTypeRepository(database, logger)
	orderStatusQuerier := repository.NewOrderStatusQuerier(database, logger)
	messageAuditQuerier := repository.NewMessageAuditQuerier(database, logger)
	deliveryLogRepo := repository.NewDeliveryLogRepository(database, logger)

	cacheManager := cache.New(redisClient, repo, orderStatusQuerier, logger)

	aiPauseChecker := cache.NewAIPauseChecker(redisClient)

	ctx := context.Background()
	if err := cacheManager.WarmupCache(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("Error en warmup de cache de notification configs - sistema continuara sin cache")
	}

	useCase := app.New(repo, notificationTypeRepo, notificationEventTypeRepo, cacheManager, messageAuditQuerier, aiPauseChecker, logger)

	configHandler := notification_config.New(useCase, logger)
	typeHandler := notification_type.New(useCase, logger)
	eventTypeHandler := notification_event_type.New(useCase, logger)
	auditHandler := message_audit.New(useCase, logger)

	templateRepo := repository.NewWhatsappTemplateRepository(database, logger)
	templateFlowRepo := repository.NewTemplateFlowRepository(database, logger)
	flowGroupRepo := repository.NewFlowRepository(database, logger)
	segmentQuerier := repository.NewSegmentQuerier(database, logger)

	var templatePublisher templates.ISubmissionPublisher
	var flowPublisher templates.IFlowPublisher
	if rabbitMQ != nil {
		templatePublisher = queue.NewTemplatePublisher(rabbitMQ, logger)
		flowPublisher = queue.NewFlowSendPublisher(rabbitMQ, logger)
	}

	templatesUseCase := templates.New(
		templateRepo,
		templateFlowRepo,
		flowGroupRepo,
		segmentQuerier,
		templatePublisher,
		flowPublisher,
		logger,
	)
	templateHandler := whatsapp_template.New(templatesUseCase, s3, logger)

	scheduledRuleRepo := repository.NewScheduledRuleRepository(database, logger)
	scheduledRunRepo := repository.NewScheduledRunRepository(database, logger)
	scheduledSendRepo := repository.NewScheduledSendRepository(database, logger)

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
		flowGroupRepo,
		templateFlowRepo,
		scheduledPublisher,
		logger,
	)
	scheduledHandler := scheduled_rule.New(scheduledUseCase, logger)

	campaignRepo := repository.NewCampaignRepository(database, logger)
	campaignSendRepo := repository.NewCampaignSendRepository(database, logger)
	campaignAudience := repository.NewCampaignAudienceQuerier(database, logger)
	campaignSender := repository.NewCampaignSenderQuerier(database, logger)

	var campaignPublisher campaigns.ICampaignPublisher
	if rabbitMQ != nil {
		campaignPublisher = queue.NewCampaignSendPublisher(rabbitMQ, logger)
	}
	campaignsUseCase := campaigns.New(
		campaignRepo,
		campaignSendRepo,
		campaignAudience,
		campaignSender,
		templateRepo,
		flowGroupRepo,
		templateFlowRepo,
		campaignPublisher,
		logger,
	)
	campaignHandler := campaign.New(campaignsUseCase, logger)

	configHandler.RegisterRoutes(router)
	typeHandler.RegisterRoutes(router)
	eventTypeHandler.RegisterRoutes(router)
	auditHandler.RegisterRoutes(router)
	templateHandler.RegisterRoutes(router)
	scheduledHandler.RegisterRoutes(router)
	campaignHandler.RegisterRoutes(router)

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

		buttonReplyConsumer := button_reply_consumer.New(rabbitMQ, templatesUseCase, logger)
		go func() {
			if err := buttonReplyConsumer.Start(context.Background()); err != nil {
				logger.Error(ctx).
					Err(err).
					Msg("Error al iniciar consumer de respuestas de botones")
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

		campaignResultConsumer := campaign_result_consumer.New(rabbitMQ, campaignsUseCase, logger)
		go func() {
			if err := campaignResultConsumer.Start(context.Background()); err != nil {
				logger.Error(ctx).
					Err(err).
					Msg("Error al iniciar consumer de resultados de campanas")
			}
		}()

		go worker.NewScheduler(scheduledUseCase, logger).Start(context.Background())
		go worker.NewCampaignDispatcher(campaignsUseCase, logger).Start(context.Background())
	}
}
