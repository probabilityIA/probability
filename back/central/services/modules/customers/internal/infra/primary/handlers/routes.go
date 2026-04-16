package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	customers := router.Group("/customers")
	{
		customers.GET("", middleware.JWT(), h.ListClients)
		customers.GET("/:id", middleware.JWT(), h.GetClient)
		customers.POST("", middleware.JWT(), h.CreateClient)
		customers.PUT("/:id", middleware.JWT(), h.UpdateClient)
		customers.DELETE("/:id", middleware.JWT(), h.DeleteClient)

		customers.GET("/:id/summary", middleware.JWT(), h.GetCustomerSummary)
		customers.GET("/:id/addresses", middleware.JWT(), h.ListCustomerAddresses)
		customers.GET("/:id/products", middleware.JWT(), h.ListCustomerProducts)
		customers.GET("/:id/order-items", middleware.JWT(), h.ListCustomerOrderItems)
	}
}
