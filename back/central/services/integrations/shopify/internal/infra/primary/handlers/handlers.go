package handlers

import (
	"net/http"
	"time"

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
	// Get Business ID from context (assuming middleware sets it)
	// businessID := c.GetUint("business_id")

	// Let's assume we sync the default shopify integration for the business.
	businessID := uint(1) // Placeholder

	// Parse optional created_at_min
	var createdMin *time.Time
	if dateStr := c.Query("since"); dateStr != "" {
		t, err := time.Parse(time.RFC3339, dateStr)
		if err == nil {
			createdMin = &t
		}
	}

	// Trigger sync (async ideally, but sync for now)
	err := h.syncUseCase.Execute(c.Request.Context(), &businessID, createdMin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync started and orders published to queue"})
}
