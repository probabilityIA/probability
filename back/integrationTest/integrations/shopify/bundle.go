package shopify

import (
	"github.com/secamc93/probability/back/integrationTest/integrations/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/integrationTest/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/integrationTest/shared/env"
	"github.com/secamc93/probability/back/integrationTest/shared/log"
)

// New inicializa el módulo de Shopify para pruebas de integración
func New(config env.IConfig, logger log.ILogger) *ShopifyIntegration {
	webhookClient := usecases.NewWebhookClient(config, logger)
	orderSimulator := usecases.NewOrderSimulator(webhookClient, config, logger)

	return &ShopifyIntegration{
		orderSimulator: orderSimulator,
		logger:         logger,
	}
}

// ShopifyIntegration representa el módulo de integración de Shopify
type ShopifyIntegration struct {
	orderSimulator *usecases.OrderSimulator
	logger         log.ILogger
}

// SimulateOrder simula una orden de Shopify y la envía como webhook
func (s *ShopifyIntegration) SimulateOrder(topic string) error {
	return s.orderSimulator.SimulateOrder(topic)
}

// GetAllOrders retorna todas las órdenes almacenadas
func (s *ShopifyIntegration) GetAllOrders() []*domain.Order {
	return s.orderSimulator.GetAllOrders()
}

// GetOrderByNumber obtiene una orden por su número
func (s *ShopifyIntegration) GetOrderByNumber(orderNumber string) (*domain.Order, bool) {
	return s.orderSimulator.GetOrderByNumber(orderNumber)
}





