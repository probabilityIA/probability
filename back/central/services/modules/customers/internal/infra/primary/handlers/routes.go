package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo customers
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	customers := router.Group("/customers")
	{
		customers.GET("", middleware.JWT(), h.ListClients)         // GET    /api/v1/customers
		customers.GET("/:id", middleware.JWT(), h.GetClient)       // GET    /api/v1/customers/:id
		customers.POST("", middleware.JWT(), h.CreateClient)       // POST   /api/v1/customers
		customers.PUT("/:id", middleware.JWT(), h.UpdateClient)    // PUT    /api/v1/customers/:id
		customers.DELETE("/:id", middleware.JWT(), h.DeleteClient) // DELETE /api/v1/customers/:id
	}
}
