package whatsapp_template

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/templates"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	RegisterRoutes(router *gin.RouterGroup)

	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Resubmit(c *gin.Context)
	Variables(c *gin.Context)
}

type handler struct {
	useCase templates.IUseCase
	logger  log.ILogger
}

func New(useCase templates.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		logger:  logger.WithModule("whatsapp_template_handler"),
	}
}

func (h *handler) resolveBusinessID(c *gin.Context) (uint, bool) {
	businessID := c.GetUint("business_id")
	if businessID > 0 {
		return businessID, true
	}
	if param := c.Query("business_id"); param != "" {
		if id, err := strconv.ParseUint(param, 10, 64); err == nil && id > 0 {
			return uint(id), true
		}
	}
	return 0, false
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return uint(id), true
}
