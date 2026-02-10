package softpymes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/domain"
	"github.com/secamc93/probability/back/testing/integrations/softpymes/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// New inicializa el módulo de SoftPymes para pruebas de integración
func New(logger log.ILogger, port string) *SoftPymesIntegration {
	// 1. Capa de aplicación (use cases)
	apiSimulator := usecases.NewAPISimulator(logger)

	// 2. Capa de infraestructura primaria (handlers HTTP)
	handler := handlers.New(apiSimulator, logger)

	return &SoftPymesIntegration{
		apiSimulator: apiSimulator,
		handler:      handler,
		logger:       logger,
		port:         port,
	}
}

// SoftPymesIntegration representa el módulo de integración de SoftPymes
type SoftPymesIntegration struct {
	apiSimulator *usecases.APISimulator
	handler      handlers.IHandler
	logger       log.ILogger
	port         string
}

// SimulateAuth simula autenticación
func (s *SoftPymesIntegration) SimulateAuth(apiKey, apiSecret, referer string) (string, error) {
	return s.apiSimulator.HandleAuth(apiKey, apiSecret, referer)
}

// SimulateInvoice simula creación de factura
func (s *SoftPymesIntegration) SimulateInvoice(token string, invoiceData map[string]interface{}) (*usecases.InvoiceWithDetails, error) {
	return s.apiSimulator.HandleCreateInvoice(token, invoiceData)
}

// SimulateCreditNote simula creación de nota de crédito
func (s *SoftPymesIntegration) SimulateCreditNote(token string, creditNoteData map[string]interface{}) (*domain.CreditNote, error) {
	return s.apiSimulator.HandleCreateCreditNote(token, creditNoteData)
}

// ListInvoices retorna todas las facturas almacenadas
func (s *SoftPymesIntegration) ListInvoices(token string) ([]domain.Invoice, error) {
	return s.apiSimulator.HandleListDocuments(token, nil)
}

// GetRepository retorna el repositorio (para listar sin token)
func (s *SoftPymesIntegration) GetRepository() *domain.InvoiceRepository {
	return s.apiSimulator.Repository
}

// Start inicia el servidor HTTP del simulador de Softpymes
func (s *SoftPymesIntegration) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware de logging
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

	// Registrar rutas usando el handler de la capa de infraestructura
	s.handler.RegisterRoutes(router)

	return router.Run(":" + s.port)
}
