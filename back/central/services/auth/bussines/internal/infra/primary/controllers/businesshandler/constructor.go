package businesshandler

import (
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/app/usecasebusiness"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// IBusinessHandler define la interfaz para el handler de Business
type IBusinessHandler interface {
	GetBusinesses(c *gin.Context)
	GetBusinessesSimple(c *gin.Context)
	GetBusinessByIDHandler(c *gin.Context)
	CreateBusinessHandler(c *gin.Context)
	UpdateBusinessHandler(c *gin.Context)
	DeleteBusinessHandler(c *gin.Context)
	GetBusinessesConfiguredResourcesHandler(c *gin.Context)
	GetBusinessConfiguredResourcesByIDHandler(c *gin.Context)
	ActivateBusinessResourceHandler(c *gin.Context)
	DeactivateBusinessResourceHandler(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, handler IBusinessHandler)
}

type BusinessHandler struct {
	usecase usecasebusiness.IUseCaseBusiness
	logger  log.ILogger
}

// New crea una nueva instancia del handler de Business
func New(usecase usecasebusiness.IUseCaseBusiness, logger log.ILogger) IBusinessHandler {
	return &BusinessHandler{
		usecase: usecase,
		logger:  logger,
	}
}
