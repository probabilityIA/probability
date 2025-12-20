package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo dashboard
func (h *DashboardHandlers) RegisterRoutes(router *gin.RouterGroup) {
	dashboard := router.Group("/dashboard")
	{
		dashboard.GET("/stats", h.GetStats) // GET /api/v1/dashboard/stats
	}
}
