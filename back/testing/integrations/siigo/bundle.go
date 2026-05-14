package siigo

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/domain"
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type SiigoIntegration struct {
	apiSimulator *usecases.APISimulator
	handler      handlers.IHandler
	logger       log.ILogger
	port         string
}

func New(logger log.ILogger, port string) *SiigoIntegration {
	apiSimulator := usecases.NewAPISimulator(logger)
	handler := handlers.New(apiSimulator, logger)

	return &SiigoIntegration{
		apiSimulator: apiSimulator,
		handler:      handler,
		logger:       logger,
		port:         port,
	}
}

func (s *SiigoIntegration) GetRepository() *domain.Repository {
	return s.apiSimulator.Repository
}

func (s *SiigoIntegration) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		s.logger.Info().Msgf("[%s] %s %s - Status: %d - Duration: %v",
			time.Now().Format("15:04:05"),
			method,
			path,
			status,
			duration,
		)
	})

	s.handler.RegisterRoutes(router)

	return router.Run(":" + s.port)
}
