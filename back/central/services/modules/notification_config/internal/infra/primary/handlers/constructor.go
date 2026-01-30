package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define todos los métodos HTTP del handler
type IHandler interface {
	// Rutas
	RegisterRoutes(router *gin.RouterGroup)

	// CRUD
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// handler contiene las dependencias compartidas
type handler struct {
	useCase ports.IUseCase
	logger  log.ILogger
}

// New crea una nueva instancia del handler
func New(useCase ports.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		logger:  logger.WithModule("notification_config_handler"),
	}
}
