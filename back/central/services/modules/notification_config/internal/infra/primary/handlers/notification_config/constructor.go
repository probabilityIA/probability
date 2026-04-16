package notification_config

import (
	"strconv"

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

	// Sync
	SyncByIntegration(c *gin.Context)
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

// resolveBusinessID obtiene el business_id efectivo.
// Para usuarios normales usa el del JWT.
// Para super admins (business_id=0 en JWT) lee el query param ?business_id=X.
func (h *handler) resolveBusinessID(c *gin.Context) (uint, bool) {
	businessID := c.GetUint("business_id")
	if businessID > 0 {
		return businessID, true
	}
	// Super admin: leer de query param
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}
