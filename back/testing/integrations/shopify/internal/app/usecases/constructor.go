package usecases

import (
	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/infra/primary/client"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// NewWebhookClient crea una nueva instancia del cliente de webhook
func NewWebhookClient(config env.IConfig, logger log.ILogger) domain.IWebhookClient {
	return client.New(config, logger)
}

// NewOrderSimulator crea una nueva instancia del simulador de órdenes
func NewOrderSimulator(webhookClient domain.IWebhookClient, config env.IConfig, logger log.ILogger) *OrderSimulator {
	return &OrderSimulator{
		webhookClient:   webhookClient,
		config:          config,
		logger:          logger,
		orderRepository: domain.NewOrderRepository(),
		dataGenerator:   NewRandomDataGenerator(),
		orderNumberSeq:  1000,
		businessConfig:  domain.DefaultTestBusinessConfig(), // Usar config del business de prueba
	}
}













