package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	ListLast(c *gin.Context)
	ListDetail(c *gin.Context)
	Record(c *gin.Context)
	Findings(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type handler struct {
	useCase domain.IUseCase
	logger  log.ILogger
}

func New(useCase domain.IUseCase, logger log.ILogger) IHandler {
	return &handler{
		useCase: useCase,
		logger:  logger.WithModule("syncruns"),
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/integrations/sync-runs", middleware.JWT())
	{
		group.GET("", h.ListLast)
		group.GET("/items", h.ListDetail)
		group.POST("", h.Record)
		group.GET("/findings", h.Findings)
	}
}

func (h *handler) resolveBusinessID(c *gin.Context, provided *uint) (uint, bool) {
	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "contexto de negocio no encontrado"})
		return 0, false
	}
	if businessID == 0 {
		if provided == nil || *provided == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
			return 0, false
		}
		businessID = *provided
	}
	return businessID, true
}
