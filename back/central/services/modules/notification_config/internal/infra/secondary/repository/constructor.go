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
