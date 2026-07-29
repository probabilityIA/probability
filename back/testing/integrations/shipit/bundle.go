package shipit

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/domain"
	"github.com/secamc93/probability/back/testing/integrations/shipit/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type StoredShipment = domain.StoredShipment

type ShipitIntegration struct {
	apiSimulator *usecases.APISimulator
	handler      handlers.IHandler
	logger       log.ILogger
	port         string
}

func New(logger log.ILogger, port, webhookTarget string) *ShipitIntegration {
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	apiSimulator := usecases.NewAPISimulator(logger, webhookTarget, baseURL)
	handler := handlers.New(apiSimulator, logger)

	return &ShipitIntegration{
		apiSimulator: apiSimulator,
		handler:      handler,
		logger:       logger,
		port:         port,
	}
}

func (s *ShipitIntegration) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		s.logger.Info().Msgf("[%s] %s %s - Status: %d - Duration: %v",
			time.Now().Format("15:04:05"),
			method,
			path,
			c.Writer.Status(),
			time.Since(start),
		)
	})

	s.handler.RegisterRoutes(router)

	return router.Run(":" + s.port)
}

func (s *ShipitIntegration) SimulateStatusChange(trackingNumber, status string) (*StoredShipment, error) {
	return s.apiSimulator.SimulateStatusChange(trackingNumber, status)
}

func (s *ShipitIntegration) GetAllShipments() []*StoredShipment {
	return s.apiSimulator.Repository.GetAll()
}
