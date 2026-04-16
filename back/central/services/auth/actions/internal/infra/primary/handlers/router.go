package handlers

import (
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas del handler de Action
func (h *ActionHandler) RegisterRoutes(router *gin.RouterGroup, handler IActionHandler, logger log.ILogger) {
	actions := router.Group("/actions")

	// Rutas de Action CRUD
	actions.GET("", middleware.JWT(), handler.GetActionsHandler)
	actions.GET("/:id", middleware.JWT(), handler.GetActionByIDHandler)
	actions.POST("", middleware.JWT(), handler.CreateActionHandler)
	actions.PUT("/:id", middleware.JWT(), handler.UpdateActionHandler)
	actions.DELETE("/:id", middleware.JWT(), handler.DeleteActionHandler)
}
