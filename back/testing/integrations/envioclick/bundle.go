package envioclick

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/domain"
	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// Re-export domain types so external packages (cmd/) can use them
// without importing internal/domain directly (Go internal restriction)
type (
	QuoteRequest     = domain.QuoteRequest
	QuoteResponse    = domain.QuoteResponse
	GenerateResponse = domain.GenerateResponse
	TrackingResponse = domain.TrackingResponse
	CancelResponse   = domain.CancelResponse
	Package          = domain.Package
	Address          = domain.Address
	StoredShipment   = domain.StoredShipment
)

// EnvioClickIntegration representa el modulo de integracion de EnvioClick
type EnvioClickIntegration struct {
	apiSimulator *usecases.APISimulator
	handler      handlers.IHandler
	logger       log.ILogger
	port         string
}

// New inicializa el modulo de EnvioClick para pruebas de integracion
func New(logger log.ILogger, port string) *EnvioClickIntegration {
	apiSimulator := usecases.NewAPISimulator(logger)
	handler := handlers.New(apiSimulator, logger)

	return &EnvioClickIntegration{
		apiSimulator: apiSimulator,
		handler:      handler,
		logger:       logger,
		port:         port,
	}
}

// Start inicia el servidor HTTP del simulador de EnvioClick
func (e *EnvioClickIntegration) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		e.logger.Info().Msgf("[%s] %s %s - Status: %d - Duration: %v",
			time.Now().Format("15:04:05"),
			method,
			path,
			status,
			duration,
		)
	})

	e.handler.RegisterRoutes(router)

	return router.Run(":" + e.port)
}

// SimulateQuote simula cotizacion de envio
func (e *EnvioClickIntegration) SimulateQuote(req QuoteRequest) (*QuoteResponse, error) {
	return e.apiSimulator.HandleQuote(req)
}

// SimulateGenerate simula generacion de guia
func (e *EnvioClickIntegration) SimulateGenerate(req QuoteRequest) (*GenerateResponse, error) {
	return e.apiSimulator.HandleGenerate(req)
}

// SimulateTrack simula rastreo de envio
func (e *EnvioClickIntegration) SimulateTrack(trackingNumber string) (*TrackingResponse, error) {
	return e.apiSimulator.HandleTrack(trackingNumber)
}

// SimulateCancel simula cancelacion de envio
func (e *EnvioClickIntegration) SimulateCancel(shipmentID string) (*CancelResponse, error) {
	return e.apiSimulator.HandleCancel(shipmentID)
}

// GetAllShipments retorna todos los envios almacenados
func (e *EnvioClickIntegration) GetAllShipments() []*StoredShipment {
	return e.apiSimulator.Repository.GetAll()
}
