package businesstypehandler

import (
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/app/usecasebusinesstype"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// IBusinessTypeHandler define la interfaz para el handler de BusinessType
type IBusinessTypeHandler interface {
	GetBusinessTypesHandler(c *gin.Context)
	GetBusinessTypeByIDHandler(c *gin.Context)
	CreateBusinessTypeHandler(c *gin.Context)
	UpdateBusinessTypeHandler(c *gin.Context)
	DeleteBusinessTypeHandler(c *gin.Context)
}

type BusinessTypeHandler struct {
	usecase usecasebusinesstype.IUseCaseBusinessType
	logger  log.ILogger
}

// New crea una nueva instancia del handler de BusinessType
func New(usecase usecasebusinesstype.IUseCaseBusinessType, logger log.ILogger) IBusinessTypeHandler {
	return &BusinessTypeHandler{
		usecase: usecase,
		logger:  logger,
	}
}
