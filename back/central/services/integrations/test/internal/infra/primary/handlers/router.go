package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra todas las rutas del módulo test
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	test := router.Group("/test")
	{
		test.POST("/generate-orders", h.GenerateOrders)
	}
}
