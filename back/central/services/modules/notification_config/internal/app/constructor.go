package app

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	repository              ports.IRepository
	notificationTypeRepo    ports.INotificationTypeRepository
	notificationEventRepo   ports.INotificationEventTypeRepository
	logger                  log.ILogger
}

// New crea una nueva instancia del caso de uso
func New(
	repository ports.IRepository,
	notificationTypeRepo ports.INotificationTypeRepository,
	notificationEventRepo ports.INotificationEventTypeRepository,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		repository:            repository,
		notificationTypeRepo:  notificationTypeRepo,
		notificationEventRepo: notificationEventRepo,
		logger:                logger.WithModule("notification_config_usecase"),
	}
}
