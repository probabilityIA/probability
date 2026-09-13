package whatsappmock

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/handlers"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/store"
	"github.com/secamc93/probability/back/testing/integrations/whatsappmock/internal/webhook"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type Mock struct {
	handler *handlers.Handler
	logger  log.ILogger
	port    string
}

func New(logger log.ILogger, port, centralURL, webhookSecret string) *Mock {
	memory := store.New()
	webhookClient := webhook.New(centralURL, webhookSecret, logger)

	return &Mock{
		handler: handlers.New(memory, webhookClient, logger),
		logger:  logger,
		port:    port,
	}
}

func (m *Mock) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		m.logger.Info().Msgf("[%s] %s %s - %d - %v",
			time.Now().Format("15:04:05"), method, path, c.Writer.Status(), time.Since(start))
	})

	m.handler.RegisterRoutes(router)

	m.logger.Info().Msgf("mock de WhatsApp Graph API escuchando en :%s", m.port)

	return router.Run(":" + m.port)
}
