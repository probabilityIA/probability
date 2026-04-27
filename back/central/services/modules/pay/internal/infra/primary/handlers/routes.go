package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	pay := router.Group("/pay")
	pay.Use(middleware.JWT())
	{
		pay.POST("/transactions", h.CreatePayment)
		pay.GET("/transactions", h.ListPayments)
		pay.GET("/transactions/:id", h.GetPayment)
	}
}
