package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New crea una nueva instancia del repositorio de configuraciones de notificaciones
func New(database db.IDatabase, logger log.ILogger) ports.IRepository {
	return &repository{
		db:     database,
		logger: logger.WithModule("notification_config_repository"),
	}
}

// NewNotificationTypeRepository crea una nueva instancia del repositorio de tipos de notificaciones
func NewNotificationTypeRepository(database db.IDatabase, logger log.ILogger) ports.INotificationTypeRepository {
	return &notificationTypeRepository{
		db:     database,
		logger: logger.WithModule("notification_type_repository"),
	}
}

// NewNotificationEventTypeRepository crea una nueva instancia del repositorio de tipos de eventos de notificación
func NewNotificationEventTypeRepository(database db.IDatabase, logger log.ILogger) ports.INotificationEventTypeRepository {
	return &notificationEventTypeRepository{
		db:     database,
		logger: logger.WithModule("notification_event_type_repository"),
	}
}

func NewWhatsappTemplateRepository(database db.IDatabase, logger log.ILogger) ports.ITemplateRepository {
	return &whatsappTemplateRepository{
		db:     database,
		logger: logger.WithModule("whatsapp_template_repository"),
	}
}

func NewScheduledRuleRepository(database db.IDatabase, logger log.ILogger) ports.IScheduledRuleRepository {
	return &scheduledRuleRepository{
		db:     database,
		logger: logger.WithModule("scheduled_rule_repository"),
	}
}

func NewScheduledRunRepository(database db.IDatabase, logger log.ILogger) ports.IScheduledRunRepository {
	return &scheduledRunRepository{
		db:     database,
		logger: logger.WithModule("scheduled_run_repository"),
	}
}

func NewScheduledSendRepository(database db.IDatabase, logger log.ILogger) ports.IScheduledSendRepository {
	return &scheduledSendRepository{
		db:     database,
		logger: logger.WithModule("scheduled_send_repository"),
	}
}

func NewSegmentQuerier(database db.IDatabase, logger log.ILogger) ports.ISegmentQuerier {
	return &segmentQuerier{
		db:     database,
		logger: logger.WithModule("segment_querier"),
	}
}

func NewCampaignRepository(database db.IDatabase, logger log.ILogger) ports.ICampaignRepository {
	return &campaignRepository{
		db:     database,
		logger: logger.WithModule("campaign_repository"),
	}
}

func NewCampaignSendRepository(database db.IDatabase, logger log.ILogger) ports.ICampaignSendRepository {
	return &campaignSendRepository{
		db:     database,
		logger: logger.WithModule("campaign_send_repository"),
	}
}

func NewCampaignAudienceQuerier(database db.IDatabase, logger log.ILogger) ports.ICampaignAudienceQuerier {
	return newCampaignAudienceQuerier(database, logger)
}

func NewCampaignSenderQuerier(database db.IDatabase, logger log.ILogger) ports.ICampaignSenderQuerier {
	return newCampaignSenderQuerier(database, logger)
}
