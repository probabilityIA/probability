package campaign

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/campaigns"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	RegisterRoutes(router *gin.RouterGroup)

	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	PreviewAudience(c *gin.Context)
	Launch(c *gin.Context)
	Pause(c *gin.Context)
	Resume(c *gin.Context)
	Cancel(c *gin.Context)
	ListSends(c *gin.Context)
}

type handler struct {
	useCase campaigns.IUseCase
	logger  log.ILogger
}

func New(useCase campaigns.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		logger:  logger.WithModule("campaign_handler"),
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

func paging(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	return page, pageSize
}

func totalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}
