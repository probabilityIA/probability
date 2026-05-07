package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/geozones")
	g.Use(middleware.JWT())
	{
		g.GET("", h.List)
		g.GET("/display", h.Display)
		g.GET("/lookup", h.Lookup)
		g.GET("/probability", h.Probability)
		g.GET("/order-zone", h.OrderZone)
		g.GET("/probability/by-carrier", h.ProbabilityByCarrier)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
		g.POST("/bulk", h.Bulk)
		g.POST("/display/flush-cache", h.FlushDisplayCache)
		g.DELETE("/:id", h.Delete)
	}
}
