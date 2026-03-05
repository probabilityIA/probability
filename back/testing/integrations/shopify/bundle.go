package shopify

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
	sharedtypes "github.com/secamc93/probability/back/testing/shared/types"
)

// New inicializa el módulo de Shopify para pruebas de integración.
// Genera un mix de órdenes: ~60% dual-currency USD/COP, ~40% single-currency COP.
// port es el puerto donde se levantará el mock Shopify API (ej: "9092").
func New(config env.IConfig, logger log.ILogger, port string) *ShopifyIntegration {
	businessConfig := domain.DefaultTestBusinessConfig()
	webhookClient := usecases.NewWebhookClient(config, logger)
	orderSimulator := usecases.NewOrderSimulator(webhookClient, config, logger, businessConfig)
	mockAPI := usecases.NewMockAPIServer(logger, businessConfig)
	handler := handlers.New(mockAPI, logger)

	return &ShopifyIntegration{
		orderSimulator: orderSimulator,
		mockAPI:        mockAPI,
		handler:        handler,
		logger:         logger,
		port:           port,
		businessConfig: businessConfig,
	}
}

// ShopifyIntegration representa el módulo de integración de Shopify
type ShopifyIntegration struct {
	orderSimulator *usecases.OrderSimulator
	mockAPI        *usecases.MockAPIServer
	handler        handlers.IHandler
	logger         log.ILogger
	port           string
	businessConfig *domain.BusinessConfig
}

// Start inicia el servidor HTTP que simula el API REST de Shopify.
// Pre-genera initialOrders órdenes distribuidas en los últimos 6 meses.
func (s *ShopifyIntegration) Start(initialOrders int) error {
	// Pre-generar órdenes si se solicita
	if initialOrders > 0 {
		dateTo := time.Now()
		dateFrom := dateTo.AddDate(0, -6, 0) // Últimos 6 meses
		s.mockAPI.GenerateOrders(initialOrders, dateFrom, dateTo)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware de logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		query := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logLine := s.logger.Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("duration", duration)
		if query != "" {
			logLine = logLine.Str("query", query)
		}
		logLine.Msg("📡 Shopify Mock API")
	})

	s.handler.RegisterRoutes(router)

	s.logger.Info().
		Str("port", s.port).
		Int("initial_orders", initialOrders).
		Str("mode", "mixed USD/COP + COP (~60/40)").
		Msg("🚀 Shopify Mock API Server iniciado")

	return router.Run(":" + s.port)
}

// GenerateOrders genera órdenes adicionales en el mock API.
func (s *ShopifyIntegration) GenerateOrders(count int, dateFrom, dateTo time.Time) {
	s.mockAPI.GenerateOrders(count, dateFrom, dateTo)
}

// GetMockTotalOrders retorna el total de órdenes en el mock API.
func (s *ShopifyIntegration) GetMockTotalOrders() int {
	return s.mockAPI.GetTotalOrders()
}

// SimulateOrder simula una orden de Shopify y la envía como webhook
func (s *ShopifyIntegration) SimulateOrder(topic string) error {
	return s.orderSimulator.SimulateOrder(topic)
}

// BuildWebhookPayload builds the webhook payload without sending it
func (s *ShopifyIntegration) BuildWebhookPayload(topic string, baseURL string) (*sharedtypes.WebhookPayload, error) {
	return s.orderSimulator.BuildWebhookPayload(topic, baseURL)
}

// GetWebhookTopics returns the list of supported webhook topics
func (s *ShopifyIntegration) GetWebhookTopics() []string {
	return s.orderSimulator.GetWebhookTopics()
}

// GetAllOrders retorna todas las órdenes almacenadas (del webhook simulator)
func (s *ShopifyIntegration) GetAllOrders() []*domain.Order {
	return s.orderSimulator.GetAllOrders()
}

// GetOrderByNumber obtiene una orden por su número (del webhook simulator)
func (s *ShopifyIntegration) GetOrderByNumber(orderNumber string) (*domain.Order, bool) {
	return s.orderSimulator.GetOrderByNumber(orderNumber)
}

