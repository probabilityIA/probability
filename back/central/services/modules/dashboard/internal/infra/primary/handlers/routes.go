package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo dashboard
func (h *DashboardHandlers) RegisterRoutes(router *gin.RouterGroup) {
	dashboard := router.Group("/dashboard")
	dashboard.Use(middleware.JWT()) // Apply Auth Middleware
	{
		dashboard.GET("/stats", h.GetStats)                         // GET /api/v1/dashboard/stats
		dashboard.GET("/top-selling-days", h.GetTopSellingDays)     // GET /api/v1/dashboard/top-selling-days
	}
}
