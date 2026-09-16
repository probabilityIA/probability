package app

import (
	"github.com/secamc93/probability/back/central/services/events/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type EventDispatcher struct {
	ssePublisher     ports.ISSEPublisher
	configCache      ports.INotificationConfigCache
	channelPublisher ports.IChannelPublisher
	alerts           ports.IAssistantAlertPublisher
	logger           log.ILogger
}

func New(
	ssePublisher ports.ISSEPublisher,
	configCache ports.INotificationConfigCache,
	channelPublisher ports.IChannelPublisher,
	alerts ports.IAssistantAlertPublisher,
	logger log.ILogger,
) ports.IEventDispatcher {
	return &EventDispatcher{
		ssePublisher:     ssePublisher,
		configCache:      configCache,
		channelPublisher: channelPublisher,
		alerts:           alerts,
		logger:           logger,
	}
}
