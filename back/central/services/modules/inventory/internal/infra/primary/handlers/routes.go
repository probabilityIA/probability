package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handlers) RegisterRoutes(router *gin.RouterGroup) {
	inventory := router.Group("/inventory")
	inventory.Use(middleware.JWT())
	{
		inventory.GET("/product/:productId", h.GetProductInventory)
		inventory.GET("/warehouse/:warehouseId", h.ListWarehouseInventory)
		inventory.POST("/adjust", h.AdjustStock)
		inventory.POST("/transfer", h.TransferStock)
		inventory.POST("/bulk-load", h.BulkLoadInventory)
		inventory.GET("/movements", h.ListMovements)
		inventory.POST("/positions/validate-cubing", h.ValidateCubing)

		lots := inventory.Group("/lots")
		{
			lots.GET("", h.ListLots)
			lots.POST("", h.CreateLot)
			lots.GET("/:id", h.GetLot)
			lots.PUT("/:id", h.UpdateLot)
			lots.DELETE("/:id", h.DeleteLot)
		}

		serials := inventory.Group("/serials")
		{
			serials.GET("", h.ListSerials)
			serials.POST("", h.CreateSerial)
			serials.GET("/:id", h.GetSerial)
			serials.PUT("/:id", h.UpdateSerial)
		}

		inventory.GET("/states", h.ListInventoryStates)
		inventory.POST("/state-transitions", h.ChangeInventoryState)

		inventory.GET("/uoms", h.ListUoMs)
		inventory.GET("/products/:productId/uoms", h.ListProductUoMs)
		inventory.POST("/products/:productId/uoms", h.CreateProductUoM)
		inventory.DELETE("/product-uoms/:id", h.DeleteProductUoM)
		inventory.POST("/uoms/convert", h.ConvertUoM)

		movTypes := inventory.Group("/movement-types")
		{
			movTypes.GET("", h.ListMovementTypes)
			movTypes.POST("", h.CreateMovementType)
			movTypes.PUT("/:id", h.UpdateMovementType)
			movTypes.DELETE("/:id", h.DeleteMovementType)
		}
	}
}
