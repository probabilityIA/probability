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
<<<<<<< HEAD
=======
	CancelRetry(c *gin.Context)
	GetInvoiceSyncLogs(c *gin.Context)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

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
<<<<<<< HEAD
=======

	// Creación masiva de facturas
	ListInvoiceableOrders(c *gin.Context)
	BulkCreateInvoices(c *gin.Context)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
}

// handler implementa IHandler
type handler struct {
	useCase ports.IUseCase
<<<<<<< HEAD
=======
	repo    ports.IRepository
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	log     log.ILogger
}

// New crea un nuevo handler de facturación
<<<<<<< HEAD
func New(useCase ports.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
=======
func New(useCase ports.IUseCase, repo ports.IRepository, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		repo:    repo,
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		log:     logger.WithModule("invoicing.handler"),
	}
}
