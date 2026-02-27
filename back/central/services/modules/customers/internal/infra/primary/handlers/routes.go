package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo customers
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	customers := router.Group("/customers")
	{
		customers.GET("", h.ListClients)         // GET    /api/v1/customers
		customers.GET("/:id", h.GetClient)       // GET    /api/v1/customers/:id
		customers.POST("", h.CreateClient)       // POST   /api/v1/customers
		customers.PUT("/:id", h.UpdateClient)    // PUT    /api/v1/customers/:id
		customers.DELETE("/:id", h.DeleteClient) // DELETE /api/v1/customers/:id
	}
}
