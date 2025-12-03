package handlers

import (
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas del handler de Resource
func RegisterRoutes(router *gin.RouterGroup, handler IResourceHandler, logger log.ILogger) {
	resources := router.Group("/resources")

	// Rutas de Resource CRUD
	resources.GET("", middleware.JWT(), handler.GetResourcesHandler)
	resources.GET("/:id", middleware.JWT(), handler.GetResourceByIDHandler)
	resources.POST("", middleware.JWT(), handler.CreateResourceHandler)
	resources.PUT("/:id", middleware.JWT(), handler.UpdateResourceHandler)
	resources.DELETE("/:id", middleware.JWT(), handler.DeleteResourceHandler)
}
