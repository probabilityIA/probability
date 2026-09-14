package message_audit

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	List(c *gin.Context)
	Stats(c *gin.Context)
	ListConversations(c *gin.Context)
	GetConversationMessages(c *gin.Context)
	MarkConversationRead(c *gin.Context)
}

type handler struct {
	useCase ports.IUseCase
	logger  log.ILogger
}

func New(useCase ports.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		logger:  logger.WithModule("message_audit_handler"),
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
