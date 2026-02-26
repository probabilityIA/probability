package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define la interfaz pública del handler
type IHandler interface {
	// Métodos HTTP
	Create(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Toggle(c *gin.Context)
	ListOrderStatuses(c *gin.Context)
	ListOrderStatusesSimple(c *gin.Context)
	ListFulfillmentStatuses(c *gin.Context)

	// Método de registro de rutas
	RegisterRoutes(router *gin.RouterGroup)
}

// handler implementa IHandler
type handler struct {
	uc  app.IUseCase
	log log.ILogger
	env env.IConfig
}

// New crea una nueva instancia del handler
func New(useCase app.IUseCase, logger log.ILogger, environment env.IConfig) IHandler {
	return &handler{
		uc:  useCase,
		log: logger.WithModule("orderstatus-handler"),
		env: environment,
	}
}

// getImageURLBase obtiene la URL base de S3 para construir URLs completas
func (h *handler) getImageURLBase() string {
	return h.env.Get("URL_BASE_DOMAIN_S3")
}
