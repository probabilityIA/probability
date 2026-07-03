package woocommerce

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/woocommerce/internal/handlers"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type WooCommerceIntegration struct {
	handler *handlers.Handler
	logger  log.ILogger
	port    string
}

func New(logger log.ILogger, port string) *WooCommerceIntegration {
	return &WooCommerceIntegration{
		handler: handlers.New(logger),
		logger:  logger,
		port:    port,
	}
}

func (s *WooCommerceIntegration) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		s.logger.Info().Msgf("[%s] %s %s - Status: %d - Duration: %v",
			time.Now().Format("15:04:05"), method, path, c.Writer.Status(), time.Since(start))
	})

	s.handler.RegisterRoutes(router)
	return router.Run(":" + s.port)
}
