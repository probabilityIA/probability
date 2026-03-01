package app

import (
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/log"
)

// EventDispatcher enruta eventos a los canales correctos según la configuración
type EventDispatcher struct {
	ssePublisher     ports.ISSEPublisher
	configCache      ports.INotificationConfigCache
	channelPublisher ports.IChannelPublisher
	emailService     email.IEmailService
	logger           log.ILogger
}

// New crea un nuevo EventDispatcher
func New(
	ssePublisher ports.ISSEPublisher,
	configCache ports.INotificationConfigCache,
	channelPublisher ports.IChannelPublisher,
	emailService email.IEmailService,
	logger log.ILogger,
) ports.IEventDispatcher {
	return &EventDispatcher{
		ssePublisher:     ssePublisher,
		configCache:      configCache,
		channelPublisher: channelPublisher,
		emailService:     emailService,
		logger:           logger,
	}
}
