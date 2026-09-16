package app

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type useCase struct {
	chatMediaSigner       ports.IChatMediaSigner
	repository            ports.IRepository
	notificationTypeRepo  ports.INotificationTypeRepository
	notificationEventRepo ports.INotificationEventTypeRepository
	cacheManager          ports.ICacheManager
	messageAuditQuerier   ports.IMessageAuditQuerier
	aiPauseChecker        ports.IAIPauseChecker
	defaultRules          ports.IDefaultRulesQuerier
	logger                log.ILogger
}

func New(
	repository ports.IRepository,
	notificationTypeRepo ports.INotificationTypeRepository,
	notificationEventRepo ports.INotificationEventTypeRepository,
	cacheManager ports.ICacheManager,
	messageAuditQuerier ports.IMessageAuditQuerier,
	aiPauseChecker ports.IAIPauseChecker,
	logger log.ILogger,
) ports.IUseCase {
	return &useCase{
		repository:            repository,
		notificationTypeRepo:  notificationTypeRepo,
		notificationEventRepo: notificationEventRepo,
		cacheManager:          cacheManager,
		messageAuditQuerier:   messageAuditQuerier,
		aiPauseChecker:        aiPauseChecker,
		logger:                logger.WithModule("notification_config_usecase"),
	}
}
