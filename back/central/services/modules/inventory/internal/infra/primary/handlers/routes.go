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

		putawayRules := inventory.Group("/putaway-rules")
		{
			putawayRules.GET("", h.ListPutawayRules)
			putawayRules.POST("", h.CreatePutawayRule)
			putawayRules.PUT("/:id", h.UpdatePutawayRule)
			putawayRules.DELETE("/:id", h.DeletePutawayRule)
		}

		putaway := inventory.Group("/putaway")
		{
			putaway.POST("/suggest", h.SuggestPutaway)
			putaway.POST("/suggestions/:id/confirm", h.ConfirmPutaway)
			putaway.GET("/suggestions", h.ListPutawaySuggestions)
		}

		replenishment := inventory.Group("/replenishment")
		{
			replenishment.GET("/tasks", h.ListReplenishmentTasks)
			replenishment.POST("/tasks", h.CreateReplenishmentTask)
			replenishment.POST("/tasks/:id/assign", h.AssignReplenishment)
			replenishment.POST("/tasks/:id/complete", h.CompleteReplenishment)
			replenishment.POST("/tasks/:id/cancel", h.CancelReplenishment)
			replenishment.POST("/detect", h.DetectReplenishment)
		}

		crossDock := inventory.Group("/cross-dock")
		{
			crossDock.GET("/links", h.ListCrossDockLinks)
			crossDock.POST("/links", h.CreateCrossDockLink)
			crossDock.POST("/links/:id/execute", h.ExecuteCrossDock)
		}

		slotting := inventory.Group("/slotting")
		{
			slotting.POST("/run", h.RunSlotting)
			slotting.GET("/velocities", h.ListVelocities)
		}
	}
}
