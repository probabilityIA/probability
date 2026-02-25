package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
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
	CancelRetry(c *gin.Context)
	EnableRetry(c *gin.Context)
	GetInvoiceSyncLogs(c *gin.Context)

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
	EnableConfig(c *gin.Context)
	DisableConfig(c *gin.Context)
	EnableAutoInvoice(c *gin.Context)
	DisableAutoInvoice(c *gin.Context)

	// Estadísticas y resúmenes
	GetSummary(c *gin.Context)
	GetStats(c *gin.Context)
	GetTrends(c *gin.Context)

	// Creación masiva de facturas
	ListInvoiceableOrders(c *gin.Context)
	BulkCreateInvoices(c *gin.Context)

	// Jobs de facturación masiva
	ListBulkJobs(c *gin.Context)
	GetBulkJobStatus(c *gin.Context)

	// Comparación con proveedor (auditoría esporádica)
	CompareInvoices(c *gin.Context)
}

// handler implementa IHandler
type handler struct {
	useCase ports.IUseCase
	repo    ports.IRepository
	config  env.IConfig
	log     log.ILogger
}

// New crea un nuevo handler de facturación
func New(useCase ports.IUseCase, repo ports.IRepository, logger log.ILogger, config env.IConfig) IHandler {
	return &handler{
		useCase: useCase,
		repo:    repo,
		config:  config,
		log:     logger.WithModule("invoicing.handler"),
	}
}

// getS3Config retorna la URL base y bucket de S3 desde la configuración
func (h *handler) getS3Config() (string, string) {
	baseURL := h.config.Get("URL_BASE_DOMAIN_S3")
	bucket := h.config.Get("S3_BUCKET")
	return baseURL, bucket
}
