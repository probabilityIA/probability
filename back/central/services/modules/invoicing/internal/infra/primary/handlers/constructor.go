package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define las operaciones HTTP del módulo de facturación
type IHandler interface {
	// Rutas
	RegisterRoutes(router *gin.RouterGroup)

	// Facturas
	CreateInvoice(c *gin.Context)
	ListInvoices(c *gin.Context)
	GetInvoice(c *gin.Context)
	CancelInvoice(c *gin.Context)
	RetryInvoice(c *gin.Context)

	// Notas de crédito
	CreateCreditNote(c *gin.Context)

	// Proveedores (DEPRECATED - Migrados a integrations/core)
	CreateProvider(c *gin.Context)
	ListProviders(c *gin.Context)
	GetProvider(c *gin.Context)
	UpdateProvider(c *gin.Context)
	TestProvider(c *gin.Context)

	// Configuraciones
	CreateConfig(c *gin.Context)
	ListConfigs(c *gin.Context)
	GetConfig(c *gin.Context)
	UpdateConfig(c *gin.Context)
	DeleteConfig(c *gin.Context)

	// Estadísticas y resúmenes
	GetSummary(c *gin.Context)
	GetStats(c *gin.Context)
	GetTrends(c *gin.Context)
}

// handler implementa IHandler
type handler struct {
	useCase ports.IUseCase
	log     log.ILogger
}

// New crea un nuevo handler de facturación
func New(useCase ports.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		log:     logger.WithModule("invoicing.handler"),
	}
}
