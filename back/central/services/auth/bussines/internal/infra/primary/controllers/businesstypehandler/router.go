package businesstypehandler

import (
	"github.com/secamc93/probability/back/central/services/auth/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas del handler de BusinessType
func RegisterRoutes(router *gin.RouterGroup, handler IBusinessTypeHandler) {
	businessTypes := router.Group("/business-types")

	// Rutas de BusinessType
	businessTypes.GET("", middleware.JWT(), handler.GetBusinessTypesHandler)
	businessTypes.GET("/:id", middleware.JWT(), handler.GetBusinessTypeByIDHandler)
	businessTypes.POST("", middleware.JWT(), handler.CreateBusinessTypeHandler)
	businessTypes.PUT("/:id", middleware.JWT(), handler.UpdateBusinessTypeHandler)
	businessTypes.DELETE("/:id", middleware.JWT(), handler.DeleteBusinessTypeHandler)
}
