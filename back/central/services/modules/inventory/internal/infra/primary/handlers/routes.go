package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo inventory
func (h *handlers) RegisterRoutes(router *gin.RouterGroup) {
	inventory := router.Group("/inventory")
	{
		inventory.GET("/product/:productId", h.GetProductInventory)
		inventory.GET("/warehouse/:warehouseId", h.ListWarehouseInventory)
		inventory.POST("/adjust", h.AdjustStock)
		inventory.POST("/transfer", h.TransferStock)
		inventory.GET("/movements", h.ListMovements)

		// Movement Types CRUD
		movTypes := inventory.Group("/movement-types")
		{
			movTypes.GET("", h.ListMovementTypes)
			movTypes.POST("", h.CreateMovementType)
			movTypes.PUT("/:id", h.UpdateMovementType)
			movTypes.DELETE("/:id", h.DeleteMovementType)
		}
	}
}
