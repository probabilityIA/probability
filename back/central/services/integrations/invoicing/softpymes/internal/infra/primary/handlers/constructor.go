package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define la interfaz de los handlers HTTP
type IHandler interface {
	// Providers
	CreateProvider(c *gin.Context)
	GetProvider(c *gin.Context)
	ListProviders(c *gin.Context)
	UpdateProvider(c *gin.Context)
	DeleteProvider(c *gin.Context)
	TestProvider(c *gin.Context)

	// Provider Types
	ListProviderTypes(c *gin.Context)

	// Routes registration (OBLIGATORIO)
	RegisterRoutes(router *gin.RouterGroup)
}

// handler implementa IHandler
type handler struct {
	useCase app.IUseCase
	log     log.ILogger
}

// New crea una nueva instancia del handler
func New(useCase app.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		log:     logger.WithModule("softpymes.handlers"),
	}
}
