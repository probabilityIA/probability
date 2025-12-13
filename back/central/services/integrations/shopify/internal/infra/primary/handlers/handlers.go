package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases"
)

type ShopifyHandlers struct {
	syncUseCase *usecases.SyncOrdersUseCase
}

func New(syncUseCase *usecases.SyncOrdersUseCase) *ShopifyHandlers {
	return &ShopifyHandlers{
		syncUseCase: syncUseCase,
	}
}

func (h *ShopifyHandlers) SyncOrders(c *gin.Context) {
	integrationID := c.Param("id")
	fmt.Printf("[SyncOrders Handler] Received request for ID: '%s'\n", integrationID)

	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "integration_id is required",
		})
		return
	}

	// Start sync in background
	if err := h.syncUseCase.SyncOrders(c.Request.Context(), integrationID); err != nil {
		fmt.Printf("[SyncOrders Handler] UseCase Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to start sync",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Order synchronization started in background. Orders from the last 15 days will be imported.",
	})
}
