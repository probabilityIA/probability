package usecases

import (
	"time"

	"github.com/secamc93/probability/back/integrationTest/integrations/whatsapp/internal/domain"
	"github.com/secamc93/probability/back/integrationTest/integrations/whatsapp/internal/infra/primary/client"
	"github.com/secamc93/probability/back/integrationTest/shared/env"
	"github.com/secamc93/probability/back/integrationTest/shared/log"
)

// NewWebhookClient crea una nueva instancia del cliente de webhook
func NewWebhookClient(config env.IConfig, logger log.ILogger) domain.IWebhookClient {
	return client.New(config, logger)
}

// NewConversationSimulator crea una nueva instancia del simulador de conversaciones
func NewConversationSimulator(webhookClient domain.IWebhookClient, config env.IConfig, logger log.ILogger) *ConversationSimulator {
	autoReplyDelay := 2 * time.Second // Default 2 segundos
	if delayStr := config.Get("WHATSAPP_AUTO_REPLY_DELAY"); delayStr != "" {
		if duration, err := time.ParseDuration(delayStr + "s"); err == nil {
			autoReplyDelay = duration
		}
	}

	return &ConversationSimulator{
		webhookClient:          webhookClient,
		config:                 config,
		logger:                 logger,
		conversationRepository: domain.NewConversationRepository(),
		autoReplyDelay:         autoReplyDelay,
	}
}
